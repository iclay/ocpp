package protocol

import (
	validator "github.com/go-playground/validator/v10"
)

func init() {

	Validate.RegisterValidation("reason", func(f validator.FieldLevel) bool {
		reason := Reason(f.Field().String())
		switch reason {
		case EmergencyStop, EVDisconnected, HardReset, Local,
			Other, PowerLoss, Reboot, Remote,
			SoftReset, UnlockCommand, DeAuthorized:
			return true
		default:
			return false
		}
	})
}

type Reason string

const (
	EmergencyStop  Reason = "EmergencyStop"
	EVDisconnected Reason = "EVDisconnected"
	HardReset      Reason = "HardReset"
	Local          Reason = "Local"
	Other          Reason = "Other"
	PowerLoss      Reason = "PowerLoss"
	Reboot         Reason = "Reboot"
	Remote         Reason = "Remote"
	SoftReset      Reason = "SoftReset"
	UnlockCommand  Reason = "UnlockCommand"
	DeAuthorized   Reason = "DeAuthorized"
)

type StopTransactionRequest struct {
	IdTag           IdToken      `json:"idTag,omitempty" validate:"max=20"`
	MeterStop       int          `json:"meterStop" validate:"required,gte=0"`
	Timestamp       string       `json:"timestamp" validate:"required,dateTime"`
	TransactionId   int          `json:"transactionId" validate:"required"`
	Reason          Reason       `json:"reason,omitempty" validate:"omitempty,reason"`
	TransactionData []MeterValue `json:"transactionData,omitempty" validate:"omitempty,dive"`
}

func (StopTransactionRequest) Action() string {
	return StopTransactionName
}

type StopTransactionResponse struct {
	IdTagInfo IdTagInfo `json:"idTagInfo,omitempty" validate:"omitempty"`
}

func (StopTransactionResponse) Action() string {
	return StopTransactionName
}
