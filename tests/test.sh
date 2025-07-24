#!/bin/bash

# =============================================================================
# SuperMac Test Suite
# =============================================================================
# Comprehensive testing for SuperMac modular architecture
# 
# Usage: bash test.sh [module]
# 
# Built by CosmoLabs - https://cosmolabs.org
# =============================================================================

set -e

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SUPERMAC_ROOT="$(dirname "$SCRIPT_DIR")"
BIN_DIR="$SUPERMAC_ROOT/bin"
LIB_DIR="$SUPERMAC_ROOT/lib"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m'

# Test counters
TESTS_PASSED=0
TESTS_FAILED=0
TESTS_TOTAL=0

# Output functions
test_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
test_success() { echo -e "${GREEN}[PASS]${NC} $1"; ((TESTS_PASSED++)); }
test_fail() { echo -e "${RED}[FAIL]${NC} $1"; ((TESTS_FAILED++)); }
test_header() { echo -e "${BOLD}${BLUE}=== $1 ===${NC}"; }

# Increment test counter
count_test() { ((TESTS_TOTAL++)); }

# =============================================================================
# Architecture Tests
# =============================================================================

test_file_structure() {
    test_header "File Structure Tests"
    
    # Test main directories exist
    count_test
    if [[ -d "$BIN_DIR" ]]; then
        test_success "bin/ directory exists"
    else
        test_fail "bin/ directory missing"
    fi
    
    count_test
    if [[ -d "$LIB_DIR" ]]; then
        test_success "lib/ directory exists"
    else
        test_fail "lib/ directory missing"
    fi
    
    # Test main files exist
    count_test
    if [[ -f "$BIN_DIR/mac" ]]; then
        test_success "Main mac script exists"
    else
        test_fail "Main mac script missing"
    fi
    
    count_test
    if [[ -x "$BIN_DIR/mac" ]]; then
        test_success "Main mac script is executable"
    else
        test_fail "Main mac script not executable"
    fi
    
    # Test core modules exist
    local modules=("utils" "finder" "wifi" "network" "system" "dev" "display")
    for module in "${modules[@]}"; do
        count_test
        if [[ -f "$LIB_DIR/$module.sh" ]]; then
            test_success "Module $module.sh exists"
        else
            test_fail "Module $module.sh missing"
        fi
    done
}

test_syntax() {
    test_header "Syntax Tests"
    
    # Test main script syntax
    count_test
    if bash -n "$BIN_DIR/mac" 2>/dev/null; then
        test_success "Main script syntax is valid"
    else
        test_fail "Main script has syntax errors"
    fi
    
    # Test module syntax
    for module_file in "$LIB_DIR"/*.sh; do
        if [[ -f "$module_file" ]]; then
            count_test
            local module_name
            module_name=$(basename "$module_file" .sh)
            if bash -n "$module_file" 2>/dev/null; then
                test_success "Module $module_name syntax is valid"
            else
                test_fail "Module $module_name has syntax errors"
            fi
        fi
    done
}

test_module_loading() {
    test_header "Module Loading Tests"
    
    # Source utils first
    count_test
    if source "$LIB_DIR/utils.sh" 2>/dev/null; then
        test_success "Utils module loads successfully"
    else
        test_fail "Utils module failed to load"
        return 1
    fi
    
    # Test each module can be sourced
    local modules=("finder" "wifi" "network" "system" "dev" "display")
    for module in "${modules[@]}"; do
        count_test
        if source "$LIB_DIR/$module.sh" 2>/dev/null; then
            test_success "Module $module loads successfully"
        else
            test_fail "Module $module failed to load"
        fi
    done
}

test_function_existence() {
    test_header "Function Existence Tests"
    
    # Source all modules
    source "$LIB_DIR/utils.sh" 2>/dev/null || return 1
    
    # Test that required functions exist in utils
    local utils_functions=("print_success" "print_error" "print_info" "print_header" "check_macos")
    for func in "${utils_functions[@]}"; do
        count_test
        if declare -f "$func" >/dev/null; then
            test_success "Utils function $func exists"
        else
            test_fail "Utils function $func missing"
        fi
    done
    
    # Test module dispatcher functions
    local modules=("finder" "wifi" "network" "system" "dev" "display")
    for module in "${modules[@]}"; do
        source "$LIB_DIR/$module.sh" 2>/dev/null || continue
        
        count_test
        if declare -f "${module}_dispatch" >/dev/null; then
            test_success "Module $module has dispatcher function"
        else
            test_fail "Module $module missing dispatcher function"
        fi
        
        count_test
        if declare -f "${module}_help" >/dev/null; then
            test_success "Module $module has help function"
        else
            test_fail "Module $module missing help function"
        fi
    done
}

# =============================================================================
# Functional Tests
# =============================================================================

test_help_system() {
    test_header "Help System Tests"
    
    # Test main help
    count_test
    if "$BIN_DIR/mac" help >/dev/null 2>&1; then
        test_success "Main help command works"
    else
        test_fail "Main help command failed"
    fi
    
    # Test category help
    local modules=("finder" "wifi" "network" "system" "dev" "display")
    for module in "${modules[@]}"; do
        count_test
        if "$BIN_DIR/mac" help "$module" >/dev/null 2>&1; then
            test_success "Help for $module works"
        else
            test_fail "Help for $module failed"
        fi
    done
}

test_version_info() {
    test_header "Version Info Tests"
    
    count_test
    if "$BIN_DIR/mac" version >/dev/null 2>&1; then
        test_success "Version command works"
    else
        test_fail "Version command failed"
    fi
    
    count_test
    if "$BIN_DIR/mac" --version >/dev/null 2>&1; then
        test_success "--version flag works"
    else
        test_fail "--version flag failed"
    fi
}

test_error_handling() {
    test_header "Error Handling Tests"
    
    # Test unknown category
    count_test
    if ! "$BIN_DIR/mac" nonexistent action >/dev/null 2>&1; then
        test_success "Unknown category properly rejected"
    else
        test_fail "Unknown category not properly handled"
    fi
    
    # Test missing action
    count_test
    if ! "$BIN_DIR/mac" finder >/dev/null 2>&1; then
        test_success "Missing action properly rejected"
    else
        test_fail "Missing action not properly handled"
    fi
    
    # Test unknown action
    count_test
    if ! "$BIN_DIR/mac" finder nonexistent >/dev/null 2>&1; then
        test_success "Unknown action properly rejected"
    else
        test_fail "Unknown action not properly handled"
    fi
}

test_safe_commands() {
    test_header "Safe Command Tests"
    
    # Test commands that should work without side effects
    
    # Finder status
    count_test
    if "$BIN_DIR/mac" finder status >/dev/null 2>&1; then
        test_success "Finder status command works"
    else
        test_fail "Finder status command failed"
    fi
    
    # Network IP
    count_test
    if "$BIN_DIR/mac" network ip >/dev/null 2>&1; then
        test_success "Network IP command works"
    else
        test_fail "Network IP command failed"
    fi
    
    # System info
    count_test
    if "$BIN_DIR/mac" system info >/dev/null 2>&1; then
        test_success "System info command works"
    else
        test_fail "System info command failed"
    fi
    
    # WiFi status
    count_test
    if "$BIN_DIR/mac" wifi status >/dev/null 2>&1; then
        test_success "WiFi status command works"
    else
        test_fail "WiFi status command failed"
    fi
    
    # Dev processes
    count_test
    if "$BIN_DIR/mac" dev processes >/dev/null 2>&1; then
        test_success "Dev processes command works"
    else
        test_fail "Dev processes command failed"
    fi
    
    # Display status
    count_test
    if "$BIN_DIR/mac" display status >/dev/null 2>&1; then
        test_success "Display status command works"
    else
        test_fail "Display status command failed"
    fi
}

# =============================================================================
# Configuration Tests
# =============================================================================

test_config_system() {
    test_header "Configuration System Tests"
    
    local config_file="$SUPERMAC_ROOT/config/config.json"
    
    count_test
    if [[ -f "$config_file" ]]; then
        test_success "Configuration file exists"
    else
        test_fail "Configuration file missing"
    fi
    
    # Test JSON validity if jq is available
    if command -v jq >/dev/null 2>&1; then
        count_test
        if jq '.' "$config_file" >/dev/null 2>&1; then
            test_success "Configuration file is valid JSON"
        else
            test_fail "Configuration file has invalid JSON"
        fi
    fi
}

# =============================================================================
# Performance Tests
# =============================================================================

test_performance() {
    test_header "Performance Tests"
    
    # Test startup time
    count_test
    local start_time end_time duration
    start_time=$(date +%s.%N)
    "$BIN_DIR/mac" help >/dev/null 2>&1
    end_time=$(date +%s.%N)
    duration=$(echo "$end_time - $start_time" | bc -l 2>/dev/null || echo "1")
    
    if (( $(echo "$duration < 2.0" | bc -l 2>/dev/null || echo "0") )); then
        test_success "Startup time acceptable (${duration}s)"
    else
        test_fail "Startup time too slow (${duration}s)"
    fi
    
    # Test multiple command execution
    count_test
    start_time=$(date +%s.%N)
    for i in {1..5}; do
        "$BIN_DIR/mac" system info >/dev/null 2>&1
    done
    end_time=$(date +%s.%N)
    duration=$(echo "$end_time - $start_time" | bc -l 2>/dev/null || echo "10")
    
    if (( $(echo "$duration < 10.0" | bc -l 2>/dev/null || echo "0") )); then
        test_success "Multiple commands execute quickly (${duration}s for 5 commands)"
    else
        test_fail "Multiple commands too slow (${duration}s for 5 commands)"
    fi
}

# =============================================================================
# Module-Specific Tests
# =============================================================================

test_specific_module() {
    local module="$1"
    
    test_header "Testing Module: $module"
    
    case "$module" in
        "finder")
            # Test finder module specific functionality
            count_test
            if "$BIN_DIR/mac" finder status >/dev/null 2>&1; then
                test_success "Finder status works"
            else
                test_fail "Finder status failed"
            fi
            ;;
        "network")
            # Test network module
            count_test
            if "$BIN_DIR/mac" network ip >/dev/null 2>&1; then
                test_success "Network IP works"
            else
                test_fail "Network IP failed"
            fi
            ;;
        "system")
            # Test system module
            count_test
            if "$BIN_DIR/mac" system info >/dev/null 2>&1; then
                test_success "System info works"
            else
                test_fail "System info failed"
            fi
            ;;
        "dev")
            # Test dev module
            count_test
            if "$BIN_DIR/mac" dev processes >/dev/null 2>&1; then
                test_success "Dev processes works"
            else
                test_fail "Dev processes failed"
            fi
            ;;
        *)
            test_info "No specific tests for module: $module"
            ;;
    esac
}

# =============================================================================
# Test Runner
# =============================================================================

run_all_tests() {
    test_header "SuperMac Test Suite"
    echo ""
    
    test_file_structure
    echo ""
    
    test_syntax
    echo ""
    
    test_module_loading
    echo ""
    
    test_function_existence
    echo ""
    
    test_help_system
    echo ""
    
    test_version_info
    echo ""
    
    test_error_handling
    echo ""
    
    test_safe_commands
    echo ""
    
    test_config_system
    echo ""
    
    test_performance
    echo ""
}

show_summary() {
    echo ""
    test_header "Test Summary"
    echo ""
    echo "Total tests: $TESTS_TOTAL"
    echo -e "Passed: ${GREEN}$TESTS_PASSED${NC}"
    echo -e "Failed: ${RED}$TESTS_FAILED${NC}"
    echo ""
    
    if [[ $TESTS_FAILED -eq 0 ]]; then
        echo -e "${GREEN}${BOLD}✅ All tests passed!${NC}"
        return 0
    else
        echo -e "${RED}${BOLD}❌ Some tests failed${NC}"
        return 1
    fi
}

# =============================================================================
# Main
# =============================================================================

main() {
    local target_module="$1"
    
    # Check if SuperMac structure exists
    if [[ ! -f "$BIN_DIR/mac" ]]; then
        echo -e "${RED}Error: SuperMac not found at $BIN_DIR/mac${NC}"
        echo "Please run tests from the SuperMac directory"
        exit 1
    fi
    
    if [[ -n "$target_module" ]]; then
        # Test specific module
        test_file_structure
        test_syntax
        test_module_loading
        test_specific_module "$target_module"
    else
        # Run all tests
        run_all_tests
    fi
    
    show_summary
}

# Run tests
main "$@"
