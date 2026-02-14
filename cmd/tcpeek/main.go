package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

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

	cmd := exec.Command(exe)
	cmd.Stdout = nil
	cmd.Stderr = nil
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
