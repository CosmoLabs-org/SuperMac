# Architecture Map — SuperMac v2.1.0

## System Layers

```
User Input Layer
    |
    v
bin/mac (349 lines) — Main Dispatcher
    |  Routes: help, version, search, debug, global shortcuts, module dispatch
    |  Strict mode: set -euo pipefail
    |
    v
lib/utils.sh (589 lines) — Shared Utilities
    |  Colors, output formatting, validation, module loading, config reader
    |  Loaded eagerly at startup (always sourced)
    |
    v
lib/*.sh (9 modules, 4,141 lines total) — Feature Modules
    |  Loaded lazily on-demand via load_module()
    |  Each follows: _dispatch() / _help() / _search() contract
    |
    v
macOS APIs — System Commands
    osascript, defaults, networksetup, system_profiler, lsof, etc.
```

## Module Boundaries

| Module | Lines | Commands | Dependencies | macOS APIs Used |
|--------|-------|----------|--------------|-----------------|
| finder | 313 | 6 | utils only | killall Finder, defaults write, open |
| wifi | 507 | 9 | utils only | networksetup, airport (broken path) |
| network | 484 | 9 | utils + optional wifi | ipconfig, ifconfig, ping, scutil, dscacheutil |
| system | 541 | 9 | utils only | system_profiler, vm_stat, df, pmset |
| dev | 610 | 13 | utils only | lsof, ps, python3, jq, openssl |
| display | 428 | 8 | utils only | osascript (brightness, dark mode) |
| dock | 573 | 8 | utils + optional dockutil | defaults write, killall Dock |
| audio | 519 | 11 | utils + optional SwitchAudioSource | osascript (volume) |
| screenshot | 590 | 10 | utils only | defaults write/read, screencapture |

## Dependency Graph

```
utils.sh (always loaded first)
    |
    +-- finder.sh (independent)
    +-- wifi.sh (independent)
    +-- network.sh (optional: probes wifi_get_current_network via declare -f)
    +-- system.sh (independent)
    +-- dev.sh (independent)
    +-- display.sh (independent)
    +-- dock.sh (independent)
    +-- audio.sh (independent)
    +-- screenshot.sh (independent)

bin/mac (dispatcher)
    |
    +-- sources utils.sh at startup
    +-- loads modules on-demand via load_module()
    +-- resolves GLOBAL_SHORTCUTS to category:action pairs

No circular dependencies exist.
Only one cross-module reference: network -> wifi (optional, via declare -f probe).
```

## Data Flow (see agent-2-core-logic.md)

### Request lifecycle: `mac system info`
1. bin/mac:347 main() receives ["system", "info"]
2. bin/mac:286 validate_environment() checks macOS + lib dir
3. bin/mac:339 route_command("system", "info")
4. bin/mac:247 validates category in CATEGORIES dict
5. bin/mac:256 load_module("system") sources lib/system.sh
6. bin/mac:270 system_dispatch("info") called
7. lib/system.sh:27 system_info() executes
8. Calls: sw_vers, system_profiler (4x, slow), uname, uptime, df, vm_stat

### Global shortcut lifecycle: `mac ip`
1. main() receives ["ip"]
2. route_command("ip") checks GLOBAL_SHORTCUTS["ip"] = "network:ip"
3. Parses to category=network, action=ip
4. load_module("network") + network_dispatch("ip")
5. network_get_local_ip() iterates en0-en3 with ipconfig getifaddr

## Critical Architecture Issues (see agent-4-architecture.md)

1. **Dead config system**: config.json defines settings but get_config() is never called by any module
2. **No auto-discovery**: CATEGORIES hardcoded in bin/mac instead of scanned from lib/*.sh
3. **God-file utils.sh**: 589 lines mixing colors, output, validation, config, loading, debug
4. **Root duplicates**: 14 identical files at project root undermine the canonical structure

## Agent Cross-References

| Finding | Source Agent | File:Line |
|---------|-------------|-----------|
| Config never consumed | agent-2-core-logic.md | utils.sh:429 (get_config unused) |
| No module re-source guard | agent-2-core-logic.md | utils.sh:510-521 |
| Flat dep tree is good | agent-4-architecture.md | All modules depend only on utils.sh |
| utils.sh is a god-file | agent-4-architecture.md | utils.sh spans 8 responsibilities |
| Root duplicates confuse contributors | agent-1-code-quality.md | Root level vs lib/ |
| dev.sh has split personality | agent-4-architecture.md | dev.sh:30-315 vs 321-461 |
