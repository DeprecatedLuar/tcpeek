package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
)

var pidFile = "/tmp/tcpeek.pid"

func stop() {
	data, err := os.ReadFile(pidFile)
	if err != nil {
		log.Fatalf("[ERROR] Failed to read PID file: %v", err)
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		log.Fatalf("[ERROR] Invalid PID in file: %v", err)
	}

	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
		log.Fatalf("[ERROR] Failed to stop process: %v", err)
	}

	os.Remove(pidFile)
	fmt.Println("tcpeek stopped")
}

func writePID() error {
	return os.WriteFile(pidFile, []byte(strconv.Itoa(os.Getpid())), 0644)
}

func removePID() {
	os.Remove(pidFile)
}
