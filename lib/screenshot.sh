#!/bin/bash

# =============================================================================
# SuperMac - Screenshot Module
# =============================================================================
# Screenshot settings and management commands
# 
# Commands:
#   location desktop/downloads/clipboard/<path>  - Set save location
#   format png/jpg/pdf/tiff                     - Set file format
#   shadows on/off                              - Toggle window shadows
#   show-cursor on/off                          - Toggle cursor in screenshots
#   thumbnail on/off                            - Toggle thumbnail preview
#   sound on/off                                - Toggle camera sound
#   name-format <format>                        - Set filename format
#   status                                      - Show current settings
#   reset                                       - Reset to defaults
#   take [area/window/screen]                   - Take screenshot now
#
# Built by CosmoLabs - https://cosmolabs.org
# =============================================================================

# =============================================================================
# Screenshot Location Functions
# =============================================================================

screenshot_get_location() {
    defaults read com.apple.screencapture location 2>/dev/null || echo "$HOME/Desktop"
}

screenshot_set_location() {
    local location="$1"
    
    if [[ -z "$location" ]]; then
        print_error "Location required"
        print_info "Usage: mac screenshot location <desktop|downloads|clipboard|path>"
        return 1
    fi
    
    local target_path
    
    case "$location" in
        "desktop")
            target_path="$HOME/Desktop"
            print_info "Setting screenshot location to Desktop..."
            ;;
        "downloads")
            target_path="$HOME/Downloads"
            print_info "Setting screenshot location to Downloads..."
            ;;
        "clipboard")
            # Special case - saves to clipboard only
            print_info "Setting screenshots to save to clipboard only..."
            defaults write com.apple.screencapture location -string "clipboard"
            screenshot_restart_service
            print_success "Screenshots will now be saved to clipboard only!"
            print_info "💡 Use Cmd+Shift+4 to capture and copy directly"
            return 0
            ;;
        "documents")
            target_path="$HOME/Documents"
            print_info "Setting screenshot location to Documents..."
            ;;
        "pictures")
            target_path="$HOME/Pictures"
            print_info "Setting screenshot location to Pictures..."
            ;;
        *)
            # Custom path
            if [[ -d "$location" ]]; then
                target_path="$location"
                print_info "Setting screenshot location to: $location"
            else
                print_error "Directory does not exist: $location"
                print_info "Create the directory first or use: desktop, downloads, clipboard"
                return 1
            fi
            ;;
    esac
    
    # Set the location (except for clipboard which is handled above)
    if [[ "$location" != "clipboard" ]]; then
        defaults write com.apple.screencapture location -string "$target_path"
        screenshot_restart_service
        print_success "Screenshot location set to: $target_path"
        print_info "💡 Screenshots will now be saved here by default"
    fi
}

# =============================================================================
# Screenshot Format Functions
# =============================================================================

screenshot_get_format() {
    defaults read com.apple.screencapture type 2>/dev/null || echo "png"
}

screenshot_set_format() {
    local format="$1"
    
    if [[ -z "$format" ]]; then
        print_error "Format required"
        print_info "Usage: mac screenshot format <png|jpg|pdf|tiff>"
        return 1
    fi
    
    case "$format" in
        "png")
            print_info "Setting screenshot format to PNG..."
            defaults write com.apple.screencapture type -string "png"
            print_success "Format set to PNG (lossless, best quality)"
            ;;
        "jpg"|"jpeg")
            print_info "Setting screenshot format to JPEG..."
            defaults write com.apple.screencapture type -string "jpg"
            print_success "Format set to JPEG (smaller files, some compression)"
            ;;
        "pdf")
            print_info "Setting screenshot format to PDF..."
            defaults write com.apple.screencapture type -string "pdf"
            print_success "Format set to PDF (vector format, scalable)"
            ;;
        "tiff"|"tif")
            print_info "Setting screenshot format to TIFF..."
            defaults write com.apple.screencapture type -string "tiff"
            print_success "Format set to TIFF (lossless, larger files)"
            ;;
        *)
            print_error "Invalid format: $format"
            print_info "Available formats: png, jpg, pdf, tiff"
            return 1
            ;;
    esac
    
    screenshot_restart_service
}

# =============================================================================
# Screenshot Behavior Settings
# =============================================================================

screenshot_get_shadows() {
    local disable_shadows
    disable_shadows=$(defaults read com.apple.screencapture disable-shadow 2>/dev/null || echo "false")
    
    if [[ "$disable_shadows" == "true" ]] || [[ "$disable_shadows" == "1" ]]; then
        echo "disabled"
    else
        echo "enabled"
    fi
}

screenshot_set_shadows() {
    local action="$1"
    
    case "$action" in
        "on"|"enable"|"true")
            print_info "Enabling window shadows in screenshots..."
            defaults write com.apple.screencapture disable-shadow -bool false
            screenshot_restart_service
            print_success "Window shadows enabled!"
            print_info "💡 Windows will include drop shadows in screenshots"
            ;;
        "off"|"disable"|"false")
            print_info "Disabling window shadows in screenshots..."
            defaults write com.apple.screencapture disable-shadow -bool true
            screenshot_restart_service
            print_success "Window shadows disabled!"
            print_info "💡 Windows will appear without shadows for cleaner look"
            ;;
        *)
            print_error "Invalid shadows action: $action"
            print_info "Usage: mac screenshot shadows [on|off]"
            return 1
            ;;
    esac
}

screenshot_get_cursor() {
    local show_cursor
    show_cursor=$(defaults read com.apple.screencapture showsCursor 2>/dev/null || echo "false")
    
    if [[ "$show_cursor" == "true" ]] || [[ "$show_cursor" == "1" ]]; then
        echo "enabled"
    else
        echo "disabled"
    fi
}

screenshot_set_cursor() {
    local action="$1"
    
    case "$action" in
        "on"|"enable"|"show"|"true")
            print_info "Enabling cursor in screenshots..."
            defaults write com.apple.screencapture showsCursor -bool true
            screenshot_restart_service
            print_success "Cursor will now appear in screenshots!"
            print_info "💡 Useful for tutorials and demonstrations"
            ;;
        "off"|"disable"|"hide"|"false")
            print_info "Disabling cursor in screenshots..."
            defaults write com.apple.screencapture showsCursor -bool false
            screenshot_restart_service
            print_success "Cursor will not appear in screenshots!"
            print_info "💡 Creates cleaner screenshots"
            ;;
        *)
            print_error "Invalid cursor action: $action"
            print_info "Usage: mac screenshot show-cursor [on|off]"
            return 1
            ;;
    esac
}

screenshot_get_thumbnail() {
    local show_thumbnail
    show_thumbnail=$(defaults read com.apple.screencapture show-thumbnail 2>/dev/null || echo "true")
    
    if [[ "$show_thumbnail" == "true" ]] || [[ "$show_thumbnail" == "1" ]]; then
        echo "enabled"
    else
        echo "disabled"
    fi
}

screenshot_set_thumbnail() {
    local action="$1"
    
    case "$action" in
        "on"|"enable"|"show"|"true")
            print_info "Enabling screenshot thumbnail preview..."
            defaults write com.apple.screencapture show-thumbnail -bool true
            screenshot_restart_service
            print_success "Thumbnail preview enabled!"
            print_info "💡 You'll see a preview thumbnail after taking screenshots"
            ;;
        "off"|"disable"|"hide"|"false")
            print_info "Disabling screenshot thumbnail preview..."
            defaults write com.apple.screencapture show-thumbnail -bool false
            screenshot_restart_service
            print_success "Thumbnail preview disabled!"
            print_info "💡 Screenshots will save immediately without preview"
            ;;
        *)
            print_error "Invalid thumbnail action: $action"
            print_info "Usage: mac screenshot thumbnail [on|off]"
            return 1
            ;;
    esac
}

screenshot_set_sound() {
    local action="$1"
    
    case "$action" in
        "on"|"enable"|"true")
            print_info "Enabling screenshot sound..."
            defaults write com.apple.screencapture disable-sound -bool false
            screenshot_restart_service
            print_success "Screenshot sound enabled!"
            print_info "🔊 You'll hear a camera sound when taking screenshots"
            ;;
        "off"|"disable"|"false")
            print_info "Disabling screenshot sound..."
            defaults write com.apple.screencapture disable-sound -bool true
            screenshot_restart_service
            print_success "Screenshot sound disabled!"
            print_info "🔇 Screenshots will be taken silently"
            ;;
        *)
            print_error "Invalid sound action: $action"
            print_info "Usage: mac screenshot sound [on|off]"
            return 1
            ;;
    esac
}

screenshot_get_sound() {
    local disable_sound
    disable_sound=$(defaults read com.apple.screencapture disable-sound 2>/dev/null || echo "false")
    
    if [[ "$disable_sound" == "true" ]] || [[ "$disable_sound" == "1" ]]; then
        echo "disabled"
    else
        echo "enabled"
    fi
}

# =============================================================================
# Screenshot Filename Functions
# =============================================================================

screenshot_set_name_format() {
    local format="$1"
    
    if [[ -z "$format" ]]; then
        print_error "Name format required"
        print_info "Usage: mac screenshot name-format <format>"
        print_info "Example: 'Screenshot %Y-%m-%d at %H.%M.%S'"
        return 1
    fi
    
    print_info "Setting screenshot name format..."
    defaults write com.apple.screencapture name -string "$format"
    screenshot_restart_service
    print_success "Screenshot name format updated!"
    print_info "💡 New screenshots will use this naming pattern"
}

screenshot_get_name_format() {
    defaults read com.apple.screencapture name 2>/dev/null || echo "Screenshot %Y-%m-%d at %H.%M.%S"
}

# =============================================================================
# Screenshot Management Functions
# =============================================================================

screenshot_restart_service() {
    # Restart the screenshot service to apply changes
    killall SystemUIServer 2>/dev/null || true
    sleep 1
}

screenshot_reset() {
    print_warning "This will reset all screenshot settings to defaults"
    
    if ! confirm "Are you sure you want to reset screenshot settings?" "n"; then
        print_info "Screenshot reset cancelled"
        return 0
    fi
    
    print_info "Resetting screenshot settings to defaults..."
    
    # Remove all screenshot-related preferences
    defaults delete com.apple.screencapture 2>/dev/null || true
    
    screenshot_restart_service
    
    print_success "Screenshot settings reset to defaults!"
    print_info "💡 Screenshots will now save to Desktop in PNG format"
}

screenshot_status() {
    print_header "📸 Screenshot Settings"
    echo ""
    
    # Location
    local location
    location=$(screenshot_get_location)
    if [[ "$location" == "clipboard" ]]; then
        echo "  Location: $(print_command "Clipboard only")"
    else
        echo "  Location: $(print_command "$location")"
    fi
    
    # Format
    local format
    format=$(screenshot_get_format)
    echo "  Format: $(print_command "$(to_upper "$format")")"
    
    # Shadows
    local shadows
    shadows=$(screenshot_get_shadows)
    echo "  Window shadows: $(print_command "$shadows")"
    
    # Cursor
    local cursor
    cursor=$(screenshot_get_cursor)
    echo "  Show cursor: $(print_command "$cursor")"
    
    # Thumbnail
    local thumbnail
    thumbnail=$(screenshot_get_thumbnail)
    echo "  Thumbnail preview: $(print_command "$thumbnail")"
    
    # Sound
    local sound
    sound=$(screenshot_get_sound)
    echo "  Camera sound: $(print_command "$sound")"
    
    # Name format
    local name_format
    name_format=$(screenshot_get_name_format)
    echo "  Name format: $(print_dim "$name_format")"
    
    echo ""
    print_subheader "Quick Actions:"
    echo "  Cmd+Shift+3         Full screen screenshot"
    echo "  Cmd+Shift+4         Select area screenshot"
    echo "  Cmd+Shift+5         Screenshot options menu"
    echo ""
    echo "  mac screenshot take         Take screenshot now"
}

# =============================================================================
# Take Screenshots
# =============================================================================

screenshot_take() {
    local type="${1:-area}"
    
    case "$type" in
        "screen"|"fullscreen")
            print_info "Taking full screen screenshot..."
            screencapture -x "/tmp/screenshot-$(date +%s).png"
            print_success "Full screen screenshot taken!"
            ;;
        "area"|"selection")
            print_info "Select area for screenshot (press Space for window mode)..."
            screencapture -i "/tmp/screenshot-$(date +%s).png"
            print_success "Area screenshot taken!"
            ;;
        "window")
            print_info "Click on a window to capture..."
            screencapture -i -w "/tmp/screenshot-$(date +%s).png"
            print_success "Window screenshot taken!"
            ;;
        *)
            print_error "Invalid screenshot type: $type"
            print_info "Available types: screen, area, window"
            return 1
            ;;
    esac
    
    local location
    location=$(screenshot_get_location)
    if [[ "$location" != "clipboard" ]]; then
        print_info "Screenshot saved to: $location"
    else
        print_info "Screenshot copied to clipboard"
    fi
}

# =============================================================================
# Module Help System
# =============================================================================

screenshot_help() {
    print_category_header "screenshot" "📸" 80
    
    printf "${PURPLE}│${NC}  %-30s %-40s ${PURPLE}│${NC}\n" "location desktop/downloads/path" "Set save location"
    printf "${PURPLE}│${NC}  %-30s %-40s ${PURPLE}│${NC}\n" "location clipboard" "Save to clipboard only"
    printf "${PURPLE}│${NC}  %-30s %-40s ${PURPLE}│${NC}\n" "format png/jpg/pdf/tiff" "Set file format"
    printf "${PURPLE}│${NC}  %-30s %-40s ${PURPLE}│${NC}\n" "shadows on/off" "Toggle window shadows"
    printf "${PURPLE}│${NC}  %-30s %-40s ${PURPLE}│${NC}\n" "show-cursor on/off" "Toggle cursor visibility"
    printf "${PURPLE}│${NC}  %-30s %-40s ${PURPLE}│${NC}\n" "thumbnail on/off" "Toggle preview thumbnail"
    printf "${PURPLE}│${NC}  %-30s %-40s ${PURPLE}│${NC}\n" "sound on/off" "Toggle camera sound"
    printf "${PURPLE}│${NC}  %-30s %-40s ${PURPLE}│${NC}\n" "name-format <format>" "Set filename pattern"
    printf "${PURPLE}│${NC}  %-30s %-40s ${PURPLE}│${NC}\n" "take [area/window/screen]" "Take screenshot now"
    printf "${PURPLE}│${NC}  %-30s %-40s ${PURPLE}│${NC}\n" "status" "Show current settings"
    printf "${PURPLE}│${NC}  %-30s %-40s ${PURPLE}│${NC}\n" "reset" "Reset to default settings"
    
    print_category_footer 80
    echo ""
    
    print_header "Examples:"
    echo "  mac screenshot location downloads   # Save to Downloads folder"
    echo "  mac screenshot format jpg           # Use JPEG format"
    echo "  mac screenshot shadows off          # Remove window shadows"
    echo "  mac screenshot location clipboard   # Copy to clipboard only"
    echo "  mac screenshot take window          # Capture specific window"
    echo ""
    
    print_header "Location Options:"
    echo "  • desktop     - ~/Desktop (default)"
    echo "  • downloads   - ~/Downloads folder"
    echo "  • clipboard   - Copy to clipboard only"
    echo "  • documents   - ~/Documents folder"
    echo "  • pictures    - ~/Pictures folder"
    echo "  • /custom/path - Any custom directory"
    echo ""
    
    print_header "Format Guidelines:"
    echo "  • PNG: Best quality, larger files (default)"
    echo "  • JPG: Smaller files, slight compression"
    echo "  • PDF: Vector format, good for text"
    echo "  • TIFF: Lossless, very large files"
    echo ""
    
    print_header "Keyboard Shortcuts:"
    echo "  • Cmd+Shift+3: Full screen"
    echo "  • Cmd+Shift+4: Select area (Space for window)"
    echo "  • Cmd+Shift+5: Screenshot utility with options"
    echo "  • Add Ctrl to copy to clipboard instead"
    echo ""
    
    print_header "Pro Tips:"
    echo "  • Use clipboard location for quick sharing"
    echo "  • Disable shadows for cleaner UI screenshots"
    echo "  • Enable cursor for tutorial screenshots"
    echo "  • JPG format good for photos, PNG for UI"
    echo ""
}

# Search function for this module
screenshot_search() {
    local search_term="$1"
    local results=""
    
    if [[ "location" == *"$search_term"* ]] || [[ "save" == *"$search_term"* ]] || [[ "folder" == *"$search_term"* ]]; then
        results+="  mac screenshot location <path>    Set save location\n"
    fi
    
    if [[ "format" == *"$search_term"* ]] || [[ "png" == *"$search_term"* ]] || [[ "jpg" == *"$search_term"* ]]; then
        results+="  mac screenshot format <type>      Set file format\n"
    fi
    
    if [[ "shadow" == *"$search_term"* ]] || [[ "window" == *"$search_term"* ]]; then
        results+="  mac screenshot shadows on/off     Toggle window shadows\n"
    fi
    
    if [[ "cursor" == *"$search_term"* ]] || [[ "mouse" == *"$search_term"* ]]; then
        results+="  mac screenshot show-cursor on/off Toggle cursor\n"
    fi
    
    if [[ "thumbnail" == *"$search_term"* ]] || [[ "preview" == *"$search_term"* ]]; then
        results+="  mac screenshot thumbnail on/off   Toggle preview\n"
    fi
    
    if [[ "sound" == *"$search_term"* ]] || [[ "camera" == *"$search_term"* ]] || [[ "audio" == *"$search_term"* ]]; then
        results+="  mac screenshot sound on/off       Toggle camera sound\n"
    fi
    
    if [[ "take" == *"$search_term"* ]] || [[ "capture" == *"$search_term"* ]]; then
        results+="  mac screenshot take [type]        Take screenshot\n"
    fi
    
    if [[ "clipboard" == *"$search_term"* ]] || [[ "copy" == *"$search_term"* ]]; then
        results+="  mac screenshot location clipboard  Copy only mode\n"
    fi
    
    if [[ -n "$results" ]]; then
        echo -e "$results"
    fi
}

# =============================================================================
# Module Dispatcher
# =============================================================================

screenshot_dispatch() {
    local action="$1"
    shift  # Remove action from arguments
    
    case "$action" in
        "location")
            screenshot_set_location "$1"
            ;;
        "format")
            screenshot_set_format "$1"
            ;;
        "shadows")
            screenshot_set_shadows "$1"
            ;;
        "show-cursor")
            screenshot_set_cursor "$1"
            ;;
        "thumbnail")
            screenshot_set_thumbnail "$1"
            ;;
        "sound")
            screenshot_set_sound "$1"
            ;;
        "name-format")
            screenshot_set_name_format "$1"
            ;;
        "take")
            screenshot_take "$1"
            ;;
        "reset")
            screenshot_reset
            ;;
        "status")
            screenshot_status
            ;;
        "help"|"-h"|"--help")
            screenshot_help
            ;;
        *)
            print_error "Unknown screenshot action: $action"
            echo ""
            print_info "Available actions: location, format, shadows, show-cursor, thumbnail, sound, name-format, take, reset, status"
            print_info "Use 'mac help screenshot' for detailed information"
            return 1
            ;;
    esac
}

print_debug "Screenshot module loaded successfully"
