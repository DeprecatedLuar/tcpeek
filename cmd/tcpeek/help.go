package main

import (
	gohelp "github.com/DeprecatedLuar/gohelp-luar"
)

var helpPage = gohelp.NewPage("tcpeek", "TCP event listener — executes shell commands based on incoming TCP messages").
	Usage("tcpeek [command] [flags]").
	Section("Commands",
		gohelp.Item("(none)", "Start listening in the foreground"),
		gohelp.Item("-d", "Start as a background daemon"),
		gohelp.Item("stop", "Send SIGTERM to the running daemon"),
		gohelp.Item("restart", "Stop then start the daemon"),
		gohelp.Item("reconnect", "Re-establish TCP connections without restarting"),
		gohelp.Item("help", "Show this help page"),
	).
	Section("Flags",
		gohelp.Item("--debug", "Enable verbose debug logging"),
	).
	Text("Run 'tcpeek help config' for config file details.")

var configPage = gohelp.NewPage("config", "config file format and location").
	Text("Files are discovered by directory structure under $TCPEEK_CONFIG_DIR/ (if set), then $XDG_CONFIG_HOME/tcpeek/, falling back to ~/.config/tcpeek/.").
	Usage("$XDG_CONFIG_HOME/tcpeek/{IP}/{PORT}.toml").
	Section("Fields",
		gohelp.Item("[events]", "Map of TCP message payload → shell command to execute"),
	).
	Text("Example path: ~/.config/tcpeek/127.0.0.1/9999.toml\n\nExample content:\n  [events]\n  \"nav\"  = \"border-ctl --color blue\"\n  \"base\" = \"border-ctl --color gray\"").
	Section("Options",
		gohelp.Item("reconnect", "Re-connect on lost connection; true by default (set to false to disable)"),
	)

func helpCmd(args []string) {
	gohelp.Run(args, helpPage, configPage)
}
