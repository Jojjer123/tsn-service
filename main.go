package main

import (
	"time"
	northboundInterface "tsn-service/pkg/northbound"
)

func main() {

	northboundInterface.Start()

	for {
		time.Sleep(time.Second * 5)
	}
}
