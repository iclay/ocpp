package protocol

import (
	validator "github.com/go-playground/validator/v10"
)

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

func (r *ResetRequest) Reset() {
	r.Type = ""
}

func (ResetRequest) Action() string {
	return ResetName
}

type ResetStatus string

const (
	resetAccepted ResetStatus = "Accepted"
	resetRejected ResetStatus = "Rejected"
)

type ResetResponse struct {
	Status ResetStatus `json:"status" validate:"required,resetStatus"`
}

func (ResetResponse) Action() string {
	return ResetName
}

func (r *ResetResponse) Reset() {
	r.Status = ""
}
