package protocol

import (
	"time"

	validator "github.com/go-playground/validator/v10"
)

var Validate = validator.New()

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

const (
	//currently, the timestamp supports rfc3339,ISO8601
	RFC3339     = time.RFC3339
	RFC3339Nano = time.RFC3339Nano
	ISO8601     = "2006-01-02T15:04:05Z"
)

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
	authAccepted     AuthorizationStatus = "Accepted"
	authBlock        AuthorizationStatus = "Blocked"
	authExpired      AuthorizationStatus = "Expired"
	authInvaliad     AuthorizationStatus = "Invalid"
	authConcurrentTx AuthorizationStatus = "ConcurrentTx"
)

type AuthorizationStatus string
type IdToken string
type IdTagInfo struct {
	ExpiryDate  string              `json:"expiryDate,omitempty" validate:"omitempty,dateTime"`
	ParentIdTag IdToken             `json:"parentIdTag,omitempty" validate:"omitempty,max=20"`
	Status      AuthorizationStatus `json:"status" validate:"required,authorizationStatus"`
}

const (
	BootNotificationName       = "BootNotification"
	HeartbeatName              = "Heartbeat"
	StatusNotificationName     = "StatusNotification"
	MeterValuesName            = "MeterValues"
	AuthorizeName              = "Authorize"
	StartTransactionName       = "StartTransaction"
	StopTransactionName        = "StopTransaction"
	ChangeConfigurationName    = "ChangeConfiguration"
	DataTransferName           = "DataTransfer"
	SetChargingProfileName     = "SetChargingProfile"
	RemoteStartTransactionName = "RemoteStartTransaction"
	RemoteStopTransactionName  = "RemoteStopTransaction"
	ResetName                  = "Reset"
	UnlockConnectorName        = "UnlockConnector"
	SendLocalListName          = "SendLocalList"
	GetLocalListVersionName    = "GetLocalListVersion"
	GetConfigurationName       = "GetConfiguration"
	CallName                   = "Call"
	CallErrorName              = "CallError"
	CallResultName             = "CallResult"
)

type MessageType int

const (
	CALL        = 2
	CALL_RESULT = 3
	CALL_ERROR  = 4
)
