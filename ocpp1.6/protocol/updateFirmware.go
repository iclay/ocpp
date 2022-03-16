package protocol

type UpdateFirmwareRequest struct {
	Location      string `json:"location" validate:"required,uri"`
	Retries       *int   `json:"retries,omitempty" validate:"omitempty,gte=0"`
	RetrieveDate  string `json:"retrieveDate" validate:"required,dateTime"`
	RetryInterval *int   `json:"retryInterval,omitempty" validate:"omitempty,gte=0"`
}

func (UpdateFirmwareRequest) Action() string {
	return UpdateFirmwareName
}

type UpdateFirmwareResponse struct{}

func (UpdateFirmwareResponse) Action() string {
	return UpdateFirmwareName
}
