package protocol

import (
	validator "github.com/go-playground/validator/v10"
)

func init() {

	Validate.RegisterValidation("unlockStatus", func(f validator.FieldLevel) bool {
		status := UnlockStatus(f.Field().String())
		switch status {
		case unlocked, unlockFailed, unlockNotSupported:
			return true
		default:
			return false
		}
	})
}

type UnlockConnectorRequest struct {
	ConnectorId int `json:"connectorId" validate:"required,gte=0"`
}

func (UnlockConnectorRequest) Action() string {
	return UnlockConnectorName
}

type UnlockStatus string

const (
	unlocked           UnlockStatus = "Unlocked"
	unlockFailed       UnlockStatus = "UnlockFailed"
	unlockNotSupported UnlockStatus = "NotSupported"
)

type UnlockConnectorResponse struct {
	Status UnlockStatus `json:"status" validate:"required,unlockStatus"`
}

func (UnlockConnectorResponse) Action() string {
	return UnlockConnectorName
}
