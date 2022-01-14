package protocol

type GetDiagnosticsRequest struct {
	Location      string `json:"location" validate:"required,uri"`
	Retries       int    `json:"retries,omitempty" validate:"omitempty,gte=0"`
	RetryInterval int    `json:"retryInterval,omitempty" validate:"omitempty,gte=0"`
	StartTime     string `json:"startTime,omitempty"`
	StopTime      string `json:"stopTime,omitempty"`
}

func (GetDiagnosticsRequest) Action() string {
	return GetDiagnosticsName
}

type GetDiagnosticsResponse struct {
	FileName string `json:"fileName,omitempty" validate:"max=255"`
}

func (GetDiagnosticsResponse) Action() string {
	return GetDiagnosticsName
}
