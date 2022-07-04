package notification

import (
	"strings"
	"tsn-service/pkg/logger"
	handler "tsn-service/pkg/notification-handler"

	"golang.org/x/net/context"
)

var log = logger.GetLogger()

type Server struct {
	UnimplementedNotificationServer
}

func (s *Server) CalcConfig(ctx context.Context, in *UUID) (*UUID, error) {
	id := strings.Split(in.GetValue(), "\"")[1] // in.GetValue() is the string: value:"some-uuid"

	log.Infof("Received notification to calculate configuration for UUID: %s", id)

	configId, err := handler.CalculateConfiguration(id)
	if err != nil {
		log.Errorf("Failed calculating configuration: %v", err)
		return nil, err
	}

	var transportConfId = &UUID{
		Value: configId,
	}

	return transportConfId, nil
}
