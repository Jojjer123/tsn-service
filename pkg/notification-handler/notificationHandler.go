package notificationHandler

import (
	"fmt"

	"tsn-service/pkg/logger"
	"tsn-service/pkg/optimizer"
	store "tsn-service/pkg/store-wrapper"
	"tsn-service/pkg/structures/configuration"

	"github.com/google/uuid"
)

var log = logger.GetLogger()

func CalculateConfiguration(ids []string) (string, error) {
	// log.Info("Calculating configuration...")

	var allRequestData []*configuration.Request

	// TODO: Get request from k/v store
	for _, requestId := range ids {
		reqData, err := store.GetRequestData(requestId)
		if err != nil {
			return "", err
		}

		allRequestData = append(allRequestData, reqData)
	}

	// TODO: Get topology
	// topology, err := store.GetTopology()
	topology, err := store.GetTopology()
	if err != nil {
		log.Errorf("Failed getting topology: %v", err)
		return "", err
	}

	// TODO: Get current configuration of the network
	oldConfig, err := store.GetConfiguration()
	if err != nil {
		log.Errorf("Failed getting configuration: %v", err)
		return "", err
	}

	// TODO: Calculate configuration
	newConfig, err := optimizer.CalculateConf(topology, oldConfig)
	if err != nil {
		log.Errorf("Failed calculating configuration: %v", err)
		return "", err
	}

	log.Infof("Config created: %v", newConfig)

	// TODO: Have template of configuration

	// TODO: Load template with random values (simulate applying configuration)

	confId := fmt.Sprintf("%v", uuid.New())

	// TODO: Store configuration in k/v store

	return confId, nil
}
