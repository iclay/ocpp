package local
import (
	"context"
	"ocpp16/proto"
	"time"
)
type LocalService struct{}

func NewLocalService() *LocalService {
	return &LocalService{}
}
func (l *LocalService) BootNotification(ctx context.Context, request proto.Request) (proto.Response, error) {
	return &proto.BootNotificationResponse{
		CurrentTime: time.Now().Format(time.RFC3339),
		Interval:    10,
		Status:      "Accepted",
	}, nil
}



func (l *LocalService) RegisterOCPPHandler() map[string]proto.RequestHandler {
	return map[string]proto.RequestHandler{
		proto.BootNotificationName: proto.RequestHandler(l.BootNotification),
	}
}