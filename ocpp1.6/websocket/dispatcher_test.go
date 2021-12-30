package websocket

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	randn "math/rand"
	"net/url"
	"ocpp16/plugin/local"
	"ocpp16/proto"
	"sync"
	"testing"
	"time"
)

var r = randn.New(randn.NewSource(time.Now().Unix()))
var addr = flag.String("addr", "127.0.0.1:8090", "websocket service address")

func RandString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}
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
		for i := 0; i < 30; i++ {
			call := &proto.Call{
				MessageTypeID: proto.CALL,
				UniqueID:      RandString(7),
				Action:        "BootNotification",
				Request:       fn(),
			}
			queue.Push(call.UniqueID)
			d.appendRequest(fmt.Sprintf("%v-%v", name, id), call)
			time.Sleep(time.Second * 1)
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
				next_uniqueid, _ := queue.Peek()
				t.Logf("ws_id(%v), res_uniqueid(%v),rep_uniqueid(%v),queue remain(%v), next_uniqueid(%v)", fmt.Sprintf("%v-%v", name, id), res_uniqueid, rep_uniqueid, queue.Len(), next_uniqueid)
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
					t.Errorf("parseMessage err(%v)", err)
					return
				}
				if fields[0].(float64) != float64(proto.CALL) {
					return
				}
				uniqueid := fields[1].(string)
				result := &proto.CallResult{
					MessageTypeID: proto.CALL_RESULT,
					UniqueID:      uniqueid,
					Response: &proto.BootNotificationResponse{
						CurrentTime: time.Now().Format(time.RFC3339),
						Interval:    10,
						Status:      "Accepted",
					},
				}
				msg, err := json.Marshal(result)
				if err != nil {
					t.Fatal(err)
				}
				time.Sleep(time.Second * time.Duration(randn.Intn(5)))
				// t.Logf("client send msg(%+v), recv_msg(%+v)", string(msg), string(message))
				c.WriteMessage(websocket.TextMessage, msg)
				ch <- result.UniqueID
			}
		}
	}()
	waitgroup.Wait()
	t.Log("grace exit gorutine")
}

func Test_DispatcherHandler(t *testing.T) {
	ctx, _ := context.WithTimeout(context.TODO(), time.Second*150)
	server := NewDefaultServer()
	plugin := local.NewLocalService()
	server.RegisterOCPPHandler(plugin)
	go func() {
		server.Serve(*addr, "/ocpp/:name/:id")
	}()
	for i := 0; i < 1; i++ {
		time.Sleep(time.Second / 10)
		go func() {
			clientHandler(ctx, t, server.dispatcher)
		}()
	}
	select {
	case <-ctx.Done():
		time.Sleep(time.Second * 1)
	}
}
