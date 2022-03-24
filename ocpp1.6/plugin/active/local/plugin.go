package local

import (
	"context"
	"fmt"
	"ocpp16/protocol"
	ocpp16server "ocpp16/server"
)

var activeCallHandler ActiveCallHandler

func NewActiveCallPlugin(handler ocpp16server.ActiveCallHandler) {
	activeCallHandler = ActiveCallHandler{
		handler: handler,
	}
}

type ActiveCallHandler struct {
	handler ocpp16server.ActiveCallHandler
}

func ActiveChangeConfiguration(ctx context.Context, id string, uniqueid string, req *protocol.ChangeConfigurationRequest) error {
	return activeCallHandler.activeChangeConfiguration(ctx, id, uniqueid, req)
}

func ActiveDataTransfer(ctx context.Context, id string, uniqueid string, req *protocol.DataTransferRequest) error {
	return activeCallHandler.activeDataTransfer(ctx, id, uniqueid, req)
}
func ActiveRemoteStartTransaction(ctx context.Context, id string, uniqueid string, req *protocol.RemoteStartTransactionRequest) error {
	return activeCallHandler.activeRemoteStartTransaction(ctx, id, uniqueid, req)
}

func ActiveRemoteStopTransaction(ctx context.Context, id string, uniqueid string, req *protocol.RemoteStopTransactionRequest) error {
	return activeCallHandler.activeRemoteStopTransaction(ctx, id, uniqueid, req)
}

func ActiveUnlockConnector(ctx context.Context, id string, uniqueid string, req *protocol.UnlockConnectorRequest) error {
	return activeCallHandler.activeUnlockConnector(ctx, id, uniqueid, req)

}

func ActiveReset(ctx context.Context, id string, uniqueid string, req *protocol.ResetRequest) error {
	return activeCallHandler.activeReset(ctx, id, uniqueid, req)

}
func ActiveGetConfiguration(ctx context.Context, id string, uniqueid string, req *protocol.GetConfigurationRequest) error {
	return activeCallHandler.activeGetConfiguration(ctx, id, uniqueid, req)
}

func ActiveSetChargingProfile(ctx context.Context, id string, uniqueid string, req *protocol.SetChargingProfileRequest) error {
	return activeCallHandler.activeSetChargingProfile(ctx, id, uniqueid, req)

}

func ActiveGetLocalListVersion(ctx context.Context, id string, uniqueid string, req *protocol.GetLocalListVersionRequest) error {
	return activeCallHandler.activeGetLocalListVersion(ctx, id, uniqueid, req)
}

func ActiveSendLocalList(ctx context.Context, id string, uniqueid string, req *protocol.SendLocalListRequest) error {
	return activeCallHandler.activeSendLocalList(ctx, id, uniqueid, req)
}

func ActiveChangeAvailability(ctx context.Context, id string, uniqueid string, req *protocol.ChangeAvailabilityRequest) error {
	return activeCallHandler.activeChangeAvailability(ctx, id, uniqueid, req)
}
func ActiveClearCache(ctx context.Context, id string, uniqueid string, req *protocol.ClearCacheRequest) error {
	return activeCallHandler.activeClearCache(ctx, id, uniqueid, req)
}

func ActiveGetCompositeSchedule(ctx context.Context, id string, uniqueid string, req *protocol.GetCompositeScheduleRequest) error {
	return activeCallHandler.activeGetCompositeSchedule(ctx, id, uniqueid, req)
}
func ActiveClearChargingProfile(ctx context.Context, id string, uniqueid string, req *protocol.ClearChargingProfileRequest) error {
	return activeCallHandler.activeClearChargingProfile(ctx, id, uniqueid, req)
}
func ActiveCancelReservation(ctx context.Context, id string, uniqueid string, req *protocol.CancelReservationRequest) error {
	return activeCallHandler.activeCancelReservation(ctx, id, uniqueid, req)

}

//ChargingCore
func (s *ActiveCallHandler) activeChangeConfiguration(ctx context.Context, id string, uniqueid string, req *protocol.ChangeConfigurationRequest) error {
	if req == nil {
		return fmt.Errorf("ActiveChangeConfiguration error: req nil, req(%+v)", req)
	}
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.ChangeConfigurationName,
		Request:       *req,
	}
	return s.handler(ctx, id, &call)
}

func (s *ActiveCallHandler) activeDataTransfer(ctx context.Context, id string, uniqueid string, req *protocol.DataTransferRequest) error {

	if req == nil {
		return fmt.Errorf("ActiveDataTransfer error: req nil, req(%+v)", req)
	}
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.DataTransferName,
		Request:       *req,
	}

	return s.handler(ctx, id, &call)

}

func (s *ActiveCallHandler) activeRemoteStartTransaction(ctx context.Context, id string, uniqueid string, req *protocol.RemoteStartTransactionRequest) error {

	if req == nil {
		return fmt.Errorf("ActiveRemoteStartTransaction error: req nil, req(%+v)", req)
	}
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.RemoteStartTransactionName,
		Request:       *req,
	}
	return s.handler(ctx, id, &call)
}

func (s *ActiveCallHandler) activeRemoteStopTransaction(ctx context.Context, id string, uniqueid string, req *protocol.RemoteStopTransactionRequest) error {

	if req == nil {
		return fmt.Errorf("ActiveRemoteStopTransaction error: req nil, req(%+v)", req)
	}
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.RemoteStopTransactionName,
		Request:       *req,
	}
	return s.handler(ctx, id, &call)
}

func (s *ActiveCallHandler) activeUnlockConnector(ctx context.Context, id string, uniqueid string, req *protocol.UnlockConnectorRequest) error {

	if req == nil {
		return fmt.Errorf("ActiveUnlockConnector error: req  nil, req(%+v)", req)
	}
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.UnlockConnectorName,
		Request:       *req,
	}
	return s.handler(ctx, id, &call)
}

func (s *ActiveCallHandler) activeReset(ctx context.Context, id string, uniqueid string, req *protocol.ResetRequest) error {
	if req == nil {
		return fmt.Errorf("ActiveReset error: req nil, req(%+v)", req)
	}
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.ResetName,
		Request:       *req,
	}
	return s.handler(ctx, id, &call)
}

func (s *ActiveCallHandler) activeGetConfiguration(ctx context.Context, id string, uniqueid string, req *protocol.GetConfigurationRequest) error {
	if req == nil {
		return fmt.Errorf("ActiveGetConfiguration error: req nil, req(%+v)", req)
	}
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.GetConfigurationName,
		Request:       *req,
	}
	return s.handler(ctx, id, &call)
}

func (s *ActiveCallHandler) activeChangeAvailability(ctx context.Context, id string, uniqueid string, req *protocol.ChangeAvailabilityRequest) error {
	if req == nil {
		return fmt.Errorf("ActiveChangeAvailability error: req nil, req(%+v)", req)
	}
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.ChangeAvailabilityName,
		Request:       *req,
	}
	return s.handler(ctx, id, &call)
}

func (s *ActiveCallHandler) activeClearCache(ctx context.Context, id string, uniqueid string, req *protocol.ClearCacheRequest) error {
	if req == nil {
		return fmt.Errorf("ActiveClearCache error: req nil, req(%+v)", req)
	}
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.ClearCacheName,
		Request:       *req,
	}
	return s.handler(ctx, id, &call)
}

//SmartCharging
func (s *ActiveCallHandler) activeSetChargingProfile(ctx context.Context, id string, uniqueid string, req *protocol.SetChargingProfileRequest) error {
	if req == nil {
		return fmt.Errorf("ActiveSetChargingProfile error: req nil, req(%+v)", req)
	}
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.SetChargingProfileName,
		Request:       *req,
	}
	return s.handler(ctx, id, &call)
}

func (s *ActiveCallHandler) activeGetCompositeSchedule(ctx context.Context, id string, uniqueid string, req *protocol.GetCompositeScheduleRequest) error {
	if req == nil {
		return fmt.Errorf("ActiveGetCompositeSchedule error: req nil, req(%+v)", req)
	}
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.GetCompositeScheduleName,
		Request:       *req,
	}
	return s.handler(ctx, id, &call)
}

func (s *ActiveCallHandler) activeClearChargingProfile(ctx context.Context, id string, uniqueid string, req *protocol.ClearChargingProfileRequest) error {
	if req == nil {
		return fmt.Errorf("ActiveGetCompositeSchedule error: req nil, req(%+v)", req)
	}
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.ClearChargingProfileName,
		Request:       *req,
	}
	return s.handler(ctx, id, &call)
}

//Reservation

func (s *ActiveCallHandler) activeCancelReservation(ctx context.Context, id string, uniqueid string, req *protocol.CancelReservationRequest) error {
	if req == nil {
		return fmt.Errorf("ActiveCancelReservation error: req nil, req(%+v)", req)
	}
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.CancelReservationName,
		Request:       *req,
	}
	return s.handler(ctx, id, &call)
}

//LocalAuthListManagement

func (s *ActiveCallHandler) activeGetLocalListVersion(ctx context.Context, id string, uniqueid string, req *protocol.GetLocalListVersionRequest) error {
	if req == nil {
		return fmt.Errorf("ActiveGetLocalListVersion error: req nil, req(%+v)", req)
	}
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.GetLocalListVersionName,
		Request:       *req,
	}
	return s.handler(ctx, id, &call)
}

func (s *ActiveCallHandler) activeSendLocalList(ctx context.Context, id string, uniqueid string, req *protocol.SendLocalListRequest) error {
	if req == nil {
		return fmt.Errorf("ActiveSendLocalList error: req nil, req(%+v)", req)
	}
	call := protocol.Call{
		MessageTypeID: protocol.CALL,
		UniqueID:      uniqueid,
		Action:        protocol.SendLocalListName,
		Request:       *req,
	}
	return s.handler(ctx, id, &call)
}
