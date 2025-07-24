#!/bin/bash

# =============================================================================
# SuperMac - Display Module
# =============================================================================
# Display and appearance settings management
# 
# Commands:
#   brightness <0-100>     - Set screen brightness percentage
#   dark-mode             - Switch to dark mode
#   light-mode            - Switch to light mode
#   toggle-mode           - Toggle between dark and light mode
#   night-shift on/off    - Control night shift
#   true-tone on/off      - Control True Tone (if supported)
#   detect                - Force display detection
#   resolution list       - List available resolutions
#
# Built by CosmoLabs - https://cosmolabs.org
# =============================================================================

# =============================================================================
# Display Brightness Functions
# =============================================================================

display_get_brightness() {
    # Get current brightness (returns value between 0.0 and 1.0)
    local brightness
    brightness=$(osascript -e 'tell application "System Events" to get brightness of first item of (displays whose built in is true)' 2>/dev/null)
    
    if [[ -n "$brightness" ]]; then
        # Convert to percentage (0-100)
        local percentage
        percentage=$(echo "$brightness * 100" | bc 2>/dev/null | cut -d. -f1)
        echo "${percentage:-0}"
    else
        echo "0"
    fi
}

display_set_brightness() {
    local target="$1"
    
    if [[ -z "$target" ]]; then
        print_error "Brightness level required"
        print_info "Usage: mac display brightness <0-100>"
        return 1
    fi
    
    if ! is_in_range "$target" 0 100; then
        print_error "Brightness must be between 0 and 100"
        return 1
    fi
    
    local current
    current=$(display_get_brightness)
    
    print_info "Setting brightness from $current% to $target%..."
    
    # Convert percentage to decimal (0.0-1.0)
    local decimal
    decimal=$(echo "scale=2; $target / 100" | bc 2>/dev/null)
    
    if osascript -e "tell application \"System Events\" to set brightness of first item of (displays whose built in is true) to $decimal" 2>/dev/null; then
        print_success "Brightness set to $target%"
        
        # Provide contextual feedback
        if [[ $target -le 20 ]]; then
            print_info "💡 Low brightness - good for nighttime use"
        elif [[ $target -ge 80 ]]; then
            print_info "💡 High brightness - good for bright environments"
        fi
    else
        print_error "Failed to set brightness"
        print_info "Make sure you have permission to control display settings"
        return 1
    fi
}

# =============================================================================
# Dark Mode Functions
# =============================================================================

display_get_appearance() {
    local appearance
    appearance=$(defaults read -g AppleInterfaceStyle 2>/dev/null || echo "Light")
    echo "$appearance"
}

display_set_dark_mode() {
    local current
    current=$(display_get_appearance)
    
    if [[ "$current" == "Dark" ]]; then
        print_info "Already in dark mode"
        return 0
    fi
    
    print_info "Switching to dark mode..."
    
    if osascript -e 'tell application "System Events" to tell appearance preferences to set dark mode to true' 2>/dev/null; then
        print_success "Switched to dark mode!"
        print_info "🌙 Dark mode is easier on the eyes, especially at night"
    else
        # Fallback method
        defaults write NSGlobalDomain AppleInterfaceStyle -string "Dark"
        print_success "Dark mode enabled (may require app restart)"
    fi
}

display_set_light_mode() {
    local current
    current=$(display_get_appearance)
    
    if [[ "$current" == "Light" ]]; then
        print_info "Already in light mode"
        return 0
    fi
    
    print_info "Switching to light mode..."
    
    if osascript -e 'tell application "System Events" to tell appearance preferences to set dark mode to false' 2>/dev/null; then
        print_success "Switched to light mode!"
        print_info "☀️ Light mode provides better contrast in bright environments"
    else
        # Fallback method
        defaults delete NSGlobalDomain AppleInterfaceStyle 2>/dev/null
        print_success "Light mode enabled (may require app restart)"
    fi
}

display_toggle_mode() {
    local current
    current=$(display_get_appearance)
    
    print_info "Current mode: $current"
    
    if [[ "$current" == "Dark" ]]; then
        display_set_light_mode
    else
        display_set_dark_mode
    fi
}

# =============================================================================
# Night Shift Functions
# =============================================================================

display_get_night_shift_status() {
    # Check if Night Shift is enabled
    local enabled
    enabled=$(defaults read com.apple.CoreBrightness "CBUser-$(id -u)" 2>/dev/null | grep -E "BlueLightReductionEnabled.*=.*1" >/dev/null && echo "enabled" || echo "disabled")
    echo "$enabled"
}

display_night_shift() {
    local action="$1"
    
    case "$action" in
        "on"|"enable")
            print_info "Enabling Night Shift..."
            if osascript -e 'tell application "System Events" to tell appearance preferences to set night shift enabled to true' 2>/dev/null; then
                print_success "Night Shift enabled!"
                print_info "🌅 Night Shift reduces blue light for better sleep"
            else
                print_error "Failed to enable Night Shift"
                print_info "You may need to enable it manually in System Preferences"
                return 1
            fi
            ;;
        "off"|"disable")
            print_info "Disabling Night Shift..."
            if osascript -e 'tell application "System Events" to tell appearance preferences to set night shift enabled to false' 2>/dev/null; then
                print_success "Night Shift disabled!"
            else
                print_error "Failed to disable Night Shift"
                return 1
            fi
            ;;
        "status")
            local status
            status=$(display_get_night_shift_status)
            echo "Night Shift: $(print_command "$status")"
            ;;
        *)
            print_error "Invalid Night Shift action: $action"
            print_info "Usage: mac display night-shift [on|off|status]"
            return 1
            ;;
    esac
}

# =============================================================================
# True Tone Functions
# =============================================================================

display_true_tone() {
    local action="$1"
    
    # Check if True Tone is supported
    if ! system_profiler SPDisplaysDataType 2>/dev/null | grep -q "True Tone"; then
        print_warning "True Tone not supported on this display"
        return 1
    fi
    
    case "$action" in
        "on"|"enable")
            print_info "Enabling True Tone..."
            if osascript -e 'tell application "System Events" to tell appearance preferences to set true tone enabled to true' 2>/dev/null; then
                print_success "True Tone enabled!"
                print_info "🎨 True Tone adjusts colors based on ambient lighting"
            else
                print_error "Failed to enable True Tone"
                return 1
            fi
            ;;
        "off"|"disable")
            print_info "Disabling True Tone..."
            if osascript -e 'tell application "System Events" to tell appearance preferences to set true tone enabled to false' 2>/dev/null; then
                print_success "True Tone disabled!"
            else
                print_error "Failed to disable True Tone"
                return 1
            fi
            ;;
        *)
            print_error "Invalid True Tone action: $action"
            print_info "Usage: mac display true-tone [on|off]"
            return 1
            ;;
    esac
}

# =============================================================================
# Display Detection & Resolution
# =============================================================================

display_detect() {
    print_info "Forcing display detection..."
    
    # Use system_profiler to trigger display detection
    system_profiler SPDisplaysDataType >/dev/null 2>&1
    
    # Also try the detect displays button equivalent
    if osascript -e 'tell application "System Events" to keystroke "d" using {command down}' 2>/dev/null; then
        print_success "Display detection triggered!"
        print_info "💡 This may help if external displays aren't detected properly"
    else
        print_warning "Display detection may not have worked"
        print_info "Try manually: System Preferences → Displays → Detect Displays"
    fi
}

display_list_resolutions() {
    print_header "🖥️ Available Display Resolutions"
    echo ""
    
    # Get display information
    system_profiler SPDisplaysDataType 2>/dev/null | grep -E "(Display Type|Resolution)" | while read -r line; do
        if [[ "$line" == *"Display Type"* ]]; then
            display_name=$(echo "$line" | cut -d: -f2 | trim)
            echo -e "${BOLD}$display_name:${NC}"
        elif [[ "$line" == *"Resolution"* ]]; then
            resolution=$(echo "$line" | cut -d: -f2 | trim)
            echo "  Current: $(print_command "$resolution")"
        fi
    done
    
    echo ""
    print_info "💡 Change resolution in System Preferences → Displays"
}

display_status() {
    print_header "🖥️ Display Status"
    echo ""
    
    # Brightness
    local brightness
    brightness=$(display_get_brightness)
    echo "  Brightness: $(print_command "$brightness%")"
    
    # Appearance mode
    local appearance
    appearance=$(display_get_appearance)
    echo "  Appearance: $(print_command "$appearance Mode")"
    
    # Night Shift
    local night_shift
    night_shift=$(display_get_night_shift_status)
    echo "  Night Shift: $(print_command "$night_shift")"
    
    # Display count
    local display_count
    display_count=$(system_profiler SPDisplaysDataType 2>/dev/null | grep -c "Display Type" || echo "1")
    echo "  Connected Displays: $(print_command "$display_count")"
    
    # Current resolution
    local resolution
    resolution=$(system_profiler SPDisplaysDataType 2>/dev/null | grep "Resolution:" | head -1 | cut -d: -f2 | trim || echo "Unknown")
    echo "  Primary Resolution: $(print_command "$resolution")"
}

# =============================================================================
# Module Help System
# =============================================================================

display_help() {
    print_category_header "display" "🖥️" 70
    
    printf "${PURPLE}│${NC}  %-25s %-35s ${PURPLE}│${NC}\n" "brightness <0-100>" "Set screen brightness percentage"
    printf "${PURPLE}│${NC}  %-25s %-35s ${PURPLE}│${NC}\n" "dark-mode" "Switch to dark appearance"
    printf "${PURPLE}│${NC}  %-25s %-35s ${PURPLE}│${NC}\n" "light-mode" "Switch to light appearance"
    printf "${PURPLE}│${NC}  %-25s %-35s ${PURPLE}│${NC}\n" "toggle-mode" "Toggle dark/light mode"
    printf "${PURPLE}│${NC}  %-25s %-35s ${PURPLE}│${NC}\n" "night-shift on/off" "Control Night Shift"
    printf "${PURPLE}│${NC}  %-25s %-35s ${PURPLE}│${NC}\n" "true-tone on/off" "Control True Tone"
    printf "${PURPLE}│${NC}  %-25s %-35s ${PURPLE}│${NC}\n" "detect" "Force display detection"
    printf "${PURPLE}│${NC}  %-25s %-35s ${PURPLE}│${NC}\n" "resolution list" "List available resolutions"
    printf "${PURPLE}│${NC}  %-25s %-35s ${PURPLE}│${NC}\n" "status" "Show display settings"
    
    print_category_footer 70
    echo ""
    
    print_header "Examples:"
    echo "  mac display brightness 75           # Set brightness to 75%"
    echo "  mac display dark-mode               # Switch to dark mode"
    echo "  mac display night-shift on          # Enable Night Shift"
    echo "  mac display detect                  # Detect external displays"
    echo "  mac display status                  # Show current settings"
    echo ""
    
    print_header "Global Shortcuts:"
    echo "  mac dark                            # Quick dark mode"
    echo "  mac light                           # Quick light mode"
    echo ""
    
    print_header "Tips:"
    echo "  • Use low brightness (20-40%) for nighttime work"
    echo "  • Night Shift reduces blue light for better sleep"
    echo "  • True Tone adjusts colors based on ambient lighting"
    echo "  • Use detect if external displays aren't recognized"
    echo ""
}

# Search function for this module
display_search() {
    local search_term="$1"
    local results=""
    
    if [[ "brightness" == *"$search_term"* ]] || [[ "bright" == *"$search_term"* ]]; then
        results+="  mac display brightness <0-100>     Set screen brightness\n"
    fi
    
    if [[ "dark" == *"$search_term"* ]] || [[ "mode" == *"$search_term"* ]] || [[ "theme" == *"$search_term"* ]]; then
        results+="  mac display dark-mode               Switch to dark mode\n"
        results+="  mac display light-mode              Switch to light mode\n"
        results+="  mac display toggle-mode             Toggle appearance\n"
    fi
    
    if [[ "night" == *"$search_term"* ]] || [[ "shift" == *"$search_term"* ]] || [[ "blue" == *"$search_term"* ]]; then
        results+="  mac display night-shift on/off      Control Night Shift\n"
    fi
    
    if [[ "true" == *"$search_term"* ]] || [[ "tone" == *"$search_term"* ]]; then
        results+="  mac display true-tone on/off        Control True Tone\n"
    fi
    
    if [[ "detect" == *"$search_term"* ]] || [[ "resolution" == *"$search_term"* ]]; then
        results+="  mac display detect                  Detect displays\n"
        results+="  mac display resolution list         List resolutions\n"
    fi
    
    if [[ -n "$results" ]]; then
        echo -e "$results"
    fi
}

# =============================================================================
# Module Dispatcher
# =============================================================================

display_dispatch() {
    local action="$1"
    shift  # Remove action from arguments
    
    case "$action" in
        "brightness")
            display_set_brightness "$1"
            ;;
        "dark-mode")
            display_set_dark_mode
            ;;
        "light-mode")
            display_set_light_mode
            ;;
        "toggle-mode")
            display_toggle_mode
            ;;
        "night-shift")
            display_night_shift "$1"
            ;;
        "true-tone")
            display_true_tone "$1"
            ;;
        "detect")
            display_detect
            ;;
        "resolution")
            case "$1" in
                "list") display_list_resolutions ;;
                *) print_error "Usage: mac display resolution list" ;;
            esac
            ;;
        "status")
            display_status
            ;;
        "help"|"-h"|"--help")
            display_help
            ;;
        *)
            print_error "Unknown display action: $action"
            echo ""
            print_info "Available actions: brightness, dark-mode, light-mode, toggle-mode, night-shift, true-tone, detect, resolution, status"
            print_info "Use 'mac help display' for detailed information"
            return 1
            ;;
    esac
}

print_debug "Display module loaded successfully"
