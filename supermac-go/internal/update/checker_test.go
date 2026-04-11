package update

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
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
		CheckedAt:      time.Now().Add(-48 * time.Hour),
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
	tests := []struct {
		a, b string
		want bool
	}{
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
