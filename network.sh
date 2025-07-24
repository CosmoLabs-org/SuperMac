#!/bin/bash

# =============================================================================
# SuperMac - Network Module
# =============================================================================
# Network information and troubleshooting commands
# 
# Commands:
#   ip             - Show local IP address
#   public-ip      - Show public IP address
#   info           - Comprehensive network information
#   flush-dns      - Clear DNS cache
#   ping <host>    - Ping with enhanced stats
#   speed-test     - Basic network speed test
#   reset          - Reset network settings
#   renew-dhcp     - Renew DHCP lease
#   locations      - Manage network locations
#
# Built by CosmoLabs - https://cosmolabs.org
# =============================================================================

# =============================================================================
# IP Address Functions
# =============================================================================

network_get_local_ip() {
    # Try multiple interfaces to find active connection
    local interfaces=("en0" "en1" "en2" "en3")
    
    for interface in "${interfaces[@]}"; do
        local ip
        ip=$(ipconfig getifaddr "$interface" 2>/dev/null)
        if [[ -n "$ip" ]]; then
            echo "$ip"
            return 0
        fi
    done
    
    return 1
}

network_get_interface_for_ip() {
    local interfaces=("en0" "en1" "en2" "en3")
    
    for interface in "${interfaces[@]}"; do
        local ip
        ip=$(ipconfig getifaddr "$interface" 2>/dev/null)
        if [[ -n "$ip" ]]; then
            echo "$interface"
            return 0
        fi
    done
    
    return 1
}

network_ip() {
    local ip
    ip=$(network_get_local_ip)
    
    if [[ -n "$ip" ]]; then
        print_success "Local IP address: $ip"
        
        # Show which interface
        local interface
        interface=$(network_get_interface_for_ip)
        if [[ -n "$interface" ]]; then
            print_info "Interface: $interface"
        fi
    else
        print_warning "No active network connection found"
        print_info "Make sure you're connected to WiFi or Ethernet"
        return 1
    fi
}

network_public_ip() {
    print_info "Fetching public IP address..."
    
    local public_ip=""
    local services=(
        "https://ifconfig.me"
        "https://ipinfo.io/ip"
        "https://api.ipify.org"
        "https://checkip.amazonaws.com"
    )
    
    for service in "${services[@]}"; do
        if command_exists curl; then
            public_ip=$(curl -s --connect-timeout 10 "$service" 2>/dev/null | trim)
        elif command_exists wget; then
            public_ip=$(wget -qO- --timeout=10 "$service" 2>/dev/null | trim)
        fi
        
        # Validate IP format
        if [[ "$public_ip" =~ ^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$ ]]; then
            break
        fi
        public_ip=""
    done
    
    if [[ -n "$public_ip" ]]; then
        print_success "Public IP address: $public_ip"
        
        # Try to get location info
        network_get_ip_location "$public_ip"
    else
        print_error "Failed to retrieve public IP address"
        print_info "Check your internet connection"
        return 1
    fi
}

network_get_ip_location() {
    local ip="$1"
    
    if ! command_exists curl; then
        return 1
    fi
    
    local location_info
    location_info=$(curl -s --connect-timeout 5 "https://ipinfo.io/$ip/json" 2>/dev/null)
    
    if [[ -n "$location_info" ]]; then
        local city country isp
        city=$(echo "$location_info" | grep '"city"' | cut -d'"' -f4 2>/dev/null)
        country=$(echo "$location_info" | grep '"country"' | cut -d'"' -f4 2>/dev/null)
        isp=$(echo "$location_info" | grep '"org"' | cut -d'"' -f4 2>/dev/null)
        
        if [[ -n "$city" && -n "$country" ]]; then
            echo "  Location: $(print_command "$city, $country")"
        fi
        
        if [[ -n "$isp" ]]; then
            echo "  ISP: $(print_dim "$isp")"
        fi
    fi
}

# =============================================================================
# Network Information
# =============================================================================

network_info() {
    print_header "📡 Network Information"
    echo ""
    
    # Local IP and interface
    local ip interface
    ip=$(network_get_local_ip)
    interface=$(network_get_interface_for_ip)
    
    if [[ -n "$ip" ]]; then
        echo "  Local IP: $(print_command "$ip")"
        echo "  Interface: $(print_command "$interface")"
    else
        echo "  Local IP: $(print_dim "Not connected")"
    fi
    
    # Gateway
    local gateway
    gateway=$(route -n get default 2>/dev/null | grep gateway | awk '{print $2}')
    if [[ -n "$gateway" ]]; then
        echo "  Gateway: $(print_command "$gateway")"
    fi
    
    # DNS servers
    local dns_servers
    dns_servers=$(scutil --dns 2>/dev/null | grep nameserver | head -3 | awk '{print $3}' | tr '\n' ' ')
    if [[ -n "$dns_servers" ]]; then
        echo "  DNS: $(print_command "$dns_servers")"
    fi
    
    # WiFi status if applicable
    if [[ "$interface" == "en0" ]] || [[ "$interface" == "en1" ]]; then
        echo ""
        print_subheader "WiFi Details:"
        
        # Load wifi module functions if available
        if declare -f wifi_get_current_network >/dev/null 2>&1; then
            local network
            network=$(wifi_get_current_network)
            if [[ "$network" != "Not connected" ]]; then
                echo "  Network: $(print_command "$network")"
            fi
        fi
    fi
    
    # Network speed test option
    echo ""
    print_info "💡 Use 'mac network speed-test' for connection speed"
    print_info "💡 Use 'mac network public-ip' for external IP"
}

# =============================================================================
# DNS Functions
# =============================================================================

network_flush_dns() {
    print_info "Flushing DNS cache..."
    
    # Clear system DNS cache
    sudo dscacheutil -flushcache
    
    # Kill and restart mDNSResponder
    sudo killall -HUP mDNSResponder 2>/dev/null
    
    # Clear additional caches on newer macOS versions
    if command_exists discoveryutil; then
        sudo discoveryutil mdnsflushcache 2>/dev/null
        sudo discoveryutil udnsflushcaches 2>/dev/null
    fi
    
    print_success "DNS cache cleared successfully!"
    print_info "💡 This can resolve DNS-related connectivity issues"
}

# =============================================================================
# Network Testing
# =============================================================================

network_ping() {
    local host="$1"
    local count="${2:-5}"
    
    if [[ -z "$host" ]]; then
        print_error "Host required"
        print_info "Usage: mac network ping <host> [count]"
        return 1
    fi
    
    print_info "Pinging $host ($count packets)..."
    echo ""
    
    if ping -c "$count" "$host" 2>/dev/null; then
        echo ""
        print_success "Ping completed successfully"
    else
        echo ""
        print_error "Ping failed - host may be unreachable"
        return 1
    fi
}

network_speed_test() {
    print_info "Running basic network speed test..."
    print_warning "This is a simple test - use dedicated tools for accurate measurements"
    echo ""
    
    # Test with a small file download
    local test_file="https://httpbin.org/bytes/1024"
    local start_time end_time duration
    
    if command_exists curl; then
        print_info "Testing download speed..."
        start_time=$(date +%s.%N)
        
        if curl -s -o /dev/null --connect-timeout 10 "$test_file"; then
            end_time=$(date +%s.%N)
            duration=$(echo "$end_time - $start_time" | bc 2>/dev/null || echo "1")
            
            print_success "Basic connectivity test passed"
            print_info "Response time: ${duration}s"
        else
            print_error "Speed test failed"
            return 1
        fi
    else
        print_error "curl not available for speed test"
        return 1
    fi
    
    echo ""
    print_info "💡 For accurate speed tests, use:"
    print_info "  • Speedtest.net"
    print_info "  • Fast.com"
    print_info "  • Network utility apps"
}

# =============================================================================
# Network Management
# =============================================================================

network_renew_dhcp() {
    local interface
    interface=$(network_get_interface_for_ip)
    
    if [[ -z "$interface" ]]; then
        print_error "No active network interface found"
        return 1
    fi
    
    print_info "Renewing DHCP lease on $interface..."
    
    if sudo ipconfig set "$interface" DHCP 2>/dev/null; then
        sleep 3
        local new_ip
        new_ip=$(ipconfig getifaddr "$interface" 2>/dev/null)
        
        if [[ -n "$new_ip" ]]; then
            print_success "DHCP lease renewed successfully"
            print_info "New IP address: $new_ip"
        else
            print_warning "DHCP renewal completed but no IP assigned"
        fi
    else
        print_error "Failed to renew DHCP lease"
        return 1
    fi
}

network_reset() {
    print_warning "This will reset all network settings to defaults"
    
    if ! confirm "Are you sure you want to reset network settings?" "n"; then
        print_info "Network reset cancelled"
        return 0
    fi
    
    print_info "Resetting network settings..."
    
    # Remove network preferences
    sudo rm -f /Library/Preferences/SystemConfiguration/NetworkInterfaces.plist 2>/dev/null
    sudo rm -f /Library/Preferences/SystemConfiguration/preferences.plist 2>/dev/null
    
    # Restart network services
    sudo launchctl unload /System/Library/LaunchDaemons/com.apple.networkd.plist 2>/dev/null
    sudo launchctl load /System/Library/LaunchDaemons/com.apple.networkd.plist 2>/dev/null
    
    print_success "Network settings reset"
    print_warning "You may need to reconfigure WiFi networks and other settings"
    print_info "Consider restarting your Mac for complete reset"
}

network_locations() {
    local action="$1"
    
    case "$action" in
        "list"|"")
            print_header "📍 Network Locations"
            echo ""
            
            networksetup -listlocations 2>/dev/null | while read -r location; do
                if [[ -n "$location" ]]; then
                    echo "  $(print_command "$location")"
                fi
            done
            ;;
        "current")
            local current
            current=$(networksetup -getcurrentlocation 2>/dev/null)
            echo "Current location: $(print_command "$current")"
            ;;
        *)
            print_error "Invalid locations action: $action"
            print_info "Usage: mac network locations [list|current]"
            return 1
            ;;
    esac
}

# =============================================================================
# Module Help System
# =============================================================================

network_help() {
    print_category_header "network" "📡" 70
    
    printf "${PURPLE}│${NC}  %-20s %-40s ${PURPLE}│${NC}\n" "ip" "Show local IP address"
    printf "${PURPLE}│${NC}  %-20s %-40s ${PURPLE}│${NC}\n" "public-ip" "Show public IP address"
    printf "${PURPLE}│${NC}  %-20s %-40s ${PURPLE}│${NC}\n" "info" "Comprehensive network information"
    printf "${PURPLE}│${NC}  %-20s %-40s ${PURPLE}│${NC}\n" "flush-dns" "Clear DNS cache"
    printf "${PURPLE}│${NC}  %-20s %-40s ${PURPLE}│${NC}\n" "ping <host>" "Ping with enhanced statistics"
    printf "${PURPLE}│${NC}  %-20s %-40s ${PURPLE}│${NC}\n" "speed-test" "Basic network speed test"
    printf "${PURPLE}│${NC}  %-20s %-40s ${PURPLE}│${NC}\n" "renew-dhcp" "Renew DHCP lease"
    printf "${PURPLE}│${NC}  %-20s %-40s ${PURPLE}│${NC}\n" "reset" "Reset network settings"
    printf "${PURPLE}│${NC}  %-20s %-40s ${PURPLE}│${NC}\n" "locations" "Manage network locations"
    
    print_category_footer 70
    echo ""
    
    print_header "Examples:"
    echo "  mac network ip                      # Quick IP lookup"
    echo "  mac network public-ip               # External IP with location"
    echo "  mac network ping google.com         # Test connectivity"
    echo "  mac network flush-dns               # Fix DNS issues"
    echo "  mac network info                    # Complete network status"
    echo ""
    
    print_header "Global Shortcuts:"
    echo "  mac ip                              # Quick local IP"
    echo ""
    
    print_header "Troubleshooting:"
    echo "  • Use flush-dns for website loading issues"
    echo "  • Use renew-dhcp for IP address problems"
    echo "  • Use ping to test connectivity to specific hosts"
    echo "  • Use reset as last resort (requires reconfiguration)"
    echo ""
}

# Search function for this module
network_search() {
    local search_term="$1"
    local results=""
    
    if [[ "ip" == *"$search_term"* ]] || [[ "address" == *"$search_term"* ]]; then
        results+="  mac network ip                   Show local IP\n"
        results+="  mac network public-ip            Show public IP\n"
    fi
    
    if [[ "dns" == *"$search_term"* ]] || [[ "cache" == *"$search_term"* ]]; then
        results+="  mac network flush-dns            Clear DNS cache\n"
    fi
    
    if [[ "ping" == *"$search_term"* ]] || [[ "test" == *"$search_term"* ]] || [[ "speed" == *"$search_term"* ]]; then
        results+="  mac network ping <host>          Test connectivity\n"
        results+="  mac network speed-test           Network speed test\n"
    fi
    
    if [[ "reset" == *"$search_term"* ]] || [[ "dhcp" == *"$search_term"* ]]; then
        results+="  mac network renew-dhcp           Renew DHCP lease\n"
        results+="  mac network reset                Reset settings\n"
    fi
    
    if [[ "info" == *"$search_term"* ]] || [[ "status" == *"$search_term"* ]]; then
        results+="  mac network info                 Network information\n"
    fi
    
    if [[ -n "$results" ]]; then
        echo -e "$results"
    fi
}

# =============================================================================
# Module Dispatcher
# =============================================================================

network_dispatch() {
    local action="$1"
    shift  # Remove action from arguments
    
    case "$action" in
        "ip")
            network_ip
            ;;
        "public-ip")
            network_public_ip
            ;;
        "info")
            network_info
            ;;
        "flush-dns")
            network_flush_dns
            ;;
        "ping")
            network_ping "$@"
            ;;
        "speed-test")
            network_speed_test
            ;;
        "renew-dhcp")
            network_renew_dhcp
            ;;
        "reset")
            network_reset
            ;;
        "locations")
            network_locations "$1"
            ;;
        "help"|"-h"|"--help")
            network_help
            ;;
        *)
            print_error "Unknown network action: $action"
            echo ""
            print_info "Available actions: ip, public-ip, info, flush-dns, ping, speed-test, renew-dhcp, reset, locations"
            print_info "Use 'mac help network' for detailed information"
            return 1
            ;;
    esac
}

print_debug "Network module loaded successfully"
