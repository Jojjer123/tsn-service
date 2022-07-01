package notificationServer

import (
	// "context"
	// "errors"
	"net"

	"tsn-service/pkg/logger"
	pb "tsn-service/pkg/structures/grpc/notification"

	// "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	// "google.golang.org/protobuf/types/known/emptypb"
)

var log = logger.GetLogger()

func CreateServer(protocol string, addr string) {
	lis, err := net.Listen(protocol, addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Infof("Now listening on %v", addr)

	// var opts []grpc.ServerOption
	// opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	s := pb.Server{}

	grpcServer := grpc.NewServer()

	log.Info("Created grpc server!")

	pb.RegisterNotificationServer(grpcServer, &s)

	log.Info("Starting to serve...")

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
