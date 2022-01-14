package protocol

type GetLocalListVersionRequest struct {
}

func (GetLocalListVersionRequest) Action() string {
	return GetLocalListVersionName
}

type GetLocalListVersionResponse struct {
	ListVersion int `json:"listVersion" validate:"gte=0"`
}

func (GetLocalListVersionResponse) Action() string {
	return GetLocalListVersionName
}

// GetConfiguration

type ConfigurationKey struct {
	Key      string `json:"key" validate:"required,max=50"`
	Readonly bool   `json:"readonly"`
	Value    string `json:"value,omitempty" validate:"max=500"`
}

type GetConfigurationRequest struct {
	Key []string `json:"key,omitempty" validate:"omitempty,unique,dive,max=50"`
}

func (GetConfigurationRequest) Action() string {
	return GetConfigurationName
}

type GetConfigurationResponse struct {
	ConfigurationKey []ConfigurationKey `json:"configurationKey,omitempty" validate:"omitempty,dive"`
	UnknownKey       []string           `json:"unknownKey,omitempty" validate:"omitempty,dive,max=50"`
}

func (GetConfigurationResponse) Action() string {
	return GetConfigurationName
}
