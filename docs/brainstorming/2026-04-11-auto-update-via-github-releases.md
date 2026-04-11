---
completed: "2026-04-11"
created: "2026-04-11"
origin: /brainplan FEAT-004
status: COMPLETED
tags:
    - auto-update
    - distribution
    - github-releases
    - cross-project
title: Auto-Update via GitHub Releases
---

# Auto-Update via GitHub Releases

## Purpose

A self-update system for SuperMac that checks GitHub Releases for new versions, notifies the user on launch, and applies updates via atomic binary swap with SHA256 verification and one-version rollback.

**Cross-project significance**: This is the first implementation of a CosmoLabs auto-update pattern. After shipping, detailed PoC feedback will be sent to CCS so a universal release system can be abstracted for all CosmoLabs CLI and desktop apps (GoRalph, future tools). The design is intentionally clean and extractable.

## Approach

**Approach A: Minimal Checker + Manual Update** — background check on launch with 24h cache, user runs `mac update` to apply. Matches the pattern used by `gh`, `rustup`, and `brew`.

### Why not automatic update?

Unpredictable behavior — background binary replacement while the CLI is running is risky on slow/shaky connections and hard to test. Manual update gives the user control.

### Why not an external updater binary?

Two binaries to distribute, complex install/uninstall. Overkill for a CLI tool. The atomic swap with `.bak` rollback is reliable enough.

---

## Section 1: Package Structure

```
supermac-go/internal/update/
├── checker.go         # GitHub Releases API client + cache
├── updater.go         # Download, verify SHA256, atomic swap
├── checker_test.go    # 5 tests
└── updater_test.go    # 5 tests
```

### Core Types

```go
// checker.go — API client + cache
type Release struct {
    Version   string    // "0.2.2" (no v prefix)
    Tag       string    // "v0.2.2"
    TarballURL string   // download URL for arch-matched tarball
    ChecksumsURL string // checksums.txt URL
    Date      time.Time // published_at
}

type Checker struct {
    repo      string        // "CosmoLabs-org/SuperMac"
    cachePath string        // ~/.supermac/update-cache.json
    cacheTTL  time.Duration // 24h
    current   string        // current version from ldflags
    arch      string        // arm64 or amd64
}

func (c *Checker) CheckAvailable(ctx context.Context) (*Release, error)
```

```go
// updater.go — download, verify, swap
type Updater struct {
    binaryPath string // path to current binary (os.Executable)
    repo       string
    arch       string
}

func (u *Updater) Update(release *Release) error
func (u *Updater) Rollback() error
```

No interfaces — SuperMac-specific. When CCS extracts this later, interfaces for `Source` and `Verifier` will be added then. YAGNI for now.

**Zero new dependencies.** Pure stdlib: `net/http`, `archive/tar`, `compress/gzip`, `crypto/sha256`, `encoding/json`.

---

## Section 2: Update Check Flow

### On Launch

Triggered in `rootCmd.PersistentPreRunE` in `cmd/mac/main.go`.

```
1. Load config -> cfg.Updates.Check
2. If false, skip entirely
3. Read cache file (~/.supermac/update-cache.json)
4. If cache exists and age < 24h, use cached result
5. If cache stale or missing, run synchronous check with 2s timeout:
   a. GET https://api.github.com/repos/CosmoLabs-org/SuperMac/releases/latest
   b. Parse tag_name -> strip "v" -> semver compare with current version
   c. Match asset: mac-{arch}.tar.gz
   d. Write Release to cache file with timestamp
   e. On timeout or error: silently fall back to cached result (if any)
6. If new version available, print one-line notice before command output:
   "  ^ SuperMac v0.2.2 available. Run 'mac update' to upgrade."
```

**Why synchronous, not background goroutine?** A background goroutine has no deterministic point to print — it could interleave with command output, creating garbled text. A synchronous check with a strict 2s timeout is fast enough to be imperceptible, and guarantees the notification prints before any command output. On timeout, the check is silent and the cache is used next time.

### Cache File Format

`~/.supermac/update-cache.json`:

```json
{
  "checked_at": "2026-04-11T10:00:00Z",
  "latest": {
    "version": "0.2.2",
    "tag": "v0.2.2",
    "tarball_url": "https://github.com/CosmoLabs-org/SuperMac/releases/download/v0.2.2/mac-arm64.tar.gz",
    "checksums_url": "https://github.com/CosmoLabs-org/SuperMac/releases/download/v0.2.2/checksums.txt",
    "published_at": "2026-04-11T10:20:00Z"
  },
  "current_version": "0.2.1"
}
```

### Key Behaviors

- Synchronous check with 2s timeout — never blocks startup meaningfully
- GitHub API errors/timeouts are silent — fall back to cache, no noisy warnings
- Rate limit aware: 24h cache = max 1 API call per day per user (well within 60/hr unauthenticated limit)
- `--quiet` flag suppresses the update notification
- Dev builds (version = "dev") always show "update available" since dev != any release

### Manual Check

`mac update --check` — forces a fresh check (ignores cache), prints result synchronously, exits.

---

## Section 3: Update & Atomic Swap Flow

### `mac update` (synchronous, no background)

```
1. Resolve current binary path: os.Executable() -> e.g. /usr/local/bin/mac
2. Verify binary is writable (or will need sudo)
3. Fetch latest release info (from cache if fresh, otherwise API call)
4. If already on latest: "SuperMac is up to date (v0.2.2)"
5. Download tarball to $TMPDIR/supermac-update/
6. Download checksums.txt to $TMPDIR/supermac-update/
7. Verify: SHA256 computed against tarball, compare with checksums.txt entry
8. Extract binary from tarball -> $TMPDIR/supermac-update/mac
9. Verify extracted binary: code-sign it (codesign -f -s -), then exec with "version" arg,
   parse output for version string to confirm it matches the expected release version
10. Atomic swap:
    a. Rename current binary: /usr/local/bin/mac -> /usr/local/bin/mac.bak
       (same filesystem — guaranteed on standard macOS installs)
    b. Copy (io.Copy, NOT rename) new binary: $TMPDIR/supermac-update/mac -> /usr/local/bin/mac
       Uses io.Copy because $TMPDIR may be on a different filesystem than install dir
    c. Preserve permissions: chmod to match old binary
    d. Code-sign the new binary: codesign -f -s - /usr/local/bin/mac
       (required for macOS Gatekeeper — matches CLAUDE.md binary deployment policy)
    e. Remove quarantine: xattr -d com.apple.quarantine on new binary
    f. Cleanup: rm -rf $TMPDIR/supermac-update/
11. Report: "SuperMac updated v0.2.1 -> v0.2.2"
12. Error recovery: if step 10b fails, rename mac.bak back to mac (restore)
```

### Rollback (`mac update --rollback`)

```
1. Check if mac.bak exists next to current binary
2. Read its version: exec mac.bak version
3. Prompt: "Rollback to v0.2.1? [Y/n]"
4. Swap: current -> mac.old, mac.bak -> current
5. Cleanup: rm mac.old
6. Report: "Rolled back to v0.2.1"
```

### Sudo Handling

If the binary is in `/usr/local/bin` and not writable by current user, swap steps use `sudo` with a clear prompt. If `$HOME/bin` or writable location, no sudo needed.

### Homebrew Detection

If the binary path contains `Cellar` or `homebrew`, warn and redirect:

```
Installed via Homebrew. Run 'brew upgrade supermac' instead.
```

### Error Handling

| Scenario | Behavior |
|----------|----------|
| Network error during download | "Download failed: <error>. Check your connection and try again." |
| SHA256 mismatch | "Checksum mismatch! Download may be corrupted. Aborting." |
| Binary verification fails | "Downloaded binary failed verification. Aborting." |
| Swap fails (permission) | "Cannot replace binary. Try: sudo mac update" |
| Swap fails (other) | Restore mac.bak, report error |
| No .bak file for rollback | "No previous version found for rollback." |
| Concurrent update runs | PID-based lock file at `$TMPDIR/supermac-update.lock` — second instance detects and exits |

---

## Section 4: CLI Commands & Config

### Prerequisite: Fix Config Display Bug

`cmd/mac/main.go` line 140 displays the wrong field for Updates status:

```go
// BUG: prints cfg.Output.Format instead of cfg.Updates.Check
fmt.Printf("  Updates:  %v (%s)\n", cfg.Output.Format, cfg.Updates.Channel)
```

Fix to:
```go
fmt.Printf("  Updates:  %v (%s)\n", cfg.Updates.Check, cfg.Updates.Channel)
```

### Commands

```
mac update              # Update to latest version
mac update --check      # Check for available update (no download)
mac update --rollback   # Restore previous version from .bak
# Note: global --yes/-y flag skips confirmation prompts (already exists on rootCmd)
```

### Config

Already wired in `internal/config/config.go`. No schema changes needed.

```yaml
updates:
  check: true        # Enable/disable background check on launch
  channel: stable    # Only checks latest non-prerelease. "beta" reserved for future use.
```

**Beta channel**: The `beta` value is reserved for future implementation. When set, the checker behaves identically to `stable` and prints a one-time notice: "Beta channel is not yet supported. Using stable."

### Version Output Enhancement

Add `--raw` flag to `mac version` for machine-parseable output (used by the updater for binary verification):

```bash
$ mac version --raw
0.2.2
```

Normal output gains an update status line:
$ mac version
+------------------+
| SuperMac v0.2.2  |
+------------------+
  Version:    0.2.2
  Build:      2026-04-11T10:20:00Z
  Modules:    apps, audio, bluetooth, dev, display, dock, finder, network, power, screenshot, system, wifi
  Update:     up to date          <-- new line (or "v0.2.3 available")
```

### Startup Notification (when update available)

```bash
$ mac network ip
  ^ SuperMac v0.2.3 available. Run 'mac update' to upgrade.    <-- before command output
  Local IP address: 192.168.8.224
```

Suppressed by: `--quiet`, `updates.check: false`, or fresh cache with no new release.

---

## Section 5: Testing Strategy

### Unit Tests (No Network)

| Test | What it verifies |
|------|-----------------|
| `TestChecker_ParseRelease` | GitHub API JSON -> Release struct, v-prefix stripped |
| `TestChecker_SemverCompare` | "0.2.2" > "0.2.1", equal, pre-release handling |
| `TestChecker_CacheRead` | Reads cache file, respects TTL |
| `TestChecker_CacheExpired` | Ignores stale cache (>24h) |
| `TestUpdater_VerifyChecksum` | SHA256 matches checksums.txt entry |
| `TestUpdater_ExtractBinary` | Tarball -> correct binary for arch |
| `TestUpdater_AtomicSwap` | Rename -> copy -> verify -> cleanup |
| `TestUpdater_SwapFails_Restore` | .bak restored on failure |
| `TestUpdater_Rollback` | .bak -> current, cleanup |

### Testing Approach

- Mock HTTP via `httptest.NewServer` serving fake GitHub API JSON responses
- Create real tarballs in temp dirs for extraction tests
- No real network calls in any test
- Tests run in parallel (`t.Parallel()`)

### Edge Cases

| Edge case | Handling |
|-----------|----------|
| Homebrew install | Detect Cellar path, redirect to `brew upgrade` |
| macOS quarantine | `xattr -d com.apple.quarantine` on new binary |
| Concurrent updates | flock on binary path |
| GitHub API rate limit | 24h cache = 1 call/day, well within 60/hr |
| Dev version ("dev") | Always shows "update available" |
| Pre-release tags | Ignored — only latest non-prerelease release |

---

## Size Estimate

| File | Lines |
|------|-------|
| `internal/update/checker.go` | ~120 |
| `internal/update/updater.go` | ~180 |
| `internal/update/checker_test.go` | ~150 |
| `internal/update/updater_test.go` | ~200 |
| `cmd/mac/main.go` (changes) | ~30 |
| **Total** | **~680** |

## Dependencies

Zero new external dependencies. All stdlib:
- `net/http` — GitHub API calls
- `archive/tar` + `compress/gzip` — tarball extraction
- `crypto/sha256` — checksum verification
- `encoding/json` — API + cache parsing
- `os` + `path/filepath` — file operations

## Future Extraction Notes (for CCS feedback)

When extracting to a shared CosmoLabs package:
1. Add `Source` interface (GitHub Releases, R2/Cloudflare, custom manifest)
2. Add `Verifier` interface (SHA256, code signature, GPG)
3. Extract `Checker` and `Updater` into `cosmo-go/update` package
4. Config stays per-project — only the update logic moves
5. The cache format and atomic swap pattern are reusable as-is
