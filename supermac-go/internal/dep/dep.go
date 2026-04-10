package dep

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

// Dependency describes an external tool required by a module.
type Dependency struct {
	Name     string   // display name: "blueutil"
	Brew     string   // brew formula name: "blueutil"
	Check    string   // binary to look for on PATH: "blueutil"
	Commands []string // which commands need this (nil = all commands in module)
}

// IsInstalled checks whether the dependency binary exists on PATH.
func (d Dependency) IsInstalled() bool {
	_, err := exec.LookPath(d.Check)
	return err == nil
}

// Install runs brew install for the dependency.
func (d Dependency) Install() error {
	brew, err := exec.LookPath("brew")
	if err != nil {
		return fmt.Errorf("Homebrew is not installed. Install it from https://brew.sh")
	}
	cmd := exec.Command(brew, "install", d.Brew)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Ensure checks if the dependency is installed. If not and interactive is true,
// it prompts the user to install via Homebrew. If interactive is false, it returns an error.
func (d Dependency) Ensure(interactive bool) error {
	if d.IsInstalled() {
		return nil
	}

	if !interactive {
		return fmt.Errorf("%s is not installed. Install with: brew install %s", d.Name, d.Brew)
	}

	// Auto-install with SUPERMAC_AUTO_INSTALL env var
	if os.Getenv("SUPERMAC_AUTO_INSTALL") == "1" {
		fmt.Printf("  Installing %s via Homebrew...\n", d.Name)
		if err := d.Install(); err != nil {
			return fmt.Errorf("failed to install %s: %w", d.Name, err)
		}
		fmt.Printf("  ✓ %s installed successfully.\n", d.Name)
		return nil
	}

	fmt.Printf("\n  %s is not installed.\n", d.Name)
	fmt.Printf("  Install via Homebrew? [Y/n]: ")

	var response string
	fmt.Scanln(&response)
	response = strings.TrimSpace(strings.ToLower(response))

	if response != "" && response != "y" && response != "yes" {
		return fmt.Errorf("%s is required. Install with: brew install %s", d.Name, d.Brew)
	}

	fmt.Printf("  Installing %s...\n", d.Name)
	if err := d.Install(); err != nil {
		return fmt.Errorf("failed to install %s: %w", d.Name, err)
	}
	fmt.Printf("  ✓ %s installed successfully.\n\n", d.Name)
	return nil
}

// AffectsCommand returns true if this dependency is needed for the given command.
func (d Dependency) AffectsCommand(commandName string) bool {
	if d.Commands == nil {
		return true // affects all commands
	}
	return slices.Contains(d.Commands, commandName)
}

// CheckBrew returns true if Homebrew is installed and provides its path.
func CheckBrew() (string, bool) {
	path, err := exec.LookPath("brew")
	if err != nil {
		return "", false
	}
	return path, true
}
