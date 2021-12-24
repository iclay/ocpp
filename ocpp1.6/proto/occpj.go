package proto

import (
	"encoding/json"
	"time"

	validator "github.com/go-playground/validator/v10"
)

type MessageType int

var Validate = validator.New()

const (
	CALL        = 2
	CALL_RESULT = 3
	CALL_ERROR  = 4
)

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
)

type CallAbstract interface {
	json.Marshaler
	QueryMessageType() MessageType
	QueryUniqueID() string
}

//Call
type Call struct {
	MessageTypeID MessageType `json:"messageTypeId" validate:"required,eq=2"`
	UniqueID      string      `json:"uniqueId" validate:"required,max=36"`
	Action        string      `json:"action" validate:"required,max=36"`
	Request       Request     `json:"payload" validate:"required"`
}

func (c *Call) QueryMessageType() MessageType {
	return c.MessageTypeID
}

func (c *Call) QueryUniqueID() string {
	return c.UniqueID
}

func (c *Call) MarshalJSON() ([]byte, error) {
	fields := make([]interface{}, 4)
	fields[0], fields[1], fields[2], fields[3] = int(c.MessageTypeID), c.UniqueID, c.Action, c.Request
	return json.Marshal(fields)

}

//CallResult
type CallResult struct {
	MessageTypeID MessageType `json:"messageTypeId" validate:"required,eq=3"`
	UniqueID      string      `json:"uniqueId" validate:"required,max=36"`
	Response      Response    `json:"payload" validate:"required"`
}

func (cr *CallResult) QueryMessageType() MessageType {
	return cr.MessageTypeID
}

func (cr *CallResult) QueryUniqueID() string {
	return cr.UniqueID
}

func (cr *CallResult) MarshalJSON() ([]byte, error) {
	fields := make([]interface{}, 3)
	fields[0], fields[1], fields[2] = int(cr.MessageTypeID), cr.UniqueID, cr.Response
	return json.Marshal(fields)
}

//CallError

func init() {
	Validate.RegisterValidation("errorCode", func(f validator.FieldLevel) bool {
		errcode := ErrCodeType(f.Field().String())
		switch errcode {
		case NotImplemented, NotSupported, CallInternalError, ProtocolError,
			SecurityError, FormationViolation, PropertyConstraintViolation, OccurenceConstraintViolation,
			TypeConstraintViolation, GenericError:
			return true
		}
		return false
	})
}

type CallError struct {
	MessageTypeID    MessageType `json:"messageTypeId" validate:"required,eq=4"`
	UniqueID         string      `json:"uniqueId" validate:"required,max=36"`
	ErrorCode        ErrCodeType `json:"errorCode" validate:"errorCode"`
	ErrorDescription string      `json:"errorDescription" validate:"required"`
	ErrorDetails     interface{} `json:"errorDetails" validate:"omitempty"`
}

func (ce *CallError) QueryMessageType() MessageType {
	return ce.MessageTypeID
}

func (ce *CallError) QueryUniqueID() string {
	return ce.UniqueID
}

func (ce *CallError) MarshalJSON() ([]byte, error) {
	fields := make([]interface{}, 5)
	fields[0], fields[1], fields[2], fields[3], fields[4] = int(ce.MessageTypeID), ce.UniqueID, ce.ErrorCode, ce.ErrorDescription, ce.ErrorDetails
	return json.Marshal(fields)
}

const (
	//currently, the timestamp supports rfc3339,ISO8601
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

//BootBotification
type RegistrationStatus string

const (
	bootAccepted RegistrationStatus = "Accepted"
	bootPending  RegistrationStatus = "Pending"
	bootRejected RegistrationStatus = "Rejected"
)

func init() {
	Validate.RegisterValidation("registrationStatus", func(f validator.FieldLevel) bool {
		status := RegistrationStatus(f.Field().String())
		switch status {
		case bootAccepted, bootPending, bootRejected:
			return true
		default:
			return false
		}
	})
}

type BootNotificationRequest struct {
	ChargePointVendor       string `json:"chargePointVendor" validate:"required,max=20"`
	ChargePointModel        string `json:"chargePointModel" validate:"required,max=20"`
	ChargePointSerialNumber string `json:"chargePointSerialNumber,omitempty" validate:"max=25"`
	ChargeBoxSerialNumber   string `json:"chargeBoxSerialNumber,omitempty" validate:"max=25"`
	FirmwareVersion         string `json:"firmwareVersion,omitempty" validate:"max=50"`
	Iccid                   string `json:"iccid,omitempty" validate:"max=20"`
	Imsi                    string `json:"imsi,omitempty" validate:"max=20"`
	MeterType               string `json:"meterType,omitempty" validate:"max=25"`
	MeterSerialNumber       string `json:"meterSerialNumber,omitempty" validate:"max=25"`
}

func (BootNotificationRequest) Action() string {
	return BootNotificationName
}

type BootNotificationResponse struct {
	CurrentTime string             `json:"currentTime" validate:"required"`
	Interval    int                `json:"interval" validate:"required,gte=0"`
	Status      RegistrationStatus `json:"status" validate:"required,registrationStatus"`
}

func (BootNotificationResponse) Action() string {
	return BootNotificationName
}

//HeartBeat
type HeartbeatRequest struct {
}

type HeartbeatResponse struct {
	CurrentTime string `json:"currentTime" validate:"required"`
}

//StatusNotification

func init() {
	Validate.RegisterValidation("chargePointErrorCode", func(f validator.FieldLevel) bool {
		errcode := ChargePointErrorCode(f.Field().String())
		switch errcode {
		case connectorLockFailure, eVCommunicationError, groundFailure, highTemperature,
			internalError, localListConflict, noError, otherError, overCurrentFailure,
			powerMeterFailure, powerSwitchFailure, readerFailure, resetFailure, underVoltage,
			overVoltage, weakSignal:
			return true
		default:
			return false
		}
	})
	Validate.RegisterValidation("chargePointStatus", func(f validator.FieldLevel) bool {
		status := ChargePointStatus(f.Field().String())
		switch status {
		case available, preparing, charging, suspendedEVSE, suspendedEV, finishing, reserved, unavailable, faulted:
			return true
		default:
			return false
		}
	})
}

type ChargePointStatus string

const (
	available     ChargePointStatus = "Available"
	preparing     ChargePointStatus = "Preparing"
	charging      ChargePointStatus = "Charging"
	suspendedEVSE ChargePointStatus = "SuspendedEVSE"
	suspendedEV   ChargePointStatus = "SuspendedEV"
	finishing     ChargePointStatus = "Finishing"
	reserved      ChargePointStatus = "Reserved"
	unavailable   ChargePointStatus = "Unavailable"
	faulted       ChargePointStatus = "Faulted"
)

type ChargePointErrorCode string

const (
	connectorLockFailure ChargePointErrorCode = "ConnectorLockFailure"
	eVCommunicationError ChargePointErrorCode = "EVCommunicationError"
	groundFailure        ChargePointErrorCode = "GroundFailure"
	highTemperature      ChargePointErrorCode = "HighTemperature"
	internalError        ChargePointErrorCode = "InternalError"
	localListConflict    ChargePointErrorCode = "LocalListConflict"
	noError              ChargePointErrorCode = "NoError"
	otherError           ChargePointErrorCode = "OtherError"
	overCurrentFailure   ChargePointErrorCode = "OverCurrentFailure"
	powerMeterFailure    ChargePointErrorCode = "PowerMeterFailure"
	powerSwitchFailure   ChargePointErrorCode = "PowerSwitchFailure"
	readerFailure        ChargePointErrorCode = "ReaderFailure"
	resetFailure         ChargePointErrorCode = "ResetFailure"
	underVoltage         ChargePointErrorCode = "UnderVoltage"
	overVoltage          ChargePointErrorCode = "OverVoltage"
	weakSignal           ChargePointErrorCode = "WeakSignal"
)

type StatusNotificationRequest struct {
	ConnectorId     int                  `json:"connectorId" validate:"required,gte=0"`
	ErrorCode       ChargePointErrorCode `json:"errorCode" validate:"required,chargePointErrorCode"`
	Info            string               `json:"info,omitempty" validate:"max=50"`
	Status          ChargePointStatus    `json:"status" validate:"required,chargePointStatus"`
	Timestamp       string               `json:"timestamp,omitempty" validate:"omitempty,dateTime"`
	VendorId        string               `json:"vendorId,omitempty" validate:"max=255"`
	VendorErrorCode string               `json:"vendorErrorCode,omitempty" validate:"max=50"`
}

type StatusNotificationResponse struct {
}

// MeterValues
func init() {
	Validate.RegisterValidation("readingContext", func(f validator.FieldLevel) bool {
		readingContext := ReadingContext(f.Field().String())
		switch readingContext {
		case interruptionBegin, interruptionEnd, sampleClock, samplePeriodic,
			transcationBegin, transcationEnd, trigger, other:
			return true
		default:
			return false
		}
	})
	Validate.RegisterValidation("valueFormat", func(fl validator.FieldLevel) bool {
		valueFormat := ValueFormat(fl.Field().String())
		switch valueFormat {
		case rawFormat, signFormat:
			return true
		default:
			return false
		}
	})

	Validate.RegisterValidation("measurand", func(f validator.FieldLevel) bool {
		measurand := Measurand(f.Field().String())
		switch measurand {
		case energyActiveExportRegister, energyActiveImportRegister, energyReactiveExportRegister, energyReactiveImportRegister,
			energyActiveExportInterval, energyActiveImportInterval, energyReactiveExportInterval, energyReactiveImportInterval,
			powerActiveExport, powerActiveImport, powerOffered, powerReactiveExport,
			powerReactiveImport, powerFactor, currentImport, currentExport,
			currentOffered, voltage, frequency, temperature, soc, rpm:
			return true
		default:
			return false
		}
	})

	Validate.RegisterValidation("phase", func(f validator.FieldLevel) bool {
		phase := Phase(f.Field().String())
		switch phase {
		case L1, L2, L3, N, L1N, L2N, L3N, L1L2, L2L3, L3L1:
			return true
		default:
			return false
		}
	})

	Validate.RegisterValidation("location", func(f validator.FieldLevel) bool {
		location := Location(f.Field().String())
		switch location {
		case cable, ev, inlet, outlet, body:
			return true
		default:
			return false
		}
	})

	Validate.RegisterValidation("unitOfMeasure", func(f validator.FieldLevel) bool {
		unitOfMeasure := UnitOfMeasure(f.Field().String())
		switch unitOfMeasure {
		case Wh, kWh, varh, kvarh,
			W, kW, VA, kVA,
			Var, kvar, A, V, K,
			Celcius, Celsius, Fahrenheit, Percent:
			return true
		default:
			return false
		}
	})
}

type UnitOfMeasure string

const (
	Wh         UnitOfMeasure = "Wh"
	kWh        UnitOfMeasure = "kWh"
	varh       UnitOfMeasure = "varh"
	kvarh      UnitOfMeasure = "kvarh"
	W          UnitOfMeasure = "W"
	kW         UnitOfMeasure = "kW"
	VA         UnitOfMeasure = "VA"
	kVA        UnitOfMeasure = "kVA"
	Var        UnitOfMeasure = "var"
	kvar       UnitOfMeasure = "kvar"
	A          UnitOfMeasure = "A"
	V          UnitOfMeasure = "V"
	K          UnitOfMeasure = "K"
	Celcius    UnitOfMeasure = "Celcius"
	Celsius    UnitOfMeasure = "Celsius"
	Fahrenheit UnitOfMeasure = "Fahrenheit"
	Percent    UnitOfMeasure = "Percent"
)

type Location string

const (
	cable  Location = "Cable"
	ev     Location = "EV"
	inlet  Location = "Inlet"
	outlet Location = "Outlet"
	body   Location = "Body"
)

type Phase string

const (
	L1   Phase = "L1"
	L2   Phase = "L2"
	L3   Phase = "L3"
	N    Phase = "N"
	L1N  Phase = "L1-N"
	L2N  Phase = "L2-N"
	L3N  Phase = "L3-N"
	L1L2 Phase = "L1-L2"
	L2L3 Phase = "L2-L3"
	L3L1 Phase = "L3-L1"
)

type Measurand string

const (
	energyActiveExportRegister   Measurand = "Energy.Active.Export.Register"
	energyActiveImportRegister   Measurand = "Energy.Active.Import.Register"
	energyReactiveExportRegister Measurand = "Energy.Reactive.Export.Register"
	energyReactiveImportRegister Measurand = "Energy.Reactive.Import.Register"
	energyActiveExportInterval   Measurand = "Energy.Active.Export.Interval"
	energyActiveImportInterval   Measurand = "Energy.Active.Import.Interval"
	energyReactiveExportInterval Measurand = "Energy.Reactive.Export.Interval"
	energyReactiveImportInterval Measurand = "Energy.Reactive.Import.Interval"
	powerActiveExport            Measurand = "Power.Active.Export"
	powerActiveImport            Measurand = "Power.Active.Import"
	powerOffered                 Measurand = "Power.Offered"
	powerReactiveExport          Measurand = "Power.Reactive.Export"
	powerReactiveImport          Measurand = "Power.Reactive.Import"
	powerFactor                  Measurand = "Power.Factor"
	currentImport                Measurand = "Current.Import"
	currentExport                Measurand = "Current.Export"
	currentOffered               Measurand = "Current.Offered"
	voltage                      Measurand = "Voltage"
	frequency                    Measurand = "Frequency"
	temperature                  Measurand = "Temperature"
	soc                          Measurand = "SoC"
	rpm                          Measurand = "RPM"
)

type ValueFormat string

const (
	rawFormat  ValueFormat = "Raw"
	signFormat ValueFormat = "SignedData"
)

type ReadingContext string

const (
	interruptionBegin ReadingContext = "Interruption.Begin"
	interruptionEnd   ReadingContext = "Interruption.End"
	sampleClock       ReadingContext = "Sample.Clock"
	samplePeriodic    ReadingContext = "Sample.Periodic"
	transcationBegin  ReadingContext = "Transaction.Begin"
	transcationEnd    ReadingContext = "Transaction.End"
	trigger           ReadingContext = "Trigger"
	other             ReadingContext = "Other"
)

type SampledValue struct {
	Value     string         `json:"value" validate:"required"`
	Context   ReadingContext `json:"context,omitempty" validate:"omitempty,readingContext"`
	Format    ValueFormat    `json:"format,omitempty" validate:"omitempty,valueFormat"`
	Measurand Measurand      `json:"measurand,omitempty" validate:"omitempty,measurand"`
	Phase     Phase          `json:"phase,omitempty" validate:"omitempty,phase"`
	Location  Location       `json:"location,omitempty" validate:"omitempty,location"`
	Unit      UnitOfMeasure  `json:"unit,omitempty" validate:"omitempty,unitOfMeasure"`
}

type MeterValue struct {
	Timestamp    string         `json:"timeStamp"    validate:"required,dateTime"`
	SampledValue []SampledValue `json:"sampledValue" validate:"required,min=1,dive"`
}

type MeterValuesRequest struct {
	ConnectorId   int          `json:"connectorId" validate:"required,gte=0"`
	TransactionId int          `json:"transactionId,omitempty"`
	MeterValue    []MeterValue `json:"meterValue"    validate:"required,min=1,dive"`
}

type MeterValuesResponse struct {
}

//Authorize

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

type IdToken string

type AuthorizeRequest struct {
	IdTag IdToken `json:"idTag" validate:"required,max=20"`
}

type AuthorizationStatus string

const (
	authAccepted     AuthorizationStatus = "Accepted"
	authBlock        AuthorizationStatus = "Blocked"
	authExpired      AuthorizationStatus = "Expired"
	authInvaliad     AuthorizationStatus = "Invalid"
	authConcurrentTx AuthorizationStatus = "ConcurrentTx"
)

type IdTagInfo struct {
	ExpiryDate  string              `json:"expiryDate,omitempty" validate:"omitempty"`
	ParentIdTag IdToken             `json:"parentIdTag,omitempty" validate:"omitempty,max=20"`
	Status      AuthorizationStatus `json:"status" validate:"required,authorizationStatus"`
}

type AuthorizeResponse struct {
	IdTagInfo IdTagInfo `json:"idTagInfo" validate:"required"`
}

// StartTransaction

type StartTransactionRequest struct {
	ConnectorId   int     `json:"type" validate:"required,gte=0"`
	IdTag         IdToken `json:"idTag" validate:"required,max=20"`
	MeterStart    int     `json:"meterStart" validate:"required,gte=0"`
	ReservationId int     `json:"reservationId,omitempty" validate:"omitempty"`
	Timestamp     string  `json:"timestamp" validate:"required,dateTime"`
}

type StartTransactionResponse struct {
	IdTagInfo     IdTagInfo `json:"idTagInfo" validate:"required"`
	TransactionId int       `json:"transactionId" validate:"required"`
}

//StopTransaction
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

type StopTransactionResponse struct {
	IdTagInfo IdTagInfo `json:"idTagInfo,omitempty" validate:"omitempty"`
}

//ChangeConfiguration
func init() {

	Validate.RegisterValidation("configurationStatus", func(f validator.FieldLevel) bool {
		conf := ConfigurationStatus(f.Field().String())
		switch conf {
		case configAccepted, configRejected, rebootRequired, notSupported:
			return true
		default:
			return false
		}
	})
}

type ChangeConfigurationRequest struct {
	Key   string `json:"key"   validate:"required,max=50"`
	Value string `json:"value" validate:"required,max=500"`
}

type ConfigurationStatus string

const (
	configAccepted ConfigurationStatus = "Accepted"
	configRejected ConfigurationStatus = "Rejected"
	rebootRequired ConfigurationStatus = "RebootRequired"
	notSupported   ConfigurationStatus = "NotSupported"
)

type ChangeConfigurationResponse struct {
	Status ConfigurationStatus `json:"status" validate:"required,configurationStatus"`
}

//DataTransfer
func init() {
	Validate.RegisterValidation("dataTransferStatus", func(f validator.FieldLevel) bool {
		status := DataTransferStatus(f.Field().String())
		switch status {
		case dataAccecpted, dataRejected, dataUnknownMessageId, dataUnknownVendorId:
			return true
		default:
			return false
		}
	})
}

type DataTransferRequest struct {
	VendorId  string `json:"vendorid"  validate:"required,max=255"`
	MessageId string `json:"messageId,omitempty" validate:"omitempty,max=50"`
	Data      string `json:"data,omitempty" validate:"omitempty"`
}

type DataTransferStatus string

const (
	dataAccecpted        DataTransferStatus = "Accepted"
	dataRejected         DataTransferStatus = "Rejected"
	dataUnknownMessageId DataTransferStatus = "UnknownMessageId"
	dataUnknownVendorId  DataTransferStatus = "UnknownVendorId"
)

type DataTransferResponse struct {
	Status DataTransferStatus `json:"status" validate:"required,dataTransferStatus"`
	Data   string             `json:"data,omitempty" validate:"omitempty"`
}

//SetChargingProfile
func init() {
	Validate.RegisterValidation("chargingProfileStatus", func(f validator.FieldLevel) bool {
		status := ChargingProfileStatus(f.Field().String())
		switch status {
		case chargingProfileAccepted, chargingProfileRejected, chargingProfileNotImplemented:
			return true
		default:
			return false
		}
	})
}

const (
	chargingProfileAccepted       ChargingProfileStatus = "Accepted"
	chargingProfileRejected       ChargingProfileStatus = "Rejected"
	chargingProfileNotImplemented ChargingProfileStatus = "NotImplemented"
)

type SetChargingProfileRequest struct {
	ConnectorId     int             `json:"connectorId" validate:"required,gte=0"`
	ChargingProfile ChargingProfile `json:"csChargingProfiles" validate:"required"`
}

type ChargingProfileStatus string

type SetChargingProfileResponse struct {
	Status ChargingProfileStatus `json:"status" validate:"required,chargingProfileStatus"`
}

//RemoteStartTransaction

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

const (
	remoteStartAccepted RemoteStartStatus = "Accepted"
	remoteStartRejected RemoteStartStatus = "Rejected"
)

type RemoteStartTransactionResponse struct {
	Status RemoteStartStatus `json:"status" validate:"required,remoteStartStopStatus"`
}

//RemoteStopTransaction

func init() {

	Validate.RegisterValidation("remoteStopStatus", func(f validator.FieldLevel) bool {
		status := RemoteStopStatus(f.Field().String())
		switch status {
		case remoteStopAccepted, remoteStopRejected:
			return true
		default:
			return false
		}
	})
}

type RemoteStopTransactionRequest struct {
	TransactionId int `json:"transactionId" validate:"required"`
}

type RemoteStopStatus string

const (
	remoteStopAccepted RemoteStopStatus = "Accepted"
	remoteStopRejected RemoteStopStatus = "Rejected"
)

type RemoteStopTransactionResponse struct {
	Status RemoteStopStatus `json:"status" validate:"required,remoteStopStatus"`
}

//Reset
func init() {
	Validate.RegisterValidation("resetType", func(f validator.FieldLevel) bool {
		typ := ResetType(f.Field().String())
		switch typ {
		case hard, soft:
			return true
		default:
			return false
		}
	})

	Validate.RegisterValidation("resetStatus", func(f validator.FieldLevel) bool {
		status := ResetStatus(f.Field().String())
		switch status {
		case resetAccepted, resetRejected:
			return true
		default:
			return false
		}
	})

}

type ResetType string

const (
	hard ResetType = "Hard"
	soft ResetType = "Soft"
)

type ResetRequest struct {
	Type ResetType `json:"type" validate:"required,resetType"`
}

type ResetStatus string

const (
	resetAccepted ResetStatus = "Accepted"
	resetRejected ResetStatus = "Rejected"
)

type ResetResponse struct {
	Status ResetStatus `json:"status" validate:"required,resetStatus"`
}

//UnlockConnector

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

type UnlockStatus string

const (
	unlocked           UnlockStatus = "Unlocked"
	unlockFailed       UnlockStatus = "UnlockFailed"
	unlockNotSupported UnlockStatus = "NotSupported"
)

type UnlockConnectorResponse struct {
	Status UnlockStatus `json:"status" validate:"required,unlockStatus"`
}
