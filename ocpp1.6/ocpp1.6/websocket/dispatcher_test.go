package websocket

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	randn "math/rand"
	"net/url"
	local "ocpp16/plugin/local"
	"ocpp16/proto"
	// registry "ocpp16/registry/rpcx"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	// "github.com/smallnest/rpcx/client"
	"ocpp16/logwriter"
	"sync"
	"testing"
	"time"
)

//go test -timeout=30m -v
var mx sync.Mutex
var r = randn.New(randn.NewSource(time.Now().Unix()))
var addr = flag.String("addr", "127.0.0.1:8090", "websocket service address")

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

var fnBootNotificationRequest = func() proto.BootNotificationRequest {
	return proto.BootNotificationRequest{
		ChargePointVendor:       "qinglianyun",
		ChargePointModel:        "sujunkang",
		ChargePointSerialNumber: RandString(15),
		ChargeBoxSerialNumber:   RandString(15),
		FirmwareVersion:         RandString(15),
		Iccid:                   RandString(15),
		Imsi:                    RandString(15),
		MeterType:               RandString(15),
		MeterSerialNumber:       RandString(15),
	}
}

var fnStatusNotificationRequest = func() proto.StatusNotificationRequest {
	return proto.StatusNotificationRequest{ //valid request
		ConnectorId:     15,
		ErrorCode:       "ConnectorLockFailure",
		Info:            RandString(40),
		Status:          "Available",
		Timestamp:       time.Now().Format(proto.ISO8601),
		VendorId:        RandString(240),
		VendorErrorCode: RandString(40),
	}
	// return proto.StatusNotificationRequest{ //invalid request
	// 	ConnectorId:     1,
	// 	ErrorCode:       "ConnectorLockFailure",
	// 	Info:            RandString(40),
	// 	Status:          "Available",
	// 	Timestamp:       time.Now().Format(proto.ISO8601),
	// 	VendorId:        RandString(260),
	// 	VendorErrorCode: RandString(55),
	// }
}

var fnAuthorizeRequest = func() proto.AuthorizeRequest {
	return proto.AuthorizeRequest{
		IdTag: "qinglianyun",
	}
}

var fnMeterValueRequest = func() proto.MeterValuesRequest {
	var meterValueReq = proto.MeterValuesRequest{
		ConnectorId:   10,
		TransactionId: 81043077757669376,
	}
	var meterValue = proto.MeterValue{
		Timestamp: time.Now().Format(time.RFC3339),
	}
	var sampledValue = proto.SampledValue{
		Value:   "100",
		Context: "",
		// Context:   "Interruption.Begin",
		Format: "Raw",
		// Measurand: "Energy.Active.Export.Register",
		Measurand: "",
		// Phase:     "L1",
		Phase: "",
		// Location:  "Cable",
		Location: "",
		Unit:     "Wh",
	}
	meterValue.SampledValue = append(meterValue.SampledValue, sampledValue)
	meterValueReq.MeterValue = append(meterValueReq.MeterValue, meterValue)
	return meterValueReq
}

var fnStartTransactionRequest = func() proto.StartTransactionRequest {
	return proto.StartTransactionRequest{
		ConnectorId:   10,
		IdTag:         "qinglianyun",
		MeterStart:    10,
		ReservationId: 10,
		Timestamp:     time.Now().Format(proto.ISO8601),
	}
}
var fnStopTransactionRequest = func() proto.StopTransactionRequest {
	var meterValue = proto.MeterValue{
		Timestamp: time.Now().Format(proto.ISO8601),
	}
	var sampledValue = proto.SampledValue{
		Value:     RandString(10),
		Context:   "Interruption.Begin",
		Format:    "Raw",
		Measurand: "Energy.Active.Export.Register",
		Phase:     "L1",
		Location:  "Cable",
		Unit:      "Wh",
	}
	meterValue.SampledValue = append(meterValue.SampledValue, sampledValue)
	return proto.StopTransactionRequest{
		IdTag:           "qinglianyun",
		MeterStop:       100,
		Timestamp:       time.Now().Format(proto.ISO8601),
		TransactionId:   80030155044556800,
		Reason:          "EmergencyStop",
		TransactionData: []proto.MeterValue{meterValue},
	}
}

func clientHandler(ctx context.Context, t *testing.T, d *dispatcher) {
	flag.Parse()
	// name, id := RandString(5), RandString(5)
	name, id := "qilianyun", "lihuaye"
	path := fmt.Sprintf("/ocpp/%v/%v", name, id)
	u := url.URL{Scheme: "ws", Host: "182.92.132.15:8090", Path: path}
	// u := url.URL{Scheme: "ws", Host: *addr, Path: path}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatal("dial:", err)
	}
	defer c.Close()
	ch := make(chan string, 10)
	defer close(ch)
	// queue := NewQueue()
	var waitgroup sync.WaitGroup
	waitgroup.Add(1)
	var mtx sync.Mutex

	go func() { //test for center request
		defer waitgroup.Done()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				// call := &proto.Call{
				// 	MessageTypeID: proto.CALL,
				// 	UniqueID:      RandString(7),
				// 	Action:        "BootNotification",
				// 	Request:       fnBootNotificationRequest(),
				// }
				// queue.Push(call.UniqueID)
				// if err := d.appendRequest(context.Background(), fmt.Sprintf("%v-%v", name, id), call); err != nil {
				// 	return
				// }
				// time.Sleep(time.Second * time.Duration(randn.Intn(3)) / 5)
				time.Sleep(time.Second * 10000)
			}
		}
	}()
	// waitgroup.Add(1)
	// go func() {
	// 	defer waitgroup.Done()
	// 	for {
	// 		select {
	// 		case <-ctx.Done():
	// 			return
	// 		case res_uniqueid := <-ch:
	// 			rep_uniqueid, _ := queue.Pop()
	// 			// next_uniqueid, _ := queue.Peek()
	// 			// t.Logf("ws_id(%v), res_uniqueid(%v),rep_uniqueid(%v),queue remain(%v), next_uniqueid(%v)", fmt.Sprintf("%v-%v", name, id), res_uniqueid, rep_uniqueid, queue.Len(), next_uniqueid)
	// 			if res_uniqueid != rep_uniqueid {
	// 				t.Errorf("ws_id(%v), res_uniqueid(%v) != rep_uniqueid(%v)", fmt.Sprintf("%v-%v", name, id), res_uniqueid, rep_uniqueid)
	// 			}
	// 		}
	// 	}
	// }()
	waitgroup.Add(1)
	go func() {
		defer waitgroup.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, message, err := c.ReadMessage()
				if err != nil {
					return
				}
				fields, err := parseMessage(message)
				if err != nil {
					return
				}
				switch fields[0].(float64) {
				case float64(proto.CALL):
					go func() {
						uniqueid := fields[1].(string)
						callResult := &proto.CallResult{
							MessageTypeID: proto.CALL_RESULT,
							UniqueID:      uniqueid,
							Response: &proto.BootNotificationResponse{
								CurrentTime: time.Now().Format(time.RFC3339),
								Interval:    10,
								Status:      "Accepted",
							},
						}
						callResultMsg, err := json.Marshal(callResult)
						if err != nil {
							return
						}
						time.Sleep(time.Second * time.Duration(randn.Intn(3)) / 10)
						t.Logf("test for center call: recv msg(%+v), resp_msg(%+v)", string(message), string(callResultMsg))
						mtx.Lock()
						// err = c.WriteMessage(websocket.TextMessage, callResultMsg)
						mtx.Unlock()
						if err != nil {
							return
						}
						// ch <- callResult.UniqueID
					}()
				case float64(proto.CALL_RESULT), float64(proto.CALL_ERROR):
					t.Logf("test for client call: recv msg(%v), ", string(message))
				default:
				}

			}
		}
	}()
	//test for client call
	waitgroup.Add(1)
	go func() {
		defer waitgroup.Done()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				var action = "MeterValues"
				call := &proto.Call{
					MessageTypeID: proto.CALL,
					UniqueID:      RandString(7),
					Action:        action,
				}
				switch action {
				case "StatusNotification":
					call.Request = fnStatusNotificationRequest()
				case "Authorize":
					call.Request = fnAuthorizeRequest()
				case "BootNotification":
					call.Request = fnBootNotificationRequest()
				case "MeterValues":
					call.Request = fnMeterValueRequest()
					t.Logf("%+v", call.Request)
				case "StartTransaction":
					call.Request = fnStartTransactionRequest()
				case "StopTransaction":
					call.Request = fnStopTransactionRequest()
				default:
				}
				callMsg, err := json.Marshal(call)
				if err != nil {
					t.Error(err)
					return
				}
				mtx.Lock()
				err = c.WriteMessage(websocket.TextMessage, callMsg)
				mtx.Unlock()
				if err != nil {
					t.Error(err)
					return
				}
				time.Sleep(time.Second * 100)
			}
		}
	}()
	waitgroup.Wait()
	t.Logf("(%v) grace exit gorutine", path)
}

func initLogger() *logrus.Logger {
	lw := &logwriter.HourlySplit{
		Dir:           "logs",
		FileFormat:    "log_2006-01-02T15",
		MaxFileNumber: 10,
		MaxDiskUsage:  20480000,
	}
	defer lw.Close()
	lg := logrus.New()
	customFormatter := &logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
		FullTimestamp:   true,
	}
	lg.SetFormatter(customFormatter)
	lg.SetReportCaller(true)
	lg.SetOutput(lw)
	lv, err := logrus.ParseLevel("trace")
	if err != nil {
		lv = logrus.WarnLevel
	}
	lg.SetLevel(lv)
	return lg
}

func TestDispatcherHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*100)
	lg := initLogger()
	SetLogger(lg)
	server := NewDefaultServer()
	plugin := local.NewActionPlugin()
	server.RegisterActionPlugin(plugin)
	// server.RegisterActiveCallHandler(server.HandleActiveCall, registry.NewActiveCallPlugin)
	go func() {
		server.Serve(*addr, "/ocpp/:name/:id")
	}()
	for i := 0; i < 1; i++ { //numbers of client
		time.Sleep(time.Second / 10)
		go func() {
			clientHandler(ctx, t, server.dispatcher)
		}()
	}
	select {
	case <-ctx.Done():
		time.Sleep(time.Second * 50)
		cancel()
	}
}
