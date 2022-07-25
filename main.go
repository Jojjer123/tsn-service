package main

import (
	"time"
	"tsn-service/pkg/logger"

	// northboundInterface "tsn-service/pkg/northbound"
	server "tsn-service/pkg/notificationServer"
)

var log = logger.GetLogger()

func main() {
	// northboundInterface.Start()

	server.CreateServer("tcp", ":5150")

	log.Info("Back in main now...")

	for {
		time.Sleep(time.Second * 5)
	}
}
