package notification

import (
	"tsn-service/pkg/logger"
	handler "tsn-service/pkg/notification-handler"

	"golang.org/x/net/context"
)

var log = logger.GetLogger()

type Server struct {
	UnimplementedNotificationServer
}

func (s *Server) CalcConfig(ctx context.Context, in *IdList) (*UUID, error) {
	// in.GetValue() is the string: value:"some-uuid"
	// ids := strings.Split(in.GetValue(), "\"")[1]

	var idStringSlice []string

	ids := in.GetValues()

	for _, id := range ids {
		idStringSlice = append(idStringSlice, id.GetValue())
	}

	log.Infof("Received notification to calculate configuration for: %s", ids)

	configId, err := handler.CalculateConfiguration(idStringSlice)
	if err != nil {
		log.Errorf("Failed calculating configuration: %v", err)
		return nil, err
	}

	var transportConfId = &UUID{
		Value: configId,
	}

	return transportConfId, nil
}
