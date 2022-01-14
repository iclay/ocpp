package protocol

import (
	validator "github.com/go-playground/validator/v10"
)

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

func (StatusNotificationRequest) Action() string {
	return StatusNotificationName
}

type StatusNotificationResponse struct{}

func (StatusNotificationResponse) Action() string {
	return StatusNotificationName
}
