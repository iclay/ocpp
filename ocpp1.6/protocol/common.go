package protocol

import (
	"time"

	validator "github.com/go-playground/validator/v10"
)

var Validate = validator.New()

const (
	//currently, the timestamp supports rfc3339,RFC3339Nano, ISO8601
	RFC3339     = time.RFC3339
	RFC3339Nano = time.RFC3339Nano
	ISO8601     = "2006-01-02T15:04:05Z"
)

func init() {
	Validate.RegisterValidation("dateTime", func(f validator.FieldLevel) bool {
		timeString := f.Field().String()
		timeFormatList := []string{time.RFC3339, ISO8601, RFC3339Nano}
		for _, v := range timeFormatList {
			if _, err := time.Parse(v, timeString); err == nil {
				return true
			}
		}
		return false
	})
}

type ChargingRateUnitType string

const (
	uintTypeA ChargingRateUnitType = "A"
	uintTypeW ChargingRateUnitType = "W"
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
}

type ChargingProfilePurposeType string

const (
	chargePointMaxProfile ChargingProfilePurposeType = "ChargePointMaxProfile"
	txDefaultProfile      ChargingProfilePurposeType = "TxDefaultProfile"
	txProfile             ChargingProfilePurposeType = "TxProfile"
)

func init() {
	Validate.RegisterValidation("chargingProfilePurpose", func(f validator.FieldLevel) bool {
		purpose := ChargingProfilePurposeType(f.Field().String())
		switch purpose {
		case chargePointMaxProfile, txDefaultProfile, txProfile:
			return true
		default:
			return false
		}
	})
}

type AuthorizationStatus string
type IdToken string

const (
	authAccepted     AuthorizationStatus = "Accepted"
	authBlock        AuthorizationStatus = "Blocked"
	authExpired      AuthorizationStatus = "Expired"
	authInvaliad     AuthorizationStatus = "Invalid"
	authConcurrentTx AuthorizationStatus = "ConcurrentTx"
)

type IdTagInfo struct {
	ExpiryDate  string              `json:"expiryDate,omitempty" validate:"omitempty,dateTime"`
	ParentIdTag IdToken             `json:"parentIdTag,omitempty" validate:"omitempty,max=20"`
	Status      AuthorizationStatus `json:"status" validate:"required,authorizationStatus"`
}

func init() {
	Validate.RegisterValidation("authorizationStatus", func(f validator.FieldLevel) bool {
		status := AuthorizationStatus(f.Field().String())
		switch status {
		case authAccepted, authBlock, authExpired, authInvaliad, authConcurrentTx:
			return true
		default:
			return false
		}
	})

}

const (
	BootNotificationName              = "BootNotification"
	HeartbeatName                     = "Heartbeat"
	StatusNotificationName            = "StatusNotification"
	MeterValuesName                   = "MeterValues"
	AuthorizeName                     = "Authorize"
	StartTransactionName              = "StartTransaction"
	StopTransactionName               = "StopTransaction"
	ChangeConfigurationName           = "ChangeConfiguration"
	DataTransferName                  = "DataTransfer"
	SetChargingProfileName            = "SetChargingProfile"
	RemoteStartTransactionName        = "RemoteStartTransaction"
	RemoteStopTransactionName         = "RemoteStopTransaction"
	ResetName                         = "Reset"
	UnlockConnectorName               = "UnlockConnector"
	SendLocalListName                 = "SendLocalList"
	GetLocalListVersionName           = "GetLocalListVersion"
	GetConfigurationName              = "GetConfiguration"
	FirmwareStatusNotificationName    = "FirmwareStatusNotification"
	DiagnosticsStatusNotificationName = "DiagnosticsStatusNotification"
	ChangeAvailabilityName            = "ChangeAvailability"
	ClearCacheName                    = "ClearCache"
	GetCompositeScheduleName          = "GetCompositeSchedule"
	ClearChargingProfileName          = "ClearChargingProfile"
	CancelReservationName             = "CancelReservation"
	ReserveNowName                    = "ReserveNow"
	TriggerMessageName                = "TriggerMessage"
	UpdateFirmwareName                = "UpdateFirmware"
	GetDiagnosticsName                = "GetDiagnostics"
	CallName                          = "Call"
	CallErrorName                     = "CallError"
	CallResultName                    = "CallResult"
)

type MessageType int

const (
	CALL        = 2
	CALL_RESULT = 3
	CALL_ERROR  = 4
)
