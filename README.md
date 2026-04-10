# SuperMac — macOS Power Tools for the CLI

<div align="center">

![macOS](https://img.shields.io/badge/macOS-12.0+-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![Go](https://img.shields.io/badge/Go-1.26+-00ADD8.svg)
![Architecture](https://img.shields.io/badge/architecture-Intel%20%26%20Apple%20Silicon-purple.svg)

**12 modules. 127 commands. 1 binary.**

Built by [CosmoLabs](https://cosmolabs.org) for macOS developers and power users.

</div>

---

## Quick Install

```bash
curl -fsSL https://raw.githubusercontent.com/CosmoLabs-org/SuperMac/master/install.sh | bash
```

Or with [Homebrew](https://brew.sh):
```bash
brew install cosmolabs-org/tap/supermac
```

## Command Structure

```
mac <category> <action> [arguments]
```

### Modules

| Module | Commands | Highlights |
|--------|----------|------------|
| **finder** | 6 | restart, toggle-hidden, reveal |
| **wifi** | 9 | on/off, scan, connect, status |
| **network** | 12 | ip, public-ip, speed-test, connections |
| **system** | 11 | info, cleanup, battery, temperature |
| **dev** | 16 | kill-port, servers, uuid, password, hash |
| **display** | 8 | brightness, dark-mode, night-shift, resolution |
| **dock** | 11 | position, autohide, size, magnification |
| **audio** | 11 | volume, mute, devices, balance |
| **screenshot** | 11 | format, location, shadow, record, sound |
| **bluetooth** | 6 | status, connect, disconnect, power |
| **apps** | 6 | list, info, cache-clear, kill |
| **power** | 21 | caffeinate, hidden-files, gatekeeper, animations |

### Global Shortcuts

```bash
mac ip                # network ip
mac cleanup           # system cleanup
mac dark              # display dark-mode on
mac light             # display dark-mode off
mac kp 3000           # dev kill-port 3000
mac vol 75            # audio volume 75
mac search wifi       # search all modules
```

## Usage Examples

```bash
# System
mac system info                  # Full system overview
mac system battery               # Battery health & status
mac system cleanup               # Deep system cleanup

# Network
mac network ip                   # Local IP address
mac network public-ip            # Public IP with geolocation
mac network speed-test           # Download speed test
mac wifi scan                    # Scan nearby networks

# Development
mac dev kill-port 3000           # Kill process on port
mac dev uuid                     # Generate UUID (copies to clipboard)
mac dev password 24              # Generate secure password
mac dev hash myfile.go sha256    # File hash

# Display & Audio
mac display brightness 75        # Set brightness
mac display dark-mode toggle     # Toggle dark mode
mac audio volume 50              # Set volume
mac audio balance center         # Center audio balance

# Bluetooth & Apps
mac bluetooth status             # Bluetooth power & devices
mac apps list                    # List installed apps
mac apps cache-clear Spotify     # Clear app cache
```

## Shell Completions

```bash
mac completion zsh  > ~/.zfunc/_mac
mac completion bash > /etc/bash_completion.d/mac
mac completion fish > ~/.config/fish/completions/mac.fish
```

## Building from Source

Requires [Go 1.26+](https://go.dev/dl/) and macOS.

```bash
git clone https://github.com/CosmoLabs-org/SuperMac.git
cd SuperMac/supermac-go
make build          # builds ./mac binary
make test           # run all tests
make install        # install to ~/bin/
```

## Architecture

```
supermac-go/
├── cmd/mac/main.go           # CLI entry point (Cobra)
├── internal/
│   ├── modules/              # 12 feature modules
│   │   ├── finder/           # Each module: struct + Commands()
│   │   ├── dock/
│   │   ├── system/
│   │   ├── wifi/
│   │   ├── network/
│   │   ├── display/
│   │   ├── audio/
│   │   ├── screenshot/
│   │   ├── bluetooth/
│   │   ├── apps/
│   │   ├── dev/
│   │   └── power/
│   ├── platform/             # macOS system call interface
│   │   ├── platform.go       # Interface (test seam)
│   │   ├── darwin.go         # Real implementation
│   │   └── mock.go           # Test mock
│   ├── module/               # Module/Command/Context types
│   ├── config/               # YAML configuration
│   ├── output/               # Text/JSON/Quiet output
│   └── version/              # ldflags-injected version
├── go.mod
└── Makefile
```

Modules never call `exec.Command` directly — they use `platform.Interface`, making everything testable with mocks.

## Bash Version

The original Bash implementation remains in the repo root (`mac` dispatcher + `lib/*.sh`). The Go rewrite lives in `supermac-go/` and is the primary implementation going forward.

## Roadmap

- **v0.2.0** — Go rewrite complete (current)
- **Website** — Interactive CLI demos, documentation
- **Desktop app** — SwiftUI wrapper around Go CLI (paid tier)

## License

MIT License. Built by [CosmoLabs](https://cosmolabs.org).
