# SuperMac Go Rewrite Design

**Status**: Approved (revised after spec review)
**Date**: 2026-04-08
**Author**: GAB (CosmoLabs)
**Scope**: Full CLI rewrite from Bash to Go
**Review**: 4 critical, 7 medium, 6 low issues found — critical addressed below

## Overview

Rewrite SuperMac from Bash to Go. 1:1 port of all 10 modules maintaining the
same CLI interface (`mac <category> <action>`). Go-native improvements:
`--json` output, auto-generated shell completions, self-update, real config
system. The Go CLI becomes the engine for a future desktop application.

## Architecture

### Approach: Cobra + Package-per-Module

Cobra for CLI backbone (subcommands, help, completions, man pages). Each
module is an isolated Go package implementing a shared interface. Platform
calls (osascript, defaults, airport) centralized in a `platform` package.

### Project Structure

```
supermac/
  cmd/
    mac/main.go              Entrypoint (wires Cobra root command)
  internal/
    module/
      module.go              Module interface: Name, Commands, Help
      registry.go            Auto-discovers and registers modules
    modules/
      wifi/wifi.go           Implements module.Interface
      system/system.go
      network/network.go
      display/display.go
      dock/dock.go
      finder/finder.go
      dev/dev.go
      audio/audio.go
      screenshot/screenshot.go
    config/
      config.go              Loads ~/.supermac/config.yaml
    output/
      output.go              Colored, JSON, quiet formatters
      colors.go              ANSI color helpers
    update/
      update.go              Check GitHub Releases, prompt/apply
    platform/
      airport.go             macOS WiFi binary wrapper
      defaults.go            defaults read/write wrappers
      osascript.go           AppleScript bridge
    version/
      version.go             Single source of truth (stamped via ldflags)
  pkg/
    supermac/                Public API (for desktop app later)
  completions/               Generated shell completion scripts
  config.example.yaml        Sample config with comments
  Makefile                   build, install, test, lint, release
  go.mod
  go.sum
```

Key boundary: `internal/` is private to the CLI (can change freely). `pkg/`
is the stable public API the desktop app will consume.

## Module Interface

```go
// internal/module/module.go

type Module interface {
    Name() string
    ShortDescription() string
    Emoji() string
    Commands() []Command
    Search(term string) []SearchResult  // Full-text search across module commands
}

type Command struct {
    Name        string
    Description string
    Aliases     []string
    Args        []Arg
    Flags       []Flag                  // Per-command flags (e.g., --sort, --force)
    Run         func(ctx *Context) error
}

type Flag struct {
    Name         string
    Shorthand    string                 // Single char, e.g., "s" for --sort
    Description  string
    DefaultValue string
    Required     bool
}

type Arg struct {
    Name        string
    Required    bool
    Description string
}

type Context struct {
    Config   *config.Config
    Output   output.Writer
    Platform platform.Interface       // Injected — real or mock
    Prompt   PromptInterface          // Interactive confirmation (mockable)
    Args     []string
    Flags    map[string]string        // Parsed flag values
    Verbose  bool
    DryRun   bool
}

type SearchResult struct {
    Command     string
    Description string
    Module      string
}
```

### Interactive Prompts

```go
// internal/module/prompt.go
type PromptInterface interface {
    Confirm(msg string) (bool, error)       // Y/n dialog
    Input(msg string) (string, error)       // Free text input
    Select(msg string, opts []string) (int, error) // Choice selection
}
```

Destructive operations (`system cleanup`, `network reset`, `dock reset`) use
`ctx.Prompt.Confirm()` which can be mocked in tests. `--yes` flag sets a
non-interactive Prompt implementation that always returns true.

### Error Handling

```go
// internal/module/errors.go
type ExitError struct {
    Code    int
    Message string
}

const (
    ExitOK          = 0
    ExitError       = 1
    ExitUsage       = 2   // Bad command/flag usage
    ExitPermission  = 3   // sudo/permission denied
    ExitNetwork     = 4   // Network unreachable
    ExitNotFound    = 5   // Resource not found
)
```

In `--json` mode, errors are emitted as `{"error": {"code": 3, "message": "..."}}`.

### Command Flow

```
user types: mac wifi on
  -> Cobra root command parses args
  -> Registry looks up "wifi" module
  -> Module.Commands() returns all wifi commands
  -> Cobra matches "on" subcommand
  -> Command.Run(ctx) executes
  -> wifi.Run() calls platform.SetWiFi(true)
  -> ctx.Output prints result (colored, --json, or quiet)
```

### Adding a New Module

1. Create `internal/modules/<name>/<name>.go`
2. Implement `Module` interface (4 methods)
3. Add blank import in `cmd/mac/main.go`: `_ "supermac/internal/modules/<name>"`
4. Auto-registered via `init()` — help text, completions, all generated

No dispatcher file edits. No category map drift (BUG-004 class bug impossible).

## Platform Abstraction

The test seam that makes 80% coverage possible. Modules never call `exec.Command`
directly — they call `ctx.Platform` methods.

```go
// internal/platform/platform.go
type Interface interface {
    // osascript
    RunOSAScript(script string) (string, error)

    // defaults read/write/delete
    ReadDefault(domain, key string) (string, error)
    WriteDefault(domain, key, value string) error
    DeleteDefault(domain, key string) error

    // Network
    SetWiFi(on bool) error
    GetWiFiStatus() (*WiFiInfo, error)
    ScanWiFiNetworks() ([]Network, error)
    FlushDNS() error                // sudo-gated
    ResetNetwork() error            // sudo-gated

    // System
    GetMemoryInfo() (*MemoryInfo, error)
    GetCPUInfo() (*CPUInfo, error)
    GetBatteryInfo() (*BatteryInfo, error)
    GetHardwareInfo() (*HardwareInfo, error)
    GetPageSize() (int, error)

    // Display
    SetBrightness(level float64) error
    GetDarkMode() (bool, error)
    SetDarkMode(on bool) error

    // Audio
    GetVolume() (int, error)
    SetVolume(level int) error
    GetAudioDevices() ([]AudioDevice, error)

    // Process management
    ListProcesses(filter string) ([]Process, error)
    KillPort(port int) error
    GetPortUser(port int) (string, error)

    // General command execution
    RunCommand(name string, args ...string) (string, error)
    RunSudoCommand(name string, args ...string) (string, error)
}
```

### sudo-gated Operations

sudo operations use `ctx.Platform.RunSudoCommand()` which:
1. Checks if the current user has passwordless sudo (`sudo -n true`)
2. If yes: runs the command directly
3. If no: shells out to `sudo` which prompts the user natively
4. In `--dry-run` mode: prints the command that would run, returns success
5. In tests: `MockPlatform.RunSudoCommand()` is a no-op that records the call

Modules do not know about sudo — they just call platform methods. The platform
layer handles privilege escalation transparently.

### External Tool Dependencies

| Bash dependency | Go replacement | Status |
|----------------|----------------|--------|
| `python3 -m http.server` | `net/http` | Eliminated by Go |
| `uuidgen` | `github.com/google/uuid` | Eliminated by Go |
| `base64` | `encoding/base64` | Eliminated by Go |
| `jq` | `encoding/json` | Eliminated by Go |
| `osascript` | `ctx.Platform.RunOSAScript()` | Kept — no Go equivalent |
| `defaults` | `ctx.Platform.ReadDefault()` | Kept — no Go equivalent |
| `networksetup` | `ctx.Platform.SetWiFi()` | Kept — no Go equivalent |
| `airport` | `ctx.Platform.GetWiFiStatus()` | Kept — no Go equivalent |
| `system_profiler` | `ctx.Platform.GetHardwareInfo()` | Kept — no Go equivalent |
| `SwitchAudioSource` | Optional — detect at runtime | Keep as optional dep |
| `dockutil` | Optional — detect at runtime | Keep as optional dep |

External tools that can't be replaced return a friendly error:
`"This command requires SwitchAudioSource. Install with: brew install switchaudio-osx"`

## Config System

```yaml
# ~/.supermac/config.yaml
version: 1

output:
  color: true
  format: text         # text | json | quiet

updates:
  check: true
  channel: stable      # stable | beta

modules:
  screenshot:
    location: Desktop
    format: PNG
    shadow: false
  audio:
    volume_step: 10
  display:
    brightness_step: 10

aliases:
  kp: "dev kill-port"
  dark: "display dark-mode"
```

- First run creates config with sensible defaults
- Modules read their section: `ctx.Config.Modules.Screenshot.Location`
- CLI flags override config: `--json` beats `output.format: text`
- `mac config edit` opens in `$EDITOR`
- `mac config get <key>` / `mac config set <key> <value>` for scripting
- Aliases are single source of truth (no drift with hardcoded map)
- YAML for comments and readability (standard in Go CLIs)

### Config Migration (Bash → Go)

On first run, if `~/.supermac/config.yaml` doesn't exist but
`~/.supermac/config.json` does (from Bash version):
1. Read old JSON config
2. Convert to new YAML schema (alias format: `category:action` → `category action`)
3. Write `~/.supermac/config.yaml`
4. Rename old `config.json` to `config.json.bak`
5. Print migration notice

### Global Shortcuts & Built-in Commands

Config aliases are converted to Cobra aliases at registration time:

```go
// During Cobra command setup, for each alias in config:
rootCmd.AddCommand(&cobra.Command{
    Use: aliasName,
    Run: func(cmd *cobra.Command, args []string) {
        // Parse "dev kill-port" → route to dev module, kill-port action
    },
})
```

Built-in commands (not part of any module):
- `mac help` / `mac help <category>` — Help system
- `mac version` / `mac -v` — Version info
- `mac search <term>` — Full-text search across modules
- `mac config edit/get/set/list` — Config management
- `mac update` / `mac update --check` — Self-update
- `mac completion <shell>` — Shell completions

## Output System

```go
type Writer interface {
    Info(msg string)
    Success(msg string)
    Warning(msg string)
    Error(msg string)
    Header(title string)
    Table(rows []Row)
    JSON(v interface{})
}
```

- `--json`: structured JSON output for all commands (scripting, desktop API)
- `--quiet`: suppress all output except errors
- `--no-color`: strip ANSI (also respects `NO_COLOR` env var)
- Non-TTY auto-detection: disables colors and spinners when piped

## Auto-Update

- Background goroutine checks GitHub Releases API on launch
- Compares semver, notifies user if newer
- `mac update --check` shows available version only
- Respects config: `updates.check: false` disables

### Atomic Self-Replace

`mac update` uses a two-step atomic swap on macOS:

1. Download new binary to temp file
2. Verify SHA256 against checksums fetched from GitHub Release
3. Checksums verified against signed provenance attestation (GitHub attestations)
4. `mv current_binary current_binary.bak` (preserve old version)
5. `mv temp_file current_binary` (atomic rename on same filesystem)
6. `chmod +x current_binary`
7. If step 4-6 fail: rollback by restoring `.bak`
8. On next run: `.bak` file cleaned up if new version starts successfully

### Code Signing & Notarization

Release binaries must be:
- Signed: `codesign -s "Developer ID Application: CosmoLabs" mac`
- Notarized: `notarytool submit mac.zip --wait`
- Stapled: `xnotary staple mac`

macOS Gatekeeper will reject unsigned binaries. Notarization is required for
users who don't Homebrew install (Homebrew handles Gatekeeper differently).

### Checksum Integrity

SHA256SUMS file is:
- Generated at release time by CI
- Signed with GPG or GitHub attestation
- Fetched from the same GitHub Release but verified against the attestation
- Not just "matching hash" — provenance must be verifiable

## Versioning

Single source of truth in `internal/version/version.go`:

```go
var Version = "dev"
```

Stamped at build time via ldflags:

```
go build -ldflags "-X internal/version.Version=0.2.0" ./cmd/mac
```

Used by `mac version`, auto-update checker, and `--version` flag.

## Testing

- Each module testable in isolation via `ctx.Platform` (mock or real)
- `platform.Interface` is the test seam with ~25 methods covering all macOS calls
- Unit tests: mock platform, test module logic in isolation
- Integration tests: real platform calls on macOS CI runner
- Target: 80%+ unit coverage for modules, integration coverage for platform/
- `go test ./...` for everything, `go test ./internal/modules/wifi/...` for one
- CI: GitHub Actions on macOS runner (needs real osascript, defaults, etc.)
- `platform/` package itself has lower coverage (makes real system calls) — this
  is acceptable. The mock boundary is at the module/platform interface.

## Distribution

### Multi-channel

1. **GitHub Releases**: Binary + checksums per platform (darwin-amd64, darwin-arm64)
2. **Install script**: `curl -fsSL cosmolabs.org/install | bash` — detects arch, verifies checksum
3. **Homebrew tap**: `brew install cosmolabs-org/tap/supermac`
4. **npm/bun**: `bun install -g supermac` — thin wrapper downloading Go binary

### Makefile

```makefile
make build          # Current OS/arch
make build-all      # Cross-compile darwin-amd64, darwin-arm64
make install        # go install to GOPATH/bin
make test           # go test ./...
make lint           # golangci-lint
make completions    # Generate bash/zsh/fish
make release        # Tag + GitHub Release with checksums
```

### Release Artifacts

- `supermac_vX.Y.Z_darwin_arm64.tar.gz` (binary + completions)
- `supermac_vX.Y.Z_darwin_amd64.tar.gz`
- `SHA256SUMS` (signed)

## Key Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/spf13/cobra` | CLI framework (subcommands, help, completions) |
| `gopkg.in/yaml.v3` | Config file parsing |
| `golang.org/x/crypto` | SHA256 verification for downloads |
| `github.com/blang/semver` | Semantic version comparison |

## Constraints

- **macOS only** for now (uses osascript, defaults, airport, system_profiler)
- **No cgo** — pure Go for easy cross-compilation
- **Single binary** — no external dependencies at runtime
- **Bash compatibility** — same command interface, users can swap seamlessly
- **`pkg/` starts empty** — public API populated only after internal API stabilizes

## Future (Out of Scope for Initial Rewrite)

- Desktop application (SwiftUI wrapping the Go CLI)
- Bubble Tea TUI for rich interactive commands
- Linux support (would need platform abstraction layer)
- Plugin system for third-party modules
- Homebrew-core submission (requires community adoption)
