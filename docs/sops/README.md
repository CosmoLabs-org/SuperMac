# SuperMac - Standard Operating Procedures

Project-specific SOPs that define how we do things in this project.

## What Are SOPs?

SOPs (Standard Operating Procedures) are executable protocols that:
- Codify project-specific workflows
- Ensure consistency across sessions
- Chain skills and tools together
- Document how this project operates

## Global vs Project SOPs

| Type | Location | Scope |
|------|----------|-------|
| **Global SOPs** | `~/.claude/skills/sop-*.md` | All CosmoLabs projects |
| **Project SOPs** | `docs/sops/` (this folder) | This project only |

## Creating Project SOPs

```bash
# Create new SOP for this project
touch docs/sops/sop-<name>.md
```

Use the template from ClaudeCodeSetup:
- `plugins/internal/sops/_template/sop.md`

## SOP Template Structure

```markdown
---
name: sop-<name>
version: 1.0.0
description: What this SOP does
category: integration|development|operations|calculation
---

# SOP Name

## Purpose
Why this exists for THIS project.

## Pre-Flight Checklist
- [ ] Requirement 1

## Phase 1: [Name]
Steps...

## Verification
How to verify success.

## Troubleshooting
Common issues.
```

## Current Project SOPs

| SOP | Category | Description |
|-----|----------|-------------|
| (none yet) | | |

## Related

- Global SOPs: `~/.claude/skills/sop-*.md`
- SOP development: `ClaudeCodeSetup/docs/instructions/sop-development.md`
- Brainstorming: `docs/brainstorming/SOPs/`
