# SuperMac v2.1.0 -- Audit Findings

*Full audit conducted 2026-04-07 across 8 dimensions.*

---

## CRITICAL

### C-1: Duplicated source files (root vs lib/)
**Severity**: Critical | **Area**: Architecture

Every module exists in two places: root directory AND `lib/`. For example, `audio.sh` (root) and `lib/audio.sh` are identical copies. The dispatcher (`mac`) sources from `lib/`, meaning the root copies are dead code.

- **Impact**: Confusion about which file to edit; changes to root copies silently ignored. Any patch applied to root `.sh` files is lost.
- **Evidence**: `wc -l` shows matching line counts across all pairs. Root `mac` hardcodes `LIB_DIR="$SUPERMAC_ROOT/lib"`.
- **Fix**: Delete all root-level module copies. Keep only `lib/` versions. Update `.gitignore` if needed.
- **Files**: `audio.sh`, `dev.sh`, `display.sh`, `dock.sh`, `finder.sh`, `network.sh`, `screenshot.sh`, `system.sh`, `utils.sh`, `wifi.sh` (10 files, ~5,300 lines of dead code)

### C-2: `curl | bash` installation with no integrity check
**Severity**: Critical | **Area**: Security

The README recommends `curl -fsSL ... | bash` with no checksum or signature verification. The install script itself downloads additional files from GitHub raw URLs.

- **Impact**: MITM attack can inject arbitrary code during installation. Any CDN compromise or DNS hijacking yields full user-level code execution.
- **Evidence**: `README.md` line 29-30; `install.sh` lines 1-20.
- **Fix**: (1) Publish SHA256 checksums alongside releases. (2) Have install.sh verify checksum before execution. (3) Consider signed commits or GPG verification for high-trust installs.

---

## HIGH

### H-1: `set -euo pipefail` missing from module files
**Severity**: High | **Area**: Code Quality

The main dispatcher uses strict mode, but individual modules in `lib/` do not. They rely on being sourced, which means failures in module functions propagate unpredictably.

- **Impact**: Silent failures in module functions. Unset variables go undetected. Pipe failures swallowed.
- **Fix**: Add `set -euo pipefail` to each module, or document that modules must be sourced only and the caller provides strict mode.

### H-2: No input sanitization on user-supplied arguments
**Severity**: High | **Area**: Security

Several functions pass user input directly to `osascript`, `defaults write`, `lsof`, and shell commands without sanitization.

- **Impact**: Command injection via crafted arguments, especially in `osascript -e` strings where user input is interpolated.
- **Evidence**: `display.sh` brightness functions interpolate `$1` into AppleScript strings. `dev.sh` kill-port interpolates `$port` into `lsof -ti:$port` (safe for numeric input, but `is_number` check exists only in some paths).
- **Fix**: Validate all user input before interpolation. Use `printf '%s'` instead of variable expansion inside command strings. Whitelist allowed characters per parameter type.

### H-3: `sudo` operations without user consent prompt or dry-run
**Severity**: High | **Area**: Security / UX

`network.sh` and `system.sh` run `sudo` commands (DNS flush, network reset, log deletion) that can fail silently or require interactive password entry, breaking scripted usage.

- **Evidence**: `lib/network.sh` lines 203-211 (flush-dns), 295 (DHCP renew), 323-328 (network reset -- deletes plist files and unloads launch daemons). `lib/system.sh` lines 155-178 (log cleanup, font cache reset).
- **Impact**: Network reset deletes system configuration plists -- destructive operation with no confirmation. Log cleanup uses `sudo find ... -delete` which can remove files unexpectedly if paths change.
- **Fix**: (1) Add confirmation prompt before destructive `sudo` operations. (2) Add `--dry-run` flag. (3) Document which operations require elevated privileges.

### H-4: No automated test execution in CI
**Severity**: High | **Area**: Testing

Tests exist (`tests/test.sh`) but there is no CI configuration (no `.github/workflows/`, no `.travis.yml`, no `Makefile` CI target).

- **Impact**: Regressions land in main without detection. The 50+ tests are only run manually.
- **Fix**: Add GitHub Actions workflow that runs `make test` on push/PR.

---

## MEDIUM

### M-1: Config aliases defined in two places
**Severity**: Medium | **Area**: Architecture

Aliases exist in both `config/config.json` (`"aliases"` key) and `mac` dispatcher (`GLOBAL_SHORTCUTS` associative array). They are manually kept in sync.

- **Impact**: Config aliases are never read at runtime. They are dead configuration. Only `GLOBAL_SHORTCUTS` is used.
- **Fix**: Either read `config.json` aliases at dispatcher startup, or remove the aliases key from config to avoid confusion.

### M-2: No shellcheck integration
**Severity**: Medium | **Area**: Code Quality

`Makefile` has a `lint` target but no shellcheck configuration (`.shellcheckrc`). No inline shellcheck directives beyond one `# shellcheck source=` comment.

- **Fix**: Add `.shellcheckrc` with `source-path=lib/`, add shellcheck to `make lint`, run in CI.

### M-3: `install.sh` duplicates color/output functions from utils.sh
**Severity**: Medium | **Area**: Code Quality

`install.sh` redefines `print_success`, `print_error`, `print_info`, `print_warning`, `print_header` because it runs before SuperMac is installed and cannot source `utils.sh`.

- **Impact**: Divergent formatting if one is updated without the other.
- **Fix**: Extract a minimal standalone output library that `install.sh` can embed or source inline.

### M-4: Hardcoded version in multiple locations
**Severity**: Medium | **Area**: Maintenance

Version "2.1.0" is hardcoded in: `mac` (line 43), `utils.sh` (line 14), `install.sh` (line 17), `config.json` (line 2), `Makefile` (line 11), `README.md` badge, `setup.sh`.

- **Impact**: Version bumps require editing 6+ files. Easy to miss one.
- **Fix**: Single source of truth. Read version from one file (e.g., a `.version` file or `config.json`) in all other locations.

### M-5: No uninstall mechanism
**Severity**: Medium | **Area**: Distribution

`install.sh` creates symlinks and writes preferences but there is no `uninstall.sh` or `mac uninstall` command.

- **Impact**: Users must manually remove symlinks, config, and installed files. Leaves orphan files.
- **Fix**: Add `uninstall.sh` or `mac self-uninstall` that reverses all install.sh operations.

---

## LOW

### L-1: macOS version claim mismatch
**Severity**: Low | **Area**: Documentation

README badge says "macOS 12.0+" but code constant says `MIN_MACOS_VERSION="10.15"`.

- **Fix**: Align both to the actual minimum tested version.

### L-2: Missing `.gitignore` for build artifacts
**Severity**: Low | **Area**: Distribution

Only `.DS_Store` is ignored. No ignore rules for `*.log`, `*.tmp`, or editor swap files.

- **Fix**: Add standard `.gitignore` entries.

### L-3: No `--quiet` or `--json` output flag
**Severity**: Low | **Area**: UX

All output is human-readable formatted text. No machine-parseable output option for scripting.

- **Fix**: Add `--json` flag for programmatic consumption. Add `--quiet` to suppress banners/formatting.

### L-4: Test suite lacks isolation
**Severity**: Low | **Area**: Testing

Tests run against the live system (checking actual WiFi state, running actual `sw_vers`, etc.) with no mocking layer.

- **Impact**: Tests pass/fail depending on host machine state. Not reproducible in CI.
- **Fix**: Abstract system commands behind functions that can be overridden in test fixtures.

### L-5: `plugins/` directory contains only a symlink
**Severity**: Low | **Area**: Architecture

`plugins/internal` symlinks to an external project (`ClaudeCodeSetup`). `plugins/registry.json` is minimal (85 bytes). No actual SuperMac plugins exist.

- **Fix**: Document the plugin system intent, or remove the empty scaffolding to reduce confusion.

---

## Summary Table

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
