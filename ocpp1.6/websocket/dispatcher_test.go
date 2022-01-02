package websocket

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	randn "math/rand"
	"net/url"
	"ocpp16/plugin/local"
	"ocpp16/proto"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
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

// func randomInt(min, max int) int {
// 	return min + r.Intn(max-min)
// }

// func RandString(len int) string {
// 	mx.Lock()
// 	defer mx.Unlock()
// 	bytes := make([]byte, len)
// 	for i := 0; i < len; i++ {
// 		bytes[i] = byte(randomInt(65, 90))
// 	}
// 	return string(bytes)
// }
func clientHandler(ctx context.Context, t *testing.T, d *dispatcher) {
	flag.Parse()
	name, id := RandString(5), RandString(5)
	path := fmt.Sprintf("/ocpp/%v/%v", name, id)
	u := url.URL{Scheme: "ws", Host: *addr, Path: path}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatal("dial:", err)
	}
	defer c.Close()
	ch := make(chan string, 10)
	defer close(ch)
	queue := NewQueue()
	var waitgroup sync.WaitGroup
	waitgroup.Add(1)
	var mtx sync.Mutex
	go func() {
		defer waitgroup.Done()
		fn := func() *proto.BootNotificationRequest {
			return &proto.BootNotificationRequest{
				ChargePointVendor:       RandString(15),
				ChargePointModel:        RandString(15),
				ChargePointSerialNumber: RandString(15),
				ChargeBoxSerialNumber:   RandString(15),
				FirmwareVersion:         RandString(15),
				Iccid:                   RandString(15),
				Imsi:                    RandString(15),
				MeterType:               RandString(15),
				MeterSerialNumber:       RandString(15),
			}
		}
		for {
			select {
			case <-ctx.Done():
				return
			default:
				call := &proto.Call{
					MessageTypeID: proto.CALL,
					UniqueID:      RandString(7),
					Action:        "BootNotification",
					Request:       fn(),
				}
				queue.Push(call.UniqueID)
				if err := d.appendRequest(fmt.Sprintf("%v-%v", name, id), call); err != nil {
					return
				}
				time.Sleep(time.Second * time.Duration(randn.Intn(3)) / 5)
			}
		}
	}()
	waitgroup.Add(1)
	go func() {
		defer waitgroup.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case res_uniqueid := <-ch:
				rep_uniqueid, _ := queue.Pop()
				// next_uniqueid, _ := queue.Peek()
				// t.Logf("ws_id(%v), res_uniqueid(%v),rep_uniqueid(%v),queue remain(%v), next_uniqueid(%v)", fmt.Sprintf("%v-%v", name, id), res_uniqueid, rep_uniqueid, queue.Len(), next_uniqueid)
				if res_uniqueid != rep_uniqueid {
					t.Errorf("ws_id(%v), res_uniqueid(%v) != rep_uniqueid(%v)", fmt.Sprintf("%v-%v", name, id), res_uniqueid, rep_uniqueid)
				}
			}
		}
	}()
	//test for center request
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
						//t.Logf("test for center call: recv msg(%+v), resp_msg(%+v)", string(message), string(callResultMsg))
						mtx.Lock()
						err = c.WriteMessage(websocket.TextMessage, callResultMsg)
						mtx.Unlock()
						if err != nil {
							return
						}
						ch <- callResult.UniqueID
					}()
				case float64(proto.CALL_RESULT), float64(proto.CALL_ERROR):
					//t.Logf("test for client call: recv msg(%v), ", string(message))
				default:
				}

			}
		}
	}()
	//test for client call
	waitgroup.Add(1)
	go func() {
		defer waitgroup.Done()
		fn := func() *proto.StatusNotificationRequest {
			return &proto.StatusNotificationRequest{ //valid request
				ConnectorId:     15,
				ErrorCode:       "ConnectorLockFailure",
				Info:            RandString(40),
				Status:          "Available",
				Timestamp:       time.Now().Format(proto.ISO8601),
				VendorId:        RandString(240),
				VendorErrorCode: RandString(40),
			}
			// return &proto.StatusNotificationRequest{ //invalid request
			// 	ConnectorId:     1,
			// 	ErrorCode:       "ConnectorLockFailure",
			// 	Info:            RandString(40),
			// 	Status:          "Available",
			// 	Timestamp:       time.Now().Format(proto.ISO8601),
			// 	VendorId:        RandString(260),
			// 	VendorErrorCode: RandString(55),
			// }
		}
		for {
			select {
			case <-ctx.Done():
				return
			default:
				call := &proto.Call{
					MessageTypeID: proto.CALL,
					UniqueID:      RandString(7),
					Action:        "StatusNotification",
					Request:       fn(),
				}
				callMsg, err := json.Marshal(call)
				if err != nil {
					return
				}
				mtx.Lock()
				err = c.WriteMessage(websocket.TextMessage, callMsg)
				mtx.Unlock()
				if err != nil {
					return
				}
				time.Sleep(time.Second * 1)
			}
		}
	}()
	waitgroup.Wait()
	t.Logf("(%v) grace exit gorutine", path)
}

func TestDispatcherHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*100)
	server := NewDefaultServer()
	plugin := local.NewLocalService()
	server.RegisterOCPPHandler(plugin)
	go func() {
		server.Serve(*addr, "/ocpp/:name/:id")
	}()
	for i := 0; i < 20; i++ { //numbers of client
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
