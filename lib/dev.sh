#!/bin/bash

# =============================================================================
# SuperMac - Developer Module
# =============================================================================
# Developer tools and utilities for enhanced productivity
# 
# Commands:
#   kill-port <port>       - Kill process on specific port
#   list-ports             - Show all processes using ports
#   servers                - List running development servers
#   localhost <port>       - Open localhost in browser
#   serve <dir>            - Start HTTP server in directory
#   processes              - Enhanced process viewer
#   cpu-hogs               - Show CPU-intensive processes
#   memory-hogs            - Show memory-intensive processes
#   json-format <file>     - Format JSON file
#   base64-encode <text>   - Base64 encode text
#   base64-decode <text>   - Base64 decode text
#   uuid                   - Generate UUID
#   password <length>      - Generate secure password
#
# Built by CosmoLabs - https://cosmolabs.org
# =============================================================================

# =============================================================================
# Port Management Functions
# =============================================================================

dev_kill_port() {
    local port="$1"
    
    if [[ -z "$port" ]]; then
        print_error "Port number required"
        print_info "Usage: mac dev kill-port <port>"
        return 1
    fi
    
    if ! is_number "$port"; then
        print_error "Invalid port number: $port"
        return 1
    fi
    
    print_info "Looking for processes on port $port..."
    
    local pids
    pids=$(lsof -ti:$port 2>/dev/null)
    
    if [[ -n "$pids" ]]; then
        echo "$pids" | while read -r pid; do
            local process_name
            process_name=$(ps -p "$pid" -o comm= 2>/dev/null | xargs basename)
            
            print_info "Found process: $process_name (PID: $pid)"
            
            if kill -9 "$pid" 2>/dev/null; then
                print_success "Killed process $process_name (PID: $pid) on port $port"
            else
                print_error "Failed to kill process $pid"
            fi
        done
    else
        print_warning "No process found running on port $port"
        return 1
    fi
}

dev_list_ports() {
    print_header "🔗 Active Network Ports"
    echo ""
    
    # Get listening ports with process information
    print_subheader "Listening Ports:"
    lsof -i -P -n | grep LISTEN | while read -r command pid user fd type device size_off node name; do
        local port
        port=$(echo "$name" | awk -F: '{print $NF}')
        printf "  %-6s %-20s %-10s %s\n" "$port" "$(echo "$command" | cut -c1-20)" "$pid" "$user"
    done | sort -n
    
    echo ""
    print_subheader "Common Development Ports:"
    local dev_ports=(3000 3001 4000 5000 8000 8080 8888 9000 9001)
    local found_any=false
    
    for port in "${dev_ports[@]}"; do
        local process_info
        process_info=$(lsof -ti:$port 2>/dev/null)
        if [[ -n "$process_info" ]]; then
            local pid="$process_info"
            local process_name
            process_name=$(ps -p "$pid" -o comm= 2>/dev/null | xargs basename)
            printf "  %-6s %-20s %s\n" "$port" "$process_name" "$pid"
            found_any=true
        fi
    done
    
    if [[ "$found_any" != true ]]; then
        print_info "No processes found on common development ports"
    fi
}

# =============================================================================
# Development Server Management
# =============================================================================

dev_servers() {
    print_header "🚀 Running Development Servers"
    echo ""
    
    # Common development ports and their typical uses
    declare -A common_ports=(
        [3000]="React/Next.js"
        [3001]="React (alt)"
        [4000]="Gatsby/Express"
        [5000]="Flask/Express"
        [5173]="Vite"
        [8000]="Django/Python"
        [8080]="Webpack/Tomcat"
        [8888]="Jupyter"
        [9000]="PHP/Node"
        [9001]="SvelteKit"
    )
    
    local found_servers=false
    
    for port in $(printf '%s\n' "${!common_ports[@]}" | sort -n); do
        local pid
        pid=$(lsof -ti:$port 2>/dev/null)
        
        if [[ -n "$pid" ]]; then
            local process_name command_line
            process_name=$(ps -p "$pid" -o comm= 2>/dev/null | xargs basename)
            command_line=$(ps -p "$pid" -o args= 2>/dev/null | cut -c1-50)
            
            printf "  %-6s %-15s %-12s %s\n" "$port" "${common_ports[$port]}" "$process_name" "$pid"
            printf "         Command: %s\n" "$(print_dim "$command_line")"
            echo ""
            found_servers=true
        fi
    done
    
    if [[ "$found_servers" != true ]]; then
        print_info "No development servers found on common ports"
        echo ""
        print_info "💡 Try 'mac dev list-ports' to see all active ports"
    fi
    
    echo ""
    print_subheader "Quick Actions:"
    echo "  mac dev kill-port <port>            Kill server on port"
    echo "  mac dev localhost <port>            Open in browser"
    echo "  mac dev serve <directory>           Start HTTP server"
}

dev_localhost() {
    local port="$1"
    local protocol="${2:-http}"
    
    if [[ -z "$port" ]]; then
        print_error "Port number required"
        print_info "Usage: mac dev localhost <port> [protocol]"
        return 1
    fi
    
    if ! is_number "$port"; then
        print_error "Invalid port number: $port"
        return 1
    fi
    
    local url="$protocol://localhost:$port"
    
    # Check if anything is running on the port
    if ! lsof -ti:$port >/dev/null 2>&1; then
        print_warning "No service detected on port $port"
        if ! confirm "Open $url anyway?" "y"; then
            return 0
        fi
    fi
    
    print_info "Opening $url in default browser..."
    
    if open "$url" 2>/dev/null; then
        print_success "Browser opened!"
    else
        print_error "Failed to open browser"
        print_info "URL: $url"
        return 1
    fi
}

dev_serve() {
    local directory="${1:-$(pwd)}"
    local port="${2:-8000}"
    
    if [[ ! -d "$directory" ]]; then
        print_error "Directory not found: $directory"
        return 1
    fi
    
    # Check if port is already in use
    if lsof -ti:$port >/dev/null 2>&1; then
        print_error "Port $port is already in use"
        print_info "Use 'mac dev kill-port $port' to free it"
        return 1
    fi
    
    print_info "Starting HTTP server on port $port..."
    print_info "Serving directory: $directory"
    print_info "URL: http://localhost:$port"
    print_info "Press Ctrl+C to stop"
    echo ""
    
    cd "$directory" || return 1
    
    # Try different Python versions
    if command_exists python3; then
        python3 -m http.server "$port"
    elif command_exists python; then
        python -m http.server "$port"
    else
        print_error "Python not found - cannot start HTTP server"
        return 1
    fi
}

# =============================================================================
# Process Management
# =============================================================================

dev_processes() {
    local sort_type="${1:-cpu}"
    local count="${2:-15}"
    
    print_header "🔄 System Processes"
    echo ""
    
    case "$sort_type" in
        "cpu")
            print_subheader "Top $count processes by CPU usage:"
            ps aux | sort -nr -k 3 | head -"$count" | while read -r user pid cpu mem vsz rss tt stat started time command; do
                local cmd_name
                cmd_name=$(echo "$command" | awk '{print $1}' | xargs basename)
                printf "  %-20s %6s%% %6s%% %8s %s\n" "$cmd_name" "$cpu" "$mem" "$pid" "$user"
            done
            ;;
        "memory"|"mem")
            print_subheader "Top $count processes by memory usage:"
            ps aux | sort -nr -k 4 | head -"$count" | while read -r user pid cpu mem vsz rss tt stat started time command; do
                local cmd_name
                cmd_name=$(echo "$command" | awk '{print $1}' | xargs basename)
                printf "  %-20s %6s%% %6s%% %8s %s\n" "$cmd_name" "$mem" "$cpu" "$pid" "$user"
            done
            ;;
        "all")
            print_subheader "All running processes (last $count):"
            ps aux | tail -"$count" | while read -r user pid cpu mem vsz rss tt stat started time command; do
                local cmd_name
                cmd_name=$(echo "$command" | awk '{print $1}' | xargs basename)
                printf "  %-20s %6s%% %6s%% %8s %s\n" "$cmd_name" "$cpu" "$mem" "$pid" "$user"
            done
            ;;
        *)
            print_error "Unknown sort type: $sort_type"
            print_info "Available types: cpu, memory, all"
            return 1
            ;;
    esac
    
    echo ""
    print_info "💡 Use 'mac dev cpu-hogs' or 'mac dev memory-hogs' for focused views"
}

dev_cpu_hogs() {
    print_header "🔥 CPU-Intensive Processes"
    echo ""
    
    ps aux | awk '$3 > 1.0' | sort -nr -k 3 | head -10 | while read -r user pid cpu mem vsz rss tt stat started time command; do
        local cmd_name
        cmd_name=$(echo "$command" | awk '{print $1}' | xargs basename)
        
        if [[ "$cpu" > 5.0 ]]; then
            printf "  %-20s ${RED}%6s%%${NC} %6s%% %8s\n" "$cmd_name" "$cpu" "$mem" "$pid"
        elif [[ "$cpu" > 2.0 ]]; then
            printf "  %-20s ${YELLOW}%6s%%${NC} %6s%% %8s\n" "$cmd_name" "$cpu" "$mem" "$pid"
        else
            printf "  %-20s %6s%% %6s%% %8s\n" "$cmd_name" "$cpu" "$mem" "$pid"
        fi
    done
    
    echo ""
    print_info "💡 High CPU usage may indicate runaway processes"
    print_info "💡 Use 'mac dev kill-port <port>' to stop development servers"
}

dev_memory_hogs() {
    print_header "💾 Memory-Intensive Processes"
    echo ""
    
    ps aux | awk '$4 > 1.0' | sort -nr -k 4 | head -10 | while read -r user pid cpu mem vsz rss tt stat started time command; do
        local cmd_name
        cmd_name=$(echo "$command" | awk '{print $1}' | xargs basename)
        
        if (( $(echo "$mem > 10.0" | bc -l) )); then
            printf "  %-20s %6s%% ${RED}%6s%%${NC} %8s\n" "$cmd_name" "$cpu" "$mem" "$pid"
        elif (( $(echo "$mem > 5.0" | bc -l) )); then
            printf "  %-20s %6s%% ${YELLOW}%6s%%${NC} %8s\n" "$cmd_name" "$cpu" "$mem" "$pid"
        else
            printf "  %-20s %6s%% %6s%% %8s\n" "$cmd_name" "$cpu" "$mem" "$pid"
        fi
    done
    
    echo ""
    print_info "💡 High memory usage may slow down your system"
    print_info "💡 Consider closing unused applications"
}

# =============================================================================
# Developer Utilities
# =============================================================================

dev_json_format() {
    local file="$1"
    
    if [[ -z "$file" ]]; then
        print_error "File path required"
        print_info "Usage: mac dev json-format <file>"
        return 1
    fi
    
    if [[ ! -f "$file" ]]; then
        print_error "File not found: $file"
        return 1
    fi
    
    print_info "Formatting JSON file: $file"
    
    if command_exists jq; then
        if jq '.' "$file" > "${file}.formatted" 2>/dev/null; then
            mv "${file}.formatted" "$file"
            print_success "JSON file formatted successfully"
        else
            rm -f "${file}.formatted"
            print_error "Invalid JSON in file"
            return 1
        fi
    elif command_exists python3; then
        if python3 -m json.tool "$file" "${file}.formatted" 2>/dev/null; then
            mv "${file}.formatted" "$file"
            print_success "JSON file formatted successfully"
        else
            rm -f "${file}.formatted"
            print_error "Invalid JSON in file"
            return 1
        fi
    else
        print_error "Neither jq nor Python available for JSON formatting"
        return 1
    fi
}

dev_base64_encode() {
    local text="$1"
    
    if [[ -z "$text" ]]; then
        print_error "Text required"
        print_info "Usage: mac dev base64-encode <text>"
        return 1
    fi
    
    local encoded
    encoded=$(echo -n "$text" | base64)
    
    print_success "Base64 encoded:"
    echo "  $encoded"
    
    # Copy to clipboard if available
    if command_exists pbcopy; then
        echo -n "$encoded" | pbcopy
        print_info "💡 Copied to clipboard"
    fi
}

dev_base64_decode() {
    local encoded="$1"
    
    if [[ -z "$encoded" ]]; then
        print_error "Encoded text required"
        print_info "Usage: mac dev base64-decode <encoded_text>"
        return 1
    fi
    
    local decoded
    if decoded=$(echo "$encoded" | base64 -d 2>/dev/null); then
        print_success "Base64 decoded:"
        echo "  $decoded"
        
        # Copy to clipboard if available
        if command_exists pbcopy; then
            echo -n "$decoded" | pbcopy
            print_info "💡 Copied to clipboard"
        fi
    else
        print_error "Invalid Base64 encoding"
        return 1
    fi
}

dev_uuid() {
    local uuid
    
    if command_exists uuidgen; then
        uuid=$(uuidgen | tr '[:upper:]' '[:lower:]')
    else
        # Fallback UUID generation
        uuid=$(cat /proc/sys/kernel/random/uuid 2>/dev/null || echo "$(date +%s)-$(shuf -i 1000-9999 -n 1)")
    fi
    
    print_success "Generated UUID:"
    echo "  $uuid"
    
    # Copy to clipboard if available
    if command_exists pbcopy; then
        echo -n "$uuid" | pbcopy
        print_info "💡 Copied to clipboard"
    fi
}

dev_password() {
    local length="${1:-16}"
    
    if ! is_number "$length"; then
        print_error "Invalid length: $length"
        return 1
    fi
    
    if [[ "$length" -lt 4 ]] || [[ "$length" -gt 128 ]]; then
        print_error "Length must be between 4 and 128"
        return 1
    fi
    
    local password
    
    # Generate secure password with mixed characters
    if command_exists openssl; then
        password=$(openssl rand -base64 "$((length * 3 / 4))" | tr -d "=+/" | cut -c1-"$length")
    else
        # Fallback method
        password=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9!@#$%^&*()_+-=' | head -c "$length")
    fi
    
    print_success "Generated password (length: $length):"
    echo "  $password"
    
    # Copy to clipboard if available
    if command_exists pbcopy; then
        echo -n "$password" | pbcopy
        print_info "💡 Copied to clipboard"
    fi
    
    print_warning "🔒 Store this password securely"
}

# =============================================================================
# Module Help System
# =============================================================================

dev_help() {
    print_category_header "dev" "💻" 75
    
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "kill-port <port>" "Kill process on specific port"
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "list-ports" "Show all processes using ports"
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "servers" "List running development servers"
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "localhost <port>" "Open localhost in browser"
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "serve <dir>" "Start HTTP server in directory"
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "processes [sort]" "Enhanced process viewer"
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "cpu-hogs" "Show CPU-intensive processes"
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "memory-hogs" "Show memory-intensive processes"
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "json-format <file>" "Format JSON file"
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "base64-encode <text>" "Base64 encode text"
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "base64-decode <text>" "Base64 decode text"
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "uuid" "Generate UUID"
    printf "${PURPLE}│${NC}  %-25s %-40s ${PURPLE}│${NC}\n" "password [length]" "Generate secure password"
    
    print_category_footer 75
    echo ""
    
    print_header "Examples:"
    echo "  mac dev kill-port 3000              # Kill React dev server"
    echo "  mac dev servers                     # See all running servers"
    echo "  mac dev localhost 8080              # Open localhost:8080"
    echo "  mac dev serve ~/my-site             # Start HTTP server"
    echo "  mac dev password 32                 # Generate 32-char password"
    echo ""
    
    print_header "Global Shortcuts:"
    echo "  mac kp 3000                         # Quick kill-port"
    echo ""
    
    print_header "Development Workflow:"
    echo "  • Use servers to see what's running"
    echo "  • Kill ports before starting new services"
    echo "  • Use serve for quick static file hosting"
    echo "  • Monitor resource usage with cpu-hogs/memory-hogs"
    echo ""
}

# Search function for this module
dev_search() {
    local search_term="$1"
    local results=""
    
    if [[ "port" == *"$search_term"* ]] || [[ "kill" == *"$search_term"* ]]; then
        results+="  mac dev kill-port <port>         Kill process on port\n"
        results+="  mac dev list-ports               Show active ports\n"
    fi
    
    if [[ "server" == *"$search_term"* ]] || [[ "localhost" == *"$search_term"* ]]; then
        results+="  mac dev servers                  List dev servers\n"
        results+="  mac dev localhost <port>         Open in browser\n"
        results+="  mac dev serve <dir>              Start HTTP server\n"
    fi
    
    if [[ "process" == *"$search_term"* ]] || [[ "cpu" == *"$search_term"* ]] || [[ "memory" == *"$search_term"* ]]; then
        results+="  mac dev processes                Process viewer\n"
        results+="  mac dev cpu-hogs                 CPU intensive processes\n"
        results+="  mac dev memory-hogs              Memory intensive processes\n"
    fi
    
    if [[ "json" == *"$search_term"* ]] || [[ "format" == *"$search_term"* ]]; then
        results+="  mac dev json-format <file>       Format JSON file\n"
    fi
    
    if [[ "base64" == *"$search_term"* ]] || [[ "encode" == *"$search_term"* ]]; then
        results+="  mac dev base64-encode <text>     Base64 encode\n"
        results+="  mac dev base64-decode <text>     Base64 decode\n"
    fi
    
    if [[ "uuid" == *"$search_term"* ]] || [[ "password" == *"$search_term"* ]]; then
        results+="  mac dev uuid                     Generate UUID\n"
        results+="  mac dev password [length]        Generate password\n"
    fi
    
    if [[ -n "$results" ]]; then
        echo -e "$results"
    fi
}

# =============================================================================
# Module Dispatcher
# =============================================================================

dev_dispatch() {
    local action="$1"
    shift  # Remove action from arguments
    
    case "$action" in
        "kill-port")
            dev_kill_port "$1"
            ;;
        "list-ports")
            dev_list_ports
            ;;
        "servers")
            dev_servers
            ;;
        "localhost")
            dev_localhost "$@"
            ;;
        "serve")
            dev_serve "$@"
            ;;
        "processes")
            dev_processes "$@"
            ;;
        "cpu-hogs")
            dev_cpu_hogs
            ;;
        "memory-hogs")
            dev_memory_hogs
            ;;
        "json-format")
            dev_json_format "$1"
            ;;
        "base64-encode")
            dev_base64_encode "$1"
            ;;
        "base64-decode")
            dev_base64_decode "$1"
            ;;
        "uuid")
            dev_uuid
            ;;
        "password")
            dev_password "$1"
            ;;
        "help"|"-h"|"--help")
            dev_help
            ;;
        *)
            print_error "Unknown dev action: $action"
            echo ""
            print_info "Available actions: kill-port, list-ports, servers, localhost, serve, processes, cpu-hogs, memory-hogs"
            print_info "Utilities: json-format, base64-encode, base64-decode, uuid, password"
            print_info "Use 'mac help dev' for detailed information"
            return 1
            ;;
    esac
}

print_debug "Developer module loaded successfully"
