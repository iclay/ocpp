package websocket

import (
	"context"
	"fmt"
	"net/http"
	"ocpp16/config"
	local "ocpp16/plugin/passive/local"
	"ocpp16/protocol"
	"reflect"
	"strings"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
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
	ocppTypePools     *ocppTypePools
	dispatcher        *dispatcher
	actionPlugin      ActionPlugin
	connectHandler    []func(id string) error
	disconnectHandler []func(id string) error
}

func logIfError(id string, err error) {
	if err != nil {
		log.Errorf("id(%s),error(%v)", id, err)
	}
}

func (s *Server) clientOnConnect(id string, ws *wsconn) {
	s.dispatcher.callStateMap.createNewRequest(id)
	s.dispatcher.requestQueueMap.createNewQueue(id)
	s.registerConn(id, ws)
	log.Debug(len(s.connectHandler))
	if s.connectHandler != nil {
		for _, handler := range s.connectHandler {
			go logIfError(id, handler(id))
		}
	}
}

func (s *Server) SetConnectHandlers(fns ...func(id string) error) {
	s.connectHandler = append(s.connectHandler, fns...)
}

func (s *Server) clientOnDisconnect(id string) {
	s.deleteConn(id)
	s.deleteDispatcherQueue(id)
	s.deleteDispatcherCallState(id)
	s.cancelContex(id)
	if s.disconnectHandler != nil {
		for _, handler := range s.disconnectHandler {
			go logIfError(id, handler(id))
		}
	}
}

func (s *Server) SetDisconnetHandlers(fns ...func(id string) error) {
	s.disconnectHandler = append(s.disconnectHandler, fns...)
}

func (s *Server) cancelContex(id string) {
	s.dispatcher.cancelContext(id)
}

func (s *Server) registerConn(id string, wsconn *wsconn) {
	s.wsconns.registerConn(id, wsconn)
}

func (s *Server) connExists(id string) bool {
	return s.wsconns.connExists(id)
}

func (s *Server) getConn(id string) (*wsconn, bool) {
	return s.wsconns.getConn(id)
}

func (s *Server) deleteConn(id string) {
	s.wsconns.deleteConn(id)
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

func (s *Server) get(t reflect.Type) interface{} {
	return s.ocppTypePools.get(t)
}

func (s *Server) put(t reflect.Type, x interface{}) {
	s.ocppTypePools.put(t, x)
}

func (s *Server) RegisterActiveCallHandler(handler ActiveCallHandler, fn func(ActiveCallHandler)) {
	fn(handler)
}

func (s *Server) RegisterActionPlugin(actionPlugin ActionPlugin) {
	s.actionPlugin = actionPlugin
}

var defaultServer = func() *Server {
	s := &Server{
		ginServer:     gin.Default(),
		upgrader:      websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
		wsconns:       newWsconns(),
		validate:      protocol.Validate,
		ocpp16map:     protocol.OCPP16M,
		ocppTypePools: typePools,
		actionPlugin:  local.NewActionPlugin(), //default action plugin
	}
	pprof.Register(s.ginServer)
	s.setDefaultDispatcher(NewDefaultDispatcher(s))
	s.initOCPPTypePools(s.ocpp16map.SupportActions())
	return s
}()

func NewDefaultServer() *Server {
	return defaultServer
}

func (s *Server) initOCPPTypePools(actions []string) {
	for _, action := range actions {
		if ocpptrait, ok := s.ocpp16map.GetTraitAction(action); ok {
			reqTyp, resTyp := ocpptrait.RequestType(), ocpptrait.ResponseType()
			s.ocppTypePools.init(reqTyp)
			s.ocppTypePools.init(resTyp)
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
	conn, err := s.upgrader.Upgrade(c.Writer, c.Request, respHeader)
	if err != nil {
		log.Error("id(%s) upgrade error", p.String())
		return
	}
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
		log.Errorf("id(%s) already connect, wait about %ds and try again", p.String(), conf.HeartbeatTimeout)
		conn.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseProtocolError,
				fmt.Sprintf("id(%s) already connect, wait a while and try again", p.String())), time.Now().Add(timeoutDuration))
		conn.Close()
		return
	}
	ws := &wsconn{
		server:  s,
		conn:    conn,
		id:      p.String(),
		timeout: timeoutDuration,
		ping:    make(chan []byte),
		closeC:  make(chan error),
	}
	s.clientOnConnect(ws.id, ws)
	go ws.read()
	go ws.write()
}
