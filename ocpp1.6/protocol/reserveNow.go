package protocol

import (
	validator "github.com/go-playground/validator/v10"
)

type ReservationStatus string

const (
	ReservationStatusAccepted    ReservationStatus = "Accepted"
	ReservationStatusFaulted     ReservationStatus = "Faulted"
	ReservationStatusOccupied    ReservationStatus = "Occupied"
	ReservationStatusRejected    ReservationStatus = "Rejected"
	ReservationStatusUnavailable ReservationStatus = "Unavailable"
)

func init() {
	Validate.RegisterValidation("reservationStatus", func(f validator.FieldLevel) bool {
		status := ReservationStatus(f.Field().String())
		switch status {
		case ReservationStatusAccepted, ReservationStatusFaulted, ReservationStatusOccupied, ReservationStatusRejected, ReservationStatusUnavailable:
			return true
		default:
			return false
		}
	})
}

type ReserveNowRequest struct {
	ConnectorId   *int   `json:"connectorId" validate:"required,gte=0"`
	ExpiryDate    string `json:"expiryDate" validate:"required,dateTime"`
	IdTag         string `json:"idTag" validate:"required,max=20"`
	ParentIdTag   string `json:"parentIdTag,omitempty" validate:"omitempty,max=20"`
	ReservationId *int   `json:"reservationId" validate:"required"`
}

func (r *ReserveNowRequest) Reset() {
	r.ConnectorId = nil
	r.ExpiryDate = ""
	r.IdTag = ""
	r.ParentIdTag = ""
	r.ReservationId = nil
}

func (ReserveNowRequest) Action() string {
	return ReserveNowName
}

type ReserveNowResponse struct {
	Status ReservationStatus `json:"status" validate:"required,reservationStatus"`
}

func (ReserveNowResponse) Action() string {
	return ReserveNowName
}

func (r *ReserveNowResponse) Reset() {
	r.Status = ""
}
