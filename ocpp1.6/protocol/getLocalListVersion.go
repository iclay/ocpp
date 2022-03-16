package protocol

type GetLocalListVersionRequest struct {
}

func (GetLocalListVersionRequest) Action() string {
	return GetLocalListVersionName
}

type GetLocalListVersionResponse struct {
	ListVersion *int `json:"listVersion" validate:"required,gte=0"`
}

func (GetLocalListVersionResponse) Action() string {
	return GetLocalListVersionName
}
