# Brainstorming

Ideas and exploration. Capture innovative ideas before they're lost.

---

## Document Types

### Type 1: Quick Ideas

Rapid capture of concepts and possibilities.

- **Naming**: `topic-ideas.md` or `feature-possibilities.md`
- Bullet-point format, capture everything
- No idea is too small or ambitious
- Document the "why" behind each idea

**Examples**:
- `feature-ideas.md`
- `api-v2-possibilities.md`
- `performance-optimization-ideas.md`

### Type 2: Collaborative Design Sessions

Full dialogue format preserving the Q&A and decision-making process.

- **Naming**: `YYYY-MM-DD-topic-design.md`
- Captures back-and-forth that refines ideas
- Preserves options explored with tradeoffs
- Documents decisions made and WHY
- **Template**: Use `brainstorming-collaborative.md`

**Required Sections**:
- **Goal**: What we're solving
- **Discovery & Questions**: Q&A that gathered context
- **Design Decisions**: What we chose and rationale
- **Resulting Design**: The design outcome
- **Next Steps**: Links to implementation plans

**Examples**:
- `2024-03-15-authentication-redesign.md`
- `2024-06-20-api-versioning-strategy.md`

---

## Guidelines

- **Capture everything first**, filter later
- Include technical considerations
- Break large ideas into components
- Mark status: concept, planned, in-progress, implemented
- **Always link** to resulting planning-mode docs when ideas become plans

---

## Cross-Referencing

When brainstorming becomes a plan:

```markdown
## Next Steps

- Implementation Plan: [../planning-mode/YYYY-MM-DD-feature.md]
```

When plans reference their origin (in planning-mode docs):

```markdown
## References

**Brainstorming**: [../brainstorming/YYYY-MM-DD-design.md]
```

---

## Related

- [Planning Mode](../planning-mode/) - Plans that come from brainstorming
- [New Features](../newfeatures/) - Feature specifications
- [Roadmap](../roadmap/) - Long-term direction
- [Collaborative Template](./brainstorming-collaborative.md) - Full dialogue format
