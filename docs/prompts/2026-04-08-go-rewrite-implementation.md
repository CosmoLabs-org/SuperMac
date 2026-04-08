---
title: "SuperMac Go Rewrite — Implementation"
created: 2026-04-08
status: PENDING
priority: critical
branch: master
origin: "/continuation-prompt"
tags: [continuation, go-rewrite, implementation]
goals_total: 5
goals_completed: 0
carried_over_from: "docs/prompts/2026-04-07-supermac-audit-followup.md"
carried_over_goals: 2
related_prompts:
  - docs/prompts/2026-04-07-supermac-audit-followup.md
---

# SuperMac Go Rewrite — Implementation

## File Scope

```yaml
files_modified: []
files_created:
  - supermac-go/go.mod
  - supermac-go/cmd/mac/main.go
  - supermac-go/internal/module/module.go
  - supermac-go/internal/module/registry.go
  - supermac-go/internal/module/errors.go
  - supermac-go/internal/module/prompt.go
  - supermac-go/internal/platform/platform.go
  - supermac-go/internal/platform/darwin.go
  - supermac-go/internal/platform/mock.go
  - supermac-go/internal/config/config.go
  - supermac-go/internal/output/output.go
  - supermac-go/internal/output/colors.go
  - supermac-go/internal/version/version.go
  - supermac-go/Makefile
  - supermac-go/config.example.yaml
```

Note: The Go rewrite lives in `supermac-go/` subdirectory within the SuperMac repo initially. Once stable, it replaces the Bash version. Each module port adds files to `supermac-go/internal/modules/<name>/`.

## Context

SuperMac is being reborn. The Bash version (v0.1.1) just completed a full audit (50.5/100) and Phase 0 critical fixes — 9 bugs fixed, 5,154 lines of dead code removed, project hygiene restored. The audit proved the concept works but Bash has hard limits (no testing, no structured output, Bash 3.2 compatibility hell, god files).

Now we're rewriting the entire CLI in Go. The design spec (`docs/brainstorming/2026-04-08-go-rewrite-design.md`) has been approved and reviewed — 4 critical issues were caught and fixed. The implementation plan (`docs/planning-mode/2026-04-08-go-rewrite-implementation.md`) is ready with a 3-phase approach: Foundation → Module Port → Distribution.

This session executes the plan. The Go CLI will be the engine that powers everything — the open-source CLI and eventually a paid desktop app. Same `mac <category> <action>` interface, but with `--json`, completions, auto-update, real config, and 80%+ test coverage.

## What Got Done (Previous Session)

- **Audit Phase 0 complete**: 9 bugs fixed (BUG-001 through BUG-009), version corrected to 0.1.1
- **Root duplicates deleted**: 14 files, 5,154 lines of dead code removed
- **Project hygiene**: MIT LICENSE, .gitignore, project CLAUDE.md created
- **Roadmap initialized**: 31 items across Bash phases + Go rewrite + desktop app
- **6 feature issues filed**: FEAT-001 through FEAT-006 (Cobra interface, config, output, auto-update, completions, distribution)
- **Go design spec approved and reviewed**: 4 critical issues addressed (per-command flags, platform interface, atomic update, sudo design)
- **Implementation plan written**: 3 phases, dependency graph, parallelizable module ports
- **Branch renamed**: `main` → `master` (per CosmoLabs standard)

## Required Reading

Read these BEFORE starting implementation:

1. `docs/brainstorming/2026-04-08-go-rewrite-design.md` — The full design spec (~200 lines). Defines interfaces, architecture, all decisions.
2. `docs/planning-mode/2026-04-08-go-rewrite-implementation.md` — Step-by-step implementation plan with code samples.

Do NOT reference the Bash source as a design guide — the design spec is the authority. The Bash code (`lib/*.sh`) is only for understanding what each command does functionally.

## Goals

### [ ] 1. Scaffold Go project (ROAD-021)

Create `supermac-go/` directory with full project structure:

```
supermac-go/
  cmd/mac/main.go          — Cobra root command with persistent flags
  internal/module/         — Module interface + registry (from design spec)
  internal/platform/       — Platform interface + darwin impl + mock
  internal/config/         — YAML config loader
  internal/output/         — Writer interface (colored/json/quiet)
  internal/version/        — Single source of truth (ldflags)
  go.mod                   — github.com/cosmolabs-org/supermac
  Makefile                 — build, test, lint, install
  config.example.yaml
```

Dependencies: `cobra`, `yaml.v3`, `semver`, `google/uuid`

Verification: `make build && ./mac --version` prints version. `./mac help` shows category list (empty stubs ok).

**Key detail from spec review**: The `Command` struct must include `Flags []Flag` (per-command flags like `--sort`, `--force`). The `Context` struct must include `Platform platform.Interface` and `Prompt PromptInterface`. See design spec Module Interface section.

### [ ] 2. Implement module interface + registry (ROAD-022)

Implement the full Module interface from the design spec:
- `Module` interface: `Name()`, `ShortDescription()`, `Emoji()`, `Commands()`, `Search(term string)`
- `Command` struct with `Flags []Flag` (critical — not just positional args)
- `Context` struct with `Platform`, `Prompt`, `Flags map[string]string`
- `Registry` with `Register()`, `All()`, `Get()` — package-level `var modules = make(map[string]Module)` (NOT in init())
- Wire into Cobra: iterate `module.All()`, create subcommand per module, nested subcommand per command

**Key detail**: Registry must use package-level var, NOT init() function, to avoid race condition with module init() registrations.

### [ ] 3. Implement platform abstraction layer (ROAD-023)

This is the test seam. Without it, no testing is possible.

- Define `platform.Interface` with all ~25 methods from design spec
- Implement `platform.DarwinPlatform{}` with real `exec.Command` calls
- Implement `platform.MockPlatform{}` with configurable responses
- Key methods: `RunOSAScript`, `ReadDefault/WriteDefault`, `SetWiFi`, `GetWiFiStatus`, `FlushDNS` (sudo), `GetMemoryInfo`, `GetCPUInfo`, `GetBatteryInfo`, `SetBrightness`, `GetVolume`, `ListProcesses`, `KillPort`, `RunSudoCommand`

**Key detail for sudo**: `RunSudoCommand` checks passwordless sudo first (`sudo -n true`), falls back to prompting via native `sudo`. In `--dry-run` mode, prints the command and returns success. Mock records calls without executing.

### [ ] 4. Implement config + output systems (ROAD-024, ROAD-025)

**Config** (`internal/config/`):
- Load `~/.supermac/config.yaml`, create with defaults if missing
- Migration: if old `config.json` exists, convert to YAML, backup old file
- `mac config edit/get/set/list` commands
- Aliases from config converted to Cobra aliases at registration

**Output** (`internal/output/`):
- `Writer` interface: Info, Success, Warning, Error, Header, Table, JSON
- `NewColoredWriter` — ANSI colors + icons, respects NO_COLOR and non-TTY
- `NewJSONWriter` — structured JSON for all output
- `NewQuietWriter` — errors only
- Auto-selection based on `--json`, `--quiet`, `NO_COLOR`, TTY detection

**Also implement** (from spec review):
- `PromptInterface` with `Confirm()`, `Input()`, `Select()` — mockable for tests
- `ExitError` type with codes 0-5 for machine-readable errors in --json mode

### [ ] 5. Write initial tests and verify foundation works

Test the foundation before starting module ports:
- `module/registry_test.go` — register, get, all
- `config/config_test.go` — load, save, defaults, migration
- `output/output_test.go` — colored, json, quiet modes
- `platform/mock_test.go` — mock returns expected values

Verification commands:
```
cd supermac-go && make build
./mac help                      # Shows module categories
./mac version                   # Shows version, arch, macOS info
./mac config list               # Shows default config
./mac config set output.format json  # Sets a config value
./mac config get output.format       # Returns "json"
make test                       # All tests pass
```

## Carry-Over Tasks

From the audit followup session (docs/prompts/2026-04-07-supermac-audit-followup.md):
- [ ] Run `/upgrade-docs` to bring README and docs up to post-fix reality
- [ ] Consider next release version bump (from 0.1.1)

## Carry-Overs

1. **[MEDIUM] SuperMac v2.1.0 Audit Follow-up** (9/11 goals complete)
   → `docs/prompts/2026-04-07-supermac-audit-followup.md`
   Remaining: `/upgrade-docs` and release bump. Lower priority than Go rewrite.

## Where We're Headed

**The big picture**: SuperMac is transitioning from a Bash prototype to a professional Go CLI that will power an open-source + paid desktop app ecosystem. This session builds the foundation.

After this session's 5 goals are done, the next session tackles module porting (ROAD-018) — 10 Bash modules ported to Go, 6 of which can run as parallel agents in worktrees. The Bash version becomes the reference implementation while the Go version achieves feature parity.

**Roadmap position**:
```
✅ Phase 0: Bash critical fixes (ROAD-001) — DONE
▶️  Go Rewrite: Foundation (ROAD-017) — THIS SESSION
⬜ Go Rewrite: Module Port (ROAD-018) — Next session, parallelizable
⬜ Go Rewrite: Distribution (ROAD-019) — After modules
⬜ Desktop Application (ROAD-020) — Future
```

**Business model reminder**: CLI is open-source (MIT). Desktop app is freemium ($4.99/mo premium features). Both powered by the same Go engine.

## Priority Order

1. **Scaffold + module interface** (Goals 1-2) — everything depends on this
2. **Platform abstraction** (Goal 3) — testability depends on this
3. **Config + output** (Goal 4) — user-facing features
4. **Tests + verification** (Goal 5) — validates foundation before module porting

## Working Conventions

- **Branch**: `master` (not main — CosmoLabs standard)
- **Author**: `GΛB <Gab@CosmoLabs.org>`
- **No AI attribution** in commits
- **CCS for commits**: Use `/commit-all`, never raw `git commit`
- **Go code**: `gofmt`, `golangci-lint`, standard Go project layout
- **Design spec is authority**: When in doubt, follow the spec. The Bash code is only a functional reference.
- **No Bash scripts ever again** — Go from here on
