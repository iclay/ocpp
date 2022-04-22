package protocol

type GetLocalListVersionRequest struct {
}

func (r *GetLocalListVersionRequest) Reset() {}

func (GetLocalListVersionRequest) Action() string {
	return GetLocalListVersionName
}

type GetLocalListVersionResponse struct {
	ListVersion *int `json:"listVersion" validate:"required,gte=0"`
}

func (GetLocalListVersionResponse) Action() string {
	return GetLocalListVersionName
}

func (r *GetLocalListVersionResponse) Reset() {
	r.ListVersion = nil
}
