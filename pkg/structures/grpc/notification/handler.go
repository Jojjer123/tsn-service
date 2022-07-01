package notification

import (
	"tsn-service/pkg/logger"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/emptypb"
)

var log = logger.GetLogger()

type Server struct {
	UnimplementedNotificationServer
}

// TODO: Change parameters in and out? At least out should contain some data for
// the main-service to know which configuration is the newly caluculated one.
func (s *Server) CalcConfig(ctx context.Context, in *empty.Empty) (*empty.Empty, error) {
	log.Info("Calculating configuration now...")

	// TODO: Have template of configuration (response file)

	// TODO: Load template with random values (simulate applying configuration)

	// TODO: Store configuration in k/v store

	return &emptypb.Empty{}, nil
}
