package notificationHandler

import (
	"fmt"

	"tsn-service/pkg/logger"
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
	_, err := store.GetTopology()
	if err != nil {
		log.Errorf("Failed getting topology: %v", err)
		return "", err
	}

	// TODO: Get current configuration of the network
	_, err = store.GetConfiguration()
	if err != nil {
		log.Errorf("Failed getting configuration: %v", err)
		return "", err
	}

	// TODO: REDO EVERYTHING BELOW THIS COMMENT, CONFIGURATION IS NOT RESPONSE FILE

	// TODO: Have template of configuration (not response file, response file should only be in UNI or handler of main-service)
	// config, err := genConfigTemplate(allRequestData)
	// // _, err = generateBaseResponse(reqData)
	// if err != nil {
	// 	log.Errorf("Failed generating response template: %v", err)
	// 	return "", nil
	// }

	// log.Info(resp)

	// TODO: Load template with random values (simulate applying configuration)

	confId := fmt.Sprintf("%v", uuid.New())

	// TODO: Store configuration in k/v store
	// if err = store.StoreConfiguration(config, confId); err != nil {
	// 	log.Errorf("Failed storing configurations: %v", err)
	// 	return "", nil
	// }

	return confId, nil
}
