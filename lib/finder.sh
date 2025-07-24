#!/bin/bash

# =============================================================================
# SuperMac - Finder Module
# =============================================================================
# File visibility and Finder management commands
# 
# Commands:
#   restart       - Restart Finder application
#   show-hidden   - Show hidden files
#   hide-hidden   - Hide hidden files  
#   toggle-hidden - Toggle hidden file visibility
#
# Built by CosmoLabs - https://cosmolabs.org
# =============================================================================

# =============================================================================
# Finder Module Commands
# =============================================================================

finder_restart() {
    print_info "Restarting Finder..."
    
    if ! pgrep -f "Finder" >/dev/null 2>&1; then
        print_warning "Finder doesn't appear to be running. Starting it..."
        open -a Finder
        print_success "Finder started successfully!"
        return 0
    fi
    
    if killall Finder 2>/dev/null; then
        sleep 1
        
        # Wait for Finder to restart
        local timeout=5
        local count=0
        while ! pgrep -f "Finder" >/dev/null 2>&1 && [[ $count -lt $timeout ]]; do
            sleep 1
            ((count++))
        done
        
        if pgrep -f "Finder" >/dev/null 2>&1; then
            print_success "Finder restarted successfully!"
            print_info "💡 Tip: This often fixes unresponsive Finder windows"
        else
            print_warning "Finder may not have restarted properly. Opening manually..."
            open -a Finder
            print_success "Finder opened manually"
        fi
    else
        print_error "Failed to restart Finder"
        print_info "You may need to restart Finder manually from Activity Monitor"
        return 1
    fi
}

# Get current hidden files state
finder_get_hidden_state() {
    local current_state
    current_state=$(defaults read com.apple.finder AppleShowAllFiles 2>/dev/null || echo "FALSE")
    
    case "$(to_upper "$current_state")" in
        "TRUE"|"1"|"YES")
            echo "visible"
            ;;
        *)
            echo "hidden"
            ;;
    esac
}

finder_show_hidden() {
    local current_state
    current_state=$(finder_get_hidden_state)
    
    if [[ "$current_state" == "visible" ]]; then
        print_info "Hidden files are already visible"
        return 0
    fi
    
    print_info "Making hidden files visible..."
    defaults write com.apple.finder AppleShowAllFiles TRUE
    
    if killall Finder 2>/dev/null; then
        sleep 1
        print_success "Hidden files are now visible!"
        print_info "💡 You can now see files like .gitignore, .env, .DS_Store"
        print_info "💡 Use 'mac finder hide-hidden' to hide them again"
    else
        print_error "Failed to restart Finder"
        print_info "Changes may not take effect until Finder is restarted"
        return 1
    fi
}

finder_hide_hidden() {
    local current_state
    current_state=$(finder_get_hidden_state)
    
    if [[ "$current_state" == "hidden" ]]; then
        print_info "Hidden files are already hidden"
        return 0
    fi
    
    print_info "Hiding hidden files..."
    defaults write com.apple.finder AppleShowAllFiles FALSE
    
    if killall Finder 2>/dev/null; then
        sleep 1
        print_success "Hidden files are now hidden!"
        print_info "💡 Clean view restored - system files are tucked away"
        print_info "💡 Use 'mac finder show-hidden' to show them again"
    else
        print_error "Failed to restart Finder"
        print_info "Changes may not take effect until Finder is restarted"
        return 1
    fi
}

finder_toggle_hidden() {
    local current_state
    current_state=$(finder_get_hidden_state)
    
    print_info "Current state: Hidden files are $current_state"
    
    if [[ "$current_state" == "visible" ]]; then
        defaults write com.apple.finder AppleShowAllFiles FALSE
        print_info "Setting hidden files to HIDDEN"
    else
        defaults write com.apple.finder AppleShowAllFiles TRUE
        print_info "Setting hidden files to VISIBLE"
    fi
    
    print_info "Applying changes..."
    if killall Finder 2>/dev/null; then
        sleep 1
        local new_state
        new_state=$(finder_get_hidden_state)
        print_success "Hidden files are now $new_state!"
        
        if [[ "$new_state" == "visible" ]]; then
            print_info "💡 You can now see system files and hidden folders"
        else
            print_info "💡 System files are now hidden for a cleaner view"
        fi
    else
        print_error "Failed to restart Finder"
        return 1
    fi
}

finder_reveal() {
    local target="$1"
    
    if [[ -z "$target" ]]; then
        print_error "Path required"
        print_info "Usage: mac finder reveal <path>"
        return 1
    fi
    
    if [[ ! -e "$target" ]]; then
        print_error "Path does not exist: $target"
        return 1
    fi
    
    print_info "Revealing '$target' in Finder..."
    open -R "$target"
    print_success "Path revealed in Finder!"
}

finder_status() {
    print_header "📁 Finder Status"
    echo ""
    
    # Check if Finder is running
    if pgrep -f "Finder" >/dev/null 2>&1; then
        print_success "Finder is running"
    else
        print_warning "Finder is not running"
    fi
    
    # Check hidden files state
    local hidden_state
    hidden_state=$(finder_get_hidden_state)
    echo "  Hidden files: $(print_command "$hidden_state")"
    
    # Check Finder preferences
    local show_extensions
    show_extensions=$(defaults read NSGlobalDomain AppleShowAllExtensions 2>/dev/null || echo "0")
    if [[ "$show_extensions" == "1" ]]; then
        echo "  File extensions: $(print_command "visible")"
    else
        echo "  File extensions: $(print_dim "hidden")"
    fi
    
    # Show default view style
    local view_style
    view_style=$(defaults read com.apple.finder FXPreferredViewStyle 2>/dev/null || echo "icnv")
    case "$view_style" in
        "icnv") echo "  Default view: $(print_command "Icon View")" ;;
        "Nlsv") echo "  Default view: $(print_command "List View")" ;;
        "clmv") echo "  Default view: $(print_command "Column View")" ;;
        "Flwv") echo "  Default view: $(print_command "Gallery View")" ;;
        *) echo "  Default view: $(print_dim "Unknown")" ;;
    esac
}

# =============================================================================
# Module Help System
# =============================================================================

finder_help() {
    print_category_header "finder" "📁" 65
    
    printf "${PURPLE}│${NC}  %-20s %-35s ${PURPLE}│${NC}\n" "restart" "Restart Finder application"
    printf "${PURPLE}│${NC}  %-20s %-35s ${PURPLE}│${NC}\n" "show-hidden" "Show hidden files and folders"
    printf "${PURPLE}│${NC}  %-20s %-35s ${PURPLE}│${NC}\n" "hide-hidden" "Hide hidden files and folders"
    printf "${PURPLE}│${NC}  %-20s %-35s ${PURPLE}│${NC}\n" "toggle-hidden" "Toggle hidden file visibility"
    printf "${PURPLE}│${NC}  %-20s %-35s ${PURPLE}│${NC}\n" "reveal <path>" "Reveal file/folder in Finder"
    printf "${PURPLE}│${NC}  %-20s %-35s ${PURPLE}│${NC}\n" "status" "Show Finder status and settings"
    
    print_category_footer 65
    echo ""
    
    print_header "Examples:"
    echo "  mac finder restart              # Fix unresponsive Finder"
    echo "  mac finder toggle-hidden        # Quick toggle hidden files"
    echo "  mac finder reveal ~/.ssh        # Show .ssh folder in Finder"
    echo "  mac finder status               # Check current settings"
    echo ""
    
    print_header "Shortcuts:"
    echo "  mac restart-finder              # Global shortcut for restart"
    echo ""
    
    print_header "Tips:"
    echo "  • Hidden files include system files like .DS_Store, .gitignore"
    echo "  • Showing hidden files is useful for development work"
    echo "  • Restart Finder if windows become unresponsive"
    echo "  • Use reveal to quickly navigate to specific files"
    echo ""
}

# Search function for this module
finder_search() {
    local search_term="$1"
    local results=""
    
    # Search through command names and descriptions
    if [[ "restart" == *"$search_term"* ]] || [[ "finder" == *"$search_term"* ]]; then
        results+="  mac finder restart              Restart Finder application\n"
    fi
    
    if [[ "hidden" == *"$search_term"* ]] || [[ "show" == *"$search_term"* ]] || [[ "hide" == *"$search_term"* ]]; then
        results+="  mac finder show-hidden          Show hidden files\n"
        results+="  mac finder hide-hidden          Hide hidden files\n"
        results+="  mac finder toggle-hidden        Toggle hidden files\n"
    fi
    
    if [[ "reveal" == *"$search_term"* ]] || [[ "open" == *"$search_term"* ]]; then
        results+="  mac finder reveal <path>        Reveal in Finder\n"
    fi
    
    if [[ "status" == *"$search_term"* ]] || [[ "info" == *"$search_term"* ]]; then
        results+="  mac finder status               Show Finder status\n"
    fi
    
    if [[ -n "$results" ]]; then
        echo -e "$results"
    fi
}

# =============================================================================
# Module Dispatcher
# =============================================================================

finder_dispatch() {
    local action="$1"
    shift  # Remove action from arguments
    
    case "$action" in
        "restart")
            finder_restart
            ;;
        "show-hidden")
            finder_show_hidden
            ;;
        "hide-hidden")
            finder_hide_hidden
            ;;
        "toggle-hidden")
            finder_toggle_hidden
            ;;
        "reveal")
            finder_reveal "$@"
            ;;
        "status")
            finder_status
            ;;
        "help"|"-h"|"--help")
            finder_help
            ;;
        *)
            print_error "Unknown finder action: $action"
            echo ""
            print_info "Available actions: restart, show-hidden, hide-hidden, toggle-hidden, reveal, status"
            print_info "Use 'mac help finder' for detailed information"
            return 1
            ;;
    esac
}

print_debug "Finder module loaded successfully"
