package protocol

type GetDiagnosticsRequest struct {
	Location      string `json:"location" validate:"required,uri"`
	Retries       *int   `json:"retries,omitempty" validate:"omitempty,gte=0"`
	RetryInterval *int   `json:"retryInterval,omitempty" validate:"omitempty,gte=0"`
	StartTime     string `json:"startTime,omitempty" validate:"omitempty"`
	StopTime      string `json:"stopTime,omitempty" validate:"omitempty"`
}

func (GetDiagnosticsRequest) Action() string {
	return GetDiagnosticsName
}

func (r *GetDiagnosticsRequest) Reset() {
	r.Location = ""
	r.Retries = nil
	r.RetryInterval = nil
	r.StartTime = ""
	r.StopTime = ""
}

type GetDiagnosticsResponse struct {
	FileName string `json:"fileName,omitempty" validate:"omitempty,max=255"`
}

func (GetDiagnosticsResponse) Action() string {
	return GetDiagnosticsName
}

func (r *GetDiagnosticsResponse) Reset() {
	r.FileName = ""
}
