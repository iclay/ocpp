package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"ocpp16/proto"
	"runtime"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type (
	request struct {
		call    *proto.Call
		reqTime string
	}

	callStateMap struct {
		pendingCallState map[string]*request
		sync.RWMutex
	}
)

func newCallStateMap() *callStateMap {
	return &callStateMap{
		pendingCallState: make(map[string]*request),
	}
}

func (m *callStateMap) createNewRequest(id string) {
	m.Lock()
	defer m.Unlock()
	m.pendingCallState[id] = &request{
		call: &proto.Call{},
	}
}

func (m *callStateMap) getPendingRequest(uniqueid string) (*request, bool) {
	m.Lock()
	defer m.Unlock()
	req, ok := m.pendingCallState[uniqueid]
	return req, ok
}

func (m *callStateMap) deleteRequest(id string) {
	m.Lock()
	defer m.Unlock()
	delete(m.pendingCallState, id)
}

func (m *callStateMap) AddRequest(id string, request *request) error {
	m.Lock()
	defer m.Unlock()
	if req, ok := m.pendingCallState[id]; ok && request.call.QueryUniqueID() != "" && req.call.QueryUniqueID() == "" {
		m.pendingCallState[id] = request
		return nil
	} else {
		return fmt.Errorf("call state is not empty or connection close or other errors, so add request failed, id(%v),request(%+v)", id, request)
	}
}

func (m *callStateMap) requestDone(id string, uniqueid string) {
	m.Lock()
	defer m.Unlock()
	if req, ok := m.pendingCallState[id]; ok {
		if req.call.QueryUniqueID() == uniqueid {
			m.pendingCallState[id] = &request{
				call: &proto.Call{},
			}
		}

	}
}

type requestQueueMap struct {
	queueMap map[string]Queue
	sync.RWMutex
}

func newRequesQueueMap() *requestQueueMap {
	return &requestQueueMap{
		queueMap: make(map[string]Queue),
	}
}

func (m *requestQueueMap) createNewQueue(id string) {
	m.Lock()
	defer m.Unlock()
	m.queueMap[id] = NewQueue()
	return
}

func (m *requestQueueMap) queueExists(id string) bool {
	m.RLock()
	defer m.RUnlock()
	_, ok := m.queueMap[id]
	return ok
}

func (m *requestQueueMap) getQueue(id string) (Queue, bool) {
	m.RLock()
	defer m.RUnlock()
	q, ok := m.queueMap[id]
	return q, ok
}

func (m *requestQueueMap) deleteQueue(id string) {
	m.Lock()
	defer m.Unlock()
	delete(m.queueMap, id)
}

func (m *requestQueueMap) pushRequset(id string, request interface{}) error {
	m.Lock()
	defer m.Unlock()
	queue, ok := m.queueMap[id]
	if !ok {
		return fmt.Errorf("push request failed, may be conn has closed down, id(%v), request(%+v)", id, request)
	}
	queue.Push(request)
	return nil
}

type dispatcher struct {
	server          *Server
	callStateMap    *callStateMap
	requestQueueMap *requestQueueMap
	requestC        chan string
	nextReadyC      chan string
	timeout         time.Duration
	timeoutC        chan timeoutFlag
	stopC           chan error
}

func NewDefaultDispatcher(s *Server) (d *dispatcher) {
	d = &dispatcher{
		server:          s,
		callStateMap:    newCallStateMap(),
		requestQueueMap: newRequesQueueMap(),
		requestC:        make(chan string, 10),
		nextReadyC:      make(chan string),
		timeout:         time.Second * 5,
		timeoutC:        make(chan timeoutFlag),
		stopC:           make(chan error),
	}
	go d.run()
	return d
}

type timeoutContext struct {
	ctx      context.Context
	cancel   context.CancelFunc
	uniqueid string
}

func (ctx *timeoutContext) isActive() bool {
	return ctx.cancel != nil
}

//don't try to modify the code unless you know what to do
func (d *dispatcher) run() {
	contextMap := make(map[string]timeoutContext)
	var allow bool
	var ok bool
	var q Queue
	var id string
	var ctx timeoutContext
	defer func() {
		close(d.requestC)
		close(d.nextReadyC)
		close(d.timeoutC)
		for _, timeoutCtx := range contextMap {
			if timeoutCtx.isActive() {
				timeoutCtx.cancel() //notify all timeout goroutine to exit
			}
		}
	}()
	for {
		select {
		case <-d.stopC:
			log.Debugf("dispatcher has stopped")
			return
		case id = <-d.requestC:
			q, ok = d.requestQueueMap.getQueue(id)
			if !ok { //the connection may have been closed
				if ctx, ok := contextMap[id]; ok { //this may happen when there has been a request but the connection has been closed
					if ctx.isActive() {
						ctx.cancel() //Cancel the timeout and notice exit the goroutine
					}
					delete(contextMap, id)
				}
				continue
			}
			if ctx, ok = contextMap[id]; !ok { //the first request, so the write can be triggered
				allow = true
			} else {
				allow = !ctx.isActive() //at this time, it is idle and can trigger write
			}
		case id = <-d.nextReadyC: //the timeout mechanism or a correct return has been triggered
			if ctx, ok = contextMap[id]; ok && ctx.isActive() {
				ctx.cancel()
				contextMap[id] = timeoutContext{}
			}
			if q, ok = d.requestQueueMap.getQueue(id); ok {
				allow = true
			}
		case timeOutFlag := <-d.timeoutC: //timeout trigger
			pendingReq, ok := d.callStateMap.getPendingRequest(timeOutFlag.uniqueid)
			if !ok {
				continue
			}
			if ctx, ok = contextMap[id]; ok && ctx.isActive() && timeOutFlag.uniqueid == pendingReq.call.QueryUniqueID() {
				ctx.cancel()
				d.requestDone(id, ctx.uniqueid)
				contextMap[id] = timeoutContext{}
			}
		}
		if allow && !q.IsEmpty() {
			contextMap[id] = d.dispatchNextRequest(id)
			allow = false
		}
	}
}

type timeoutFlag struct {
	id       string
	uniqueid string
}

func (d *dispatcher) dispatchNextRequest(id string) (timeoutCtx timeoutContext) {
	q, ok := d.requestQueueMap.getQueue(id)
	if !ok {
		log.Errorf("get queue error, conn may be close, id(%v)", id)
		return
	}
	req, ok := q.Peek()
	if !ok {
		log.Errorf("queue peek is empty,id(%v)", id)
		return
	}
	request := req.(*request)
	request.reqTime = time.Now().Format(time.RFC3339)
	call := request.call
	uniqueid := call.QueryUniqueID()
	message, err := json.Marshal(call)
	if err != nil {
		log.Errorf("json Marshal error is error, id(%v),uniqueid(%v), request(%+v)", id, uniqueid, request)
		return
	}
	if err = d.callStateMap.AddRequest(id, request); err != nil {
		log.Error(err)
		return
	}
	ws, ok := d.server.getConn(id)
	if !ok {
		d.requestDone(id, uniqueid)
		log.Errorf("get ws conn error, conn may be close, id(%v), uniqueid(%v), request(%+v)", id, uniqueid, request)
		return
	}
	err = ws.writeMessage(websocket.TextMessage, message)
	if err != nil {
		d.requestDone(id, uniqueid)
		log.Errorf("write message error, conn may be close or other errors, id(%v), uniqueid(%v), request(%+v), err(%v)", id, uniqueid, request, err)
		return
	}
	ctx := func() (timeoutCtx timeoutContext) {
		if d.timeout < 0 {
			return
		}
		ctx, cancel := context.WithTimeout(context.TODO(), d.timeout)
		timeoutCtx = timeoutContext{ctx: ctx, cancel: cancel, uniqueid: uniqueid}
		go func() {
			runtime.Gosched()
			select {
			case <-ctx.Done():
				switch ctx.Err() {
				case context.DeadlineExceeded:
					d.timeoutC <- timeoutFlag{
						id:       id,
						uniqueid: uniqueid,
					}
				default:
					log.Debugf("timeoutC has cancald due to valid response or connection close id(%v), uniqueid(%v), request(%+v)", id, uniqueid, request)
				}
			}
		}()
		return
	}()
	return ctx
}

func (d *dispatcher) requestDone(id string, uniqueid string) {
	var q Queue
	q, ok := d.requestQueueMap.getQueue(id)
	if !ok {
		log.Errorf("get queue error, conn may be close, id(%v), uniqueid(%v)", id, uniqueid)
		return
	}
	requestQueue := q.(*requestQueue)
	req, ok := requestQueue.Peek()
	if !ok {
		log.Errorf("queue peek is empty,id(%v), uniqueid(%v)", id, uniqueid)
		return
	}
	request := req.(*request)
	if request.call.QueryUniqueID() != uniqueid {
		log.Errorf("requestid is not equal to uniqueid,maybe due to request timeout, id(%v), uniqueid(%v),latest request=(%+v)", id, uniqueid, request)
		return
	}
	requestQueue.Pop()
	d.callStateMap.requestDone(id, uniqueid)
	log.Debug("request has already complete, id(%v), uniqueid(%v), request(%+v)", id, uniqueid, request)
	d.nextReadyC <- id
}

func (d *dispatcher) appendRequest(id string, call *proto.Call) error {
	if call == nil || call.UniqueID == "" {
		return fmt.Errorf("append request failed,call is nil or uniqueid is nil,id(%v),call(%+v)", id, call)
	}
	if err := d.server.validate.Struct(call); err != nil {
		return fmt.Errorf("append request failed because of vaild call error,id(%v),call(%+v), err(%v)", id, call, checkValidatorError(err, call.Action))
	}
	if err := d.requestQueueMap.pushRequset(id, &request{
		call: call,
	}); err != nil {
		return err
	}
	d.requestC <- id
	return nil
}
