package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"golang.org/x/net/websocket"
	"math/rand"
	"net/url"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	concurrency  uint64 = 1
	Num          uint64 = 1
	websocketurl string = ""
	keepalive    bool   = true
	duration     uint64 = 1
)

func init() {
	flag.Uint64Var(&concurrency, "c", concurrency, "concurrent number")
	flag.Uint64Var(&Num, "n", Num, "number of requests sent per second per connection")
	flag.StringVar(&websocketurl, "u", websocketurl, "websocket url")
	flag.BoolVar(&keepalive, "k", keepalive, "save long connection")             //whether to maintain long connection status
	flag.Uint64Var(&duration, "d", duration, "duration of pressure measurement") //stress test duration
	flag.Parse()
}

type parameters struct {
	concurrency uint64
	Num         uint64
	url         string
	keepalive   bool
	duration    uint64
}

type requestResults struct {
	id          string //messageID
	connid      uint64 //connectionID
	requestTime uint64 //rtt time
	isSucceed   bool   //is the request successful
	errCode     int    //error code
}

func main() {
	runtime.GOMAXPROCS(1)
	if concurrency == 0 || Num == 0 || websocketurl == "" {
		fmt.Println("example:go run pressure.go -c 1 -n 1 -u wss://ip:port/pressure_test")
		flag.Usage()
		return
	}
	p := &parameters{
		concurrency: concurrency,
		Num:         Num,
		url:         websocketurl,
		keepalive:   keepalive,
		duration:    duration,
	}
	startOcppPressureTest(p)
}

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

func newWSConn(wsURL string) (*websocket.Conn, error) {
	retry := 3
	u, err := url.Parse(wsURL)
	if err != nil {
		panic(err)
	}
	for i := 0; i < retry; i++ {
		conn, err := websocket.Dial(wsURL, "ocpp1.6", fmt.Sprintf("%s%s/", "http://", u.Host))
		if err != nil {
			continue
		}
		return conn, err
	}
	return nil, fmt.Errorf("connection(%s) establishment failed", wsURL)
}

func startOcppPressureTest(p *parameters) {
	ch := make(chan *requestResults, 1000)
	var wg sync.WaitGroup
	var wgRecv sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	wgRecv.Add(1)
	go receiveTestResults(ctx, p, ch, &wgRecv)
	for i := 0; uint64(i) < p.concurrency; i++ {
		path := fmt.Sprintf("/ocpp/%s/%s", randString(5), randString(5))
		conn, err := newWSConn(p.url + path)
		if err != nil {
			fmt.Printf("connection(%s) establishment success, connid(%d), err(%v)\n", path, i, err)
			return
		}
		wg.Add(1)
		go process(ctx, uint64(i), ch, &wg, p, conn)
		if uint64(i) == p.concurrency-1 {
			cancel()
		}
		time.Sleep(2 * time.Millisecond)
	}
	wg.Wait()
	time.Sleep(1 * time.Millisecond)
	close(ch)
	wgRecv.Wait()
	return
}

const (
	firstTime              = 1 * time.Second //the time when data was first requested after the connection
	intervalTime           = 1 * time.Second //time interval for requesting data
	intervalStatisticsTime = 1 * time.Second //statistical time interval of pressure measurement results
)

func process(ctx context.Context, connid uint64, ch chan<- *requestResults, wg *sync.WaitGroup, p *parameters, conn *websocket.Conn) {
	select {
	case <-ctx.Done(): //wait until all connections are established before testing
		break
	}
	timeoutCtx, _ := context.WithTimeout(context.Background(), time.Duration(p.duration)*time.Minute)
	defer func() {
		wg.Done()
		_ = conn.Close()
	}()
	var i uint64
	var err error
	t := time.NewTimer(firstTime)
	defer func() {
		t.Stop()
		if err == nil {
			if p.keepalive {
				wait := make(chan int, 0)
				<-wait
			}
		}
	}()
	for {
		select {
		case <-timeoutCtx.Done():
			err = errNormalShutDown
			return
		case <-t.C:
			t.Reset(intervalTime / time.Duration(p.Num))
			if err = simulationRequest(connid, i, ch, p, conn); err == errShutDown {
				return
			}
		}
	}
}

const (
	requestSuccess = 200
	requestError   = 509
)

var errShutDown = errors.New("conn should close")
var errNormalShutDown = errors.New("conn NormalShutDown")

func simulationRequest(connid uint64, messageid uint64, ch chan<- *requestResults, p *parameters, conn *websocket.Conn) error {
	var (
		isSucceed bool
		errCode   int
	)
	//28 byte
	// writemsg := fmt.Sprintf("[2,\"%v\",\"Heartbeat\",{}]", randString(7))
	//179 byte
	writemsg := fmt.Sprintf("[2,\"%v\",\"MeterValues\",{\"connectorId\":10,\"transactionId\":22,\"meterValue\":[{\"timeStamp\":\"2022-03-01T11:05:00Z\",\"sampledValue\":[{\"value\":\"50000\",\"format\":\"Raw\",\"unit\":\"Wh\"}]}]}]", randString(7))
	start := time.Now()
	_, err := conn.Write([]byte(writemsg))
	if err != nil {
		isSucceed, errCode = false, requestError
	} else {
		msg := make([]byte, 512)
		_, err = conn.Read(msg)
		if err != nil {
			isSucceed, errCode = false, requestError
			err = errShutDown
		} else {
			isSucceed, errCode = true, requestSuccess
		}
	}
	ch <- &requestResults{
		id:          fmt.Sprintf("%d_%d", connid, messageid),
		connid:      connid,
		requestTime: uint64(time.Since(start)),
		isSucceed:   isSucceed,
		errCode:     errCode,
	}
	return err
}

func receiveTestResults(ctx context.Context, p *parameters, ch <-chan *requestResults, wg *sync.WaitGroup) {
	select {
	case <-ctx.Done(): //wait until all connections are established before testing
		break
	}
	defer func() {
		wg.Done()
	}()
	var stop = make(chan struct{})
	var (
		processingTime  uint64
		requestTime     uint64
		maxTime         uint64
		minTime         uint64
		successNum      uint64
		failureNum      uint64
		runConcurrency  uint64
		connidMap       = make(map[uint64]bool)
		mutex           = sync.RWMutex{}
		requestTimeList []uint64
	)
	start := uint64(time.Now().UnixNano())
	errCodeMap := &sync.Map{}
	ticker := time.NewTimer(intervalStatisticsTime)
	defer ticker.Stop()
	go func() {
		for {
			select {
			case <-ticker.C:
				ticker.Reset(intervalStatisticsTime * 10)
				mutex.Lock()
				resultStatistics(p.concurrency, processingTime, uint64(time.Now().UnixNano())-start, maxTime, minTime, successNum, failureNum, runConcurrency, errCodeMap)
				mutex.Unlock()
			case <-stop:
				return
			}
		}
	}()
	printHeader()
	for data := range ch {
		mutex.Lock()
		processingTime += data.requestTime
		if maxTime <= data.requestTime {
			maxTime = data.requestTime
		}
		if minTime == 0 || minTime > data.requestTime {
			minTime = data.requestTime
		}
		if data.isSucceed == true {
			successNum += 1
		} else {
			failureNum += 1
		}
		if value, ok := errCodeMap.Load(data.errCode); ok {
			valueInt, _ := value.(int)
			errCodeMap.Store(data.errCode, valueInt+1)
		} else {
			errCodeMap.Store(data.errCode, 1)
		}
		if _, ok := connidMap[data.connid]; !ok {
			connidMap[data.connid] = true
			runConcurrency = uint64(len(connidMap))
		}
		requestTimeList = append(requestTimeList, data.requestTime)
		mutex.Unlock()
	}
	stop <- struct{}{}
	requestTime = uint64(time.Now().UnixNano()) - start
	resultStatistics(p.concurrency, processingTime, requestTime, maxTime, minTime, successNum, failureNum, runConcurrency, errCodeMap)
	fmt.Println("\n\n")
	fmt.Println("*************************  结果统计  ****************************")

	fmt.Println("请求总数:", successNum+failureNum, "总请求时间:",
		fmt.Sprintf("%.3f", float64(requestTime)/1e9),
		"秒", "successNum:", successNum, "failureNum:", failureNum)
	printTop(requestTimeList)
	fmt.Println("*************************  结果 end   ****************************")
	fmt.Printf("\n\n")

}

func resultStatistics(concurrent, processingTime, requestTime, maxTime, minTime, successNum, failureNum,
	runConcurrency uint64, errCodeMap *sync.Map) {
	if processingTime == 0 {
		processingTime = 1
	}
	var (
		// qps              float64
		averageTime      float64
		maxTimeFloat     float64
		minTimeFloat     float64
		requestTimeFloat float64
	)
	if processingTime != 0 {
		// qps = float64(successNum*1e9*concurrent) / float64(processingTime)
	}
	if successNum != 0 && concurrent != 0 {
		averageTime = float64(processingTime) / float64(successNum*1e6)
	}
	maxTimeFloat = float64(maxTime) / 1e6
	minTimeFloat = float64(minTime) / 1e6
	requestTimeFloat = float64(requestTime) / 1e9
	result := fmt.Sprintf("%6.0fs│%8d│%10d│%10d│%12.2f│%13.2f│%13.2f│%v",
		requestTimeFloat, runConcurrency, successNum, failureNum, maxTimeFloat, minTimeFloat, averageTime, printMap(errCodeMap))
	fmt.Println(result)
}

func printMap(errCodeMap *sync.Map) (mapStr string) {
	var mapArr []string
	errCodeMap.Range(func(key, value interface{}) bool {
		mapArr = append(mapArr, fmt.Sprintf("%v:%v", key, value))
		return true
	})
	sort.Strings(mapArr)
	mapStr = strings.Join(mapArr, ";")
	return
}

func printHeader() {
	fmt.Printf("\n\n")
	fmt.Println(" 耗时(s)│ 连接数│请求成功数│请求失败数│最长耗时(ms)│最短耗时(ms) │平均耗时(ms) │状态码")
	return
}

type uint64List []uint64

func (u64 uint64List) Len() int           { return len(u64) }
func (u64 uint64List) Swap(i, j int)      { u64[i], u64[j] = u64[j], u64[i] }
func (u64 uint64List) Less(i, j int) bool { return u64[i] < u64[j] }

func printTop(requestTimeList []uint64) {
	if requestTimeList == nil {
		return
	}
	all := uint64List{}
	all = requestTimeList
	sort.Sort(all)
	var numFirst, numSecond, numThree, numFour uint64
	for _, v := range all {
		switch v / 1e9 {
		case 0:
			numFirst++
		case 1:
			numSecond++
		case 2:
			numThree++
		default:
			numFour++
		}
	}
	fmt.Printf("0-1s内返回成功请求数(%d)\n", numFirst)
	fmt.Printf("1-2s内返回成功请求数(%d)\n", numSecond)
	fmt.Printf("2-3s内返回成功请求数(%d)\n", numThree)
	fmt.Printf("大于3s返回成功请求数(%d)\n", numFour)
}
