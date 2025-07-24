#!/bin/bash

# =============================================================================
# SuperMac - Audio Module
# =============================================================================
# Audio control and device management commands
# 
# Commands:
#   volume <0-100>         - Set system volume percentage
#   mute                   - Mute system audio
#   unmute                 - Unmute system audio
#   toggle-mute            - Toggle mute state
#   devices                - List audio devices
#   input <device>         - Set audio input device
#   output <device>        - Set audio output device
#   balance left/right/center - Set audio balance
#   status                 - Show audio status
#
# Built by CosmoLabs - https://cosmolabs.org
# =============================================================================

# =============================================================================
# Volume Control Functions
# =============================================================================

audio_get_volume() {
    # Get current volume (0-100)
    osascript -e "output volume of (get volume settings)" 2>/dev/null || echo "0"
}

audio_get_mute_state() {
    # Check if audio is muted
    local muted
    muted=$(osascript -e "output muted of (get volume settings)" 2>/dev/null || echo "false")
    
    if [[ "$muted" == "true" ]]; then
        echo "muted"
    else
        echo "unmuted"
    fi
}

audio_set_volume() {
    local target="$1"
    
    if [[ -z "$target" ]]; then
        print_error "Volume level required"
        print_info "Usage: mac audio volume <0-100>"
        return 1
    fi
    
    if ! is_in_range "$target" 0 100; then
        print_error "Volume must be between 0 and 100"
        return 1
    fi
    
    local current
    current=$(audio_get_volume)
    
    print_info "Setting volume from $current% to $target%..."
    
    if osascript -e "set volume output volume $target" 2>/dev/null; then
        print_success "Volume set to $target%"
        
        # Provide contextual feedback
        if [[ $target -eq 0 ]]; then
            print_info "🔇 Audio is now silent"
        elif [[ $target -le 25 ]]; then
            print_info "🔈 Low volume"
        elif [[ $target -le 75 ]]; then
            print_info "🔉 Medium volume"
        else
            print_info "🔊 High volume"
        fi
    else
        print_error "Failed to set volume"
        return 1
    fi
}

audio_mute() {
    local current_state
    current_state=$(audio_get_mute_state)
    
    if [[ "$current_state" == "muted" ]]; then
        print_info "Audio is already muted"
        return 0
    fi
    
    print_info "Muting audio..."
    
    if osascript -e "set volume with output muted" 2>/dev/null; then
        print_success "Audio muted 🔇"
        print_info "💡 Use 'mac audio unmute' to restore sound"
    else
        print_error "Failed to mute audio"
        return 1
    fi
}

audio_unmute() {
    local current_state
    current_state=$(audio_get_mute_state)
    
    if [[ "$current_state" == "unmuted" ]]; then
        print_info "Audio is already unmuted"
        return 0
    fi
    
    print_info "Unmuting audio..."
    
    if osascript -e "set volume without output muted" 2>/dev/null; then
        local volume
        volume=$(audio_get_volume)
        print_success "Audio unmuted 🔊 (Volume: $volume%)"
    else
        print_error "Failed to unmute audio"
        return 1
    fi
}

audio_toggle_mute() {
    local current_state
    current_state=$(audio_get_mute_state)
    
    print_info "Current state: Audio is $current_state"
    
    if [[ "$current_state" == "muted" ]]; then
        audio_unmute
    else
        audio_mute
    fi
}

# =============================================================================
# Audio Device Management
# =============================================================================

audio_list_devices() {
    local device_type="${1:-all}"
    
    print_header "🎵 Audio Devices"
    echo ""
    
    case "$device_type" in
        "output"|"out")
            print_subheader "Output Devices:"
            system_profiler SPAudioDataType 2>/dev/null | grep -A 20 "Audio Devices:" | grep -E "^      [A-Z]" | while read -r device; do
                device=$(echo "$device" | sed 's/^      //')
                echo "  $(print_command "$device")"
            done
            ;;
        "input"|"in")
            print_subheader "Input Devices:"
            system_profiler SPAudioDataType 2>/dev/null | grep -A 10 "Input:" | grep -E "Default Input Source" | while read -r line; do
                device=$(echo "$line" | awk -F': ' '{print $2}')
                echo "  $(print_command "$device")"
            done
            ;;
        "all"|*)
            # List all audio devices
            if command_exists SwitchAudioSource; then
                print_subheader "Output Devices:"
                SwitchAudioSource -a -t output 2>/dev/null | while read -r device; do
                    echo "  $(print_command "$device")"
                done
                
                echo ""
                print_subheader "Input Devices:"
                SwitchAudioSource -a -t input 2>/dev/null | while read -r device; do
                    echo "  $(print_command "$device")"
                done
            else
                # Fallback method using system_profiler
                system_profiler SPAudioDataType 2>/dev/null | grep -E "Audio Devices:|Default Output Source|Default Input Source" | while read -r line; do
                    if [[ "$line" == *"Audio Devices:"* ]]; then
                        continue
                    elif [[ "$line" == *"Default Output Source:"* ]]; then
                        device=$(echo "$line" | awk -F': ' '{print $2}')
                        echo "  Output: $(print_command "$device")"
                    elif [[ "$line" == *"Default Input Source:"* ]]; then
                        device=$(echo "$line" | awk -F': ' '{print $2}')
                        echo "  Input: $(print_command "$device")"
                    fi
                done
            fi
            ;;
    esac
    
    if ! command_exists SwitchAudioSource; then
        echo ""
        print_info "💡 Install SwitchAudioSource for enhanced device switching:"
        print_info "    brew install switchaudio-osx"
    fi
}

audio_get_current_output() {
    if command_exists SwitchAudioSource; then
        SwitchAudioSource -c -t output 2>/dev/null
    else
        system_profiler SPAudioDataType 2>/dev/null | grep "Default Output Source:" | awk -F': ' '{print $2}' | trim
    fi
}

audio_get_current_input() {
    if command_exists SwitchAudioSource; then
        SwitchAudioSource -c -t input 2>/dev/null
    else
        system_profiler SPAudioDataType 2>/dev/null | grep "Default Input Source:" | awk -F': ' '{print $2}' | trim
    fi
}

audio_set_output_device() {
    local device_name="$1"
    
    if [[ -z "$device_name" ]]; then
        print_error "Device name required"
        print_info "Usage: mac audio output <device_name>"
        print_info "Use 'mac audio devices' to see available devices"
        return 1
    fi
    
    if command_exists SwitchAudioSource; then
        print_info "Setting output device to: $device_name"
        
        if SwitchAudioSource -s "$device_name" -t output 2>/dev/null; then
            print_success "Output device set to: $device_name"
            
            # Show current volume for new device
            local volume
            volume=$(audio_get_volume)
            print_info "Current volume: $volume%"
        else
            print_error "Failed to set output device"
            print_info "Make sure the device name is correct"
            print_info "Use 'mac audio devices' to see available devices"
            return 1
        fi
    else
        print_error "SwitchAudioSource not found"
        print_info "Install with: brew install switchaudio-osx"
        return 1
    fi
}

audio_set_input_device() {
    local device_name="$1"
    
    if [[ -z "$device_name" ]]; then
        print_error "Device name required"
        print_info "Usage: mac audio input <device_name>"
        print_info "Use 'mac audio devices' to see available devices"
        return 1
    fi
    
    if command_exists SwitchAudioSource; then
        print_info "Setting input device to: $device_name"
        
        if SwitchAudioSource -s "$device_name" -t input 2>/dev/null; then
            print_success "Input device set to: $device_name"
        else
            print_error "Failed to set input device"
            print_info "Make sure the device name is correct"
            return 1
        fi
    else
        print_error "SwitchAudioSource not found"
        print_info "Install with: brew install switchaudio-osx"
        return 1
    fi
}

# =============================================================================
# Audio Balance and Advanced Settings
# =============================================================================

audio_set_balance() {
    local balance="$1"
    
    case "$balance" in
        "left"|"l")
            print_info "Setting audio balance to left channel..."
            osascript -e "set volume output volume (output volume of (get volume settings)) with output balance -1" 2>/dev/null
            print_success "Audio balance set to left channel"
            ;;
        "right"|"r")
            print_info "Setting audio balance to right channel..."
            osascript -e "set volume output volume (output volume of (get volume settings)) with output balance 1" 2>/dev/null
            print_success "Audio balance set to right channel"
            ;;
        "center"|"c"|"middle")
            print_info "Setting audio balance to center..."
            osascript -e "set volume output volume (output volume of (get volume settings)) with output balance 0" 2>/dev/null
            print_success "Audio balance set to center"
            ;;
        *)
            print_error "Invalid balance setting: $balance"
            print_info "Valid settings: left, right, center"
            return 1
            ;;
    esac
}

audio_status() {
    print_header "🎵 Audio Status"
    echo ""
    
    # Volume and mute status
    local volume mute_state
    volume=$(audio_get_volume)
    mute_state=$(audio_get_mute_state)
    
    echo "  Volume: $(print_command "$volume%")"
    echo "  Status: $(print_command "$mute_state")"
    
    # Current devices
    local current_output current_input
    current_output=$(audio_get_current_output)
    current_input=$(audio_get_current_input)
    
    if [[ -n "$current_output" ]]; then
        echo "  Output device: $(print_command "$current_output")"
    fi
    
    if [[ -n "$current_input" ]]; then
        echo "  Input device: $(print_command "$current_input")"
    fi
    
    # Additional audio info
    echo ""
    print_subheader "System Audio:"
    
    # Check if sound effects are enabled
    local sound_effects
    sound_effects=$(defaults read NSGlobalDomain com.apple.sound.uiaudio.enabled 2>/dev/null || echo "1")
    if [[ "$sound_effects" == "1" ]]; then
        echo "  Sound effects: $(print_command "enabled")"
    else
        echo "  Sound effects: $(print_dim "disabled")"
    fi
    
    # Check alert volume
    local alert_volume
    alert_volume=$(osascript -e "alert volume of (get volume settings)" 2>/dev/null || echo "unknown")
    if [[ "$alert_volume" != "unknown" ]]; then
        echo "  Alert volume: $(print_command "$alert_volume%")"
    fi
}

# =============================================================================
# Quick Audio Actions
# =============================================================================

audio_volume_up() {
    local step="${1:-10}"
    local current
    current=$(audio_get_volume)
    local new_volume=$((current + step))
    
    if [[ $new_volume -gt 100 ]]; then
        new_volume=100
    fi
    
    audio_set_volume "$new_volume"
}

audio_volume_down() {
    local step="${1:-10}"
    local current
    current=$(audio_get_volume)
    local new_volume=$((current - step))
    
    if [[ $new_volume -lt 0 ]]; then
        new_volume=0
    fi
    
    audio_set_volume "$new_volume"
}

# =============================================================================
# Module Help System
# =============================================================================

audio_help() {
    print_category_header "audio" "🔊" 75
    
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "volume <0-100>" "Set system volume percentage"
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "mute" "Mute system audio"
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "unmute" "Unmute system audio"
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "toggle-mute" "Toggle mute state"
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "devices [type]" "List audio devices"
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "output <device>" "Set audio output device"
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "input <device>" "Set audio input device"
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "balance <pos>" "Set audio balance (left/right/center)"
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "up [step]" "Increase volume by step (default: 10)"
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "down [step]" "Decrease volume by step (default: 10)"
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "status" "Show audio status"
    
    print_category_footer 75
    echo ""
    
    print_header "Examples:"
    echo "  mac audio volume 75                 # Set volume to 75%"
    echo "  mac audio mute                      # Mute all audio"
    echo "  mac audio devices                   # List all audio devices"
    echo "  mac audio output \"Built-in Speakers\" # Switch output device"
    echo "  mac audio balance left              # Set balance to left"
    echo "  mac audio up 5                      # Increase volume by 5%"
    echo ""
    
    print_header "Global Shortcuts:"
    echo "  mac vol 50                          # Quick volume set"
    echo ""
    
    print_header "Device Types:"
    echo "  all      - Show both input and output devices"
    echo "  output   - Show only output devices"
    echo "  input    - Show only input devices"
    echo ""
    
    print_header "Tips:"
    echo "  • Use quotes around device names with spaces"
    echo "  • Install switchaudio-osx for enhanced device switching"
    echo "  • Balance adjustment helps with accessibility needs"
    echo "  • Volume 0 is silent, 100 is maximum system volume"
    echo ""
}

# Search function for this module
audio_search() {
    local search_term="$1"
    local results=""
    
    if [[ "volume" == *"$search_term"* ]] || [[ "sound" == *"$search_term"* ]]; then
        results+="  mac audio volume <0-100>         Set volume\n"
        results+="  mac audio up/down [step]         Adjust volume\n"
    fi
    
    if [[ "mute" == *"$search_term"* ]] || [[ "silent" == *"$search_term"* ]]; then
        results+="  mac audio mute                   Mute audio\n"
        results+="  mac audio unmute                 Unmute audio\n"
        results+="  mac audio toggle-mute            Toggle mute\n"
    fi
    
    if [[ "device" == *"$search_term"* ]] || [[ "speaker" == *"$search_term"* ]] || [[ "headphone" == *"$search_term"* ]]; then
        results+="  mac audio devices                List audio devices\n"
        results+="  mac audio output <device>        Set output device\n"
        results+="  mac audio input <device>         Set input device\n"
    fi
    
    if [[ "balance" == *"$search_term"* ]] || [[ "left" == *"$search_term"* ]] || [[ "right" == *"$search_term"* ]]; then
        results+="  mac audio balance <pos>          Set audio balance\n"
    fi
    
    if [[ "status" == *"$search_term"* ]] || [[ "info" == *"$search_term"* ]]; then
        results+="  mac audio status                 Show audio status\n"
    fi
    
    if [[ -n "$results" ]]; then
        echo -e "$results"
    fi
}

# =============================================================================
# Module Dispatcher
# =============================================================================

audio_dispatch() {
    local action="$1"
    shift  # Remove action from arguments
    
    case "$action" in
        "volume")
            audio_set_volume "$1"
            ;;
        "mute")
            audio_mute
            ;;
        "unmute")
            audio_unmute
            ;;
        "toggle-mute")
            audio_toggle_mute
            ;;
        "devices")
            audio_list_devices "$1"
            ;;
        "output")
            audio_set_output_device "$1"
            ;;
        "input")
            audio_set_input_device "$1"
            ;;
        "balance")
            audio_set_balance "$1"
            ;;
        "up")
            audio_volume_up "$1"
            ;;
        "down")
            audio_volume_down "$1"
            ;;
        "status")
            audio_status
            ;;
        "help"|"-h"|"--help")
            audio_help
            ;;
        *)
            print_error "Unknown audio action: $action"
            echo ""
            print_info "Available actions: volume, mute, unmute, toggle-mute, devices, output, input, balance, up, down, status"
            print_info "Use 'mac help audio' for detailed information"
            return 1
            ;;
    esac
}

print_debug "Audio module loaded successfully"
