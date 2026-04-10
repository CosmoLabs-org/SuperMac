---
title: "SuperMac Dependency System — Implementation Plan"
created: 2026-04-10
status: approved
origin: "/brainplan"
brainstorm_ref: docs/brainstorming/2026-04-10-dependency-system-auto-install.md
---

# Implementation Plan

## Step 1: Create dep package
Create `internal/dep/dep.go` with Dependency type, IsInstalled, Install, Ensure.
Create `internal/dep/dep_test.go` with unit tests.

## Step 2: Update Module interface
Add `Dependencies() []dep.Dependency` to Module interface in `internal/module/module.go`.

## Step 3: Add Dependencies() to all 12 modules
9 modules return nil. 3 modules (bluetooth, dock, audio) return their deps.

## Step 4: Remove hardcoded LookPath checks
Remove 6 exec.LookPath calls from bluetooth, dock, audio. Replace with trust in the framework.

## Step 5: Add dep checking to registerModules
In `cmd/mac/main.go`, add checkModuleDeps() before cmd.Run(ctx).

## Step 6: Add mac doctor command
In `cmd/mac/main.go`, add doctor command that iterates modules, collects deps, reports status. Add --fix flag.

## Step 7: Update tests
Run full test suite, fix any broken tests from interface change.

## Execution: Sequential (each step depends on previous)
