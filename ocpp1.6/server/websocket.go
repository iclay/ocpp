package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"ocpp16/protocol"
	"reflect"
	"sync"
	"time"
)

type Wsconn struct {
	server  *Server
	conn    *websocket.Conn
	id      string
	fd      int
	timeout time.Duration
	ping    chan []byte
	closeC  chan error
	closed  bool
	sync.Mutex
}

type ErrorDetails struct{}

type wsconns struct {
	wsmap   map[string]*Wsconn
	wsfdmap map[int]*Wsconn
	sync.RWMutex
}

func newWsconns() *wsconns {
	return &wsconns{
		wsmap:   make(map[string]*Wsconn),
		wsfdmap: make(map[int]*Wsconn),
	}
}

func (ws *wsconns) deleteConn(id string, fd int) {
	ws.Lock()
	defer ws.Unlock()
	delete(ws.wsmap, id)
	delete(ws.wsfdmap, fd)
}

func (ws *wsconns) registerConn(id string, fd int, wsconn *Wsconn) {
	ws.Lock()
	defer ws.Unlock()
	ws.wsmap[id] = wsconn
	ws.wsfdmap[fd] = wsconn
}
func (ws *wsconns) getConn(id string) (*Wsconn, bool) {
	ws.RLock()
	defer ws.RUnlock()
	conn, ok := ws.wsmap[id]
	return conn, ok
}

func (ws *wsconns) getConnByFD(fd int) (*Wsconn, bool) {
	ws.RLock()
	defer ws.RUnlock()
	conn, ok := ws.wsfdmap[fd]
	return conn, ok
}

func (ws *wsconns) connExists(id string) bool {
	ws.RLock()
	defer ws.RUnlock()
	_, ok := ws.wsmap[id]
	return ok
}

func (ws *Wsconn) ID() string {
	return ws.id
}

func (ws *Wsconn) stop(err error) {
	if !ws.closed {
		ws.server.clientOnDisconnect(ws)
		ws.closeC <- err
		ws.closed = true
		ws.conn.Close()
	}
}

func (ws *Wsconn) readdump() {
	for {
		if err := ws.read(); err != nil {
			ws.Lock()
			ws.stop(err)
			ws.Unlock()
			break
		}
	}
}

func (ws *Wsconn) read() error {
	typ, message, err := ws.conn.ReadMessage()
	if err != nil {
		log.Errorf("read error, id(%s), err(%v)", ws.id, err)
		return err
	}
	log.Debugf("read: id(%s), recv(%s), messagetype(%d)", ws.id, String(message), typ)
	ws.setReadDeadTimeout(ws.timeout)
	go ws.messageHandler(message)
	return nil
}

func (ws *Wsconn) responseHandler(uniqueid string, action string, res protocol.Response) {
	if handler, ok := ws.server.actionPlugin.ResponseHandler(action); ok {
		log.Debugf("client response, id(%s), uniqueid(%s),action(%s), response(%+v)", ws.id, uniqueid, action, res)
		if err := handler(context.Background(), ws.id, uniqueid, res); err != nil {
			log.Errorf("client response handler failed, id:(%s), uniqueid:(%s),action:(%s),err:(%v)", ws.id, uniqueid, action, err)
		}
		return
	}
	log.Errorf("not support action:(%s) current, id:(%s), uniqueid:(%s)", action, ws.id, uniqueid)
}

func (ws *Wsconn) requestHandler(uniqueid string, action string, req protocol.Request) {
	if handler, ok := ws.server.actionPlugin.RequestHandler(action); ok {
		log.Debugf("client request, id(%s), uniqueid(%s),action(%s), request(%+v)", ws.id, uniqueid, action, req)
		res, err := handler(context.Background(), ws.id, uniqueid, req)
		if err != nil {
			log.Errorf("client request handler failed, id:(%s), uniqueid:(%s),action:(%s),err:(%v)", ws.id, uniqueid, action, err)
			return
		}
		callResult := &protocol.CallResult{
			MessageTypeID: protocol.CALL_RESULT,
			UniqueID:      uniqueid,
			Response:      res,
		}
		log.Debugf("server response, id(%s), uniqueid(%s),action(%s), callResult(%+v)", ws.id, uniqueid, action, callResult)
		if err = ws.server.validate.Struct(callResult); err != nil {
			log.Errorf("validate callResult invalid, id:(%s), uniqueid:(%s),action:(%s),err:(%v)", ws.id, uniqueid, action, checkValidatorError(err, action))
			return
		}
		result, err := json.Marshal(callResult)
		if err != nil {
			log.Errorf("marshal result error, id:(%s), uniqueid:(%s),action:(%s),err:(%v)", ws.id, uniqueid, action, err)
			return
		}
		if err = ws.writeMessage(websocket.TextMessage, result); err != nil {
			log.Errorf("write message error, id:(%s), uniqueid:(%s),action:(%s),err:(%v)", ws.id, uniqueid, action, err)
		}
		return
	}
	log.Errorf("not support action:(%s) current, id:(%s), uniqueid:(%s)", action, ws.id, uniqueid)
}

func (ws *Wsconn) setReadDeadTimeout(readTimeout time.Duration) error {
	ws.Lock()
	defer ws.Unlock()
	return ws.conn.SetReadDeadline(time.Now().Add(readTimeout))
}

func (ws *Wsconn) setWriteDeadTimeout(readTimeout time.Duration) error {
	//this function is always accompanied by writemessage, so there is no need to lock it
	return ws.conn.SetWriteDeadline(time.Now().Add(readTimeout))
}

func (ws *Wsconn) writeMessage(messageType int, data []byte) error {
	ws.Lock()
	defer ws.Unlock()
	if ws.closed {
		return fmt.Errorf("conn has closed down, id(%s)", ws.id)
	}
	ws.setWriteDeadTimeout(ws.timeout)
	var err error
	if err = ws.conn.WriteMessage(messageType, data); err != nil {
		ws.stop(err)
	}
	return err
}

func (ws *Wsconn) sendCallError(uniqueID string, e *Error) error {
	callError := &protocol.CallError{
		MessageTypeID:    protocol.CALL_ERROR,
		UniqueID:         uniqueID,
		ErrorCode:        e.ErrorCode,
		ErrorDescription: e.ErrorDescription,
		ErrorDetails:     e.ErrorDetails,
	}
	if err := ws.server.validate.Struct(callError); err != nil {
		return err
	}
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.Encode(callError)
	return ws.writeMessage(websocket.TextMessage, bf.Bytes())
}

func parseMessage(wsmsg []byte) ([]interface{}, error) {
	var fields []interface{}
	if err := json.Unmarshal(wsmsg, &fields); err != nil {
		return nil, err
	}
	return fields, nil
}

const (
	Call       = protocol.CallName
	CallResult = protocol.CallResultName
	CallError  = protocol.CallErrorName
)

func (ws *Wsconn) callHandler(uniqueid string, wsmsg []byte, fields []interface{}) {

	if len(fields) != 4 {
		log.Errorf("invalid num of call fields(%+v),exptect 4 fields, id(%s), wsmsg(%s),wsmsg_type(%s)", fields, ws.id, String(wsmsg), Call)
		if err := ws.sendCallError(uniqueid, &Error{
			ErrorCode:        protocol.FormationViolation,
			ErrorDescription: fmt.Sprintf("invalid num of call fields(%+v),exptect 4 fields,uniqueid(%s)", fields, uniqueid),
			ErrorDetails:     ErrorDetails{}}); err != nil {
			log.Errorf("send CallError error(%v), id(%s), wsmsg(%s),wsmsg_type(%s)", err, ws.id, String(wsmsg), Call)
		}
		return
	}
	action, ok := fields[2].(string)
	if !ok {
		log.Errorf("invalid call action(%s) type,must be string, id(%s), wsmsg(%s), wsmsg_type(%s)", fields[2], ws.id, String(wsmsg), Call)
		if err := ws.sendCallError(uniqueid, &Error{
			ErrorCode:        protocol.TypeConstraintViolation,
			ErrorDescription: fmt.Sprintf("invalid call action(%s) type,must be string,uniqueid(%s)", fields[2], uniqueid),
			ErrorDetails:     ErrorDetails{}}); err != nil {
			log.Errorf("send CallError error(%v),id(%s),wsmsg(%s),wsmsg_type(%s)", err, ws.id, String(wsmsg), Call)
		}
		return
	}
	ocpptrait, ok := ws.server.ocpp16map.GetTraitAction(action)
	if !ok {
		log.Errorf("not support action(%s) current,id(%s),wsmsg(%s),wsmsg_type(%s)", action, ws.id, String(wsmsg), Call)
		if err := ws.sendCallError(uniqueid, &Error{
			ErrorCode:        protocol.NotSupported,
			ErrorDescription: fmt.Sprintf("action(%s) not support current,uniqueid(%s)", action, uniqueid),
			ErrorDetails:     ErrorDetails{}}); err != nil {
			log.Errorf("send CallError error(%v),id(%s),wsmsg(%s),wsmsg_type(%s)", err, ws.id, String(wsmsg), Call)
		}
		return
	}
	reqType := ocpptrait.RequestType()
	reqByte, err := json.Marshal(fields[3])
	if err != nil {
		log.Errorf("json Marshal error(%v),id(%s),wsmsg(%s),wsmsg_type(%s)", err, ws.id, String(wsmsg), Call)
		if err = ws.sendCallError(uniqueid, &Error{
			ErrorCode:        protocol.CallInternalError,
			ErrorDescription: fmt.Sprintf("json Marshal error(%v),uniqueid(%s)", err, uniqueid),
			ErrorDetails:     ErrorDetails{}}); err != nil {
			log.Errorf("send CallError error(%v),id(%s),wsmsg(%s),wsmsg_type(%s)", err, ws.id, String(wsmsg), Call)
		}
		return
	}
	req := get(reqType)
	defer put(reqType, req)
	if err = json.Unmarshal(reqByte, &req); err != nil {
		log.Errorf("json Unmarshal error(%v),id(%s),wsmsg(%s),wsmsg_type(%s)", err, ws.id, String(wsmsg), Call)
		if err = ws.sendCallError(uniqueid, &Error{
			ErrorCode:        protocol.CallInternalError,
			ErrorDescription: fmt.Sprintf("json Marshal error(%v),uniqueid(%s)", err, uniqueid),
			ErrorDetails:     ErrorDetails{}}); err != nil {
			log.Errorf("send CallError error(%v),id(%s),wsmsg(%s),wsmsg_type(%s)", err, ws.id, String(wsmsg), Call)
		}
		return
	}
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        action,
		Request:       req.(protocol.Request),
	}
	if err = ws.server.validate.Struct(call); err != nil {
		log.Errorf("validate Call error(%v),id(%s),wsmsg(%s),wsmsg_type(%s)", checkValidatorError(err, action), ws.id, String(wsmsg), Call)
		if err = ws.sendCallError(uniqueid, checkValidatorError(err, action)); err != nil {
			log.Errorf("send CallError error(%v),id(%s),wsmsg(%s),wsmsg_type(%s)", err, ws.id, String(wsmsg), Call)
		}
		return
	}
	ws.requestHandler(uniqueid, action, req.(protocol.Request))
}
func (ws *Wsconn) callResultHandler(uniqueid string, wsmsg []byte, fields []interface{}) {
	if len(fields) != 3 {
		log.Errorf("invalid num of call fields(%+v),exptect 3 fields,id(%s),msg(%s),wsmsg_type(%s)", fields, ws.id, String(wsmsg), CallResult)
		if err := ws.sendCallError(uniqueid, &Error{
			ErrorCode:        protocol.FormationViolation,
			ErrorDescription: fmt.Sprintf("invalid num of callresult fields(%+v),exptect 3 fields,uniqueid(%s)", fields, uniqueid),
			ErrorDetails:     ErrorDetails{}}); err != nil {
			log.Errorf("send CallError error(%v),id(%s),wsmsg(%s),wsmsg_type(%s)", err, ws.id, String(wsmsg), CallResult)
		}
		return
	}
	pendingReq, ok := ws.server.getPendingRequest(ws.id)
	if !ok {
		log.Errorf("ignoring this message, may be conn close,id(%s), wsmsg(%s),wsmsg_type(%s)", ws.id, String(wsmsg), CallResult)
		return
	}
	action := pendingReq.call.Action
	if action == "" {
		log.Errorf("action is nil, may be client response timeout or center never request,id(%s),wsmsg(%s),wsmsg_type(%s)", ws.id, String(wsmsg), CallResult)
		if err := ws.sendCallError(uniqueid, &Error{
			ErrorCode:        protocol.CallInternalError,
			ErrorDescription: fmt.Sprintf("may be client response timeout or center never request,uniqueid(%s)", uniqueid),
			ErrorDetails:     ErrorDetails{}}); err != nil {
			log.Errorf("send CallError error(%v),id(%s),wsmsg(%s),wsmsg_type(%s)", err, ws.id, String(wsmsg), CallResult)
		}
		return
	}
	ocpptrait, ok := ws.server.ocpp16map.GetTraitAction(action)
	if !ok {
		log.Errorf("not support action(%s) current,id(%s),wsmsg(%s),wsmsg_type(%s)", action, ws.id, String(wsmsg), CallResult)
		if err := ws.sendCallError(uniqueid, &Error{
			ErrorCode:        protocol.NotSupported,
			ErrorDescription: fmt.Sprintf("action(%s) not support current,uniqueid(%s)", action, uniqueid),
			ErrorDetails:     ErrorDetails{}}); err != nil {
			log.Errorf("send CallError error(%v),id(%s),wsmsg(%s),wsmsg_type(%s)", err, ws.id, String(wsmsg), CallResult)
		}
		return
	}
	resType := ocpptrait.ResponseType()
	resByte, err := json.Marshal(fields[2])
	if err != nil {
		log.Errorf("json Marshal error(%v),id(%s),wsmsg(%s),wsmsg_type(%s)", err, ws.id, String(wsmsg), CallResult)
		if err = ws.sendCallError(uniqueid, &Error{
			ErrorCode:        protocol.CallInternalError,
			ErrorDescription: fmt.Sprintf("json Marshal error(%v),uniqueid(%s)", err, uniqueid),
			ErrorDetails:     ErrorDetails{}}); err != nil {
			log.Errorf("send CallError error(%v),id(%s),wsmsg(%s),wsmsg_type(%s)", err, ws.id, String(wsmsg), CallResult)
		}
		return
	}
	res := get(resType)
	defer put(resType, res)
	if err = json.Unmarshal(resByte, &res); err != nil {
		log.Errorf("json Unmarshal error(%v),id(%s),wsmsg(%s),wsmsg_type(%s)", err, ws.id, String(wsmsg), CallResult)
		if err = ws.sendCallError(uniqueid, &Error{
			ErrorCode:        protocol.CallInternalError,
			ErrorDescription: fmt.Sprintf("json Marshal error(%v),uniqueid(%s)", err, uniqueid),
			ErrorDetails:     ErrorDetails{}}); err != nil {
			log.Errorf("send CallError error(%v),id(%s),wsmsg(%s),wsmsg_type(%s)", err, ws.id, String(wsmsg), CallResult)
		}
		return
	}
	callResult := protocol.CallResult{
		MessageTypeID: protocol.CALL_RESULT,
		UniqueID:      uniqueid,
		Response:      res.(protocol.Response),
	}
	if err = ws.server.validate.Struct(callResult); err != nil {
		log.Errorf("validate CallResult error(%v),id(%s),wsmsg(%s),wsmsg_type(%s)", checkValidatorError(err, action), ws.id, String(wsmsg), CallResult)
		if err = ws.sendCallError(uniqueid, checkValidatorError(err, action)); err != nil {
			log.Errorf("send CallError error(%v),id(%s),wsmsg(%s),wsmsg_type(%s)", err, ws.id, String(wsmsg), CallResult)
		}
		return
	}
	ws.server.requestDone(ws.id, uniqueid)
	ws.responseHandler(uniqueid, action, res.(protocol.Response))
}

func (ws *Wsconn) callErrorHandler(uniqueid string, wsmsg []byte, fields []interface{}) {
	if len(fields) != 5 {
		log.Errorf("invalid num of call fields(%+v), id(%s), wsmsg(%s),wsg_type(%d),exptect 5 fields", fields, ws.id, String(wsmsg), CallError)
		if err := ws.sendCallError(uniqueid, &Error{
			ErrorCode:        protocol.FormationViolation,
			ErrorDescription: fmt.Sprintf("invalid num of callresult fields(%+v),exptect 5 fields,uniqueid(%s)", fields, uniqueid),
			ErrorDetails:     ErrorDetails{}}); err != nil {
			log.Errorf("send CallError error(%v),id(%s),wsmsg(%s),wsmsg_type(%s)", err, ws.id, String(wsmsg), CallError)
		}
		return
	}
	errCode, ok := fields[2].(string)
	if !ok {
		log.Errorf("invalid CallError ErrCode(%v) type,must be string, id(%s), wsmsg(%s), wsmsg_type(%s)", fields[2], ws.id, String(wsmsg), CallError)
		if err := ws.sendCallError(uniqueid, &Error{
			ErrorCode:        protocol.TypeConstraintViolation,
			ErrorDescription: fmt.Sprintf("invalid CallError errCode(%v) type,must be string,uniqueid(%s)", fields[2], uniqueid),
			ErrorDetails:     ErrorDetails{}}); err != nil {
			log.Errorf("send CallError error(%v),id(%s),wsmsg(%s),wsmsg_type(%s)", err, ws.id, String(wsmsg), CallError)
		}
		return
	}
	errorDescription, ok := fields[2].(string)
	if !ok {
		log.Errorf("invalid CallError errorDescription(%v) type,must be string, id(%s), wsmsg(%s), wsmsg_type(%s)", fields[2], ws.id, String(wsmsg), CallError)
		if err := ws.sendCallError(uniqueid, &Error{
			ErrorCode:        protocol.TypeConstraintViolation,
			ErrorDescription: fmt.Sprintf("invalid CallError errorDescription(%v) type,must be string,uniqieid(%s)", fields[2], uniqueid),
			ErrorDetails:     ErrorDetails{}}); err != nil {
			log.Errorf("send CallError error(%v),id(%s),wsmsg(%s),wsmsg_type(%s)", err, ws.id, String(wsmsg), CallError)
		}
		return
	}
	pendingReq, ok := ws.server.getPendingRequest(ws.id)
	if !ok {
		log.Errorf("ignoring this message, may be conn close,id(%s), wsmsg(%s),wsmsg_type(%s)", ws.id, String(wsmsg), CallError)
		return
	}
	action := pendingReq.call.Action
	if action == "" {
		log.Errorf("action is nil, may be client response timeout or center never request,id(%s),wsmsg(%s),wsmsg_type(%s)", ws.id, String(wsmsg), CallError)
		if err := ws.sendCallError(uniqueid, &Error{
			ErrorCode:        protocol.CallInternalError,
			ErrorDescription: fmt.Sprintf("may be you client response timeout or center never request,uniqueid(%s)", uniqueid),
			ErrorDetails:     ErrorDetails{}}); err != nil {
			log.Errorf("send CallError error(%v),id(%s),wsmsg(%s),wsmsg_type(%s)", err, ws.id, String(wsmsg), CallError)
		}
		return
	}
	callError := protocol.CallError{
		MessageTypeID:    protocol.CALL_ERROR,
		UniqueID:         uniqueid,
		ErrorCode:        protocol.ErrCodeType(errCode),
		ErrorDescription: errorDescription,
		ErrorDetails:     fields[4],
	}
	if err := ws.server.validate.Struct(callError); err != nil {
		log.Errorf("validate CallError error(%v),id(%s),wsmsg(%s),wsmsg_type(%s)", checkValidatorError(err, action), ws.id, String(wsmsg), CallError)
		action := pendingReq.call.Action
		if err = ws.sendCallError(uniqueid, checkValidatorError(err, action)); err != nil {
			log.Errorf("send CallError error(%v),id(%s),wsmsg(%s),wsmsg_type(%s)", err, ws.id, String(wsmsg), CallError)
		}
		return
	}
	ws.server.requestDone(ws.id, uniqueid)
	ws.responseHandler(uniqueid, protocol.CallErrorName, &callError)
}

func (ws *Wsconn) messageHandler(wsmsg []byte) {
	fields, err := parseMessage(wsmsg)
	if err != nil {
		log.Errorf("parse wsmessage error, id(%s), wsmsg(%s), err(%v)", ws.id, String(wsmsg), err)
		return
	}
	if len(fields) < 3 {
		log.Errorf("invalid wsmessage because of fields < 3, id(%s), msg(%s)", ws.id, String(wsmsg))
		return
	}

	wsmsgTypeid, ok := fields[0].(float64)
	if !ok {
		log.Errorf("invalid wsmsgTypeID(%v), type(%s),must be float64, id(%s), msg(%s)", wsmsgTypeid, reflect.TypeOf(wsmsgTypeid).String(), ws.id, String(wsmsg))
		return
	}
	uniqueid, ok := fields[1].(string)
	if !ok {
		log.Errorf("invalid uniqueid(%s), type(%s),must be string, id(%s),msg(%s)", fields[1], reflect.TypeOf(fields[1]).String(), ws.id, String(wsmsg))
		return
	}
	switch wsmsgTypeid {
	case protocol.CALL:
		ws.callHandler(uniqueid, wsmsg, fields)
	case protocol.CALL_RESULT:
		ws.callResultHandler(uniqueid, wsmsg, fields)
	case protocol.CALL_ERROR:
		ws.callErrorHandler(uniqueid, wsmsg, fields)
	default:
		log.Errorf("not support wsmsgTypeID(%v) current, id(%s)", wsmsgTypeid, ws.id)
	}
	return
}

func (ws *Wsconn) writedump() {
	for {
		select {
		case ping := <-ws.ping:
			log.Debugf("id(%s) recv ping message(%s)", ws.id, String(ping))
			if err := ws.writeMessage(websocket.PongMessage, ping); err != nil {
				log.Errorf("id(%s) write pong message(%s) error", ws.id, String(ping))
			}
		case closeError := <-ws.closeC:
			log.Errorf("id(%s) closed, err(%v)", ws.id, closeError)
			return
		}
	}
}
