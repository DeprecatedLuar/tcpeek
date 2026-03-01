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
			for i := 0; i < 10; i++ {
				time.Sleep(200 * time.Millisecond)
				if getRunningPid() == 0 {
					break
				}
			}
		}
		os.Remove(pidFile)
	}

	daemonize()
}
