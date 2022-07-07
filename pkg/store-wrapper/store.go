package storewrapper

import (
	"tsn-service/pkg/logger"
	"tsn-service/pkg/structures/configuration"

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
