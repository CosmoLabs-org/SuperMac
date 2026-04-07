# Codebase Patterns — SuperMac v2.1.0

## Module Contract Pattern (consistent across all 9 modules)

Every module MUST implement:
```bash
{module}_dispatch()   # Routes action to handler via case statement
{module}_help()       # Formatted help output with box drawing
{module}_search()     # Keyword search via hardcoded substring matching
```

Optional per-action functions:
```bash
{module}_{action}()   # Individual action handler
{module}_get_state()  # State query helper
```

This pattern is followed consistently across all modules. (see agent-1-code-quality.md, agent-8-ux-cli.md)

## Error Handling Patterns

### Pattern 1: Three-part error message (GOOD)
```bash
print_error "Unknown action: $action"
echo ""
print_info "Available actions: on, off, toggle, status, ..."
print_info "Use 'mac help $category' for detailed information"
```
Used in all 9 module dispatchers. Consistent and helpful. (see agent-8-ux-cli.md)

### Pattern 2: Silent error suppression (PROBLEMATIC)
```bash
command 2>/dev/null || true
```
Used extensively in system.sh cleanup (lines 128-178). Hides permission errors, disk corruption, and real failures. (see agent-1-code-quality.md, agent-3-security.md)

### Pattern 3: Confirmation before destruction (GOOD)
```bash
if ! confirm "Continue with system cleanup?" "N"; then
    print_info "Cleanup cancelled"
    return 0
fi
```
Used for cleanup, network reset, dock reset, screenshot reset. Safe defaults. (see agent-3-security.md)

## State Management

### Pattern: Idempotent state detection
Most state-changing commands detect current state:
```bash
wifi_get_power_state()  # Checks before toggling
display_get_dark_mode() # Checks before switching
audio_get_mute_state()  # Checks before toggling
```
Reports "already enabled/disabled" rather than re-applying. Good UX. (see agent-8-ux-cli.md)

### Anti-pattern: No config persistence
Config changes made via defaults write are persistent, but the tool's own config.json is never read or written. Settings like volume_step and screenshot_location are defined but ignored. (see agent-2-core-logic.md)

## Naming Conventions

| Pattern | Convention | Example |
|---------|-----------|---------|
| Functions | snake_case | system_get_battery() |
| Constants | UPPER_SNAKE | SUPERMAC_VERSION |
| Local vars | snake_case | local port_number |
| Actions | kebab-case | kill-port, dark-mode |
| Categories | lowercase | finder, wifi, network |
| Globals | UPPER_SNAKE | BOLD, RED, NC |

Inconsistencies found:
- Box widths vary: 50-80 across modules (see agent-8-ux-cli.md)
- Boolean synonyms differ per module: dock accepts 8, wifi accepts 2 (see agent-8-ux-cli.md)
- Some aliases undocumented: wifi "enable", dock "l" for "left" (see agent-8-ux-cli.md)

## File Organization Pattern

### Canonical (correct):
```
bin/mac          - dispatcher
lib/*.sh         - modules
config/config.json - settings
tests/test.sh    - tests
Makefile         - build automation
```

### Actual (problematic):
```
ROOT/*.sh        - exact copies of lib/* (5,154 duplicate lines)
ROOT/mac         - exact copy of bin/mac
ROOT/config.json - exact copy of config/config.json
```
Every root copy confirmed identical via diff. (see agent-1-code-quality.md, agent-4-architecture.md, agent-6-distribution.md)

## Shellcheck Patterns

Known issues recurring across modules:
- SC2086: Unquoted $port in lsof calls (dev.sh x5)
- SC2072: String comparison for decimals (dev.sh x2)
- SC2155: local+assign in one line (utils.sh, bin/mac x5)
- SC2034: Unused exported constants (utils.sh)

No .shellcheckrc exists. Make lint runs shellcheck but fails on first warning. (see agent-1-code-quality.md)

## Agent Cross-References

| Pattern | Source Agent | Evidence |
|---------|-------------|----------|
| Module contract | agent-4-architecture.md | All 9 modules follow |
| Silent error suppression | agent-1-code-quality.md, agent-3-security.md | system.sh:128-178 |
| Three-part error messages | agent-8-ux-cli.md | All dispatchers |
| Idempotent state checks | agent-8-ux-cli.md | wifi, display, audio |
| Dead config | agent-2-core-logic.md | get_config() zero calls |
| Boolean synonym inconsistency | agent-8-ux-cli.md | dock vs wifi vs screenshot |
| Shellcheck issues | agent-1-code-quality.md | dev.sh, utils.sh, bin/mac |
