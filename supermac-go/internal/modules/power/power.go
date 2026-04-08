package power

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/cosmolabs-org/supermac/internal/module"
)

func init() {
	module.Register(&PowerModule{})
}

type PowerModule struct{}

func (p *PowerModule) Name() string            { return "power" }
func (p *PowerModule) ShortDescription() string { return "Power user toggles for macOS" }
func (p *PowerModule) Emoji() string            { return "⚡" }

// toggleDef defines a single toggle command.
type toggleDef struct {
	name        string
	desc        string
	getState    func() (string, error)           // returns "on" or "off"
	setState    func(on bool) error               // set the state
	restarter   func()                            // restart process (nil = none)
	restartMsg  string                            // message if logout needed
	needsSudo   bool
}

func (p *PowerModule) Commands() []module.Command {
	cmds := []module.Command{
		{
			Name:        "status",
			Description: "Show all power toggles and their current state",
			Run:         p.status,
		},
	}

	for _, t := range allToggles() {
		// Capture for closure
		toggle := t
		cmds = append(cmds, module.Command{
			Name:        toggle.name,
			Description: toggle.desc,
			Args: []module.Arg{
				{Name: "action", Required: false, Description: "on, off, or toggle (omit to show state)"},
			},
			Run: func(ctx *module.Context) error {
				return runToggle(ctx, &toggle)
			},
		})
	}
	return cmds
}

func (p *PowerModule) Search(term string) []module.SearchResult {
	var results []module.SearchResult
	for _, cmd := range p.Commands() {
		if strings.Contains(cmd.Name, term) || strings.Contains(strings.ToLower(cmd.Description), term) {
			results = append(results, module.SearchResult{
				Command:     cmd.Name,
				Description: cmd.Description,
				Module:      p.Name(),
			})
		}
	}
	return results
}

// ============================================================================
// All Toggles
// ============================================================================

func allToggles() []toggleDef {
	return []toggleDef{
		{
			name:       "caffeinate",
			desc:       "Prevent system/display sleep (caffeinate)",
			getState:   getCaffeinateState,
			setState:   setCaffeinateState,
			restarter:  nil,
			restartMsg: "",
			needsSudo:  false,
		},
		{
			name:       "hidden-files",
			desc:       "Show hidden files in Finder",
			getState:   func() (string, error) { return readDefault("com.apple.finder", "AppleShowAllFiles") },
			setState:   func(on bool) error { return writeDefaultBool("com.apple.finder", "AppleShowAllFiles", on) },
			restarter:  restartFinder,
			restartMsg: "",
			needsSudo:  false,
		},
		{
			name:       "file-extensions",
			desc:       "Show all file extensions in Finder",
			getState:   func() (string, error) { return readDefault("NSGlobalDomain", "AppleShowAllExtensions") },
			setState:   func(on bool) error { return writeDefaultBool("NSGlobalDomain", "AppleShowAllExtensions", on) },
			restarter:  restartFinder,
			restartMsg: "",
			needsSudo:  false,
		},
		{
			name:       "desktop-icons",
			desc:       "Show icons on desktop",
			getState:   func() (string, error) { return readDefault("com.apple.finder", "CreateDesktop") },
			setState:   func(on bool) error { return writeDefaultBool("com.apple.finder", "CreateDesktop", on) },
			restarter:  restartFinder,
			restartMsg: "",
			needsSudo:  false,
		},
		{
			name:       "gatekeeper",
			desc:       "Allow apps from unidentified developers",
			getState:   getGatekeeperState,
			setState:   setGatekeeperState,
			restarter:  nil,
			restartMsg: "",
			needsSudo:  true,
		},
		{
			name:       "crash-reporter",
			desc:       "Show crash reporter dialog (vs silent)",
			getState:   func() (string, error) { return readDefault("com.apple.CrashReporter", "DialogType") },
			setState:   setCrashReporterState,
			restarter:  nil,
			restartMsg: "",
			needsSudo:  false,
		},
		{
			name:       "function-keys",
			desc:       "Use F1-F12 as standard function keys",
			getState:   func() (string, error) { return readDefault("com.apple.keyboard", "fnState") },
			setState:   func(on bool) error { return writeDefaultBool("com.apple.keyboard", "fnState", on) },
			restarter:  nil,
			restartMsg: "",
			needsSudo:  false,
		},
		{
			name:       "spotlight-indexing",
			desc:       "Enable/disable Spotlight indexing",
			getState:   getSpotlightState,
			setState:   setSpotlightState,
			restarter:  nil,
			restartMsg: "",
			needsSudo:  true,
		},
		{
			name:       "key-repeat",
			desc:       "Fast key repeat rate (2x default)",
			getState:   getKeyRepeatState,
			setState:   setKeyRepeatState,
			restarter:  nil,
			restartMsg: "Log out and back in for changes to take effect",
			needsSudo:  false,
		},
		{
			name:       "smooth-scrolling",
			desc:       "Smooth scrolling animation",
			getState:   func() (string, error) { return readDefault("NSGlobalDomain", "AppleScrollAnimationEnabled") },
			setState:   func(on bool) error { return writeDefaultBool("NSGlobalDomain", "AppleScrollAnimationEnabled", on) },
			restarter:  nil,
			restartMsg: "",
			needsSudo:  false,
		},
		{
			name:       "animations",
			desc:       "Window open/close animations",
			getState:   func() (string, error) { return readDefault("NSGlobalDomain", "NSAutomaticWindowAnimationsEnabled") },
			setState:   func(on bool) error { return writeDefaultBool("NSGlobalDomain", "NSAutomaticWindowAnimationsEnabled", on) },
			restarter:  nil,
			restartMsg: "Log out and back in for changes to take effect",
			needsSudo:  false,
		},
		{
			name:       "transparency",
			desc:       "Reduce UI transparency",
			getState:   func() (string, error) { return readDefault("com.apple.universalaccess", "reduceTransparency") },
			setState:   func(on bool) error { return writeDefaultBool("com.apple.universalaccess", "reduceTransparency", on) },
			restarter:  nil,
			restartMsg: "",
			needsSudo:  false,
		},
		{
			name:       "dock-bounce",
			desc:       "Dock bounce animation on app launch",
			getState:   func() (string, error) { return readDefault("com.apple.dock", "launchanim") },
			setState:   func(on bool) error { return writeDefaultBool("com.apple.dock", "launchanim", on) },
			restarter:  restartDock,
			restartMsg: "",
			needsSudo:  false,
		},
		{
			name:       "firewall-stealth",
			desc:       "Firewall stealth mode (no ping response)",
			getState:   getFirewallStealthState,
			setState:   setFirewallStealthState,
			restarter:  nil,
			restartMsg: "",
			needsSudo:  true,
		},
		{
			name:       "remote-login",
			desc:       "SSH remote login server",
			getState:   getRemoteLoginState,
			setState:   setRemoteLoginState,
			restarter:  nil,
			restartMsg: "",
			needsSudo:  true,
		},
		{
			name:       "quarantine",
			desc:       "File quarantine warnings on downloads",
			getState:   func() (string, error) { return readDefault("com.apple.LaunchServices", "LSQuarantine") },
			setState:   func(on bool) error { return writeDefaultBool("com.apple.LaunchServices", "LSQuarantine", on) },
			restarter:  nil,
			restartMsg: "",
			needsSudo:  false,
		},
		{
			name:       "developer-dir",
			desc:       "Xcode developer tools path",
			getState:   getDeveloperDirState,
			setState:   setDeveloperDirState,
			restarter:  nil,
			restartMsg: "",
			needsSudo:  true,
		},
		{
			name:       "login-items",
			desc:       "Show login window instead of auto-login",
			getState:   func() (string, error) { return readDefault("com.apple.loginwindow", "autoLoginUser") },
			setState:   setLoginItemsState,
			restarter:  nil,
			restartMsg: "",
			needsSudo:  false,
		},
		{
			name:       "save-panel",
			desc:       "Expanded save dialogs by default",
			getState:   func() (string, error) { return readDefault("NSGlobalDomain", "NSNavPanelExpandedStateForSaveMode") },
			setState:   func(on bool) error { return writeDefaultBool("NSGlobalDomain", "NSNavPanelExpandedStateForSaveMode", on) },
			restarter:  nil,
			restartMsg: "",
			needsSudo:  false,
		},
		{
			name:       "print-dialog",
			desc:       "Expanded print dialogs by default",
			getState:   func() (string, error) { return readDefault("NSGlobalDomain", "PMPrintingExpandedStateForPrint") },
			setState:   func(on bool) error { return writeDefaultBool("NSGlobalDomain", "PMPrintingExpandedStateForPrint", on) },
			restarter:  nil,
			restartMsg: "",
			needsSudo:  false,
		},
	}
}

// ============================================================================
// Status Overview
// ============================================================================

func (p *PowerModule) status(ctx *module.Context) error {
	ctx.Output.Header("Power Toggles")
	fmt.Println()

	for _, t := range allToggles() {
		state, err := t.getState()
		label := formatState(state, err)
		fmt.Printf("  %-18s %s\n", t.name, label)
	}

	return nil
}

func formatState(state string, err error) string {
	if err != nil {
		return "unknown"
	}
	s := strings.TrimSpace(strings.ToLower(state))
	switch {
	case s == "1", s == "true", s == "yes", s == "on", s == "enabled", s == "active":
		return "enabled"
	case s == "0", s == "false", s == "no", s == "off", s == "disabled", s == "":
		return "disabled"
	default:
		return state
	}
}

// ============================================================================
// Toggle Runner (shared state machine)
// ============================================================================

func runToggle(ctx *module.Context, t *toggleDef) error {
	// Show state
	if len(ctx.Args) == 0 {
		state, err := t.getState()
		if err != nil {
			ctx.Output.Info("%s: unknown", t.name)
			return nil
		}
		state = formatState(state, err)
		ctx.Output.Info("%s: %s", t.name, state)
		ctx.Output.Info("Use 'mac power %s <on|off|toggle>' to change", t.name)
		return nil
	}

	action := strings.ToLower(ctx.Args[0])

	// Determine desired state
	current, _ := t.getState()
	currentOn := isOn(current)

	var wantOn bool
	switch action {
	case "on":
		wantOn = true
	case "off":
		wantOn = false
	case "toggle":
		wantOn = !currentOn
	default:
		return module.NewExitError(module.ExitUsage, "Action must be: on, off, or toggle")
	}

	// Sudo check
	if t.needsSudo {
		if err := exec.Command("sudo", "-n", "true").Run(); err != nil {
			return module.NewExitError(module.ExitPermission,
				fmt.Sprintf("%s requires admin privileges. Run: sudo mac power %s %s", t.name, t.name, action))
		}
	}

	// Apply
	desired := "off"
	if wantOn {
		desired = "on"
	}
	ctx.Output.Info("Setting %s to %s...", t.name, desired)

	if err := t.setState(wantOn); err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to set %s: %v", t.name, err))
	}

	// Restart if needed
	if t.restarter != nil {
		t.restarter()
	}

	if wantOn {
		ctx.Output.Success("%s enabled", t.name)
	} else {
		ctx.Output.Success("%s disabled", t.name)
	}

	if t.restartMsg != "" {
		ctx.Output.Info("%s", t.restartMsg)
	}

	return nil
}

// ============================================================================
// State Getters (special cases)
// ============================================================================

func getCaffeinateState() (string, error) {
	data, err := os.ReadFile("/tmp/supermac-caffeinate.pid")
	if err != nil {
		return "0", nil
	}
	pid, _ := strconv.Atoi(strings.TrimSpace(string(data)))
	if pid <= 0 {
		return "0", nil
	}
	if err := exec.Command("kill", "-0", strconv.Itoa(pid)).Run(); err != nil {
		return "0", nil // process gone
	}
	return "1", nil
}

func getGatekeeperState() (string, error) {
	out, err := exec.Command("spctl", "--status").CombinedOutput()
	if err != nil {
		return "", err
	}
	if strings.Contains(string(out), "assessments disabled") {
		return "1", nil // gatekeeper OFF = allow anywhere
	}
	return "0", nil
}

func getSpotlightState() (string, error) {
	out, err := exec.Command("mdutil", "-s", "/").CombinedOutput()
	if err != nil {
		return "", err
	}
	if strings.Contains(string(out), "disabled") {
		return "0", nil
	}
	return "1", nil
}

func getKeyRepeatState() (string, error) {
	val, err := readDefault("NSGlobalDomain", "KeyRepeat")
	if err != nil {
		return "", err
	}
	// Default KeyRepeat is ~6 (30ms). Fast = 2 (15ms)
	if strings.TrimSpace(val) == "2" {
		return "1", nil
	}
	return "0", nil
}

func getFirewallStealthState() (string, error) {
	out, err := exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--getstealthmode").CombinedOutput()
	if err != nil {
		return "", err
	}
	if strings.Contains(string(out), "enabled") {
		return "1", nil
	}
	return "0", nil
}

func getRemoteLoginState() (string, error) {
	out, err := exec.Command("systemsetup", "-getremotelogin").CombinedOutput()
	if err != nil {
		return "", err
	}
	if strings.Contains(string(out), "Off") {
		return "0", nil
	}
	return "1", nil
}

func getDeveloperDirState() (string, error) {
	out, err := exec.Command("xcode-select", "-p").Output()
	if err != nil {
		return "", err
	}
	// Just return the path as "on" since it exists
	path := strings.TrimSpace(string(out))
	if path != "" {
		return "1", nil
	}
	return "0", nil
}

// ============================================================================
// State Setters (special cases)
// ============================================================================

func setCaffeinateState(on bool) error {
	if on {
		cmd := exec.Command("caffeinate", "-d")
		if err := cmd.Start(); err != nil {
			return err
		}
		return os.WriteFile("/tmp/supermac-caffeinate.pid", []byte(strconv.Itoa(cmd.Process.Pid)), 0644)
	}
	// Off: kill the process
	data, err := os.ReadFile("/tmp/supermac-caffeinate.pid")
	if err != nil {
		return nil // nothing to stop
	}
	pid := strings.TrimSpace(string(data))
	if pid != "" {
		exec.Command("kill", pid).Run()
		os.Remove("/tmp/supermac-caffeinate.pid")
	}
	return nil
}

func setGatekeeperState(on bool) error {
	if on {
		return exec.Command("sudo", "spctl", "--master-disable").Run()
	}
	return exec.Command("sudo", "spctl", "--master-enable").Run()
}

func setCrashReporterState(on bool) error {
	val := "none"
	if on {
		val = "prompt"
	}
	return exec.Command("defaults", "write", "com.apple.CrashReporter", "DialogType", "-string", val).Run()
}

func setSpotlightState(on bool) error {
	if on {
		return exec.Command("sudo", "mdutil", "-i", "on", "/").Run()
	}
	return exec.Command("sudo", "mdutil", "-i", "off", "/").Run()
}

func setKeyRepeatState(on bool) error {
	if on {
		exec.Command("defaults", "write", "NSGlobalDomain", "KeyRepeat", "-int", "2").Run()
		return exec.Command("defaults", "write", "NSGlobalDomain", "InitialKeyRepeat", "-int", "15").Run()
	}
	exec.Command("defaults", "write", "NSGlobalDomain", "KeyRepeat", "-int", "6").Run()
	return exec.Command("defaults", "write", "NSGlobalDomain", "InitialKeyRepeat", "-int", "25").Run()
}

func setFirewallStealthState(on bool) error {
	if on {
		return exec.Command("sudo", "/usr/libexec/ApplicationFirewall/socketfilterfw", "--setstealthmode", "on").Run()
	}
	return exec.Command("sudo", "/usr/libexec/ApplicationFirewall/socketfilterfw", "--setstealthmode", "off").Run()
}

func setRemoteLoginState(on bool) error {
	if on {
		return exec.Command("sudo", "systemsetup", "-setremotelogin", "on").Run()
	}
	return exec.Command("sudo", "systemsetup", "-setremotelogin", "off").Run()
}

func setDeveloperDirState(on bool) error {
	if !on {
		// Can't really "disable" developer dir, just show current
		return nil
	}
	return exec.Command("sudo", "xcode-select", "--install").Run()
}

func setLoginItemsState(on bool) error {
	if on {
		return exec.Command("defaults", "write", "com.apple.loginwindow", "autoLoginUser", "-string", "").Run()
	}
	return exec.Command("defaults", "delete", "com.apple.loginwindow", "autoLoginUser").Run()
}

// ============================================================================
// Helpers
// ============================================================================

func readDefault(domain, key string) (string, error) {
	out, err := exec.Command("defaults", "read", domain, key).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func writeDefaultBool(domain, key string, on bool) error {
	val := "0"
	if on {
		val = "1"
	}
	return exec.Command("defaults", "write", domain, key, "-bool", val).Run()
}

func isOn(state string) bool {
	s := strings.TrimSpace(strings.ToLower(state))
	return s == "1" || s == "true" || s == "yes" || s == "on" || s == "enabled" || s == "active"
}

func restartFinder() {
	exec.Command("killall", "Finder").Run()
}

func restartDock() {
	exec.Command("killall", "Dock").Run()
}
