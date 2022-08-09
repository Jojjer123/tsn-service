package storewrapper

import (
	"tsn-service/pkg/logger"
	"tsn-service/pkg/structures/configuration"
	"tsn-service/pkg/structures/schedule"

	"github.com/golang/protobuf/proto"
	pb "github.com/openconfig/gnmi/proto/gnmi"
	// "google.golang.org/protobuf/proto"
)

var log = logger.GetLogger()

func GetRequestData(configId string) (*configuration.Request, error) {
	// Build the URN for the request data
	urn := "streams.requests." + configId

	// log.Info(urn)

	// Send request to specific path in k/v store "streams"
	rawData, err := getFromStore(urn)
	if err != nil {
		log.Errorf("Failed getting request data from store: %v", err)
		return &configuration.Request{}, err
	}

	// Unmarshal the byte slice from the store into request data
	var req = &configuration.Request{}
	err = proto.Unmarshal(rawData, req)
	if err != nil {
		log.Errorf("Failed to unmarshal request data from store: %v", err)
		return nil, err
	}

	return req, nil
}

func StoreConfiguration(config *pb.SetRequest, confId string) error {
	// Create a URN where the serialized request will be stored
	urn := "configurations.tsn-configuration." + confId

	// Serialize config
	rawConf, err := proto.Marshal(config)
	if err != nil {
		log.Errorf("Failed marshaling config: %v", err)
		return err
	}

	// Send serialized request to it's specific path in a store
	if err = sendToStore(rawConf, urn); err != nil {
		log.Errorf("Failed storing configuration: %v", err)
		return err
	}

	return nil
}

// func StoreConfiguration(conf []byte, confId string) error {
// 	// Create a URN where the serialized request will be stored
// 	urn := "configurations.tsn-configuration." + confId

// 	// Send serialized request to it's specific path in a store
// 	err := sendToStore(conf, urn)
// 	if err != nil {
// 		log.Errorf("Failed storing configuration: %v", err)
// 		return err
// 	}

// 	return nil
// }

func StoreSchedule(sched []byte, schedId string) error {
	// Create a URN where the serialized request will be stored
	urn := "configurations.schedules." + schedId

	// Send serialized request to it's specific path in a store
	err := sendToStore(sched, urn)
	if err != nil {
		log.Errorf("Failed storing schedule: %v", err)
		return err
	}

	return nil
}

func GetSchedule(schedId string) (*schedule.Schedule, error) {
	// Build the URN for the request data
	urn := "configurations.schedules." + schedId

	// Send request to specific path in k/v store "configurations"
	rawData, err := getFromStore(urn)
	if err != nil {
		log.Errorf("Failed getting request data from store: %v", err)
		return &schedule.Schedule{}, err
	}

	var sched = &schedule.Schedule{}
	if err = proto.Unmarshal(rawData, sched); err != nil {
		log.Errorf("Failed unmarshaling schedule: %v", err)
		return &schedule.Schedule{}, err
	}

	return sched, nil
}

// func StoreConfiguration(req *configuration.ConfigResponse, confId string) error {
// 	// Serialize request
// 	obj, err := proto.Marshal(req)
// 	if err != nil {
// 		log.Errorf("Failed to marshal request: %v", err)
// 		return err
// 	}

// 	// Create a URN where the serialized request will be stored
// 	urn := "configurations.tsn-configuration." + confId

// 	log.Infof("Storing config response at: %s", urn)

// 	// TODO: Generate or use some ID to keep track of the specific stream request
// 	// urn += fmt.Sprintf("%v", uuid.New())

// 	// Send serialized request to it's specific path in a store
// 	err = sendToStore(obj, urn)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

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
