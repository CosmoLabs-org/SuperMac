---
title: "SuperMac Dependency System — Auto-Install + mac doctor"
created: 2026-04-10
status: approved
origin: "/brainplan"
tags: [dependency-management, homebrew, doctor, architecture]
---

# SuperMac Dependency System

## Problem

Three SuperMac modules depend on external CLI tools (blueutil, dockutil, SwitchAudioSource). Each module does ad-hoc `exec.LookPath` checks with hardcoded error messages. There is no centralized dependency management, no `mac doctor` health check, and no auto-install capability.

## Design

### Dependency Types

New package `internal/dep/dep.go`:

```go
type Dependency struct {
    Name     string   // display name: "blueutil"
    Brew     string   // brew formula name: "blueutil"
    Check    string   // binary to look for: "blueutil"
    Commands []string // which commands need this (nil = all commands in module)
}

func (d Dependency) IsInstalled() bool
func (d Dependency) Install() error
func (d Dependency) Ensure(interactive bool) error
```

Behaviors:
- `IsInstalled()` — uses `exec.LookPath(d.Check)`
- `Ensure(false)` — check only, error if missing (for `mac doctor`)
- `Ensure(true)` — check, prompt user, install if yes (for individual commands)
- `Install()` — runs `brew install <brew>`, streams output
- `Commands` field: `nil` = entire module needs this dep. Set = only those specific commands.

### Module Interface Change

Add `Dependencies()` to the Module interface:

```go
type Module interface {
    Name() string
    ShortDescription() string
    Emoji() string
    Commands() []Command
    Search(term string) []SearchResult
    Dependencies() []dep.Dependency
}
```

Modules with no external deps return `nil`. Only 3 modules return non-empty:

| Module | Dep | Brew Formula | Commands |
|--------|-----|-------------|----------|
| bluetooth | blueutil | blueutil | all |
| dock | dockutil | dockutil | add, remove |
| audio | SwitchAudioSource | switchaudio-osx | input-device, output-device |

### mac doctor Command

Registered in `main.go` as a built-in command. Iterates all modules, collects deps, checks each.

Output:
```
$ mac doctor

  SuperMac Health Check

  ✓ brew              installed (/opt/homebrew/bin/brew)

  Module Dependencies:
    ✓ blueutil           installed
    ✗ dockutil           missing     brew install dockutil
    ✓ SwitchAudioSource  installed

  3/4 checks passed. Run 'mac doctor --fix' to install missing dependencies.
```

`--fix` flag auto-installs all missing deps (skipping prompts). Without `--fix`, report only.

Also checks whether `brew` itself is installed.

### Auto-Install Flow

When a command runs and its dep is missing:

```
$ mac bluetooth status

  blueutil is not installed.

  Install via Homebrew? [Y/n]: y

  ==> Downloading blueutil...
  ==> Installing blueutil
  ✓ blueutil installed successfully.

  Retrying command...
```

**Hook point**: In `main.go`'s `registerModules()`, before calling `cmd.Run(ctx)`:

```go
RunE: func(subCmd *cobra.Command, args []string) error {
    if err := checkModuleDeps(mod, cmd.Name); err != nil {
        return err
    }
    return cmd.Run(ctx)
}
```

`checkModuleDeps` iterates `mod.Dependencies()`, checks if the current command name matches (or dep.Commands is nil for whole-module deps), and calls `dep.Ensure(true)`.

### What Gets Replaced

All 6 hardcoded `exec.LookPath` checks removed from:
- `bluetooth/bluetooth.go` (1 check)
- `dock/dock.go` (2 checks — add, remove)
- `audio/audio.go` (3 checks — input-device, output-device, fallback)

### Edge Cases

- `--yes` flag (`mac -y bluetooth status`) — skips prompt, auto-installs
- `--quiet` flag — suppresses prompt, just errors
- `brew` not installed — tells user to install Homebrew first with canonical install command
- Install fails — shows error, does not retry the command

### Testing

- `dep_test.go` — test `IsInstalled()` with known system tools, test `Ensure(false)` behavior
- Mock `exec.LookPath` for install tests
- Doctor command test with mock modules returning various dep states

### File Scope

```
Files created:
  - supermac-go/internal/dep/dep.go         (~80 lines)
  - supermac-go/internal/dep/dep_test.go    (~60 lines)

Files modified:
  - supermac-go/internal/module/module.go   (add Dependencies() to interface)
  - supermac-go/cmd/mac/main.go             (add doctor cmd, add dep checking in registerModules)
  - supermac-go/internal/modules/bluetooth/bluetooth.go  (remove LookPath, add Dependencies())
  - supermac-go/internal/modules/dock/dock.go            (remove LookPath, add Dependencies())
  - supermac-go/internal/modules/audio/audio.go          (remove LookPath, add Dependencies())
  - supermac-go/internal/modules/finder/finder.go        (add Dependencies() returning nil)
  - supermac-go/internal/modules/wifi/wifi.go            (add Dependencies() returning nil)
  - supermac-go/internal/modules/network/network.go      (add Dependencies() returning nil)
  - supermac-go/internal/modules/system/system.go        (add Dependencies() returning nil)
  - supermac-go/internal/modules/dev/dev.go              (add Dependencies() returning nil)
  - supermac-go/internal/modules/display/display.go      (add Dependencies() returning nil)
  - supermac-go/internal/modules/screenshot/screenshot.go (add Dependencies() returning nil)
  - supermac-go/internal/modules/apps/apps.go            (add Dependencies() returning nil)
  - supermac-go/internal/modules/power/power.go          (add Dependencies() returning nil)
```
