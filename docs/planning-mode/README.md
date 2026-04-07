# Planning Mode

Implementation plans and design documents. Use for detailed specifications before coding.

---

## File Format

- **Naming**: `YYYY-MM-DD-title.md`
- Start with clear goal or objective
- End with specific implementation steps

---

## Required Sections

### Header Template

```markdown
# Implementation Plan: {Title}

**Date**: YYYY-MM-DD
**Branch**: {target branch}
**Status**: Draft | Approved | In Progress | Complete

---

## References

**Brainstorming**: [../brainstorming/YYYY-MM-DD-design.md](link) ← REQUIRED if from brainstorming
**Related**: [other relevant docs]

---

## Goal

{Clear objective - what will be achieved}

---

## Implementation Steps

1. Step one
2. Step two
3. Step three

---

## Success Criteria

- [ ] Criterion one
- [ ] Criterion two
```

### Brainstorming Reference (Mandatory)

**Every plan that originated from a brainstorming session MUST reference it.**

This preserves the "how we got here" context - decisions, alternatives considered, and rationale are in the brainstorming doc; the plan focuses on execution.

```markdown
## References

**Brainstorming**: [Design Session - Feature Name](../brainstorming/2024-03-15-feature-design.md)
```

If the plan didn't come from brainstorming, note why:
```markdown
## References

**Origin**: Direct request / Bug fix / Tech debt
```

---

## Contents

- Feature specifications
- Architecture decisions
- Implementation strategies
- Technical requirements

---

## Examples

- `2024-01-15-user-authentication.md`
- `2024-02-03-api-redesign.md`

---

## Related

- [Brainstorming](../brainstorming/) - Ideas that become plans
- [Sessions](../sessions/) - Sessions that implement plans
- [Architecture](../architecture/) - System design context
