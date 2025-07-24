#!/bin/bash

# =============================================================================
# SuperMac - Dock Module
# =============================================================================
# Dock management and positioning commands
# 
# Commands:
#   position left/bottom/right  - Move dock position
#   autohide on/off            - Toggle dock auto-hide
#   size small/medium/large    - Set dock size
#   magnification on/off       - Toggle magnification
#   reset                      - Reset dock to defaults
#   add <app>                  - Add application to dock
#   remove <app>               - Remove application from dock
#   status                     - Show dock settings
#
# Built by CosmoLabs - https://cosmolabs.org
# =============================================================================

# =============================================================================
# Dock Position Functions
# =============================================================================

dock_get_position() {
    local position
    position=$(defaults read com.apple.dock orientation 2>/dev/null || echo "bottom")
    echo "$position"
}

dock_set_position() {
    local position="$1"
    
    case "$position" in
        "left"|"l")
            position="left"
            ;;
        "bottom"|"b")
            position="bottom"
            ;;
        "right"|"r")
            position="right"
            ;;
        *)
            print_error "Invalid position: $position"
            print_info "Valid positions: left, bottom, right"
            return 1
            ;;
    esac
    
    local current_position
    current_position=$(dock_get_position)
    
    if [[ "$current_position" == "$position" ]]; then
        print_info "Dock is already positioned on the $position"
        return 0
    fi
    
    print_info "Moving dock to $position..."
    defaults write com.apple.dock orientation "$position"
    
    if dock_restart; then
        print_success "Dock moved to $position!"
        
        case "$position" in
            "left"|"right")
                print_info "💡 Vertical dock gives you more horizontal screen space"
                ;;
            "bottom")
                print_info "💡 Bottom dock is the traditional macOS layout"
                ;;
        esac
    else
        print_error "Failed to restart dock"
        return 1
    fi
}

# =============================================================================
# Dock Auto-Hide Functions
# =============================================================================

dock_get_autohide() {
    local autohide
    autohide=$(defaults read com.apple.dock autohide 2>/dev/null || echo "0")
    
    if [[ "$autohide" == "1" ]]; then
        echo "enabled"
    else
        echo "disabled"
    fi
}

dock_set_autohide() {
    local action="$1"
    
    case "$action" in
        "on"|"enable"|"true"|"1")
            action="enable"
            ;;
        "off"|"disable"|"false"|"0")
            action="disable"
            ;;
        *)
            print_error "Invalid autohide action: $action"
            print_info "Valid actions: on, off"
            return 1
            ;;
    esac
    
    local current_state
    current_state=$(dock_get_autohide)
    
    if [[ "$current_state" == "enabled" && "$action" == "enable" ]]; then
        print_info "Dock auto-hide is already enabled"
        return 0
    elif [[ "$current_state" == "disabled" && "$action" == "disable" ]]; then
        print_info "Dock auto-hide is already disabled"
        return 0
    fi
    
    if [[ "$action" == "enable" ]]; then
        print_info "Enabling dock auto-hide..."
        defaults write com.apple.dock autohide -bool true
    else
        print_info "Disabling dock auto-hide..."
        defaults write com.apple.dock autohide -bool false
    fi
    
    if dock_restart; then
        print_success "Dock auto-hide $action"d!"
        
        if [[ "$action" == "enable" ]]; then
            print_info "💡 Move cursor to screen edge to show dock"
        else
            print_info "💡 Dock will always be visible"
        fi
    else
        print_error "Failed to restart dock"
        return 1
    fi
}

# =============================================================================
# Dock Size Functions
# =============================================================================

dock_get_size() {
    local tilesize
    tilesize=$(defaults read com.apple.dock tilesize 2>/dev/null || echo "64")
    
    if [[ "$tilesize" -le 40 ]]; then
        echo "small"
    elif [[ "$tilesize" -le 80 ]]; then
        echo "medium"
    else
        echo "large"
    fi
}

dock_set_size() {
    local size="$1"
    local tilesize
    
    case "$size" in
        "small"|"s")
            tilesize=32
            size="small"
            ;;
        "medium"|"m")
            tilesize=64
            size="medium"
            ;;
        "large"|"l")
            tilesize=96
            size="large"
            ;;
        *)
            print_error "Invalid size: $size"
            print_info "Valid sizes: small, medium, large"
            return 1
            ;;
    esac
    
    local current_size
    current_size=$(dock_get_size)
    
    if [[ "$current_size" == "$size" ]]; then
        print_info "Dock size is already set to $size"
        return 0
    fi
    
    print_info "Setting dock size to $size (${tilesize}px)..."
    defaults write com.apple.dock tilesize -int "$tilesize"
    
    if dock_restart; then
        print_success "Dock size set to $size!"
        
        case "$size" in
            "small")
                print_info "💡 Small dock saves screen space"
                ;;
            "large")
                print_info "💡 Large dock is easier to see and click"
                ;;
        esac
    else
        print_error "Failed to restart dock"
        return 1
    fi
}

# =============================================================================
# Dock Magnification Functions
# =============================================================================

dock_get_magnification() {
    local magnification
    magnification=$(defaults read com.apple.dock magnification 2>/dev/null || echo "0")
    
    if [[ "$magnification" == "1" ]]; then
        echo "enabled"
    else
        echo "disabled"
    fi
}

dock_set_magnification() {
    local action="$1"
    
    case "$action" in
        "on"|"enable"|"true"|"1")
            action="enable"
            ;;
        "off"|"disable"|"false"|"0")
            action="disable"
            ;;
        *)
            print_error "Invalid magnification action: $action"
            print_info "Valid actions: on, off"
            return 1
            ;;
    esac
    
    local current_state
    current_state=$(dock_get_magnification)
    
    if [[ "$current_state" == "enabled" && "$action" == "enable" ]]; then
        print_info "Dock magnification is already enabled"
        return 0
    elif [[ "$current_state" == "disabled" && "$action" == "disable" ]]; then
        print_info "Dock magnification is already disabled"
        return 0
    fi
    
    if [[ "$action" == "enable" ]]; then
        print_info "Enabling dock magnification..."
        defaults write com.apple.dock magnification -bool true
        # Set a reasonable magnification size
        defaults write com.apple.dock largesize -int 128
    else
        print_info "Disabling dock magnification..."
        defaults write com.apple.dock magnification -bool false
    fi
    
    if dock_restart; then
        print_success "Dock magnification ${action}d!"
        
        if [[ "$action" == "enable" ]]; then
            print_info "💡 Hover over dock icons to see magnification effect"
        fi
    else
        print_error "Failed to restart dock"
        return 1
    fi
}

# =============================================================================
# Dock Management Functions
# =============================================================================

dock_restart() {
    print_debug "Restarting dock..."
    if killall Dock 2>/dev/null; then
        sleep 2
        return 0
    else
        return 1
    fi
}

dock_reset() {
    print_warning "This will reset all dock settings to defaults"
    
    if ! confirm "Are you sure you want to reset the dock?" "n"; then
        print_info "Dock reset cancelled"
        return 0
    fi
    
    print_info "Resetting dock to default settings..."
    
    # Reset all dock preferences
    defaults delete com.apple.dock 2>/dev/null || true
    
    # Set some sensible defaults
    defaults write com.apple.dock orientation "bottom"
    defaults write com.apple.dock autohide -bool false
    defaults write com.apple.dock tilesize -int 64
    defaults write com.apple.dock magnification -bool false
    defaults write com.apple.dock show-recents -bool true
    
    if dock_restart; then
        print_success "Dock reset to default settings!"
        print_info "💡 Position: bottom, Size: medium, Auto-hide: off"
    else
        print_error "Failed to restart dock"
        return 1
    fi
}

dock_add_app() {
    local app_name="$1"
    
    if [[ -z "$app_name" ]]; then
        print_error "Application name required"
        print_info "Usage: mac dock add <application_name>"
        return 1
    fi
    
    # Find application path
    local app_path=""
    
    # Check common locations
    local search_paths=(
        "/Applications/$app_name.app"
        "/Applications/$app_name"
        "/System/Applications/$app_name.app"
        "/System/Applications/$app_name"
    )
    
    for path in "${search_paths[@]}"; do
        if [[ -d "$path" ]]; then
            app_path="$path"
            break
        fi
    done
    
    # If not found, search more broadly
    if [[ -z "$app_path" ]]; then
        app_path=$(find /Applications -maxdepth 2 -iname "*$app_name*.app" -type d | head -1)
    fi
    
    if [[ -z "$app_path" ]]; then
        print_error "Application not found: $app_name"
        print_info "Make sure the application is installed in /Applications"
        return 1
    fi
    
    print_info "Adding $(basename "$app_path" .app) to dock..."
    
    # Add to dock using dockutil if available, otherwise use defaults
    if command_exists dockutil; then
        if dockutil --add "$app_path" 2>/dev/null; then
            print_success "Added $(basename "$app_path" .app) to dock!"
        else
            print_error "Failed to add application to dock"
            return 1
        fi
    else
        print_warning "dockutil not found - using fallback method"
        print_info "Application path: $app_path"
        print_info "You can manually drag the app to your dock"
    fi
}

dock_remove_app() {
    local app_name="$1"
    
    if [[ -z "$app_name" ]]; then
        print_error "Application name required"
        print_info "Usage: mac dock remove <application_name>"
        return 1
    fi
    
    print_info "Removing $app_name from dock..."
    
    if command_exists dockutil; then
        if dockutil --remove "$app_name" 2>/dev/null; then
            print_success "Removed $app_name from dock!"
        else
            print_error "Failed to remove application from dock"
            print_info "Make sure the application name is correct"
            return 1
        fi
    else
        print_warning "dockutil not found - cannot remove apps programmatically"
        print_info "You can manually remove apps by dragging them out of the dock"
    fi
}

dock_status() {
    print_header "🚢 Dock Status"
    echo ""
    
    # Position
    local position
    position=$(dock_get_position)
    echo "  Position: $(print_command "$position")"
    
    # Auto-hide
    local autohide
    autohide=$(dock_get_autohide)
    echo "  Auto-hide: $(print_command "$autohide")"
    
    # Size
    local size tilesize
    size=$(dock_get_size)
    tilesize=$(defaults read com.apple.dock tilesize 2>/dev/null || echo "64")
    echo "  Size: $(print_command "$size") (${tilesize}px)"
    
    # Magnification
    local magnification
    magnification=$(dock_get_magnification)
    echo "  Magnification: $(print_command "$magnification")"
    
    # Show recent apps
    local show_recents
    show_recents=$(defaults read com.apple.dock show-recents 2>/dev/null || echo "1")
    if [[ "$show_recents" == "1" ]]; then
        echo "  Recent apps: $(print_command "shown")"
    else
        echo "  Recent apps: $(print_dim "hidden")"
    fi
    
    # Minimize effect
    local minimize_effect
    minimize_effect=$(defaults read com.apple.dock mineffect 2>/dev/null || echo "genie")
    echo "  Minimize effect: $(print_command "$minimize_effect")"
}

# =============================================================================
# Module Help System
# =============================================================================

dock_help() {
    print_category_header "dock" "🚢" 75
    
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "position <pos>" "Move dock (left/bottom/right)"
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "autohide on/off" "Toggle dock auto-hide"
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "size <size>" "Set dock size (small/medium/large)"
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "magnification on/off" "Toggle icon magnification"
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "reset" "Reset dock to defaults"
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "add <app>" "Add application to dock"
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "remove <app>" "Remove application from dock"
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "status" "Show current dock settings"
    
    print_category_footer 75
    echo ""
    
    print_header "Examples:"
    echo "  mac dock position left              # Move dock to left edge"
    echo "  mac dock autohide on                # Enable auto-hide"
    echo "  mac dock size small                 # Make dock smaller"
    echo "  mac dock add Safari                 # Add Safari to dock"
    echo "  mac dock reset                      # Reset to defaults"
    echo ""
    
    print_header "Position Options:"
    echo "  left     - Dock on left edge (saves horizontal space)"
    echo "  bottom   - Traditional bottom position"
    echo "  right    - Dock on right edge"
    echo ""
    
    print_header "Size Options:"
    echo "  small    - Compact dock (32px icons)"
    echo "  medium   - Standard size (64px icons)"
    echo "  large    - Large icons (96px icons)"
    echo ""
    
    print_header "Tips:"
    echo "  • Auto-hide gives you more screen real estate"
    echo "  • Left/right positions work well on widescreen displays"
    echo "  • Magnification helps with small dock sizes"
    echo "  • Reset dock if settings become corrupted"
    echo ""
}

# Search function for this module
dock_search() {
    local search_term="$1"
    local results=""
    
    if [[ "position" == *"$search_term"* ]] || [[ "move" == *"$search_term"* ]]; then
        results+="  mac dock position <pos>          Move dock position\n"
    fi
    
    if [[ "hide" == *"$search_term"* ]] || [[ "auto" == *"$search_term"* ]]; then
        results+="  mac dock autohide on/off         Toggle auto-hide\n"
    fi
    
    if [[ "size" == *"$search_term"* ]] || [[ "small" == *"$search_term"* ]] || [[ "large" == *"$search_term"* ]]; then
        results+="  mac dock size <size>             Set dock size\n"
    fi
    
    if [[ "magnif" == *"$search_term"* ]] || [[ "zoom" == *"$search_term"* ]]; then
        results+="  mac dock magnification on/off    Toggle magnification\n"
    fi
    
    if [[ "add" == *"$search_term"* ]] || [[ "remove" == *"$search_term"* ]] || [[ "app" == *"$search_term"* ]]; then
        results+="  mac dock add <app>               Add application\n"
        results+="  mac dock remove <app>            Remove application\n"
    fi
    
    if [[ "reset" == *"$search_term"* ]] || [[ "default" == *"$search_term"* ]]; then
        results+="  mac dock reset                   Reset to defaults\n"
    fi
    
    if [[ "status" == *"$search_term"* ]] || [[ "info" == *"$search_term"* ]]; then
        results+="  mac dock status                  Show dock settings\n"
    fi
    
    if [[ -n "$results" ]]; then
        echo -e "$results"
    fi
}

# =============================================================================
# Module Dispatcher
# =============================================================================

dock_dispatch() {
    local action="$1"
    shift  # Remove action from arguments
    
    case "$action" in
        "position")
            dock_set_position "$1"
            ;;
        "autohide")
            dock_set_autohide "$1"
            ;;
        "size")
            dock_set_size "$1"
            ;;
        "magnification")
            dock_set_magnification "$1"
            ;;
        "reset")
            dock_reset
            ;;
        "add")
            dock_add_app "$1"
            ;;
        "remove")
            dock_remove_app "$1"
            ;;
        "status")
            dock_status
            ;;
        "help"|"-h"|"--help")
            dock_help
            ;;
        *)
            print_error "Unknown dock action: $action"
            echo ""
            print_info "Available actions: position, autohide, size, magnification, reset, add, remove, status"
            print_info "Use 'mac help dock' for detailed information"
            return 1
            ;;
    esac
}

print_debug "Dock module loaded successfully"
