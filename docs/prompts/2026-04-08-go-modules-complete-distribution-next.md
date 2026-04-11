---
branch: master
completed: "2026-04-11"
created: "2026-04-08"
goals_completed: 5
goals_total: 5
origin: /continuation-prompt
priority: critical
related_prompts:
    - docs/prompts/2026-04-08-go-modules-and-fixes.md
    - docs/prompts/2026-04-08-go-rewrite-implementation.md
    - docs/prompts/2026-04-07-supermac-audit-followup.md
started: "2026-04-08"
status: COMPLETED
tags:
    - continuation
    - go-rewrite
    - distribution
    - release
title: SuperMac Go Rewrite — Modules Complete, Distribution Phase Next
---

# SuperMac Go Rewrite — Modules Complete, Distribution Phase Next

## File Scope

```yaml
files_modified:
  - supermac-go/cmd/mac/main.go
  - supermac-go/internal/platform/darwin.go
  - supermac-go/docs/USAGE.md
  - docs/roadmap/items/ROAD-018.yaml
files_created:
  - supermac-go/internal/modules/power/power.go
  - supermac-go/internal/modules/power/power_test.go
  - supermac-go/internal/modules/bluetooth/bluetooth.go
  - supermac-go/internal/modules/apps/apps.go
  - docs/brainstorming/2026-04-08-power-toggles-module.md
```

## Context

The Go rewrite is essentially feature-complete. This session was a massive sprint: started with 9 unported Bash modules and a MockPlatform that returned empty strings. Ended with 12 modules, 127 commands, 40 passing tests, a real DarwinPlatform implementation, and full 1:1 parity with the original Bash CLI — plus 15+ Go-only extras. The tool works end-to-end with real system data. All code is committed on master. The energy is high and the architecture is solid.

What's left is the distribution layer: CI/CD, install script, Homebrew tap, shell completions. Then a v0.2.0 release. The original Bash CLI stays as `mac` while the Go binary is at `supermac-go/mac` — distribution will handle the transition.

## What Got Done (This Session)

**11 commits on master:**

1. **9 module ports** (5ca02d6) — Ported dock, system, wifi, network, display, dev, audio, screenshot from Bash to Go. Created DarwinPlatform with real macOS system calls. All 9 modules wired into main.go via blank imports. 4180 insertions.

2. **Finder tests** (7e923a5) — 9 unit tests with stub platform for the finder module.

3. **Roadmap completion** (334c70d) — Marked ROAD-018 and sub-items ROAD-026 through ROAD-031 as completed.

4. **21 missing commands** (50d03a6) — disk-usage, processes, uptime, thumbnail, sound, take, speed-test, renew-dhcp, locations, json-format, base64-encode/decode, password, add/remove (dock), detect, resolution-list, balance (audio). 8 global shortcuts (ip, cleanup, restart-finder, kp, vol, dark, light, search).

5. **Backward-compatible aliases** (9a45ec2) — display light-mode/toggle-mode, network info, audio input/output.

6. **USAGE.md rewrite** (ea45478) — Expanded from 369 to 740 lines covering all 88 commands.

7. **Bluetooth + apps modules** (239c0e5) — Two new modules via parallel agents: bluetooth (6 commands via blueutil), apps (6 commands: list, info, cache-clear, recent, kill, open). Plus 7 new commands in existing modules: updates, temperature, hash, timestamp, connections, dock list, screenshot record. 1450 insertions.

8. **Power module** (717f9a3) — New `mac power` module with 20 developer toggle commands (caffeinate, hidden-files, gatekeeper, function-keys, animations, etc.). Consistent on/off/toggle UX, sudo awareness, PID tracking for caffeinate. 692 lines.

9. **Power toggles design spec** (acf2321) — Saved to docs/brainstorming/.

## Goals

### [ ] 1. Distribution: GitHub Actions CI/CD (ROAD-019)

Set up GitHub Actions workflow for the Go binary:
- Build on push to master and PRs
- Cross-compile for arm64 and amd64
- Run `go test ./...` and `go vet ./...`
- Upload binaries as artifacts
- Tag-triggered releases with goreleaser or manual release workflow

Files: `.github/workflows/ci.yml`, `.github/workflows/release.yml`, `.goreleaser.yml`

### [ ] 2. Distribution: Install script and Homebrew tap (ROAD-019)

Create install script and Homebrew formula:
- `install.sh` — curl-able install script that detects arch, downloads latest release binary, installs to `/usr/local/bin/mac`
- Homebrew tap at `cosmolabs-org/homebrew-tap` with formula for `supermac`
- Shell completions generation (zsh, bash, fish) via `mac completion <shell>`

Files: `install.sh`, `Formula/supermac.rb`

### [ ] 3. Write tests for all new modules (ROAD-018 follow-up)

Current test coverage: finder (9 tests), power (7 tests). Need tests for:
- bluetooth (mock blueutil output)
- apps (mock system_profiler/mdfind)
- dock (add/remove with mock dockutil)
- system (disk-usage, processes, uptime)
- network (speed-test, connections)
- screenshot (record, thumbnail, sound)
- dev (hash, timestamp, password)
- audio (balance)

Pattern: use stub platform (see `finder_test.go`) or mock `exec.Command` output.

### [ ] 4. Prepare v0.2.0 release

- Update version in `internal/version/version.go` or ldflags
- Write CHANGELOG entry summarizing the Go rewrite
- Tag `v0.2.0` and push
- Verify goreleaser produces correct binaries
- Update README.md to reflect Go rewrite (reference `supermac-go/` as the new home)

### [ ] 5. Update roadmap and plan desktop app phase (ROAD-020)

- Mark ROAD-019 items as in-progress/completed as applicable
- Plan ROAD-020 (SwiftUI desktop app wrapping Go CLI — paid tier)
- Update `docs/prompts/2026-04-08-go-modules-and-fixes.md` status to COMPLETED
- Archive completed prompts

## Where We're Headed

The Go rewrite is done. 12 modules, 127 commands, real system calls, 40 tests. The next unlock is **distribution** — getting the binary into users' hands via Homebrew and a one-line install script. After that, the v0.2.0 release marks the transition from Bash to Go as the primary implementation.

Beyond distribution, the roadmap has two major horizons:
1. **Website** (interactive CLI demos, marketing for open-source + paid desktop model)
2. **Desktop app** (ROAD-020 — SwiftUI wrapping Go CLI, the paid tier)

The Bash version remains functional at the repo root. The Go binary lives in `supermac-go/`. Distribution will handle the naming transition.

## Priority Order

1. **CI/CD** (Goal 1) — enables everything else, catches breakage early
2. **Tests** (Goal 3) — quality gate before release
3. **Install script + Homebrew** (Goal 2) — distribution channel
4. **v0.2.0 release** (Goal 4) — milestone marker
5. **Roadmap update** (Goal 5) — planning for next phase
