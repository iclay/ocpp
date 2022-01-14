package protocol

import (
	validator "github.com/go-playground/validator/v10"
)

func init() {

	Validate.RegisterValidation("chargingRateUnit", func(f validator.FieldLevel) bool {
		typ := ChargingRateUnitType(f.Field().String())
		switch typ {
		case uintTypeA, uintTypeW:
			return true
		default:
			return false
		}
	})

	Validate.RegisterValidation("chargingProfilePurpose", func(f validator.FieldLevel) bool {
		purpose := ChargingProfilePurposeType(f.Field().String())
		switch purpose {
		case chargePointMaxProfile, txDefaultProfile, txProfile:
			return true
		default:
			return false
		}
	})

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
	ChargingRateUnitType       string
	ChargingProfileKindType    string
	ChargingProfilePurposeType string
	RecurrencyKindType         string
	RemoteStartStatus          string
)

const (
	uintTypeA ChargingRateUnitType = "A"
	uintTypeW ChargingRateUnitType = "W"
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

const (
	chargePointMaxProfile ChargingProfilePurposeType = "ChargePointMaxProfile"
	txDefaultProfile      ChargingProfilePurposeType = "TxDefaultProfile"
	txProfile             ChargingProfilePurposeType = "TxProfile"
)

type ChargingSchedulePeriod struct {
	StartPeriod  int     `json:"startPeriod" validate:"required,gte=0"`
	Limit        float64 `json:"limit" validate:"required,gte=0"`
	NumberPhases int     `json:"numberPhases,omitempty" validate:"omitempty,gte=0"`
}

type ChargingSchedule struct {
	Duration               int                      `json:"duration,omitempty" validate:"omitempty,gte=0"`
	StartSchedule          string                   `json:"startSchedule,omitempty" validate:"omitempty"`
	ChargingRateUnit       ChargingRateUnitType     `json:"chargingRateUnit" validate:"required,chargingRateUnit"`
	ChargingSchedulePeriod []ChargingSchedulePeriod `json:"chargingSchedulePeriod" validate:"required,min=1"`
	MinChargingRate        float64                  `json:"minChargingRate,omitempty" validate:"omitempty"`
}

type ChargingProfile struct {
	ChargingProfiled       int                        `json:"chargingProfileId" validate:"required"`
	TransactionId          int                        `json:"transactionId,omitempty" validate:"omitempty"`
	StackLevel             int                        `json:"stackLevel" validate:"required,gte=0"`
	ChargingProfilePurpose ChargingProfilePurposeType `json:"chargingProfilePurpose" validate:"required,chargingProfilePurpose"`
	ChargingProfileKind    ChargingProfileKindType    `json:"chargingProfileKind" validate:"required,chargingProfileKind"`
	RecurrencyKind         RecurrencyKindType         `json:"recurrencyKind,omitempty" validate:"omitempty,recurrencyKind"`
	ValidFrom              string                     `json:"validFrom,omitempty" validate:"omitempty,dateTime"`
	ValidTo                string                     `json:"validTo,omitempty" validate:"omitempty,dateTime"`
	ChargingSchedule       ChargingSchedule           `json:"chargingSchedule" validate:"required"`
}

type RemoteStartTransactionRequest struct {
	ConnectorId     int             `json:"connectorId,omitempty" validate:"omitempty,gte=0"`
	IdTag           IdToken         `json:"idTag" validate:"required,max=20"`
	ChargingProfile ChargingProfile `json:"chargingProfile,omitempty" validate:"omitempty"`
}

func (RemoteStartTransactionRequest) Action() string {
	return RemoteStartTransactionName
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
