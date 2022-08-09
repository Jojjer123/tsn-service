package main

import (
	"time"
	"tsn-service/pkg/logger"
	"tsn-service/pkg/optimizer"

	// northboundInterface "tsn-service/pkg/northbound"
	server "tsn-service/pkg/notificationServer"
)

var log = logger.GetLogger()

func main() {
	// northboundInterface.Start()

	if err := optimizer.CreateDefaultSchedule(); err != nil {
		log.Fatalf("Failed creating default schedule: %v", err)
		return
	}

	server.CreateServer("tcp", ":5150")

	log.Info("Back in main now...")

	for {
		time.Sleep(time.Second * 5)
	}
}
