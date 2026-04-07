# Risk Assessment — SuperMac v2.1.0

## Tech Debt Hotspots

| File | Lines | Complexity | Tests | Risk |
|------|-------|------------|-------|------|
| lib/utils.sh | 589 | 8 responsibilities | 0 behavioral | HIGH - god file |
| lib/dev.sh | 610 | 2 unrelated domains | 1 smoke test | HIGH - scope creep |
| lib/system.sh | 541 | Destructive operations | 1 smoke test | CRITICAL - rm -rf |
| lib/wifi.sh | 507 | airport binary broken | 1 smoke test | HIGH - degraded features |
| lib/screenshot.sh | 590 | Predictable temp files | 0 tests | HIGH - symlink attacks |
| install.sh | 356 | Missing 4 modules | 0 tests | CRITICAL - broken install |

## Fragile Areas (high coupling, low tests)

1. **system_cleanup()** — system.sh:109-185. Single function handles cache cleanup, download pruning, trash, logs, Safari cache, /tmp deletion, font cache. One confirm() prompt for 7+ destructive operations. No --dry-run option. (see agent-3-security.md, agent-4-architecture.md)

2. **network_reset()** — network.sh:312-332. Deletes system network config plists with sudo, no backup. If interrupted, machine has no network config. (see agent-3-security.md, agent-2-core-logic.md)

3. **Module loading** — utils.sh:510-521. No re-source guard, no interface validation after load, fragile path resolution with readlink -f (not on macOS). (see agent-2-core-logic.md)

## Security Surface

### Input vectors (see agent-3-security.md)
- User CLI arguments → module dispatch functions
- Network names → wifi.sh:331 networksetup
- File paths → finder.sh:167 open -R
- App names → dock.sh:350 find command
- Hostnames → network.sh:235 ping command
- Brightness/volume values → osascript string interpolation

### Privilege escalation points (see agent-3-security.md)
- system.sh:155-158 — sudo -n to test access silently
- network.sh:199-211 — sudo for DNS flush
- network.sh:322-328 — sudo rm of network config plists
- system.sh:173-178 — sudo for system log/font cache cleanup

### Install attack surface (see agent-3-security.md, agent-6-distribution.md)
- curl|bash with no integrity verification
- Downloads from GitHub main branch (not pinned commit)
- Modifies shell RC files without explicit consent
- Prepends ~/bin to PATH (command shadowing risk)

## Chained Risks (compounding failures)

1. **install.sh missing modules + no verification = guaranteed broken commands**
   - install.sh downloads 6/10 modules (agent-6-distribution.md)
   - No checksum verification means no way to detect incomplete install
   - Users get success message but `mac wifi`, `mac dock`, `mac audio`, `mac screenshot` fail
   - Impact: 35+ commands silently broken on fresh install

2. **No tests + silent error suppression = bugs ship invisibly**
   - Only 6 of 77 commands tested (agent-5-testing.md)
   - Tests only check exit codes, never output content
   - `|| true` pattern hides all errors in system cleanup
   - Impact: bugs like wrong memory stats on Apple Silicon persist undetected

3. **Root duplicates + no CI = drift is undetectable**
   - 14 duplicate files at root (agent-1-code-quality.md)
   - No CI pipeline to verify consistency
   - Contributor edits root copy, divergence happens silently
   - Impact: installation uses lib/ but contributor tested root copy

## Single Points of Failure

1. **utils.sh** — All 9 modules depend on it. Any breaking change affects everything. Currently a 589-line god file with 8 responsibilities. (see agent-4-architecture.md)
2. **bin/mac CATEGORIES dict** — Hardcoded category registry. If lib/*.sh and CATEGORIES get out of sync, modules exist but are unreachable. (see agent-4-architecture.md)
3. **install.sh download list** — Maintained separately from actual module list. Already drifted (missing 4 modules). No validation. (see agent-6-distribution.md)

## Risk Summary Table

| Risk | Severity | Source Agents | Evidence |
|------|----------|---------------|----------|
| rm -rf /tmp/* in cleanup | CRITICAL | security, core-logic | system.sh:173 |
| No LICENSE file | CRITICAL | distribution | Missing from repo |
| Install missing 4 modules | CRITICAL | distribution, documentation | install.sh:154-162 |
| config.json dead code | CRITICAL | core-logic, architecture | get_config() never called |
| Root duplicate files | HIGH | code-quality, architecture, distribution | 14 identical files |
| curl|bash no verification | HIGH | security, distribution | install.sh |
| airport binary not found | HIGH | core-logic | wifi.sh:174 |
| Wrong Apple Silicon memory | HIGH | core-logic | system.sh page_size=4096 |
| No module re-source guard | HIGH | core-logic | utils.sh:510-521 |
| 3 modules untested | HIGH | testing | dock, audio, screenshot |
| Destructive network reset | HIGH | security, core-logic | network.sh:322-328 |
| Silent sudo operations | HIGH | security | system.sh:155 |
