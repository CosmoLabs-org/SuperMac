package display

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/cosmolabs-org/supermac/internal/module"
)

func init() {
	module.Register(&DisplayModule{})
}

type DisplayModule struct{}

func (d *DisplayModule) Name() string                { return "display" }
func (d *DisplayModule) ShortDescription() string     { return "Display and appearance settings management" }
func (d *DisplayModule) Emoji() string                { return "🖥️" }

func (d *DisplayModule) Commands() []module.Command {
	return []module.Command{
		{
			Name:        "brightness",
			Description: "Set screen brightness percentage (0-100)",
			Args: []module.Arg{
				{Name: "level", Required: true, Description: "Brightness level 0-100"},
			},
			Run: d.brightness,
		},
		{
			Name:        "dark-mode",
			Description: "Control dark mode (on/off/toggle)",
			Args: []module.Arg{
				{Name: "action", Required: true, Description: "on, off, or toggle"},
			},
			Run: d.darkMode,
		},
		{
			Name:        "night-shift",
			Description: "Control Night Shift (on/off)",
			Args: []module.Arg{
				{Name: "action", Required: true, Description: "on or off"},
			},
			Run: d.nightShift,
		},
		{
			Name:        "true-tone",
			Description: "Control True Tone (on/off)",
			Args: []module.Arg{
				{Name: "action", Required: true, Description: "on or off"},
			},
			Run: d.trueTone,
		},
		{
			Name:        "wallpaper",
			Description: "Set desktop wallpaper",
			Args: []module.Arg{
				{Name: "path", Required: true, Description: "Path to image file"},
			},
			Run: d.wallpaper,
		},
		{
			Name:        "status",
			Description: "Show display settings and connected displays",
			Run:         d.status,
		},
	}
}

func (d *DisplayModule) Search(term string) []module.SearchResult {
	var results []module.SearchResult
	for _, cmd := range d.Commands() {
		if strings.Contains(cmd.Name, term) || strings.Contains(strings.ToLower(cmd.Description), term) {
			results = append(results, module.SearchResult{
				Command:     cmd.Name,
				Description: cmd.Description,
				Module:      d.Name(),
			})
		}
	}
	return results
}

// --- Brightness ---

func (d *DisplayModule) brightness(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		return module.NewExitError(module.ExitUsage, "Brightness level required: mac display brightness <0-100>")
	}

	level, err := strconv.Atoi(ctx.Args[0])
	if err != nil || level < 0 || level > 100 {
		return module.NewExitError(module.ExitUsage, "Brightness must be between 0 and 100")
	}

	ctx.Output.Info("Setting brightness to %d%%...", level)

	if err := ctx.Platform.SetBrightness(float64(level) / 100.0); err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to set brightness: %v", err))
	}

	ctx.Output.Success("Brightness set to %d%%", level)

	if level <= 20 {
		ctx.Output.Info("Low brightness - good for nighttime use")
	} else if level >= 80 {
		ctx.Output.Info("High brightness - good for bright environments")
	}

	return nil
}

// --- Dark Mode ---

func (d *DisplayModule) darkMode(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		return module.NewExitError(module.ExitUsage, "Action required: mac display dark-mode <on|off|toggle>")
	}

	action := strings.ToLower(ctx.Args[0])

	switch action {
	case "on":
		return d.setDarkMode(ctx, true)
	case "off":
		return d.setDarkMode(ctx, false)
	case "toggle":
		current, err := ctx.Platform.GetDarkMode()
		if err != nil {
			return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to read dark mode state: %v", err))
		}
		ctx.Output.Info("Current mode: %s", modeString(current))
		return d.setDarkMode(ctx, !current)
	default:
		return module.NewExitError(module.ExitUsage, fmt.Sprintf("Invalid action %q: use on, off, or toggle", action))
	}
}

func (d *DisplayModule) setDarkMode(ctx *module.Context, on bool) error {
	current, _ := ctx.Platform.GetDarkMode()
	if current == on {
		ctx.Output.Info("Already in %s mode", modeString(on))
		return nil
	}

	mode := modeString(on)
	ctx.Output.Info("Switching to %s mode...", mode)

	if err := ctx.Platform.SetDarkMode(on); err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to set %s mode: %v", mode, err))
	}

	ctx.Output.Success("Switched to %s mode", mode)
	return nil
}

// --- Night Shift ---

func (d *DisplayModule) nightShift(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		return module.NewExitError(module.ExitUsage, "Action required: mac display night-shift <on|off>")
	}

	action := strings.ToLower(ctx.Args[0])

	var enable bool
	switch action {
	case "on", "enable":
		enable = true
	case "off", "disable":
		enable = false
	default:
		return module.NewExitError(module.ExitUsage, fmt.Sprintf("Invalid action %q: use on or off", action))
	}

	state := "enabled"
	if !enable {
		state = "disabled"
	}

	script := fmt.Sprintf(
		`tell application "System Events" to tell appearance preferences to set night shift enabled to %v`,
		enable,
	)
	ctx.Output.Info("Setting Night Shift %s...", state)

	if _, err := ctx.Platform.RunOSAScript(script); err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to set Night Shift: %v", err))
	}

	ctx.Output.Success("Night Shift %s", state)
	if enable {
		ctx.Output.Info("Night Shift reduces blue light for better sleep")
	}
	return nil
}

// --- True Tone ---

func (d *DisplayModule) trueTone(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		return module.NewExitError(module.ExitUsage, "Action required: mac display true-tone <on|off>")
	}

	action := strings.ToLower(ctx.Args[0])

	var enable bool
	switch action {
	case "on", "enable":
		enable = true
	case "off", "disable":
		enable = false
	default:
		return module.NewExitError(module.ExitUsage, fmt.Sprintf("Invalid action %q: use on or off", action))
	}

	// Check True Tone support
	out, err := exec.Command("system_profiler", "SPDisplaysDataType").Output()
	if err != nil || !strings.Contains(string(out), "True Tone") {
		return module.NewExitError(module.ExitGeneral, "True Tone is not supported on this display")
	}

	state := "enabled"
	if !enable {
		state = "disabled"
	}

	script := fmt.Sprintf(
		`tell application "System Events" to tell appearance preferences to set true tone enabled to %v`,
		enable,
	)
	ctx.Output.Info("Setting True Tone %s...", state)

	if _, err := ctx.Platform.RunOSAScript(script); err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to set True Tone: %v", err))
	}

	ctx.Output.Success("True Tone %s", state)
	if enable {
		ctx.Output.Info("True Tone adjusts colors based on ambient lighting")
	}
	return nil
}

// --- Wallpaper ---

func (d *DisplayModule) wallpaper(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		return module.NewExitError(module.ExitUsage, "Path required: mac display wallpaper <path>")
	}

	path := ctx.Args[0]
	ctx.Output.Info("Setting wallpaper to %s...", path)

	script := fmt.Sprintf(
		`tell application "System Events" to set picture of every desktop to %q`,
		path,
	)

	if _, err := ctx.Platform.RunOSAScript(script); err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to set wallpaper: %v", err))
	}

	ctx.Output.Success("Wallpaper set")
	return nil
}

// --- Status ---

func (d *DisplayModule) status(ctx *module.Context) error {
	ctx.Output.Header("Display Status")

	// Brightness via osascript
	brightnessOut, err := ctx.Platform.RunOSAScript(
		`tell application "System Events" to get brightness of first item of (displays whose built in is true)`,
	)
	if err != nil {
		fmt.Printf("  Brightness:          unknown\n")
	} else {
		val, _ := strconv.ParseFloat(strings.TrimSpace(brightnessOut), 64)
		fmt.Printf("  Brightness:          %d%%\n", int(val*100))
	}

	// Dark mode
	darkMode, err := ctx.Platform.GetDarkMode()
	if err != nil {
		fmt.Printf("  Appearance:          unknown\n")
	} else {
		fmt.Printf("  Appearance:          %s Mode\n", modeString(darkMode))
	}

	// Night Shift status via CoreBrightness defaults
	nsStatus := "unknown"
	cbOut, err := ctx.Platform.ReadDefault(
		"com.apple.CoreBrightness",
		fmt.Sprintf("CBUser-%d", uid()),
	)
	if err == nil && strings.Contains(cbOut, "BlueLightReductionEnabled") {
		if strings.Contains(cbOut, "BlueLightReductionEnabled = 1") {
			nsStatus = "enabled"
		} else {
			nsStatus = "disabled"
		}
	}
	fmt.Printf("  Night Shift:         %s\n", nsStatus)

	// Connected displays and resolution via system_profiler
	spOut, err := exec.Command("system_profiler", "SPDisplaysDataType").Output()
	if err == nil {
		displayCount := strings.Count(string(spOut), "Display Type")
		if displayCount == 0 {
			displayCount = 1
		}
		fmt.Printf("  Connected Displays:  %d\n", displayCount)

		for _, line := range strings.Split(string(spOut), "\n") {
			if strings.Contains(line, "Resolution:") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					fmt.Printf("  Primary Resolution: %s\n", strings.TrimSpace(parts[1]))
				}
				break
			}
		}
	}

	return nil
}

// --- Helpers ---

func modeString(dark bool) string {
	if dark {
		return "Dark"
	}
	return "Light"
}

func uid() int {
	out, err := exec.Command("id", "-u").Output()
	if err != nil {
		return 0
	}
	val, _ := strconv.Atoi(strings.TrimSpace(string(out)))
	return val
}
