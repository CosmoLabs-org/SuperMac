# SuperMac v2.1.0 - Comprehensive Audit Report

Date: 2026-04-07 | Mode: Fresh | Agents: 8 | Duration: ~10min

## Scorecard

| Agent | Score | Grade | Priority |
|-------|-------|-------|----------|
| Code Quality | 52/100 | D | High |
| Core Logic | 62/100 | C | Medium |
| Security | 38/100 | D- | CRITICAL |
| Architecture | 62/100 | C | Medium |
| Testing | 28/100 | F | CRITICAL |
| Distribution | 38/100 | D- | CRITICAL |
| Documentation | 52/100 | D | High |
| UX & CLI | 72/100 | C+ | Maintain |
| **Overall** | **50.5/100** | **D** | **High** |

Grade scale: A (86-100), B (61-85), C (41-60), D (21-40), F (0-20)

## Top 5 Strengths

1. **Polished UX/CLI design** (72/100) — Consistent module pattern, beautiful help system with box drawing, helpful error messages, smart global shortcuts. (see agent-8-ux-cli.md)
2. **Clean module architecture** (coupling: 75/100) — Flat dependency tree with no circular deps, lazy module loading, consistent _dispatch/_help/_search contract. (see agent-4-architecture.md)
3. **Inline help system** (80/100) — Every module has formatted help output with examples, tips, and command tables. Search functionality across all modules. (see agent-7-documentation.md)
4. **Zero eval/exec usage** — No eval or exec statements anywhere. Command dispatch via case statements. (see agent-3-security.md)
5. **Good input validation** (68/100) — Port numbers, volume, brightness, and ranges validated with is_number() and is_in_range(). (see agent-1-code-quality.md)

## Critical Bugs (must fix immediately)

1. **`rm -rf /tmp/*` in system cleanup** — system.sh:173. Deletes ALL files in /tmp including active Unix sockets, PID files, and temp databases for running applications. Can crash apps and corrupt data.
2. **Airport binary never found** — wifi.sh:174. The `airport` command lives off-PATH, so `command_exists airport` always fails, silently disabling WiFi scanning, signal strength, and detailed connection info.
3. **Wrong memory stats on Apple Silicon** — system.sh uses hardcoded 4096-byte page size. Apple Silicon uses 16KB pages. All memory numbers displayed by `mac system memory` are 4x too low on M1/M2/M3/M4 Macs.
4. **Installer missing 4 modules** — install.sh:154-162 only downloads 6 of 10 library files. Missing: wifi.sh, dock.sh, audio.sh, screenshot.sh. Users installing via curl|bash get broken commands.
5. **config.json is completely dead** — get_config() defined but never called. All config settings (volume_step, screenshot_location, etc.) are ignored at runtime.
6. **No LICENSE file** — README claims MIT license but no LICENSE file exists. Default copyright (all rights reserved) applies.

## Top 5 Weaknesses

1. **Massive code duplication** (DRY: 15/100) — 5,154 lines of identical files at project root, duplicating all lib/ files. 71.6% of the codebase is duplicated. (see agent-1-code-quality.md)
2. **Testing is nearly non-existent** (28/100) — Only 6 of 77 commands tested. Three modules have zero coverage. No output validation. Tests only check exit codes, not behavior. (see agent-5-testing.md)
3. **Security vulnerabilities** (38/100) — curl|bash install with no integrity verification, predictable temp files, unsanitized input to osascript, silent sudo usage. (see agent-3-security.md)
4. **Documentation accuracy** (40/100) — README documents commands that don't exist (mac config, mac modules), omits the screenshot module entirely, and architecture diagram is wrong. (see agent-7-documentation.md)
5. **Distribution readiness** (38/100) — No LICENSE, no integrity verification, broken install script, no update mechanism, inadequate .gitignore. (see agent-6-distribution.md)

## Cross-Agent Patterns (systemic issues found by 3+ agents)

| Pattern | Agents | Severity |
|---------|--------|----------|
| Root-level duplicate files | Code Quality, Architecture, Distribution, Core Logic | CRITICAL |
| config.json never consumed | Core Logic, Architecture, UX | CRITICAL |
| Missing modules from test/install lists | Testing, Distribution, Documentation | HIGH |
| Inconsistent version in 7+ locations | Code Quality, Distribution, Architecture | HIGH |
| No NO_COLOR/accessibility support | UX, Security | HIGH |
| set -e conflicts with pipe patterns | Code Quality, Core Logic, Testing | MEDIUM |

## Phased Upgrade Plan

### Phase 0: Critical Fixes (1-2 days)
- Delete all root-level duplicate files
- Fix rm -rf /tmp/* → targeted cleanup
- Fix airport binary path resolution
- Fix vm_stat page size for Apple Silicon
- Add missing modules to install.sh download list
- Create LICENSE file
- Fix dock.sh:131 string concatenation bug

### Phase 1: Foundation (3-5 days)
- Add shellcheck to CI
- Expand test suite (output validation, missing modules)
- Add NO_COLOR support
- Make config.json functional (or remove it)
- Single-source version string
- Fix .gitignore

### Phase 2: Quality (1 week)
- Split utils.sh into focused files
- Add module auto-discovery
- Add --yes flag for non-interactive use
- Consolidate boolean parsing
- Add integrity verification to install
- Improve uninstall support

### Phase 3: Growth (1-2 weeks)
- Homebrew formula
- GitHub Actions CI/CD
- Auto-update mechanism
- Plugin/module system
- Performance optimization (cache system_profiler calls)

## Metrics Snapshot

| Metric | Value |
|--------|-------|
| Total lines (all files) | 12,797 |
| Unique functional lines | ~7,200 |
| Duplicate lines (root copies) | 5,154 (71.6%) |
| Modules | 9 |
| Commands | 77+ |
| Test coverage (commands) | 7.8% (6/77) |
| Security vulnerabilities | 17 |
| Missing LICENSE | Yes |
| CI/CD | None |
| Files in audit | 24 |

## Audit Files

| File | Description |
|------|-------------|
| [agent-1-code-quality.md](agent-1-code-quality.md) | Code quality analysis (400KB) |
| [agent-2-core-logic.md](agent-2-core-logic.md) | Core logic & execution paths (386KB) |
| [agent-3-security.md](agent-3-security.md) | Security audit (356KB) |
| [agent-4-architecture.md](agent-4-architecture.md) | Architecture analysis (331KB) |
| [agent-5-testing.md](agent-5-testing.md) | Test coverage analysis (311KB) |
| [agent-6-distribution.md](agent-6-distribution.md) | Distribution & packaging (126KB) |
| [agent-7-documentation.md](agent-7-documentation.md) | Documentation audit (412KB) |
| [agent-8-ux-cli.md](agent-8-ux-cli.md) | UX & CLI design (344KB) |
| [scorecard.json](scorecard.json) | Machine-readable scores |
| [brief.md](brief.md) | Project brief for future sessions |
| [architecture.md](architecture.md) | Architecture map |
| [patterns.md](patterns.md) | Code patterns analysis |
| [risk-map.md](risk-map.md) | Risk assessment |
| [upgrade-plan.md](upgrade-plan.md) | Phased improvement plan |
| [surprises.md](surprises.md) | Non-obvious findings |
