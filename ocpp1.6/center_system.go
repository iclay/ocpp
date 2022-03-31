package main

import (
	"fmt"
	"ocpp16/config"
	"ocpp16/logwriter"
	// active "ocpp16/plugin/active/local"
	// passive "ocpp16/plugin/passive/local"
	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
	active "ocpp16/plugin/active/rpcx"
	passive "ocpp16/plugin/passive/rpcx"
	ocpp16server "ocpp16/server"
	"os"
	"time"
)

var Version = "manual build has no version"

func main() {
	app := &cli.App{
		EnableBashCompletion: true,
		Name:                 "OCPP16",
		Usage:                "OCPP16 Protocol",
		Commands: []*cli.Command{
			{
				Name:   "serve",
				Usage:  "start the ocpp16 server",
				Action: serve,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "config",
						Usage:    "config file",
						Required: true,
						Aliases:  []string{"c"},
						EnvVars:  []string{"OCPP16_SERVER_CONFIG"},
					},
				},
			},
		},
		Authors: []*cli.Author{
			{
				Name:  "Tsinglink tech",
				Email: "16499111504li@gmail.com",
			},
		},
		Copyright: "Beijing Tsinglink Cloud Technology Co., Ltd (2021)",
		Version:   Version,
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func initLogger() *log.Logger {
	conf := config.GCONF
	lw := &logwriter.HourlySplit{
		Dir:           conf.LogPath,
		FileFormat:    "log_2006-01-02T15",
		MaxFileNumber: conf.LogMaxFileNum,
		MaxDiskUsage:  conf.LogMaxDiskUsage,
	}
	defer lw.Close()
	lg := log.New()
	customFormatter := &log.TextFormatter{
		TimestampFormat: time.RFC3339,
		FullTimestamp:   true,
	}
	lg.SetFormatter(customFormatter)
	lg.SetReportCaller(true)
	lg.SetOutput(lw)
	lv, err := log.ParseLevel(conf.LogLevel)
	if err != nil {
		lv = log.WarnLevel
	}
	lg.SetLevel(lv)
	return lg
}
func serve(c *cli.Context) error {
	config.ParseFile(c.String("config"))
	config.Print()
	conf := config.GCONF
	lg := initLogger()
	ocpp16server.SetLogger(lg)
	ocpp16server.WithOptions(ocpp16server.SupportCustomConversion(conf.UseConvert), ocpp16server.SupportObjectPool(conf.UsePool))
	server := ocpp16server.NewDefaultServer()
	defer server.Stop()
	actionPlugin := passive.NewActionPlugin()
	server.RegisterActionPlugin(actionPlugin)
	server.SetConnectHandlers(func(ws *ocpp16server.Wsconn) error {
		lg.Debugf("id(%s) connect,time(%s)", ws.ID(), time.Now().Format(time.RFC3339))
		return nil
	}, func(ws *ocpp16server.Wsconn) error {
		return actionPlugin.ChargingPointOnline(ws.ID())
	})
	server.SetDisconnetHandlers(func(ws *ocpp16server.Wsconn) error {
		lg.Debugf("id(%s) disconnect,time(%s)", ws.ID(), time.Now().Format(time.RFC3339))
		return nil
	}, func(ws *ocpp16server.Wsconn) error {
		return actionPlugin.ChargingPointOffline(ws.ID())
	})
	server.RegisterActiveCallHandler(server.HandleActiveCall, active.NewActiveCallPlugin)
	ServiceAddr, ServiceURI := conf.ServiceAddr, conf.ServiceURI
	if conf.WsEnable {
		wsAddr := fmt.Sprintf("%s:%d", ServiceAddr, conf.WsPort)
		server.Serve(wsAddr, ServiceURI)
	}
	if conf.WssEnable && conf.TLSCertificate != "" && conf.TLSCertificateKey != "" {
		wssAddr := fmt.Sprintf("%s:%d", ServiceAddr, conf.WssPort)
		server.ServeTLS(wssAddr, ServiceURI, conf.TLSCertificate, conf.TLSCertificateKey)
	}
	return nil
}
