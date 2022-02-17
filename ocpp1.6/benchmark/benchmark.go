package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"ocpp16/protocol"
	"sync"
	"time"
)

type quota struct {
	connectsuccess int
	connectfail    int
	sendsuccess    int
	sendfail       int
	acceptsuccess  int
	acceptfail     int
	maxtime        int64
	mintime        int64
	m              map[float64]int
}

var benchmark_quota = quota{
	m: make(map[float64]int),
}

var spend = make(chan int64, 10000)
var interval = 5
var waitgroup sync.WaitGroup

var mx sync.Mutex
var r = rand.New(rand.NewSource(time.Now().Unix()))

func randString(len int) string {
	mx.Lock()
	defer mx.Unlock()
	bytes := make([]byte, len, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

var fnBootNotificationRequest = func() protocol.BootNotificationRequest {
	return protocol.BootNotificationRequest{
		ChargePointVendor:       "qinglianyun",
		ChargePointModel:        "sujunkang",
		ChargePointSerialNumber: randString(15),
		ChargeBoxSerialNumber:   randString(15),
		FirmwareVersion:         randString(15),
		Iccid:                   randString(15),
		Imsi:                    randString(15),
		MeterType:               randString(15),
		MeterSerialNumber:       randString(15),
	}
}

var fnStatusNotificationRequest = func() protocol.StatusNotificationRequest {
	return protocol.StatusNotificationRequest{ //valid request
		ConnectorId:     10,
		ErrorCode:       "ConnectorLockFailure",
		Info:            randString(40),
		Status:          "Available",
		Timestamp:       time.Now().Format(protocol.ISO8601),
		VendorId:        randString(240),
		VendorErrorCode: randString(40),
	}
}

var fnAuthorizeRequest = func() protocol.AuthorizeRequest {
	return protocol.AuthorizeRequest{
		IdTag: "qinglianyun",
	}
}

var fnMeterValueRequest = func() protocol.MeterValuesRequest {
	var meterValueReq = protocol.MeterValuesRequest{
		ConnectorId:   10,
		TransactionId: 22,
	}
	var meterValue = protocol.MeterValue{
		Timestamp: time.Now().Format(protocol.ISO8601),
	}
	var sampledValue = protocol.SampledValue{
		Value:   "50000",
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

var fnStartTransactionRequest = func() protocol.StartTransactionRequest {
	meterStart := 10
	return protocol.StartTransactionRequest{
		ConnectorId:   10,
		IdTag:         "qinglianyun",
		MeterStart:    &meterStart,
		ReservationId: 10,
		Timestamp:     time.Now().Format(protocol.ISO8601),
	}
}
var fnStopTransactionRequest = func() protocol.StopTransactionRequest {
	var meterValue = protocol.MeterValue{
		Timestamp: time.Now().Format(protocol.ISO8601),
	}
	var sampledValue = protocol.SampledValue{
		Value:     randString(10),
		Context:   "Interruption.Begin",
		Format:    "Raw",
		Measurand: "Energy.Active.Export.Register",
		Phase:     "L1",
		Location:  "Cable",
		Unit:      "Wh",
	}
	meterValue.SampledValue = append(meterValue.SampledValue, sampledValue)
	return protocol.StopTransactionRequest{
		IdTag:           "qinglianyun",
		MeterStop:       50000,
		Timestamp:       time.Now().Format(protocol.ISO8601),
		TransactionId:   22,
		Reason:          "EmergencyStop",
		TransactionData: []protocol.MeterValue{meterValue},
	}
}

var fnResetRequest = func() *protocol.ResetRequest {
	return &protocol.ResetRequest{
		Type: protocol.ResetType("Hard"),
	}
}

var fnHeartbeatRequest = func() *protocol.HeartbeatRequest {
	return &protocol.HeartbeatRequest{}
}

func main() {
	var flag_n = flag.Int("n", 0, "total connections")
	var flag_c = flag.Int("c", 0, "number of concurrent connections")
	var flag_host = flag.String("h", "", "ws connection address")
	flag.Parse()
	if *flag_n == 0 {
		log.Fatal("parameter -n must be specified")
	}
	n := *flag_n
	if *flag_c == 0 {
		log.Fatal("parameter -n must be specified")
	}
	c := *flag_c
	if n <= 0 || c <= 0 {
		log.Fatalln("parameter -n or -c must > 0")
	}
	if *flag_host == "" {
		log.Fatal("parameter -url must be specified")
	}
	host := *flag_host
	var action = "Heartbeat"
	call := &protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      randString(7),
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
	case "StartTransaction":
		call.Request = fnStartTransactionRequest()
	case "StopTransaction":
		call.Request = fnStopTransactionRequest()
	case "Heartbeat":
		call.Request = fnHeartbeatRequest()
	}
	msg, err := json.Marshal(call)
	if err != nil {
		return
	}
	start := time.Now().UnixNano()
	for i := 1; i <= n/c; i++ {
		for j := 1; j <= c; j++ {
			waitgroup.Add(1)
			go connect(host, msg)
		}
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
	total_request := 0
	var total_consume int64 = 0
	mux := sync.Mutex{}
	for i := 1; i <= n; i++ {
		go func() {
			select {
			case request_time := <-spend:
				total_request++
				mux.Lock()
				if count, ok := benchmark_quota.m[math.Ceil(float64(request_time/1000000))]; ok {
					count = count + 1
					benchmark_quota.m[math.Ceil(float64(request_time/1000000))] = count
				} else {
					benchmark_quota.m[math.Ceil(float64(request_time/1000000))] = 1
				}
				mux.Unlock()
				total_consume += request_time
				if benchmark_quota.maxtime < request_time {
					benchmark_quota.maxtime = request_time
				}
				if benchmark_quota.mintime == 0 || benchmark_quota.mintime > request_time {
					benchmark_quota.mintime = request_time
				}
			}
		}()
	}
	waitgroup.Wait()
	end := time.Now().UnixNano()
	// fmt.Println("total request : ", total_request)
	totalTime := float64(end-start)/1000000000 - float64(interval*(n/c)/1000)
	fmt.Println("spend time : ", totalTime, "s")
	fmt.Printf("qps : %.2f [#/sec] \r\n", float64(total_request)/totalTime)
	fmt.Println("avg request min time : ", float64(total_consume)/float64(total_request)/1000000, "ms")
	fmt.Println("per request min time : ", float64(benchmark_quota.mintime)/1000000, "ms")
	fmt.Println("per request max time : ", float64(benchmark_quota.maxtime)/1000000, "ms")

	fmt.Println("connect success: ", benchmark_quota.connectsuccess)
	fmt.Println("connect fail: ", benchmark_quota.connectfail)
	fmt.Println("send message success: ", benchmark_quota.sendsuccess)
	fmt.Println("send message fail: ", benchmark_quota.sendfail)
	fmt.Println("accept message success: ", benchmark_quota.acceptsuccess)
	fmt.Println("accept message fail: ", benchmark_quota.acceptfail)
	for k, v := range benchmark_quota.m {
		fmt.Printf("%vms:%v\n", k, v)
	}
}

func connect(host string, msg []byte) {
	defer waitgroup.Done()
	name, id := randString(5), randString(5)
	path := fmt.Sprintf("/ocpp/%s/%s", name, id)
	u := url.URL{Scheme: "ws", Host: host, Path: path}
	dialer := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
		Subprotocols:     []string{"ocpp1.5", "ocpp1.6"},
	}
	ws, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		benchmark_quota.connectfail++
		return
	}
	defer ws.Close()
	benchmark_quota.connectsuccess++
	var start int64
	go func() {
		start = time.Now().UnixNano()
		err = ws.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			benchmark_quota.sendfail++
			return
		}
		benchmark_quota.sendsuccess++
	}()
	_, _, err = ws.ReadMessage()
	if err != nil {
		benchmark_quota.acceptfail++
		return
	}
	benchmark_quota.acceptsuccess++
	spend <- time.Now().UnixNano() - start
}
