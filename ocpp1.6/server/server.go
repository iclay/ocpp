package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"net/http"
	"ocpp16/config"
	local "ocpp16/plugin/passive/local"
	"ocpp16/protocol"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

var serverSubprotocols = []string{"ocpp1.6"}

type ActionPlugin interface {
	RequestHandler(action string) (protocol.RequestHandler, bool)
	ResponseHandler(action string) (protocol.ResponseHandler, bool)
}

type Server struct {
	ginServer         *gin.Engine
	upgrader          websocket.Upgrader
	wsconns           *wsconns
	validate          *validator.Validate
	ocpp16map         *protocol.OCPP16Map
	dispatcher        *dispatcher
	loadBalancer      LoadBalancer
	once              sync.Once
	cond              *sync.Cond
	wg                sync.WaitGroup
	actionPlugin      ActionPlugin
	connectHandler    []func(ws *Wsconn) error
	disconnectHandler []func(ws *Wsconn) error
}

func (s *Server) clientOnConnect(ws *Wsconn) {
	s.dispatcher.callStateMap.createNewRequest(ws.id)
	s.dispatcher.requestQueueMap.createNewQueue(ws.id)
	s.registerConn(ws.id, ws.fd, ws)
	if s.connectHandler != nil {
		for _, handler := range s.connectHandler {
			go func() {
				if err := handler(ws); err != nil {
					log.Errorf("id(%s),error(%v)", ws.id, err)
				}
			}()
			time.Sleep(time.Millisecond)
		}
	}
}

func (s *Server) clientOnDisconnect(ws *Wsconn) {
	s.deleteConn(ws.id, ws.fd)
	s.deleteDispatcherQueue(ws.id)
	s.deleteDispatcherCallState(ws.id)
	s.cancelContex(ws.id)
	if s.disconnectHandler != nil {
		for _, handler := range s.disconnectHandler {
			go func() {
				if err := handler(ws); err != nil {
					log.Errorf("id(%s),error(%v)", ws.id, err)
				}
			}()
			time.Sleep(time.Millisecond)
		}
	}
}

func (s *Server) SetConnectHandlers(fns ...func(ws *Wsconn) error) {
	s.connectHandler = append(s.connectHandler, fns...)
}

func (s *Server) SetDisconnetHandlers(fns ...func(ws *Wsconn) error) {
	s.disconnectHandler = append(s.disconnectHandler, fns...)
}

func (s *Server) cancelContex(id string) {
	s.dispatcher.cancelContext(id)
}

func (s *Server) registerConn(id string, fd int, wsconn *Wsconn) {
	s.wsconns.registerConn(id, fd, wsconn)
}

func (s *Server) deleteConn(id string, fd int) {
	s.wsconns.deleteConn(id, fd)
}

func (s *Server) getConn(id string) (*Wsconn, bool) {
	return s.wsconns.getConn(id)
}

func (s *Server) connExists(id string) bool {
	return s.wsconns.connExists(id)
}

func (s *Server) setDefaultDispatcher(d *dispatcher) {
	s.dispatcher = d
}

func (s *Server) requestDone(id string, uniqueid string) {
	s.dispatcher.requestDone(id, uniqueid)
}

func (s *Server) deleteDispatcherCallState(id string) {
	s.dispatcher.callStateMap.deleteRequest(id)
}

func (s *Server) getPendingRequest(uniqueid string) (*request, bool) {
	return s.dispatcher.callStateMap.getPendingRequest(uniqueid)
}

func (s *Server) deleteDispatcherQueue(id string) {
	s.dispatcher.requestQueueMap.deleteQueue(id)
}

func (s *Server) HandleActiveCall(ctx context.Context, id string, call *protocol.Call) error {
	return s.dispatcher.appendRequest(ctx, id, call)
}
func (s *Server) RegisterActiveCallHandler(handler ActiveCallHandler, fn func(ActiveCallHandler)) {
	fn(handler)
}

func (s *Server) RegisterActionPlugin(actionPlugin ActionPlugin) {
	s.actionPlugin = actionPlugin
}

func (s *Server) waitStopSignal() {
	s.cond.L.Lock()
	s.cond.Wait()
	s.cond.L.Unlock()
}

func (s *Server) Stop() {
	conf := config.GCONF
	s.dispatcher.stop(errors.New("stop dispatcher"))
	if conf.UseEpoll {
		s.waitStopSignal()
		//trigger all reactor to stop
		s.loadBalancer.iterate(func(_ int, react *reactor) error {
			if err := react.epoller.trigger(func(_ interface{}) error { return ErrReactorShutdown }, nil); err != nil {
				log.Errorf("failed to call trigger on reactor%d when stopping server", react.index)
			}
			return nil
		})
		s.wg.Wait() //waiting for all reactor to stop
	}
}

var defaultServer = func(useEpoll bool) *Server {
	s := &Server{
		ginServer:    gin.Default(),
		upgrader:     websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
		wsconns:      newWsconns(),
		validate:     protocol.Validate,
		ocpp16map:    protocol.OCPP16M,
		wg:           sync.WaitGroup{},
		actionPlugin: local.NewActionPlugin(), //default action plugin
	}
	pprof.Register(s.ginServer)
	s.setDefaultDispatcher(NewDefaultDispatcher(s))
	s.initOCPPTypePools(s.ocpp16map.SupportActions())
	if useEpoll {
		var epoller *epoller
		var err error
		s.cond = sync.NewCond(&sync.Mutex{})
		loadBalancer := new(roundRobinLoadBalancer)
		defer func() {
			if err != nil {
				loadBalancer.iterate(func(_ int, react *reactor) error {
					react.epoller.close()
					react.connections = nil
					react.index = -1
					return nil
				})
			}
		}()
		numCPU := runtime.NumCPU()
		for index := 0; index < numCPU; index++ {
			if epoller, err = createEpoller(); err != nil {
				panic(err)
			} else {
				loadBalancer.register(&reactor{
					server:      s,
					connections: newWsconns(),
					epoller:     epoller,
				})
			}
		}
		loadBalancer.iterate(func(_ int, r *reactor) error {
			s.wg.Add(1)
			go func() {
				r.activity()
				s.wg.Done()
			}()
			return nil
		})
		s.loadBalancer = loadBalancer
	}
	return s
}

func NewDefaultServer() *Server {
	useEpoll := config.GCONF.UseEpoll
	return defaultServer(useEpoll)
}

func (s *Server) initOCPPTypePools(actions []string) {
	for _, action := range actions {
		if ocpptrait, ok := s.ocpp16map.GetTraitAction(action); ok {
			reqTyp, resTyp := ocpptrait.RequestType(), ocpptrait.ResponseType()
			options.object.init(reqTyp)
			options.object.init(resTyp)
		}
	}
}

type Point struct {
	Name string `uri:"name" binding:"required,uuid"`
	ID   string `uri:"id" binding:"required"`
}

func (c *Point) String() string {
	return fmt.Sprintf("%s-%s", c.Name, c.ID)
}
func (s *Server) Serve(addr string, path string) {
	s.ginServer.GET(path, s.wsHandler)
	s.ginServer.Run(addr)
}

func (s *Server) ServeTLS(addr string, path string, tlsCertificate string, tlsCertificateKey string) {
	s.ginServer.GET(path, s.wsHandler)
	s.ginServer.RunTLS(addr, tlsCertificate, tlsCertificateKey)
}

func (s *Server) wsHandler(c *gin.Context) {
	conf := config.GCONF
	var p Point
	c.ShouldBindUri(&p)
	var ocppProto string
	clientSubprotocols := websocket.Subprotocols(c.Request)
	for _, cproto := range clientSubprotocols {
		for _, sproto := range serverSubprotocols {
			if strings.EqualFold(cproto, sproto) {
				ocppProto = cproto
				break
			}
		}
	}
	respHeader := http.Header{}
	if ocppProto != "" {
		respHeader.Add("Sec-WebSocket-Protocol", ocppProto)
	}
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, respHeader)
	if err != nil {
		log.Error("id(%s) upgrade error", p.String())
		return
	}
	fd := websocketFD2(conn)
	timeoutDuration := time.Second * time.Duration(conf.HeartbeatTimeout)
	//The protocol does not support
	if ocppProto == "" {
		log.Errorf("not support protocol(%+v) current, id(%s)", clientSubprotocols, p.String())
		conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseProtocolError,
			fmt.Sprintf("not support protocol(%+v) current, id(%s)", clientSubprotocols, p.String())), time.Now().Add(timeoutDuration))
		conn.Close()
		return
	}
	//The situation may occur when the charging pile has been disconnected, but the cloud heartbeat mechanism has not responded.
	//When the charging pile initiates a connection, it needs to wait for the cloud to trigger the heartbeat mechanism to close the last connection
	if s.connExists(p.String()) {
		log.Errorf("id(%s) already connect, please wait about %ds and try again", p.String(), conf.HeartbeatTimeout)
		conn.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseProtocolError,
				fmt.Sprintf("id(%s) already connect, wait about %ds and try again", p.String(), conf.HeartbeatTimeout)), time.Now().Add(timeoutDuration))
		conn.Close()
		return
	}
	// _ = unix.SetNonblock(fd, true)
	ws := &Wsconn{
		server:  s,
		conn:    conn,
		id:      p.String(),
		fd:      fd,
		timeout: timeoutDuration,
		ping:    make(chan []byte),
		closeC:  make(chan error, 1),
	}
	ws.setReadDeadTimeout(ws.timeout)
	ws.conn.SetPingHandler(func(appData string) error {
		ws.ping <- Bytes(appData)
		return ws.setReadDeadTimeout(ws.timeout)
	})
	if conf.UseEpoll {
		reactor := s.loadBalancer.next()
		log.Debug("reactor index", reactor.index)
		if err = reactor.epoller.trigger(reactor.registerConn, ws); err != nil {
			log.Errorf("id=%s, error=%v", ws.id, err)
			return
		}
	} else {
		s.clientOnConnect(ws)
		go ws.readdump()
		go ws.writedump()
	}
}

func websocketFD2(conn *websocket.Conn) int {
	connVal := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn").Elem()
	tcpConn := reflect.Indirect(connVal).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}
