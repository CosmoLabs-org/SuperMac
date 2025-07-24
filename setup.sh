#!/bin/bash

# =============================================================================
# SuperMac Setup & Demo Script
# =============================================================================
# Sets up SuperMac for development and demonstrates functionality
# 
# Usage: bash setup.sh [demo|test|install]
# 
# Built by CosmoLabs - https://cosmolabs.org
# =============================================================================

set -e

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SUPERMAC_ROOT="$(dirname "$SCRIPT_DIR")"
BIN_DIR="$SUPERMAC_ROOT/bin"
LIB_DIR="$SUPERMAC_ROOT/lib"

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

show_banner() {
    echo -e "${PURPLE}"
    echo "╭─────────────────────────────────────────────────────────╮"
    echo "│               🚀 SuperMac v2.1.0 Setup                │"
    echo "│                                                         │"
    echo "│              Built by CosmoLabs                         │"
    echo "│            https://cosmolabs.org                        │"
    echo "╰─────────────────────────────────────────────────────────╯"
    echo -e "${NC}"
    echo ""
}

# =============================================================================
# Setup Functions
# =============================================================================

setup_permissions() {
    print_header "Setting up file permissions..."
    
    # Make main script executable
    chmod +x "$BIN_DIR/mac" 2>/dev/null || print_warning "Could not make mac script executable"
    
    # Make all library files readable
    chmod +r "$LIB_DIR"/*.sh 2>/dev/null || print_warning "Could not set library permissions"
    
    # Make test script executable
    chmod +x "$SUPERMAC_ROOT/tests/test.sh" 2>/dev/null || print_warning "Could not make test script executable"
    
    print_success "Permissions configured"
}

validate_structure() {
    print_header "Validating project structure..."
    
    local all_good=true
    
    # Check main directories
    for dir in "bin" "lib" "config" "docs" "tests"; do
        if [[ -d "$SUPERMAC_ROOT/$dir" ]]; then
            print_success "Directory $dir/ exists"
        else
            print_error "Directory $dir/ missing"
            all_good=false
        fi
    done
    
    # Check main script
    if [[ -f "$BIN_DIR/mac" ]]; then
        print_success "Main script bin/mac exists"
    else
        print_error "Main script bin/mac missing"
        all_good=false
    fi
    
    # Check core modules
    local modules=("utils" "finder" "wifi" "network" "system" "dev" "display" "dock" "audio")
    for module in "${modules[@]}"; do
        if [[ -f "$LIB_DIR/$module.sh" ]]; then
            print_success "Module lib/$module.sh exists"
        else
            print_error "Module lib/$module.sh missing"
            all_good=false
        fi
    done
    
    if [[ "$all_good" == true ]]; then
        print_success "Project structure is valid!"
        return 0
    else
        print_error "Project structure has issues"
        return 1
    fi
}

test_syntax() {
    print_header "Testing script syntax..."
    
    # Test main script
    if bash -n "$BIN_DIR/mac" 2>/dev/null; then
        print_success "Main script syntax is valid"
    else
        print_error "Main script has syntax errors"
        return 1
    fi
    
    # Test all modules
    for module_file in "$LIB_DIR"/*.sh; do
        if [[ -f "$module_file" ]]; then
            local module_name
            module_name=$(basename "$module_file" .sh)
            if bash -n "$module_file" 2>/dev/null; then
                print_success "Module $module_name syntax is valid"
            else
                print_error "Module $module_name has syntax errors"
                return 1
            fi
        fi
    done
    
    print_success "All syntax checks passed!"
}

test_basic_functionality() {
    print_header "Testing basic functionality..."
    
    # Test help command
    if "$BIN_DIR/mac" help >/dev/null 2>&1; then
        print_success "Help command works"
    else
        print_error "Help command failed"
        return 1
    fi
    
    # Test version command
    if "$BIN_DIR/mac" version >/dev/null 2>&1; then
        print_success "Version command works"
    else
        print_error "Version command failed"
        return 1
    fi
    
    # Test module help
    if "$BIN_DIR/mac" help finder >/dev/null 2>&1; then
        print_success "Module help works"
    else
        print_error "Module help failed"
        return 1
    fi
    
    # Test safe commands
    local safe_commands=(
        "finder status"
        "network ip"
        "system info"
        "display status"
        "audio status"
        "dock status"
    )
    
    for cmd in "${safe_commands[@]}"; do
        if timeout 10 "$BIN_DIR/mac" $cmd >/dev/null 2>&1; then
            print_success "Command 'mac $cmd' works"
        else
            print_warning "Command 'mac $cmd' had issues (may be expected)"
        fi
    done
    
    print_success "Basic functionality tests completed!"
}

# =============================================================================
# Demo Functions
# =============================================================================

demo_help_system() {
    print_header "🎨 Beautiful Help System Demo"
    echo ""
    
    print_info "Showing main help..."
    "$BIN_DIR/mac" help
    
    echo ""
    print_info "Press Enter to see category-specific help..."
    read -r
    
    print_info "Showing finder help..."
    "$BIN_DIR/mac" help finder
    
    echo ""
    print_info "Press Enter to see search functionality..."
    read -r
    
    print_info "Searching for 'network' commands..."
    "$BIN_DIR/mac" search network
}

demo_safe_commands() {
    print_header "🔍 Safe Command Demonstrations"
    echo ""
    
    local demos=(
        "system info:Show comprehensive system information"
        "network ip:Display current IP address"
        "display status:Show display settings"
        "audio status:Show audio configuration"
        "dock status:Show dock settings"
        "finder status:Show Finder configuration"
        "dev processes:Show running processes"
    )
    
    for demo in "${demos[@]}"; do
        IFS=':' read -r command description <<< "$demo"
        
        print_info "$description"
        echo "  Command: mac $command"
        echo ""
        
        if timeout 15 "$BIN_DIR/mac" $command 2>/dev/null; then
            echo ""
        else
            print_warning "Command timed out or had issues"
            echo ""
        fi
        
        echo "Press Enter for next demo..."
        read -r
    done
}

demo_search_system() {
    print_header "🔍 Search System Demo"
    echo ""
    
    local search_terms=("volume" "network" "dark" "finder" "port")
    
    for term in "${search_terms[@]}"; do
        print_info "Searching for: $term"
        "$BIN_DIR/mac" search "$term" 2>/dev/null || print_warning "Search failed for $term"
        echo ""
        echo "Press Enter for next search..."
        read -r
    done
}

# =============================================================================
# Installation Functions
# =============================================================================

install_to_system() {
    print_header "Installing SuperMac to system..."
    
    # Create ~/bin if it doesn't exist
    mkdir -p "$HOME/bin"
    
    # Copy main script
    cp "$BIN_DIR/mac" "$HOME/bin/mac"
    chmod +x "$HOME/bin/mac"
    
    # Create SuperMac directory in home
    mkdir -p "$HOME/.supermac/lib"
    mkdir -p "$HOME/.supermac/config"
    
    # Copy libraries
    cp "$LIB_DIR"/*.sh "$HOME/.supermac/lib/"
    
    # Copy configuration
    cp "$SUPERMAC_ROOT/config/config.json" "$HOME/.supermac/config/"
    
    # Check if ~/bin is in PATH
    if echo "$PATH" | grep -q "$HOME/bin"; then
        print_success "SuperMac installed successfully!"
        print_info "You can now use 'mac' command from anywhere"
    else
        print_success "SuperMac installed successfully!"
        print_warning "~/bin is not in your PATH"
        print_info "Add this line to your ~/.zshrc or ~/.bash_profile:"
        print_info "export PATH=\"\$HOME/bin:\$PATH\""
    fi
    
    # Test installation
    if command -v mac >/dev/null 2>&1; then
        print_success "Installation verified - 'mac' command is available"
        echo ""
        print_info "Try these commands:"
        echo "  mac help"
        echo "  mac system info"
        echo "  mac network ip"
    else
        print_warning "Installation complete but 'mac' command not in PATH"
        print_info "Restart your terminal or source your shell config file"
    fi
}

# =============================================================================
# Main Functions
# =============================================================================

run_setup() {
    show_banner
    
    print_header "🔧 SuperMac Development Setup"
    echo ""
    
    setup_permissions
    echo ""
    
    validate_structure || exit 1
    echo ""
    
    test_syntax || exit 1
    echo ""
    
    test_basic_functionality
    echo ""
    
    print_success "Setup completed successfully!"
    echo ""
    print_info "Next steps:"
    echo "  • Run 'bash setup.sh demo' to see SuperMac in action"
    echo "  • Run 'bash setup.sh test' to run full test suite"
    echo "  • Run 'bash setup.sh install' to install to your system"
    echo "  • Check docs/DEVELOPMENT.md for development guide"
}

run_demo() {
    show_banner
    
    print_header "🎬 SuperMac Interactive Demo"
    echo ""
    
    print_info "This demo will showcase SuperMac's capabilities"
    print_info "Press Enter to continue, Ctrl+C to exit"
    read -r
    
    demo_help_system
    echo ""
    echo "Press Enter to continue to command demonstrations..."
    read -r
    
    demo_safe_commands
    echo ""
    echo "Press Enter to see search system..."
    read -r
    
    demo_search_system
    echo ""
    print_success "Demo completed!"
    print_info "Try 'bash setup.sh install' to install SuperMac to your system"
}

run_tests() {
    show_banner
    
    print_header "🧪 Running SuperMac Test Suite"
    echo ""
    
    if [[ -x "$SUPERMAC_ROOT/tests/test.sh" ]]; then
        "$SUPERMAC_ROOT/tests/test.sh"
    else
        print_error "Test script not found or not executable"
        exit 1
    fi
}

show_help() {
    show_banner
    
    echo "Usage: bash setup.sh [command]"
    echo ""
    echo "Commands:"
    echo "  setup     Set up development environment (default)"
    echo "  demo      Interactive demonstration"
    echo "  test      Run test suite"
    echo "  install   Install to system (~/.bin)"
    echo "  help      Show this help"
    echo ""
    echo "Examples:"
    echo "  bash setup.sh           # Setup development environment"
    echo "  bash setup.sh demo      # See SuperMac in action"
    echo "  bash setup.sh install   # Install to your system"
}

# =============================================================================
# Main
# =============================================================================

main() {
    local command="${1:-setup}"
    
    case "$command" in
        "setup"|"")
            run_setup
            ;;
        "demo")
            run_demo
            ;;
        "test")
            run_tests
            ;;
        "install")
            install_to_system
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            print_error "Unknown command: $command"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

# Run main function
main "$@"
