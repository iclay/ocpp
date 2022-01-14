package protocol

import (
	validator "github.com/go-playground/validator/v10"
)

func init() {
	Validate.RegisterValidation("clearChargingProfileStatus", func(f validator.FieldLevel) bool {
		status := ClearChargingProfileStatus(f.Field().String())
		switch status {
		case ClearChargingProfileStatusAccepted, ClearChargingProfileStatusUnknown:
			return true
		default:
			return false
		}
	})
}

type ClearChargingProfileStatus string

const (
	ClearChargingProfileStatusAccepted ClearChargingProfileStatus = "Accepted"
	ClearChargingProfileStatusUnknown  ClearChargingProfileStatus = "Unknown"
)

func isValidClearChargingProfileStatus(f validator.FieldLevel) bool {
	status := ClearChargingProfileStatus(f.Field().String())
	switch status {
	case ClearChargingProfileStatusAccepted, ClearChargingProfileStatusUnknown:
		return true
	default:
		return false
	}
}

type ClearChargingProfileRequest struct {
	Id                     int                        `json:"id,omitempty" validate:"omitempty"`
	ConnectorId            int                        `json:"connectorId,omitempty" validate:"omitempty,gte=0"`
	ChargingProfilePurpose ChargingProfilePurposeType `json:"chargingProfilePurpose,omitempty" validate:"omitempty,chargingProfilePurpose"`
	StackLevel             int                        `json:"stackLevel,omitempty" validate:"omitempty,gte=0"`
}

func (ClearChargingProfileRequest) Action() string {
	return ClearChargingProfileName
}

type ClearChargingProfileResponse struct {
	Status ClearChargingProfileStatus `json:"status" validate:"required,clearChargingProfileStatus"`
}

func (ClearChargingProfileResponse) Action() string {
	return ClearChargingProfileName
}
