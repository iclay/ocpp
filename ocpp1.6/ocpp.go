package main

import (
	"fmt"
	rpcx "ocpp16/plugin/rpcx"
	registry "ocpp16/registry/rpcx"

	// "ocpp16/plugin/local"
	"ocpp16/config"
	"ocpp16/logwriter"
	"ocpp16/websocket"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
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
	websocket.SetLogger(lg)
	server := websocket.NewDefaultServer()
	server.RegisterActionPlugin(rpcx.NewActionPlugin())
	server.RegisterActiveCallHandler(server.HandleActiveCall, registry.NewActiveCallPlugin)
	serveAddr, serveURI := fmt.Sprintf("%v:%v", conf.WebsocketAddr, conf.WebsocketPort), conf.WebsocketURI
	server.Serve(serveAddr, serveURI)
	return nil
}
