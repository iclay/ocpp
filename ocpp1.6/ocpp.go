package main

import (
	"fmt"
	// "ocpp16/plugin/rpcx"
	"ocpp16/plugin/local"
	"ocpp16/websocket"
	"os"

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
				Name:  "lihuaye",
				Email: "16499111504li@gmail.com",
			},
		},
		Copyright: "Beijing Tsinglink Cloud Technology Co., Ltd (2021)",
		Version:   Version,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// func serve(c *cli.Context) error {
// 	server := websocket.NewDefaultServer()
// 	client := rpcx.NewRPCXClient([]string{}, "")
// 	server.RegisterOCPPHandler(client)
// 	server.Serve("127.0.0.1:8090", "/ocpp/:name/:id")
// 	return nil
// }

func serve(c *cli.Context) error {
	server := websocket.NewDefaultServer()
	plugin := local.NewLocalService()
	server.RegisterOCPPHandler(plugin)
	server.Serve("127.0.0.1:8090", "/ocpp/:name/:id")
	return nil
}