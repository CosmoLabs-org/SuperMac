---
status: PENDING
type: continuation
priority: high
created: 2026-04-07
project: SuperMac
version: v2.1.0
branch: main
audit_score: 50.5/100
bugs: BUG-001 through BUG-008
---

# SuperMac v2.1.0 -- Audit Follow-up Session

## Context

Comprehensive audit completed 2026-04-07 across 8 dimensions (code quality, core logic, security, architecture, testing, distribution, documentation, UX/CLI). Score: 50.5/100. Eight bugs filed. This session executes Phase 0 of the upgrade plan: critical fixes and project hygiene.

## Required Reading

Read these files FIRST before making any changes:

1. `docs/audit/latest/brief.md` -- Project architecture, patterns, tech stack (~70 lines)
2. `docs/audit/latest/risk-map.md` -- What NOT to touch without tests (~86 lines)
3. `docs/audit/latest/upgrade-plan.md` -- Full phased plan (~140 lines)

Key risk: `lib/utils.sh` is a god file (589 lines, 8 responsibilities) that ALL modules depend on. Do not restructure it this session. Also do NOT touch `system_cleanup()` beyond BUG-001 -- it is a single function handling 7+ destructive operations and needs its own dedicated session with tests first.

## Project Structure

```
SuperMac/
  mac                    -- Main dispatcher (350 lines, routes to lib/)
  lib/                   -- Modular category libraries (CANONICAL source)
  bin/
    install.sh           -- Installation script (MISSING 4 modules from download list)
  tests/
    test.sh              -- Test harness (50+ tests, only 7.8% command coverage)
  config/
    config.json          -- Settings (decorative -- get_config() is never called)
```

Root-level `.sh` files are DUPLICATES of `lib/` counterparts. They must be deleted, not edited.

## Task List (Priority Order)

### 1. BUG-001: Replace rm -rf /tmp/* with targeted cleanup
- **File**: `lib/system.sh` line 173
- **Current**: `rm -rf /tmp/*` (wipes ALL temp files, not just SuperMac's)
- **Fix**: `find "${TMPDIR:-/tmp}" -type f -user "$USER" -atime +7 -delete`
- **Risk**: LOW -- scoped to user-owned files older than 7 days
- **Test**: Run `mac system cleanup` and verify no unrelated temp files are deleted

### 2. BUG-002: Use full airport binary path
- **File**: `lib/wifi.sh` lines 174, 237, 288
- **Current**: `which airport` or relies on PATH lookup
- **Fix**: Define constant at top of file:
  ```bash
  AIRPORT_BIN="/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport"
  ```
  Use `$AIRPORT_BIN` in all three locations. Add a guard: if `[[ ! -x "$AIRPORT_BIN" ]]`, print error and return 1.
- **Risk**: LOW -- path is stable across macOS versions

### 3. BUG-003: Read page size from sysctl for Apple Silicon
- **File**: `lib/system.sh` (memory display functions)
- **Current**: `page_size=4096` hardcoded (wrong on Apple Silicon -- 16384 bytes)
- **Fix**: `page_size=$(sysctl -n hw.pagesize)`
- **Risk**: LOW -- sysctl is available on all macOS versions

### 4. BUG-004: Add missing modules to install.sh download list
- **File**: `bin/install.sh` lines 154-162
- **Current**: Downloads 6 of 10 modules
- **Missing**: `wifi.sh`, `dock.sh`, `audio.sh`, `screenshot.sh`
- **Fix**: Add the four missing module URLs to the download array
- **Risk**: CRITICAL if wrong -- this is the install script. Verify URLs match the actual raw GitHub format.

### 5. BUG-005: Fix dock.sh string concatenation
- **File**: `lib/dock.sh` line 131
- **Current**: `"Dock auto-hide $action"d!` -- malformed concatenation
- **Fix**: `"Dock auto-hide ${action}d!"`
- **Risk**: LOW -- single line fix, test with `mac dock autohide on` and `mac dock autohide off`

### 6. BUG-006: Fix dev.sh decimal comparison with bc
- **File**: `lib/dev.sh` line 281-283
- **Current**: `[[ "$cpu" > 5.0 ]]` -- uses STRING comparison (broken for cpu >= 10.0)
- **Fix**: Use `bc -l` for numeric comparison, matching the pattern already used in `dev_memory_hogs`:
  ```bash
  if echo "$cpu > 5.0" | bc -l | grep -q 1; then
  ```
- **Risk**: LOW -- same pattern already validated elsewhere in the same file

### 7. BUG-008: Delete root-level duplicate files
- **Files**: 14 files at project root (audio.sh, dev.sh, display.sh, dock.sh, finder.sh, network.sh, screenshot.sh, system.sh, utils.sh, wifi.sh, config.json, mac, and any others)
- **Action**: `git rm` each duplicate. Verify they are identical to `lib/` counterparts first.
- **Verification**: After deletion, `mac help` should still work (dispatcher sources from `lib/`)

### 8. Create LICENSE file
- **Type**: MIT
- **Copyright**: CosmoLabs
- **Year**: 2026
- **Content**: Standard MIT license text

### 9. Fix .gitignore
- **Current**: Only `.DS_Store` (and even that may be missing from committed .gitignore)
- **Add**: `.DS_Store`, `dist/`, `*.tar.gz`, `*.log`, `*.tmp`, `.idea/`, `.vscode/`

### 10. Create project CLAUDE.md
- **Location**: `/Users/gab/PROJECTS/SuperMac/CLAUDE.md`
- **Content**: Project-specific instructions for Claude Code sessions. Include:
  - Project overview (Bash CLI tool, modular architecture)
  - Architecture: dispatcher + lib/ modules pattern
  - How to test (`make test` or `bash tests/test.sh`)
  - How to lint (`make lint` or `shellcheck`)
  - Key patterns: `*_dispatch` functions, `lib/utils.sh` sourcing, `GLOBAL_SHORTCUTS`
  - Module structure conventions
  - Known constraints (no strict mode in modules, config.json is decorative)
  - DO NOT edit root-level .sh files (they were deleted; only lib/ is canonical)

### 11. Set up roadmap
- Use `ccs roadmap` or create `docs/roadmap/` with initial milestones:
  - v2.1.1: Phase 0 critical fixes (this session)
  - v2.2.0: Phase 1 foundation (shellcheck, expanded tests, single-source version, config decision)
  - v2.3.0: Phase 2 quality (utils.sh split, module auto-discovery, --yes flag)
  - v3.0.0: Phase 3 growth (Homebrew formula, CI/CD, auto-update)

### 12. Consider first official release
- After Phase 0 fixes are merged and tested
- Tag as `v2.1.1` (patch release with critical fixes) or `v0.1.0` (if rebranding as initial stable)
- Ensure LICENSE is in place before any release
- Create GitHub release with changelog

## Working Conventions

- **All edits to `lib/` files only** -- root duplicates are dead code
- **Test after each fix**: `bash tests/test.sh` and manual smoke test the affected command
- **Commit per bug fix** -- one commit per BUG-xxx with conventional commit format:
  - `fix(system): replace rm -rf /tmp/* with targeted user-owned cleanup`
  - `fix(wifi): use full airport binary path instead of PATH lookup`
  - etc.
- **Do NOT add `set -euo pipefail` to modules** -- that is a Phase 1 change requiring careful testing
- **Do NOT restructure utils.sh** -- that is a Phase 2 change
- **Do NOT modify the installer download logic** beyond adding the 4 missing module URLs (BUG-004)

## Success Criteria

After this session:
- [ ] All 8 bugs closed with verified fixes
- [ ] LICENSE file committed
- [ ] .gitignore updated and committed
- [ ] Project CLAUDE.md created
- [ ] Roadmap initialized
- [ ] `make test` passes
- [ ] `mac help` lists all categories correctly
- [ ] `mac system cleanup` no longer nukes /tmp
- [ ] `mac wifi status` works on both Intel and Apple Silicon
- [ ] Fresh install would download all 10 modules
- [ ] No duplicate source files in repository
