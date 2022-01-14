package protocol

import (
	validator "github.com/go-playground/validator/v10"
)

func init() {
	Validate.RegisterValidation("diagnosticsStatus", func(f validator.FieldLevel) bool {
		status := DiagnosticsStatus(f.Field().String())
		switch status {
		case DiagnosticsStatusIdle, DiagnosticsStatusUploaded, DiagnosticsStatusUploadFailed, DiagnosticsStatusUploading:
			return true
		default:
			return false
		}
	})
}

type DiagnosticsStatus string

const (
	DiagnosticsStatusIdle         DiagnosticsStatus = "Idle"
	DiagnosticsStatusUploaded     DiagnosticsStatus = "Uploaded"
	DiagnosticsStatusUploadFailed DiagnosticsStatus = "UploadFailed"
	DiagnosticsStatusUploading    DiagnosticsStatus = "Uploading"
)

type DiagnosticsStatusNotificationRequest struct {
	Status DiagnosticsStatus `json:"status" validate:"required,diagnosticsStatus"`
}

func (DiagnosticsStatusNotificationRequest) Action() string {
	return DiagnosticsStatusNotificationName
}

type DiagnosticsStatusNotificationResponse struct {
}

func (DiagnosticsStatusNotificationResponse) Action() string {
	return DiagnosticsStatusNotificationName
}
