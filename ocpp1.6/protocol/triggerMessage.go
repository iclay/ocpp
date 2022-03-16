package protocol

import (
	validator "github.com/go-playground/validator/v10"
)

func init() {
	Validate.RegisterValidation("triggerMessageStatus", func(f validator.FieldLevel) bool {
		status := TriggerMessageStatus(f.Field().String())
		switch status {
		case TriggerMessageStatusAccepted, TriggerMessageStatusRejected, TriggerMessageStatusNotImplemented:
			return true
		default:
			return false
		}
	})
	Validate.RegisterValidation("messageTrigger", func(f validator.FieldLevel) bool {
		trigger := MessageTrigger(f.Field().String())
		switch trigger {
		case BootNotificationName, DiagnosticsStatusNotificationName, FirmwareStatusNotificationName, HeartbeatName, MeterValuesName, StatusNotificationName:
			return true
		default:
			return false
		}
	})
}

type TriggerMessageStatus string

type MessageTrigger string

const (
	TriggerMessageStatusAccepted       TriggerMessageStatus = "Accepted"
	TriggerMessageStatusRejected       TriggerMessageStatus = "Rejected"
	TriggerMessageStatusNotImplemented TriggerMessageStatus = "NotImplemented"
)

type TriggerMessageRequest struct {
	RequestedMessage MessageTrigger `json:"requestedMessage" validate:"required,messageTrigger"`
	ConnectorId      *int           `json:"connectorId,omitempty" validate:"omitempty,gt=0"`
}

func (TriggerMessageRequest) Action() string {
	return TriggerMessageName
}

type TriggerMessageResponse struct {
	Status TriggerMessageStatus `json:"status" validate:"required,triggerMessageStatus"`
}

func (TriggerMessageResponse) Action() string {
	return TriggerMessageName
}
