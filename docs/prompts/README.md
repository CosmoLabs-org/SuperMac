# Prompts Directory

Session handoff prompts, implementation guides, and debug continuations.

## Naming Convention

**Format**: `YYYY-MM-DD-descriptive-slug.md`

| Suffix | Purpose |
|--------|---------|
| `-continuation` | Session handoff |
| `-implementation` | Feature build guide |
| `-debug` | Troubleshooting |
| `-handoff` | Cross-model handoff |

## Template

```markdown
# [Task Name] - [Type] Prompt

**Date**: YYYY-MM-DD
**Branch**: [current branch]

## Current State
[What exists, what's working]

## Goal
[What this prompt accomplishes]

## Tasks
1. [ ] Task one
2. [ ] Task two

## Reference Files
- `path/to/file.md`
```

## Examples

- `2025-12-13-orchestra-auto-continuation.md`
- `2025-12-13-auth-system-implementation.md`

## Status Tracking (v2026.01)

The `.prompt-status.json` file tracks implementation status:

```json
{
  "prompts": {
    "YYYY-MM-DD-prompt-name.md": {
      "status": "pending|in_progress|implemented|deferred",
      "implemented_date": "YYYY-MM-DD",
      "implemented_in": "branch-name",
      "summary": "Brief description"
    }
  }
}
```

**Status values:**
| Status | Meaning |
|--------|---------|
| `pending` | Not yet started |
| `in_progress` | Currently being worked on |
| `implemented` | Completed |
| `deferred` | Postponed for later |

## Related

- [Sessions](../sessions/) - Sessions that created these prompts
- [Planning Mode](../planning-mode/) - Plans being continued
