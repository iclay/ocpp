package protocol

import (
	validator "github.com/go-playground/validator/v10"
)

func init() {

	Validate.RegisterValidation("chargingProfileKind", func(f validator.FieldLevel) bool {
		kind := ChargingProfileKindType(f.Field().String())
		switch kind {
		case absolute, recurring, relative:
			return true
		default:
			return false
		}
	})

	Validate.RegisterValidation("recurrencyKind", func(f validator.FieldLevel) bool {
		kind := RecurrencyKindType(f.Field().String())
		switch kind {
		case daily, weekly:
			return true
		default:
			return false
		}
	})

	Validate.RegisterValidation("remoteStartStopStatus", func(f validator.FieldLevel) bool {
		status := RemoteStartStatus(f.Field().String())
		switch status {
		case remoteStartAccepted, remoteStartRejected:
			return true
		default:
			return false
		}
	})

}

type (
	ChargingProfileKindType string
	RecurrencyKindType      string
	RemoteStartStatus       string
)

const (
	daily  RecurrencyKindType = "Daily"
	weekly RecurrencyKindType = "Weekly"
)
const (
	absolute  ChargingProfileKindType = "Absolute"
	recurring ChargingProfileKindType = "Recurring"
	relative  ChargingProfileKindType = "Relative"
)

type ChargingProfile struct {
	ChargingProfileId      *int                       `json:"chargingProfileId" validate:"required"`
	TransactionId          *int                       `json:"transactionId,omitempty" validate:"omitempty"`
	StackLevel             *int                       `json:"stackLevel" validate:"required,gte=0"`
	ChargingProfilePurpose ChargingProfilePurposeType `json:"chargingProfilePurpose" validate:"required,chargingProfilePurpose"`
	ChargingProfileKind    ChargingProfileKindType    `json:"chargingProfileKind" validate:"required,chargingProfileKind"`
	RecurrencyKind         RecurrencyKindType         `json:"recurrencyKind,omitempty" validate:"omitempty,recurrencyKind"`
	ValidFrom              string                     `json:"validFrom,omitempty" validate:"omitempty,dateTime"`
	ValidTo                string                     `json:"validTo,omitempty" validate:"omitempty,dateTime"`
	ChargingSchedule       ChargingSchedule           `json:"chargingSchedule" validate:"required"`
}

type RemoteStartTransactionRequest struct {
	ConnectorId     *int             `json:"connectorId,omitempty" validate:"omitempty,gte=0"`
	IdTag           IdToken          `json:"idTag" validate:"required,max=20"`
	ChargingProfile *ChargingProfile `json:"chargingProfile,omitempty" validate:"omitempty"`
}

func (RemoteStartTransactionRequest) Action() string {
	return RemoteStartTransactionName
}

func (r *RemoteStartTransactionRequest) Reset() {
	r.ConnectorId = nil
	r.IdTag = ""
	r.ChargingProfile = nil
}

const (
	remoteStartAccepted RemoteStartStatus = "Accepted"
	remoteStartRejected RemoteStartStatus = "Rejected"
)

type RemoteStartTransactionResponse struct {
	Status RemoteStartStatus `json:"status" validate:"required,remoteStartStopStatus"`
}

func (RemoteStartTransactionResponse) Action() string {
	return RemoteStartTransactionName
}

func (r *RemoteStartTransactionResponse) Reset() {
	r.Status = ""
}
