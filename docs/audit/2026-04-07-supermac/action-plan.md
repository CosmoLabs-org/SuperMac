# SuperMac v2.1.0 -- Audit Findings and Action Plan

*Full audit conducted 2026-04-07 across 8 dimensions.*

---

## Findings

### CRITICAL

**C-1: Duplicated source files (root vs lib/)** | Architecture

Every module exists in both root and `lib/`. The dispatcher sources from `lib/` only, making 10 root-level copies (~5,300 lines) dead code. Editing the wrong file silently loses changes.

**C-2: `curl | bash` installation with no integrity check** | Security

README recommends `curl -fsSL ... | bash` with no checksum or signature verification. MITM during install yields arbitrary code execution.

### HIGH

**H-1: `set -euo pipefail` missing from module files** | Code Quality

Dispatcher uses strict mode but sourced modules do not. Silent failures in module functions go undetected.

**H-2: No input sanitization on user-supplied arguments** | Security

User input interpolated directly into `osascript -e` strings, `defaults write`, and `lsof` calls without sanitization. Command injection possible via crafted arguments.

**H-3: `sudo` operations without confirmation** | Security/UX

`network.sh` runs `sudo rm -f` on system plists and `sudo launchctl unload/load`. `system.sh` runs `sudo find ... -delete`. No confirmation prompts or dry-run mode.

**H-4: No automated test execution in CI** | Testing

Test suite exists (`tests/test.sh`, 50+ tests) but no CI pipeline. No `.github/workflows/`, no CI Makefile target.

### MEDIUM

**M-1: Config aliases in two places, one unused** | Architecture
**M-2: No shellcheck integration** | Code Quality
**M-3: `install.sh` duplicates output functions from `utils.sh`** | Code Quality
**M-4: Hardcoded version in 6+ files** | Maintenance
**M-5: No uninstall mechanism** | Distribution

### LOW

**L-1: macOS version claim mismatch** (README: 12.0+, code: 10.15) | Documentation
**L-2: Incomplete `.gitignore`** | Distribution
**L-3: No `--quiet` or `--json` output mode** | UX
**L-4: Test suite lacks isolation/mocking** | Testing
**L-5: Empty plugin scaffolding** | Architecture

| ID | Severity | Area | Title |
|----|----------|------|-------|
| C-1 | Critical | Architecture | Duplicated source files (root vs lib/) |
| C-2 | Critical | Security | `curl \| bash` with no integrity check |
| H-1 | High | Code Quality | No strict mode in module files |
| H-2 | High | Security | No input sanitization on user arguments |
| H-3 | High | Security/UX | Destructive `sudo` without confirmation |
| H-4 | High | Testing | No CI pipeline for test automation |
| M-1 | Medium | Architecture | Config aliases in two places, one unused |
| M-2 | Medium | Code Quality | No shellcheck integration |
| M-3 | Medium | Code Quality | Duplicated output functions in install.sh |
| M-4 | Medium | Maintenance | Hardcoded version in 6+ files |
| M-5 | Medium | Distribution | No uninstall mechanism |
| L-1 | Low | Documentation | macOS version claim mismatch |
| L-2 | Low | Distribution | Incomplete .gitignore |
| L-3 | Low | UX | No --quiet or --json output mode |
| L-4 | Low | Testing | Test suite lacks isolation/mocking |
| L-5 | Low | Architecture | Empty plugin scaffolding |

---

## Action Plan

*Prioritized remediation roadmap.*

---

## Phase 1: Critical Fixes (Day 1)

### 1A. Remove duplicated source files (C-1)
**Why first**: Every edit risks hitting the wrong file. This blocks all other work.

Steps:
1. Verify root `.sh` files are byte-identical to `lib/` counterparts: `diff audio.sh lib/audio.sh` for each pair
2. `git rm` all 10 root-level module copies (audio.sh, dev.sh, display.sh, dock.sh, finder.sh, network.sh, screenshot.sh, system.sh, utils.sh, wifi.sh)
3. Verify `mac` dispatcher still sources correctly from `lib/`
4. Run `make test` to confirm nothing breaks
5. Commit: `refactor: remove duplicated module files from project root`

### 1B. Secure the install pipeline (C-2)
**Why first**: Security issue affects every user who installs.

Steps:
1. Add `SHA256SUMS` file to the repository
2. Generate checksums for all distributable scripts: `shasum -a 256 install.sh mac lib/*.sh > SHA256SUMS`
3. Update `install.sh` to download `SHA256SUMS` and verify before executing
4. Add GPG sign-off to release tags
5. Update README to show verification instructions
6. Commit: `security: add integrity verification to installation pipeline`

---

## Phase 2: High-Priority Fixes (Days 2-3)

### 2A. Input validation and sanitization (H-2)
Steps:
1. Create `lib/validate.sh` with reusable input sanitization functions:
   - `validate_port()` -- reject non-numeric, enforce 1-65535 range
   - `validate_path()` -- reject paths with shell metacharacters
   - `validate_percent()` -- reject non-numeric, enforce 0-100 range
   - `sanitize_string()` -- escape dangerous characters for osascript
2. Apply to all `osascript -e` calls in display.sh, audio.sh
3. Apply to all `lsof` calls in dev.sh
4. Apply to all `defaults write` calls in dock.sh, finder.sh, screenshot.sh
5. Commit: `security: add input validation across all modules`

### 2B. Safe sudo operations (H-3)
Steps:
1. Add `confirm_action()` to `utils.sh` with Yes/No prompt and `--force` bypass
2. Wrap all destructive `sudo` operations in confirmation:
   - `network:reset` -- warn that this deletes network preferences
   - `system:cleanup` -- list what will be deleted before proceeding
   - `network:flush-dns` -- low risk but still add `--yes` flag for scripting
3. Add `--dry-run` flag to system cleanup that shows what would be deleted
4. Commit: `fix: add confirmation prompts for destructive sudo operations`

### 2C. Strict mode for modules (H-1)
Steps:
1. Add header comment to each module: `# This file is sourced by the SuperMac dispatcher. Strict mode is provided by the caller.`
2. OR add `set -uo pipefail` to each module (not `-e` since sourced files should not exit the parent)
3. Run full test suite after change
4. Commit: `fix: add error handling directives to module files`

### 2D. CI pipeline (H-4)
Steps:
1. Create `.github/workflows/test.yml`:
   - Trigger on push to main, pull requests
   - Run `make lint` (after adding shellcheck, see M-2)
   - Run `make test`
   - Cache shellcheck binary
2. Add status badge to README
3. Commit: `ci: add GitHub Actions test pipeline`

---

## Phase 3: Medium-Priority Improvements (Week 1)

### 3A. Single source of truth for version (M-4)
Steps:
1. Create `.version` file containing `2.1.0`
2. Update `mac` dispatcher to read: `SUPERMAC_VERSION=$(cat "$SUPERMAC_ROOT/.version")`
3. Update `utils.sh` similarly
4. Update `Makefile`: `VERSION := $(shell cat .version)`
5. Update `install.sh` to read version from repo
6. Remove hardcoded version strings everywhere else
7. Commit: `refactor: single source of truth for version number`

### 3B. Unify alias config (M-1)
Steps:
1. Make `mac` dispatcher read aliases from `config.json` at startup using `jq` or simple bash JSON parser
2. Remove `GLOBAL_SHORTCUTS` hardcoded map
3. Fall back to hardcoded map if `config.json` is missing or `jq` unavailable
4. Commit: `refactor: load aliases from config.json at runtime`

### 3C. Shellcheck integration (M-2)
Steps:
1. Create `.shellcheckrc`:
   ```
   source-path=lib/
   external-sources=true
   ```
2. Update `Makefile` lint target:
   ```makefile
   lint:
       shellcheck lib/*.sh mac bin/install.sh tests/test.sh
   ```
3. Fix all shellcheck warnings
4. Commit: `quality: add shellcheck linting`

### 3D. Add uninstall command (M-5)
Steps:
1. Create `uninstall.sh` that reverses `install.sh`:
   - Remove `~/bin/mac` symlink
   - Remove installed files from `~/.supermac/`
   - Optionally remove config preferences
   - Confirm before deleting anything
2. Add `mac self-uninstall` route
3. Commit: `feat: add uninstall capability`

### 3E. Reduce install.sh duplication (M-3)
Steps:
1. Extract minimal output functions into `lib/output-minimal.sh` (no dependencies)
2. Inline-source this in `install.sh` via heredoc or include
3. Commit: `refactor: shared minimal output library for install.sh`

---

## Phase 4: Low-Priority Polish (Week 2+)

| ID | Task | Estimate |
|----|------|----------|
| L-1 | Align macOS version claim (README badge vs code constant) | 5 min |
| L-2 | Expand `.gitignore` with standard entries | 10 min |
| L-3 | Add `--quiet` and `--json` output flags | 2 hours |
| L-4 | Add test mocking layer for system commands | 4 hours |
| L-5 | Document plugin system or remove empty scaffolding | 1 hour |

---

## Effort Estimate

| Phase | Items | Estimated Time |
|-------|-------|---------------|
| Phase 1 (Critical) | 1A, 1B | 3-4 hours |
| Phase 2 (High) | 2A-2D | 6-8 hours |
| Phase 3 (Medium) | 3A-3E | 4-6 hours |
| Phase 4 (Low) | L-1 to L-5 | 7-8 hours |
| **Total** | **16 findings** | **20-26 hours** |

## Execution Order Rationale

1. **Remove duplicates first** -- prevents all future work from hitting the wrong file
2. **Secure install pipeline** -- blocks real-world attack vector for users
3. **Input validation before strict mode** -- security beats code quality
4. **CI before refactoring** -- gives a safety net for subsequent changes
5. **Version unification before other refactors** -- reduces merge conflicts
6. **Low-priority items last** -- nice-to-have, not blockers
