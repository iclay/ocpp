package protocol

import (
	validator "github.com/go-playground/validator/v10"
)

func init() {

	Validate.RegisterValidation("configurationStatus", func(f validator.FieldLevel) bool {
		conf := ConfigurationStatus(f.Field().String())
		switch conf {
		case configAccepted, configRejected, rebootRequired, notSupported:
			return true
		default:
			return false
		}
	})
}

type ChangeConfigurationRequest struct {
	Key   string `json:"key"   validate:"required,max=50"`
	Value string `json:"value" validate:"required,max=500"`
}

func (ChangeConfigurationRequest) Action() string {
	return ChangeConfigurationName
}

type ConfigurationStatus string

const (
	configAccepted ConfigurationStatus = "Accepted"
	configRejected ConfigurationStatus = "Rejected"
	rebootRequired ConfigurationStatus = "RebootRequired"
	notSupported   ConfigurationStatus = "NotSupported"
)

type ChangeConfigurationResponse struct {
	Status ConfigurationStatus `json:"status" validate:"required,configurationStatus"`
}

func (ChangeConfigurationResponse) Action() string {
	return ChangeConfigurationName
}
