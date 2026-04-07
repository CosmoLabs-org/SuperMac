# Handoffs

Session handoff documents for cross-session context preservation.

## Purpose

Handoffs capture the complete context needed to continue work in a fresh Claude Code session. Unlike ephemeral clipboard handoffs, saved handoffs provide:

- **Session Continuity** - Fresh sessions access previous work context
- **Cross-Reference** - Link to related plans, prompts, and sessions
- **Orchestration Tracking** - Document what led to GOrchestra parallel sessions

## Naming Convention

```
YYYY-MM-DD-descriptive-slug.md
```

**Examples:**
- `2026-01-07-auth-implementation.md`
- `2026-01-07-api-refactoring.md`

## Creating Handoffs

```bash
# Clipboard only (ephemeral)
ccs handoff "focus area"

# Save to docs/handoffs/
ccs handoff --save "focus area"

# Save + GOrchestra reference (for orchestration)
ccs handoff --save-plan "focus area"
```

## Template

```markdown
# Handoff: [Focus Area]

**Created**: YYYY-MM-DD HH:MM:SS
**Branch**: [branch-name]
**Version**: [version]

## Context

[What was being worked on]

## Recent Work

[Summary of commits and changes]

## Next Steps

[What should be done next]

## Key Files

- `path/to/file.go` - [description]

## Related

- [Planning Document](../planning-mode/YYYY-MM-DD-*.md)
- [Session Summary](../sessions/Session-XXX-*.md)
```

## Related

- `docs/planning-mode/` - Implementation plans from planning mode
- `docs/prompts/` - Continuation prompts for resuming work
- `docs/sessions/` - Session summaries documenting completed work
- `GOrchestra/` - Orchestration workspaces reference handoffs
