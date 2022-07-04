package notificationHandler

import (
	"fmt"

	"tsn-service/pkg/logger"
	store "tsn-service/pkg/store-wrapper"

	"github.com/google/uuid"
)

var log = logger.GetLogger()

func CalculateConfiguration(id string) (string, error) {
	// log.Info("Calculating configuration...")

	// TODO: Get request from k/v store
	reqData, err := store.GetRequestData(id)
	if err != nil {
		return "", err
	}

	// log.Info(reqData)

	// TODO: Have template of configuration (response file)
	resp, err := generateBaseResponse(reqData)
	// _, err = generateBaseResponse(reqData)
	if err != nil {
		log.Errorf("Failed generating response template: %v", err)
		return "", nil
	}

	// log.Info(resp)

	// TODO: Load template with random values (simulate applying configuration)

	confId := fmt.Sprintf("%v", uuid.New())

	// TODO: Store configuration in k/v store
	if err = store.StoreConfiguration(resp); err != nil {
		log.Errorf("Failed storing configurations: %v", err)
		return "", nil
	}

	return confId, nil
}
