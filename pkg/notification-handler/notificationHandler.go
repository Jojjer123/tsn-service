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

	// log.Info(reqData)

	// TODO: Have template of configuration (response file)
	config, err := genConfigTemplate(allRequestData)
	// _, err = generateBaseResponse(reqData)
	if err != nil {
		log.Errorf("Failed generating response template: %v", err)
		return "", nil
	}

	// log.Info(resp)

	// TODO: Load template with random values (simulate applying configuration)

	confId := fmt.Sprintf("%v", uuid.New())

	// TODO: Store configuration in k/v store
	if err = store.StoreConfiguration(config); err != nil {
		log.Errorf("Failed storing configurations: %v", err)
		return "", nil
	}

	return confId, nil
}
