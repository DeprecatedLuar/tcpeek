# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What is tcpeek

TCP event listener daemon that executes shell commands based on incoming TCP messages. Primary use case: reacting to Kanata keyboard layer changes to trigger commands (e.g., changing border colors).

## Build & Run

```bash
go build -o tcpeek ./cmd/tcpeek
go build ./...   # verify all packages compile
go mod tidy      # after adding/removing imports
```

## Architecture

```
cmd/tcpeek/     # CLI entry point
  main.go            # arg parsing, routes to commands
  start.go           # loads config, starts listeners
  stop.go            # reads PID file, sends SIGTERM
  restart.go         # stop + start -d

internal/
  config/            # loads TOML configs from $XDG_CONFIG_HOME/tcpeek/{IP}/{PORT}.toml
  listener/          # TCP server per IP:port, each runs in goroutine
  executor/          # runs shell commands via sh -c
```

## Data Flow

1. `config.Load()` discovers config files by directory structure (IP folders, port.toml files)
2. `start.go` creates `listener.New(ip, port, events)` for each config
3. Each listener spawns goroutine on `Start()`, accepts connections
4. On TCP message: match payload → lookup command → `executor.Run(command)`

## Config Location

`$XDG_CONFIG_HOME/tcpeek/{IP}/{PORT}.toml` (falls back to `~/.config`)

Example: `~/.config/tcpeek/127.0.0.1/9999.toml`

```toml
[events]
"nav" = "border-ctl --color blue"
"base" = "border-ctl --color gray"
```

## Code Conventions

- Constants and package-level vars at top of file
- Each CLI command self-contained in its own file
- Listeners are self-contained with their own events map
