package rpcx

import (
	"context"
	"fmt"
	"ocpp16/config"
	"ocpp16/protocol"
	"ocpp16/websocket"
	"time"

	metrics "github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"github.com/smallnest/rpcx/share"
)

type ActiveCallServer struct {
	ChargingCore            *ChargingCoreServer
	SmartCharging           *SmartChargingServer
	LocalAuthListManagement *LocalAuthListManagementServer
	Reservation             *ReservationServer
	RemoteTrigger           *RemoteTriggerServer
	FirmwareManagement      *FirmwareManagementServer
}

type ChargingCoreServer struct {
	handler websocket.ActiveCallHandler
}

type SmartChargingServer struct {
	handler websocket.ActiveCallHandler
}

type LocalAuthListManagementServer struct {
	handler websocket.ActiveCallHandler
}
type ReservationServer struct {
	handler websocket.ActiveCallHandler
}

type RemoteTriggerServer struct {
	handler websocket.ActiveCallHandler
}

type FirmwareManagementServer struct {
	handler websocket.ActiveCallHandler
}

func NewActiveCallPlugin(handler websocket.ActiveCallHandler) {
	s := &ActiveCallServer{
		ChargingCore:            &ChargingCoreServer{handler: handler},
		SmartCharging:           &SmartChargingServer{handler: handler},
		LocalAuthListManagement: &LocalAuthListManagementServer{handler: handler},
		Reservation:             &ReservationServer{handler: handler},
		RemoteTrigger:           &RemoteTriggerServer{handler: handler},
		FirmwareManagement:      &FirmwareManagementServer{handler: handler},
	}
	go s.run()
}
func (s *ActiveCallServer) run() {
	conf := config.GCONF
	r := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: conf.RPCAddress,
		EtcdServers:    conf.ETCDList,
		BasePath:       conf.ETCDBasePath,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
	}
	err := r.Start()
	if err != nil {
		panic(err)
	}
	rpcxServer := server.NewServer()
	rpcxServer.Plugins.Add(r)
	rpcxServer.RegisterName("ChargingCoreServer", s.ChargingCore, "")
	rpcxServer.RegisterName("SmartChargingServer", s.SmartCharging, "")
	rpcxServer.RegisterName("LocalAuthListManagementServer", s.LocalAuthListManagement, "")
	rpcxServer.RegisterName("ReservationServer", s.Reservation, "")
	rpcxServer.RegisterName("RemoteTriggerServer", s.RemoteTrigger, "")
	rpcxServer.RegisterName("FirmwareManagementServer", s.FirmwareManagement, "")
	rpcxServer.Serve("tcp", conf.RPCAddress)
}

type Reply struct {
	Err error
}

//ChargingCore
func (o *ChargingCoreServer) ActiveChangeConfiguration(ctx context.Context, req *protocol.ChangeConfigurationRequest, res *Reply) error {
	if req == nil || res == nil {
		return fmt.Errorf("ActiveChangeConfiguration error: req  or res nil, req(%+v), res(%+v)", req, res)
	}
	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	var uniqueid, id string
	id, uniqueid = m["chargingPointIdentify"], m["messageId"]
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.ChangeConfigurationName,
		Request:       *req,
	}
	err := o.handler(ctx, id, &call)
	res.Err = err
	return err
}

func (o *ChargingCoreServer) ActiveDataTransfer(ctx context.Context, req *protocol.DataTransferRequest, res *Reply) error {

	if req == nil || res == nil {
		return fmt.Errorf("ActiveDataTransfer error: req  or res nil, req(%+v), res(%+v)", req, res)
	}
	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	var uniqueid, id string
	id, uniqueid = m["chargingPointIdentify"], m["messageId"]
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.DataTransferName,
		Request:       *req,
	}
	err := o.handler(ctx, id, &call)
	res.Err = err
	return err
}

func (o *ChargingCoreServer) ActiveRemoteStartTransaction(ctx context.Context, req *protocol.RemoteStartTransactionRequest, res *Reply) error {

	if req == nil || res == nil {
		return fmt.Errorf("ActiveRemoteStartTransaction error: req  or res nil, req(%+v), res(%+v)", req, res)
	}
	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	var uniqueid, id string
	id, uniqueid = m["chargingPointIdentify"], m["messageId"]
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.RemoteStartTransactionName,
		Request:       *req,
	}
	err := o.handler(ctx, id, &call)
	res.Err = err
	return err

}

func (o *ChargingCoreServer) ActiveRemoteStopTransaction(ctx context.Context, req *protocol.RemoteStopTransactionRequest, res *Reply) error {

	if req == nil || res == nil {
		return fmt.Errorf("ActiveRemoteStopTransaction error: req  or res nil, req(%+v), res(%+v)", req, res)
	}
	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	var uniqueid, id string
	id, uniqueid = m["chargingPointIdentify"], m["messageId"]
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.RemoteStopTransactionName,
		Request:       *req,
	}
	err := o.handler(ctx, id, &call)
	res.Err = err
	return err
}

func (o *ChargingCoreServer) ActiveUnlockConnector(ctx context.Context, req *protocol.UnlockConnectorRequest, res *Reply) error {

	if req == nil || res == nil {
		return fmt.Errorf("ActiveUnlockConnector error: req  or res nil, req(%+v), res(%+v)", req, res)
	}
	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	var uniqueid, id string
	id, uniqueid = m["chargingPointIdentify"], m["messageId"]
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.UnlockConnectorName,
		Request:       *req,
	}
	err := o.handler(ctx, id, &call)
	res.Err = err
	return err
}

func (o *ChargingCoreServer) ActiveReset(ctx context.Context, req *protocol.ResetRequest, res *Reply) error {
	if req == nil || res == nil {
		return fmt.Errorf("ActiveReset error: req  or res nil, req(%+v), res(%+v)", req, res)
	}
	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	var uniqueid, id string
	id, uniqueid = m["chargingPointIdentify"], m["messageId"]
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.ResetName,
		Request:       *req,
	}
	err := o.handler(ctx, id, &call)
	res.Err = err
	return err

}

func (o *ChargingCoreServer) ActiveGetConfiguration(ctx context.Context, req *protocol.GetConfigurationRequest, res *Reply) error {
	if req == nil || res == nil {
		return fmt.Errorf("ActiveGetConfiguration error: req  or res nil, req(%+v), res(%+v)", req, res)
	}
	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	var uniqueid, id string
	id, uniqueid = m["chargingPointIdentify"], m["messageId"]
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.GetConfigurationName,
		Request:       *req,
	}
	err := o.handler(ctx, id, &call)
	res.Err = err
	return err
}

func (o *ChargingCoreServer) ActiveChangeAvailability(ctx context.Context, req *protocol.ChangeAvailabilityRequest, res *Reply) error {
	if req == nil || res == nil {
		return fmt.Errorf("ActiveChangeAvailability error: req  or res nil, req(%+v), res(%+v)", req, res)
	}
	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	var uniqueid, id string
	id, uniqueid = m["chargingPointIdentify"], m["messageId"]
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.ChangeAvailabilityName,
		Request:       *req,
	}
	err := o.handler(ctx, id, &call)
	res.Err = err
	return err
}

func (o *ChargingCoreServer) ActiveClearCache(ctx context.Context, req *protocol.ClearCacheRequest, res *Reply) error {
	if req == nil || res == nil {
		return fmt.Errorf("ActiveClearCache error: req  or res nil, req(%+v), res(%+v)", req, res)
	}
	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	var uniqueid, id string
	id, uniqueid = m["chargingPointIdentify"], m["messageId"]
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.ClearCacheName,
		Request:       *req,
	}
	err := o.handler(ctx, id, &call)
	res.Err = err
	return err
}

//SmartCharging
func (s *SmartChargingServer) ActiveSetChargingProfile(ctx context.Context, req *protocol.SetChargingProfileRequest, res *Reply) error {
	if req == nil || res == nil {
		return fmt.Errorf("ActiveSetChargingProfile error: req  or res nil, req(%+v), res(%+v)", req, res)
	}
	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	var uniqueid, id string
	id, uniqueid = m["chargingPointIdentify"], m["messageId"]
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.SetChargingProfileName,
		Request:       *req,
	}
	err := s.handler(ctx, id, &call)
	res.Err = err
	return err
}

func (s *SmartChargingServer) ActiveGetCompositeSchedule(ctx context.Context, req *protocol.GetCompositeScheduleRequest, res *Reply) error {
	if req == nil || res == nil {
		return fmt.Errorf("ActiveGetCompositeSchedule error: req  or res nil, req(%+v), res(%+v)", req, res)
	}
	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	var uniqueid, id string
	id, uniqueid = m["chargingPointIdentify"], m["messageId"]
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.GetCompositeScheduleName,
		Request:       *req,
	}
	err := s.handler(ctx, id, &call)
	res.Err = err
	return err
}

func (s *SmartChargingServer) ActiveClearChargingProfile(ctx context.Context, req *protocol.ClearChargingProfileRequest, res *Reply) error {
	if req == nil || res == nil {
		return fmt.Errorf("ActiveClearChargingProfile error: req  or res nil, req(%+v), res(%+v)", req, res)
	}
	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	var uniqueid, id string
	id, uniqueid = m["chargingPointIdentify"], m["messageId"]
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.ClearChargingProfileName,
		Request:       *req,
	}
	err := s.handler(ctx, id, &call)
	res.Err = err
	return err
}

//Reservation
func (s *ReservationServer) ActiveCancelReservation(ctx context.Context, req *protocol.CancelReservationRequest, res *Reply) error {
	if req == nil || res == nil {
		return fmt.Errorf("ActiveCancelReservation error: req  or res nil, req(%+v), res(%+v)", req, res)
	}
	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	var uniqueid, id string
	id, uniqueid = m["chargingPointIdentify"], m["messageId"]
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.CancelReservationName,
		Request:       *req,
	}
	err := s.handler(ctx, id, &call)
	res.Err = err
	return err
}

func (s *ReservationServer) ActiveReserveNow(ctx context.Context, req *protocol.ReserveNowRequest, res *Reply) error {
	if req == nil || res == nil {
		return fmt.Errorf("ActiveReserveNow error: req  or res nil, req(%+v), res(%+v)", req, res)
	}
	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	var uniqueid, id string
	id, uniqueid = m["chargingPointIdentify"], m["messageId"]
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.ReserveNowName,
		Request:       *req,
	}
	err := s.handler(ctx, id, &call)
	res.Err = err
	return err
}

//LocalAuthListManagement

func (s *LocalAuthListManagementServer) ActiveGetLocalListVersion(ctx context.Context, req *protocol.GetLocalListVersionRequest, res *Reply) error {
	if req == nil || res == nil {
		return fmt.Errorf("ActiveGetLocalListVersion error: req  or res nil, req(%+v), res(%+v)", req, res)
	}
	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	var uniqueid, id string
	id, uniqueid = m["chargingPointIdentify"], m["messageId"]
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.GetLocalListVersionName,
		Request:       *req,
	}
	err := s.handler(ctx, id, &call)
	res.Err = err
	return err
}

func (s *LocalAuthListManagementServer) ActiveSendLocalList(ctx context.Context, req *protocol.SendLocalListRequest, res *Reply) error {
	if req == nil || res == nil {
		return fmt.Errorf("ActiveSendLocalList error: req  or res nil, req(%+v), res(%+v)", req, res)
	}
	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	var uniqueid, id string
	id, uniqueid = m["chargingPointIdentify"], m["messageId"]
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.SendLocalListName,
		Request:       *req,
	}
	err := s.handler(ctx, id, &call)
	res.Err = err
	return err
}

//RemoteTrigger
func (s *RemoteTriggerServer) ActiveTriggerMessage(ctx context.Context, req *protocol.TriggerMessageRequest, res *Reply) error {
	if req == nil || res == nil {
		return fmt.Errorf("ActiveTriggerMessage error: req  or res nil, req(%+v), res(%+v)", req, res)
	}
	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	var uniqueid, id string
	id, uniqueid = m["chargingPointIdentify"], m["messageId"]
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.TriggerMessageName,
		Request:       *req,
	}
	err := s.handler(ctx, id, &call)
	res.Err = err
	return err
}

//FirmwareManagement
func (s *FirmwareManagementServer) ActiveUpdateFirmware(ctx context.Context, req *protocol.UpdateFirmwareRequest, res *Reply) error {
	if req == nil || res == nil {
		return fmt.Errorf("ActiveUpdateFirmware error: req  or res nil, req(%+v), res(%+v)", req, res)
	}
	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	var uniqueid, id string
	id, uniqueid = m["chargingPointIdentify"], m["messageId"]
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.UpdateFirmwareName,
		Request:       *req,
	}
	err := s.handler(ctx, id, &call)
	res.Err = err
	return err
}

func (s *FirmwareManagementServer) ActiveGetDiagnostics(ctx context.Context, req *protocol.GetDiagnosticsRequest, res *Reply) error {
	if req == nil || res == nil {
		return fmt.Errorf("ActiveGetDiagnostics error: req  or res nil, req(%+v), res(%+v)", req, res)
	}
	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	var uniqueid, id string
	id, uniqueid = m["chargingPointIdentify"], m["messageId"]
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.GetDiagnosticsName,
		Request:       *req,
	}
	err := s.handler(ctx, id, &call)
	res.Err = err
	return err
}
