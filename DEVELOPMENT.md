# 🛠️ SuperMac Development Guide

## 📋 Overview

SuperMac is a modular, professional command-line tool for macOS built with a clean architecture that makes it easy to add new commands and maintain existing functionality.

## 🏗️ Architecture

### File Structure
```
SuperMac/
├── bin/
│   ├── mac                 # Main dispatcher (entry point)
│   └── install.sh          # Installation script
├── lib/                    # Modular command libraries
│   ├── utils.sh           # Shared utilities & formatting
│   ├── finder.sh          # Finder commands
│   ├── display.sh         # Display commands
│   ├── network.sh         # Network commands
│   ├── system.sh          # System commands
│   ├── dev.sh             # Developer commands
│   ├── dock.sh            # Dock management
│   └── audio.sh           # Audio control
├── config/
│   └── config.json        # User configuration
├── docs/
│   ├── README.md          # Main documentation
│   └── DEVELOPMENT.md     # This file
└── tests/
    └── test.sh            # Test suite
```

### Design Principles

1. **Modular**: Each category is a separate module
2. **Consistent**: All modules follow the same patterns
3. **Discoverable**: Beautiful help system with search
4. **Safe**: Input validation and confirmation for destructive operations
5. **Fast**: Minimal overhead, optimized for quick execution
6. **Professional**: Enterprise-ready with proper error handling

## 🧩 Module Architecture

### Module Template

Every module follows this structure:

```bash
#!/bin/bash

# =============================================================================
# SuperMac - [ModuleName] Module
# =============================================================================
# Brief description of module functionality
# 
# Commands:
#   command1 - Description
#   command2 - Description
#
# Built by CosmoLabs - https://cosmolabs.org
# =============================================================================

# =============================================================================
# Module Functions
# =============================================================================

modulename_function1() {
    # Implementation
}

modulename_function2() {
    # Implementation
}

# =============================================================================
# Module Help System
# =============================================================================

modulename_help() {
    print_category_header "modulename" "🔥" 70
    
    printf "${PURPLE}│${NC}  %-20s %-40s ${PURPLE}│${NC}\n" "command1" "Description"
    printf "${PURPLE}│${NC}  %-20s %-40s ${PURPLE}│${NC}\n" "command2" "Description"
    
    print_category_footer 70
    echo ""
    
    print_header "Examples:"
    echo "  mac modulename command1             # Example usage"
    echo ""
}

# Search function for this module
modulename_search() {
    local search_term="$1"
    local results=""
    
    if [[ "keyword" == *"$search_term"* ]]; then
        results+="  mac modulename command1          Description\n"
    fi
    
    if [[ -n "$results" ]]; then
        echo -e "$results"
    fi
}

# =============================================================================
# Module Dispatcher
# =============================================================================

modulename_dispatch() {
    local action="$1"
    shift  # Remove action from arguments
    
    case "$action" in
        "command1")
            modulename_function1 "$@"
            ;;
        "command2")
            modulename_function2 "$@"
            ;;
        "help"|"-h"|"--help")
            modulename_help
            ;;
        *)
            print_error "Unknown modulename action: $action"
            echo ""
            print_info "Available actions: command1, command2"
            print_info "Use 'mac help modulename' for detailed information"
            return 1
            ;;
    esac
}

print_debug "ModuleName module loaded successfully"
```

### Required Functions

Every module MUST implement these functions:

1. **`modulename_dispatch()`** - Routes actions to appropriate functions
2. **`modulename_help()`** - Provides detailed help for the module
3. **`modulename_search()`** - Enables search functionality

### Optional Functions

Modules can implement these for enhanced functionality:

- **`modulename_status()`** - Show current status/settings
- **`modulename_reset()`** - Reset module settings to defaults
- **`modulename_validate()`** - Validate module requirements

## 🎨 Styling Guidelines

### Output Functions

Use the standard output functions from `utils.sh`:

```bash
print_success "Operation completed!"       # Green checkmark
print_error "Something went wrong"         # Red X
print_info "Information message"           # Blue info icon
print_warning "Proceed with caution"       # Yellow warning
print_header "Section Header"              # Bold blue text
print_debug "Debug information"            # Gray text (only in debug mode)
```

### Visual Formatting

Use the box drawing functions for beautiful output:

```bash
print_category_header "category" "🔥" 70   # Category header with emoji
print_category_footer 70                  # Matching footer
print_banner "SuperMac v2.1.0"           # Centered banner
print_box "Important message"             # Simple box around text
```

### Color Usage

Colors are available as constants:
- `$RED`, `$GREEN`, `$YELLOW`, `$BLUE`, `$PURPLE`, `$CYAN`
- `$BOLD`, `$DIM`, `$UNDERLINE`
- `$NC` (No Color) - always end colored text with this

## 🔧 Adding a New Module

### Step 1: Create Module File

Create `lib/newmodule.sh` following the module template above.

### Step 2: Update Main Dispatcher

Add your module to `bin/mac`:

```bash
# In the CATEGORIES array
declare -A CATEGORIES=(
    # ... existing categories ...
    ["newmodule"]="🔥:Description of new module"
)

# In the route_command function, the module will be automatically loaded
```

### Step 3: Update Configuration

Add your module to `config/config.json`:

```json
{
  "categories": {
    "newmodule": {
      "enabled": true,
      "description": "Description of new module"
    }
  }
}
```

### Step 4: Add Tests

Add tests for your module in `tests/test.sh`:

```bash
test_newmodule() {
    test_header "Testing New Module"
    
    count_test
    if "$BIN_DIR/mac" newmodule help >/dev/null 2>&1; then
        test_success "New module help works"
    else
        test_fail "New module help failed"
    fi
}
```

### Step 5: Update Documentation

Add your module to the main README.md and any relevant documentation.

## 🧪 Testing

### Running Tests

```bash
# Run all tests
bash tests/test.sh

# Test specific module
bash tests/test.sh newmodule

# Test with debug output
bash tests/test.sh --debug
```

### Test Types

1. **Syntax Tests** - Verify bash syntax is valid
2. **Module Loading** - Ensure modules can be sourced
3. **Function Existence** - Check required functions exist
4. **Functional Tests** - Test actual command execution
5. **Performance Tests** - Verify startup time and speed

### Writing Good Tests

```bash
test_new_feature() {
    test_header "Feature Tests"
    
    # Test success case
    count_test
    if expected_command_succeeds; then
        test_success "Feature works correctly"
    else
        test_fail "Feature failed"
    fi
    
    # Test error case
    count_test
    if ! expected_command_fails; then
        test_success "Error handling works"
    else
        test_fail "Error handling broken"
    fi
}
```

## 📝 Code Standards

### Naming Conventions

- **Functions**: `modulename_functionname()` (underscore separated)
- **Variables**: `local_variable_name` (lowercase with underscores)
- **Constants**: `CONSTANT_NAME` (uppercase with underscores)
- **Files**: `modulename.sh` (lowercase)

### Error Handling

Always include proper error handling:

```bash
function_with_validation() {
    local input="$1"
    
    # Validate input
    if [[ -z "$input" ]]; then
        print_error "Input required"
        print_info "Usage: mac module function <input>"
        return 1
    fi
    
    # Validate input format
    if ! is_valid_format "$input"; then
        print_error "Invalid input format: $input"
        return 1
    fi
    
    # Execute with error checking
    if ! some_command "$input"; then
        print_error "Command failed"
        return 1
    fi
    
    print_success "Operation completed successfully"
}
```

### Input Validation

Use utility functions for validation:

```bash
# Number validation
if ! is_number "$port"; then
    print_error "Invalid port number: $port"
    return 1
fi

# Range validation
if ! is_in_range "$volume" 0 100; then
    print_error "Volume must be between 0 and 100"
    return 1
fi

# File existence
if ! file_exists "$config_file"; then
    print_error "Configuration file not found: $config_file"
    return 1
fi
```

### User Interaction

For destructive operations, always confirm:

```bash
dangerous_operation() {
    print_warning "This will delete important data"
    
    if ! confirm "Are you sure you want to continue?" "n"; then
        print_info "Operation cancelled"
        return 0
    fi
    
    # Proceed with operation
}
```

## 🚀 Performance Guidelines

### Startup Time

- Keep module loading fast
- Avoid expensive operations in global scope
- Use lazy loading where possible

### Command Execution

- Optimize for common use cases
- Cache expensive computations
- Provide progress feedback for long operations

### Memory Usage

- Use local variables where possible
- Clean up temporary files
- Avoid large arrays in global scope

## 🔒 Security Considerations

### Input Sanitization

Always sanitize user input:

```bash
safe_filename() {
    local filename="$1"
    
    # Remove dangerous characters
    filename=$(echo "$filename" | tr -d ';<>|&`$(){}[]')
    
    # Prevent path traversal
    filename=$(basename "$filename")
    
    echo "$filename"
}
```

### Privilege Escalation

- Only use `sudo` when absolutely necessary
- Always explain why privileges are needed
- Provide alternatives when possible

### File Operations

- Validate file paths
- Check file permissions
- Create backups for important files

## 📚 Best Practices

### Documentation

1. **Inline Comments**: Explain complex logic
2. **Function Headers**: Describe purpose and parameters
3. **Usage Examples**: Show how to use functions
4. **Error Messages**: Be helpful and actionable

### User Experience

1. **Consistent Interface**: Follow established patterns
2. **Clear Feedback**: Always tell users what happened
3. **Helpful Errors**: Suggest solutions, not just problems
4. **Progressive Disclosure**: Start simple, allow complexity

### Maintainability

1. **Small Functions**: Keep functions focused and small
2. **Reusable Code**: Use utility functions
3. **Clear Structure**: Organize code logically
4. **Version Control**: Use meaningful commit messages

## 🎯 Common Patterns

### Status Commands

```bash
module_status() {
    print_header "📊 Module Status"
    echo ""
    
    local setting1 setting2
    setting1=$(get_setting1)
    setting2=$(get_setting2)
    
    echo "  Setting 1: $(print_command "$setting1")"
    echo "  Setting 2: $(print_command "$setting2")"
    
    if [[ "$setting1" == "optimal" ]]; then
        print_success "Configuration looks good"
    else
        print_warning "Consider adjusting settings"
    fi
}
```

### Configuration Commands

```bash
module_set_option() {
    local option="$1"
    local value="$2"
    
    case "$option" in
        "setting1")
            if validate_setting1 "$value"; then
                apply_setting1 "$value"
                print_success "Setting 1 updated to: $value"
            else
                print_error "Invalid value for setting1: $value"
                return 1
            fi
            ;;
        *)
            print_error "Unknown option: $option"
            return 1
            ;;
    esac
}
```

### List Commands

```bash
module_list_items() {
    print_header "📋 Available Items"
    echo ""
    
    get_items | while read -r item; do
        if is_active_item "$item"; then
            echo "  $(print_command "$item") $(print_dim "(active)")"
        else
            echo "  $(print_dim "$item")"
        fi
    done
}
```

## 🐛 Debugging

### Debug Mode

Enable debug output:

```bash
# Enable debug globally
export SUPERMAC_DEBUG=1

# Or per command
mac --debug system info
```

### Debug Functions

Use debug functions liberally:

```bash
complex_function() {
    print_debug "Starting complex operation"
    print_debug "Parameter 1: $1"
    print_debug "Parameter 2: $2"
    
    # ... operation ...
    
    print_debug "Operation completed successfully"
}
```

### Common Issues

1. **Module not loading**: Check syntax with `bash -n module.sh`
2. **Function not found**: Ensure function is exported
3. **Path issues**: Use absolute paths or validate relative paths
4. **Permission errors**: Check file permissions and ownership

## 📦 Deployment

### Pre-release Checklist

- [ ] All tests pass
- [ ] Documentation updated
- [ ] Version numbers incremented
- [ ] Install script tested
- [ ] Backward compatibility verified

### Release Process

1. Update version in all relevant files
2. Run comprehensive tests
3. Update documentation
4. Create release notes
5. Tag release in git
6. Update installation script

## 🤝 Contributing

### Pull Request Process

1. Fork the repository
2. Create feature branch: `git checkout -b feature/awesome-feature`
3. Follow coding standards
4. Add tests for new functionality
5. Update documentation
6. Submit pull request

### Code Review

All code changes require review focusing on:
- Functionality and correctness
- Code style and consistency
- Test coverage
- Documentation quality
- Performance impact

---

**Happy coding! 🚀**

*Built with ❤️ by CosmoLabs*
