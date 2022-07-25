package northboundInterface

import (
	"time"

	"tsn-service/pkg/logger"
)

var log = logger.GetLogger()

func Start() {
	// Starts gNMI servers.
	go startServer(false, ":11161")
	go startServer(true, ":10161")

	for {
		time.Sleep(10 * time.Second)
	}
}
