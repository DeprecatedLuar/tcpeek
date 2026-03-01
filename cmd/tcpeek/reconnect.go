package main

import (
	"fmt"
	"log"
	"syscall"
)

func reconnectCmd() {
	pid := getRunningPid()
	if pid == 0 {
		fmt.Println("tcpeek not running, starting...")
		start(false)
		return
	}

	if err := syscall.Kill(pid, syscall.SIGUSR1); err != nil {
		log.Fatalf("[ERROR] Failed to send reconnect signal: %v", err)
	}

	fmt.Println("reconnect signal sent")
}
