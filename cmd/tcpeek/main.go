package main

import (
	"fmt"
	"os"

	daemoninator "github.com/deprecatedluar/luar-daemonator"
)

const appName = "tcpeek"

var daemon = daemoninator.New(appName)

func main() {
	var debug bool
	var args []string
	for _, a := range os.Args[1:] {
		if a == "--debug" {
			debug = true
		} else {
			args = append(args, a)
		}
	}

	if len(args) == 0 {
		start(debug)
		return
	}

	switch args[0] {
	case "stop":
		stop()
	case "reconnect":
		reconnectCmd()
	case "help", "-h", "--help":
		helpCmd(args)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", args[0])
		fmt.Fprintln(os.Stderr, "Run 'tcpeek help' for usage.")
		os.Exit(1)
	}
}
