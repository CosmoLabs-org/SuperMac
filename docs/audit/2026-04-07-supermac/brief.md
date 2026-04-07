# SuperMac v2.1.0 -- Project Brief

*Condensed intelligence for future sessions. Load this instead of re-auditing.*

## Tech Stack
- **Language**: Bash (shell scripts)
- **Platform**: macOS 12.0+ (claimed), 10.15+ (in code constant `MIN_MACOS_VERSION`)
- **Architecture**: Intel and Apple Silicon
- **Entry point**: `mac <category> <action> [args]`
- **Install method**: `curl | bash` one-liner, or local `install.sh`

## Project Stats
- **Total source lines**: ~13,500 across all `.sh` / `Makefile`
- **Shell scripts**: 25 files
- **Disk size**: 3.6 MB (including docs/config)
- **Version**: 2.1.0

## Architecture

```
SuperMac/
  mac                    -- Main dispatcher (350 lines, routes to lib/)
  lib/                   -- Modular category libraries (canonical source)
    utils.sh             -- Shared utilities: colors, output, validation
    finder.sh            -- File visibility, Finder restart
    display.sh           -- Brightness, dark mode, Night Shift, True Tone
    wifi.sh              -- WiFi on/off/toggle/status
    network.sh           -- IP, DNS flush, connectivity
    system.sh            -- Info, cleanup, battery, memory, CPU
    dev.sh               -- kill-port, servers, serve, processes, UUID
    dock.sh              -- Position, autohide, size, magnification
    audio.sh             -- Volume, mute, devices
    screenshot.sh        -- Location, format, shadow, cursor, naming
  bin/
    mac                  -- Symlinked/copy of root `mac` dispatcher
    install.sh           -- Installation script
  tests/
    test.sh              -- Test harness (50+ tests claimed)
  config/
    config.json          -- Settings, aliases, category registry
  Makefile               -- Dev automation (setup, test, lint, install-dev)
```

Root-level `.sh` files (audio.sh, dev.sh, etc.) are **duplicates** of `lib/` counterparts. This is the single most important structural issue.

## Key Patterns
- Every module sources `lib/utils.sh` for shared functions
- Every module exports a `<category>_dispatch` function and optionally `<category>_help`
- Global shortcuts (e.g., `mac ip` -> `network:ip`) defined in dispatcher `GLOBAL_SHORTCUTS` map
- Colors degrade gracefully when stdout is not a TTY
- `set -euo pipefail` on the main dispatcher; modules rely on caller error handling
- `osascript` used heavily for display/audio/AppleScript operations
- `defaults write` used for persistent system preference changes
- `sudo` used in network and system modules (DNS flush, log cleanup, network reset)

## External Dependencies
- `lsof` -- port management (dev module)
- `networksetup` -- WiFi/network control
- `system_profiler` -- hardware info
- `pmset` -- battery/power info
- `osascript` -- AppleScript bridge (display brightness, audio, dark mode)
- `defaults` -- macOS preferences
- `sqlite3` -- possibly for history (config.json references command_history)
- Standard Unix: `grep`, `awk`, `sed`, `tr`, `cut`, `sort`, `find`, `du`

## Config System
- `config/config.json` holds: version, command name, default settings, aliases, category registry, user preferences
- Aliases defined in both `config.json` and `GLOBAL_SHORTCUTS` in the dispatcher (potential drift)
- No runtime config reload -- read at startup only
