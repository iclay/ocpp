package rpcx

import (
	"context"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/share"
	"ocpp16/config"
	"ocpp16/protocol"
	"time"
)

type RPCXPlugin struct {
	etcdAddr                []string
	basePath                string
	chargingCore            client.XClient
	smartCharging           client.XClient
	firmwareManagement      client.XClient
	Reservation             client.XClient
	RemoteTrigger           client.XClient
	LocalAuthListManagement client.XClient
	requestHandlerMap       map[string]protocol.RequestHandler
	responseHandlerMap      map[string]protocol.ResponseHandler
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
	d3 := client.NewEtcdV3Discovery(c.basePath, "FirmwareManagementClient", c.etcdAddr, nil)
	c.firmwareManagement = client.NewXClient("FirmwareManagementClient", client.Failtry, client.RandomSelect, d3, client.DefaultOption)
	d4 := client.NewEtcdV3Discovery(c.basePath, "ReservationClient", c.etcdAddr, nil)
	c.Reservation = client.NewXClient("ReservationClient", client.Failtry, client.RandomSelect, d4, client.DefaultOption)
	d5 := client.NewEtcdV3Discovery(c.basePath, "RemoteTriggerClient", c.etcdAddr, nil)
	c.RemoteTrigger = client.NewXClient("RemoteTriggerClient", client.Failtry, client.RandomSelect, d5, client.DefaultOption)
	d6 := client.NewEtcdV3Discovery(c.basePath, "LocalAuthListManagementClient", c.etcdAddr, nil)
	c.LocalAuthListManagement = client.NewXClient("LocalAuthListManagementClient", client.Failtry, client.RandomSelect, d6, client.DefaultOption)
}
func (c *RPCXPlugin) Heartbeat(ctx context.Context, id string, uniqueid string, request protocol.Request) (protocol.Response, error) {
	reply := &protocol.HeartbeatResponse{
		CurrentTime: time.Now().Format(protocol.ISO8601),
	}
	return reply, nil
}

// chargingCore - request
func (c *RPCXPlugin) BootNotification(ctx context.Context, id string, uniqueid string, request protocol.Request) (protocol.Response, error) {
	reply := &protocol.BootNotificationResponse{}
	ctx = context.WithValue(ctx, share.ReqMetaDataKey, map[string]string{
		"chargingPointIdentify": id,
		"messageId":             uniqueid,
	})
	err := c.chargingCore.Call(ctx, "BootNotification", request.(*protocol.BootNotificationRequest), reply)
	return reply, err
}

func (c *RPCXPlugin) StatusNotification(ctx context.Context, id string, uniqueid string, request protocol.Request) (protocol.Response, error) {
	reply := &protocol.StatusNotificationResponse{}
	ctx = context.WithValue(ctx, share.ReqMetaDataKey, map[string]string{
		"chargingPointIdentify": id,
		"messageId":             uniqueid,
	})
	err := c.chargingCore.Call(ctx, "StatusNotification", request.(*protocol.StatusNotificationRequest), reply)
	return reply, err
}

func (c *RPCXPlugin) MeterValues(ctx context.Context, id string, uniqueid string, request protocol.Request) (protocol.Response, error) {
	reply := &protocol.MeterValuesResponse{}
	ctx = context.WithValue(ctx, share.ReqMetaDataKey, map[string]string{
		"chargingPointIdentify": id,
		"messageId":             uniqueid,
	})
	err := c.chargingCore.Call(ctx, "MeterValues", request.(*protocol.MeterValuesRequest), reply)
	return reply, err
}

func (c *RPCXPlugin) Authorize(ctx context.Context, id string, uniqueid string, request protocol.Request) (protocol.Response, error) {
	reply := &protocol.AuthorizeResponse{}
	ctx = context.WithValue(ctx, share.ReqMetaDataKey, map[string]string{
		"chargingPointIdentify": id,
		"messageId":             uniqueid,
	})
	err := c.chargingCore.Call(ctx, "Authorize", request.(*protocol.AuthorizeRequest), reply)
	return reply, err
}

func (c *RPCXPlugin) StartTransaction(ctx context.Context, id string, uniqueid string, request protocol.Request) (protocol.Response, error) {
	reply := &protocol.StartTransactionResponse{}
	ctx = context.WithValue(ctx, share.ReqMetaDataKey, map[string]string{
		"chargingPointIdentify": id,
		"messageId":             uniqueid,
	})
	err := c.chargingCore.Call(ctx, "StartTransaction", request.(*protocol.StartTransactionRequest), reply)
	return reply, err
}

func (c *RPCXPlugin) StopTransaction(ctx context.Context, id string, uniqueid string, request protocol.Request) (protocol.Response, error) {
	reply := &protocol.StopTransactionResponse{}
	ctx = context.WithValue(ctx, share.ReqMetaDataKey, map[string]string{
		"chargingPointIdentify": id,
		"messageId":             uniqueid,
	})
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

// firmwareManagement - request
func (c *RPCXPlugin) FirmwareStatusNotification(ctx context.Context, id string, uniqueid string, request protocol.Request) (protocol.Response, error) {
	reply := &protocol.FirmwareStatusNotificationResponse{}
	ctx = context.WithValue(ctx, share.ReqMetaDataKey, map[string]string{
		"chargingPointIdentify": id,
		"messageId":             uniqueid,
	})
	err := c.firmwareManagement.Call(ctx, "FirmwareStatusNotification", request.(*protocol.FirmwareStatusNotificationRequest), reply)
	return reply, err
}

func (c *RPCXPlugin) DiagnosticsStatusNotification(ctx context.Context, id string, uniqueid string, request protocol.Request) (protocol.Response, error) {
	reply := &protocol.DiagnosticsStatusNotificationResponse{}
	ctx = context.WithValue(ctx, share.ReqMetaDataKey, map[string]string{
		"chargingPointIdentify": id,
		"messageId":             uniqueid,
	})
	err := c.firmwareManagement.Call(ctx, "DiagnosticsStatusNotification", request.(*protocol.DiagnosticsStatusNotificationRequest), reply)
	return reply, err
}

func (c *RPCXPlugin) registerRequestHandler() {
	c.requestHandlerMap = map[string]protocol.RequestHandler{
		protocol.BootNotificationName:           protocol.RequestHandler(c.BootNotification),
		protocol.StatusNotificationName:         protocol.RequestHandler(c.StatusNotification),
		protocol.MeterValuesName:                protocol.RequestHandler(c.MeterValues),
		protocol.AuthorizeName:                  protocol.RequestHandler(c.Authorize),
		protocol.StartTransactionName:           protocol.RequestHandler(c.StartTransaction),
		protocol.StopTransactionName:            protocol.RequestHandler(c.StopTransaction),
		protocol.FirmwareStatusNotificationName: protocol.RequestHandler(c.FirmwareStatusNotification),
		protocol.HeartbeatName:                  protocol.RequestHandler(c.Heartbeat),
	}
}

//RequestHandler represent device active request Center
func (c *RPCXPlugin) RequestHandler(action string) (protocol.RequestHandler, bool) {
	handler, ok := c.requestHandlerMap[action]
	return handler, ok
}

type Reply struct {
	Err error
}

// chargingCore-response
func (c *RPCXPlugin) ChangeConfigurationResponse(ctx context.Context, id string, uniqueid string, res protocol.Response) error {
	reply := &Reply{}
	ctx = context.WithValue(ctx, share.ReqMetaDataKey, map[string]string{
		"chargingPointIdentify": id,
		"messageId":             uniqueid,
	})
	err := c.chargingCore.Call(ctx, "ChangeConfigurationResponse", res.(*protocol.ChangeConfigurationResponse), reply)
	return err
}

func (c *RPCXPlugin) DataTransferResponse(ctx context.Context, id string, uniqueid string, res protocol.Response) error {
	reply := &Reply{}
	ctx = context.WithValue(ctx, share.ReqMetaDataKey, map[string]string{
		"chargingPointIdentify": id,
		"messageId":             uniqueid,
	})
	err := c.chargingCore.Call(ctx, "DataTransferResponse", res.(*protocol.DataTransferResponse), reply)
	return err
}

func (c *RPCXPlugin) RemoteStartTransactionResponse(ctx context.Context, id string, uniqueid string, res protocol.Response) error {
	reply := &Reply{}
	ctx = context.WithValue(ctx, share.ReqMetaDataKey, map[string]string{
		"chargingPointIdentify": id,
		"messageId":             uniqueid,
	})
	err := c.chargingCore.Call(ctx, "RemoteStartTransactionResponse", res.(*protocol.RemoteStartTransactionResponse), reply)
	return err
}

func (c *RPCXPlugin) ResetResponse(ctx context.Context, id string, uniqueid string, res protocol.Response) error {
	reply := &Reply{}
	ctx = context.WithValue(ctx, share.ReqMetaDataKey, map[string]string{
		"chargingPointIdentify": id,
		"messageId":             uniqueid,
	})
	err := c.chargingCore.Call(ctx, "ResetResponse", res.(*protocol.ResetResponse), reply)
	return err
}

func (c *RPCXPlugin) RemoteStopTransactionResponse(ctx context.Context, id string, uniqueid string, res protocol.Response) error {
	reply := &Reply{}
	ctx = context.WithValue(ctx, share.ReqMetaDataKey, map[string]string{
		"chargingPointIdentify": id,
		"messageId":             uniqueid,
	})
	err := c.chargingCore.Call(ctx, "RemoteStopTransactionResponse", res.(*protocol.RemoteStopTransactionResponse), reply)
	return err
}

func (c *RPCXPlugin) UnlockConnectorResponse(ctx context.Context, id string, uniqueid string, res protocol.Response) error {
	reply := &Reply{}
	ctx = context.WithValue(ctx, share.ReqMetaDataKey, map[string]string{
		"chargingPointIdentify": id,
		"messageId":             uniqueid,
	})
	err := c.chargingCore.Call(ctx, "UnlockConnectorResponse", res.(*protocol.UnlockConnectorResponse), reply)
	return err
}

func (c *RPCXPlugin) GetConfigurationResponse(ctx context.Context, id string, uniqueid string, res protocol.Response) error {
	reply := &Reply{}
	ctx = context.WithValue(ctx, share.ReqMetaDataKey, map[string]string{
		"chargingPointIdentify": id,
		"messageId":             uniqueid,
	})
	err := c.chargingCore.Call(ctx, "GetConfigurationResponse", res.(*protocol.GetConfigurationResponse), reply)
	return err
}

func (c *RPCXPlugin) CallError(ctx context.Context, id string, uniqueid string, res protocol.Response) error {
	reply := &Reply{}
	err := c.chargingCore.Call(ctx, "CallError", res.(*protocol.CallError), reply)
	return err
}

// smartCharging - repsonse
func (c *RPCXPlugin) SetChargingProfileResponse(ctx context.Context, id string, uniqueid string, res protocol.Response) error {
	reply := &Reply{}
	ctx = context.WithValue(ctx, share.ReqMetaDataKey, map[string]string{
		"chargingPointIdentify": id,
		"messageId":             uniqueid,
	})
	err := c.smartCharging.Call(ctx, "SetChargingProfileResponse", res.(*protocol.SetChargingProfileResponse), reply)
	return err
}

func (c *RPCXPlugin) ClearChargingProfileResponse(ctx context.Context, id string, uniqueid string, res protocol.Response) error {
	reply := &Reply{}
	ctx = context.WithValue(ctx, share.ReqMetaDataKey, map[string]string{
		"chargingPointIdentify": id,
		"messageId":             uniqueid,
	})
	err := c.smartCharging.Call(ctx, "ClearChargingProfileResponse", res.(*protocol.ClearChargingProfileResponse), reply)
	return err
}

func (c *RPCXPlugin) GetCompositeScheduleResponse(ctx context.Context, id string, uniqueid string, res protocol.Response) error {
	reply := &Reply{}
	ctx = context.WithValue(ctx, share.ReqMetaDataKey, map[string]string{
		"chargingPointIdentify": id,
		"messageId":             uniqueid,
	})
	err := c.smartCharging.Call(ctx, "GetCompositeScheduleResponse", res.(*protocol.GetCompositeScheduleResponse), reply)
	return err
}

// firmwareManagement - response
func (c *RPCXPlugin) GetDiagnosticsResponse(ctx context.Context, id string, uniqueid string, res protocol.Response) error {
	reply := &Reply{}
	ctx = context.WithValue(ctx, share.ReqMetaDataKey, map[string]string{
		"chargingPointIdentify": id,
		"messageId":             uniqueid,
	})
	err := c.firmwareManagement.Call(ctx, "GetDiagnosticsResponse", res.(*protocol.GetDiagnosticsResponse), reply)
	return err
}

func (c *RPCXPlugin) UpdateFirmWareResponse(ctx context.Context, id string, uniqueid string, res protocol.Response) error {
	reply := &Reply{}
	ctx = context.WithValue(ctx, share.ReqMetaDataKey, map[string]string{
		"chargingPointIdentify": id,
		"messageId":             uniqueid,
	})
	err := c.firmwareManagement.Call(ctx, "UpdateFirmWareResponse", res.(*protocol.UpdateFirmwareResponse), reply)
	return err
}

//Reservation - response

func (c *RPCXPlugin) ReserveNowResponse(ctx context.Context, id string, uniqueid string, res protocol.Response) error {
	reply := &Reply{}
	ctx = context.WithValue(ctx, share.ReqMetaDataKey, map[string]string{
		"chargingPointIdentify": id,
		"messageId":             uniqueid,
	})
	err := c.Reservation.Call(ctx, "ReserveNowResponse", res.(*protocol.ReserveNowResponse), reply)
	return err
}

func (c *RPCXPlugin) CancelReservationResponse(ctx context.Context, id string, uniqueid string, res protocol.Response) error {
	reply := &Reply{}
	ctx = context.WithValue(ctx, share.ReqMetaDataKey, map[string]string{
		"chargingPointIdentify": id,
		"messageId":             uniqueid,
	})
	err := c.Reservation.Call(ctx, "CancelReservationResponse", res.(*protocol.CancelReservationResponse), reply)
	return err
}

//RemoteTrigger -response
func (c *RPCXPlugin) TriggerMessageResponse(ctx context.Context, id string, uniqueid string, res protocol.Response) error {
	reply := &Reply{}
	ctx = context.WithValue(ctx, share.ReqMetaDataKey, map[string]string{
		"chargingPointIdentify": id,
		"messageId":             uniqueid,
	})
	err := c.RemoteTrigger.Call(ctx, "TriggerMessageResponse", res.(*protocol.TriggerMessageResponse), reply)
	return err
}

//LocalAuthListManagement -response
func (c *RPCXPlugin) SendLocalListResponse(ctx context.Context, id string, uniqueid string, res protocol.Response) error {
	reply := &Reply{}
	ctx = context.WithValue(ctx, share.ReqMetaDataKey, map[string]string{
		"chargingPointIdentify": id,
		"messageId":             uniqueid,
	})
	err := c.LocalAuthListManagement.Call(ctx, "SendLocalListResponse", res.(*protocol.SendLocalListResponse), reply)
	return err
}

func (c *RPCXPlugin) GetLocalListVersionResponse(ctx context.Context, id string, uniqueid string, res protocol.Response) error {
	reply := &Reply{}
	ctx = context.WithValue(ctx, share.ReqMetaDataKey, map[string]string{
		"chargingPointIdentify": id,
		"messageId":             uniqueid,
	})
	err := c.LocalAuthListManagement.Call(ctx, "GetLocalListVersionResponse", res.(*protocol.GetLocalListVersionResponse), reply)
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
		protocol.GetConfigurationName:       protocol.ResponseHandler(c.GetConfigurationResponse),
		protocol.SetChargingProfileName:     protocol.ResponseHandler(c.SetChargingProfileResponse),
		protocol.ClearChargingProfileName:   protocol.ResponseHandler(c.ClearChargingProfileResponse),
		protocol.GetCompositeScheduleName:   protocol.ResponseHandler(c.GetCompositeScheduleResponse),
		protocol.ReserveNowName:             protocol.ResponseHandler(c.ReserveNowResponse),
		protocol.CancelReservationName:      protocol.ResponseHandler(c.CancelReservationResponse),
		protocol.TriggerMessageName:         protocol.ResponseHandler(c.TriggerMessageResponse),
		protocol.SendLocalListName:          protocol.ResponseHandler(c.SendLocalListResponse),
		protocol.GetLocalListVersionName:    protocol.ResponseHandler(c.GetLocalListVersionResponse),
		protocol.GetDiagnosticsName:         protocol.ResponseHandler(c.GetDiagnosticsResponse),
		protocol.UpdateFirmwareName:         protocol.ResponseHandler(c.UpdateFirmWareResponse),
		protocol.CallErrorName:              protocol.ResponseHandler(c.CallError),
	}
}

//ResponseHandler represent The device reply to the center request
func (c *RPCXPlugin) ResponseHandler(action string) (protocol.ResponseHandler, bool) {
	handler, ok := c.responseHandlerMap[action]
	return handler, ok
}
