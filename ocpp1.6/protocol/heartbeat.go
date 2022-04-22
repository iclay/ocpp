package protocol

type HeartbeatRequest struct{}

func (HeartbeatRequest) Action() string {
	return HeartbeatName
}

func (r *HeartbeatRequest) Reset() {}

type HeartbeatResponse struct {
	CurrentTime string `json:"currentTime" validate:"required"`
}

func (HeartbeatResponse) Action() string {
	return HeartbeatName
}

func (r *HeartbeatResponse) Reset() {
	r.CurrentTime = ""
}
