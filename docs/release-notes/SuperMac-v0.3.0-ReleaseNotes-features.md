---
project: SuperMac
version: 0.3.0
date: 2026-04-11
previous: 0.2.0
slug: features
title: "features Release"
---

# SuperMac v0.3.0 Release Notes

**Release Date**: April 11, 2026

**Previous**: v0.2.0

## Overview

This release brings 1 new feature.

## Highlights

Auto-update system via GitHub Releases with SHA256 verification, atomic binary swap, and 24h cache. Includes checker package (GitHub API client, semver compare), updater package (download, verify, extract, swap, rollback), mac update command with --check/--rollback flags, version --raw for machine-parseable output, and pre-launch update notification. Zero new dependencies — pure Go stdlib.

## What's New

- add self-update system via GitHub Releases (commit:76acf310)

## Breaking Changes

> _None in this release_

## Upgrade Instructions

No breaking changes in this release. Standard upgrade applies.

## Stats

| Metric | Value |
|--------|-------|
| Commits | 15 |
| Files changed | 47 |
| New features | 1 |

---
_Full changelog: CHANGELOG.md_
