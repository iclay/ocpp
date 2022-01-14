package protocol

import (
	"encoding/json"
)

type CallResult struct {
	MessageTypeID MessageType `json:"messageTypeId" validate:"required,eq=3"`
	UniqueID      string      `json:"uniqueId" validate:"required,max=36"`
	Response      Response    `json:"payload" validate:"required"`
}

func (cr *CallResult) MessageType() MessageType {
	return cr.MessageTypeID
}

func (cr *CallResult) UID() string {
	return cr.UniqueID
}

func (cr *CallResult) SpecificResponse() Response {
	switch cr.Response.(type) {
	case BootNotificationResponse:
		cr.Response = cr.Response.(BootNotificationResponse)
	case HeartbeatResponse:
		cr.Response = cr.Response.(HeartbeatResponse)
	case StatusNotificationResponse:
		cr.Response = cr.Response.(StatusNotificationResponse)
	case MeterValuesResponse:
		cr.Response = cr.Response.(MeterValuesResponse)
	case AuthorizeResponse:
		cr.Response = cr.Response.(AuthorizeResponse)
	case StartTransactionResponse:
		cr.Response = cr.Response.(StartTransactionResponse)
	case StopTransactionResponse:
		cr.Response = cr.Response.(StopTransactionResponse)
	case ChangeConfigurationResponse:
		cr.Response = cr.Response.(ChangeConfigurationResponse)
	case DataTransferResponse:
		cr.Response = cr.Response.(DataTransferResponse)
	case SetChargingProfileResponse:
		cr.Response = cr.Response.(SetChargingProfileResponse)
	case RemoteStartTransactionResponse:
		cr.Response = cr.Response.(RemoteStartTransactionResponse)
	case RemoteStopTransactionResponse:
		cr.Response = cr.Response.(RemoteStopTransactionResponse)
	case ResetResponse:
		cr.Response = cr.Response.(ResetResponse)
	case UnlockConnectorResponse:
		cr.Response = cr.Response.(UnlockConnectorResponse)
	case SendLocalListResponse:
		cr.Response = cr.Response.(SendLocalListResponse)
	case GetLocalListVersionResponse:
		cr.Response = cr.Response.(GetLocalListVersionResponse)
	case GetConfigurationResponse:
		cr.Response = cr.Response.(GetConfigurationResponse)
	case FirmwareStatusNotificationResponse:
		cr.Response = cr.Response.(FirmwareStatusNotificationResponse)
	case DiagnosticsStatusNotificationResponse:
		cr.Response = cr.Response.(DiagnosticsStatusNotificationResponse)
	case ChangeAvailabilityResponse:
		cr.Response = cr.Response.(ChangeAvailabilityResponse)
	case ClearCacheResponse:
		cr.Response = cr.Response.(ClearCacheResponse)
	case GetCompositeScheduleResponse:
		cr.Response = cr.Response.(GetCompositeScheduleResponse)
	case ClearChargingProfileResponse:
		cr.Response = cr.Response.(ClearChargingProfileResponse)
	case CancelReservationResponse:
		cr.Response = cr.Response.(CancelReservationResponse)
	case ReserveNowResponse:
		cr.Response = cr.Response.(ReserveNowResponse)
	case TriggerMessageResponse:
		cr.Response = cr.Response.(TriggerMessageResponse)
	case UpdateFirmwareResponse:
		cr.Response = cr.Response.(UpdateFirmwareResponse)
	case GetDiagnosticsResponse:
		cr.Response = cr.Response.(GetDiagnosticsResponse)
	default:
	}
	return cr.Response
}

func (cr *CallResult) String() string {
	cr.SpecificResponse()
	callResBytes, _ := json.Marshal(cr)
	return string(callResBytes)
}

func (cr *CallResult) MarshalJSON() ([]byte, error) {
	fields := make([]interface{}, 3)
	fields[0], fields[1], fields[2] = int(cr.MessageTypeID), cr.UniqueID, cr.Response
	return json.Marshal(fields)
}
