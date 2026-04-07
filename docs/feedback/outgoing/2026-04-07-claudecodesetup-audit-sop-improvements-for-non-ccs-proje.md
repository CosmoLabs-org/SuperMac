---
id: FB-386
title: Audit SOP improvements for non-CCS projects
type: improvement
status: pending
priority: high
complexity: medium
from_project: SuperMac
from_path: /Users/gab/PROJECTS/SuperMac
to_project: ClaudeCodeSetup
to_target: project
created: "2026-04-07T17:58:49.394613+02:00"
updated: "2026-04-07T17:58:49.394613+02:00"
suggested_conversion: feature
converted_to: null
related_issues: []
brainstorm_ref: null
suggested_workflow:
  - brainstorming
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

# FB-386: Audit SOP improvements for non-CCS projects

Ran /audit SOP on a standalone shell project (SuperMac) with no CCS infrastructure. Gaps identified:

1. NO-CCS FALLBACK: SOP assumes ccs audit detect/context/metrics/agents are available. No fallback for standalone projects. Should add manual context gathering section with basic shell commands (wc -l, diff, git log, file listing).

2. TEMPLATE AGENT PROMPTS: SOP references ccs audit agents for prompts but provides none inline. Should include 8-10 standard agent prompt templates in the SOP or companion file for non-CCS projects.

3. ROADMAP HEALTH CONDITIONAL: Agent should be conditional on ccs roadmap availability. SOP doesn't mark it as optional.

4. OUTPUT FILE LIST: Background synthesis agents added extra files (action-plan.md, findings.md) beyond SOP structure. Standardize or explicitly allow extensions.

5. CROSSREF FALLBACK: ccs audit crossref not available outside CCS. Add manual cross-referencing approach.

6. MINIMUM REPORT SIZE: Define minimum expected agent report size (50KB+) to catch thin reports.

7. POST-AUDIT CHECKLIST: Add verification step: file count >= 16, agent reports > 50KB, all JSON parses.

Context: Full audit produced 29 files (2.7MB), 8 agents. Project scored 50.5/100. All agent reports verbatim (123-403KB each).

## Suggested Workflow

1. brainstorming
2. implementation

