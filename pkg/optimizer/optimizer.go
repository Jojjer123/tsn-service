package optimizer

import (
	"io/ioutil"
	"tsn-service/pkg/logger"
	"tsn-service/pkg/structures/schedule"

	"github.com/onosproject/onos-api/go/onos/config/diags"
	"github.com/onosproject/onos-api/go/onos/topo"
	"gopkg.in/yaml.v2"
)

var log = logger.GetLogger()

func CalculateConf(topology []topo.Object, oldConfig []*diags.ListNetworkChangeResponse) ([]byte, error) {
	// TODO: Call optimizer

	// If optimizer isn't available
	// TODO: Get default scheduling
	//     Read and load file
	//     Store configuration based on default schedule in k/v store
	defaultSched, err := getDefaultSchedule()
	if err != nil {
		log.Errorf("Failed getting default schedule: %v", err)
		return nil, err
	}

	log.Infof("Default schedule: %v", defaultSched)

	// TODO: Calculate config for schedule
	var config = []byte("config here")

	return config, nil
}

func getDefaultSchedule() (*schedule.Schedule, error) {
	schedBytes, err := ioutil.ReadFile("./schedules/default-schedule.yaml")
	if err != nil {
		log.Errorf("Failed reading default schedule from file: %v", err)
		return nil, err
	}

	var defaultSched schedule.Schedule

	err = yaml.Unmarshal(schedBytes, &defaultSched)
	if err != nil {
		log.Errorf("Failed unmarshaling default schedule: %v", err)
		return nil, err
	}

	return &defaultSched, nil
}
