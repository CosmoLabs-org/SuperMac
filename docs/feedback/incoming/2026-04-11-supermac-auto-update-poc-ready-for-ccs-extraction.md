---
id: FB-002
title: Auto-Update PoC ready for CCS extraction
type: idea
status: pending
priority: medium
complexity: ""
from_project: SuperMac
from_path: /Users/gab/PROJECTS/SuperMac
to_project: SuperMac
to_target: project
created: "2026-04-11T17:45:05.875804-03:00"
updated: "2026-04-11T17:45:05.875804-03:00"
suggested_conversion: feature
converted_to: null
related_issues: []
brainstorm_ref: null
session: 2027
suggested_workflow: []
response:
  acknowledged: null
  acknowledged_by: null
  started: null
  implemented: null
  rejected: null
  rejection_reason: null
  notes: ""
---

# FB-002: Auto-Update PoC ready for CCS extraction

SuperMac's auto-update system (internal/update/) is now shipped and working. It's designed as a reusable CosmoLabs pattern: checker (GitHub API + cache + semver) + updater (download + verify + atomic swap + rollback). Zero deps, pure stdlib. When extracting to cosmo-go/update: add Source interface (GitHub/R2/custom), Verifier interface (SHA256/GPG), and keep the cache/swap patterns as-is. Key files: supermac-go/internal/update/checker.go, updater.go. Live test: mac update --check works against CosmoLabs-org/SuperMac GitHub releases.

