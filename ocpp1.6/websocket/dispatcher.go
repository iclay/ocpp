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

func (m *callStateMap) CreateNewRequest(id string) {
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
		return fmt.Errorf("call state is not empty or connection close or other errors, so add request failed, id=%v,request=%+v", id, request)
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

func (m *requestQueueMap) CreateNewQueue(id string) {
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

func (m *requestQueueMap) GetQueue(id string) (Queue, bool) {
	m.RLock()
	defer m.RUnlock()
	q, ok := m.queueMap[id]
	return q, ok
}

func (m *requestQueueMap) DeleteQueue(id string) {
	m.Lock()
	defer m.Unlock()
	delete(m.queueMap, id)
}

func (m *requestQueueMap) pushRequset(id string, request interface{}) error {
	m.Lock()
	defer m.Unlock()
	queue, ok := m.queueMap[id]
	if !ok { //不存在，说明可能已经删除
		return fmt.Errorf("push request failed, may be conn has closed down, id=%v, request=%+v", id, request)
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
	timeoutC        chan string
	stopC           chan error
}

func NewDefaultDispatcher(s *Server) (d *dispatcher) {
	d = &dispatcher{
		server:          s,
		callStateMap:    newCallStateMap(),
		requestQueueMap: newRequesQueueMap(),
		requestC:        make(chan string),
		nextReadyC:      make(chan string),
		timeout:         time.Second * 5,
		timeoutC:        make(chan string),
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
				timeoutCtx.cancel() //通知所有的超时协程退出
			}
		}
	}()
	for {
		select {
		case <-d.stopC:
			log.Debugf("dispatcher has stopped")
			return
		case id = <-d.requestC: //clientid触发,如果队列不为空，则取出
			q, ok = d.requestQueueMap.GetQueue(id)
			if !ok { //不存在，说明连接已经关闭，已经删除
				if ctx, ok := contextMap[id]; ok { //可能出现的情况是，有过下行数据，但是连接关闭，但是clientid key值还在,需要删除
					if ctx.cancel != nil {
						ctx.cancel() //立刻触发超时机制
					}
					delete(contextMap, id)
				}
				continue
			}
			if ctx, ok = contextMap[id]; !ok { //第一次request
				allow = true
			} else {
				allow = !ctx.isActive() //说明此时空闲
			}
		case id = <-d.nextReadyC: //此时说明已经完成上一次请求，此处必须是有正确的返回
			ctx = contextMap[id]
			if ctx.isActive() {
				ctx.cancel()
				contextMap[id] = timeoutContext{}
			}
			q, ok = d.requestQueueMap.GetQueue(id) //说明连接还在，此时由于是nextReady说明contextMap一定也有，不用判断
			if ok {
				allow = true
			}
		case id = <-d.timeoutC: //id已经超时，需要删除
			ctx = contextMap[id]
			if ctx.isActive() {
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

func (d *dispatcher) dispatchNextRequest(id string) (timeoutCtx timeoutContext) {
	q, ok := d.requestQueueMap.GetQueue(id)
	if !ok {
		log.Debugf("get queue error, conn may be close, id=%v", id)
		return
	}
	req, ok := q.Peek()
	if !ok {
		log.Debugf("queue peek is empty,id=%v", id)
		return
	}
	request := req.(*request)
	request.reqTime = time.Now().Format(time.RFC3339)
	call := request.call
	uniqueid := call.QueryUniqueID()
	message, err := json.Marshal(call)
	if err != nil {
		log.Debugf("json Marshal error is error, id=%v,uniqueid=%v, request=%+v", id, uniqueid, request)
		return
	}
	if err = d.callStateMap.AddRequest(id, request); err != nil {
		log.Error(err)
		return
	}
	ws, ok := d.server.getConn(id)
	if !ok {
		d.requestDone(id, uniqueid)
		log.Errorf("get ws conn error, conn may be close, id=%v, uniqueid=%v, request=%+v", id, uniqueid, request)
		return
	}
	ws.writeMessage(websocket.TextMessage, message)
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
					d.timeoutC <- id
				default:
					log.Debugf("timeoutC has cancald due to valid response or connection close id=%v, uniqueid=%v, request=%+v", id, uniqueid, request)
				}
			}
		}()
		return
	}()
	return ctx
}

func (d *dispatcher) requestDone(id string, uniqueid string) {
	var q Queue
	q, ok := d.requestQueueMap.GetQueue(id)
	if !ok {
		log.Debugf("get queue error, conn may be close, id=%v, uniqueid=%v", id, uniqueid)
		return
	}
	requestQueue := q.(*requestQueue)
	req, ok := requestQueue.Peek()
	if !ok {
		log.Debugf("queue peek is empty,id=%v, uniqueid=%v", id, uniqueid)
		return
	}
	request := req.(*request)
	if request.call.QueryUniqueID() != uniqueid {
		log.Debugf("requestid is not equal to uniqueid,maybe due to request timeout, id=%v, uniqueid=%v,latest request=%+v", id, uniqueid, request)
		return
	}
	requestQueue.Pop()
	d.callStateMap.requestDone(id, uniqueid)
	log.Debugf("request has already complete, id=%v, uniqueid=%v, request=%+v", id, uniqueid, request)
	d.nextReadyC <- id
}

func (d *dispatcher) appendRequest(id string, call *proto.Call) error {
	if !d.requestQueueMap.queueExists(id) {
		return fmt.Errorf("queue not exists, id=%v,may be conn already close, call=%+v", id, call)
	}
	request := &request{call: call}
	if err := d.requestQueueMap.pushRequset(id, request); err != nil {
		return err
	}
	d.requestC <- id
	return nil
}
