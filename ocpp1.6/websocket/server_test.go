package websocket

import (
	"github.com/sirupsen/logrus"
	"ocpp16/config"
	"ocpp16/logwriter"
	"sync"
	"testing"
	"time"
)

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

func TestServer(t *testing.T) {
	config.GCONF = config.GConf{
		HeartbeatTimeout: 30,
	}
	wsEnable, wssEnable := false, true
	// wsEnable, wssEnable := true, false
	waitGroup := &sync.WaitGroup{}
	if wsEnable {
		waitGroup.Add(1)
		go WsHandler(t, waitGroup)
	}
	if wssEnable {
		waitGroup.Add(1)
		go WssHandler(t, waitGroup)
	}
	waitGroup.Wait()
}
