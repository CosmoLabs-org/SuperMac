---
title: "SuperMac Power Toggles Module"
created: 2026-04-08
status: approved
priority: high
branch: master
---

# Power Toggles Module — Design Spec

## Goal

Add a new `power` module to SuperMac Go that exposes 20 lesser-known macOS system toggles in one discoverable place. Every toggle follows the same UX pattern, making it easy for beginners and power users alike to find and flip system settings.

## Module Definition

```
mac power <toggle> [on|off|toggle]
```

- Module name: `power`
- Emoji: `⚡`
- All commands: show current state with no args, set with `on`/`off`, flip with `toggle`

## Commands (20)

### System Behavior

| Command | What it toggles | Implementation | Needs sudo |
|---------|----------------|----------------|------------|
| `caffeinate` | Prevent system/display sleep | Start `caffeinate -d` in background (track PID), kill on off | No |
| `spotlight-indexing` | Enable/disable Spotlight indexing | `mdutil -i on/off /` | Yes |
| `gatekeeper` | Allow apps from unidentified developers | `spctl --master-disable` / `--master-enable` | Yes |
| `crash-reporter` | Crash dialog mode (show/silent) | `defaults write com.apple.CrashReporter DialogType` | No |
| `quarantine` | File quarantine warnings on downloads | `defaults write com.apple.LaunchServices LSQuarantine` | No |

### Finder & Desktop

| Command | What it toggles | Implementation | Restarts |
|---------|----------------|----------------|----------|
| `hidden-files` | Show hidden files in Finder | `defaults write com.apple.finder AppleShowAllFiles` | Finder |
| `file-extensions` | Show all file extensions | `defaults write NSGlobalDomain AppleShowAllExtensions` | Finder |
| `desktop-icons` | Show icons on desktop | `defaults write com.apple.finder CreateDesktop` | Finder |
| `save-panel` | Expanded save dialogs by default | `defaults write NSGlobalDomain NSNavPanelExpandedStateForSaveMode` | None |
| `print-dialog` | Expanded print dialogs by default | `defaults write NSGlobalDomain PMPrintingExpandedStateForPrint` | None |

### Performance & UI

| Command | What it toggles | Implementation | Restarts |
|---------|----------------|----------------|----------|
| `animations` | Window open/close animations | `defaults write NSGlobalDomain NSAutomaticWindowAnimationsEnabled` | Log out |
| `smooth-scrolling` | Smooth scrolling animation | `defaults write NSGlobalDomain AppleScrollAnimationEnabled` | None |
| `transparency` | Reduce UI transparency | `defaults write com.apple.universalaccess reduceTransparency` | None |
| `dock-bounce` | Dock bounce animation on launch | `defaults write com.apple.dock launchanim` | Dock |

### Input

| Command | What it toggles | Implementation | Restarts |
|---------|----------------|----------------|----------|
| `function-keys` | F1-F12 as standard function keys | `defaults write com.apple.keyboard fnState` | None |
| `key-repeat` | Fast key repeat (2x default speed) | `defaults write NSGlobalDomain KeyRepeat` + `InitialKeyRepeat` | None |

### Security & Network

| Command | What it toggles | Implementation | Needs sudo |
|---------|----------------|----------------|------------|
| `firewall-stealth` | Firewall stealth mode (no ping response) | `/usr/libexec/ApplicationFirewall/socketfilterfw --setstealthmode` | Yes |
| `remote-login` | SSH remote login server | `systemsetup -setremotelogin on/off` | Yes |

### Developer

| Command | What it toggles | Implementation | Needs sudo |
|---------|----------------|----------------|------------|
| `developer-dir` | Xcode developer tools path | `xcode-select -p` / `--switch` | Yes (to switch) |
| `login-items` | Auto-show login window items | `defaults write com.apple.loginwindow autoLoginUser` | None |

## UX Pattern

Every command follows identical state machine:

```
mac power <name>        -> Show current state
mac power <name> on     -> Enable
mac power <name> off    -> Disable
mac power <name> toggle -> Flip current state
```

### Output Examples

**Show state:**
```
$ mac power hidden-files
  Hidden files: visible
  Use 'mac power hidden-files off' to hide
```

**Set state:**
```
$ mac power animations off
  Disabling window animations...
  Log out and back in for changes to take effect.
```

**Toggle:**
```
$ mac power dock-bounce toggle
  Setting dock bounce to off...
  Dock restarted. Bounce animation disabled.
```

**Sudo required:**
```
$ mac power gatekeeper off
  Error: Gatekeeper requires admin privileges
  Run: sudo mac power gatekeeper off
```

## Special Command: `status`

`mac power status` shows all 20 toggles and their current state:

```
$ mac power status
  Power Toggles
  ─────────────────────────────────────────────
  caffeinate       * active (PID 12345)
  hidden-files     - hidden
  file-extensions  * visible
  desktop-icons    * visible
  gatekeeper       * enabled
  crash-reporter   * dialog
  function-keys    - standard (media keys)
  spotlight        * indexing
  animations       * enabled
  smooth-scrolling  * enabled
  transparency     * enabled
  dock-bounce      * enabled
  ...
```

## Implementation Notes

### File Structure

- `supermac-go/internal/modules/power/power.go` — Single file, ~600-700 lines
- All toggles read/write via `defaults` commands or direct system calls
- No new platform interface methods needed — uses `exec.Command` directly

### Helper Pattern

```go
type toggle struct {
    name       string
    getter     func() (string, error)
    setter     func(on bool) error
    restarter  func()  // nil if no restart needed
    needsSudo  bool
    restartMsg string // "Log out and back in" etc.
}
```

Each command is a `toggle` struct. A shared `runToggle(ctx, toggle)` function handles the state machine logic (show/set/flip), reducing all 20 commands to ~10 lines each.

### Caffeinate Special Case

Caffeinate is unique — it's a long-running process, not a defaults toggle:
- `on`: start `caffeinate -d` in background, save PID to `/tmp/supermac-caffeinate.pid`
- `off`: read PID from file, kill process
- `status`: check if PID file exists and process is running

### Restarter Functions

Some toggles need a process restart to take effect:
- **Finder**: `killall Finder`
- **Dock**: `killall Dock`
- **SystemUIServer**: `killall SystemUIServer`
- **None**: No restart needed
- **Log out**: Print "Log out and back in for changes to take effect"

### Sudo Handling

Commands that need sudo check `sudo -n true` first. If it fails, print a helpful message with the exact command to re-run with sudo.

## Testing

- Unit tests with MockPlatform for defaults read/write
- Test state machine: show/on/off/toggle for a sample toggle
- Test caffeinate PID file management
- Test sudo detection

## File Scope

```yaml
files_created:
  - supermac-go/internal/modules/power/power.go
  - supermac-go/internal/modules/power/power_test.go
files_modified:
  - supermac-go/cmd/mac/main.go  # add blank import
  - supermac-go/docs/USAGE.md     # document new module
```

## Out of Scope

- SIP (System Integrity Protection) — too dangerous to toggle from CLI
- FileVault — encryption toggle requires reboot and recovery key management
- Night Shift scheduling — already in display module
- Time Machine — deserves its own module if needed
