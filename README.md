# 🚀 SuperMac v2.1.0 - Professional macOS CLI Tool

<div align="center">

![macOS](https://img.shields.io/badge/macOS-12.0+-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![Version](https://img.shields.io/badge/version-2.1.0-orange.svg)
![Architecture](https://img.shields.io/badge/architecture-Intel%20%26%20Apple%20Silicon-purple.svg)

**Organized, Powerful, Professional**

Built by [CosmoLabs](https://cosmolabs.org) for macOS developers and power users

![SuperMac](https://i.imgur.com/HU2yib3.jpeg)

</div>

## ✨ What's New in v2.1.0

🎨 **Beautiful Help System** - Stunning terminal output with box drawing and colors  
🧩 **Modular Architecture** - Clean, maintainable codebase split into focused modules  
🖥️ **Display Commands** - Brightness, dark mode, Night Shift, and True Tone control  
🎯 **Enhanced UX** - Better error messages, input validation, and user feedback  
🔧 **Configuration System** - Customizable settings and user preferences  
🔍 **Smart Search** - Find commands quickly with `mac search <term>`  

## 🚀 Quick Install

```bash
# One-line installation
curl -fsSL https://raw.githubusercontent.com/CosmoLabs-org/SuperMac/main/install.sh | bash

# Test installation
mac help
mac finder restart
```

## 🎯 Command Structure

SuperMac uses an organized, hierarchical command structure:

```bash
mac <category> <action> [arguments]
```

### 📋 Available Categories

| Category | Icon | Description | Example Commands |
|----------|------|-------------|------------------|
| **finder** | 📁 | File visibility and Finder management | `restart`, `toggle-hidden`, `reveal` |
| **wifi** | 🌐 | WiFi control and management | `on`, `off`, `toggle`, `status` |
| **network** | 📡 | Network information and troubleshooting | `ip`, `public-ip`, `flush-dns` |
| **system** | 🖥️ | System information and maintenance | `info`, `cleanup`, `battery` |
| **dev** | 💻 | Developer tools and utilities | `kill-port`, `servers`, `localhost` |
| **display** | 🖥️ | Display and appearance settings | `brightness`, `dark-mode`, `night-shift` |
| **dock** | 🚢 | Dock management and positioning | `position`, `autohide`, `size` |
| **audio** | 🔊 | Audio control and device management | `volume`, `mute`, `devices` |

## 🎨 Beautiful Help System

SuperMac features a stunning help system with visual hierarchy:

```bash
╭─────────────────────────────────────────────────────────╮
│                  🚀 SuperMac v2.1.0                    │
│                Built by CosmoLabs                       │
╰─────────────────────────────────────────────────────────╯

┌─ 📁 FINDER ──────────────────────────────────────────────┐
│  restart               Restart Finder application       │
│  toggle-hidden         Toggle hidden file visibility    │
│  show-hidden           Show system files               │
│  hide-hidden           Hide system files               │
│  reveal <path>         Reveal file/folder in Finder    │
└──────────────────────────────────────────────────────────┘
```

## ⚡ Quick Commands

### Global Shortcuts
```bash
mac ip                    # Quick IP lookup
mac cleanup               # System cleanup
mac dark                  # Switch to dark mode
mac light                 # Switch to light mode
mac kp 3000               # Kill process on port 3000
```

### Essential Commands
```bash
# Finder Management
mac finder restart                # Fix unresponsive Finder
mac finder toggle-hidden          # Show/hide system files
mac finder reveal ~/.ssh          # Open .ssh folder in Finder

# Display Control
mac display brightness 75         # Set brightness to 75%
mac display dark-mode             # Switch to dark mode
mac display night-shift on        # Enable Night Shift

# Development
mac dev kill-port 3000            # Kill process on port
mac dev servers                   # List running servers
mac dev localhost 8080            # Open localhost in browser

# Network
mac network ip                    # Show local IP
mac network public-ip             # Show public IP
mac network flush-dns             # Clear DNS cache

# System Maintenance
mac system info                   # System information
mac system cleanup                # Deep cleanup
mac system battery                # Battery status
```

## 🔍 Discovery Features

### Smart Help System
```bash
mac help                          # Show all categories
mac help display                  # Show display commands
mac search brightness             # Find brightness-related commands
```

### Command Search
```bash
mac search "night"                # Find all night-related commands
mac search "wifi"                 # Find all WiFi commands
mac search "dark"                 # Find dark mode commands
```

## 🛠️ Advanced Features

### Configuration System
```bash
# Customize command name
mac config set command-name "supermac"

# Set default screenshot location
mac config set screenshot-location "Downloads"

# Configure volume step
mac config set volume-step 15
```

### Debug Mode
```bash
mac --debug system info           # Run with detailed debugging
mac debug                        # Enable debug mode
```

### Module Information
```bash
mac modules list                  # Show available modules
mac modules info finder           # Show module details
```

## 📊 Performance & Compatibility

- **Startup Time**: < 0.5 seconds
- **Memory Usage**: < 5MB
- **Dependencies**: None (pure bash + macOS APIs)
- **Compatibility**: macOS 12.0+ (Intel & Apple Silicon)
- **Terminals**: Terminal.app, iTerm2, and all standard terminals

## 🏗️ Architecture

SuperMac v2.1.0 features a clean, modular architecture:

```
SuperMac/
├── bin/
│   ├── mac                 # Main dispatcher (entry point)
│   └── install.sh          # Installation script
├── lib/                    # Modular command libraries
│   ├── utils.sh           # Shared utilities & formatting
│   ├── finder.sh          # Finder commands
│   ├── display.sh         # Display commands
│   ├── network.sh         # Network commands
│   ├── system.sh          # System commands
│   └── dev.sh             # Developer commands
├── config/
│   └── config.json        # User configuration
└── docs/
    ├── README.md          # This file
    ├── commands.md        # Command reference
    └── development.md     # Development guide
```

## 🚀 Development Workflow Examples

### Frontend Developer
```bash
# Start development
mac dev servers                   # Check running services
mac dev kill-port 3000           # Clean up old process
npm start                        # Start your app
mac dev localhost 3000           # Open in browser

# Appearance work
mac display dark-mode             # Test dark theme
mac display light-mode            # Test light theme
mac display brightness 50        # Adjust for design work
```

### System Administration
```bash
# Daily maintenance
mac system info                   # Check system status
mac system cleanup                # Clean temporary files
mac network flush-dns             # Refresh network
mac finder restart                # Fix any Finder issues

# Network troubleshooting
mac network info                  # Check connectivity
mac wifi status                   # Check WiFi details
mac network public-ip             # Verify external access
```

### Development Environment Setup
```bash
# Clean slate
mac system cleanup                # Clear caches
mac dev kill-port 3000 8080 9000 # Clean up ports
mac finder toggle-hidden          # Show config files

# Optimize display
mac display brightness 70         # Comfortable brightness
mac display night-shift off       # Accurate colors
mac display dark-mode             # Reduce eye strain
```

## 🔧 Installation Options

### Automatic Installation (Recommended)
```bash
curl -fsSL https://raw.githubusercontent.com/CosmoLabs-org/SuperMac/main/install.sh | bash
```

### Manual Installation
```bash
# Clone repository
git clone https://github.com/CosmoLabs-org/SuperMac.git
cd SuperMac

# Run installer
chmod +x install.sh
./install.sh

# Test installation
mac help
```

### Development Installation
```bash
# Clone for development
git clone https://github.com/CosmoLabs-org/SuperMac.git
cd SuperMac

# Link to development version
ln -sf "$(pwd)/bin/mac" ~/bin/mac-dev

# Test development version
mac-dev help
```

## 🎯 Use Cases

### For Developers
- **Quick port management** - Kill processes instantly
- **Localhost testing** - Open development servers
- **Environment setup** - Configure display and system settings
- **System maintenance** - Keep development machine clean

### For Power Users
- **System optimization** - Regular cleanup and maintenance
- **Display management** - Perfect brightness and appearance
- **Network troubleshooting** - Quick connectivity diagnostics
- **File management** - Advanced Finder operations

### For System Administrators
- **Bulk operations** - Manage multiple settings quickly
- **Troubleshooting** - Comprehensive system information
- **Automation** - Script-friendly command structure
- **Monitoring** - System status and health checks

## 🤝 Contributing

We welcome contributions! Here's how to help:

### Adding New Commands
1. Create/modify module in `lib/`
2. Follow the existing pattern
3. Add help and search functions
4. Update configuration
5. Test thoroughly

### Reporting Issues
- Use GitHub Issues with detailed information
- Include macOS version and system details
- Provide steps to reproduce

### Feature Requests
- Check existing issues first
- Describe the use case clearly
- Consider implementation complexity

## 📄 License

MIT License - see [LICENSE](../LICENSE) file for details.

## 🏢 About CosmoLabs

SuperMac is built by **CosmoLabs**, creating tools that make developers and power users more productive on macOS.

- 🌐 Website: [cosmolabs.org](https://cosmolabs.org)
- 🐦 Twitter: [@CosmoLabsHQ](https://twitter.com/CosmoLabsHQ)
- 📧 Contact: hello@cosmolabs.org
- 🔗 GitHub: [github.com/cosmolabs-org](https://github.com/cosmolabs-org)

## ⭐ Show Your Support

If SuperMac saves you time and improves your workflow:
- ⭐ Star this repository
- 🐛 Report issues or request features
- 🤝 Contribute new commands
- 📢 Share with other macOS users
- ☕ [Buy us a coffee](https://buymeacoffee.com/cosmolabs)

## 📈 Roadmap

### v2.2.0 - Security & Privacy
- Privacy controls and permissions
- Security settings management
- Keychain operations

### v2.3.0 - File Operations
- Advanced file management
- Compression utilities
- Permission management

### v2.4.0 - Cloud & Backup
- iCloud management
- Time Machine controls
- Cloud service integration

### v2.5.0 - Automation
- Command scripting
- Workflow automation
- Plugin system

---

<div align="center">

**Built with ❤️ for the macOS community by CosmoLabs**

*SuperMac - Organized, Powerful, Professional*

</div>
