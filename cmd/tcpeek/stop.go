package main

import (
	"fmt"
	"log"
)

func stop() {
	if err := daemon.Stop(); err != nil {
		log.Fatalf("[ERROR] Failed to stop: %v", err)
	}

	fmt.Println("tcpeek stopped")
}
