#!/bin/bash

# =============================================================================
# SuperMac v2.1.0 - Installation Script
# =============================================================================
# Automatically sets up the SuperMac modular system
# 
# Usage: 
#   curl -fsSL https://raw.githubusercontent.com/CosmoLabs-org/SuperMac/main/install.sh | bash
#   or: bash install.sh
#
# Built by CosmoLabs - https://cosmolabs.org
# =============================================================================

set -e  # Exit on any error

VERSION="2.1.0"
REPO_URL="https://raw.githubusercontent.com/CosmoLabs-org/SuperMac/main"
REPO_BASE="https://github.com/CosmoLabs-org/SuperMac"

# Colors for output
if [[ -t 1 ]]; then
    RED='\033[0;31m'
    GREEN='\033[0;32m'
    YELLOW='\033[1;33m'
    BLUE='\033[0;34m'
    PURPLE='\033[0;35m'
    BOLD='\033[1m'
    NC='\033[0m'
else
    RED='' GREEN='' YELLOW='' BLUE='' PURPLE='' BOLD='' NC=''
fi

# Output functions
print_success() { echo -e "${GREEN}✓${NC} $1"; }
print_error() { echo -e "${RED}✗${NC} $1" >&2; }
print_info() { echo -e "${BLUE}ℹ${NC} $1"; }
print_warning() { echo -e "${YELLOW}⚠${NC} $1"; }
print_header() { echo -e "${BOLD}${BLUE}$1${NC}"; }

# Banner
show_banner() {
    echo -e "${PURPLE}"
    echo "╭─────────────────────────────────────────────────────────╮"
    echo "│                  🚀 SuperMac v$VERSION                    │"
    echo "│                                                         │"
    echo "│               Built by CosmoLabs                        │"
    echo "│             https://cosmolabs.org                       │"
    echo "╰─────────────────────────────────────────────────────────╯"
    echo -e "${NC}"
    echo ""
}

# Check if running on macOS
check_macos() {
    if [[ "$OSTYPE" != "darwin"* ]]; then
        print_error "This installer is for macOS only."
        print_info "Detected OS: $OSTYPE"
        exit 1
    fi
    
    # Check macOS version
    local macos_version
    macos_version=$(sw_vers -productVersion)
    print_info "macOS version: $macos_version"
    
    # Warn if older than 12.0
    if [[ "$(echo "$macos_version" | cut -d. -f1)" -lt "12" ]]; then
        print_warning "macOS 12.0+ recommended. Some features may not work on older versions."
        echo -n "Continue anyway? [y/N] "
        read -r response
        if [[ ! "$response" =~ ^[Yy]$ ]]; then
            print_info "Installation cancelled."
            exit 0
        fi
    fi
}

# Detect shell configuration file
detect_shell_config() {
    local shell_config=""
    
    if [[ "$SHELL" == *"zsh"* ]]; then
        shell_config="$HOME/.zshrc"
    elif [[ "$SHELL" == *"bash"* ]]; then
        if [[ -f "$HOME/.bash_profile" ]]; then
            shell_config="$HOME/.bash_profile"
        else
            shell_config="$HOME/.bashrc"
        fi
    else
        print_warning "Unknown shell: $SHELL - using ~/.profile"
        shell_config="$HOME/.profile"
    fi
    
    echo "$shell_config"
}

# Check if PATH is already configured
check_path_configured() {
    local shell_config="$1"
    
    if [[ -f "$shell_config" ]]; then
        grep -q 'export PATH="$HOME/bin:$PATH"' "$shell_config" 2>/dev/null || \
        grep -q 'export PATH=$HOME/bin:$PATH' "$shell_config" 2>/dev/null
    else
        return 1
    fi
}

# Download file with fallback
download_file() {
    local url="$1"
    local output="$2"
    local description="$3"
    
    print_info "Downloading $description..."
    
    if command -v curl >/dev/null 2>&1; then
        if ! curl -fsSL "$url" -o "$output" 2>/dev/null; then
            print_error "Failed to download $description from GitHub"
            return 1
        fi
    elif command -v wget >/dev/null 2>&1; then
        if ! wget -q "$url" -O "$output" 2>/dev/null; then
            print_error "Failed to download $description from GitHub"
            return 1
        fi
    else
        print_error "Neither curl nor wget found. Cannot download files."
        return 1
    fi
    
    return 0
}

# Create directory structure
setup_directories() {
    print_info "Setting up SuperMac directories..."
    
    # Create main directories
    mkdir -p "$HOME/bin"
    mkdir -p "$HOME/.supermac/lib"
    mkdir -p "$HOME/.supermac/config"
    
    print_success "Directory structure created"
}

# Download core files
download_supermac_files() {
    local base_url="$REPO_URL"
    
    # Core files to download
    local files=(
        "bin/mac:$HOME/bin/mac:Main SuperMac script"
        "lib/utils.sh:$HOME/.supermac/lib/utils.sh:Utilities library"
        "lib/finder.sh:$HOME/.supermac/lib/finder.sh:Finder module"
        "lib/display.sh:$HOME/.supermac/lib/display.sh:Display module"
        "lib/network.sh:$HOME/.supermac/lib/network.sh:Network module"
        "lib/system.sh:$HOME/.supermac/lib/system.sh:System module"
        "lib/dev.sh:$HOME/.supermac/lib/dev.sh:Developer module"
        "config/config.json:$HOME/.supermac/config/config.json:Configuration file"
    )
    
    for file_info in "${files[@]}"; do
        IFS=':' read -r remote_path local_path description <<< "$file_info"
        
        if download_file "$base_url/$remote_path" "$local_path" "$description"; then
            print_success "$description installed"
        else
            print_warning "Failed to download $description - using fallback"
            # Here we could include embedded fallback versions
        fi
    done
}

# Create embedded fallback (minimal version)
create_fallback_mac() {
    print_info "Creating fallback SuperMac installation..."
    
    cat > "$HOME/bin/mac" << 'EOF'
#!/bin/bash
# SuperMac v2.1.0 - Minimal Installation
# For full features, reinstall from: https://github.com/CosmoLabs-org/SuperMac

VERSION="2.1.0"

# Basic colors
if [[ -t 1 ]]; then
    RED='\033[0;31m'
    GREEN='\033[0;32m'
    BLUE='\033[0;34m'
    NC='\033[0m'
else
    RED='' GREEN='' BLUE='' NC=''
fi

print_success() { echo -e "${GREEN}✓${NC} $1"; }
print_error() { echo -e "${RED}✗${NC} $1" >&2; }
print_info() { echo -e "${BLUE}ℹ${NC} $1"; }

# Basic help
show_help() {
    echo "SuperMac v$VERSION - Minimal Installation"
    echo ""
    echo "Available commands:"
    echo "  restart-finder     Restart Finder"
    echo "  toggle-hidden      Toggle hidden files"
    echo "  help              Show this help"
    echo ""
    echo "For full SuperMac features, reinstall from:"
    echo "https://github.com/CosmoLabs-org/SuperMac"
}

# Basic commands
case "$1" in
    "restart-finder"|"rf")
        print_info "Restarting Finder..."
        killall Finder 2>/dev/null && print_success "Finder restarted!"
        ;;
    "toggle-hidden"|"th")
        current=$(defaults read com.apple.finder AppleShowAllFiles 2>/dev/null || echo "FALSE")
        if [[ "$current" == "TRUE" ]]; then
            defaults write com.apple.finder AppleShowAllFiles FALSE
            print_info "Hidden files are now HIDDEN"
        else
            defaults write com.apple.finder AppleShowAllFiles TRUE
            print_info "Hidden files are now VISIBLE"
        fi
        killall Finder 2>/dev/null
        ;;
    "help"|"-h"|"--help"|"")
        show_help
        ;;
    *)
        print_error "Unknown command: $1"
        echo "Use 'mac help' for available commands"
        ;;
esac
EOF
    
    chmod +x "$HOME/bin/mac"
    print_success "Fallback SuperMac created"
}

# Main installation function
install_supermac() {
    show_banner
    
    print_header "Installing SuperMac v$VERSION..."
    echo ""
    
    # Validate environment
    check_macos
    
    # Setup directories
    setup_directories
    
    # Get shell configuration file
    local shell_config
    shell_config=$(detect_shell_config)
    print_info "Using shell config: $shell_config"
    
    # Check if PATH is already configured
    if check_path_configured "$shell_config"; then
        print_success "PATH already configured in $shell_config"
    else
        print_info "Adding ~/bin to PATH..."
        echo '' >> "$shell_config"
        echo '# Added by SuperMac installer (CosmoLabs)' >> "$shell_config"
        echo 'export PATH="$HOME/bin:$PATH"' >> "$shell_config"
        print_success "PATH configuration added to $shell_config"
    fi
    
    # Try to download full SuperMac
    if download_supermac_files; then
        print_success "Full SuperMac installation completed!"
    else
        print_warning "Network download failed. Installing minimal version..."
        create_fallback_mac
    fi
    
    # Make scripts executable
    chmod +x "$HOME/bin/mac" 2>/dev/null
    chmod +x "$HOME/.supermac/lib/"*.sh 2>/dev/null
    
    # Update PATH for current session
    export PATH="$HOME/bin:$PATH"
    
    # Test installation
    print_info "Testing installation..."
    
    if command -v mac >/dev/null 2>&1; then
        print_success "Installation completed successfully!"
        echo ""
        print_header "🎉 SuperMac is ready to use!"
        echo ""
        echo -e "${BOLD}Try these commands:${NC}"
        echo "  mac help                      # Show all available categories"
        echo "  mac finder restart            # Restart Finder"
        echo "  mac display dark-mode         # Switch to dark mode"
        echo "  mac network ip                # Show IP address"
        echo "  mac system info               # Show system information"
        echo ""
        print_info "💡 Tip: Use 'mac help <category>' to see commands in each category"
        echo ""
        
        # Show version if available
        mac version 2>/dev/null || true
        
    else
        print_warning "Installation completed, but 'mac' command not immediately available."
        print_info "Please restart your terminal or run: source $shell_config"
        print_info "Then test with: mac help"
    fi
    
    echo ""
    print_header "🚀 What's Next?"
    echo "• Explore commands: mac help"
    echo "• Quick IP check: mac ip"
    echo "• System cleanup: mac cleanup"
    echo "• Dark mode: mac dark"
    echo "• Star the repo: $REPO_BASE"
    echo ""
    echo -e "${BOLD}${PURPLE}Built with ❤️ by CosmoLabs${NC}"
    echo -e "${BLUE}https://cosmolabs.org${NC}"
    echo ""
}

# Uninstall function
uninstall_supermac() {
    print_header "Uninstalling SuperMac..."
    echo ""
    
    # Remove files
    rm -f "$HOME/bin/mac"
    rm -rf "$HOME/.supermac"
    
    print_success "SuperMac files removed"
    print_info "PATH configuration left in shell config (safe to leave)"
    print_info "To complete removal, restart your terminal"
}

# Handle command line arguments
case "${1:-install}" in
    "install")
        install_supermac
        ;;
    "uninstall")
        uninstall_supermac
        ;;
    *)
        echo "Usage: $0 [install|uninstall]"
        exit 1
        ;;
esac
