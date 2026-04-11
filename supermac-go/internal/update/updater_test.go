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
	dir := t.TempDir()
	content := []byte("test binary content")
	tarballPath := filepath.Join(dir, "mac-arm64.tar.gz")
	os.WriteFile(tarballPath, content, 0644)

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

func TestAtomicSwap(t *testing.T) {
	dir := t.TempDir()
	currentBin := filepath.Join(dir, "mac")
	os.WriteFile(currentBin, []byte("old binary"), 0755)

	newBin := filepath.Join(dir, "new-mac")
	os.WriteFile(newBin, []byte("new binary"), 0755)

	err := atomicSwap(currentBin, newBin)
	if err != nil {
		t.Fatalf("atomicSwap: %v", err)
	}

	got, _ := os.ReadFile(currentBin)
	if string(got) != "new binary" {
		t.Errorf("binary = %q, want %q", string(got), "new binary")
	}
	bak, _ := os.ReadFile(currentBin + ".bak")
	if string(bak) != "old binary" {
		t.Errorf("bak = %q, want %q", string(bak), "old binary")
	}
	if _, err := os.Stat(newBin); !os.IsNotExist(err) {
		t.Error("new binary should be removed after swap")
	}
}

func TestSwapFailsRestore(t *testing.T) {
	dir := t.TempDir()
	currentBin := filepath.Join(dir, "mac")
	os.WriteFile(currentBin, []byte("old binary"), 0755)

	err := atomicSwap(currentBin, filepath.Join(dir, "nonexistent"))
	if err == nil {
		t.Fatal("expected error for missing new binary")
	}

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
