#!/bin/bash

# =============================================================================
# SuperMac - System Module
# =============================================================================
# System information and maintenance commands
# 
# Commands:
#   info           - Comprehensive system information
#   cleanup        - Deep system cleanup
#   battery        - Battery status and health
#   memory         - Memory usage statistics
#   cpu            - CPU usage and information
#   disk-usage     - Disk usage by directory
#   temperature    - System temperature sensors
#   uptime         - System uptime with details
#   login-items    - Manage startup items
#   processes      - Top processes by resource usage
#
# Built by CosmoLabs - https://cosmolabs.org
# =============================================================================

# =============================================================================
# System Information Functions
# =============================================================================

system_info() {
    print_header "🖥️ System Information"
    echo ""
    
    # macOS version and build
    local macos_name version build
    macos_name=$(sw_vers -productName)
    version=$(sw_vers -productVersion)
    build=$(sw_vers -buildVersion)
    
    echo "  OS: $(print_command "$macos_name $version")"
    echo "  Build: $(print_dim "$build")"
    
    # Hardware information
    local model_name chip_info memory_info
    model_name=$(system_profiler SPHardwareDataType 2>/dev/null | grep 'Model Name' | awk -F': ' '{print $2}' | trim)
    chip_info=$(system_profiler SPHardwareDataType 2>/dev/null | grep 'Chip\|Processor Name' | awk -F': ' '{print $2}' | trim | head -1)
    memory_info=$(system_profiler SPHardwareDataType 2>/dev/null | grep 'Memory' | awk -F': ' '{print $2}' | trim)
    
    if [[ -n "$model_name" ]]; then
        echo "  Model: $(print_command "$model_name")"
    fi
    
    if [[ -n "$chip_info" ]]; then
        echo "  Processor: $(print_command "$chip_info")"
    fi
    
    if [[ -n "$memory_info" ]]; then
        echo "  Memory: $(print_command "$memory_info")"
    fi
    
    # Architecture
    local arch
    arch=$(uname -m)
    echo "  Architecture: $(print_command "$arch")"
    
    # Uptime
    local uptime_info
    uptime_info=$(uptime | awk -F'up ' '{print $2}' | awk -F', load' '{print $1}' | trim)
    echo "  Uptime: $(print_command "$uptime_info")"
    
    # Shell
    echo "  Shell: $(print_command "$(basename "$SHELL")")"
    
    # Disk usage summary
    echo ""
    print_subheader "Storage:"
    df -h / | tail -1 | while read -r filesystem size used available capacity mount; do
        echo "  Used: $(print_command "$used") of $(print_command "$size") ($(print_dim "$capacity"))"
        echo "  Free: $(print_command "$available")"
    done
}

system_detailed_info() {
    system_info
    echo ""
    
    # Additional detailed information
    print_subheader "Hardware Details:"
    
    # CPU cores
    local cpu_cores
    cpu_cores=$(sysctl -n hw.ncpu 2>/dev/null)
    echo "  CPU Cores: $(print_command "$cpu_cores")"
    
    # Memory details
    local physical_memory
    physical_memory=$(sysctl -n hw.memsize 2>/dev/null | awk '{print int($1/1024/1024/1024) " GB"}')
    echo "  Physical Memory: $(print_command "$physical_memory")"
    
    # System serial number
    local serial
    serial=$(system_profiler SPHardwareDataType 2>/dev/null | grep 'Serial Number' | awk -F': ' '{print $2}' | trim)
    if [[ -n "$serial" ]]; then
        echo "  Serial Number: $(print_dim "$serial")"
    fi
}

# =============================================================================
# System Cleanup Functions
# =============================================================================

system_cleanup() {
    print_header "🧹 System Cleanup"
    echo ""
    
    print_warning "This will clean caches, logs, and temporary files"
    if ! confirm "Continue with system cleanup?" "y"; then
        print_info "Cleanup cancelled"
        return 0
    fi
    
    echo ""
    local cleaned_size=0
    
    # User caches
    if [[ -d "$HOME/Library/Caches" ]]; then
        print_info "Cleaning user caches..."
        local cache_size
        cache_size=$(du -sk "$HOME/Library/Caches" 2>/dev/null | awk '{print $1}')
        
        find "$HOME/Library/Caches" -type f -atime +7 -delete 2>/dev/null || true
        find "$HOME/Library/Caches" -type d -empty -delete 2>/dev/null || true
        
        cleaned_size=$((cleaned_size + cache_size))
        print_success "User caches cleaned"
    fi
    
    # Downloads cleanup (files older than 30 days)
    if [[ -d "$HOME/Downloads" ]]; then
        print_info "Cleaning old downloads (30+ days)..."
        local downloads_count
        downloads_count=$(find "$HOME/Downloads" -type f -mtime +30 2>/dev/null | wc -l | trim)
        
        if [[ "$downloads_count" -gt 0 ]]; then
            find "$HOME/Downloads" -type f -mtime +30 -delete 2>/dev/null || true
            print_success "Cleaned $downloads_count old download files"
        else
            print_info "No old downloads to clean"
        fi
    fi
    
    # Trash
    print_info "Emptying trash..."
    osascript -e 'tell application "Finder" to empty trash' 2>/dev/null || true
    print_success "Trash emptied"
    
    # System logs (requires admin privileges)
    if [[ "$EUID" -eq 0 ]] || sudo -n true 2>/dev/null; then
        print_info "Cleaning system logs (admin required)..."
        sudo find /var/log -name "*.log" -mtime +7 -delete 2>/dev/null || true
        sudo find /private/var/log -name "*.log" -mtime +7 -delete 2>/dev/null || true
        print_success "System logs cleaned"
    else
        print_info "Skipping system logs (no admin privileges)"
    fi
    
    # Browser caches (Safari)
    if [[ -d "$HOME/Library/Caches/com.apple.Safari" ]]; then
        print_info "Cleaning Safari cache..."
        rm -rf "$HOME/Library/Caches/com.apple.Safari"/* 2>/dev/null || true
        print_success "Safari cache cleaned"
    fi
    
    # Temporary files
    print_info "Cleaning temporary files..."
    rm -rf /tmp/* 2>/dev/null || true
    print_success "Temporary files cleaned"
    
    # Font caches
    print_info "Clearing font caches..."
    sudo atsutil databases -remove 2>/dev/null || true
    print_success "Font caches cleared"
    
    echo ""
    print_success "System cleanup completed!"
    print_info "💡 Consider restarting applications to free additional memory"
    print_info "💡 Use 'mac system memory' to check current memory usage"
}

# =============================================================================
# Battery Information
# =============================================================================

system_battery() {
    # Check if this is a laptop
    if ! system_profiler SPPowerDataType >/dev/null 2>&1; then
        print_info "No battery information available (desktop Mac)"
        return 0
    fi
    
    print_header "🔋 Battery Information"
    echo ""
    
    # Get battery info using pmset
    local battery_info
    battery_info=$(pmset -g batt 2>/dev/null)
    
    if [[ -n "$battery_info" ]]; then
        # Parse battery percentage and status
        local percentage status time_remaining
        percentage=$(echo "$battery_info" | grep -E "InternalBattery" | awk '{print $3}' | sed 's/;//')
        status=$(echo "$battery_info" | grep -E "InternalBattery" | awk '{print $4}' | sed 's/;//')
        time_remaining=$(echo "$battery_info" | grep -E "InternalBattery" | awk '{for(i=5;i<=NF;i++) printf "%s ", $i; print ""}' | trim)
        
        echo "  Charge: $(print_command "$percentage")"
        echo "  Status: $(print_command "$status")"
        
        if [[ -n "$time_remaining" && "$time_remaining" != "(no estimate)" ]]; then
            echo "  Time remaining: $(print_command "$time_remaining")"
        fi
    fi
    
    # Battery health information
    local cycle_count max_capacity
    cycle_count=$(system_profiler SPPowerDataType 2>/dev/null | grep "Cycle Count" | awk -F': ' '{print $2}' | trim)
    max_capacity=$(system_profiler SPPowerDataType 2>/dev/null | grep "Maximum Capacity" | awk -F': ' '{print $2}' | trim)
    
    if [[ -n "$cycle_count" ]]; then
        echo "  Cycle count: $(print_command "$cycle_count")"
    fi
    
    if [[ -n "$max_capacity" ]]; then
        echo "  Maximum capacity: $(print_command "$max_capacity")"
        
        # Health assessment
        local capacity_num
        capacity_num=$(echo "$max_capacity" | sed 's/%//')
        if [[ "$capacity_num" -gt 80 ]]; then
            print_success "Battery health: Good"
        elif [[ "$capacity_num" -gt 60 ]]; then
            print_warning "Battery health: Fair"
        else
            print_warning "Battery health: Poor (consider replacement)"
        fi
    fi
    
    # Power adapter status
    local adapter_info
    adapter_info=$(pmset -g ac 2>/dev/null)
    if [[ "$adapter_info" == *"AC Power"* ]]; then
        echo "  Power adapter: $(print_command "Connected")"
    else
        echo "  Power adapter: $(print_dim "Not connected")"
    fi
}

# =============================================================================
# Memory Information
# =============================================================================

system_memory() {
    print_header "💾 Memory Usage"
    echo ""
    
    # Get memory statistics
    local memory_info
    memory_info=$(vm_stat)
    
    local page_size=4096
    local pages_free pages_active pages_inactive pages_wired pages_compressed
    
    pages_free=$(echo "$memory_info" | grep "Pages free" | awk '{print $3}' | sed 's/\.//')
    pages_active=$(echo "$memory_info" | grep "Pages active" | awk '{print $3}' | sed 's/\.//')
    pages_inactive=$(echo "$memory_info" | grep "Pages inactive" | awk '{print $3}' | sed 's/\.//')
    pages_wired=$(echo "$memory_info" | grep "Pages wired down" | awk '{print $4}' | sed 's/\.//')
    pages_compressed=$(echo "$memory_info" | grep "Pages stored in compressor" | awk '{print $5}' | sed 's/\.//' 2>/dev/null || echo "0")
    
    # Convert to MB
    local free_mb active_mb inactive_mb wired_mb compressed_mb total_mb used_mb
    free_mb=$((pages_free * page_size / 1024 / 1024))
    active_mb=$((pages_active * page_size / 1024 / 1024))
    inactive_mb=$((pages_inactive * page_size / 1024 / 1024))
    wired_mb=$((pages_wired * page_size / 1024 / 1024))
    compressed_mb=$((pages_compressed * page_size / 1024 / 1024))
    
    used_mb=$((active_mb + inactive_mb + wired_mb))
    total_mb=$((free_mb + used_mb))
    
    echo "  Total: $(print_command "${total_mb} MB")"
    echo "  Used: $(print_command "${used_mb} MB")"
    echo "  Free: $(print_command "${free_mb} MB")"
    echo ""
    echo "  Active: $(print_command "${active_mb} MB")"
    echo "  Inactive: $(print_command "${inactive_mb} MB")"
    echo "  Wired: $(print_command "${wired_mb} MB")"
    
    if [[ $compressed_mb -gt 0 ]]; then
        echo "  Compressed: $(print_command "${compressed_mb} MB")"
    fi
    
    # Memory pressure
    local memory_pressure
    memory_pressure=$(memory_pressure 2>/dev/null | grep "System-wide memory free percentage" | awk '{print $NF}' | sed 's/%//')
    
    if [[ -n "$memory_pressure" ]]; then
        echo ""
        if [[ "$memory_pressure" -gt 20 ]]; then
            print_success "Memory pressure: Low ($memory_pressure% free)"
        elif [[ "$memory_pressure" -gt 10 ]]; then
            print_warning "Memory pressure: Medium ($memory_pressure% free)"
        else
            print_warning "Memory pressure: High ($memory_pressure% free)"
            print_info "💡 Consider closing some applications"
        fi
    fi
}

# =============================================================================
# CPU Information
# =============================================================================

system_cpu() {
    print_header "⚡ CPU Information"
    echo ""
    
    # CPU model and specs
    local cpu_name cpu_cores cpu_threads
    cpu_name=$(sysctl -n machdep.cpu.brand_string 2>/dev/null)
    cpu_cores=$(sysctl -n hw.physicalcpu 2>/dev/null)
    cpu_threads=$(sysctl -n hw.logicalcpu 2>/dev/null)
    
    if [[ -n "$cpu_name" ]]; then
        echo "  Processor: $(print_command "$cpu_name")"
    fi
    
    echo "  Physical cores: $(print_command "$cpu_cores")"
    echo "  Logical cores: $(print_command "$cpu_threads")"
    
    # CPU usage
    echo ""
    print_subheader "Current Usage:"
    
    # Get CPU usage from top
    local cpu_usage
    cpu_usage=$(top -l 1 -n 0 | grep "CPU usage" | awk '{print $3, $5, $7}')
    echo "  $cpu_usage"
    
    # Load average
    local load_avg
    load_avg=$(uptime | awk -F'load averages: ' '{print $2}')
    echo "  Load average: $(print_command "$load_avg")"
}

# =============================================================================
# Disk Usage
# =============================================================================

system_disk_usage() {
    local target_dir="${1:-$HOME}"
    
    if [[ ! -d "$target_dir" ]]; then
        print_error "Directory not found: $target_dir"
        return 1
    fi
    
    print_header "💿 Disk Usage: $target_dir"
    echo ""
    
    # Overall disk usage
    df -h "$target_dir" | tail -1 | while read -r filesystem size used available capacity mount; do
        echo "  Volume: $(print_command "$filesystem")"
        echo "  Total: $(print_command "$size")"
        echo "  Used: $(print_command "$used")"
        echo "  Available: $(print_command "$available")"
        echo "  Usage: $(print_command "$capacity")"
    done
    
    echo ""
    print_subheader "Largest directories in $target_dir:"
    
    # Find largest directories
    du -sh "$target_dir"/* 2>/dev/null | sort -hr | head -10 | while read -r size path; do
        printf "  %-10s %s\n" "$size" "$(basename "$path")"
    done
}

# =============================================================================
# System Processes
# =============================================================================

system_processes() {
    local sort_by="${1:-cpu}"
    
    print_header "🔄 Top Processes"
    echo ""
    
    case "$sort_by" in
        "cpu")
            print_subheader "By CPU Usage:"
            ps aux | sort -nr -k 3 | head -10 | while read -r user pid cpu mem vsz rss tt stat started time command; do
                printf "  %-20s %5s%% %5s%% %s\n" "$(echo "$command" | cut -d' ' -f1 | xargs basename)" "$cpu" "$mem" "$pid"
            done
            ;;
        "memory"|"mem")
            print_subheader "By Memory Usage:"
            ps aux | sort -nr -k 4 | head -10 | while read -r user pid cpu mem vsz rss tt stat started time command; do
                printf "  %-20s %5s%% %5s%% %s\n" "$(echo "$command" | cut -d' ' -f1 | xargs basename)" "$mem" "$cpu" "$pid"
            done
            ;;
        *)
            print_error "Unknown sort option: $sort_by"
            print_info "Available options: cpu, memory"
            return 1
            ;;
    esac
}

# =============================================================================
# Module Help System
# =============================================================================

system_help() {
    print_category_header "system" "🖥️" 70
    
    printf "${PURPLE}│${NC}  %-20s %-40s ${PURPLE}│${NC}\n" "info" "Comprehensive system information"
    printf "${PURPLE}│${NC}  %-20s %-40s ${PURPLE}│${NC}\n" "cleanup" "Deep system cleanup"
    printf "${PURPLE}│${NC}  %-20s %-40s ${PURPLE}│${NC}\n" "battery" "Battery status and health"
    printf "${PURPLE}│${NC}  %-20s %-40s ${PURPLE}│${NC}\n" "memory" "Memory usage statistics"
    printf "${PURPLE}│${NC}  %-20s %-40s ${PURPLE}│${NC}\n" "cpu" "CPU usage and information"
    printf "${PURPLE}│${NC}  %-20s %-40s ${PURPLE}│${NC}\n" "disk-usage [dir]" "Disk usage analysis"
    printf "${PURPLE}│${NC}  %-20s %-40s ${PURPLE}│${NC}\n" "processes [sort]" "Top processes (cpu/memory)"
    printf "${PURPLE}│${NC}  %-20s %-40s ${PURPLE}│${NC}\n" "uptime" "System uptime information"
    
    print_category_footer 70
    echo ""
    
    print_header "Examples:"
    echo "  mac system info                     # Complete system overview"
    echo "  mac system cleanup                  # Clean temporary files"
    echo "  mac system battery                  # Battery health check"
    echo "  mac system memory                   # Memory usage details"
    echo "  mac system disk-usage ~/Downloads   # Analyze Downloads folder"
    echo ""
    
    print_header "Global Shortcuts:"
    echo "  mac cleanup                         # Quick system cleanup"
    echo ""
    
    print_header "Maintenance Tips:"
    echo "  • Run cleanup weekly to free disk space"
    echo "  • Monitor battery health on laptops"
    echo "  • Check memory usage if system feels slow"
    echo "  • Use disk-usage to find large files"
    echo ""
}

# Search function for this module
system_search() {
    local search_term="$1"
    local results=""
    
    if [[ "cleanup" == *"$search_term"* ]] || [[ "clean" == *"$search_term"* ]]; then
        results+="  mac system cleanup               Deep system cleanup\n"
    fi
    
    if [[ "battery" == *"$search_term"* ]] || [[ "power" == *"$search_term"* ]]; then
        results+="  mac system battery               Battery status\n"
    fi
    
    if [[ "memory" == *"$search_term"* ]] || [[ "ram" == *"$search_term"* ]]; then
        results+="  mac system memory                Memory usage\n"
    fi
    
    if [[ "cpu" == *"$search_term"* ]] || [[ "processor" == *"$search_term"* ]]; then
        results+="  mac system cpu                   CPU information\n"
    fi
    
    if [[ "disk" == *"$search_term"* ]] || [[ "storage" == *"$search_term"* ]]; then
        results+="  mac system disk-usage            Disk usage\n"
    fi
    
    if [[ "info" == *"$search_term"* ]] || [[ "status" == *"$search_term"* ]]; then
        results+="  mac system info                  System information\n"
    fi
    
    if [[ "process" == *"$search_term"* ]] || [[ "top" == *"$search_term"* ]]; then
        results+="  mac system processes             Top processes\n"
    fi
    
    if [[ -n "$results" ]]; then
        echo -e "$results"
    fi
}

# =============================================================================
# Module Dispatcher
# =============================================================================

system_dispatch() {
    local action="$1"
    shift  # Remove action from arguments
    
    case "$action" in
        "info")
            system_info
            ;;
        "detailed-info")
            system_detailed_info
            ;;
        "cleanup")
            system_cleanup
            ;;
        "battery")
            system_battery
            ;;
        "memory"|"mem")
            system_memory
            ;;
        "cpu")
            system_cpu
            ;;
        "disk-usage")
            system_disk_usage "$1"
            ;;
        "processes")
            system_processes "$1"
            ;;
        "uptime")
            uptime
            ;;
        "help"|"-h"|"--help")
            system_help
            ;;
        *)
            print_error "Unknown system action: $action"
            echo ""
            print_info "Available actions: info, cleanup, battery, memory, cpu, disk-usage, processes, uptime"
            print_info "Use 'mac help system' for detailed information"
            return 1
            ;;
    esac
}

print_debug "System module loaded successfully"
