# Project Lessons

Persistent learnings captured during development sessions. These insights help Claude grow smarter about this specific project over time.

## Purpose

This folder stores lessons learned from debugging, architectural decisions, and recurring patterns. Unlike session summaries (which document what happened), lessons capture **insights worth remembering** for future development.

## File Format

- **Naming**: `LESSON-NNN-category-description.md`
- **Index**: All lessons summarized in `LESSONS.md`
- **Categories**: Emerge from lesson topics (api, swiftui, git, state, testing, etc.)

## When to Add a Lesson

- Bug that took significant debugging
- Pattern that should be followed or avoided
- Workaround for a recurring issue
- Architectural insight worth remembering
- Performance optimization discovered
- Integration quirk to remember

## Template

```markdown
# LESSON-NNN: Title

**Category**: [category]
**Learned**: YYYY-MM-DD
**Session**: [link to originating session]

## The Lesson
[Core insight in 1-3 sentences]

## Context
[What happened that led to this learning]

## Prevention/Solution
[How to avoid or handle this in future]
```

## Examples

- `LESSON-001-api-rate-limiting.md` - Exponential backoff patterns
- `LESSON-002-swiftui-state-management.md` - Singleton AppState pattern
- `LESSON-003-git-worktree-isolation.md` - Parallel development strategy

## Workflow

1. At end of session, `/session-summary` prompts for lessons
2. If lesson identified, creates `LESSON-NNN-*.md` file
3. Updates `LESSONS.md` index with entry under category
4. Recurring lessons get count updates, not duplicates

## Related

- [Sessions](../sessions/) - Where lessons originate
- [Feedback](../feedback/) - Issues that may become lessons
- [Bugfixes](../bugfixes/) - Fixes that inform lessons
