---
project: SuperMac
version: 0.1.2
date: 2026-04-08
previous: 0.1.1
slug: features-and-fixes
title: "features-and-fixes Release"
---

# SuperMac v0.1.2 Release Notes

**Release Date**: April 8, 2026

**Previous**: v0.1.1

## Overview

This release brings 3 new features, 6 bug fixes, and 1 improvement.

## Highlights

Phase 0 audit fixes complete: 9 bugs fixed (BUG-001 through BUG-009), 14 root duplicate files removed (5,154 lines), version corrected to 0.1.1. Project hygiene restored with MIT LICENSE, .gitignore, CLAUDE.md, and 31-item roadmap initialized. Go rewrite design spec approved and implementation plan ready for next session.

## What's New

- Add MIT LICENSE, .gitignore, project CLAUDE.md
- Initialize roadmap with 31 items across 5 phases
- Add Go rewrite design spec and implementation plan

## Bug Fixes

- Apply 6 critical bug fixes from audit (BUG-001 through BUG-006)
- Replace declare -A with Bash 3.2-compatible lookup (BUG-009)
- Add 4 missing modules to install script (BUG-004)
- replace declare -A with Bash 3.2-compatible lookup (commit:18efc0d1)
- add 4 missing modules and deduplicate version (commit:c9c08b83)
- apply 6 critical bug fixes from audit (commit:7cf780fc)

## Removed

- Remove 14 root-level duplicate files (5,154 lines)

## Breaking Changes

> _None in this release_

## Upgrade Instructions

No breaking changes in this release. Standard upgrade applies.

## Stats

| Metric | Value |
|--------|-------|
| Commits | 12 |
| Files changed | 79 |
| New features | 3 |
| Bug fixes | 6 |

---
_Full changelog: CHANGELOG.md_
