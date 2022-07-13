package storewrapper

import (
	"context"
	"crypto/tls"
	"io"
	"tsn-service/pkg/logger"
	"tsn-service/pkg/structures/configuration"

	"github.com/onosproject/onos-api/go/onos/config/diags"
	"github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/certs"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/proto"
)

var log = logger.GetLogger()

func GetRequestData(configId string) (*configuration.Request, error) {
	// Build the URN for the request data
	urn := "streams.requests." + configId

	// log.Info(urn)

	// Send request to specific path in k/v store "streams"
	reqData, err := getFromStore(urn)
	if err != nil {
		log.Errorf("Failed getting request data from store: %v", err)
		return &configuration.Request{}, err
	}

	return reqData, nil
}

func StoreConfiguration(req *configuration.ConfigResponse, confId string) error {
	// Serialize request
	obj, err := proto.Marshal(req)
	if err != nil {
		log.Errorf("Failed to marshal request: %v", err)
		return err
	}

	// Create a URN where the serialized request will be stored
	urn := "configurations.tsn-configuration." + confId

	log.Infof("Storing config response at: %s", urn)

	// TODO: Generate or use some ID to keep track of the specific stream request
	// urn += fmt.Sprintf("%v", uuid.New())

	// Send serialized request to it's specific path in a store
	err = sendToStore(obj, urn)
	if err != nil {
		return err
	}

	return nil
}

// TODO: Define topology structure or get from other repo
func GetTopology() ([]topo.Object, error) {
	ctx := context.Background()

	cert, err := tls.X509KeyPair([]byte(certs.DefaultClientCrt), []byte(certs.DefaultClientKey))
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	conn, err := grpc.Dial("onos-topo:5150", grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	if err != nil {
		log.Fatalf("Failed dialing onos-topo: %v", err)
		return nil, err
	}

	defer conn.Close()

	client := topo.CreateTopoClient(conn)

	// TODO: Make the correct ListRequest with correct filters to get everything necessary
	var filter = &topo.Filters{
		ObjectTypes: []topo.Object_Type{
			topo.Object_ENTITY,   // Switches?
			topo.Object_RELATION, // Links?
		},
	}

	resp, err := client.List(ctx, &topo.ListRequest{Filters: filter})
	if err != nil {
		log.Fatalf("Failed listing topo object: %v", errors.FromGRPC(err))
		return nil, err
	}

	log.Infof("Topo objects: %v", resp.Objects)

	return resp.Objects, nil
}

// TODO: Define topology structure or get from other repo
func GetConfiguration() (error, error) {
	ctx := context.Background()

	cert, err := tls.X509KeyPair([]byte(certs.DefaultClientCrt), []byte(certs.DefaultClientKey))
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	conn, err := grpc.Dial("onos-config:5150", grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	if err != nil {
		log.Fatalf("Failed dialing onos-config: %v", err)
		return nil, err
	}

	client := diags.CreateChangeServiceClient(conn)

	var req = &diags.ListNetworkChangeRequest{
		Subscribe:     false,
		WithoutReplay: false,
	}

	stream, err := client.ListNetworkChanges(ctx, req)
	if err != nil {
		log.Errorf("Failed getting list of network changes: %v", err)
		return nil, err
	}

	log.Infof("Successfully requested network changes from onos-config!")

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Warn(err)
				break
			}
			log.Infof("Type: %v", resp.Type)
		}

		conn.Close()
	}()

	return nil, nil

	// resp, err := client.List(ctx, &topo.ListRequest{Filters: filter})
	// if err != nil {
	// 	log.Fatalf("Failed listing topo object: %v", errors.FromGRPC(err))
	// 	return nil, err
	// }

	// log.Infof("Topo objects: %v", resp.Objects)

	// return resp.Objects, nil

	// watchClient, err := client.Watch(ctx, &topo.WatchRequest{Noreplay: false})
	// if err != nil {
	// 	log.Fatalf("Failed to watch topo for updates: %v", errors.FromGRPC(err))
	// 	return
	// }

	// go func() {
	// 	for {
	// 		resp, err := watchClient.Recv()
	// 		if err == io.EOF {
	// 			break
	// 		}
	// 		if err != nil {
	// 			log.Warn(err)
	// 			break
	// 		}
	// 		log.Infof("Event: %v", resp.Event)
	// 	}

	// 	conn.Close()
	// }()
}

//////////////////////////////////////////////////
/*                   TEMPLATES                  */
//////////////////////////////////////////////////
/*

func PublicFunctionName(req structureType) error {
	// Serialize request
	obj, err := proto.Marshal(req)
	if err != nil {
		log.Errorf("Failed to marshal request: %v", err)
		return err
	}

	// Create a URN where the serialized request will be stored
	urn := "store.type."

	// TODO: Generate or use some ID to keep track of the specific stream request
	urn += fmt.Sprintf("%v", uuid.New())

	// Send serialized request to it's specific path in a store
	err = sendToStore(obj, urn)
	if err != nil {
		return err
	}

	return nil
}

*/
