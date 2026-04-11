---
id: IDEA-MNTQLPBN
title: ccs homebrew publish — auto-generate and update Homebrew formulas in cosmolabs-org/homebrew-tap
created: "2026-04-10T23:50:05.555074-03:00"
status: seed
source: human
origin:
    session: 2027
tags:
    - distribution
    - homebrew
    - automation
    - ccs-command
---

# ccs homebrew publish — auto-generate and update Homebrew formulas in cosmolabs-org/homebrew-tap

Standardize Homebrew distribution across all CosmoLabs projects. Command reads release artifacts from GitHub, generates or updates the formula in `cosmolabs-org/homebrew-tap` with real sha256 hashes, and pushes. Supports Go binaries, Rust CLIs, and JS CLI tools. Every CosmoLabs project that distributes a binary should use this.

## Proposed Commands

1. **`ccs homebrew init`** — scaffold `Formula/xxx.rb` in the tap repo for the current project
2. **`ccs homebrew publish`** — after a GitHub release, download artifacts, compute sha256, update formula, commit and push
3. **`ccs homebrew verify`** — `brew install` and smoke-test the formula

## Integration

Should integrate with `release.yml` as a post-release step or run locally after tagging. Could also be triggered by `ccs release` as a final step.

## Pattern Established

SuperMac v0.2.0 — `cosmolabs-org/homebrew-tap` with `Formula/supermac.rb` using per-architecture sha256 hashes from GitHub Release artifacts.

