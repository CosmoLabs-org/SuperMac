---
status: PENDING
type: audit-followup
priority: high
source: audit/2026-04-07-supermac
---

# Audit Followup: SuperMac 2026-04-07

## Overall Score: 50.5/100 (Developing)

## Critical Bugs to Fix First

1. **rm -rf /tmp/*** in system cleanup — system.sh:173. Replace with targeted find/delete.
2. **Airport binary path broken** — wifi.sh:174. Use full path `/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport`.
3. **Wrong memory on Apple Silicon** — system.sh hardcodes page_size=4096. Apple Silicon uses 16384. Use `sysctl -n hw.pagesize`.
4. **Install missing 4 modules** — install.sh:154-162. Add wifi.sh, dock.sh, audio.sh, screenshot.sh.
5. **dock.sh:131 string bug** — `"Dock auto-hide $action"d!` → `"Dock auto-hide ${action}d!"`
6. **dev.sh:281 decimal comparison** — String comparison used for CPU threshold. Use bc -l.
7. **No LICENSE file** — Create MIT LICENSE.
8. **config.json dead code** — Either connect get_config() to modules or remove.

## Quick Wins (small effort, high impact)

1. Delete all 14 root-level duplicate files — 5,154 lines of dead code
2. Fix .gitignore — Add .DS_Store, dist/, *.tar.gz, *.log
3. Add missing modules to tests/test.sh arrays — dock, audio, screenshot
4. Fix README accuracy — Remove non-existent commands, add screenshot module
5. Add missing modules to install.sh download list

## Roadmap Items to Start

1. Make config.json functional or remove it entirely
2. Add shellcheck to CI with .shellcheckrc
3. Add NO_COLOR support and --no-color flag
4. Split utils.sh into focused files (colors, output, validation, etc.)
5. Add module auto-discovery (replace hardcoded CATEGORIES)
6. Add --yes flag for non-interactive use
7. Add SHA256 checksums to installer

## Files to Read First

- docs/audit/latest/brief.md — Project context (3.2KB, load instead of re-auditing)
- docs/audit/latest/risk-map.md — Where NOT to touch without tests
- docs/audit/latest/upgrade-plan.md — Phased improvement plan
- docs/audit/latest/README.md — Executive summary with scorecard
