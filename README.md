# tcpeek

TCP event listener daemon that executes shell commands based on incoming TCP messages.

## Use Case

React to events from tools like [Kanata](https://github.com/jtroo/kanata) that broadcast state changes over TCP. Example: change window border colors when keyboard layers change.

## Install

```bash
go build -o tcpeek ./cmd/tcpeek
```

## Configuration

Config files live at `$XDG_CONFIG_HOME/tcpeek/{IP}/{PORT}.toml`

Example: `~/.config/tcpeek/127.0.0.1/9999.toml`

```toml
[events]
"nav" = "border-ctl --color blue"
"base" = "border-ctl --color gray"
```

When tcpeek receives `nav` on port 9999, it runs `border-ctl --color blue`.

## Usage

```bash
tcpeek start      # foreground
tcpeek start -d   # daemon mode
tcpeek stop       # stop daemon
tcpeek restart    # restart daemon
```

## License

MIT
