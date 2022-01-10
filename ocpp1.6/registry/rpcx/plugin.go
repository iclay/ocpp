package rpcx

import (
	"context"
	"fmt"
	metrics "github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"github.com/smallnest/rpcx/share"
	"ocpp16/config"
	"ocpp16/proto"
	"ocpp16/websocket"
	"time"
)

type ActiveCallServer struct {
	ChargingCore  *ChargingCoreServer
	SmartCharging *SmartChargingServer
}

type ChargingCoreServer struct {
	handler websocket.ActiveCallHandler
}

type SmartChargingServer struct {
	handler websocket.ActiveCallHandler
}

func NewActiveCallPlugin(handler websocket.ActiveCallHandler) {
	s := &ActiveCallServer{
		ChargingCore:  &ChargingCoreServer{handler: handler},
		SmartCharging: &SmartChargingServer{handler: handler},
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
	rpcxServer.Serve("tcp", conf.RPCAddress)
}

type Reply struct {
	Err error
}

func (o *ChargingCoreServer) ActiveChangeConfiguration(ctx context.Context, req *proto.ChangeConfigurationRequest, res *Reply) error {
	if req == nil || res == nil {
		return fmt.Errorf("ActiveChangeConfiguration error: req  or res nil, req(%+v), res(%+v)", req, res)
	}
	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	var uniqueid, id string
	id, uniqueid = m["chargingPointIdentify"], m["messageId"]
	call := proto.Call{
		MessageTypeID: proto.CALL,
		UniqueID:      uniqueid,
		Action:        proto.ChangeConfigurationName,
		Request:       *req,
	}
	err := o.handler(ctx, id, &call)
	res.Err = err
	return err
}

func (o *ChargingCoreServer) ActiveDataTransfer(ctx context.Context, req *proto.DataTransferRequest, res *Reply) error {

	if req == nil || res == nil {
		return fmt.Errorf("ActiveDataTransfer error: req  or res nil, req(%+v), res(%+v)", req, res)
	}
	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	var uniqueid, id string
	id, uniqueid = m["chargingPointIdentify"], m["messageId"]
	call := proto.Call{
		MessageTypeID: proto.CALL,
		UniqueID:      uniqueid,
		Action:        proto.DataTransferName,
		Request:       *req,
	}
	err := o.handler(ctx, id, &call)
	res.Err = err
	return err
}

func (o *ChargingCoreServer) ActiveRemoteStartTransaction(ctx context.Context, req *proto.RemoteStartTransactionRequest, res *Reply) error {

	if req == nil || res == nil {
		return fmt.Errorf("ActiveRemoteStartTransaction error: req  or res nil, req(%+v), res(%+v)", req, res)
	}
	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	var uniqueid, id string
	id, uniqueid = m["chargingPointIdentify"], m["messageId"]
	call := proto.Call{
		MessageTypeID: proto.CALL,
		UniqueID:      uniqueid,
		Action:        proto.RemoteStartTransactionName,
		Request:       *req,
	}
	err := o.handler(ctx, id, &call)
	res.Err = err
	return err

}

func (o *ChargingCoreServer) ActiveRemoteStopTransaction(ctx context.Context, req *proto.RemoteStopTransactionRequest, res *Reply) error {

	if req == nil || res == nil {
		return fmt.Errorf("ActiveRemoteStopTransaction error: req  or res nil, req(%+v), res(%+v)", req, res)
	}
	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	var uniqueid, id string
	id, uniqueid = m["chargingPointIdentify"], m["messageId"]
	call := proto.Call{
		MessageTypeID: proto.CALL,
		UniqueID:      uniqueid,
		Action:        proto.RemoteStopTransactionName,
		Request:       *req,
	}
	err := o.handler(ctx, id, &call)
	res.Err = err
	return err
}

func (o *ChargingCoreServer) ActiveUnlockConnector(ctx context.Context, req *proto.UnlockConnectorRequest, res *Reply) error {

	if req == nil || res == nil {
		return fmt.Errorf("ActiveUnlockConnector error: req  or res nil, req(%+v), res(%+v)", req, res)
	}
	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	var uniqueid, id string
	id, uniqueid = m["chargingPointIdentify"], m["messageId"]
	call := proto.Call{
		MessageTypeID: proto.CALL,
		UniqueID:      uniqueid,
		Action:        proto.UnlockConnectorName,
		Request:       *req,
	}
	err := o.handler(ctx, id, &call)
	res.Err = err
	return err
}

func (o *ChargingCoreServer) ActiveReset(ctx context.Context, req *proto.ResetRequest, res *Reply) error {
	if req == nil || res == nil {
		return fmt.Errorf("ActiveReset error: req  or res nil, req(%+v), res(%+v)", req, res)
	}
	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	var uniqueid, id string
	id, uniqueid = m["chargingPointIdentify"], m["messageId"]
	call := proto.Call{
		MessageTypeID: proto.CALL,
		UniqueID:      uniqueid,
		Action:        proto.ResetName,
		Request:       *req,
	}
	err := o.handler(ctx, id, &call)
	res.Err = err
	return err

}

func (s *SmartChargingServer) ActiveSetChargingProfile(ctx context.Context, req *proto.SetChargingProfileRequest, res *Reply) error {
	if req == nil || res == nil {
		return fmt.Errorf("ActiveReset error: req  or res nil, req(%+v), res(%+v)", req, res)
	}
	m := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	var uniqueid, id string
	id, uniqueid = m["chargingPointIdentify"], m["messageId"]
	call := proto.Call{
		MessageTypeID: proto.CALL,
		UniqueID:      uniqueid,
		Action:        proto.SetChargingProfileName,
		Request:       *req,
	}
	err := s.handler(ctx, id, &call)
	res.Err = err
	return err
}
