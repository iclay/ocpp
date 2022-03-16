package protocol

import (
	validator "github.com/go-playground/validator/v10"
)

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
	ConnectorId     *int            `json:"connectorId" validate:"required,gte=0"`
	ChargingProfile ChargingProfile `json:"csChargingProfiles" validate:"required"`
}

func (SetChargingProfileRequest) Action() string {
	return SetChargingProfileName
}

type ChargingProfileStatus string

type SetChargingProfileResponse struct {
	Status ChargingProfileStatus `json:"status" validate:"required,chargingProfileStatus"`
}

func (SetChargingProfileResponse) Action() string {
	return SetChargingProfileName
}
