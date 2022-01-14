package rpcx

import (
	"context"
	"ocpp16/config"
	"ocpp16/protocol"

	"github.com/smallnest/rpcx/client"
)

type RPCXPlugin struct {
	etcdAddr           []string
	basePath           string
	chargingCore       client.XClient
	smartCharging      client.XClient
	requestHandlerMap  map[string]protocol.RequestHandler
	responseHandlerMap map[string]protocol.ResponseHandler
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

func (c *RPCXPlugin) BootNotification(ctx context.Context, request protocol.Request) (protocol.Response, error) {
	reply := &protocol.BootNotificationResponse{}
	err := c.chargingCore.Call(ctx, "BootNotification", request.(*protocol.BootNotificationRequest), reply)
	return reply, err
}

func (c *RPCXPlugin) StatusNotification(ctx context.Context, request protocol.Request) (protocol.Response, error) {
	reply := &protocol.StatusNotificationResponse{}
	err := c.chargingCore.Call(ctx, "StatusNotification", request.(*protocol.StatusNotificationRequest), reply)
	return reply, err
}

func (c *RPCXPlugin) MeterValues(ctx context.Context, request protocol.Request) (protocol.Response, error) {
	reply := &protocol.MeterValuesResponse{}
	err := c.chargingCore.Call(ctx, "MeterValues", request.(*protocol.MeterValuesRequest), reply)
	return reply, err
}

func (c *RPCXPlugin) Authorize(ctx context.Context, request protocol.Request) (protocol.Response, error) {
	reply := &protocol.AuthorizeResponse{}
	err := c.chargingCore.Call(ctx, "Authorize", request.(*protocol.AuthorizeRequest), reply)
	return reply, err
}

func (c *RPCXPlugin) StartTransaction(ctx context.Context, request protocol.Request) (protocol.Response, error) {
	reply := &protocol.StartTransactionResponse{}
	err := c.chargingCore.Call(ctx, "StartTransaction", request.(*protocol.StartTransactionRequest), reply)
	return reply, err
}

func (c *RPCXPlugin) StopTransaction(ctx context.Context, request protocol.Request) (protocol.Response, error) {
	reply := &protocol.StopTransactionResponse{}
	err := c.chargingCore.Call(ctx, "StopTransaction", request.(*protocol.StopTransactionRequest), reply)
	return reply, err
}

func (c *RPCXPlugin) ChargingPointOffline(id string) error {
	type OfflineNotice struct {
		ChargingPointIdentify string `json:"ChargingPointIdentify"`
	}
	req := &OfflineNotice{ChargingPointIdentify: id}
	reply := &Reply{}
	err := c.chargingCore.Call(context.TODO(), "ChargingPointOffline", req, reply)
	return err
}

// func (c *RPCXPlugin) ChangeConfiguration(ctx context.Context, request protocol.Request) (protocol.Response, error) {
// 	reply := &protocol.ChangeConfigurationResponse{}
// 	err := c.chargingCore.Call(ctx, "ChangeConfiguration", request.(*protocol.ChangeConfigurationRequest), reply)
// 	return reply, err
// }
// func (c *RPCXPlugin) DataTransfer(ctx context.Context, request protocol.Request) (protocol.Response, error) {
// 	reply := &protocol.DataTransferResponse{}
// 	err := c.chargingCore.Call(ctx, "DataTransfer", request.(*protocol.DataTransferRequest), reply)
// 	return reply, err
// }

// func (c *RPCXPlugin) SetChargingProfile(ctx context.Context, request protocol.Request) (protocol.Response, error) {
// 	reply := &protocol.SetChargingProfileResponse{}
// 	err := c.chargingCore.Call(ctx, "SetChargingProfile", request.(*protocol.SetChargingProfileRequest), reply)
// 	return reply, err
// }

// func (c *RPCXPlugin) RemoteStartTransaction(ctx context.Context, request protocol.Request) (protocol.Response, error) {
// 	reply := &protocol.RemoteStartTransactionResponse{}
// 	err := c.chargingCore.Call(ctx, "RemoteStartTransaction", request.(*protocol.RemoteStartTransactionRequest), reply)
// 	return reply, err
// }

// func (c *RPCXPlugin) RemoteStopTransaction(ctx context.Context, request protocol.Request) (protocol.Response, error) {
// 	reply := &protocol.RemoteStopTransactionResponse{}
// 	err := c.chargingCore.Call(ctx, "RemoteStopTransaction", request.(*protocol.RemoteStopTransactionRequest), reply)
// 	return reply, err
// }

// func (c *RPCXPlugin) Reset(ctx context.Context, request protocol.Request) (protocol.Response, error) {
// 	reply := &protocol.ResetResponse{}
// 	err := c.chargingCore.Call(ctx, "Reset", request.(*protocol.ResetRequest), reply)
// 	return reply, err
// }

// func (c *RPCXPlugin) UnlockConnector(ctx context.Context, request protocol.Request) (protocol.Response, error) {
// 	reply := &protocol.UnlockConnectorResponse{}
// 	err := c.chargingCore.Call(ctx, "UnlockConnector", request.(*protocol.UnlockConnectorRequest), reply)
// 	return reply, err
// }

func (c *RPCXPlugin) registerRequestHandler() {
	c.requestHandlerMap = map[string]protocol.RequestHandler{
		protocol.BootNotificationName:   protocol.RequestHandler(c.BootNotification),
		protocol.StatusNotificationName: protocol.RequestHandler(c.StatusNotification),
		protocol.MeterValuesName:        protocol.RequestHandler(c.MeterValues),
		protocol.AuthorizeName:          protocol.RequestHandler(c.Authorize),
		protocol.StartTransactionName:   protocol.RequestHandler(c.StartTransaction),
		protocol.StopTransactionName:    protocol.RequestHandler(c.StopTransaction),
		// protocol.ChangeConfigurationName:    protocol.RequestHandler(c.ChangeConfiguration),
		// protocol.DataTransferName:           protocol.RequestHandler(c.DataTransfer),
		// protocol.SetChargingProfileName:     protocol.RequestHandler(c.SetChargingProfile),
		// protocol.RemoteStartTransactionName: protocol.RequestHandler(c.RemoteStartTransaction),
		// protocol.RemoteStopTransactionName:  protocol.RequestHandler(c.RemoteStopTransaction),
		// protocol.ResetName:                  protocol.RequestHandler(c.Reset),
		// protocol.UnlockConnectorName:        protocol.RequestHandler(c.UnlockConnector),
	}
}

//RequestHandler represent device active request Center
func (c *RPCXPlugin) RequestHandler(action string) (protocol.RequestHandler, bool) {
	handler, ok := c.requestHandlerMap[action]
	return handler, ok
}

// func (c *Client) UnlockConnector(ctx context.Context, request protocol.Request) (protocol.Response, error) {
// 	reply := &protocol.UnlockConnectorResponse{}
// 	err := c.chargingCore.Call(ctx, "UnlockConnector", request.(*protocol.UnlockConnectorRequest), reply)
// 	return reply, err
// }
type Reply struct {
	Err error
}

func (c *RPCXPlugin) ChangeConfigurationResponse(ctx context.Context, res protocol.Response) error {
	reply := &Reply{}
	err := c.chargingCore.Call(ctx, "ChangeConfigurationResponse", res.(*protocol.ChangeConfigurationResponse), reply)
	return err
}

func (c *RPCXPlugin) DataTransferResponse(ctx context.Context, res protocol.Response) error {
	reply := &Reply{}
	err := c.chargingCore.Call(ctx, "DataTransferResponse", res.(*protocol.DataTransferResponse), reply)
	return err
}

func (c *RPCXPlugin) RemoteStartTransactionResponse(ctx context.Context, res protocol.Response) error {
	reply := &Reply{}
	err := c.chargingCore.Call(ctx, "RemoteStartTransactionResponse", res.(*protocol.RemoteStartTransactionResponse), reply)
	return err
}

func (c *RPCXPlugin) ResetResponse(ctx context.Context, res protocol.Response) error {
	reply := &Reply{}
	err := c.chargingCore.Call(ctx, "ResetResponse", res.(*protocol.ResetResponse), reply)
	return err
}

func (c *RPCXPlugin) RemoteStopTransactionResponse(ctx context.Context, res protocol.Response) error {
	reply := &Reply{}
	err := c.chargingCore.Call(ctx, "RemoteStopTransactionResponse", res.(*protocol.RemoteStopTransactionResponse), reply)
	return err
}

func (c *RPCXPlugin) UnlockConnectorResponse(ctx context.Context, res protocol.Response) error {
	reply := &Reply{}
	err := c.chargingCore.Call(ctx, "UnlockConnectorResponse", res.(*protocol.UnlockConnectorResponse), reply)
	return err
}

func (c *RPCXPlugin) SendLocalListResponse(ctx context.Context, res protocol.Response) error {
	reply := &Reply{}
	err := c.chargingCore.Call(ctx, "SendLocalListResponse", res.(*protocol.SendLocalListResponse), reply)
	return err
}

func (c *RPCXPlugin) GetLocalListVersionResponse(ctx context.Context, res protocol.Response) error {
	reply := &Reply{}
	err := c.chargingCore.Call(ctx, "GetLocalListVersionResponse", res.(*protocol.GetLocalListVersionResponse), reply)
	return err
}

func (c *RPCXPlugin) GetConfigurationResponse(ctx context.Context, res protocol.Response) error {
	reply := &Reply{}
	err := c.chargingCore.Call(ctx, "GetConfigurationResponse", res.(*protocol.GetConfigurationResponse), reply)
	return err
}

func (c *RPCXPlugin) CallError(ctx context.Context, res protocol.Response) error {
	reply := &Reply{}
	err := c.chargingCore.Call(ctx, "CallError", res.(*protocol.CallError), reply)
	return err
}

func (c *RPCXPlugin) SetChargingProfileResponse(ctx context.Context, res protocol.Response) error {
	reply := &Reply{}
	err := c.smartCharging.Call(ctx, "SetChargingProfileResponse", res.(*protocol.SetChargingProfileResponse), reply)
	return err
}

func (c *RPCXPlugin) registerResponseHandler() {
	c.responseHandlerMap = map[string]protocol.ResponseHandler{
		protocol.ChangeConfigurationName:    protocol.ResponseHandler(c.ChangeConfigurationResponse),
		protocol.DataTransferName:           protocol.ResponseHandler(c.DataTransferResponse),
		protocol.RemoteStartTransactionName: protocol.ResponseHandler(c.RemoteStartTransactionResponse),
		protocol.ResetName:                  protocol.ResponseHandler(c.ResetResponse),
		protocol.RemoteStopTransactionName:  protocol.ResponseHandler(c.RemoteStopTransactionResponse),
		protocol.UnlockConnectorName:        protocol.ResponseHandler(c.UnlockConnectorResponse),
		protocol.SendLocalListName:          protocol.ResponseHandler(c.SendLocalListResponse),
		protocol.GetLocalListVersionName:    protocol.ResponseHandler(c.GetLocalListVersionResponse),
		protocol.GetConfigurationName:       protocol.ResponseHandler(c.GetConfigurationResponse),
		protocol.SetChargingProfileName:     protocol.ResponseHandler(c.SetChargingProfileResponse),
		protocol.CallErrorName:              protocol.ResponseHandler(c.CallError),
	}
}

//ResponseHandler represent The device reply to the center request
func (c *RPCXPlugin) ResponseHandler(action string) (protocol.ResponseHandler, bool) {
	handler, ok := c.responseHandlerMap[action]
	return handler, ok
}
