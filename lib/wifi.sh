#!/bin/bash

# =============================================================================
# SuperMac - WiFi Module
# =============================================================================
# WiFi control and management commands
# 
# Commands:
#   on             - Turn WiFi on
#   off            - Turn WiFi off
#   toggle         - Toggle WiFi state
#   status         - Show WiFi connection status
#   scan           - Scan for available networks
#   connect <name> - Connect to network
#   forget <name>  - Forget saved network
#   info           - Detailed connection information
#
# Built by CosmoLabs - https://cosmolabs.org
# =============================================================================

# =============================================================================
# WiFi Power Management
# =============================================================================

wifi_get_interface() {
    # Get the WiFi interface name (usually en0 or en1)
    networksetup -listallhardwareports | awk '/Wi-Fi|AirPort/{getline; print $2}' | head -1
}

wifi_get_power_state() {
    local interface
    interface=$(wifi_get_interface)
    
    if [[ -z "$interface" ]]; then
        echo "unknown"
        return 1
    fi
    
    local state
    state=$(networksetup -getairportpower "$interface" 2>/dev/null | awk '{print $NF}')
    echo "$state"
}

wifi_on() {
    local interface
    interface=$(wifi_get_interface)
    
    if [[ -z "$interface" ]]; then
        print_error "WiFi interface not found"
        return 1
    fi
    
    local current_state
    current_state=$(wifi_get_power_state)
    
    if [[ "$current_state" == "On" ]]; then
        print_info "WiFi is already on"
        return 0
    fi
    
    print_info "Turning WiFi on..."
    networksetup -setairportpower "$interface" on
    
    # Wait for WiFi to come online
    local timeout=10
    local count=0
    while [[ $count -lt $timeout ]]; do
        sleep 1
        if [[ "$(wifi_get_power_state)" == "On" ]]; then
            break
        fi
        ((count++))
    done
    
    if [[ "$(wifi_get_power_state)" == "On" ]]; then
        print_success "WiFi is now on!"
        print_info "💡 Scanning for available networks..."
        sleep 2
        wifi_show_current_connection
    else
        print_error "Failed to turn WiFi on"
        return 1
    fi
}

wifi_off() {
    local interface
    interface=$(wifi_get_interface)
    
    if [[ -z "$interface" ]]; then
        print_error "WiFi interface not found"
        return 1
    fi
    
    local current_state
    current_state=$(wifi_get_power_state)
    
    if [[ "$current_state" == "Off" ]]; then
        print_info "WiFi is already off"
        return 0
    fi
    
    print_info "Turning WiFi off..."
    networksetup -setairportpower "$interface" off
    
    # Wait for WiFi to turn off
    sleep 2
    
    if [[ "$(wifi_get_power_state)" == "Off" ]]; then
        print_success "WiFi is now off!"
        print_info "💡 Use 'mac wifi on' to turn it back on"
    else
        print_error "Failed to turn WiFi off"
        return 1
    fi
}

wifi_toggle() {
    local current_state
    current_state=$(wifi_get_power_state)
    
    print_info "Current WiFi state: $current_state"
    
    case "$current_state" in
        "On")
            wifi_off
            ;;
        "Off")
            wifi_on
            ;;
        *)
            print_error "Unable to determine WiFi state"
            return 1
            ;;
    esac
}

# =============================================================================
# WiFi Network Information
# =============================================================================

wifi_get_current_network() {
    local interface
    interface=$(wifi_get_interface)
    
    if [[ -z "$interface" ]]; then
        return 1
    fi
    
    local network
    network=$(networksetup -getairportnetwork "$interface" 2>/dev/null | sed 's/Current Wi-Fi Network: //')
    
    if [[ "$network" == "You are not associated with an AirPort network." ]]; then
        echo "Not connected"
    else
        echo "$network"
    fi
}

wifi_show_current_connection() {
    local interface
    interface=$(wifi_get_interface)
    local network
    network=$(wifi_get_current_network)
    
    if [[ "$network" == "Not connected" ]]; then
        print_warning "Not connected to any WiFi network"
        return 0
    fi
    
    echo "  Connected to: $(print_command "$network")"
    
    # Get signal strength if available
    if command_exists airport; then
        local signal
        signal=$(airport -I 2>/dev/null | awk '/agrCtlRSSI/{print $2}')
        if [[ -n "$signal" ]]; then
            echo "  Signal strength: $(print_command "${signal} dBm")"
        fi
    fi
}

wifi_status() {
    print_header "🌐 WiFi Status"
    echo ""
    
    local interface
    interface=$(wifi_get_interface)
    
    if [[ -z "$interface" ]]; then
        print_error "WiFi interface not found"
        return 1
    fi
    
    echo "  Interface: $(print_command "$interface")"
    
    local power_state
    power_state=$(wifi_get_power_state)
    echo "  Power: $(print_command "$power_state")"
    
    if [[ "$power_state" == "On" ]]; then
        wifi_show_current_connection
        
        # Show IP address if connected
        local ip
        ip=$(ipconfig getifaddr "$interface" 2>/dev/null)
        if [[ -n "$ip" ]]; then
            echo "  IP Address: $(print_command "$ip")"
        fi
    fi
}

wifi_info() {
    print_header "🌐 Detailed WiFi Information"
    echo ""
    
    local interface
    interface=$(wifi_get_interface)
    
    if [[ -z "$interface" ]]; then
        print_error "WiFi interface not found"
        return 1
    fi
    
    # Basic status
    wifi_status
    echo ""
    
    # Detailed information if connected
    local network
    network=$(wifi_get_current_network)
    
    if [[ "$network" != "Not connected" ]]; then
        print_header "Connection Details:"
        
        # Try to get detailed info using airport command line tool
        if command_exists airport; then
            airport -I 2>/dev/null | while read -r line; do
                case "$line" in
                    *"SSID"*) echo "  Network: $(echo "$line" | cut -d: -f2 | trim)" ;;
                    *"BSSID"*) echo "  Router: $(echo "$line" | cut -d: -f2 | trim)" ;;
                    *"channel"*) echo "  Channel: $(echo "$line" | cut -d: -f2 | trim)" ;;
                    *"CC"*) echo "  Country: $(echo "$line" | cut -d: -f2 | trim)" ;;
                esac
            done
        fi
        
        # Network configuration
        echo ""
        print_header "Network Configuration:"
        
        local gateway
        gateway=$(route -n get default 2>/dev/null | grep gateway | awk '{print $2}')
        if [[ -n "$gateway" ]]; then
            echo "  Gateway: $(print_command "$gateway")"
        fi
        
        local dns
        dns=$(scutil --dns | grep nameserver | head -3 | awk '{print $3}' | tr '\n' ' ')
        if [[ -n "$dns" ]]; then
            echo "  DNS: $(print_command "$dns")"
        fi
    fi
}

# =============================================================================
# Network Scanning and Management
# =============================================================================

wifi_scan() {
    local interface
    interface=$(wifi_get_interface)
    
    if [[ -z "$interface" ]]; then
        print_error "WiFi interface not found"
        return 1
    fi
    
    if [[ "$(wifi_get_power_state)" != "On" ]]; then
        print_error "WiFi is turned off"
        print_info "Turn on WiFi first: mac wifi on"
        return 1
    fi
    
    print_info "Scanning for available WiFi networks..."
    
    # Use airport command for scanning if available
    if command_exists airport; then
        airport -s 2>/dev/null | head -20 | while read -r line; do
            if [[ -n "$line" && "$line" != *"SSID"* ]]; then
                local ssid
                ssid=$(echo "$line" | awk '{print $1}')
                local signal
                signal=$(echo "$line" | awk '{print $3}')
                printf "  %-30s %s\n" "$ssid" "$(print_dim "$signal dBm")"
            fi
        done
    else
        # Fallback method
        networksetup -listpreferredwirelessnetworks "$interface" 2>/dev/null | while read -r network; do
            if [[ -n "$network" && "$network" != *"Preferred networks"* ]]; then
                echo "  $(print_command "$network")"
            fi
        done
    fi
}

wifi_connect() {
    local network_name="$1"
    
    if [[ -z "$network_name" ]]; then
        print_error "Network name required"
        print_info "Usage: mac wifi connect <network_name>"
        return 1
    fi
    
    local interface
    interface=$(wifi_get_interface)
    
    if [[ -z "$interface" ]]; then
        print_error "WiFi interface not found"
        return 1
    fi
    
    if [[ "$(wifi_get_power_state)" != "On" ]]; then
        print_error "WiFi is turned off"
        print_info "Turn on WiFi first: mac wifi on"
        return 1
    fi
    
    print_info "Connecting to network: $network_name"
    
    # Try to connect (will prompt for password if needed)
    if networksetup -setairportnetwork "$interface" "$network_name" 2>/dev/null; then
        print_success "Connected to $network_name!"
    else
        print_error "Failed to connect to $network_name"
        print_info "Make sure the network name is correct and you have the password"
        return 1
    fi
}

wifi_forget() {
    local network_name="$1"
    
    if [[ -z "$network_name" ]]; then
        print_error "Network name required"
        print_info "Usage: mac wifi forget <network_name>"
        return 1
    fi
    
    local interface
    interface=$(wifi_get_interface)
    
    if [[ -z "$interface" ]]; then
        print_error "WiFi interface not found"
        return 1
    fi
    
    print_info "Forgetting network: $network_name"
    
    if networksetup -removepreferredwirelessnetwork "$interface" "$network_name" 2>/dev/null; then
        print_success "Forgot network: $network_name"
        print_info "💡 You'll need to re-enter the password to reconnect"
    else
        print_error "Failed to forget network: $network_name"
        print_info "Make sure the network name is correct"
        return 1
    fi
}

wifi_list_saved() {
    local interface
    interface=$(wifi_get_interface)
    
    if [[ -z "$interface" ]]; then
        print_error "WiFi interface not found"
        return 1
    fi
    
    print_header "💾 Saved WiFi Networks"
    echo ""
    
    networksetup -listpreferredwirelessnetworks "$interface" 2>/dev/null | while read -r line; do
        if [[ -n "$line" && "$line" != *"Preferred networks"* ]]; then
            # Remove leading whitespace and tabs
            network=$(echo "$line" | sed 's/^[[:space:]]*//')
            echo "  $(print_command "$network")"
        fi
    done
}

# =============================================================================
# Module Help System
# =============================================================================

wifi_help() {
    print_category_header "wifi" "🌐" 70
    
    printf "${PURPLE}│${NC}  %-20s %-40s ${PURPLE}│${NC}\n" "on" "Turn WiFi on"
    printf "${PURPLE}│${NC}  %-20s %-40s ${PURPLE}│${NC}\n" "off" "Turn WiFi off"
    printf "${PURPLE}│${NC}  %-20s %-40s ${PURPLE}│${NC}\n" "toggle" "Toggle WiFi state"
    printf "${PURPLE}│${NC}  %-20s %-40s ${PURPLE}│${NC}\n" "status" "Show WiFi status"
    printf "${PURPLE}│${NC}  %-20s %-40s ${PURPLE}│${NC}\n" "info" "Detailed connection information"
    printf "${PURPLE}│${NC}  %-20s %-40s ${PURPLE}│${NC}\n" "scan" "Scan for available networks"
    printf "${PURPLE}│${NC}  %-20s %-40s ${PURPLE}│${NC}\n" "connect <name>" "Connect to WiFi network"
    printf "${PURPLE}│${NC}  %-20s %-40s ${PURPLE}│${NC}\n" "forget <name>" "Forget saved network"
    printf "${PURPLE}│${NC}  %-20s %-40s ${PURPLE}│${NC}\n" "list-saved" "Show saved networks"
    
    print_category_footer 70
    echo ""
    
    print_header "Examples:"
    echo "  mac wifi toggle                     # Quick WiFi toggle"
    echo "  mac wifi scan                       # See available networks"
    echo "  mac wifi connect \"Coffee Shop\"      # Connect to network"
    echo "  mac wifi forget \"Old Network\"       # Remove saved network"
    echo "  mac wifi info                       # Detailed connection info"
    echo ""
    
    print_header "Tips:"
    echo "  • Toggle is fastest way to reconnect to problematic networks"
    echo "  • Scan shows signal strength for available networks"
    echo "  • Forget networks you no longer need for security"
    echo "  • Info shows detailed technical information"
    echo ""
}

# Search function for this module
wifi_search() {
    local search_term="$1"
    local results=""
    
    if [[ "wifi" == *"$search_term"* ]] || [[ "wireless" == *"$search_term"* ]]; then
        results+="  mac wifi on/off/toggle          WiFi power control\n"
        results+="  mac wifi status                 WiFi status\n"
    fi
    
    if [[ "network" == *"$search_term"* ]] || [[ "connect" == *"$search_term"* ]]; then
        results+="  mac wifi connect <name>         Connect to network\n"
        results+="  mac wifi scan                   Scan networks\n"
    fi
    
    if [[ "forget" == *"$search_term"* ]] || [[ "remove" == *"$search_term"* ]]; then
        results+="  mac wifi forget <name>          Forget network\n"
    fi
    
    if [[ "info" == *"$search_term"* ]] || [[ "status" == *"$search_term"* ]]; then
        results+="  mac wifi info                   Detailed WiFi info\n"
        results+="  mac wifi status                 WiFi status\n"
    fi
    
    if [[ -n "$results" ]]; then
        echo -e "$results"
    fi
}

# =============================================================================
# Module Dispatcher
# =============================================================================

wifi_dispatch() {
    local action="$1"
    shift  # Remove action from arguments
    
    case "$action" in
        "on"|"enable")
            wifi_on
            ;;
        "off"|"disable")
            wifi_off
            ;;
        "toggle")
            wifi_toggle
            ;;
        "status")
            wifi_status
            ;;
        "info")
            wifi_info
            ;;
        "scan")
            wifi_scan
            ;;
        "connect")
            wifi_connect "$1"
            ;;
        "forget")
            wifi_forget "$1"
            ;;
        "list-saved"|"saved")
            wifi_list_saved
            ;;
        "help"|"-h"|"--help")
            wifi_help
            ;;
        *)
            print_error "Unknown wifi action: $action"
            echo ""
            print_info "Available actions: on, off, toggle, status, info, scan, connect, forget, list-saved"
            print_info "Use 'mac help wifi' for detailed information"
            return 1
            ;;
    esac
}

print_debug "WiFi module loaded successfully"
