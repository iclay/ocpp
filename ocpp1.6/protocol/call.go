package protocol

import (
	"encoding/json"
)

type Call struct {
	MessageTypeID MessageType `json:"messageTypeId" validate:"required,eq=2"`
	UniqueID      string      `json:"uniqueId" validate:"required,max=36"`
	Action        string      `json:"action" validate:"required,max=36"`
	Request       Request     `json:"payload" validate:"required"`
}

func (c *Call) MessageType() MessageType {
	return c.MessageTypeID
}

func (c *Call) UID() string {
	return c.UniqueID
}

func (c *Call) SpecificRequest() Request {
	switch c.Request.(type) {
	case BootNotificationRequest:
		c.Request = c.Request.(BootNotificationRequest)
	case HeartbeatRequest:
		c.Request = c.Request.(HeartbeatRequest)
	case StatusNotificationRequest:
		c.Request = c.Request.(StatusNotificationRequest)
	case MeterValuesRequest:
		c.Request = c.Request.(MeterValuesRequest)
	case AuthorizeRequest:
		c.Request = c.Request.(AuthorizeRequest)
	case StartTransactionRequest:
		c.Request = c.Request.(StartTransactionRequest)
	case StopTransactionRequest:
		c.Request = c.Request.(StopTransactionRequest)
	case ChangeConfigurationRequest:
		c.Request = c.Request.(ChangeConfigurationRequest)
	case DataTransferRequest:
		c.Request = c.Request.(DataTransferRequest)
	case SetChargingProfileRequest:
		c.Request = c.Request.(SetChargingProfileRequest)
	case RemoteStartTransactionRequest:
		c.Request = c.Request.(RemoteStartTransactionRequest)
	case RemoteStopTransactionRequest:
		c.Request = c.Request.(RemoteStopTransactionRequest)
	case ResetRequest:
		c.Request = c.Request.(ResetRequest)
	case UnlockConnectorRequest:
		c.Request = c.Request.(UnlockConnectorRequest)
	case SendLocalListRequest:
		c.Request = c.Request.(SendLocalListRequest)
	case GetLocalListVersionRequest:
		c.Request = c.Request.(GetLocalListVersionRequest)
	case GetConfigurationRequest:
		c.Request = c.Request.(GetConfigurationRequest)
	default:
	}
	return c.Request
}
func (c *Call) String() string {
	c.SpecificRequest()
	callBytes, _ := json.Marshal(c)
	return string(callBytes)
}

func (c *Call) MarshalJSON() ([]byte, error) {
	fields := make([]interface{}, 4)
	fields[0], fields[1], fields[2], fields[3] = int(c.MessageTypeID), c.UniqueID, c.Action, c.Request
	return json.Marshal(fields)
}
