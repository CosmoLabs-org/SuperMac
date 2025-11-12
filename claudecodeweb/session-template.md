# Claude Code Session Template

Use this template to document each Claude Code session for the SuperMac project.

---

## Session Information

**Date:** YYYY-MM-DD
**Session ID:** [Generate unique ID]
**Duration:** HH:MM
**Claude Model:** Sonnet 4.5
**Session Type:** [Code Review / Feature Development / Bug Fix / Refactoring / Documentation]

---

## Session Goals

### Primary Objectives
1. [Primary goal 1]
2. [Primary goal 2]
3. [Primary goal 3]

### Secondary Objectives
- [Secondary goal 1]
- [Secondary goal 2]

---

## Work Completed

### Code Changes

#### Files Modified
- `path/to/file1.sh` - [Brief description of changes]
- `path/to/file2.sh` - [Brief description of changes]

#### Files Created
- `path/to/newfile.sh` - [Purpose and description]

#### Files Deleted
- `path/to/oldfile.sh` - [Reason for deletion]

### Features Implemented
1. **[Feature Name]**
   - Description: [What it does]
   - Location: `path/to/file:line`
   - Testing: [How it was tested]
   - Status: ✅ Complete / ⚠️ Partial / ❌ Incomplete

### Bugs Fixed
1. **[Bug Description]**
   - Issue: [What was wrong]
   - Root Cause: [Why it happened]
   - Solution: [How it was fixed]
   - Location: `path/to/file:line`
   - Status: ✅ Fixed / 🔍 Needs Verification

### Refactoring Done
1. **[Refactoring Name]**
   - Before: [Original approach]
   - After: [New approach]
   - Benefits: [Why it's better]
   - Impact: [What changed]

---

## Technical Decisions

### Architecture Decisions
1. **[Decision Topic]**
   - Context: [Why we needed to make a decision]
   - Options Considered:
     - Option A: [Pros/Cons]
     - Option B: [Pros/Cons]
   - Decision: [What we chose]
   - Rationale: [Why we chose it]

### Implementation Choices
1. **[Choice Description]**
   - Approach: [How we implemented it]
   - Alternatives: [Other ways we could have done it]
   - Trade-offs: [What we sacrificed/gained]

---

## Issues Discovered

### New Bugs Found
1. **[Bug Description]**
   - Severity: 🔴 Critical / 🟡 High / 🟢 Medium / 🔵 Low
   - Location: `path/to/file:line`
   - Impact: [What breaks]
   - Workaround: [Temporary solution if any]
   - Status: ⏳ Pending / 🔧 In Progress / ✅ Fixed

### Technical Debt Identified
1. **[Debt Description]**
   - Type: [Code Smell / Duplication / Performance / Security]
   - Location: [Where it exists]
   - Impact: [Why it matters]
   - Estimated Effort: [Time to fix]
   - Priority: [When to address]

### Blockers Encountered
1. **[Blocker Description]**
   - Issue: [What blocked progress]
   - Impact: [What couldn't be completed]
   - Resolution: [How it was resolved or needs to be]
   - Status: ⏳ Blocked / 🔓 Unblocked

---

## Testing & Validation

### Tests Added
- [ ] Unit tests for [component]
- [ ] Integration tests for [feature]
- [ ] Manual testing performed

### Test Results
```bash
# Command run
mac test command

# Results
✅ Pass: X tests
❌ Fail: Y tests
⚠️ Skip: Z tests
```

### Validation Steps
1. [Step 1 - what was tested]
   - Result: ✅ Pass / ❌ Fail
   - Notes: [Any observations]

2. [Step 2 - what was tested]
   - Result: ✅ Pass / ❌ Fail
   - Notes: [Any observations]

---

## Documentation Updates

### Documentation Added
- `path/to/doc.md` - [What it documents]

### Documentation Updated
- `README.md` - [What sections changed]
- `docs/DEVELOPMENT.md` - [What was updated]

### Documentation Needed
- [ ] [Topic that needs documentation]
- [ ] [Another topic]

---

## Code Quality Metrics

### Lines of Code
- Added: XXX lines
- Removed: YYY lines
- Modified: ZZZ lines
- Net Change: ±NNN lines

### Complexity
- Functions Added: X
- Average Function Length: Y lines
- Cyclomatic Complexity: [High/Medium/Low]

### Code Quality Checks
- [ ] shellcheck passed
- [ ] Syntax validation passed
- [ ] Naming conventions followed
- [ ] Error handling added
- [ ] Input validation added

---

## Performance Impact

### Benchmarks
```bash
# Before
Startup time: X.XXs
Command execution: Y.YYs

# After
Startup time: X.XXs (±Z%)
Command execution: Y.YYs (±Z%)
```

### Optimizations Made
1. [Optimization description]
   - Impact: [Performance improvement]
   - Trade-off: [Any downsides]

---

## Security Considerations

### Security Improvements
- [What security issues were addressed]

### Security Concerns
- [Any new security considerations introduced]

### Vulnerabilities Fixed
1. **[Vulnerability Type]**
   - Issue: [What was vulnerable]
   - Fix: [How it was secured]
   - Verification: [How we know it's fixed]

---

## Dependencies & Compatibility

### Dependencies Added
- `dependency-name` - [Why it was added]

### Dependencies Removed
- `dependency-name` - [Why it was removed]

### Compatibility Notes
- macOS Versions: [Tested on X.X, X.X]
- Bash Versions: [Tested with bash X.X]
- Breaking Changes: [Yes/No - describe if yes]

---

## Git Activity

### Commits Made
```bash
git log --oneline --since="YYYY-MM-DD" --until="YYYY-MM-DD"
```

### Branches
- Working Branch: `branch-name`
- Merged To: `main` / `develop` / [other]
- Merge Status: ⏳ Pending / ✅ Merged

### Pull Requests
- PR #XXX: [Title and description]
- Status: 🔄 Draft / 👀 Review / ✅ Merged

---

## Lessons Learned

### What Went Well
1. [Positive outcome or approach]
2. [Something that worked effectively]

### What Could Be Improved
1. [Challenge faced or inefficiency]
2. [Area for improvement]

### Key Insights
1. [Important learning or realization]
2. [Technical insight gained]

---

## Next Steps

### Immediate Actions
- [ ] [Action item 1] - Assignee: [Name/Self]
- [ ] [Action item 2] - Assignee: [Name/Self]
- [ ] [Action item 3] - Assignee: [Name/Self]

### Short-term Goals (Next Session)
1. [Goal for next session]
2. [Another goal]

### Long-term Considerations
- [Strategic consideration for future]
- [Technical debt to address later]

---

## Resources & References

### Documentation Consulted
- [URL or file reference]
- [Another reference]

### Tools Used
- Claude Code Web
- [Other tool]

### External References
- [Stack Overflow link]
- [GitHub issue/discussion]
- [Apple developer documentation]

---

## Session Notes

### Challenges
[Describe any significant challenges encountered during the session]

### Discoveries
[Note any surprising findings or important discoveries]

### Questions for Future
1. [Question that needs investigation]
2. [Another question]

### Random Notes
- [Any other notes worth recording]

---

## Session Metrics

### Productivity
- Features Completed: X
- Bugs Fixed: Y
- Tests Added: Z
- Documentation Pages: N

### Code Quality
- Code Coverage: X%
- Static Analysis Score: Y/10
- Technical Debt Added: [High/Medium/Low]
- Technical Debt Removed: [High/Medium/Low]

---

## Approval & Sign-off

### Reviewer Notes
[Space for code reviewer comments]

### Approval Status
- [ ] Code review completed
- [ ] Tests passed
- [ ] Documentation updated
- [ ] Ready to merge

**Reviewed By:** [Name]
**Date:** YYYY-MM-DD

---

## Appendix

### Code Snippets

#### Before
```bash
# Original code
```

#### After
```bash
# Refactored code
```

### Screenshots
[If applicable, note any screenshots taken]

### Error Messages
```
[Any important error messages encountered]
```

---

**Session End Time:** HH:MM
**Total Session Duration:** X hours Y minutes
**Next Session Scheduled:** YYYY-MM-DD

---

*Template Version: 1.0*
*Last Updated: 2025-11-12*
