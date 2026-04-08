# SuperMac -- Usage Guide

Professional macOS power tools for the command line.

```
mac <category> <action> [arguments] [flags]
```

---

## Table of Contents

- [Quick Start](#quick-start)
- [Installation](#installation)
- [Global Flags](#global-flags)
- [Global Shortcuts](#global-shortcuts)
- [Built-in Commands](#built-in-commands)
- [Modules](#modules)
  - [apps](#apps--application-management)
  - [audio](#audio--audio-control--device-management)
  - [bluetooth](#bluetooth--bluetooth-control--device-management)
  - [dev](#dev--developer-tools--utilities)
  - [power](#power--power-user-toggles)  - [display](#display--display--appearance-settings)
  - [dock](#dock--dock-management)
  - [finder](#finder--file-visibility--finder-management)
  - [network](#network--network-information--troubleshooting)
  - [screenshot](#screenshot--screenshot-settings--management)
  - [system](#system--system-information--maintenance)
  - [wifi](#wifi--wifi-control--management)
- [Output Formats](#output-formats)
- [Configuration](#configuration)
- [Shell Completions](#shell-completions)
- [Real-World Examples](#real-world-examples)
- [Project](#project)

---

## Quick Start

```bash
# See everything available
mac help

# Check your IP address
mac ip

# Toggle dark mode
mac dark

# Show system info
mac system info

# Check battery health
mac system battery

# List installed applications
mac apps list

# Check Bluetooth status
mac bluetooth status

# Kill a process on port 3000
mac kp 3000

# Search for any command by keyword
mac search volume
```

---

## Installation

### From Source

```bash
git clone https://github.com/cosmolabs-org/supermac.git
cd supermac/supermac-go

# Build the binary
make build

# Install to ~/bin/
make install

# Verify
mac version
```

### Build Details

The Makefile embeds version and build date at compile time using `-ldflags`:

```bash
make build                    # Uses git describe for version
VERSION=1.0.0 make build      # Override version manually
```

### Dependencies

Most commands use built-in macOS tools (`defaults`, `osascript`, `system_profiler`). A few commands benefit from optional third-party tools:

| Tool | Used by | Install |
|------|---------|---------|
| `SwitchAudioSource` | `audio input-device`, `audio output-device` | `brew install switchaudio-osx` |
| `blueutil` | `bluetooth` (all commands) | `brew install blueutil` |
| `dockutil` | `dock add`, `dock remove` | `brew install dockutil` |
| `python3` | `dev json-format` | Pre-installed on macOS |
| `curl` | `network public-ip`, `network speed-test` | Pre-installed on macOS |

---

## Global Flags

These flags apply to every command.

| Flag | Short | Description |
|------|-------|-------------|
| `--json` | | Output in JSON format for scripting |
| `--quiet` | | Suppress all output except errors |
| `--no-color` | | Disable colorized output |
| `--verbose` | `-v` | Enable verbose output |
| `--dry-run` | | Show what would be done without executing |
| `--yes` | `-y` | Skip confirmation prompts |

```bash
mac system memory --json          # JSON output for scripts
mac system cleanup --dry-run      # Preview without making changes
mac system cleanup -y             # Skip "are you sure?" prompt
mac wifi status --quiet           # Only errors, nothing else
```

---

## Global Shortcuts

Top-level convenience commands that delegate to module subcommands. These mirror the original Bash CLI's global shortcuts.

| Shortcut | Expands to | What it does |
|----------|------------|--------------|
| `mac ip` | `mac network ip` | Show local IP address |
| `mac cleanup` | `mac system cleanup` | Clean caches, logs, temp files |
| `mac restart-finder` | `mac finder restart` | Restart the Finder process |
| `mac kp <port>` | `mac dev kill-port <port>` | Kill process on a port |
| `mac vol` | `mac audio volume` | Show or set volume |
| `mac dark` | `mac display dark-mode on` | Enable dark mode |
| `mac light` | `mac display dark-mode off` | Switch to light mode |
| `mac search <term>` | (built-in) | Search all commands by keyword |

```bash
mac ip                            # Quick IP check
mac kp 3000                       # Free up port 3000
mac dark                          # Switch to dark mode
mac vol 50                        # Set volume to 50%
mac search battery                # Find battery-related commands
```

---

## Built-in Commands

```bash
mac version                       # Show version, build date, loaded modules
mac help                          # Show top-level help with all categories
mac help <category>               # Show commands for a specific category
mac config list                   # Show current configuration
mac config get <key>              # Get a specific config value
mac search <term>                 # Search commands by keyword across all modules
```

---

## Modules

### apps -- Application Management

List, inspect, kill, and open installed macOS applications.

**6 commands.**

| Command | Aliases | Args | What it does |
|---------|---------|------|--------------|
| `list` | | `[filter]` | List installed applications. Optionally filter by name. |
| `info` | | `<appname>` | Show detailed app info: version, size, path, architecture |
| `cache-clear` | | `<appname>` | Clear app cache/data (prompts for confirmation) |
| `recent` | | | Show recently used applications |
| `kill` | | `<appname>` | Force-quit an application |
| `open` | | `<appname>` | Open an application |

**Examples:**

```bash
mac apps list                       # List all installed apps
mac apps list chrome                # Filter apps containing "chrome"
mac apps info "Visual Studio Code"  # Version, size, path, arch
mac apps cache-clear Safari         # Clear Safari cache (prompts first)
mac apps cache-clear Safari -y      # Skip confirmation
mac apps recent                     # Recently used apps
mac apps kill Safari                # Force-quit Safari
mac apps open "Visual Studio Code"  # Launch VS Code
```

---

### audio -- Audio Control and Device Management

Control volume, mute, balance, and audio device switching.

**11 commands.**

| Command | Aliases | Args | What it does |
|---------|---------|------|--------------|
| `volume` | | `[level]` | Get or set system volume (0-100). Omit arg to read current. |
| `up` | | `[step]` | Increase volume by step (default: 10) |
| `down` | | `[step]` | Decrease volume by step (default: 10) |
| `mute` | | | Mute system audio |
| `unmute` | | | Unmute system audio |
| `toggle-mute` | | | Toggle mute state |
| `devices` | | `[type]` | List audio devices. Type: `all` (default), `input`, `output` |
| `input-device` | `input` | `<name>` | Switch input audio device (requires SwitchAudioSource) |
| `output-device` | `output` | `<name>` | Switch output audio device (requires SwitchAudioSource) |
| `status` | | | Show volume, mute state, active devices, sound effects |
| `balance` | | `[position]` | Set audio balance: `left`, `right`, `center`, or `0-100` |

**Examples:**

```bash
mac audio volume                  # Show current volume
mac audio volume 50               # Set volume to 50%
mac audio up                      # Increase by 10
mac audio up 25                   # Increase by 25
mac audio mute                    # Mute
mac audio toggle-mute             # Toggle mute
mac audio devices                 # List all audio devices
mac audio devices output          # List only output devices
mac audio output "External Speakers"  # Switch output device
mac audio input "Blue Yeti"       # Switch input device
mac audio balance left            # Pan audio fully left
mac audio balance 75              # Pan 75% right
mac audio status                  # Full audio status overview
```

---

### bluetooth -- Bluetooth Control and Device Management

Bluetooth power, device pairing, connections, and discoverability. Requires `blueutil`.

**6 commands.**

| Command | Aliases | Args | What it does |
|---------|---------|------|--------------|
| `status` | | | Show Bluetooth power state and connected devices |
| `devices` | | | List all paired Bluetooth devices |
| `connect` | | `<mac>` | Connect to a paired device by MAC address |
| `disconnect` | | `<mac>` | Disconnect a device by MAC address |
| `power` | | `[on/off/toggle]` | Control Bluetooth power. Omit arg to read current state. |
| `discoverable` | | `<on/off>` | Set Bluetooth discoverable mode |

**Examples:**

```bash
mac bluetooth status               # Power state + connected devices
mac bluetooth devices              # All paired devices
mac bluetooth connect "AA:BB:CC:DD:EE:FF"  # Connect to device
mac bluetooth disconnect "AA:BB:CC:DD:EE:FF"  # Disconnect device
mac bluetooth power                # Read current power state
mac bluetooth power on             # Turn Bluetooth on
mac bluetooth power toggle         # Toggle Bluetooth power
mac bluetooth discoverable on      # Make discoverable
mac bluetooth discoverable off     # Hide from scanning
```

---

### dev -- Developer Tools and Utilities

Port management, process monitoring, HTTP serving, encoding utilities, and code generation.

**16 commands.**

| Command | Aliases | Args | Flags | What it does |
|---------|---------|------|-------|--------------|
| `kill-port` | `kp` | `<port>` | | Kill the process listening on a port |
| `ports` | `list-ports` | | | Show all processes using network ports, with common dev port highlights |
| `servers` | | | | List running dev servers on common ports (3000, 5000, 8080, etc.) |
| `localhost` | | `<port>` | `--protocol, -p` (default: `http`) | Open localhost URL in default browser |
| `serve` | | `[dir]` | `--port, -p` (default: `8000`) | Start HTTP server for a directory |
| `processes` | | | `--sort, -s` (default: `cpu`), `--count, -n` (default: `15`) | Enhanced process viewer |
| `cpu-hogs` | | | | Show top 10 CPU-consuming processes (>1% CPU) |
| `memory-hogs` | | | | Show top 10 memory-consuming processes (>1% memory) |
| `uuid` | | | | Generate a UUID v4 and copy to clipboard |
| `env` | | | | Show language versions, tools, shell, and editor |
| `json-format` | `jf` | `<file>` | | Pretty-print a JSON file in-place |
| `base64-encode` | `b64e` | `<text>` | | Base64 encode text and copy to clipboard |
| `base64-decode` | `b64d` | `<text>` | | Base64 decode a string and copy to clipboard |
| `password` | `pw` | `[length]` | | Generate a secure random password (default: 20 chars) and copy to clipboard |
| `hash` | | `<file> [algorithm]` | | Compute file hash (`md5`, `sha1`, `sha256`, `sha512`). Default: `sha256`. Copies to clipboard. |
| `timestamp` | `ts` | `[value]` | | Convert between unix timestamps and dates. No args = now. Number = timestamp to date. Date string = to unix timestamp. |

**Examples:**

```bash
# Port management
mac dev kill-port 3000            # Kill whatever is on port 3000
mac kp 3000                       # Same thing via global shortcut
mac dev ports                     # Show all listening ports
mac dev servers                   # See which dev servers are running

# HTTP serving
mac dev serve                     # Serve current directory on port 8000
mac dev serve ./dist -p 3000      # Serve ./dist on port 3000
mac dev localhost 3000            # Open localhost:3000 in browser
mac dev localhost 443 -p https    # Open https://localhost:443

# Process inspection
mac dev processes                 # Top 15 by CPU
mac dev processes -s memory -n 20 # Top 20 by memory
mac dev cpu-hogs                  # Processes eating CPU
mac dev memory-hogs               # Processes eating memory

# Utilities
mac dev uuid                      # Generate + copy UUID
mac dev json-format data.json     # Pretty-print JSON in place
mac dev base64-encode "hello"     # Encode and copy
mac dev base64-decode "aGVsbG8="  # Decode and copy
mac dev password                  # 20-char password copied to clipboard
mac dev password 32               # 32-char password
mac dev env                       # Show dev environment overview

# Hash and timestamp
mac dev hash release.tar.gz       # SHA256 hash (default) copied to clipboard
mac dev hash data.bin md5         # MD5 hash
mac dev hash data.bin sha512      # SHA512 hash
mac dev timestamp                 # Current unix timestamp
mac dev timestamp 1712592000      # Convert timestamp to readable date
mac dev timestamp "2026-04-08"    # Convert date string to unix timestamp
```

---

### power -- Power User Toggles

Developer-focused system toggles that most users don't know about. Every command follows the same pattern: show state with no args, set with `on`/`off`, flip with `toggle`.

```bash
$ mac power status                    # Show all 20 toggles at a glance
$ mac power caffeinate on             # Prevent system sleep
$ mac power hidden-files toggle       # Toggle hidden files in Finder
$ mac power gatekeeper off            # Allow apps from anywhere (requires sudo)
```

| Command | Description | Sudo |
|---------|-------------|------|
| `status` | Show all toggles and their current state | No |
| `caffeinate` | Prevent system/display sleep (background process) | No |
| `hidden-files` | Show hidden files in Finder | No |
| `file-extensions` | Show all file extensions | No |
| `desktop-icons` | Show icons on desktop | No |
| `gatekeeper` | Allow apps from unidentified developers | Yes |
| `crash-reporter` | Show crash dialog (vs silent) | No |
| `function-keys` | Use F1-F12 as standard function keys | No |
| `spotlight-indexing` | Enable/disable Spotlight indexing | Yes |
| `key-repeat` | Fast key repeat rate (2x default) | No |
| `smooth-scrolling` | Smooth scrolling animation | No |
| `animations` | Window open/close animations | No |
| `transparency` | Reduce UI transparency | No |
| `dock-bounce` | Dock bounce animation on launch | No |
| `firewall-stealth` | Firewall stealth mode (no ping) | Yes |
| `remote-login` | SSH remote login server | Yes |
| `quarantine` | File quarantine warnings on downloads | No |
| `developer-dir` | Xcode developer tools path | Yes |
| `login-items` | Show login window items | No |
| `save-panel` | Expanded save dialogs by default | No |
| `print-dialog` | Expanded print dialogs by default | No |

---

### display -- Display and Appearance Settings

Brightness, dark mode, Night Shift, True Tone, wallpaper, and resolution.

**8 commands.**

| Command | Aliases | Args | What it does |
|---------|---------|------|--------------|
| `brightness` | | `<0-100>` | Set screen brightness percentage |
| `dark-mode` | `light-mode`, `toggle-mode` | `on/off/toggle` | Control dark mode appearance |
| `night-shift` | | `on/off` | Enable or disable Night Shift |
| `true-tone` | | `on/off` | Enable or disable True Tone |
| `wallpaper` | | `<path>` | Set desktop wallpaper from image file |
| `status` | | | Show brightness, appearance mode, Night Shift, display count, resolution |
| `detect` | | | Force detect connected displays |
| `resolution` | `res` | | List current and available display resolutions |

**Examples:**

```bash
mac display brightness 50         # Set brightness to 50%
mac display dark-mode on          # Enable dark mode
mac display dark-mode off         # Switch to light mode
mac display dark-mode toggle      # Toggle between dark and light
mac dark                          # Shortcut: dark mode on
mac light                         # Shortcut: light mode off
mac display night-shift on        # Reduce blue light
mac display true-tone off         # Disable True Tone
mac display wallpaper ~/Photos/bg.jpg  # Set wallpaper
mac display status                # Show all display settings
mac display resolution            # Show current resolution
mac display detect                # Re-detect displays
```

---

### dock -- Dock Management

Position, auto-hide, icon size, magnification, minimize effects, and app management.

**11 commands.**

| Command | Aliases | Args | What it does |
|---------|---------|------|--------------|
| `position` | | `<left/bottom/right>` | Set dock position. Also accepts `l`, `b`, `r`. |
| `autohide` | | `<on/off/toggle>` | Toggle dock auto-hide behavior |
| `size` | | `<value>` | Set icon size: `small` (32px), `medium` (64px), `large` (96px), or pixel count |
| `magnification` | | `<on/off/toggle>` | Toggle dock icon magnification on hover |
| `magnification-size` | | `<pixels>` | Set magnified icon size (16-256 pixels) |
| `minimize-effect` | | `<genie/scale>` | Set window minimize animation |
| `add` | | `<app>` | Add application to dock (requires `dockutil`) |
| `remove` | | `<app>` | Remove application from dock (requires `dockutil`) |
| `list` | | | List all items currently in the Dock |
| `status` | | | Show all current dock settings |
| `reset` | | | Reset dock to factory defaults |

**Examples:**

```bash
mac dock position left            # Move dock to left side
mac dock position bottom          # Move dock to bottom
mac dock autohide on              # Enable auto-hide
mac dock autohide toggle          # Toggle auto-hide
mac dock size small               # Small icons (32px)
mac dock size 48                  # Custom 48px icons
mac dock magnification on         # Enable magnification on hover
mac dock magnification-size 128   # Set magnified size to 128px
mac dock minimize-effect scale    # Use scale minimize effect
mac dock add "Visual Studio Code" # Add app to dock
mac dock remove "Chess"           # Remove app from dock
mac dock list                     # List all Dock items
mac dock status                   # Show all dock settings
mac dock reset                    # Reset everything to defaults
```

---

### finder -- File Visibility and Finder Management

Hidden files, Finder restarts, file revealing.

**6 commands.**

| Command | Aliases | Args | What it does |
|---------|---------|------|--------------|
| `show-hidden` | | | Show hidden (dot) files in Finder |
| `hide-hidden` | | | Hide hidden files in Finder |
| `toggle-hidden` | | | Toggle hidden files visibility |
| `reveal` | | `<path>` | Reveal a file or folder in Finder |
| `restart` | | | Restart the Finder process |
| `status` | | | Show Finder version and hidden files state |

**Examples:**

```bash
mac finder show-hidden            # Show dotfiles in Finder
mac finder hide-hidden            # Hide dotfiles
mac finder toggle-hidden          # Quick toggle
mac finder reveal ~/Desktop/file.txt  # Open Finder to file
mac finder restart                # Restart Finder
mac restart-finder                # Same via global shortcut
mac finder status                 # Show Finder settings
```

---

### network -- Network Information and Troubleshooting

IP addresses, DNS, connectivity checks, port scanning, speed tests, and network locations.

**12 commands.**

| Command | Aliases | Args | Flags | What it does |
|---------|---------|------|-------|--------------|
| `ip` | | | | Show local IP address and interface |
| `public-ip` | | | | Show public IP with geolocation (city, country, ISP) |
| `dns-flush` | `flush-dns` | | | Flush DNS cache (requires sudo) |
| `ping` | | `<host>` | `--count, -c` (default: `5`) | Ping a host with enhanced output |
| `ports` | | | | Show listening ports and owning processes |
| `interfaces` | | | | List all network interfaces with status and addresses |
| `connections` | | | | Show all active network connections (lsof) |
| `status` | `info` | | | Network overview: local IP, gateway, DNS, public IP |
| `speed-test` | | | | Download speed test via Cloudflare |
| `renew-dhcp` | | `[interface]` | | Renew DHCP lease (default: en0, requires sudo) |
| `reset` | | | | Reset all network settings (requires sudo) |
| `locations` | `loc` | | | List network locations and highlight current |

**Examples:**

```bash
mac network ip                    # Local IP and interface
mac ip                            # Same via global shortcut
mac network public-ip             # Public IP with geo info
mac network status                # Full network overview
mac network dns-flush             # Flush DNS (sudo)
mac network ping google.com       # Ping 5 packets
mac network ping 8.8.8.8 -c 3    # Ping 3 packets
mac network ports                 # All listening ports
mac network connections           # Active network connections
mac network interfaces            # Interface list
mac network speed-test            # Download speed test
mac network renew-dhcp            # Renew DHCP on en0
mac network renew-dhcp en1        # Renew DHCP on en1
mac network locations             # Show network locations
mac network reset                 # Full network reset (destructive)
```

---

### screenshot -- Screenshot Settings and Management

Save location, format, shadows, cursor, naming, thumbnails, sound, and capture.

**11 commands.**

| Command | Aliases | Args | What it does |
|---------|---------|------|--------------|
| `location` | `loc` | `[destination]` | Get or set screenshot save location. Named destinations: `desktop`, `downloads`, `clipboard`, `documents`, `pictures`. Or a custom path. |
| `format` | `type` | `[format]` | Get or set file format: `png`, `jpg`, `tiff`, `gif` |
| `shadow` | `shadows` | `[on/off/toggle]` | Control window shadow in screenshots. Omit arg to read current. |
| `cursor` | `show-cursor` | `[show/hide/toggle]` | Control cursor visibility in screenshots. Omit arg to read current. |
| `naming` | `name-format`, `name` | `[mode]` | Set naming: `sequential` or `timestamp`. Or a custom format string. |
| `thumbnail` | | `[on/off/toggle]` | Toggle screenshot thumbnail preview. Omit arg to read current. |
| `sound` | | `[on/off/toggle]` | Toggle camera shutter sound. Omit arg to read current. |
| `take` | | `[type]` | Take a screenshot now: `area` (default), `window`, or `screen` |
| `record` | | `[stop]` | Start screen recording. Use `stop` to end recording. |
| `status` | | | Show all screenshot settings plus keyboard shortcuts |
| `reset` | | | Reset all screenshot settings to macOS defaults |

**Examples:**

```bash
# Location
mac screenshot location           # Show current save location
mac screenshot location desktop   # Save to Desktop
mac screenshot location ~/Pictures/Screenshots  # Custom path
mac screenshot location clipboard # Save to clipboard only

# Format
mac screenshot format             # Show current format
mac screenshot format jpg         # Set JPEG (smaller files)
mac screenshot format png         # Set PNG (best quality)

# Shadows and cursor
mac screenshot shadow off         # No shadow on window captures
mac screenshot shadow toggle      # Toggle shadow
mac screenshot cursor show        # Include cursor in screenshots
mac screenshot cursor toggle      # Toggle cursor

# Naming
mac screenshot naming sequential  # Screenshot, Screenshot 1, Screenshot 2...
mac screenshot naming timestamp   # Screenshot 2026-04-08 at 14.30.00

# Other settings
mac screenshot thumbnail off      # Disable thumbnail preview
mac screenshot sound off          # Disable shutter sound
mac screenshot status             # See all current settings

# Capture
mac screenshot take               # Area selection screenshot
mac screenshot take window        # Click a window to capture
mac screenshot take screen        # Full screen screenshot

# Screen recording
mac screenshot record             # Start screen recording
mac screenshot record stop        # Stop screen recording

# Reset
mac screenshot reset              # Back to factory defaults
```

---

### system -- System Information and Maintenance

Hardware info, memory, CPU, battery, disk usage, processes, cleanup, and uptime.

**11 commands.**

| Command | Aliases | Args | What it does |
|---------|---------|------|--------------|
| `info` | | | Comprehensive system overview: OS, arch, hostname, uptime, shell, storage |
| `cleanup` | | | Deep cleanup: user caches, old downloads (30+ days), trash, logs, Safari cache, temp files, font caches |
| `battery` | | | Battery charge, charging status, time remaining, cycle count, health |
| `memory` | `mem` | | Memory usage: total, used, free, active, inactive, wired, compressed, swap, pressure assessment |
| `cpu` | | | CPU model, cores, threads, current usage, load average |
| `hardware` | | | Full hardware specs: model, chip, memory, serial, OS, build, architecture |
| `disk-usage` | | `[path]` | Disk usage analysis. Default: home directory. Shows volume info and top 10 largest subdirectories. |
| `processes` | | `[sort]` | Top processes by `cpu` (default) or `memory` |
| `uptime` | | | System uptime with user count and load averages |
| `updates` | | | Check for available macOS software updates |
| `temperature` | `temp` | | Show CPU temperature (requires sudo) |

**Examples:**

```bash
mac system info                   # Full system overview
mac system hardware               # Hardware specs
mac system battery                # Battery health check
mac system memory                 # Memory usage breakdown
mac system cpu                    # CPU info and current usage
mac system disk-usage             # Home directory disk usage
mac system disk-usage /Applications  # Specific directory
mac system processes              # Top processes by CPU
mac system processes memory       # Top processes by memory
mac system uptime                 # How long since last boot
mac system cleanup                # Interactive cleanup (prompts first)
mac system cleanup -y             # Cleanup without confirmation
mac system cleanup --dry-run      # Preview what would be cleaned
mac cleanup                       # Same via global shortcut
mac system updates                # Check for macOS updates
mac system temperature            # Show CPU temperature (sudo)
```

---

### wifi -- WiFi Control and Management

WiFi power, connections, scanning, saved networks, and signal strength.

**9 commands.**

| Command | Aliases | Args | What it does |
|---------|---------|------|--------------|
| `on` | | | Turn WiFi on. Shows current connection after enabling. |
| `off` | | | Turn WiFi off |
| `toggle` | | | Toggle WiFi on/off |
| `status` | | | Show interface, power state, connected network, signal strength, IP |
| `scan` | | | Scan for available WiFi networks (shows SSID, signal, security) |
| `connect` | `join` | `<network> [password]` | Connect to a WiFi network, optionally with password |
| `forget` | | `<network>` | Remove a saved WiFi network |
| `info` | | | Detailed connection info: SSID, BSSID, channel, country, gateway, DNS |
| `list-saved` | `saved` | | List all preferred/saved WiFi networks |

**Examples:**

```bash
mac wifi on                       # Turn WiFi on
mac wifi off                      # Turn WiFi off
mac wifi toggle                   # Toggle WiFi
mac wifi status                   # Current connection info
mac wifi scan                     # Scan for nearby networks
mac wifi connect "OfficeWiFi"     # Connect to known network
mac wifi connect "Guest" "p@ss"   # Connect with password
mac wifi forget "OldNetwork"      # Remove saved network
mac wifi info                     # Detailed connection details
mac wifi list-saved               # All saved networks
```

---

## Output Formats

### JSON Output

Every command supports `--json` for machine-readable output:

```bash
mac system memory --json
mac wifi status --json
mac audio devices --json
mac network ip --json
mac config list --json
```

Example output from `mac system memory --json`:

```json
{
  "type": "success",
  "message": "Memory: 16 GB total, 8 GB free"
}
```

### Quiet Mode

Use `--quiet` to suppress all output except errors. Useful in scripts where you only care about success or failure:

```bash
mac display dark-mode on --quiet
mac audio volume 50 --quiet
```

### Verbose Mode

Use `--verbose` or `-v` for additional detail during execution:

```bash
mac network ping google.com --verbose
mac system cleanup -v
```

---

## Configuration

SuperMac stores configuration in `~/.supermac/config.yaml`.

```bash
mac config list                   # Show all settings
mac config get output.format      # Get specific value
mac config get updates.channel    # Check update channel
```

**Supported config keys for `config get`:**

| Key | Values | Default |
|-----|--------|---------|
| `output.format` | `text`, `json`, `quiet` | `text` |
| `output.color` | `true`, `false` | `true` |
| `updates.check` | `true`, `false` | `true` |
| `updates.channel` | `stable`, `beta` | `stable` |

**Example configuration file:**

```yaml
version: 1

output:
  color: true
  format: text

updates:
  check: true
  channel: stable

aliases:
  kp: "dev kill-port"
  dark: "display dark-mode"
  ip: "network ip"
  cleanup: "system cleanup"
```

---

## Shell Completions

Generate completion scripts for your shell:

```bash
# Zsh
mac completion zsh > ~/.zfunc/_mac

# Bash
mac completion bash > /etc/bash_completion.d/mac

# Fish
mac completion fish > ~/.config/fish/completions/mac.fish
```

After generating, restart your shell or source the file to activate completions.

---

## Real-World Examples

### Daily Workflow

```bash
# Morning system check
mac system info                   # OS version, uptime, storage
mac system battery                # Battery health
mac wifi status                   # Am I connected?

# Quick adjustments
mac dark                          # Switch to dark mode
mac display brightness 60         # Comfortable brightness
mac audio volume 40               # Reasonable volume
```

### Development Setup

```bash
# Check what's running
mac dev servers                   # Running dev servers
mac dev ports                     # All listening ports

# Free up ports
mac kp 3000                       # Kill React dev server
mac kp 5000                       # Kill Flask server

# Start a local server
mac dev serve ./dist -p 8080      # Serve static files
mac dev localhost 8080            # Open in browser

# Utilities
mac dev uuid                      # Generate UUID for a new component
mac dev password 24               # Generate API key
mac dev json-format config.json   # Fix messy JSON
```

### Troubleshooting Networking

```bash
mac network status                # Full network overview
mac network ip                    # Local IP
mac network public-ip             # External IP and location
mac network ping google.com       # Test connectivity
mac network speed-test            # Measure download speed
mac network dns-flush             # Fix DNS issues (sudo)
mac network renew-dhcp            # Get fresh IP (sudo)
mac wifi info                     # Detailed WiFi connection
```

### System Cleanup

```bash
# Preview first
mac cleanup --dry-run             # See what would be deleted

# Run cleanup
mac system cleanup                # Interactive (prompts before proceeding)
mac system cleanup -y             # No confirmation prompt

# Check before and after
mac system disk-usage             # See disk usage before
mac system cleanup -y             # Clean caches, logs, trash, temp files
mac system disk-usage             # See disk usage after
```

### Dock Customization

```bash
mac dock position left            # Side dock for more screen space
mac dock autohide on              # Auto-hide for maximum space
mac dock size small               # Compact icons
mac dock magnification on         # Enlarge on hover
mac dock magnification-size 96    # Moderate magnification
mac dock add "Visual Studio Code" # Add your editor
mac dock remove "Chess"           # Remove unused apps
mac dock status                   # Verify settings
```

### Screenshot Configuration

```bash
mac screenshot location ~/Pictures/Screenshots  # Custom save path
mac screenshot format png       # Best quality
mac screenshot shadow off       # Clean window captures
mac screenshot cursor hide      # No cursor in shots
mac screenshot naming timestamp # Timestamped filenames
mac screenshot thumbnail off    # No floating preview
mac screenshot status           # Verify all settings
```

---

## Project

- **License**: MIT
- **Author**: CosmoLabs
- **Repository**: https://github.com/cosmolabs-org/supermac
- **Website**: https://cosmolabs.org
- **Modules**: 11 (apps, audio, bluetooth, dev, display, dock, finder, network, screenshot, system, wifi)
- **Commands**: 107
