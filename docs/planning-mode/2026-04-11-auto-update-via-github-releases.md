---
completed: "2026-04-11"
created: "2026-04-11"
goals_completed: 44
goals_total: 44
status: COMPLETED
title: Auto-Update via GitHub Releases — Implementation Plan
---

# Auto-Update via GitHub Releases — Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add self-update capability to SuperMac that checks GitHub Releases on launch and applies updates via atomic binary swap with SHA256 verification.

**Architecture:** Two-package design — `internal/update/checker.go` handles GitHub API queries with filesystem cache, `internal/update/updater.go` handles download/verify/swap. Wired into `cmd/mac/main.go` as a new `update` subcommand and a synchronous pre-run check on every command.

**Tech Stack:** Go stdlib only — `net/http`, `archive/tar`, `compress/gzip`, `crypto/sha256`, `encoding/json`, `os`, `syscall` (for PID lock).

**Design spec:** `docs/brainstorming/2026-04-11-auto-update-via-github-releases.md`

---

## File Structure

```
CREATE:
  supermac-go/internal/update/checker.go          # GitHub Releases API client + cache (~120 lines)
  supermac-go/internal/update/checker_test.go     # 5 tests (~150 lines)
  supermac-go/internal/update/updater.go          # Download, verify, swap, rollback (~180 lines)
  supermac-go/internal/update/updater_test.go     # 5 tests (~200 lines)

MODIFY:
  supermac-go/cmd/mac/main.go                     # Add update cmd, version --raw, pre-run check (~40 lines)
```

---

## Chunk 1: Prerequisites & Checker

### Task 1: Fix config list display bug

**Files:**
- Modify: `supermac-go/cmd/mac/main.go:140`

- [ ] **Step 1: Fix the wrong field in config list output**

Line 140 of `main.go` prints `cfg.Output.Format` instead of `cfg.Updates.Check`:

```go
// Before:
fmt.Printf("  Updates:  %v (%s)\n", cfg.Output.Format, cfg.Updates.Channel)
// After:
fmt.Printf("  Updates:  %v (%s)\n", cfg.Updates.Check, cfg.Updates.Channel)
```

- [ ] **Step 2: Verify**

Run: `go build -C supermac-go ./cmd/mac && ./supermac-go/mac config list`
Expected: Updates line shows `true (stable)` not `text (stable)`

- [ ] **Step 3: Commit**

```bash
git add supermac-go/cmd/mac/main.go
git commit -m "fix(config): display Updates.Check instead of Output.Format in config list"
```

---

### Task 2: Add `mac version --raw` flag

**Files:**
- Modify: `supermac-go/cmd/mac/main.go:95-115`

- [ ] **Step 1: Write the test**

Create test in `supermac-go/cmd/mac/main_test.go` (if not exists, add to existing):

```go
func TestVersionRaw(t *testing.T) {
	cmd := exec.Command("./mac", "version", "--raw")
	cmd.Dir = filepath.Join("..", "..") // adjust to build dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Skipf("binary not built: %v", err)
	}
	// Should be just the version string, no box, no labels
	got := strings.TrimSpace(string(output))
	if got != version.Version {
		t.Errorf("version --raw = %q, want %q", got, version.Version)
	}
}
```

- [ ] **Step 2: Implement --raw flag on versionCmd**

Replace the `versionCmd()` function in `main.go`:

```go
func versionCmd() *cobra.Command {
	var raw bool
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			if raw {
				fmt.Println(version.Version)
				return
			}
			w := getOutput()
			w.Header(fmt.Sprintf("SuperMac v%s", version.Version))
			fmt.Printf("  Version:    %s\n", version.Version)
			fmt.Printf("  Build:      %s\n", version.BuildDate)
			fmt.Println()
			modules := module.All()
			names := make([]string, 0, len(modules))
			for name := range modules {
				names = append(names, name)
			}
			sort.Strings(names)
			fmt.Printf("  Modules:    %s\n", strings.Join(names, ", "))
		},
	}
	cmd.Flags().BoolVar(&raw, "raw", false, "Print version string only (machine-parseable)")
	return cmd
}
```

- [ ] **Step 3: Build and test**

Run: `go build -C supermac-go -o mac ./cmd/mac && ./supermac-go/mac version --raw`
Expected: prints just `0.2.1` (or current version)

- [ ] **Step 4: Commit**

```bash
git add supermac-go/cmd/mac/main.go
git commit -m "feat(version): add --raw flag for machine-parseable version output"
```

---

### Task 3: Build checker — cache layer

**Files:**
- Create: `supermac-go/internal/update/checker.go`
- Create: `supermac-go/internal/update/checker_test.go`

- [ ] **Step 1: Write failing tests for cache read/write**

```go
// checker_test.go
package update

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCacheRead(t *testing.T) {
	dir := t.TempDir()
	cachePath := filepath.Join(dir, "update-cache.json")

	release := &cachedRelease{
		CheckedAt:      time.Now(),
		Version:        "0.2.2",
		Tag:            "v0.2.2",
		TarballURL:     "https://example.com/mac-arm64.tar.gz",
		ChecksumsURL:   "https://example.com/checksums.txt",
		PublishedAt:    time.Now().Add(-1 * time.Hour),
		CurrentVersion: "0.2.1",
	}
	data, _ := json.Marshal(release)
	os.WriteFile(cachePath, data, 0644)

	c := &Checker{cachePath: cachePath, cacheTTL: 24 * time.Hour}
	got, err := c.readCache()
	if err != nil {
		t.Fatalf("readCache: %v", err)
	}
	if got.Version != "0.2.2" {
		t.Errorf("Version = %q, want %q", got.Version, "0.2.2")
	}
}

func TestCacheExpired(t *testing.T) {
	dir := t.TempDir()
	cachePath := filepath.Join(dir, "update-cache.json")

	release := &cachedRelease{
		CheckedAt:      time.Now().Add(-48 * time.Hour), // 2 days ago
		Version:        "0.2.2",
		CurrentVersion: "0.2.1",
	}
	data, _ := json.Marshal(release)
	os.WriteFile(cachePath, data, 0644)

	c := &Checker{cachePath: cachePath, cacheTTL: 24 * time.Hour}
	got, err := c.readCache()
	if err != nil {
		t.Fatalf("readCache: %v", err)
	}
	if got != nil {
		t.Error("expired cache should return nil")
	}
}

func TestCacheWrite(t *testing.T) {
	dir := t.TempDir()
	cachePath := filepath.Join(dir, "update-cache.json")

	c := &Checker{cachePath: cachePath, cacheTTL: 24 * time.Hour}
	release := &cachedRelease{
		CheckedAt:      time.Now(),
		Version:        "0.2.3",
		Tag:            "v0.2.3",
		TarballURL:     "https://example.com/mac-arm64.tar.gz",
		ChecksumsURL:   "https://example.com/checksums.txt",
		PublishedAt:    time.Now(),
		CurrentVersion: "0.2.2",
	}
	if err := c.writeCache(release); err != nil {
		t.Fatalf("writeCache: %v", err)
	}
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		t.Fatal("cache file not created")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test -C supermac-go ./internal/update/ -run TestCache -v`
Expected: FAIL — package doesn't exist

- [ ] **Step 3: Implement checker cache types and methods**

```go
// checker.go
package update

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

// Release represents a GitHub release with download info.
type Release struct {
	Version      string
	Tag          string
	TarballURL   string
	ChecksumsURL string
	Date         time.Time
}

// cachedRelease is the cache file format.
type cachedRelease struct {
	CheckedAt      time.Time `json:"checked_at"`
	Version        string    `json:"version"`
	Tag            string    `json:"tag"`
	TarballURL     string    `json:"tarball_url"`
	ChecksumsURL   string    `json:"checksums_url"`
	PublishedAt    time.Time `json:"published_at"`
	CurrentVersion string    `json:"current_version"`
}

// Checker queries GitHub Releases for updates.
type Checker struct {
	repo      string
	cachePath string
	cacheTTL  time.Duration
	current   string
	arch      string
	client    *http.Client
}

// NewChecker creates a Checker for the given repo and current version.
func NewChecker(repo, currentVersion string) *Checker {
	home, _ := os.UserHomeDir()
	arch := runtime.GOARCH
	if arch == "x86_64" {
		arch = "amd64"
	}
	return &Checker{
		repo:      repo,
		cachePath: fmt.Sprintf("%s/.supermac/update-cache.json", home),
		cacheTTL:  24 * time.Hour,
		current:   currentVersion,
		arch:      arch,
		client:    &http.Client{Timeout: 2 * time.Second},
	}
}

func (c *Checker) readCache() (*cachedRelease, error) {
	data, err := os.ReadFile(c.cachePath)
	if err != nil {
		return nil, err
	}
	var cr cachedRelease
	if err := json.Unmarshal(data, &cr); err != nil {
		return nil, err
	}
	if time.Since(cr.CheckedAt) > c.cacheTTL {
		return nil, nil
	}
	return &cr, nil
}

func (c *Checker) writeCache(cr *cachedRelease) error {
	dir := filepath.Dir(c.cachePath)  //nolint:staticcheck // filepath.Dir is fine
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.Marshal(cr)
	if err != nil {
		return err
	}
	return os.WriteFile(c.cachePath, data, 0644)
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test -C supermac-go ./internal/update/ -run TestCache -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add supermac-go/internal/update/
git commit -m "feat(update): add checker cache layer with TTL-based expiration"
```

---

### Task 4: Build checker — GitHub API client + semver compare

**Files:**
- Modify: `supermac-go/internal/update/checker.go`
- Modify: `supermac-go/internal/update/checker_test.go`

- [ ] **Step 1: Write failing tests for API parsing and semver**

```go
func TestParseRelease(t *testing.T) {
	apiResp := `{"tag_name":"v0.2.2","published_at":"2026-04-11T10:20:00Z","assets":[{"name":"mac-arm64.tar.gz","browser_download_url":"https://github.com/CosmoLabs-org/SuperMac/releases/download/v0.2.2/mac-arm64.tar.gz"},{"name":"checksums.txt","browser_download_url":"https://github.com/CosmoLabs-org/SuperMac/releases/download/v0.2.2/checksums.txt"}]}`

	got, err := parseRelease([]byte(apiResp), "arm64")
	if err != nil {
		t.Fatalf("parseRelease: %v", err)
	}
	if got.Version != "0.2.2" {
		t.Errorf("Version = %q, want %q", got.Version, "0.2.2")
	}
	if got.Tag != "v0.2.2" {
		t.Errorf("Tag = %q, want %q", got.Tag, "v0.2.2")
	}
	if !strings.Contains(got.TarballURL, "mac-arm64.tar.gz") {
		t.Errorf("TarballURL = %q, want arm64 tarball", got.TarballURL)
	}
}

func TestSemverCompare(t *testing.T) {
	tests := []struct{ a, b string; want bool }{
		{"0.2.2", "0.2.1", true},
		{"0.2.1", "0.2.1", false},
		{"0.3.0", "0.2.9", true},
		{"1.0.0", "0.9.9", true},
		{"0.2.1", "0.2.2", false},
	}
	for _, tt := range tests {
		got := isNewer(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("isNewer(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test -C supermac-go ./internal/update/ -run "TestParseRelease|TestSemverCompare" -v`
Expected: FAIL — functions not defined

- [ ] **Step 3: Implement parseRelease and isNewer**

Add to `checker.go`:

```go
// githubRelease matches the GitHub API response shape.
type githubRelease struct {
	TagName     string `json:"tag_name"`
	PublishedAt string `json:"published_at"`
	Assets      []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func parseRelease(data []byte, arch string) (*Release, error) {
	var gr githubRelease
	if err := json.Unmarshal(data, &gr); err != nil {
		return nil, err
	}
	version := strings.TrimPrefix(gr.TagName, "v")
	publishedAt, _ := time.Parse(time.RFC3339, gr.PublishedAt)

	var tarballURL, checksumsURL string
	for _, a := range gr.Assets {
		if a.Name == "mac-"+arch+".tar.gz" {
			tarballURL = a.BrowserDownloadURL
		}
		if a.Name == "checksums.txt" {
			checksumsURL = a.BrowserDownloadURL
		}
	}
	if tarballURL == "" {
		return nil, fmt.Errorf("no tarball found for arch %s", arch)
	}

	return &Release{
		Version:      version,
		Tag:          gr.TagName,
		TarballURL:   tarballURL,
		ChecksumsURL: checksumsURL,
		Date:         publishedAt,
	}, nil
}

func isNewer(a, b string) bool {
	var aMaj, aMin, aPatch, bMaj, bMin, bPatch int
	fmt.Sscanf(a, "%d.%d.%d", &aMaj, &aMin, &aPatch)
	fmt.Sscanf(b, "%d.%d.%d", &bMaj, &bMin, &bPatch)
	if aMaj != bMaj {
		return aMaj > bMaj
	}
	if aMin != bMin {
		return aMin > bMin
	}
	return aPatch > bPatch
}
```

- [ ] **Step 4: Run tests**

Run: `go test -C supermac-go ./internal/update/ -run "TestParseRelease|TestSemverCompare" -v`
Expected: PASS

- [ ] **Step 5: Implement CheckAvailable method**

Add to `checker.go`:

```go
// CheckAvailable returns the latest release if a newer version exists.
// Uses cache if fresh, otherwise queries GitHub API with a short timeout.
func (c *Checker) CheckAvailable(ctx context.Context) (*Release, error) {
	// Try cache first
	cached, err := c.readCache()
	if err == nil && cached != nil {
		rel := &Release{
			Version:      cached.Version,
			Tag:          cached.Tag,
			TarballURL:   cached.TarballURL,
			ChecksumsURL: cached.ChecksumsURL,
			Date:         cached.PublishedAt,
		}
		if isNewer(rel.Version, c.current) {
			return rel, nil
		}
		return nil, nil
	}

	// Query GitHub API
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", c.repo)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	data := make([]byte, 1024*64) // 64KB should be more than enough
	n, _ := resp.Body.Read(data)
	data = data[:n]

	rel, err := parseRelease(data, c.arch)
	if err != nil {
		return nil, err
	}

	// Write cache
	cr := &cachedRelease{
		CheckedAt:      time.Now(),
		Version:        rel.Version,
		Tag:            rel.Tag,
		TarballURL:     rel.TarballURL,
		ChecksumsURL:   rel.ChecksumsURL,
		PublishedAt:    rel.Date,
		CurrentVersion: c.current,
	}
	c.writeCache(cr) // best effort

	if isNewer(rel.Version, c.current) {
		return rel, nil
	}
	return nil, nil
}
```

- [ ] **Step 6: Run all checker tests**

Run: `go test -C supermac-go ./internal/update/ -v`
Expected: All PASS

- [ ] **Step 7: Commit**

```bash
git add supermac-go/internal/update/
git commit -m "feat(update): add GitHub API client with semver compare and cache"
```

---

## Chunk 2: Updater & Wiring

### Task 5: Build updater — download, verify, extract

**Files:**
- Create: `supermac-go/internal/update/updater.go`
- Create: `supermac-go/internal/update/updater_test.go`

- [ ] **Step 1: Write failing tests for verify and extract**

```go
// updater_test.go
package update

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVerifyChecksum(t *testing.T) {
	// Create a test file with known content
	dir := t.TempDir()
	content := []byte("test binary content")
	tarballPath := filepath.Join(dir, "mac-arm64.tar.gz")
	os.WriteFile(tarballPath, content, 0644)

	// Compute expected hash
	expected := fmt.Sprintf("%x", sha256.Sum256(content))
	checksums := fmt.Sprintf("%s  mac-arm64.tar.gz\n", expected)
	checksumsPath := filepath.Join(dir, "checksums.txt")
	os.WriteFile(checksumsPath, []byte(checksums), 0644)

	err := verifyChecksum(tarballPath, checksumsPath, "mac-arm64.tar.gz")
	if err != nil {
		t.Fatalf("verifyChecksum: %v", err)
	}
}

func TestVerifyChecksumMismatch(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "mac-arm64.tar.gz"), []byte("content"), 0644)
	os.WriteFile(filepath.Join(dir, "checksums.txt"), []byte("0000invalidhash  mac-arm64.tar.gz\n"), 0644)

	err := verifyChecksum(filepath.Join(dir, "mac-arm64.tar.gz"), filepath.Join(dir, "checksums.txt"), "mac-arm64.tar.gz")
	if err == nil {
		t.Fatal("expected error for checksum mismatch")
	}
}

func TestExtractBinary(t *testing.T) {
	dir := t.TempDir()
	// Create a valid tar.gz with a "mac-arm64" binary
	binaryContent := "#!/bin/sh\necho 0.2.2"
	tarballPath := filepath.Join(dir, "mac-arm64.tar.gz")
	createTestTarball(t, tarballPath, "mac-arm64", binaryContent)

	outPath := filepath.Join(dir, "mac")
	if err := extractBinary(tarballPath, "mac-arm64", outPath); err != nil {
		t.Fatalf("extractBinary: %v", err)
	}
	got, _ := os.ReadFile(outPath)
	if strings.TrimSpace(string(got)) != strings.TrimSpace(binaryContent) {
		t.Errorf("extracted = %q, want %q", string(got), binaryContent)
	}
}

func createTestTarball(t *testing.T, path, entryName, content string) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	gw := gzip.NewWriter(f)
	tw := tar.NewWriter(gw)
	tw.WriteHeader(&tar.Header{Name: entryName, Size: int64(len(content)), Mode: 0755})
	tw.Write([]byte(content))
	tw.Close()
	gw.Close()
	f.Close()
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test -C supermac-go ./internal/update/ -run "TestVerify|TestExtract" -v`
Expected: FAIL — functions not defined

- [ ] **Step 3: Implement verifyChecksum and extractBinary**

```go
// updater.go
package update

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Updater handles downloading and applying updates.
type Updater struct {
	binaryPath string
	arch       string
}

// NewUpdater creates an Updater for the current binary.
func NewUpdater() (*Updater, error) {
	bin, err := os.Executable()
	if err != nil {
		return nil, err
	}
	bin, err = filepath.EvalSymlinks(bin)
	if err != nil {
		return nil, err
	}
	arch := "arm64"
	if strings.Contains(strings.ToLower(runtime.GOARCH), "amd") {
		arch = "amd64"
	}
	return &Updater{binaryPath: bin, arch: arch}, nil
}

func verifyChecksum(tarballPath, checksumsPath, filename string) error {
	data, err := os.ReadFile(tarballPath)
	if err != nil {
		return err
	}
	actual := fmt.Sprintf("%x", sha256.Sum256(data))

	checksums, err := os.ReadFile(checksumsPath)
	if err != nil {
		return err
	}
	for _, line := range strings.Split(string(checksums), "\n") {
		parts := strings.SplitN(line, "  ", 2)
		if len(parts) == 2 && parts[1] == filename {
			if parts[0] == actual {
				return nil
			}
			return fmt.Errorf("checksum mismatch: expected %s, got %s", parts[0], actual)
		}
	}
	return fmt.Errorf("%s not found in checksums.txt", filename)
}

func extractBinary(tarballPath, entryName, outPath string) error {
	f, err := os.Open(tarballPath)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return fmt.Errorf("%s not found in tarball", entryName)
		}
		if err != nil {
			return err
		}
		if hdr.Name == entryName {
			out, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				return err
			}
			return out.Close()
		}
	}
}
```

- [ ] **Step 4: Run tests**

Run: `go test -C supermac-go ./internal/update/ -run "TestVerify|TestExtract" -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add supermac-go/internal/update/updater.go supermac-go/internal/update/updater_test.go
git commit -m "feat(update): add checksum verification and tarball extraction"
```

---

### Task 6: Build updater — atomic swap and rollback

**Files:**
- Modify: `supermac-go/internal/update/updater.go`
- Modify: `supermac-go/internal/update/updater_test.go`

- [ ] **Step 1: Write failing tests for swap and rollback**

```go
func TestAtomicSwap(t *testing.T) {
	dir := t.TempDir()
	// Create fake "current" binary
	currentBin := filepath.Join(dir, "mac")
	os.WriteFile(currentBin, []byte("old binary"), 0755)

	// Create fake "new" binary
	newBin := filepath.Join(dir, "new-mac")
	os.WriteFile(newBin, []byte("new binary"), 0755)

	err := atomicSwap(currentBin, newBin)
	if err != nil {
		t.Fatalf("atomicSwap: %v", err)
	}

	// Current binary should have new content
	got, _ := os.ReadFile(currentBin)
	if string(got) != "new binary" {
		t.Errorf("binary = %q, want %q", string(got), "new binary")
	}
	// .bak should have old content
	bak, _ := os.ReadFile(currentBin + ".bak")
	if string(bak) != "old binary" {
		t.Errorf("bak = %q, want %q", string(bak), "old binary")
	}
	// new binary should be gone
	if _, err := os.Stat(newBin); !os.IsNotExist(err) {
		t.Error("new binary should be removed after swap")
	}
}

func TestSwapFailsRestore(t *testing.T) {
	dir := t.TempDir()
	currentBin := filepath.Join(dir, "mac")
	os.WriteFile(currentBin, []byte("old binary"), 0755)

	// New binary doesn't exist — swap should fail and restore
	err := atomicSwap(currentBin, filepath.Join(dir, "nonexistent"))
	if err == nil {
		t.Fatal("expected error for missing new binary")
	}

	// Original should be restored
	got, _ := os.ReadFile(currentBin)
	if string(got) != "old binary" {
		t.Errorf("binary after failed swap = %q, want %q", string(got), "old binary")
	}
}

func TestRollback(t *testing.T) {
	dir := t.TempDir()
	currentBin := filepath.Join(dir, "mac")
	os.WriteFile(currentBin, []byte("new binary"), 0755)
	os.WriteFile(currentBin+".bak", []byte("old binary"), 0755)

	err := rollback(currentBin)
	if err != nil {
		t.Fatalf("rollback: %v", err)
	}

	got, _ := os.ReadFile(currentBin)
	if string(got) != "old binary" {
		t.Errorf("binary after rollback = %q, want %q", string(got), "old binary")
	}
	if _, err := os.Stat(currentBin + ".bak"); !os.IsNotExist(err) {
		t.Error(".bak should be removed after rollback")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test -C supermac-go ./internal/update/ -run "TestAtomic|TestSwap|TestRollback" -v`
Expected: FAIL — functions not defined

- [ ] **Step 3: Implement atomicSwap and rollback**

Add to `updater.go`:

```go
func atomicSwap(currentBin, newBin string) error {
	bakPath := currentBin + ".bak"

	// Step 1: Rename current to .bak (same filesystem, guaranteed atomic)
	if err := os.Rename(currentBin, bakPath); err != nil {
		return fmt.Errorf("rename current to .bak: %w", err)
	}

	// Step 2: Copy new binary to current location (may be cross-filesystem)
	if err := copyFile(newBin, currentBin); err != nil {
		// Restore .bak on failure
		os.Rename(bakPath, currentBin)
		return fmt.Errorf("copy new binary: %w (restored backup)", err)
	}

	// Step 3: Preserve permissions from old binary
	info, _ := os.Stat(bakPath)
	if info != nil {
		os.Chmod(currentBin, info.Mode())
	}

	// Step 4: Code-sign the new binary (macOS requirement)
	exec.Command("codesign", "-f", "-s", "-", currentBin).Run()

	// Step 5: Remove quarantine attribute
	exec.Command("xattr", "-d", "com.apple.quarantine", currentBin).Run()

	// Step 6: Cleanup new binary source
	os.Remove(newBin)

	return nil
}

func rollback(currentBin string) error {
	bakPath := currentBin + ".bak"
	if _, err := os.Stat(bakPath); os.IsNotExist(err) {
		return fmt.Errorf("no backup found at %s", bakPath)
	}
	if err := os.Rename(bakPath, currentBin); err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return nil
}
```

- [ ] **Step 4: Run tests**

Run: `go test -C supermac-go ./internal/update/ -run "TestAtomic|TestSwap|TestRollback" -v`
Expected: PASS

- [ ] **Step 5: Run all update tests**

Run: `go test -C supermac-go ./internal/update/ -v`
Expected: All PASS

- [ ] **Step 6: Commit**

```bash
git add supermac-go/internal/update/
git commit -m "feat(update): add atomic swap with rollback and code signing"
```

---

### Task 7: Wire into main.go — update command + pre-run check

**Files:**
- Modify: `supermac-go/cmd/mac/main.go`

- [ ] **Step 1: Add update subcommand**

Add import for the update package and a new command function after `doctorCmd()`:

```go
import (
	// ... existing imports ...
	"github.com/cosmolabs-org/supermac/internal/update"
)
```

Add after `rootCmd.AddCommand(doctorCmd())`:

```go
	rootCmd.AddCommand(updateCmd())
```

Add the updateCmd function:

```go
func updateCmd() *cobra.Command {
	var checkOnly bool
	var doRollback bool

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update SuperMac to the latest version",
		RunE: func(cmd *cobra.Command, args []string) {
			if doRollback {
				return runRollback()
			}
			return runUpdate(checkOnly)
		},
	}
	cmd.Flags().BoolVar(&checkOnly, "check", false, "Check for available update without installing")
	cmd.Flags().BoolVar(&doRollback, "rollback", false, "Restore previous version from backup")
	return cmd
}

func runUpdate(checkOnly bool) error {
	checker := update.NewChecker("CosmoLabs-org/SuperMac", version.Version)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rel, err := checker.CheckAvailable(ctx)
	if err != nil {
		return fmt.Errorf("checking for updates: %w", err)
	}
	if rel == nil {
		fmt.Println("✓ SuperMac is up to date (" + version.Version + ")")
		return nil
	}
	if checkOnly {
		fmt.Printf("⬆ SuperMac %s available (current: %s). Run 'mac update' to upgrade.\n", rel.Version, version.Version)
		return nil
	}

	updater, err := update.NewUpdater()
	if err != nil {
		return err
	}
	return updater.Update(rel)
}

func runRollback() error {
	updater, err := update.NewUpdater()
	if err != nil {
		return err
	}
	return updater.Rollback()
}
```

- [ ] **Step 2: Add Update and Rollback methods to Updater**

The `Update` and `Rollback` methods need to be public. Add to `updater.go`:

```go
func (u *Updater) Update(rel *Release) error {
	// Check for Homebrew install
	if strings.Contains(u.binaryPath, "Cellar") || strings.Contains(u.binaryPath, "homebrew") {
		return fmt.Errorf("installed via Homebrew. Run 'brew upgrade supermac' instead")
	}

	fmt.Printf("Updating SuperMac %s → %s...\n", version.Version, rel.Version)

	// Create temp dir
	tmpDir, err := os.MkdirTemp("", "supermac-update")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	// Download tarball
	tarballPath := filepath.Join(tmpDir, "mac-"+u.arch+".tar.gz")
	if err := downloadFile(rel.TarballURL, tarballPath); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	// Download checksums
	checksumsPath := filepath.Join(tmpDir, "checksums.txt")
	if err := downloadFile(rel.ChecksumsURL, checksumsPath); err != nil {
		return fmt.Errorf("download checksums failed: %w", err)
	}

	// Verify checksum
	if err := verifyChecksum(tarballPath, checksumsPath, "mac-"+u.arch+".tar.gz"); err != nil {
		return err
	}

	// Extract binary
	newBin := filepath.Join(tmpDir, "mac")
	if err := extractBinary(tarballPath, "mac-"+u.arch, newBin); err != nil {
		return err
	}

	// Code-sign extracted binary for verification
	exec.Command("codesign", "-f", "-s", "-", newBin).Run()

	// Verify binary runs
	out, err := exec.Command(newBin, "version", "--raw").Output()
	if err != nil {
		return fmt.Errorf("downloaded binary failed verification: %w", err)
	}
	gotVersion := strings.TrimSpace(string(out))
	if gotVersion != rel.Version {
		return fmt.Errorf("version mismatch: expected %s, got %s", rel.Version, gotVersion)
	}

	// Atomic swap
	if err := atomicSwap(u.binaryPath, newBin); err != nil {
		return err
	}

	fmt.Printf("✓ SuperMac updated %s → %s\n", version.Version, rel.Version)
	return nil
}

func (u *Updater) Rollback() error {
	bakPath := u.binaryPath + ".bak"
	if _, err := os.Stat(bakPath); os.IsNotExist(err) {
		return fmt.Errorf("no previous version found for rollback")
	}

	// Verify backup binary
	out, err := exec.Command(bakPath, "version", "--raw").Output()
	bakVersion := "unknown"
	if err == nil {
		bakVersion = strings.TrimSpace(string(out))
	}

	fmt.Printf("Rollback to %s? [Y/n] ", bakVersion)
	if !getYesFlag() {
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		if strings.TrimSpace(strings.ToLower(input)) != "y" && strings.TrimSpace(input) != "" {
			fmt.Println("Rollback cancelled.")
			return nil
		}
	}

	if err := rollback(u.binaryPath); err != nil {
		return err
	}
	fmt.Printf("✓ Rolled back to %s\n", bakVersion)
	return nil
}

func downloadFile(url, path string) error {
	resp, err := http.Get(url) //nolint:gosec // URL is from GitHub Releases
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}
```

- [ ] **Step 3: Add pre-run update check**

Add to `main()` after `rootCmd.AddCommand` lines, before `rootCmd.Execute()`:

```go
	// Check for updates before command execution
	if !quietFlag {
		go func() {
			cfg, _ := config.Load()
			if cfg == nil || !cfg.Updates.Check {
				return
			}
			checker := update.NewChecker("CosmoLabs-org/SuperMac", version.Version)
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			rel, err := checker.CheckAvailable(ctx)
			if err != nil || rel == nil {
				return
			}
			fmt.Fprintf(os.Stderr, "  ⬆ SuperMac %s available. Run 'mac update' to upgrade.\n", rel.Version)
		}()
	}
```

Wait — the brainstorming spec says synchronous with 2s timeout, not goroutine. Let me fix:

```go
	// Check for updates (synchronous, 2s timeout, silent on failure)
	if !quietFlag {
		cfg, _ := config.Load()
		if cfg != nil && cfg.Updates.Check {
			checker := update.NewChecker("CosmoLabs-org/SuperMac", version.Version)
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			rel, _ := checker.CheckAvailable(ctx)
			cancel()
			if rel != nil {
				fmt.Fprintf(os.Stderr, "  ⬆ SuperMac %s available. Run 'mac update' to upgrade.\n", rel.Version)
			}
		}
	}
```

- [ ] **Step 4: Build and smoke test**

Run: `go build -C supermac-go -o mac ./cmd/mac && ./supermac-go/mac update --check`
Expected: either "up to date" or "X available" message

- [ ] **Step 5: Run all tests**

Run: `go test -C supermac-go ./...`
Expected: All PASS

- [ ] **Step 6: Commit**

```bash
git add supermac-go/cmd/mac/main.go supermac-go/internal/update/
git commit -m "feat(update): add mac update command with pre-launch update check"
```

---

### Task 8: Update version output + cleanup

**Files:**
- Modify: `supermac-go/cmd/mac/main.go` (version command)

- [ ] **Step 1: Add update status to version output**

In the `versionCmd()` function (non-raw path), add after the Modules line:

```go
			// Update status
			checker := update.NewChecker("CosmoLabs-org/SuperMac", version.Version)
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			rel, _ := checker.CheckAvailable(ctx)
			cancel()
			updateStatus := "up to date"
			if rel != nil {
				updateStatus = rel.Version + " available"
			}
			fmt.Printf("  Update:     %s\n", updateStatus)
```

- [ ] **Step 2: Build and verify**

Run: `go build -C supermac-go -o mac ./cmd/mac && ./supermac-go/mac version`
Expected: Shows version with Update status line

- [ ] **Step 3: Run full test suite**

Run: `go test -C supermac-go ./...`
Expected: All PASS

- [ ] **Step 4: Commit**

```bash
git add supermac-go/cmd/mac/main.go
git commit -m "feat(version): show update status in version output"
```

---

### Task 9: Final integration test + cleanup

- [ ] **Step 1: Test the full flow**

```bash
./supermac-go/mac version          # shows update status
./supermac-go/mac version --raw    # prints just version number
./supermac-go/mac update --check   # checks for update
./supermac-go/mac update --rollback # reports no backup (first install)
./supermac-go/mac config list       # shows correct Updates field
```

- [ ] **Step 2: Run full test suite**

Run: `go test -C supermac-go ./...`
Expected: All PASS

- [ ] **Step 3: Push and verify CI**

```bash
git push origin master
# Monitor CI run
gh run list -R CosmoLabs-org/SuperMac --limit 1
```

- [ ] **Step 4: Tag and release** (optional, if CI passes)

```bash
git tag v0.3.0
git push origin v0.3.0
```
