package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"tcpeek/internal/config"
)

var logFile = filepath.Join(config.ConfigDir, "tcpeek.log")

func main() {
	if len(os.Args) < 2 {
		start()
		return
	}

	switch os.Args[1] {
	case "-d":
		daemonize()
	case "stop":
		stop()
	case "restart":
		restart()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func daemonize() {
	// Check if already running
	if pid := getRunningPid(); pid > 0 {
		fmt.Fprintf(os.Stderr, "tcpeek is already running (pid %d)\n", pid)
		os.Exit(1)
	}

	// Clean up stale pidfile
	os.Remove(pidFile)

	exe, _ := os.Executable()

	logF, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open log file: %v\n", err)
		os.Exit(1)
	}

	cmd := exec.Command(exe)
	cmd.Stdout = logF
	cmd.Stderr = logF
	cmd.Dir, _ = os.UserHomeDir()

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to start daemon: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("tcpeek started (pid %d)\n", cmd.Process.Pid)
}

func getRunningPid() int {
	out, err := exec.Command("pgrep", "-x", "tcpeek").Output()
	if err != nil {
		return 0
	}

	currentPid := os.Getpid()
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if pid, err := strconv.Atoi(line); err == nil && pid != currentPid {
			return pid
		}
	}
	return 0
}
