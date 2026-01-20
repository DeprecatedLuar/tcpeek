package main

import (
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func restart() {
	// Try to stop existing process (ignore errors if not running)
	if data, err := os.ReadFile(pidFile); err == nil {
		if pid, err := strconv.Atoi(strings.TrimSpace(string(data))); err == nil {
			syscall.Kill(pid, syscall.SIGTERM)
			time.Sleep(500 * time.Millisecond)
		}
		os.Remove(pidFile)
	}

	daemonize()
}
