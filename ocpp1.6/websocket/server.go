package websocket

import (
	"fmt"
	"net/http"
	"ocpp16/proto"
	"time"

	"reflect"

	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
)

type HandleFuncs interface {
	RegisterOCPPHandler() map[string]proto.RequestHandler
}

type Server struct {
	ginServer      *gin.Engine
	upgrader       websocket.Upgrader
	wsconns        *wsconns
	validate       *validator.Validate
	ocpp16map      *proto.OCPP16Map
	ocppTypePools  *ocppTypePools
	dispatcher     *dispatcher
	ocppHandlerMap map[string]proto.RequestHandler
}

func (s *Server) clientOnConnect(id string, ws *wsconn) {
	s.dispatcher.callStateMap.createNewRequest(id)
	s.dispatcher.requestQueueMap.createNewQueue(id)
	s.registerConn(id, ws)
}
func (s *Server) clientOnDisconnect(id string) {
	s.deleteConn(id)
	s.deleteDispatcherQueue(id)
	s.deleteDispatcherCallState(id)
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

func (s *Server) get(t reflect.Type) interface{} {
	return s.ocppTypePools.get(t)
}
func (s *Server) put(t reflect.Type, x interface{}) {
	s.ocppTypePools.put(t, x)
}

func (s *Server) RegisterOCPPHandler(ocppHandlers HandleFuncs) {
	s.ocppHandlerMap = ocppHandlers.RegisterOCPPHandler()
}

var defaultServer = func() *Server {
	s := &Server{
		ginServer: gin.Default(),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		wsconns:        newWsconns(),
		validate:       proto.Validate,
		ocpp16map:      proto.OCPP16M,
		ocppTypePools:  typePools,
		ocppHandlerMap: make(map[string]proto.RequestHandler),
	}
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
			reqTyp := ocpptrait.RequestType()
			resTyp := ocpptrait.ResponseType()
			s.ocppTypePools.init(reqTyp)
			s.ocppTypePools.init(resTyp)
		}
	}
}

func (s *Server) setDefaultDispatcher(d *dispatcher) {
	s.dispatcher = d
}

type ChargerPoint struct {
	Name string `uri:"name" binding:"required,uuid"`
	ID   string `uri:"id" binding:"required"`
}

func (c *ChargerPoint) String() string {
	return fmt.Sprintf("%s-%s", c.Name, c.ID)
}
func (s *Server) Serve(addr string, path string) {
	s.ginServer.GET(path, s.wsHandler)
	s.ginServer.Run(addr)
}

func (s *Server) wsHandler(c *gin.Context) {
	var p ChargerPoint
	c.ShouldBindUri(&p)
	var ocppProto string
	clientSubprotocols := websocket.Subprotocols(c.Request)
	for _, cproti := range clientSubprotocols {
		for _, sproto := range clientSubprotocols /*need modify server*/ {
			if cproti == sproto {
				ocppProto = sproto
				break
			}
		}
	}
	respHeader := http.Header{}
	respHeader.Add("Sec-WebSocket-Protocol", ocppProto)
	conn, err := s.upgrader.Upgrade(c.Writer, c.Request, respHeader)
	if err != nil {
		return
	}
	// if ocppProto == "" { //The protocol does not support closing connections
	// 	conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseProtocolError,
	// 		fmt.Sprintf("not support protocol for chargegun(%v), protocol(%+v)", p.String(), clientSubprotocols)), time.Now().Add(time.Second) /*时间需要写到配置参数中*/)
	// 	conn.Close()
	// 	return
	// }
	//The situation may occur when the charging pile has been disconnected, but the cloud heartbeat mechanism has not responded.
	//When the charging pile initiates a connection, it needs to wait for the cloud to trigger the heartbeat mechanism to close the last connection
	if s.connExists(p.String()) {
		conn.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseProtocolError,
				fmt.Sprintf("id(%v) already connect, wait a while and try again", p.String())), time.Now().Add(time.Second) /**时间需要写到配置参数中*/)
		conn.Close()
		return
	}
	ws := &wsconn{
		server:  s,
		conn:    conn,
		id:      p.String(),
		timeOut: time.Second * 20,
		ping:    make(chan []byte),
		closeC:  make(chan error),
	}
	s.clientOnConnect(ws.id, ws)
	go ws.read()
	go ws.write()
}
