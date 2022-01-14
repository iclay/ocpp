package protocol

type HeartbeatRequest struct{}

func (HeartbeatRequest) Action() string {
	return HeartbeatName
}

type HeartbeatResponse struct {
	CurrentTime string `json:"currentTime" validate:"required"`
}

func (HeartbeatResponse) Action() string {
	return HeartbeatName
}
