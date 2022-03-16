package protocol

import (
	validator "github.com/go-playground/validator/v10"
)

func init() {
	Validate.RegisterValidation("availabilityType", func(f validator.FieldLevel) bool {
		status := AvailabilityType(f.Field().String())
		switch status {
		case AvailabilityTypeOperative, AvailabilityTypeInoperative:
			return true
		default:
			return false
		}
	})

	Validate.RegisterValidation("availabilityStatus", func(f validator.FieldLevel) bool {
		status := AvailabilityStatus(f.Field().String())
		switch status {
		case AvailabilityStatusAccepted, AvailabilityStatusRejected, AvailabilityStatusScheduled:
			return true
		default:
			return false
		}
	})
}

type AvailabilityType string

const (
	AvailabilityTypeOperative   AvailabilityType = "Operative"
	AvailabilityTypeInoperative AvailabilityType = "Inoperative"
)

func isValidAvailabilityType(fl validator.FieldLevel) bool {
	status := AvailabilityType(fl.Field().String())
	switch status {
	case AvailabilityTypeOperative, AvailabilityTypeInoperative:
		return true
	default:
		return false
	}
}

type AvailabilityStatus string

const (
	AvailabilityStatusAccepted  AvailabilityStatus = "Accepted"
	AvailabilityStatusRejected  AvailabilityStatus = "Rejected"
	AvailabilityStatusScheduled AvailabilityStatus = "Scheduled"
)

type ChangeAvailabilityRequest struct {
	ConnectorId *int             `json:"connectorId" validate:"required,gte=0"`
	Type        AvailabilityType `json:"type" validate:"required,availabilityType"`
}

func (ChangeAvailabilityRequest) Action() string {
	return ChangeAvailabilityName
}

type ChangeAvailabilityResponse struct {
	Status AvailabilityStatus `json:"status" validate:"required,availabilityStatus"`
}

func (ChangeAvailabilityResponse) Action() string {
	return ChangeAvailabilityName
}
