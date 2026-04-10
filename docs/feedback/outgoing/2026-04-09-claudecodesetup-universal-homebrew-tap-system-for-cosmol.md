---
id: FB-398
title: Universal Homebrew tap system for CosmoLabs open-source projects
type: feature
status: pending
priority: high
complexity: complex
from_project: SuperMac
from_path: /Users/gab/PROJECTS/SuperMac
to_project: ClaudeCodeSetup
to_target: project
created: "2026-04-09T08:24:10.202384-03:00"
updated: "2026-04-09T08:24:10.202384-03:00"
suggested_conversion: feature
converted_to: null
related_issues: []
brainstorm_ref: null
suggested_workflow:
  - brainstorming
  - plan-mode
  - implementation
response:
  acknowledged: null
  acknowledged_by: null
  started: null
  implemented: null
  rejected: null
  rejection_reason: null
  notes: ""
---

# FB-398: Universal Homebrew tap system for CosmoLabs open-source projects

## What
CosmoLabs needs a universal Homebrew tap repo (cosmolabs-org/homebrew-tap) that serves as the single distribution channel for all open-source Go/CLI tools. This should be automated and integrated with CCS.

## Why
- SuperMac is the first of many CosmoLabs open-source tools needing Homebrew distribution
- Centralized tap: users run brew tap cosmolabs-org/tap once, get access to everything
- Industry standard: hashicorp/tap, grafana/tap, bufbuild/tap, cloudflare/tap all do this
- When org renames from cosmolabs-org to cosmolabs (pending trademark), one repo to update

## Proposed Solution
1. New GitHub repo cosmolabs-org/homebrew-tap with Formula/ directory
2. Auto-update GitHub Action triggered by release events from any cosmolabs-org/* repo — downloads tarball, computes sha256, updates formula
3. CCS commands: ccs brew formula <repo>, ccs brew bump <repo> <version>, ccs brew list
4. Formula generator template for new projects
5. Integration with CCS Distribution SOP (when it exists) and CosmoLabs OpenSource system (when implemented)

### Auto-update Action Flow
Any cosmolabs-org repo publishes release (tag v*)
  -> GitHub Action in homebrew-tap triggers
  -> Downloads tarball, computes sha256
  -> Updates Formula/<tool>.rb with new version + hash
  -> Commits + pushes (or opens PR for review)

### Org Rename Plan
When cosmolabs-org becomes cosmolabs:
- GitHub creates permanent redirects automatically
- Update formula download URLs (cosmolabs-org -> cosmolabs)
- Update go.mod paths across all repos
- Users re-tap: brew untap cosmolabs-org/tap && brew tap cosmolabs/tap
- Document migration SOP for smooth transition

## Integration Points
- CCS Distribution SOP — Homebrew tap as a standard distribution channel
- CosmoLabs OpenSource — tap formulas auto-generated when projects go open source
- SuperMac — pilot implementation, validates the workflow end-to-end
- Future projects — zero-friction onboarding: ccs brew formula <new-project>

## Suggested Workflow

1. brainstorming
2. plan-mode
3. implementation

