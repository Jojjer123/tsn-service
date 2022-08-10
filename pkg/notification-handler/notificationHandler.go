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

// Calculates configuration and stores it as a set request in k/v store, returns ID of configuration set request
func CalculateConfiguration(ids []string) (string, error) {
	// Not yet used when calculating configuration
	var allRequestData []*configuration.Request

	// Get request from k/v store
	for _, requestId := range ids {
		reqData, err := store.GetRequestData(requestId)
		if err != nil {
			return "", err
		}

		allRequestData = append(allRequestData, reqData)
	}

	// Get topology
	topology, err := getTopology()
	if err != nil {
		log.Errorf("Failed getting topology: %v", err)
		return "", err
	}

	// NOT TESTED???
	// Get current configuration of the network
	oldConfig, err := getConfiguration()
	if err != nil {
		log.Errorf("Failed getting configuration: %v", err)
		return "", err
	}

	// Calculate configuration set request
	newConfig, err := optimizer.CalculateConf(topology, oldConfig)
	if err != nil {
		log.Errorf("Failed calculating configuration: %v", err)
		return "", err
	}

	// log.Infof("Config set request created: %v", newConfig)

	// Generate an ID for configuration set request
	confId := fmt.Sprint(uuid.New())

	// Store configuration set request in k/v store
	if err := store.StoreConfiguration(newConfig, confId); err != nil {
		log.Errorf("Failed storing configuration: %v", err)
		return "", err
	}

	return confId, nil
}
