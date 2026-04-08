# SuperMac Go Rewrite — Implementation Plan

**Created**: 2026-04-08
**Status**: Ready for execution
**Design ref**: `docs/brainstorming/2026-04-08-go-rewrite-design.md`
**Roadmap**: ROAD-017 (Foundation) → ROAD-018 (Modules) → ROAD-019 (Distribution)

## Execution Model

Phases 1-2 can use parallel GLM agents for independent worktrees.
Phase 3 (distribution) is sequential after modules are complete.

## Phase 1: Foundation (ROAD-017)

### Step 1: Scaffold Go project (ROAD-021)

```bash
mkdir -p supermac-go && cd supermac-go
go mod init github.com/cosmolabs-org/supermac
go get github.com/spf13/cobra
go get gopkg.in/yaml.v3
go get github.com/blang/semver
mkdir -p cmd/mac internal/module internal/modules internal/config internal/output internal/update internal/platform internal/version pkg/supermac completions
```

Create `cmd/mac/main.go`:
- Root command: `mac`
- Persistent flags: `--json`, `--quiet`, `--no-color`, `--debug`, `--version`
- Wire config loader, output writer, module registry

Create `internal/version/version.go`:
```go
package version

var (
    Version   = "dev"
    BuildDate = "unknown"
)
```

Create `Makefile`:
```makefile
BINARY = mac
VERSION ?= $(shell git describe --tags --always --dirty)

.PHONY: build test lint install clean

build:
	go build -ldflags "-X internal/version.Version=$(VERSION)" -o $(BINARY) ./cmd/mac

test:
	go test ./...

lint:
	golangci-lint run

install: build
	cp $(BINARY) ~/bin/

clean:
	rm -f $(BINARY)
```

Verify: `make build && ./mac --version` should print version.

### Step 2: Module interface (ROAD-022)

Create `internal/module/module.go`:
```go
package module

type Module interface {
    Name() string
    ShortDescription() string
    Emoji() string
    Commands() []Command
}

type Command struct {
    Name        string
    Description string
    Aliases     []string
    Args        []Arg
    Run         func(ctx *Context) error
}

type Arg struct {
    Name     string
    Required bool
}

type Context struct {
    Config  *config.Config
    Output  output.Writer
    Args    []string
    Verbose bool
    DryRun  bool
}
```

Create `internal/module/registry.go`:
```go
package module

import "sync"

var (
    modulesMu sync.RWMutex
    modules   = make(map[string]Module)
)

func Register(m Module) {
    modulesMu.Lock()
    defer modulesMu.Unlock()
    modules[m.Name()] = m
}

func All() map[string]Module {
    modulesMu.RLock()
    defer modulesMu.RUnlock()
    out := make(map[string]Module, len(modules))
    for k, v := range modules {
        out[k] = v
    }
    return out
}

func Get(name string) (Module, bool) {
    modulesMu.RLock()
    defer modulesMu.RUnlock()
    m, ok := modules[name]
    return m, ok
}
```

Wire into Cobra: iterate `module.All()`, create subcommand per module, nested subcommand per `Module.Commands()`.

### Step 3: Platform abstraction (ROAD-023)

Create `internal/platform/platform.go`:
```go
package platform

type Platform interface {
    RunOSAScript(script string) (string, error)
    ReadDefault(domain, key string) (string, error)
    WriteDefault(domain, key, value string) error
    RunCommand(name string, args ...string) (string, error)
}
```

Create `internal/platform/darwin.go` — real implementation using `os/exec`.
Create `internal/platform/mock.go` — test mock.

Each macOS tool becomes a function:
- `platform.GetAirportPath()` — returns the known airport binary path
- `platform.SetWiFi(on bool)` — calls networksetup
- `platform.GetWiFiStatus()` — calls airport -I
- `platform.SetVolume(level int)` — calls osascript
- `platform.GetMemoryInfo()` — calls sysctl and vm_stat
- etc.

### Step 4: Config system (ROAD-024)

Create `internal/config/config.go`:
```go
package config

type Config struct {
    Version int            `yaml:"version"`
    Output  OutputConfig   `yaml:"output"`
    Updates UpdatesConfig  `yaml:"updates"`
    Modules ModulesConfig  `yaml:"modules"`
    Aliases map[string]string `yaml:"aliases"`
}

type OutputConfig struct {
    Color  bool   `yaml:"color"`
    Format string `yaml:"format"` // text, json, quiet
}

type UpdatesConfig struct {
    Check   bool   `yaml:"check"`
    Channel string `yaml:"channel"` // stable, beta
}

type ModulesConfig struct {
    Screenshot ScreenshotConfig `yaml:"screenshot"`
    Audio      AudioConfig      `yaml:"audio"`
    Display    DisplayConfig    `yaml:"display"`
}
```

Functions:
- `Load() (*Config, error)` — reads ~/.supermac/config.yaml, creates with defaults if missing
- `Save(c *Config) error`
- `Default() *Config` — sensible defaults
- Add `mac config` command group: `edit`, `get <key>`, `set <key> <value>`, `list`

### Step 5: Output system (ROAD-025)

Create `internal/output/output.go`:
```go
package output

type Writer interface {
    Info(msg string, args ...interface{})
    Success(msg string, args ...interface{})
    Warning(msg string, args ...interface{})
    Error(msg string, args ...interface{})
    Header(title string)
    Table(headers []string, rows [][]string)
    JSON(v interface{}) error
}
```

Three implementations:
- `NewColoredWriter(w io.Writer)` — ANSI colors, icons, headers
- `NewJSONWriter(w io.Writer)` — structured JSON output
- `NewQuietWriter(w io.Writer)` — only errors

Auto-selection: check `--json` flag, `--quiet` flag, `NO_COLOR` env, TTY detection.

**Phase 1 deliverable**: `make build && ./mac help` shows all module categories (even if modules are stubs). `./mac version` works. `./mac config list` works.

---

## Phase 2: Module Porting (ROAD-018)

All 6 module port tasks are independent — can run as parallel agents in worktrees.

### Execution pattern (per module):

For each module (e.g., `wifi`):

1. Create `internal/modules/wifi/wifi.go`:
```go
package wifi

import "github.com/cosmolabs-org/supermac/internal/module"

type WifiModule struct{}

func init() {
    module.Register(&WifiModule{})
}

func (w *WifiModule) Name() string { return "wifi" }
func (w *WifiModule) ShortDescription() string { return "WiFi control and management" }
func (w *WifiModule) Emoji() string { return "🌐" }

func (w *WifiModule) Commands() []module.Command {
    return []module.Command{
        {Name: "on", Description: "Turn WiFi on", Run: w.on},
        {Name: "off", Description: "Turn WiFi off", Run: w.off},
        {Name: "toggle", Description: "Toggle WiFi state", Run: w.toggle},
        {Name: "status", Description: "Show WiFi status", Run: w.status},
        {Name: "scan", Description: "Scan for networks", Run: w.scan},
        // ...
    }
}
```

2. Create `internal/modules/wifi/wifi_test.go` — mock platform, test each command

3. Add blank import in `cmd/mac/main.go`:
```go
import (
    _ "github.com/cosmolabs-org/supermac/internal/modules/wifi"
)
```

4. Verify: `./mac wifi on` works, `./mac wifi status --json` returns JSON

### Module port order and parallelism:

**Batch 1** (3 agents in parallel):
- wifi module (ROAD-027) — uses airport binary, networksetup
- system module (ROAD-026) — uses sysctl, vm_stat, pmset. Most complex
- network module (ROAD-028) — uses networksetup, sudo for DNS flush

**Batch 2** (3 agents in parallel):
- display module (ROAD-029) — osascript for brightness, defaults for dark mode
- dev module (ROAD-030) — lsof, ps, most complex logic
- dock/finder/audio/screenshot bundle (ROAD-031) — all defaults read/write

### Each module agent receives:

- The Bash source file (`lib/<module>.sh`)
- The module interface spec
- The platform abstraction API
- Instructions to implement each function, write tests, add to main.go

**Phase 2 deliverable**: `./mac <category> <action>` works for all 10 modules. `make test` passes. Feature parity with Bash version.

---

## Phase 3: Distribution (ROAD-019)

Sequential — depends on all modules being complete.

### Step 1: GitHub Actions CI

`.github/workflows/ci.yml`:
- Trigger: push to main, PRs
- Steps: go test ./..., golangci-lint, build for darwin-arm64 and darwin-amd64

### Step 2: GitHub Actions Release

`.github/workflows/release.yml`:
- Trigger: tag push v*
- Steps: cross-compile, generate SHA256SUMS, create GitHub Release, upload artifacts

### Step 3: Install script

`scripts/install.sh`:
- Detect arch (uname -m)
- Download correct binary from GitHub Releases
- Verify SHA256 checksum
- Install to ~/bin/ (or configurable path)
- `curl -fsSL https://cosmolabs.org/install | bash`

### Step 4: Homebrew tap

Create `cosmolabs-org/homebrew-tap` repo:
- Formula: `supermac.rb` pointing to GitHub Release URLs with SHA256
- `brew install cosmolabs-org/tap/supermac`

### Step 5: Auto-update (FEAT-004)

`internal/update/update.go`:
- `Check()` — GET GitHub Releases API, compare semver
- `Apply()` — download new binary to temp, verify checksum, atomically replace
- Background check on launch (goroutine)
- `mac update` and `mac update --check` commands

### Step 6: Shell completions (FEAT-005)

- `mac completion bash > /etc/bash_completion.d/mac`
- `mac completion zsh > "${fpath[1]}/_mac"`
- `mac completion fish > ~/.config/fish/completions/mac.fish`
- Distribute in release tarball

**Phase 3 deliverable**: Full CI/CD pipeline. One-command install. Auto-update works. Homebrew tap live.

---

## Dependency Graph

```
ROAD-021 (Scaffold)
  └─→ ROAD-022 (Module interface)
     └─→ ROAD-023 (Platform layer)
     └─→ ROAD-024 (Config system)
     └─→ ROAD-025 (Output system)
        └─→ ROAD-026-031 (Module ports — all parallel)
           └─→ ROAD-019 (Distribution — sequential)
```

## Estimated Timeline

| Phase | Items | Duration | Parallelizable |
|-------|-------|----------|----------------|
| Foundation | ROAD-021 to 025 | 3-5 days | Partial (config + output parallel after scaffold) |
| Module Port | ROAD-026 to 031 | 5-8 days | Fully (6 parallel agents) |
| Distribution | ROAD-019 sub-items | 3-5 days | Partial |
| **Total** | **15 items** | **11-18 days** | |

## Verification Checklist

After each phase:

**Phase 1**: `make build && ./mac help && ./mac version && ./mac config list`
**Phase 2**: `make test` passes, `./mac wifi status --json` works, all 10 modules respond
**Phase 3**: Fresh install via curl script works, `brew install` works, `mac update --check` works

## Notes

- No Bash scripts ever again — Go from here on
- The Go CLI is the engine for the future desktop app (ROAD-020)
- Open-source CLI (MIT) + paid desktop app (freemium) business model
- Each module agent needs: the Bash source, the interface spec, platform API docs
