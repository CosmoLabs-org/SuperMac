# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2026-04-08

### Added

**Go rewrite** — Complete rewrite from Bash to Go with 1:1 parity plus extras.

- **12 modules, 127 commands** — finder, wifi, network, system, dev, display, dock, audio, screenshot, bluetooth, apps, power
- **DarwinPlatform** — Real macOS system calls via `platform.Interface`
- **92 unit tests** — Stub platform + mock pattern for full module coverage
- **Power module** — 20 developer toggles (caffeinate, hidden-files, gatekeeper, animations, etc.)
- **Bluetooth module** — 6 commands via blueutil (status, connect, disconnect, power, discoverable)
- **Apps module** — 6 commands (list, info, cache-clear, recent, kill, open)
- **21 new commands** across existing modules (disk-usage, processes, uptime, speed-test, connections, etc.)
- **8 global shortcuts** — ip, cleanup, restart-finder, kp, vol, dark, light, search
- **Shell completions** — `mac completion bash|zsh|fish|powershell`
- **CI/CD** — GitHub Actions for test + build on push, tag-triggered releases (darwin arm64/amd64)
- **Install script** — curl-able installer with arch detection and checksum verification
- **Homebrew formula** — `Formula/supermac.rb` for tap-based installation
- **Backward-compatible aliases** — display light-mode/toggle-mode, network info, audio input/output
- Cobra-based CLI with `--json`, `--quiet`, `--no-color`, `--verbose`, `--dry-run`, `--yes` flags

## [0.1.2] - 2026-04-08

### Added
- Add MIT LICENSE, .gitignore, project CLAUDE.md
- Initialize roadmap with 31 items across 5 phases
- Add Go rewrite design spec and implementation plan

### Removed
- Remove 14 root-level duplicate files (5,154 lines)

### Fixed
- Apply 6 critical bug fixes from audit (BUG-001 through BUG-006)
- Replace declare -A with Bash 3.2-compatible lookup (BUG-009)
- Add 4 missing modules to install script (BUG-004)
- replace declare -A with Bash 3.2-compatible lookup (commit:18efc0d1)
- add 4 missing modules and deduplicate version (commit:c9c08b83)
- apply 6 critical bug fixes from audit (commit:7cf780fc)

## [0.1.1] - 2026-04-07

### Added
- file 8 bugs from audit findings (commit:15449267)
- comprehensive 360° project audit (8 agents, 50.5/100) (commit:a2e63703)

