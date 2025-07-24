#!/bin/bash

# =============================================================================
# SuperMac Shared Utilities
# =============================================================================
# Common functions, colors, and utilities used across all SuperMac modules
# 
# Built by CosmoLabs - https://cosmolabs.org
# =============================================================================

# =============================================================================
# Version & Constants
# =============================================================================
readonly SUPERMAC_VERSION="2.1.0"
readonly SUPERMAC_NAME="SuperMac"
readonly SUPERMAC_AUTHOR="CosmoLabs"
readonly SUPERMAC_URL="https://cosmolabs.org"
readonly SUPERMAC_REPO="https://github.com/CosmoLabs-org/SuperMac"

# =============================================================================
# Color Configuration
# =============================================================================
# Only use colors if terminal supports them
if [[ -t 1 ]]; then
    # Standard colors
    readonly RED='\033[0;31m'
    readonly GREEN='\033[0;32m'
    readonly YELLOW='\033[1;33m'
    readonly BLUE='\033[0;34m'
    readonly PURPLE='\033[0;35m'
    readonly CYAN='\033[0;36m'
    readonly WHITE='\033[1;37m'
    readonly GRAY='\033[0;90m'
    
    # Text formatting
    readonly BOLD='\033[1m'
    readonly DIM='\033[2m'
    readonly UNDERLINE='\033[4m'
    readonly BLINK='\033[5m'
    readonly REVERSE='\033[7m'
    readonly NC='\033[0m' # No Color
    
    # Background colors
    readonly BG_RED='\033[41m'
    readonly BG_GREEN='\033[42m'
    readonly BG_YELLOW='\033[43m'
    readonly BG_BLUE='\033[44m'
    readonly BG_PURPLE='\033[45m'
    readonly BG_CYAN='\033[46m'
    readonly BG_WHITE='\033[47m'
    
    # Box drawing characters for beautiful output
    readonly BOX_TOP_LEFT='╭'
    readonly BOX_TOP_RIGHT='╮'
    readonly BOX_BOTTOM_LEFT='╰'
    readonly BOX_BOTTOM_RIGHT='╯'
    readonly BOX_HORIZONTAL='─'
    readonly BOX_VERTICAL='│'
    readonly BOX_CROSS='┼'
    readonly BOX_TOP_TEE='┬'
    readonly BOX_BOTTOM_TEE='┴'
    readonly BOX_LEFT_TEE='├'
    readonly BOX_RIGHT_TEE='┤'
    
    # Simple box drawing
    readonly SIMPLE_TOP_LEFT='┌'
    readonly SIMPLE_TOP_RIGHT='┐'
    readonly SIMPLE_BOTTOM_LEFT='└'
    readonly SIMPLE_BOTTOM_RIGHT='┘'
    readonly SIMPLE_HORIZONTAL='─'
    readonly SIMPLE_VERTICAL='│'
    
else
    # No color support - disable all formatting
    readonly RED='' GREEN='' YELLOW='' BLUE='' PURPLE='' CYAN='' WHITE='' GRAY=''
    readonly BOLD='' DIM='' UNDERLINE='' BLINK='' REVERSE='' NC=''
    readonly BG_RED='' BG_GREEN='' BG_YELLOW='' BG_BLUE='' BG_PURPLE='' BG_CYAN='' BG_WHITE=''
    readonly BOX_TOP_LEFT='+' BOX_TOP_RIGHT='+' BOX_BOTTOM_LEFT='+' BOX_BOTTOM_RIGHT='+'
    readonly BOX_HORIZONTAL='-' BOX_VERTICAL='|' BOX_CROSS='+' BOX_TOP_TEE='+' BOX_BOTTOM_TEE='+'
    readonly BOX_LEFT_TEE='+' BOX_RIGHT_TEE='+'
    readonly SIMPLE_TOP_LEFT='+' SIMPLE_TOP_RIGHT='+' SIMPLE_BOTTOM_LEFT='+' SIMPLE_BOTTOM_RIGHT='+'
    readonly SIMPLE_HORIZONTAL='-' SIMPLE_VERTICAL='|'
fi

# =============================================================================
# Output Functions
# =============================================================================

# Basic output functions with consistent formatting
print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1" >&2
}

print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_debug() {
    if [[ "${SUPERMAC_DEBUG:-}" == "1" ]]; then
        echo -e "${GRAY}[DEBUG]${NC} $1" >&2
    fi
}

print_header() {
    echo -e "${BOLD}${BLUE}$1${NC}"
}

print_subheader() {
    echo -e "${BOLD}$1${NC}"
}

print_category() {
    echo -e "${PURPLE}$1${NC}"
}

print_command() {
    echo -e "${CYAN}$1${NC}"
}

print_dim() {
    echo -e "${DIM}$1${NC}"
}

# =============================================================================
# Advanced Output Functions
# =============================================================================

# Print a horizontal line
print_line() {
    local length=${1:-60}
    local char=${2:-${BOX_HORIZONTAL}}
    printf "%*s\n" "$length" "" | tr ' ' "$char"
}

# Print a centered text in a box
print_banner() {
    local text="$1"
    local width=${2:-60}
    local padding=$(( (width - ${#text} - 2) / 2 ))
    
    echo -e "${BOLD}${BLUE}"
    printf "%s" "$BOX_TOP_LEFT"
    printf "%*s" "$((width-2))" "" | tr ' ' "$BOX_HORIZONTAL"
    printf "%s\n" "$BOX_TOP_RIGHT"
    
    printf "%s" "$BOX_VERTICAL"
    printf "%*s" "$padding" ""
    printf "%s" "$text"
    printf "%*s" "$((width - padding - ${#text} - 2))" ""
    printf "%s\n" "$BOX_VERTICAL"
    
    printf "%s" "$BOX_BOTTOM_LEFT"
    printf "%*s" "$((width-2))" "" | tr ' ' "$BOX_HORIZONTAL"
    printf "%s\n" "$BOX_BOTTOM_RIGHT"
    echo -e "${NC}"
}

# Print a simple box around text
print_box() {
    local text="$1"
    local color="${2:-$BLUE}"
    local width=${3:-$((${#text} + 4))}
    
    echo -e "${color}"
    printf "%s" "$SIMPLE_TOP_LEFT"
    printf "%*s" "$((width-2))" "" | tr ' ' "$SIMPLE_HORIZONTAL"
    printf "%s\n" "$SIMPLE_TOP_RIGHT"
    
    printf "%s %s" "$SIMPLE_VERTICAL" "$text"
    printf "%*s" "$((width - ${#text} - 3))" ""
    printf "%s\n" "$SIMPLE_VERTICAL"
    
    printf "%s" "$SIMPLE_BOTTOM_LEFT"
    printf "%*s" "$((width-2))" "" | tr ' ' "$SIMPLE_HORIZONTAL"
    printf "%s\n" "$SIMPLE_BOTTOM_RIGHT"
    echo -e "${NC}"
}

# Progress bar function
print_progress() {
    local current=$1
    local total=$2
    local width=${3:-50}
    local fill_char="█"
    local empty_char="░"
    
    local filled=$(( current * width / total ))
    local empty=$(( width - filled ))
    
    printf "\r${BLUE}["
    printf "%*s" "$filled" "" | tr ' ' "$fill_char"
    printf "%*s" "$empty" "" | tr ' ' "$empty_char"
    printf "] %d%% (%d/%d)${NC}" "$((current * 100 / total))" "$current" "$total"
    
    if [[ $current -eq $total ]]; then
        echo
    fi
}

# =============================================================================
# Validation Functions
# =============================================================================

# Check if running on macOS
check_macos() {
    if [[ "$OSTYPE" != "darwin"* ]]; then
        print_error "SuperMac is designed for macOS only."
        print_info "Detected OS: $OSTYPE"
        return 1
    fi
}

# Check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check if a number is valid
is_number() {
    [[ "$1" =~ ^[0-9]+$ ]]
}

# Check if a number is in range
is_in_range() {
    local num=$1
    local min=$2
    local max=$3
    
    is_number "$num" && [[ $num -ge $min ]] && [[ $num -le $max ]]
}

# Validate file exists
file_exists() {
    [[ -f "$1" ]]
}

# Validate directory exists
dir_exists() {
    [[ -d "$1" ]]
}

# =============================================================================
# System Information Functions
# =============================================================================

# Get macOS version
get_macos_version() {
    sw_vers -productVersion
}

# Get macOS build version
get_macos_build() {
    sw_vers -buildVersion
}

# Get system architecture
get_architecture() {
    uname -m
}

# Check if running on Apple Silicon
is_apple_silicon() {
    [[ "$(uname -m)" == "arm64" ]]
}

# Get current shell
get_shell() {
    basename "$SHELL"
}

# =============================================================================
# User Interaction Functions
# =============================================================================

# Ask user for confirmation
confirm() {
    local prompt="${1:-Are you sure?}"
    local default="${2:-n}"
    
    if [[ "$default" == "y" ]]; then
        prompt="$prompt [Y/n]"
    else
        prompt="$prompt [y/N]"
    fi
    
    echo -e "${YELLOW}$prompt${NC} "
    read -r response
    
    case "$response" in
        [yY]|[yY][eE][sS])
            return 0
            ;;
        [nN]|[nN][oO])
            return 1
            ;;
        "")
            [[ "$default" == "y" ]] && return 0 || return 1
            ;;
        *)
            print_warning "Please answer yes or no."
            confirm "$1" "$default"
            ;;
    esac
}

# Prompt for input with validation
prompt_input() {
    local prompt="$1"
    local validator="$2"  # Optional validation function
    local default="$3"    # Optional default value
    
    while true; do
        echo -e "${BLUE}$prompt${NC}"
        if [[ -n "$default" ]]; then
            echo -e "${DIM}(default: $default)${NC}"
        fi
        echo -n "> "
        read -r input
        
        # Use default if no input provided
        if [[ -z "$input" && -n "$default" ]]; then
            input="$default"
        fi
        
        # Validate input if validator function provided
        if [[ -n "$validator" ]]; then
            if "$validator" "$input"; then
                echo "$input"
                return 0
            else
                print_error "Invalid input. Please try again."
                continue
            fi
        else
            echo "$input"
            return 0
        fi
    done
}

# =============================================================================
# String Manipulation Functions
# =============================================================================

# Convert string to lowercase
to_lower() {
    echo "$1" | tr '[:upper:]' '[:lower:]'
}

# Convert string to uppercase
to_upper() {
    echo "$1" | tr '[:lower:]' '[:upper:]'
}

# Trim whitespace from string
trim() {
    local str="$1"
    str="${str#"${str%%[![:space:]]*}"}"  # Remove leading whitespace
    str="${str%"${str##*[![:space:]]}"}"  # Remove trailing whitespace
    echo "$str"
}

# Pad string to specified length
pad_string() {
    local str="$1"
    local length="$2"
    local char="${3:- }"
    
    printf "%-*s" "$length" "$str" | tr ' ' "$char"
}

# =============================================================================
# Array Functions
# =============================================================================

# Check if array contains element
array_contains() {
    local element="$1"
    shift
    local array=("$@")
    
    for item in "${array[@]}"; do
        [[ "$item" == "$element" ]] && return 0
    done
    return 1
}

# Join array elements with delimiter
array_join() {
    local delimiter="$1"
    shift
    local array=("$@")
    
    local result=""
    for item in "${array[@]}"; do
        if [[ -z "$result" ]]; then
            result="$item"
        else
            result="$result$delimiter$item"
        fi
    done
    echo "$result"
}

# =============================================================================
# Configuration Functions
# =============================================================================

# Get SuperMac installation directory
get_supermac_dir() {
    local script_path="$(readlink -f "${BASH_SOURCE[0]}" 2>/dev/null || echo "${BASH_SOURCE[0]}")"
    dirname "$(dirname "$script_path")"
}

# Get configuration directory
get_config_dir() {
    echo "$(get_supermac_dir)/config"
}

# Load configuration value
get_config() {
    local key="$1"
    local default="$2"
    local config_file="$(get_config_dir)/config.json"
    
    if file_exists "$config_file" && command_exists jq; then
        local value
        value=$(jq -r ".$key // \"$default\"" "$config_file" 2>/dev/null)
        echo "${value:-$default}"
    else
        echo "$default"
    fi
}

# =============================================================================
# Error Handling
# =============================================================================

# Exit with error message
die() {
    print_error "$1"
    exit "${2:-1}"
}

# Trap function for cleanup
cleanup() {
    # Add any cleanup tasks here
    print_debug "Cleaning up..."
}

# Set up error handling
setup_error_handling() {
    set -e  # Exit on error
    trap cleanup EXIT
}

# =============================================================================
# Debugging Functions
# =============================================================================

# Enable debug mode
enable_debug() {
    export SUPERMAC_DEBUG=1
    print_debug "Debug mode enabled"
}

# Print stack trace
print_stack_trace() {
    print_error "Stack trace:"
    local frame=0
    while caller $frame; do
        ((frame++))
    done
}

# =============================================================================
# Performance Functions
# =============================================================================

# Time a command execution
time_command() {
    local start_time
    start_time=$(date +%s.%N)
    
    "$@"
    local exit_code=$?
    
    local end_time
    end_time=$(date +%s.%N)
    local duration
    duration=$(echo "$end_time - $start_time" | bc 2>/dev/null || echo "0")
    
    print_debug "Command completed in ${duration}s"
    return $exit_code
}

# =============================================================================
# Module Loading
# =============================================================================

# Load a SuperMac module
load_module() {
    local module_name="$1"
    local module_path="$(get_supermac_dir)/lib/${module_name}.sh"
    
    if file_exists "$module_path"; then
        print_debug "Loading module: $module_name"
        # shellcheck source=/dev/null
        source "$module_path"
    else
        die "Module not found: $module_name ($module_path)"
    fi
}

# Get list of available modules
get_available_modules() {
    local lib_dir="$(get_supermac_dir)/lib"
    if dir_exists "$lib_dir"; then
        find "$lib_dir" -name "*.sh" -not -name "utils.sh" -exec basename {} .sh \; | sort
    fi
}

# =============================================================================
# Help System Functions
# =============================================================================

# Print SuperMac header
print_supermac_header() {
    echo -e "${BOLD}${PURPLE}"
    echo "╭─────────────────────────────────────────────────────────╮"
    echo "│                  🚀 SuperMac v$SUPERMAC_VERSION                    │"
    echo "│                Built by CosmoLabs                       │"
    echo "╰─────────────────────────────────────────────────────────╯"
    echo -e "${NC}"
}

# Print category header
print_category_header() {
    local category="$1"
    local emoji="$2"
    local width=${3:-60}
    
    echo -e "${PURPLE}"
    printf "┌─ %s %s " "$emoji" "$(to_upper "$category")"
    local header_len=$((${#category} + ${#emoji} + 4))
    local remaining=$((width - header_len - 1))
    printf "%*s" "$remaining" "" | tr ' ' '─'
    printf "┐\n"
    echo -e "${NC}"
}

# Print category footer
print_category_footer() {
    local width=${1:-60}
    echo -e "${PURPLE}"
    printf "└"
    printf "%*s" "$((width - 2))" "" | tr ' ' '─'
    printf "┘\n"
    echo -e "${NC}"
}

# =============================================================================
# Export Functions for Global Use
# =============================================================================

# Make all functions available to other scripts
export -f print_success print_error print_info print_warning print_debug
export -f print_header print_subheader print_category print_command print_dim
export -f print_line print_banner print_box print_progress
export -f check_macos command_exists is_number is_in_range file_exists dir_exists
export -f get_macos_version get_macos_build get_architecture is_apple_silicon get_shell
export -f confirm prompt_input
export -f to_lower to_upper trim pad_string
export -f array_contains array_join
export -f get_supermac_dir get_config_dir get_config
export -f die cleanup setup_error_handling
export -f enable_debug print_stack_trace time_command
export -f load_module get_available_modules
export -f print_supermac_header print_category_header print_category_footer

print_debug "SuperMac utilities loaded successfully"
