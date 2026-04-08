# SuperMac Go Rewrite Design

**Status**: Approved
**Date**: 2026-04-08
**Author**: GAB (CosmoLabs)
**Scope**: Full CLI rewrite from Bash to Go

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
    Commands() []Command
}

type Command struct {
    Name        string
    Description string
    Aliases     []string
    Run         func(ctx Context) error
    Args        []Arg
}

type Context struct {
    Config    *config.Config
    Output    output.Writer
    Args      []string
    Verbose   bool
    DryRun    bool
}

type Arg struct {
    Name        string
    Required    bool
    Description string
}
```

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
- `mac update` downloads binary, replaces self, verifies SHA256
- `mac update --check` shows available version only
- Respects config: `updates.check: false` disables
- Integrity: SHA256 checksums on all downloads (fixes curl|bash vulnerability)

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

- Each module testable in isolation with `platform.Mock` interface
- `platform/` layer is the test seam: real calls in prod, mocked in tests
- Target: 80%+ coverage (current Bash: 7.8%)
- `go test ./...` for everything, `go test ./internal/modules/wifi/...` for one
- CI: GitHub Actions on macOS runner (needs real osascript, defaults, etc.)

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

## Future (Out of Scope for Initial Rewrite)

- Desktop application (SwiftUI wrapping the Go CLI)
- Bubble Tea TUI for rich interactive commands
- Linux support (would need platform abstraction layer)
- Plugin system for third-party modules
- Homebrew-core submission (requires community adoption)
