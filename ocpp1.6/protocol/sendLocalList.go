package protocol

import (
	validator "github.com/go-playground/validator/v10"
)

func init() {
	Validate.RegisterValidation("updateType", func(f validator.FieldLevel) bool {
		typ := UpdateType(f.Field().String())
		switch typ {
		case UpdateTypeDifferential, UpdateTypeFull:
			return true
		default:
			return false
		}
	})

	Validate.RegisterValidation("updateStatus", func(f validator.FieldLevel) bool {
		status := UpdateStatus(f.Field().String())
		switch status {
		case UpdateStatusAccepted, UpdateStatusFailed, UpdateStatusNotSupported, UpdateStatusVersionMismatch:
			return true
		default:
			return false
		}
	})
}

const (
	UpdateTypeDifferential UpdateType = "Differential"
	UpdateTypeFull         UpdateType = "Full"
)
const (
	UpdateStatusAccepted        UpdateStatus = "Accepted"
	UpdateStatusFailed          UpdateStatus = "Failed"
	UpdateStatusNotSupported    UpdateStatus = "NotSupported"
	UpdateStatusVersionMismatch UpdateStatus = "VersionMismatch"
)

type UpdateType string
type UpdateStatus string

type AuthorizationData struct {
	IdTag     string    `json:"idTag" validate:"required,max=20"`
	IdTagInfo IdTagInfo `json:"idTagInfo,omitempty" validate:"omitempty"` //TODO: validate required if update type is Full
}

type SendLocalListRequest struct {
	ListVersion            *int                `json:"listVersion" validate:"required,gte=0"`
	LocalAuthorizationList []AuthorizationData `json:"localAuthorizationList,omitempty" validate:"omitempty,dive"`
	UpdateType             UpdateType          `json:"updateType" validate:"required,updateType"`
}

func (SendLocalListRequest) Action() string {
	return SendLocalListName
}

func (r *SendLocalListRequest) Reset() {
	r.ListVersion = nil
	r.LocalAuthorizationList = nil
	r.UpdateType = ""
}

type SendLocalListResponse struct {
	Status UpdateStatus `json:"status" validate:"required,updateStatus"`
}

func (SendLocalListResponse) Action() string {
	return SendLocalListName
}

func (r *SendLocalListResponse) Reset() {
	r.Status = ""
}
