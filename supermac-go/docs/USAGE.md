# SuperMac — Usage Guide

**Professional macOS power tools for the command line.**

`mac <category> <action> [arguments]`

---

## Installation

```bash
# Homebrew (recommended)
brew install cosmolabs-org/tap/supermac

# One-line install
curl -fsSL https://cosmolabs.org/install | bash

# From source
git clone https://github.com/cosmolabs-org/supermac.git
cd supermac/supermac-go && make install
```

## Global Flags

| Flag | Description |
|------|-------------|
| `--json` | Output in JSON format (for scripting) |
| `--quiet` | Suppress all output except errors |
| `--no-color` | Disable color output |
| `--verbose` | Verbose output |
| `--dry-run` | Show what would be done without executing |
| `--yes`, `-y` | Skip confirmation prompts |

## Built-in Commands

```bash
mac version              # Show version and system info
mac help                 # Show help
mac help <category>      # Show commands for a category
mac config list          # Show current configuration
mac config get <key>     # Get a config value
mac config set <key> <val> # Set a config value
mac config edit          # Open config in $EDITOR
```

---

## Modules

### 📁 finder — File Visibility & Finder Management

Control hidden files, Finder restarts, and file revealing.

```bash
mac finder show-hidden       # Show hidden files in Finder
mac finder hide-hidden       # Hide hidden files
mac finder toggle-hidden     # Toggle hidden files visibility
mac finder reveal <path>     # Reveal file in Finder
mac finder restart           # Restart Finder
mac finder status            # Show Finder settings
```

**Examples:**
```bash
mac finder toggle-hidden     # Quick toggle to see/hide dotfiles
mac finder reveal ~/Desktop/screenshot.png  # Open Finder to file
```

---

### 🌐 wifi — WiFi Control & Management

Manage WiFi connections, scan networks, check signal strength.

```bash
mac wifi on                  # Turn WiFi on
mac wifi off                 # Turn WiFi off
mac wifi toggle              # Toggle WiFi state
mac wifi status              # Show current connection info
mac wifi scan                # Scan for available networks
mac wifi connect <network>   # Connect to a network
mac wifi connect <network> <password>  # Connect with password
mac wifi forget <network>    # Forget a saved network
mac wifi info                # Detailed connection information
```

**Examples:**
```bash
mac wifi toggle              # Quick WiFi on/off
mac wifi status --json       # Machine-readable status
mac wifi scan | grep Office  # Find Office network
```

---

### 📡 network — Network Information & Troubleshooting

IP addresses, DNS, connectivity checks, port scanning.

```bash
mac network ip               # Show local and public IP
mac network public-ip        # Show public IP only
mac network dns-flush        # Flush DNS cache (requires sudo)
mac network ping <host>      # Ping a host
mac network ports            # Show listening ports
mac network interfaces       # List network interfaces
mac network status           # Network status overview
mac network reset            # Reset network settings (requires sudo)
```

**Examples:**
```bash
mac network ip --json        # {"local": "192.168.1.5", "public": "203.0.113.1"}
mac network ports            # Show what's listening
mac network ping google.com  # Test connectivity
```

---

### 🖥️ system — System Information & Maintenance

Hardware info, memory stats, battery status, system cleanup.

```bash
mac system info              # System overview (hostname, OS, uptime)
mac system cleanup           # Clean caches, logs, temp files
mac system battery           # Battery status and health
mac system memory            # Memory usage breakdown
mac system cpu               # CPU info and usage
mac system hardware          # Detailed hardware specs
```

**Examples:**
```bash
mac system memory --json     # Machine-readable memory stats
mac system battery           # Check battery health
mac system cleanup --dry-run # See what would be cleaned
mac system cleanup -y        # Skip confirmation prompts
```

---

### 💻 dev — Developer Tools & Utilities

Port management, process monitoring, HTTP serving, UUIDs.

```bash
mac dev kill-port <port>     # Kill process on a port
mac dev ports                # Show all listening ports
mac dev servers              # Show running servers
mac dev localhost <port>     # Start localhost server
mac dev serve <dir> <port>   # Serve a directory over HTTP
mac dev processes            # List running processes
mac dev processes --sort cpu # Sort by CPU usage
mac dev cpu-hogs             # Show top CPU consumers
mac dev memory-hogs          # Show top memory consumers
mac dev uuid                 # Generate a UUID
mac dev env                  # Show environment variables
```

**Examples:**
```bash
mac dev kill-port 3000       # Free up port 3000
mac dev serve . 8080         # Serve current directory
mac dev uuid                 # Generate UUID for scripts
mac dev processes --sort memory --json  # JSON process list
```

---

### 🖥️ display — Display & Appearance Settings

Brightness, dark mode, Night Shift, True Tone, wallpaper.

```bash
mac display brightness <0-100>  # Set brightness level
mac display dark-mode on      # Enable dark mode
mac display dark-mode off     # Disable dark mode
mac display dark-mode toggle  # Toggle dark mode
mac display night-shift on    # Enable Night Shift
mac display night-shift off   # Disable Night Shift
mac display true-tone on      # Enable True Tone
mac display true-tone off     # Disable True Tone
mac display wallpaper <path>  # Set desktop wallpaper
mac display status            # Show display settings
```

**Examples:**
```bash
mac display brightness 50     # Set to 50%
mac display dark-mode toggle  # Quick dark/light toggle
mac display status --json     # All display settings as JSON
```

---

### 🚢 dock — Dock Management & Positioning

Dock position, auto-hide, icon size, magnification, minimize effects.

```bash
mac dock position <left|bottom|right|top>  # Set dock position
mac dock autohide on          # Enable dock auto-hide
mac dock autohide off         # Disable dock auto-hide
mac dock autohide toggle      # Toggle auto-hide
mac dock size <value>         # Set icon size (1-128)
mac dock magnification on     # Enable magnification
mac dock magnification off    # Disable magnification
mac dock magnification toggle # Toggle magnification
mac dock magnification-size <value>  # Set magnified size
mac dock minimize-effect <genie|scale>  # Set minimize effect
mac dock status               # Show dock settings
mac dock reset                # Reset dock to defaults
```

**Examples:**
```bash
mac dock position bottom      # Move dock to bottom
mac dock size 48              # Set icon size
mac dock status               # See all current dock settings
```

---

### 🔊 audio — Audio Control & Device Management

Volume control, mute, device switching.

```bash
mac audio volume              # Show current volume
mac audio volume <0-100>      # Set volume level
mac audio volume up           # Increase volume
mac audio volume down         # Decrease volume
mac audio volume mute         # Mute audio
mac audio volume unmute       # Unmute audio
mac audio volume toggle-mute  # Toggle mute
mac audio devices             # List audio devices
mac audio input-device <name> # Switch input device
mac audio output-device <name> # Switch output device
mac audio status              # Show audio settings
```

**Examples:**
```bash
mac audio volume 50           # Set volume to 50%
mac audio volume up           # Increase by configured step (default: 10)
mac audio devices --json      # List devices as JSON
mac audio output-device "External Speakers"  # Switch output
```

---

### 📸 screenshot — Screenshot Settings & Management

Screenshot location, format, shadow, cursor settings.

```bash
mac screenshot location       # Show save location
mac screenshot location <path> # Set save location
mac screenshot format         # Show current format
mac screenshot format <png|jpg|tiff|gif>  # Set format
mac screenshot shadow on      # Include window shadow
mac screenshot shadow off     # Exclude window shadow
mac screenshot shadow toggle  # Toggle shadow
mac screenshot cursor show    # Show cursor in screenshots
mac screenshot cursor hide    # Hide cursor in screenshots
mac screenshot cursor toggle  # Toggle cursor visibility
mac screenshot naming sequential # Sequential naming (Screenshot 1, 2...)
mac screenshot naming timestamp # Timestamp naming
mac screenshot status         # Show all screenshot settings
mac screenshot reset          # Reset to defaults
```

**Examples:**
```bash
mac screenshot location ~/Pictures/Screenshots  # Custom save path
mac screenshot format jpg     # Use JPEG format
mac screenshot shadow off     # No shadow on window screenshots
mac screenshot status         # See all settings
```

---

## Configuration

SuperMac stores configuration in `~/.supermac/config.yaml`.

```yaml
version: 1

output:
  color: true
  format: text         # text | json | quiet

updates:
  check: true
  channel: stable      # stable | beta

modules:
  screenshot:
    location: Desktop
    format: PNG
    shadow: false
  audio:
    volume_step: 10
  display:
    brightness_step: 10

aliases:
  kp: "dev kill-port"
  dark: "display dark-mode"
  light: "display light-mode"
  ip: "network ip"
  cleanup: "system cleanup"
```

### Config Commands

```bash
mac config list               # Show all settings
mac config get output.format  # Get specific value
mac config set output.format json  # Set value
mac config edit               # Open in $EDITOR
```

---

## Shell Completions

```bash
# Zsh
mac completion zsh > ~/.zfunc/_mac

# Bash
mac completion bash > /etc/bash_completion.d/mac

# Fish
mac completion fish > ~/.config/fish/completions/mac.fish
```

---

## JSON Output

Every command supports `--json` for scripting:

```bash
mac system memory --json
mac wifi status --json
mac audio devices --json
mac network ip --json
```

Example output:
```json
{
  "type": "success",
  "message": "Memory: 16 GB total, 8 GB free"
}
```

---

## Project

- **License**: MIT
- **Author**: CosmoLabs
- **Repository**: https://github.com/cosmolabs-org/supermac
- **Website**: https://cosmolabs.org
