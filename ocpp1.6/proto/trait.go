package proto

import (
	"context"
	"reflect"
	"sync"
)

func init() {
	OCPP16M.register(
		BootNotificationTrait{},
		HeartbeatTrait{},
		StatusNotificationTrait{},
		MeterValuesTrait{},
		AuthorizeTrait{},
		StartTransactionTrait{},
		StopTransactionTrait{},
		ChangeConfigurationTrait{},
		DataTransferTrait{},
		SetChargingProfileTrait{},
		RemoteStartTransactionTrait{},
		RemoteStopTransactionTrait{},
		ResetTrait{},
		UnlockConnectorTrait{},
	)
}

var OCPP16M = &OCPP16Map{
	traitMap: make(map[string]ocpptrait),
}

type Response interface {
	Action() string
}
type Request interface {
	Action() string
}
type RequestHandler func(context.Context, Request) (Response, error)
type ocpptrait interface {
	Action() string
	RequestType() reflect.Type
	ResponseType() reflect.Type
}
type traitMap map[string]ocpptrait

type OCPP16Map struct {
	sync.RWMutex
	traitMap
}

func (m *OCPP16Map) register(traits ...ocpptrait) {
	m.Lock()
	defer m.Unlock()
	for _, trait := range traits {
		m.traitMap[trait.Action()] = trait
	}
}

func (m *OCPP16Map) GetTraitAction(action string) (ocpptrait, bool) {
	m.RLock()
	defer m.RUnlock()
	if v, ok := m.traitMap[action]; ok {
		return v, ok
	}
	return nil, false
}

//BootBotification

type BootNotificationTrait struct{}

func (BootNotificationTrait) Action() string {
	return BootNotificationName
}
func (BootNotificationTrait) RequestType() reflect.Type {
	return reflect.TypeOf(BootNotificationRequest{})
}
func (BootNotificationTrait) ResponseType() reflect.Type {
	return reflect.TypeOf(BootNotificationResponse{})
}

//HeartBeat
type HeartbeatTrait struct{}

func (HeartbeatTrait) Action() string {
	return HeartbeatName
}

func (HeartbeatTrait) RequestType() reflect.Type {
	return reflect.TypeOf(HeartbeatRequest{})
}

func (HeartbeatTrait) ResponseType() reflect.Type {
	return reflect.TypeOf(HeartbeatResponse{})
}

//StatusNotification
type StatusNotificationTrait struct{}

func (StatusNotificationTrait) Action() string {
	return StatusNotificationName
}

func (StatusNotificationTrait) RequestType() reflect.Type {
	return reflect.TypeOf(StatusNotificationRequest{})
}

func (StatusNotificationTrait) ResponseType() reflect.Type {
	return reflect.TypeOf(StatusNotificationResponse{})
}

// MeterValues
type MeterValuesTrait struct{}

func (MeterValuesTrait) Action() string {
	return MeterValuesName
}

func (MeterValuesTrait) RequestType() reflect.Type {
	return reflect.TypeOf(MeterValuesRequest{})
}

func (MeterValuesTrait) ResponseType() reflect.Type {
	return reflect.TypeOf(MeterValuesResponse{})
}

//Authorize
type AuthorizeTrait struct{}

func (AuthorizeTrait) Action() string {
	return AuthorizeName
}

func (AuthorizeTrait) RequestType() reflect.Type {
	return reflect.TypeOf(AuthorizeRequest{})
}

func (AuthorizeTrait) ResponseType() reflect.Type {
	return reflect.TypeOf(AuthorizeResponse{})
}

// StartTransaction
type StartTransactionTrait struct{}

func (StartTransactionTrait) Action() string {
	return StartTransactionName
}

func (StartTransactionTrait) RequestType() reflect.Type {
	return reflect.TypeOf(StartTransactionRequest{})
}

func (StartTransactionTrait) ResponseType() reflect.Type {
	return reflect.TypeOf(StartTransactionResponse{})
}

//StopTransaction
type StopTransactionTrait struct{}

func (StopTransactionTrait) Action() string {
	return StopTransactionName
}

func (StopTransactionTrait) RequestType() reflect.Type {
	return reflect.TypeOf(StopTransactionRequest{})
}

func (StopTransactionTrait) ResponseType() reflect.Type {
	return reflect.TypeOf(StopTransactionResponse{})
}

//ChangeConfiguration
type ChangeConfigurationTrait struct{}

func (ChangeConfigurationTrait) Action() string {
	return ChangeConfigurationName
}

func (ChangeConfigurationTrait) RequestType() reflect.Type {
	return reflect.TypeOf(ChangeConfigurationRequest{})
}

func (ChangeConfigurationTrait) ResponseType() reflect.Type {
	return reflect.TypeOf(ChangeConfigurationResponse{})
}

//DataTransfer
type DataTransferTrait struct{}

func (DataTransferTrait) Action() string {
	return DataTransferName
}

func (DataTransferTrait) RequestType() reflect.Type {
	return reflect.TypeOf(DataTransferRequest{})
}

func (DataTransferTrait) ResponseType() reflect.Type {
	return reflect.TypeOf(DataTransferResponse{})
}

//SetChargingProfile
type SetChargingProfileTrait struct{}

func (SetChargingProfileTrait) Action() string {
	return SetChargingProfileName
}

func (SetChargingProfileTrait) RequestType() reflect.Type {
	return reflect.TypeOf(SetChargingProfileRequest{})
}

func (SetChargingProfileTrait) ResponseType() reflect.Type {
	return reflect.TypeOf(SetChargingProfileResponse{})
}

//RemoteStartTransaction
type RemoteStartTransactionTrait struct{}

func (RemoteStartTransactionTrait) Action() string {
	return RemoteStartTransactionName
}

func (RemoteStartTransactionTrait) RequestType() reflect.Type {
	return reflect.TypeOf(RemoteStartTransactionRequest{})
}

func (RemoteStartTransactionTrait) ResponseType() reflect.Type {
	return reflect.TypeOf(RemoteStartTransactionResponse{})
}

//RemoteStopTransaction
type RemoteStopTransactionTrait struct{}

func (RemoteStopTransactionTrait) Action() string {
	return RemoteStopTransactionName
}

func (RemoteStopTransactionTrait) RequestType() reflect.Type {
	return reflect.TypeOf(RemoteStopTransactionRequest{})
}

func (RemoteStopTransactionTrait) ResponseType() reflect.Type {
	return reflect.TypeOf(RemoteStopTransactionResponse{})
}

//Reset
type ResetTrait struct{}

func (ResetTrait) Action() string {
	return ResetName
}

func (ResetTrait) RequestType() reflect.Type {
	return reflect.TypeOf(ResetRequest{})
}

func (ResetTrait) ResponseType() reflect.Type {
	return reflect.TypeOf(ResetResponse{})
}

//UnlockConnector

type UnlockConnectorTrait struct{}

func (UnlockConnectorTrait) Action() string {
	return UnlockConnectorName
}

func (UnlockConnectorTrait) RequestType() reflect.Type {
	return reflect.TypeOf(UnlockConnectorRequest{})
}

func (UnlockConnectorTrait) ResponseType() reflect.Type {
	return reflect.TypeOf(UnlockConnectorResponse{})
}
