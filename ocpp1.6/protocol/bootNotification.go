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
	ChargePointSerialNumber string `json:"chargePointSerialNumber,omitempty" validate:"omitempty,max=25"`
	ChargeBoxSerialNumber   string `json:"chargeBoxSerialNumber,omitempty" validate:"omitempty,max=25"`
	FirmwareVersion         string `json:"firmwareVersion,omitempty" validate:"omitempty,max=50"`
	Iccid                   string `json:"iccid,omitempty" validate:"omitempty,max=20"`
	Imsi                    string `json:"imsi,omitempty" validate:"omitempty,max=20"`
	MeterType               string `json:"meterType,omitempty" validate:"omitempty,max=25"`
	MeterSerialNumber       string `json:"meterSerialNumber,omitempty" validate:"omitempty,max=25"`
}

func (BootNotificationRequest) Action() string {
	return BootNotificationName
}

func (r *BootNotificationRequest) Reset() {
	r.ChargePointVendor = ""
	r.ChargePointModel = ""
	r.ChargePointSerialNumber = ""
	r.ChargeBoxSerialNumber = ""
	r.FirmwareVersion = ""
	r.Iccid = ""
	r.Imsi = ""
	r.MeterType = ""
	r.MeterSerialNumber = ""
}

type BootNotificationResponse struct {
	CurrentTime string             `json:"currentTime" validate:"required,dateTime"`
	Interval    *int               `json:"interval" validate:"required,gte=0"`
	Status      RegistrationStatus `json:"status" validate:"required,registrationStatus"`
}

func (BootNotificationResponse) Action() string {
	return BootNotificationName
}

func (r *BootNotificationResponse) Reset() {
	r.CurrentTime = ""
	r.Interval = nil
	r.Status = ""
}
