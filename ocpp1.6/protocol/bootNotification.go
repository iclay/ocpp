package protocol

import (
	validator "github.com/go-playground/validator/v10"
)

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
