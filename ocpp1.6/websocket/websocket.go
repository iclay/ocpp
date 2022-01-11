package websocket

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"ocpp16/proto"
	"reflect"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/smallnest/rpcx/share"
)

type wsconn struct {
	server  *Server
	conn    *websocket.Conn
	id      string
	timeout time.Duration
	ping    chan []byte
	closeC  chan error
	closed  bool
	sync.Mutex
}

type wsconns struct {
	wsmap map[string]*wsconn
	sync.RWMutex
}

func newWsconns() *wsconns {
	return &wsconns{
		wsmap: make(map[string]*wsconn),
	}
}

func (ws *wsconns) deleteConn(id string) {
	ws.Lock()
	defer ws.Unlock()
	delete(ws.wsmap, id)
}

func (ws *wsconns) registerConn(id string, wsconn *wsconn) {
	ws.Lock()
	defer ws.Unlock()
	ws.wsmap[id] = wsconn
}
func (ws *wsconns) getConn(id string) (*wsconn, bool) {
	ws.RLock()
	defer ws.RUnlock()
	conn, ok := ws.wsmap[id]
	return conn, ok
}
func (ws *wsconns) connExists(id string) bool {
	ws.RLock()
	defer ws.RUnlock()
	_, ok := ws.wsmap[id]
	return ok
}

func (ws *wsconn) stop(err error) {
	ws.server.clientOnDisconnect(ws.id)
	ws.closeC <- err
	ws.closed = true
	ws.conn.Close()
}

func (ws *wsconn) read() {
	conn := ws.conn
	ws.setReadDeadTimeout(ws.timeout)
	conn.SetPingHandler(func(appData string) error {
		ws.ping <- []byte(appData)
		return ws.setReadDeadTimeout(ws.timeout)
	})
	for {
		typ, message, err := conn.ReadMessage()
		if err != nil {
			log.Errorf("read error, id(%v), err(%v)", ws.id, err)
			ws.Lock()
			ws.stop(err)
			ws.Unlock()
			return
		}
		log.Debugf("read: id(%v), recv(%v), messagetype(%v)", ws.id, string(message), typ)
		ws.setReadDeadTimeout(ws.timeout)
		go ws.messageHandler(message)
	}
}

func (ws *wsconn) responseHandler(uniqueid string, action string, res proto.Response) {
	if handler, ok := ws.server.actionPlugin.ResponseHandler(action); ok {
		ctx := context.WithValue(context.Background(), share.ResMetaDataKey, map[string]string{
			"chargingPointIdentify": ws.id,
			"messageId":             uniqueid,
		})
		log.Debugf("client response, id(%v), uniqueid(%v),action(%v), response(%+v)", ws.id, uniqueid, action, res)
		if err := handler(ctx, res); err != nil {
			log.Errorf("response handler failed, id:(%v), uniqueid:(%v),action:(%v),err:(%v)", ws.id, uniqueid, action, err)
		}
		return
	}
	log.Errorf("not support action:(%v) current, id:(%v), uniqueid:(%v)", action, ws.id, uniqueid)
}

func (ws *wsconn) requestHandler(uniqueid string, action string, req proto.Request) {

	if handler, ok := ws.server.actionPlugin.RequestHandler(action); ok {
		ctx := context.WithValue(context.Background(), share.ReqMetaDataKey, map[string]string{
			"chargingPointIdentify": ws.id,
			"messageId":             uniqueid,
		})
		log.Debugf("client request, id(%v), uniqueid(%v),action(%v), request(%+v)", ws.id, uniqueid, action, req)
		res, err := handler(ctx, req)
		if err != nil {
			log.Errorf("request handler failed, id:(%v), uniqueid:(%v),action:(%v),err:(%v)", ws.id, uniqueid, action, err)
			return
		}
		callResult := &proto.CallResult{
			MessageTypeID: proto.CALL_RESULT,
			UniqueID:      uniqueid,
			Response:      res,
		}
		log.Debugf("server response, id(%v), uniqueid(%v),action(%v), callResult(%+v)", ws.id, uniqueid, action, callResult)
		if err = ws.server.validate.Struct(callResult); err != nil {
			log.Errorf("validate callResult invalid, id:(%v), uniqueid:(%v),action:(%v),err:(%v)", ws.id, uniqueid, action, checkValidatorError(err, action))
			return
		}
		result, err := json.Marshal(callResult)
		if err != nil {
			log.Errorf("marshal result error, id:(%v), uniqueid:(%v),action:(%v),err:(%v)", ws.id, uniqueid, action, err)
			return
		}
		if err = ws.writeMessage(websocket.TextMessage, result); err != nil {
			log.Errorf("write message error, id:(%v), uniqueid:(%v),action:(%v),err:(%v)", ws.id, uniqueid, action, err)
		}
		return
	}
	log.Errorf("not support action:(%v) current, id:(%v), uniqueid:(%v)", action, ws.id, uniqueid)
}

func (ws *wsconn) setReadDeadTimeout(readTimeout time.Duration) error {
	ws.Lock()
	defer ws.Unlock()
	return ws.conn.SetReadDeadline(time.Now().Add(readTimeout))
}

func (ws *wsconn) setWriteDeadTimeout(readTimeout time.Duration) error {
	//this function is always accompanied by writemessage, so there is no need to lock it
	return ws.conn.SetWriteDeadline(time.Now().Add(readTimeout))
}

func (ws *wsconn) writeMessage(messageType int, data []byte) error {
	ws.Lock()
	defer ws.Unlock()
	if ws.closed {
		return fmt.Errorf("conn has closed down, id(%v)", ws.id)
	}
	ws.setWriteDeadTimeout(ws.timeout)
	var err error
	if err = ws.conn.WriteMessage(messageType, data); err != nil {
		ws.stop(err)
	}
	return err
}

func (ws *wsconn) sendCallError(uniqueID string, e *Error) error {
	callError := &proto.CallError{
		MessageTypeID:    proto.CALL_ERROR,
		UniqueID:         uniqueID,
		ErrorCode:        e.ErrorCode,
		ErrorDescription: e.ErrorDescription,
		ErrorDetails:     nil,
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
	Call       = proto.CallName
	CallResult = proto.CallResultName
	CallError  = proto.CallErrorName
)

func (ws *wsconn) callHandler(uniqueid string, wsmsg []byte, fields []interface{}) {

	if len(fields) != 4 {
		log.Errorf("invalid num of call fields(%+v),exptect 4 fields, id(%v), wsmsg(%v),wsmsg_type(%v)", fields, ws.id, string(wsmsg), Call)
		if err := ws.sendCallError(uniqueid, &Error{
			ErrorCode:        proto.FormationViolation,
			ErrorDescription: fmt.Sprintf("invalid num of call fields(%+v),exptect 4 fields,uniqueid(%v)", fields, uniqueid),
			ErrorDetails:     nil}); err != nil {
			log.Errorf("send CallError error(%v), id(%v), wsmsg(%v),wsmsg_type(%v)", err, ws.id, string(wsmsg), Call)
		}
		return
	}
	action, ok := fields[2].(string)
	if !ok {
		log.Errorf("invalid call action(%v) type,must be string, id(%v), wsmsg(%v), wsmag_type(%v)", fields[2], ws.id, string(wsmsg), Call)
		if err := ws.sendCallError(uniqueid, &Error{
			ErrorCode:        proto.TypeConstraintViolation,
			ErrorDescription: fmt.Sprintf("invalid call action(%v) type,must be string,uniqueid(%v)", fields[2], uniqueid),
			ErrorDetails:     nil}); err != nil {
			log.Errorf("send CallError error(%v),id(%v),wsmsg(%v),wsmsg_type(%v)", err, ws.id, string(wsmsg), Call)
		}
		return
	}
	ocpptrait, ok := ws.server.ocpp16map.GetTraitAction(action)
	if !ok {
		log.Errorf("not support action(%v) current,id(%v),wsmsg(%v),wsmsg_type(%v)", action, ws.id, string(wsmsg), Call)
		if err := ws.sendCallError(uniqueid, &Error{
			ErrorCode:        proto.NotSupported,
			ErrorDescription: fmt.Sprintf("action(%v) not support current,uniqueid(%v)", action, uniqueid),
			ErrorDetails:     nil}); err != nil {
			log.Errorf("send CallError error(%v),id(%v),wsmsg(%v),wsmsg_type(%v)", err, ws.id, string(wsmsg), Call)
		}
		return
	}
	reqType := ocpptrait.RequestType()
	reqByte, err := json.Marshal(fields[3])
	if err != nil {
		log.Errorf("json Marshal error(%v),id(%v),wsmsg(%v),wsmsg_type(%v)", err, ws.id, string(wsmsg), Call)
		if err = ws.sendCallError(uniqueid, &Error{
			ErrorCode:        proto.CallInternalError,
			ErrorDescription: fmt.Sprintf("json Marshal error(%v),uniqueid(%v)", err, uniqueid),
			ErrorDetails:     nil}); err != nil {
			log.Errorf("send CallError error(%v),id(%v),wsmsg(%v),wsmsg_type(%v)", err, ws.id, string(wsmsg), Call)
		}
		return
	}
	req := ws.server.get(reqType)
	defer ws.server.put(reqType, req)
	if err = json.Unmarshal(reqByte, &req); err != nil {
		log.Errorf("json Unmarshal error(%v),id(%v),wsmsg(%v),wsmsg_type(%v)", err, ws.id, string(wsmsg), Call)
		if err = ws.sendCallError(uniqueid, &Error{
			ErrorCode:        proto.CallInternalError,
			ErrorDescription: fmt.Sprintf("json Marshal error(%v),uniqueid(%v)", err, uniqueid),
			ErrorDetails:     nil}); err != nil {
			log.Errorf("send CallError error(%v),id(%v),wsmsg(%v),wsmsg_type(%v)", err, ws.id, string(wsmsg), Call)
		}
		return
	}
	call := proto.Call{
		MessageTypeID: proto.CALL,
		UniqueID:      uniqueid,
		Action:        action,
		Request:       req.(proto.Request),
	}
	if err = ws.server.validate.Struct(call); err != nil {
		log.Errorf("validate Call error(%v),id(%v),wsmsg(%v),wsmsg_type(%v)", checkValidatorError(err, action), ws.id, string(wsmsg), Call)
		if err = ws.sendCallError(uniqueid, checkValidatorError(err, action)); err != nil {
			log.Errorf("send CallError error(%v),id(%v),wsmsg(%v),wsmsg_type(%v)", err, ws.id, string(wsmsg), Call)
		}
		return
	}
	ws.requestHandler(uniqueid, action, req.(proto.Request))
}
func (ws *wsconn) callResultHandler(uniqueid string, wsmsg []byte, fields []interface{}) {
	if len(fields) != 3 {
		log.Errorf("invalid num of call fields(%+v),exptect 3 fields,id(%v),msg(%v),wsmsg_type(CALL)", fields, ws.id, string(wsmsg))
		if err := ws.sendCallError(uniqueid, &Error{
			ErrorCode:        proto.FormationViolation,
			ErrorDescription: fmt.Sprintf("invalid num of callresult fields(%+v),exptect 3 fields,uniqueid(%v)", fields, uniqueid),
			ErrorDetails:     nil}); err != nil {
			log.Errorf("send CallError error(%v),id(%v),wsmsg(%v),wsmsg_type(%v)", err, ws.id, string(wsmsg), CallResult)
		}
		return
	}
	pendingReq, ok := ws.server.getPendingRequest(ws.id)
	if !ok {
		log.Errorf("ignoring this message may request have timed out or no request before,id(%v), wsmsg(%v),wsmsg_type(%v)", ws.id, string(wsmsg), CallResult)
		return
	}
	action := pendingReq.call.Action
	if action == "" {
		log.Errorf("action is nil, may be client response timeout or center never request,id(%v),wsmsg(%v),wsmsg_type(%v)", ws.id, string(wsmsg), CallResult)
		if err := ws.sendCallError(uniqueid, &Error{
			ErrorCode:        proto.CallInternalError,
			ErrorDescription: fmt.Sprintf("may be client response timeout or center never request,uniqueid(%v)", uniqueid),
			ErrorDetails:     nil}); err != nil {
			log.Errorf("send CallError error(%v),id(%v),wsmsg(%v),wsmsg_type(%v)", err, ws.id, string(wsmsg), CallResult)
		}
		return
	}
	ocpptrait, ok := ws.server.ocpp16map.GetTraitAction(action)
	if !ok {
		log.Errorf("not support action(%v) current,id(%v),wsmsg(%v),wsmsg_type(%v)", action, ws.id, string(wsmsg), CallResult)
		if err := ws.sendCallError(uniqueid, &Error{
			ErrorCode:        proto.NotSupported,
			ErrorDescription: fmt.Sprintf("action(%v) not support current,uniqueid(%v)", action, uniqueid),
			ErrorDetails:     nil}); err != nil {
			log.Errorf("send CallError error(%v),id(%v),wsmsg(%v),wsmsg_type(%v)", err, ws.id, string(wsmsg), CallResult)
		}
		return
	}
	resType := ocpptrait.ResponseType()
	resByte, err := json.Marshal(fields[2])
	if err != nil {
		log.Errorf("json Marshal error(%v),id(%v),wsmsg(%v),wsmsg_type(%v)", err, ws.id, string(wsmsg), CallResult)
		if err = ws.sendCallError(uniqueid, &Error{
			ErrorCode:        proto.CallInternalError,
			ErrorDescription: fmt.Sprintf("json Marshal error(%v),uniqueid(%v)", err, uniqueid),
			ErrorDetails:     nil}); err != nil {
			log.Errorf("send CallError error(%v),id(%v),wsmsg(%v),wsmsg_type(%v)", err, ws.id, string(wsmsg), CallResult)
		}
		return
	}
	res := ws.server.get(resType)
	defer ws.server.put(resType, res)
	if err = json.Unmarshal(resByte, &res); err != nil {
		log.Errorf("json Unmarshal error(%v),id(%v),wsmsg(%v),wsmsg_type(%v)", err, ws.id, string(wsmsg), CallResult)
		if err = ws.sendCallError(uniqueid, &Error{
			ErrorCode:        proto.CallInternalError,
			ErrorDescription: fmt.Sprintf("json Marshal error(%v),uniqueid(%v)", err, uniqueid),
			ErrorDetails:     nil}); err != nil {
			log.Errorf("send CallError error(%v),id(%v),wsmsg(%v),wsmsg_type(%v)", err, ws.id, string(wsmsg), CallResult)
		}
		return
	}
	callResult := proto.CallResult{
		MessageTypeID: proto.CALL_RESULT,
		UniqueID:      uniqueid,
		Response:      res.(proto.Response),
	}
	if err = ws.server.validate.Struct(callResult); err != nil {
		log.Errorf("validate CallResult error(%v),id(%v),wsmsg(%v),wsmsg_type(%v)", checkValidatorError(err, action), ws.id, string(wsmsg), CallResult)
		if err = ws.sendCallError(uniqueid, checkValidatorError(err, action)); err != nil {
			log.Errorf("send CallError error(%v),id(%v),wsmsg(%v),wsmsg_type(%v)", err, ws.id, string(wsmsg), CallResult)
		}
		return
	}
	ws.server.requestDone(ws.id, uniqueid)
	ws.responseHandler(uniqueid, action, res.(proto.Response))
}

func (ws *wsconn) callErrorHandler(uniqueid string, wsmsg []byte, fields []interface{}) {
	if len(fields) != 5 {
		log.Errorf("invalid num of call fields(%+v), id(%v), wsmsg(%v),wsg_type(%v),exptect 5 fields", fields, ws.id, string(wsmsg), CallError)
		if err := ws.sendCallError(uniqueid, &Error{
			ErrorCode:        proto.FormationViolation,
			ErrorDescription: fmt.Sprintf("invalid num of callresult fields(%+v),exptect 5 fields,uniqueid(%v)", fields, uniqueid),
			ErrorDetails:     nil}); err != nil {
			log.Errorf("send CallError error(%v),id(%v),wsmsg(%v),wsmsg_type(%v)", err, ws.id, string(wsmsg), CallError)
		}
		return
	}
	errCode, ok := fields[2].(string)
	if !ok {
		log.Errorf("invalid CallError ErrCode(%v) type,must be string, id(%v), wsmsg(%v), wsmag_type(%v)", fields[2], ws.id, string(wsmsg), CallError)
		if err := ws.sendCallError(uniqueid, &Error{
			ErrorCode:        proto.TypeConstraintViolation,
			ErrorDescription: fmt.Sprintf("invalid CallError errCode(%v) type,must be string,uniqueid(%v)", fields[2], uniqueid),
			ErrorDetails:     nil}); err != nil {
			log.Errorf("send CallError error(%v),id(%v),wsmsg(%v),wsmsg_type(%v)", err, ws.id, string(wsmsg), CallError)
		}
		return
	}
	errorDescription, ok := fields[2].(string)
	if !ok {
		log.Errorf("invalid CallError errorDescription(%v) type,must be string, id(%v), wsmsg(%v), wsmag_type(%v)", fields[2], ws.id, string(wsmsg), CallError)
		if err := ws.sendCallError(uniqueid, &Error{
			ErrorCode:        proto.TypeConstraintViolation,
			ErrorDescription: fmt.Sprintf("invalid CallError errorDescription(%v) type,must be string,uniqieid(%v)", fields[2], uniqueid),
			ErrorDetails:     nil}); err != nil {
			log.Errorf("send CallError error(%v),id(%v),wsmsg(%v),wsmsg_type(%v)", err, ws.id, string(wsmsg), CallError)
		}
		return
	}
	pendingReq, ok := ws.server.getPendingRequest(ws.id)
	if !ok {
		log.Errorf("ignoring this message may request have timed out or no request before,id(%v), wsmsg(%v),wsmsg_type(%v)", ws.id, string(wsmsg), CallError)
		return
	}
	action := pendingReq.call.Action
	if action == "" {
		log.Errorf("action is nil, may be client response timeout or center never request,id(%v),wsmsg(%v),wsmsg_type(%v)", ws.id, string(wsmsg), CallError)
		if err := ws.sendCallError(uniqueid, &Error{
			ErrorCode:        proto.CallInternalError,
			ErrorDescription: fmt.Sprintf("may be client response timeout or center never request,uniqueid(%v)", uniqueid),
			ErrorDetails:     nil}); err != nil {
			log.Errorf("send CallError error(%v),id(%v),wsmsg(%v),wsmsg_type(%v)", err, ws.id, string(wsmsg), CallError)
		}
		return
	}
	callError := proto.CallError{
		MessageTypeID:    proto.CALL,
		UniqueID:         uniqueid,
		ErrorCode:        proto.ErrCodeType(errCode),
		ErrorDescription: errorDescription,
		ErrorDetails:     fields[4],
	}
	if err := ws.server.validate.Struct(callError); err != nil {
		log.Errorf("validate CallError error(%v),id(%v),wsmsg(%v),wsmsg_type(%v)", checkValidatorError(err, action), ws.id, string(wsmsg), CallError)
		action := pendingReq.call.Action
		if err = ws.sendCallError(uniqueid, checkValidatorError(err, action)); err != nil {
			log.Errorf("send CallError error(%v),id(%v),wsmsg(%v),wsmsg_type(%v)", err, ws.id, string(wsmsg), CallError)
		}
		return
	}
	ws.server.requestDone(ws.id, uniqueid)
	ws.responseHandler(uniqueid, proto.CallErrorName, &callError)
}

func (ws *wsconn) messageHandler(wsmsg []byte) {
	fields, err := parseMessage(wsmsg)
	if err != nil {
		log.Errorf("parse wsmessage error, id(%v), wsmsg(%v), err(%v)", ws.id, string(wsmsg), err)
		return
	}
	if len(fields) < 3 {
		log.Errorf("invalid wsmessage because of fields < 3, id(%v), msg(%v)", ws.id, string(wsmsg))
		return
	}

	wsmsgTypeid, ok := fields[0].(float64)
	if !ok {
		log.Errorf("invalid wsmsgTypeID(%v), type(%v),must be float64, id(%v), msg(%v)", wsmsgTypeid, reflect.TypeOf(wsmsgTypeid).String(), ws.id, string(wsmsg))
		return
	}
	uniqueid, ok := fields[1].(string)
	if !ok {
		log.Errorf("invalid uniqueid(%v), type(%v),must be string, id(%v),msg(%v)", fields[1], reflect.TypeOf(fields[1]).String(), ws.id, string(wsmsg))
		return
	}
	switch wsmsgTypeid {
	case proto.CALL:
		ws.callHandler(uniqueid, wsmsg, fields)
	case proto.CALL_RESULT:
		ws.callResultHandler(uniqueid, wsmsg, fields)
	case proto.CALL_ERROR:
		ws.callErrorHandler(uniqueid, wsmsg, fields)
	default:
		log.Errorf("not support wsmsgTypeID(%v) current", wsmsgTypeid)
	}
	return
}

func (ws *wsconn) write() {
	for {
		select {
		case ping := <-ws.ping:
			ws.writeMessage(websocket.PongMessage, ping)
		case closeError := <-ws.closeC:
			log.Errorf("id(%v) closed, err(%v)", ws.id, closeError)
			return
		}
	}
}
