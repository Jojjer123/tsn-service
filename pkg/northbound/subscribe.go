package northboundInterface

import (
	"github.com/google/gnxi/utils/credentials"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/openconfig/gnmi/proto/gnmi"
)

func (s *server) Subscribe(stream pb.GNMI_SubscribeServer) error {
	msg, ok := credentials.AuthorizeUser(stream.Context())
	if !ok {
		log.Infof("Denied a Subscribe request: %v", msg)

		return status.Error(codes.PermissionDenied, msg)
	}

	log.Infof("Allowed a Subscribe request: %v", msg)

	return nil
}
