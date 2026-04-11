package update

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

// Update downloads, verifies, and applies the given release.
func (u *Updater) Update(rel *Release) error {
	// Check for Homebrew install
	if strings.Contains(u.binaryPath, "Cellar") || strings.Contains(u.binaryPath, "homebrew") {
		return fmt.Errorf("installed via Homebrew. Run 'brew upgrade supermac' instead")
	}

	fmt.Printf("Updating SuperMac %s → %s...\n", u.currentVersion(), rel.Version)

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

	fmt.Printf("SuperMac updated %s → %s\n", u.currentVersion(), rel.Version)
	return nil
}

// Rollback restores the previous version from the .bak file.
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
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	val := strings.TrimSpace(strings.ToLower(input))
	if val != "y" && val != "" {
		fmt.Println("Rollback cancelled.")
		return nil
	}

	if err := rollback(u.binaryPath); err != nil {
		return err
	}
	fmt.Printf("Rolled back to %s\n", bakVersion)
	return nil
}

func (u *Updater) currentVersion() string {
	out, err := exec.Command(u.binaryPath, "version", "--raw").Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

// verifyChecksum computes SHA256 of the tarball and compares against checksums.txt.
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

// extractBinary extracts the named entry from a tar.gz archive.
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

// atomicSwap replaces the current binary with the new one, keeping a .bak backup.
func atomicSwap(currentBin, newBin string) error {
	bakPath := currentBin + ".bak"

	// Step 1: Rename current to .bak (same filesystem, guaranteed atomic)
	if err := os.Rename(currentBin, bakPath); err != nil {
		return fmt.Errorf("rename current to .bak: %w", err)
	}

	// Step 2: Copy new binary to current location
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

// rollback restores the .bak file to replace the current binary.
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

// copyFile copies src to dst using io.Copy.
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

// downloadFile downloads a URL to a local file path.
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
