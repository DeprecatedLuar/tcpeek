package main

import (
	"fmt"
	"log"
	"syscall"
)

func reconnectCmd() {
	pid, running := daemon.PID()
	if !running {
		log.Fatalf("[ERROR] tcpeek is not running")
	}

	if err := syscall.Kill(pid, syscall.SIGUSR1); err != nil {
		log.Fatalf("[ERROR] Failed to send reconnect signal: %v", err)
	}

	fmt.Println("reconnect signal sent")
}
