package rpcx

import (
	"context"
	"ocpp16/proto"

	"github.com/smallnest/rpcx/client"
)

type Client struct {
	ETCDAddr   []string
	BasePath   string
	ocppclient client.XClient
}

func NewRPCXClient(ETCDList []string, BasePath string) *Client {
	client := &Client{
		ETCDAddr: ETCDList,
		BasePath: BasePath,
	}
	client.Init()
	return client
}

func (c *Client) Init() {
	d1 := client.NewEtcdV3Discovery(c.BasePath, "ocpp", c.ETCDAddr, nil)
	c.ocppclient = client.NewXClient("ocpp", client.Failtry, client.RandomSelect, d1, client.DefaultOption)
}

func (c *Client) BootNotification(ctx context.Context, request proto.Request) (proto.Response, error) {
	reply := &proto.BootNotificationResponse{}
	err := c.ocppclient.Call(ctx, "BootNotification", request.(*proto.BootNotificationRequest), reply)
	return reply, err
}

func (c *Client) RegisterOCPPHandler() map[string]proto.RequestHandler {
	return map[string]proto.RequestHandler{
		proto.BootNotificationName: proto.RequestHandler(c.BootNotification),
	}
}
