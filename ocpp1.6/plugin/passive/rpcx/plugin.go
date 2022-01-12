package rpcx

import (
	"context"
	"github.com/smallnest/rpcx/client"
	"ocpp16/config"
	"ocpp16/proto"
)

type RPCXPlugin struct {
	etcdAddr           []string
	basePath           string
	chargingCore       client.XClient
	smartCharging      client.XClient
	requestHandlerMap  map[string]proto.RequestHandler
	responseHandlerMap map[string]proto.ResponseHandler
}

func NewActionPlugin() *RPCXPlugin {
	conf := config.GCONF
	plugin := &RPCXPlugin{
		etcdAddr: conf.ETCDList,
		basePath: conf.ETCDBasePath,
	}
	plugin.init()
	plugin.registerRequestHandler()
	plugin.registerResponseHandler()
	return plugin
}

func (c *RPCXPlugin) init() {
	d1 := client.NewEtcdV3Discovery(c.basePath, "ChargingCoreClient", c.etcdAddr, nil)
	c.chargingCore = client.NewXClient("ChargingCoreClient", client.Failtry, client.RandomSelect, d1, client.DefaultOption)
	d2 := client.NewEtcdV3Discovery(c.basePath, "SmartChargingClient", c.etcdAddr, nil)
	c.smartCharging = client.NewXClient("SmartChargingClient", client.Failtry, client.RandomSelect, d2, client.DefaultOption)
}

func (c *RPCXPlugin) BootNotification(ctx context.Context, request proto.Request) (proto.Response, error) {
	reply := &proto.BootNotificationResponse{}
	err := c.chargingCore.Call(ctx, "BootNotification", request.(*proto.BootNotificationRequest), reply)
	return reply, err
}

func (c *RPCXPlugin) StatusNotification(ctx context.Context, request proto.Request) (proto.Response, error) {
	reply := &proto.StatusNotificationResponse{}
	err := c.chargingCore.Call(ctx, "StatusNotification", request.(*proto.StatusNotificationRequest), reply)
	return reply, err
}

func (c *RPCXPlugin) MeterValues(ctx context.Context, request proto.Request) (proto.Response, error) {
	reply := &proto.MeterValuesResponse{}
	err := c.chargingCore.Call(ctx, "MeterValues", request.(*proto.MeterValuesRequest), reply)
	return reply, err
}

func (c *RPCXPlugin) Authorize(ctx context.Context, request proto.Request) (proto.Response, error) {
	reply := &proto.AuthorizeResponse{}
	err := c.chargingCore.Call(ctx, "Authorize", request.(*proto.AuthorizeRequest), reply)
	return reply, err
}

func (c *RPCXPlugin) StartTransaction(ctx context.Context, request proto.Request) (proto.Response, error) {
	reply := &proto.StartTransactionResponse{}
	err := c.chargingCore.Call(ctx, "StartTransaction", request.(*proto.StartTransactionRequest), reply)
	return reply, err
}

func (c *RPCXPlugin) StopTransaction(ctx context.Context, request proto.Request) (proto.Response, error) {
	reply := &proto.StopTransactionResponse{}
	err := c.chargingCore.Call(ctx, "StopTransaction", request.(*proto.StopTransactionRequest), reply)
	return reply, err
}

// func (c *RPCXPlugin) ChangeConfiguration(ctx context.Context, request proto.Request) (proto.Response, error) {
// 	reply := &proto.ChangeConfigurationResponse{}
// 	err := c.chargingCore.Call(ctx, "ChangeConfiguration", request.(*proto.ChangeConfigurationRequest), reply)
// 	return reply, err
// }
// func (c *RPCXPlugin) DataTransfer(ctx context.Context, request proto.Request) (proto.Response, error) {
// 	reply := &proto.DataTransferResponse{}
// 	err := c.chargingCore.Call(ctx, "DataTransfer", request.(*proto.DataTransferRequest), reply)
// 	return reply, err
// }

// func (c *RPCXPlugin) SetChargingProfile(ctx context.Context, request proto.Request) (proto.Response, error) {
// 	reply := &proto.SetChargingProfileResponse{}
// 	err := c.chargingCore.Call(ctx, "SetChargingProfile", request.(*proto.SetChargingProfileRequest), reply)
// 	return reply, err
// }

// func (c *RPCXPlugin) RemoteStartTransaction(ctx context.Context, request proto.Request) (proto.Response, error) {
// 	reply := &proto.RemoteStartTransactionResponse{}
// 	err := c.chargingCore.Call(ctx, "RemoteStartTransaction", request.(*proto.RemoteStartTransactionRequest), reply)
// 	return reply, err
// }

// func (c *RPCXPlugin) RemoteStopTransaction(ctx context.Context, request proto.Request) (proto.Response, error) {
// 	reply := &proto.RemoteStopTransactionResponse{}
// 	err := c.chargingCore.Call(ctx, "RemoteStopTransaction", request.(*proto.RemoteStopTransactionRequest), reply)
// 	return reply, err
// }

// func (c *RPCXPlugin) Reset(ctx context.Context, request proto.Request) (proto.Response, error) {
// 	reply := &proto.ResetResponse{}
// 	err := c.chargingCore.Call(ctx, "Reset", request.(*proto.ResetRequest), reply)
// 	return reply, err
// }

// func (c *RPCXPlugin) UnlockConnector(ctx context.Context, request proto.Request) (proto.Response, error) {
// 	reply := &proto.UnlockConnectorResponse{}
// 	err := c.chargingCore.Call(ctx, "UnlockConnector", request.(*proto.UnlockConnectorRequest), reply)
// 	return reply, err
// }

func (c *RPCXPlugin) registerRequestHandler() {
	c.requestHandlerMap = map[string]proto.RequestHandler{
		proto.BootNotificationName:   proto.RequestHandler(c.BootNotification),
		proto.StatusNotificationName: proto.RequestHandler(c.StatusNotification),
		proto.MeterValuesName:        proto.RequestHandler(c.MeterValues),
		proto.AuthorizeName:          proto.RequestHandler(c.Authorize),
		proto.StartTransactionName:   proto.RequestHandler(c.StartTransaction),
		proto.StopTransactionName:    proto.RequestHandler(c.StopTransaction),
		// proto.ChangeConfigurationName:    proto.RequestHandler(c.ChangeConfiguration),
		// proto.DataTransferName:           proto.RequestHandler(c.DataTransfer),
		// proto.SetChargingProfileName:     proto.RequestHandler(c.SetChargingProfile),
		// proto.RemoteStartTransactionName: proto.RequestHandler(c.RemoteStartTransaction),
		// proto.RemoteStopTransactionName:  proto.RequestHandler(c.RemoteStopTransaction),
		// proto.ResetName:                  proto.RequestHandler(c.Reset),
		// proto.UnlockConnectorName:        proto.RequestHandler(c.UnlockConnector),
	}
}

//RequestHandler represent device active request Center
func (c *RPCXPlugin) RequestHandler(action string) (proto.RequestHandler, bool) {
	handler, ok := c.requestHandlerMap[action]
	return handler, ok
}

// func (c *Client) UnlockConnector(ctx context.Context, request proto.Request) (proto.Response, error) {
// 	reply := &proto.UnlockConnectorResponse{}
// 	err := c.chargingCore.Call(ctx, "UnlockConnector", request.(*proto.UnlockConnectorRequest), reply)
// 	return reply, err
// }
type Reply struct {
	Err error
}

func (c *RPCXPlugin) ChangeConfigurationResponse(ctx context.Context, res proto.Response) error {
	reply := &Reply{}
	err := c.chargingCore.Call(ctx, "ChangeConfigurationResponse", res.(*proto.ChangeConfigurationResponse), reply)
	return err
}

func (c *RPCXPlugin) DataTransferResponse(ctx context.Context, res proto.Response) error {
	reply := &Reply{}
	err := c.chargingCore.Call(ctx, "DataTransferResponse", res.(*proto.DataTransferResponse), reply)
	return err
}

func (c *RPCXPlugin) RemoteStartTransactionResponse(ctx context.Context, res proto.Response) error {
	reply := &Reply{}
	err := c.chargingCore.Call(ctx, "RemoteStartTransactionResponse", res.(*proto.RemoteStartTransactionResponse), reply)
	return err
}

func (c *RPCXPlugin) ResetResponse(ctx context.Context, res proto.Response) error {
	reply := &Reply{}
	err := c.chargingCore.Call(ctx, "ResetResponse", res.(*proto.ResetResponse), reply)
	return err
}

func (c *RPCXPlugin) RemoteStopTransactionResponse(ctx context.Context, res proto.Response) error {
	reply := &Reply{}
	err := c.chargingCore.Call(ctx, "RemoteStopTransactionResponse", res.(*proto.RemoteStopTransactionResponse), reply)
	return err
}

func (c *RPCXPlugin) UnlockConnectorResponse(ctx context.Context, res proto.Response) error {
	reply := &Reply{}
	err := c.chargingCore.Call(ctx, "UnlockConnectorResponse", res.(*proto.UnlockConnectorResponse), reply)
	return err
}

func (c *RPCXPlugin) CallError(ctx context.Context, res proto.Response) error {
	reply := &Reply{}
	err := c.chargingCore.Call(ctx, "CallError", res.(*proto.CallError), reply)
	return err
}

func (c *RPCXPlugin) SetChargingProfileResponse(ctx context.Context, res proto.Response) error {
	reply := &Reply{}
	err := c.smartCharging.Call(ctx, "SetChargingProfileResponse", res.(*proto.SetChargingProfileResponse), reply)
	return err
}

func (c *RPCXPlugin) registerResponseHandler() {
	c.responseHandlerMap = map[string]proto.ResponseHandler{
		proto.ChangeConfigurationName:    proto.ResponseHandler(c.ChangeConfigurationResponse),
		proto.DataTransferName:           proto.ResponseHandler(c.DataTransferResponse),
		proto.RemoteStartTransactionName: proto.ResponseHandler(c.RemoteStartTransactionResponse),
		proto.ResetName:                  proto.ResponseHandler(c.ResetResponse),
		proto.RemoteStopTransactionName:  proto.ResponseHandler(c.RemoteStopTransactionResponse),
		proto.UnlockConnectorName:        proto.ResponseHandler(c.UnlockConnectorResponse),
		proto.SetChargingProfileName:     proto.ResponseHandler(c.SetChargingProfileResponse),
		proto.CallErrorName:              proto.ResponseHandler(c.CallError),
	}
}

//ResponseHandler represent The device reply to the center request
func (c *RPCXPlugin) ResponseHandler(action string) (proto.ResponseHandler, bool) {
	handler, ok := c.responseHandlerMap[action]
	return handler, ok
}
