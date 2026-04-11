package update

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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

// githubRelease matches the GitHub API response shape.
type githubRelease struct {
	TagName     string `json:"tag_name"`
	PublishedAt string `json:"published_at"`
	Assets      []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
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
		cachePath: filepath.Join(home, ".supermac", "update-cache.json"),
		cacheTTL:  24 * time.Hour,
		current:   currentVersion,
		arch:      arch,
		client:    &http.Client{Timeout: 2 * time.Second},
	}
}

// readCache reads the cached release info. Returns nil if expired or missing.
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

// writeCache persists the release info to disk.
func (c *Checker) writeCache(cr *cachedRelease) error {
	dir := filepath.Dir(c.cachePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.Marshal(cr)
	if err != nil {
		return err
	}
	return os.WriteFile(c.cachePath, data, 0644)
}

// parseRelease extracts a Release from GitHub API JSON.
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

// isNewer returns true if semver a > b.
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

	data := make([]byte, 1024*64)
	n, _ := resp.Body.Read(data)
	data = data[:n]

	rel, err := parseRelease(data, c.arch)
	if err != nil {
		return nil, err
	}

	// Write cache (best effort)
	cr := &cachedRelease{
		CheckedAt:      time.Now(),
		Version:        rel.Version,
		Tag:            rel.Tag,
		TarballURL:     rel.TarballURL,
		ChecksumsURL:   rel.ChecksumsURL,
		PublishedAt:    rel.Date,
		CurrentVersion: c.current,
	}
	c.writeCache(cr)

	if isNewer(rel.Version, c.current) {
		return rel, nil
	}
	return nil, nil
}
