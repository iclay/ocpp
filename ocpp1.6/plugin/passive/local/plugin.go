package local

import (
	"context"
	rand "math/rand"
	"ocpp16/proto"
	"sync"
	"time"
)

var mx sync.Mutex
var r = rand.New(rand.NewSource(time.Now().Unix()))

func RandString(len int) string {
	mx.Lock()
	defer mx.Unlock()
	bytes := make([]byte, len, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

type LocalActionPlugin struct {
	requestHandlerMap  map[string]proto.RequestHandler
	responseHandlerMap map[string]proto.ResponseHandler
}

func NewActionPlugin() *LocalActionPlugin {
	plugin := &LocalActionPlugin{}
	plugin.registerRequestHandler()
	plugin.registerResponseHandler()
	return plugin
}

func (l *LocalActionPlugin) BootNotification(ctx context.Context, request proto.Request) (proto.Response, error) {
	return &proto.BootNotificationResponse{
		CurrentTime: time.Now().Format(time.RFC3339),
		Interval:    10,
		Status:      "Accepted",
	}, nil
}

func (l *LocalActionPlugin) StatusNotification(ctx context.Context, request proto.Request) (proto.Response, error) {
	return &proto.StatusNotificationRequest{
		ConnectorId:     15,
		ErrorCode:       "ConnectorLockFailure",
		Info:            RandString(40),
		Status:          "Available",
		Timestamp:       time.Now().Format(proto.ISO8601),
		VendorId:        RandString(240),
		VendorErrorCode: RandString(40),
	}, nil
}

func (l *LocalActionPlugin) registerRequestHandler() {
	l.requestHandlerMap = map[string]proto.RequestHandler{
		proto.BootNotificationName:   proto.RequestHandler(l.BootNotification),
		proto.StatusNotificationName: proto.RequestHandler(l.StatusNotification),
	}
}

//RequestHandler represent device active request Center
func (l *LocalActionPlugin) RequestHandler(action string) (proto.RequestHandler, bool) {
	handler, ok := l.requestHandlerMap[action]
	return handler, ok
}

func (l *LocalActionPlugin) registerResponseHandler() {
	l.responseHandlerMap = map[string]proto.ResponseHandler{}
}

//ResponseHandler represent The device reply to the center request
func (l *LocalActionPlugin) ResponseHandler(action string) (proto.ResponseHandler, bool) {
	handler, ok := l.responseHandlerMap[action]
	return handler, ok
}
