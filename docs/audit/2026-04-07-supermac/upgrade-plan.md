# Upgrade Plan — SuperMac v2.1.0

Phased action plan based on findings from 8 audit agents.

## Phase 0: Critical Fixes (1-2 days)

These are must-fix items that represent broken functionality or serious risk.

### 0.1 Delete root-level duplicate files
- **What**: Remove 14 identical files at project root (mac, *.sh, config.json)
- **Why**: 5,154 lines of dead code creates maintenance drift and confusion
- **Files**: All root-level .sh files + config.json + mac
- **Source**: agent-1-code-quality.md, agent-4-architecture.md, agent-6-distribution.md

### 0.2 Fix rm -rf /tmp/* in system cleanup
- **What**: Replace `rm -rf /tmp/*` with targeted cleanup
- **Fix**: `find "${TMPDIR:-/tmp}" -type f -user "$USER" -atime +7 -delete`
- **File**: lib/system.sh:173
- **Source**: agent-3-security.md, agent-2-core-logic.md

### 0.3 Fix airport binary path
- **What**: Use full path instead of PATH lookup
- **Fix**: `AIRPORT="/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport"`
- **File**: lib/wifi.sh:174, 237, 288
- **Source**: agent-2-core-logic.md

### 0.4 Fix Apple Silicon memory stats
- **What**: Read page size from sysctl instead of hardcoding 4096
- **Fix**: `page_size=$(sysctl -n hw.pagesize)`
- **File**: lib/system.sh memory functions
- **Source**: agent-2-core-logic.md

### 0.5 Add missing modules to installer
- **What**: Add wifi.sh, dock.sh, audio.sh, screenshot.sh to download list
- **File**: bin/install.sh:154-162
- **Source**: agent-6-distribution.md, agent-7-documentation.md

### 0.6 Create LICENSE file
- **What**: Add MIT LICENSE with copyright to CosmoLabs
- **Source**: agent-6-distribution.md

### 0.7 Fix dock.sh string concatenation bug
- **What**: `"Dock auto-hide $action"d!` → `"Dock auto-hide ${action}d!"`
- **File**: lib/dock.sh:131
- **Source**: agent-1-code-quality.md

### 0.8 Fix dev.sh decimal comparison bug
- **What**: `[[ "$cpu" > 5.0 ]]` uses string comparison, broken for values >= 10
- **Fix**: Use `bc -l` like dev_memory_hogs already does
- **File**: lib/dev.sh:281-283
- **Source**: agent-1-code-quality.md, agent-2-core-logic.md

## Phase 1: Foundation (3-5 days)

### 1.1 Make config.json functional or remove it
- **Decision needed**: Connect get_config() to modules, or delete the decorative config
- **If keeping**: Each module reads defaults from config on first call
- **If removing**: Delete config.json, get_config(), and all config references
- **Source**: agent-2-core-logic.md, agent-4-architecture.md

### 1.2 Add shellcheck to CI
- **What**: Add .shellcheckrc, fix all SC warnings, add to Makefile check target
- **Source**: agent-1-code-quality.md

### 1.3 Expand test suite
- **What**: Add dock, audio, screenshot to module lists. Add output validation. Add edge case tests.
- **Target**: 50% command coverage (up from 7.8%)
- **Source**: agent-5-testing.md

### 1.4 Single-source version string
- **What**: Define version in one place, read from all others
- **Source**: agent-1-code-quality.md, agent-6-distribution.md

### 1.5 Fix .gitignore
- **What**: Add .DS_Store, dist/, *.tar.gz, *.log, *.tmp
- **Source**: agent-6-distribution.md

### 1.6 Add NO_COLOR support
- **What**: Check `${NO_COLOR:-}` and add `--no-color` flag
- **Source**: agent-8-ux-cli.md

## Phase 2: Quality (1 week)

### 2.1 Split utils.sh into focused files
- **What**: Extract colors.sh, output.sh, validation.sh, interaction.sh, loader.sh
- **Source**: agent-4-architecture.md

### 2.2 Add module auto-discovery
- **What**: Replace hardcoded CATEGORIES with runtime scan of lib/*.sh
- **Source**: agent-4-architecture.md

### 2.3 Add --yes flag for non-interactive use
- **What**: Global flag to bypass confirm() prompts
- **Source**: agent-8-ux-cli.md

### 2.4 Consolidate boolean parsing
- **What**: Shared parse_boolean() in utils.sh
- **Source**: agent-8-ux-cli.md

### 2.5 Add install integrity verification
- **What**: SHA256 checksums for all downloaded files
- **Source**: agent-3-security.md, agent-6-distribution.md

### 2.6 Improve uninstall
- **What**: Confirmation prompt, config backup, PATH cleanup, documented in README
- **Source**: agent-6-distribution.md

### 2.7 Fix documentation accuracy
- **What**: Remove non-existent commands from README, add screenshot module, fix architecture diagram
- **Source**: agent-7-documentation.md

## Phase 3: Growth (1-2 weeks)

### 3.1 Homebrew formula
- **Prerequisites**: LICENSE file, clean repo, configurable install paths
- **Source**: agent-6-distribution.md

### 3.2 GitHub Actions CI/CD
- **What**: Lint, syntax check, test on macOS runner
- **Source**: agent-6-distribution.md, agent-5-testing.md

### 3.3 Auto-update mechanism
- **What**: Check GitHub releases API, compare version, notify user
- **Source**: agent-6-distribution.md (auto_update_check config field)

### 3.4 Plugin/module system
- **What**: Drop-in modules with auto-registration
- **Source**: agent-4-architecture.md

### 3.5 Performance optimization
- **What**: Cache system_profiler calls, batch API calls
- **Source**: agent-2-core-logic.md

### 3.6 Add IPv6 support
- **Source**: agent-2-core-logic.md

### 3.7 Accessibility improvements
- **What**: Text fallbacks for emojis, Unicode detection, progress bar ASCII fallback
- **Source**: agent-8-ux-cli.md
