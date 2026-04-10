---
title: "SuperMac — Release Finalization and Distribution Validation"
created: 2026-04-10
status: PENDING
priority: high
branch: master
origin: "/continuation-prompt"
tags: [continuation, distribution, release, ci-cd]
goals_total: 7
goals_completed: 0
related_prompts:
  - docs/prompts/2026-04-08-go-modules-complete-distribution-next.md
---

# SuperMac — Release Finalization and Distribution Validation

## File Scope

```yaml
files_modified:
  - Formula/supermac.rb                          # sha256 placeholder → real hash
  - install.sh                                   # test end-to-end, fix issues
  - .github/workflows/ci.yml                     # verify on first push
  - .github/workflows/release.yml                # verify on v0.2.0 tag
  - supermac-go/docs/USAGE.md                    # add mac doctor section
  - docs/roadmap/items/ROAD-019.yaml             # update status
files_created: []                                # tap repo created externally
```

## Context

All 5 goals from the previous prompt (`2026-04-08-go-modules-complete-distribution-next.md`) are DONE. The Go rewrite is feature-complete with 12 modules, 127 commands, and 98 tests all passing. Four commits landed on master this session covering distribution infrastructure (CI/CD, install script, Homebrew formula), the declarative dependency system with `mac doctor`, 52 new unit tests, and documentation updates (README rewrite, INSTALL.md, CHANGELOG v0.2.0 entry).

**What's been built but NOT yet validated:**
- `.github/workflows/ci.yml` — builds on push to master/PR, cross-compiles arm64+amd64, runs test+vet
- `.github/workflows/release.yml` — tag-triggered, produces tarballs + checksums, creates GitHub Release
- `install.sh` — curl-able installer with arch detection and sha256 verification
- `Formula/supermac.rb` — Homebrew formula with placeholder sha256 (`PLACEHOLDER_UPDATE_ON_RELEASE`)
- `mac doctor` command — system health check with `--fix` flag for auto-installing missing deps
- `mac completion` command — generates zsh/bash/fish/powershell completions

**What's NOT done:**
- No push to GitHub yet — CI workflows are untested
- No `cosmolabs-org/homebrew-tap` repo exists yet
- No `v0.2.0` tag — release workflow has never fired
- Formula sha256 is a placeholder
- `install.sh` has never been tested end-to-end
- ROAD-019 status is `in_progress`

## Goals

### [ ] 1. Push to GitHub and verify CI passes

Push master to `origin` and confirm the CI workflow runs green:
- `git push origin master`
- Monitor the Actions tab: both `test` (go test + go vet) and `build` (arm64 + amd64 matrix) must pass
- Fix any CI failures (common: path issues in working-directory, Go version mismatch, test flakiness on macOS runner)
- This is the gate — nothing else proceeds until CI is green

### [ ] 2. Create cosmolabs-org/homebrew-tap repository

Set up the Homebrew tap that `brew install cosmolabs-org/tap/supermac` will use:
- Create repo at `github.com/CosmoLabs-org/homebrew-tap`
- Clone it locally, add `Formula/supermac.rb` from the SuperMac repo (keep in sync)
- The formula currently has `sha256 "PLACEHOLDER_UPDATE_ON_RELEASE"` — leave it until Goal 4
- Add a README explaining `brew install cosmolabs-org/tap/supermac`

### [ ] 3. Tag v0.2.0 and verify release workflow

Trigger the release pipeline and validate its output:
- `git tag v0.2.0 && git push origin v0.2.0`
- Monitor release.yml: test job, then build binaries, tarballs, checksums, GitHub Release
- Verify the GitHub Release page has: `mac-arm64.tar.gz`, `mac-amd64.tar.gz`, `checksums.txt`
- Download both binaries and smoke-test on local machine: `./mac --version`, `./mac doctor`, `./mac network ip`

### [ ] 4. Update Homebrew formula with real sha256

After the release artifact exists:
- Download `mac-arm64.tar.gz` from the GitHub Release
- Run `shasum -a 256 mac-arm64.tar.gz` to get the real hash
- Update `Formula/supermac.rb` in both SuperMac repo and homebrew-tap repo
- Commit and push to both repos
- Test: `brew install cosmolabs-org/tap/supermac` on a clean machine (or `brew reinstall`)

### [ ] 5. Test install.sh end-to-end

Validate the curl-able installer:
- Run `bash install.sh` on the local machine (or a clean macOS environment)
- Verify: correct arch detected, binary downloaded, checksum verified, installed to target path
- Test with missing dependencies, existing install (upgrade path), and offline scenarios
- Fix any issues found — install.sh is the primary distribution channel for non-Homebrew users

### [ ] 6. Write USAGE.md update for mac doctor command

Document the new dependency system and doctor command:
- Add `mac doctor` section to `supermac-go/docs/USAGE.md`
- Document: basic health check, `--fix` flag for auto-install, what it checks (external tool deps)
- Document the dependency system: which modules declare deps, how auto-install prompts work
- Include example output

### [ ] 7. Update ROAD-019 status

Close out the distribution roadmap item:
- Mark completed deliverables in ROAD-019.yaml
- Update status to `completed` (or note remaining items like auto-update and npm wrapper as future work)
- Update this prompt's status to `COMPLETED`
- Update the previous prompt's status to `COMPLETED`

## Priority Order

1. **Push + CI validation** (Goal 1) — unblocks everything, catches breakage immediately
2. **Tag v0.2.0 + release workflow** (Goal 3) — produces the artifacts Goals 4-5 depend on
3. **Homebrew tap + formula** (Goals 2 + 4) — primary distribution channel
4. **install.sh testing** (Goal 5) — secondary distribution channel
5. **USAGE.md doctor docs** (Goal 6) — documentation for the new feature
6. **Roadmap closure** (Goal 7) — housekeeping, do last

## Where We're Headed

This session takes SuperMac from "code complete" to "shipped." After these 7 goals, the project has:
- A green CI pipeline catching regressions on every push
- A v0.2.0 release with signed binaries on GitHub
- Two working install paths: `brew install` and `curl | bash`
- 98 tests, 127 commands, 12 modules — all documented

The next horizon beyond this session is ROAD-020 (SwiftUI desktop app wrapping the Go CLI — the paid tier). The open-source CLI will be stable and distributable, ready for users to discover via Homebrew or the install script. Future work not in this session: auto-update mechanism (`internal/update/`), npm/bun wrapper, website with interactive CLI demos.
