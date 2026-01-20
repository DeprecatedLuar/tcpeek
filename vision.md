# tcpeek Vision

## What It Is

tcpeek is a TCP event listener daemon that executes shell commands based on incoming TCP messages. It bridges the gap between applications that emit TCP events (like Kanata keyboard remapper) and command-line tools that need to react to those events.

## Primary Use Case

Kanata sends layer change notifications over TCP. tcpeek listens on a port, receives these notifications, and triggers corresponding commands (e.g., changing border colors to reflect the active keyboard layer).

## User Experience

### Installation & Setup

1. User places tcpeek binary somewhere in PATH
2. User creates config files in `~/.config/tcpeek/{IP}/{PORT}.toml`
3. Each config file defines what commands to run for specific TCP payloads

### Configuration Structure
```
~/.config/tcpeek/
└── 127.0.0.1/
    ├── 9999.toml
    ├── 8080.toml
    └── 5555.toml
```

- Directory name = IP address to listen on
- File name = port number to listen on (e.g., `9999.toml` → listens on port 9999)
- One daemon process manages all configs

### Config File Format

Each `{PORT}.toml` file contains an `[events]` section mapping payloads to commands:
```toml
[events]
# Simple string matching (for plain text payloads)
"nav" = "border-ctl --color blue"
"base" = "border-ctl --color gray"
"symbols" = "border-ctl --color purple"

# JSON path extraction (for structured payloads)
"LayerChange.new" = "border-ctl --layer {value}"
"status.brightness" = "notify-send 'Brightness: {value}'"
```

### Payload Matching Logic

When a TCP message arrives:

1. **Attempt JSON parsing**
   - If valid JSON, check all event keys as JSON paths
   - Extract value at that path
   - If key matches and path exists, substitute `{value}` in command and execute

2. **Fall back to exact string match**
   - If not JSON or no JSON path matched
   - Compare entire payload (trimmed) against event keys
   - Execute command if exact match found

3. **No match**
   - Log the unmatched payload
   - Continue listening

### Example Message Flow

**Example 1: Simple string**
```
TCP receives: "nav"
Matches: "nav" = "border-ctl --color blue"
Executes: border-ctl --color blue
```

**Example 2: JSON with path extraction**
```
TCP receives: {"LayerChange":{"new":"symbols"}}
Matches: "LayerChange.new" 
Extracts: "symbols"
Executes: border-ctl --layer symbols
```

**Example 3: Kanata push-msg**
```
Kanata config: (push-msg "active")
TCP receives: "active"
Matches: "active" = "border-ctl --state active"
Executes: border-ctl --state active
```

## Command Line Interface
```bash
tcpeek          # Run in foreground, log to stdout
tcpeek -d       # Daemonize (background), log to file
tcpeek stop     # Stop the daemon
tcpeek restart  # Restart the daemon (stop + start -d)
```

### Daemon Behavior

**When daemonized (`-d` flag):**
- Forks to background immediately
- Writes PID to `~/.config/tcpeek/tcpeek.pid`
- Logs to `~/.config/tcpeek/tcpeek.log`
- Spawns one TCP listener thread per config file found

**When run in foreground:**
- Logs to stdout
- Blocks until SIGINT/SIGTERM
- Still writes PID file (for stop/restart commands)

### Process Management

- **PID file location:** `~/.config/tcpeek/tcpeek.pid`
- **stop command:** Reads PID file, sends SIGTERM to that process
- **restart command:** Executes stop, waits for process to die, then runs start -d

## Error Handling

### Command Execution Failures

When an executed command returns non-zero exit code:
```bash
notify-send "tcpeek: command failed" "border-ctl --color blue (exit code: 1)"
```

The daemon continues running and processing other events.

### Configuration Errors

**Invalid TOML syntax:**
- Log error with filename and line number
- Skip that file
- Continue loading other configs

**Invalid port number in filename:**
- Log error
- Skip that file
- Continue loading other configs

**Port already in use:**
- Log error with port number
- Skip that port
- Continue with other ports

### Invalid Payloads

**Malformed JSON:**
- Attempt to parse as JSON fails
- Fall back to exact string matching
- Log warning if no match found

**Connection errors:**
- Log client connection/disconnection events
- Continue listening for new connections

## Logging Output

### Startup
```
[INFO] tcpeek starting
[INFO] Loading configs from ~/.config/tcpeek/
[INFO] Listening on 127.0.0.1:9999 (3 events configured)
[INFO] Listening on 127.0.0.1:8080 (5 events configured)
```

### Event Execution
```
[INFO] [127.0.0.1:9999] Received: "nav"
[INFO] [127.0.0.1:9999] Executing: border-ctl --color blue
```

### Errors
```
[ERROR] [127.0.0.1:9999] Command failed (exit 1): border-ctl --color blue
[ERROR] Failed to parse config: 9999.toml:5 - invalid TOML syntax
[WARN] [127.0.0.1:9999] No match found for payload: "unknown-layer"
```

## Technical Constraints

### Networking
- All listeners default to localhost (127.0.0.1) only
- Port range: 1-65535
- Supports multiple simultaneous connections per port
- Each connection handled in separate goroutine/thread

### Command Execution
- Commands executed through shell: `sh -c "command"`
- Inherits parent process environment variables
- Working directory: user's home directory
- No timeout on command execution (commands can run indefinitely)

### Performance Expectations
- Handle at least 100 events/second per port
- Near-instant response (<10ms latency from TCP receive to command exec)
- Minimal CPU usage when idle (<1%)
- Memory usage scales with number of configured ports (~1-2MB per listener)

## Out of Scope (Not Included)

- No authentication/authorization on TCP connections
- No built-in rate limiting
- No hot-reload of configs (must restart daemon)
- No wildcard/regex matching on payloads (only exact or JSON path)
- No multi-field JSON logic (AND/OR conditions)
- No UDP support
- No systemd service file (user manages process manually)
- No connection over network (localhost only via folder structure)

## Future Considerations (Not Now, But Keep in Mind)

- Hot reload on SIGHUP
- Config validation command: `tcpeek validate`
- Systemd service file generation
- Wildcard/regex matching in separate config sections
- Listen on 0.0.0.0 by adding IP folders beyond 127.0.0.1
- Metrics/stats endpoint

## Success Criteria

The project is successful when:

1. A user can run `tcpeek -d` and it daemonizes properly
2. Kanata layer changes (`{"LayerChange":{"new":"nav"}}`) trigger the correct border commands
3. Simple string payloads (`"active"`) work as expected
4. The daemon survives command failures and continues processing
5. `tcpeek stop` reliably terminates the daemon
6. Config errors don't crash the entire daemon
7. A user can have 3+ ports configured simultaneously without issues

## Non-Goals

- This is not a general-purpose message queue system
- This is not a replacement for systemd socket activation
- This is not a web server or HTTP API
- This is not a scripting engine (no conditionals/loops in config)
