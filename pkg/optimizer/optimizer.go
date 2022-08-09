package optimizer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"tsn-service/pkg/logger"
	store "tsn-service/pkg/store-wrapper"
	"tsn-service/pkg/structures/schedule"

	"github.com/ghodss/yaml"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/google/uuid"
	"github.com/onosproject/onos-api/go/onos/config/diags"
	"github.com/onosproject/onos-api/go/onos/topo"
	pb "github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/protobuf/proto"
)

var log = logger.GetLogger()

// var defaultSchedConfigID string
var defaultSchedID string

// Calculates configuration set request using optimizer, if that failes build configuration set request from default schedule
func CalculateConf(topology []topo.Object, oldConfig []*diags.ListNetworkChangeResponse) (*pb.SetRequest, error) {
	// TODO: Call optimizer

	// If optimizer isn't available

	// Load default schedule from k/v store
	sched, err := store.GetSchedule(defaultSchedID)
	if err != nil {
		log.Errorf("Failed getting default schedule: %v", err)
		return nil, err
	}

	// Create configuration set request based on default schedule and topology
	configSetReq, err := createConfigurationFromSchedule(sched, topology)
	if err != nil {
		log.Errorf("Failed creating default configuraiton set request: %v", err)
		return nil, err
	}

	return configSetReq, nil
}

// Reads default schedule config file and stores configuration for schedule in Atomix
func CreateDefaultSchedule() error {
	// Read default schedule from file
	schedBytes, err := ioutil.ReadFile("./schedules/default-schedule.yaml")
	if err != nil {
		log.Errorf("Failed reading default schedule from file: %v", err)
		return err
	}

	var defaultSched = &schedule.Schedule{}

	jsonBytes, err := yaml.YAMLToJSON(schedBytes)
	if err != nil {
		log.Errorf("Failed converting file content from yaml to json: %v", err)
		return err
	}

	if err = jsonpb.Unmarshal(bytes.NewReader(jsonBytes), defaultSched); err != nil {
		log.Errorf("Failed unmarshaling json to protobuf: %v", err)
		return err
	}

	// // Create config from default schedule
	// conf, err := createConfigurationFromSchedule(defaultSched)
	// if err != nil {
	// 	log.Errorf("Failed creating configuration from schedule: %v", err)
	// 	return err
	// }

	data, err := proto.Marshal(defaultSched)
	if err != nil {
		log.Errorf("Failed marshaling default schedule: %v", err)
		return err
	}

	// log.Infof("Data to be stored: %v", data)

	// Generate ID for default schedule configuration
	// defaultSchedConfigID = fmt.Sprintf("%v", uuid.New())
	// Generate ID for default schedule
	defaultSchedID = fmt.Sprintf("%v", uuid.New())

	// Store configuration
	// store.StoreConfiguration(conf, defaultSchedConfigID)
	// Store schedule
	err = store.StoreSchedule(data, defaultSchedID)
	if err != nil {
		log.Errorf("Failed storing default schedule: %v", err)
		return err
	}

	log.Infof("Successfully stored default schedule with ID: %v", defaultSchedID)

	return nil
}
