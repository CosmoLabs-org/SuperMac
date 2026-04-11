---
title: "SuperMac — Auto-Update Implementation"
created: "2026-04-11"
status: PENDING
priority: high
branch: master
origin: "/brainplan"
tags: [continuation, implementation, auto-update, distribution]
goals_total: 9
goals_completed: 0
related_prompts:
  - docs/brainstorming/2026-04-11-auto-update-via-github-releases.md
  - docs/planning-mode/2026-04-11-auto-update-via-github-releases.md
brainstorm_ref: docs/brainstorming/2026-04-11-auto-update-via-github-releases.md
plan_ref: docs/planning-mode/2026-04-11-auto-update-via-github-releases.md
glm_tasks_ref: docs/prompts/2026-04-11-auto-update-glm-tasks.yaml
---

# SuperMac — Auto-Update Implementation

## Context

Design approved for a self-update system: background check on launch (2s timeout, 24h cache), `mac update` command with atomic binary swap, SHA256 verification, and one-version rollback. All stdlib, zero new deps. This is the first CosmoLabs auto-update implementation — after shipping, send detailed PoC feedback to CCS for universal release system abstraction.

Design spec: `docs/brainstorming/2026-04-11-auto-update-via-github-releases.md`
Implementation plan: `docs/planning-mode/2026-04-11-auto-update-via-github-releases.md`

## Goals

- [ ] 1. Fix config list display bug (wrong field for Updates)
- [ ] 2. Add `mac version --raw` flag
- [ ] 3. Build checker cache layer (read/write/TTL)
- [ ] 4. Build checker GitHub API client + semver compare
- [ ] 5. Build updater download + verify + extract
- [ ] 6. Build updater atomic swap + rollback
- [ ] 7. Wire into main.go — update command + pre-run check
- [ ] 8. Update version output with update status
- [ ] 9. Integration test + push + tag v0.3.0

## Execution Strategy

Sequential execution — each task builds on the previous. Tasks 1-2 are prerequisites. Tasks 3-4 are checker. Tasks 5-6 are updater. Task 7 wires everything together.

### Parallel GLM Dispatch (Recommended)

Tasks 1-3 are independent and can run in parallel via GLM agents. Task 4 (wiring) depends on tasks 2+3 completing first.

```bash
# Dispatch all independent tasks in parallel:
ccs glm-agent exec-batch docs/prompts/2026-04-11-auto-update-glm-tasks.yaml

# Or dispatch individually:
ccs glm-agent exec "Fix config bug + version --raw" --task-id 1
ccs glm-agent exec "Build checker package" --task-id 2
ccs glm-agent exec "Build updater package" --task-id 3

# Then after 1-3 complete, Opus reviews and runs task 4 (wiring):
# Review all worktrees, merge, then execute Task 7-9 in main session
```

**Task dependency graph:**
```
Task 1 (prereqs) ──┐
Task 2 (checker) ──┼── Task 4 (wiring) ── Task 8 (version output) ── Task 9 (integration)
Task 3 (updater) ──┘
```

## File Scope

```yaml
files_created:
  - supermac-go/internal/update/checker.go
  - supermac-go/internal/update/checker_test.go
  - supermac-go/internal/update/updater.go
  - supermac-go/internal/update/updater_test.go
files_modified:
  - supermac-go/cmd/mac/main.go
```
