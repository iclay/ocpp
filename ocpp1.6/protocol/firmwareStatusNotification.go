package protocol

import (
	validator "github.com/go-playground/validator/v10"
)

type FirmwareStatus string

const (
	FirmwareStatusDownloaded         FirmwareStatus = "Downloaded"
	FirmwareStatusDownloadFailed     FirmwareStatus = "DownloadFailed"
	FirmwareStatusDownloading        FirmwareStatus = "Downloading"
	FirmwareStatusIdle               FirmwareStatus = "Idle"
	FirmwareStatusInstallationFailed FirmwareStatus = "InstallationFailed"
	FirmwareStatusInstalling         FirmwareStatus = "Installing"
	FirmwareStatusInstalled          FirmwareStatus = "Installed"
)

func init() {
	Validate.RegisterValidation("firmwareStatus", func(f validator.FieldLevel) bool {
		status := FirmwareStatus(f.Field().String())
		switch status {
		case FirmwareStatusDownloaded, FirmwareStatusDownloadFailed, FirmwareStatusDownloading, FirmwareStatusIdle, FirmwareStatusInstallationFailed, FirmwareStatusInstalling, FirmwareStatusInstalled:
			return true
		default:
			return false
		}
	})
}

type FirmwareStatusNotificationRequest struct {
	Status FirmwareStatus `json:"status" validate:"required,firmwareStatus"`
}

func (FirmwareStatusNotificationRequest) Action() string {
	return FirmwareStatusNotificationName
}

type FirmwareStatusNotificationResponse struct {
}

func (FirmwareStatusNotificationResponse) Action() string {
	return FirmwareStatusNotificationName
}
