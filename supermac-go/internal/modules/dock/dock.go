package dock

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/cosmolabs-org/supermac/internal/module"
)

const domain = "com.apple.dock"

func init() {
	module.Register(&DockModule{})
}

type DockModule struct{}

func (d *DockModule) Name() string            { return "dock" }
func (d *DockModule) ShortDescription() string { return "Dock position, size, and appearance management" }
func (d *DockModule) Emoji() string            { return "🚢" }

func (d *DockModule) Commands() []module.Command {
	return []module.Command{
		{
			Name:        "position",
			Description: "Set dock position (left/bottom/right)",
			Args: []module.Arg{
				{Name: "placement", Required: true, Description: "Position: left, bottom, or right"},
			},
			Run: d.position,
		},
		{
			Name:        "autohide",
			Description: "Toggle dock auto-hide (on/off/toggle)",
			Args: []module.Arg{
				{Name: "state", Required: true, Description: "on, off, or toggle"},
			},
			Run: d.autohide,
		},
		{
			Name:        "size",
			Description: "Set dock icon size (small/medium/large or pixel value)",
			Args: []module.Arg{
				{Name: "value", Required: true, Description: "small, medium, large, or pixel count (e.g. 48)"},
			},
			Run: d.size,
		},
		{
			Name:        "magnification",
			Description: "Toggle dock magnification (on/off/toggle)",
			Args: []module.Arg{
				{Name: "state", Required: true, Description: "on, off, or toggle"},
			},
			Run: d.magnification,
		},
		{
			Name:        "magnification-size",
			Description: "Set magnified icon size in pixels",
			Args: []module.Arg{
				{Name: "value", Required: true, Description: "Pixel size when magnified (e.g. 128)"},
			},
			Run: d.magnificationSize,
		},
		{
			Name:        "minimize-effect",
			Description: "Set window minimize effect (genie/scale)",
			Args: []module.Arg{
				{Name: "effect", Required: true, Description: "genie or scale"},
			},
			Run: d.minimizeEffect,
		},
		{
			Name:        "status",
			Description: "Show current dock settings",
			Run:         d.status,
		},
		{
			Name:        "reset",
			Description: "Reset dock to default settings",
			Run:         d.reset,
		},
	}
}

func (d *DockModule) Search(term string) []module.SearchResult {
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

// ---------------------------------------------------------------------------
// Command implementations
// ---------------------------------------------------------------------------

func (d *DockModule) position(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		return module.NewExitError(module.ExitUsage, "Position required: mac dock position <left|bottom|right>")
	}

	pos := normalizePosition(ctx.Args[0])
	if pos == "" {
		return module.NewExitError(module.ExitUsage,
			fmt.Sprintf("Invalid position: %s (valid: left, bottom, right)", ctx.Args[0]))
	}

	current, _ := readString(ctx, "orientation")
	if current == pos {
		ctx.Output.Info("Dock is already positioned on the %s", pos)
		return nil
	}

	ctx.Output.Info("Moving dock to %s...", pos)
	if err := ctx.Platform.WriteDefault(domain, "orientation", pos); err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to set dock position: %v", err))
	}

	if err := restartDock(ctx); err != nil {
		return err
	}

	ctx.Output.Success("Dock moved to %s", pos)
	if pos == "left" || pos == "right" {
		ctx.Output.Info("Vertical dock gives you more horizontal screen space")
	}
	return nil
}

func (d *DockModule) autohide(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		return module.NewExitError(module.ExitUsage, "State required: mac dock autohide <on|off|toggle>")
	}

	want, err := parseOnOffToggle(ctx.Args[0], ctx, "autohide")
	if err != nil {
		return err
	}

	ctx.Output.Info("Setting dock auto-hide %s...", stateString(want))
	if err := ctx.Platform.WriteDefault(domain, "autohide", "-bool "+strconv.FormatBool(want)); err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to set auto-hide: %v", err))
	}

	if err := restartDock(ctx); err != nil {
		return err
	}

	ctx.Output.Success("Dock auto-hide %s", stateString(want))
	if want {
		ctx.Output.Info("Move cursor to screen edge to show dock")
	} else {
		ctx.Output.Info("Dock will always be visible")
	}
	return nil
}

func (d *DockModule) size(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		return module.NewExitError(module.ExitUsage, "Size required: mac dock size <small|medium|large|pixel-value>")
	}

	tilesize, label := parseSize(ctx.Args[0])
	if tilesize == 0 {
		return module.NewExitError(module.ExitUsage,
			fmt.Sprintf("Invalid size: %s (valid: small, medium, large, or pixel count)", ctx.Args[0]))
	}

	ctx.Output.Info("Setting dock size to %s (%dpx)...", label, tilesize)
	if err := ctx.Platform.WriteDefault(domain, "tilesize", fmt.Sprintf("-int %d", tilesize)); err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to set dock size: %v", err))
	}

	if err := restartDock(ctx); err != nil {
		return err
	}

	ctx.Output.Success("Dock size set to %s (%dpx)", label, tilesize)
	return nil
}

func (d *DockModule) magnification(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		return module.NewExitError(module.ExitUsage, "State required: mac dock magnification <on|off|toggle>")
	}

	want, err := parseOnOffToggle(ctx.Args[0], ctx, "magnification")
	if err != nil {
		return err
	}

	ctx.Output.Info("Setting dock magnification %s...", stateString(want))
	if err := ctx.Platform.WriteDefault(domain, "magnification", "-bool "+strconv.FormatBool(want)); err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to set magnification: %v", err))
	}

	// When enabling magnification, also set a reasonable magnified size
	if want {
		_ = ctx.Platform.WriteDefault(domain, "largesize", "-int 128")
	}

	if err := restartDock(ctx); err != nil {
		return err
	}

	ctx.Output.Success("Dock magnification %s", stateString(want))
	if want {
		ctx.Output.Info("Hover over dock icons to see magnification effect")
	}
	return nil
}

func (d *DockModule) magnificationSize(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		return module.NewExitError(module.ExitUsage, "Value required: mac dock magnification-size <pixels>")
	}

	val, err := strconv.Atoi(ctx.Args[0])
	if err != nil || val < 16 || val > 256 {
		return module.NewExitError(module.ExitUsage,
			fmt.Sprintf("Invalid magnification size: %s (must be 16-256 pixels)", ctx.Args[0]))
	}

	ctx.Output.Info("Setting magnification size to %dpx...", val)
	if err := ctx.Platform.WriteDefault(domain, "largesize", fmt.Sprintf("-int %d", val)); err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to set magnification size: %v", err))
	}

	if err := restartDock(ctx); err != nil {
		return err
	}

	ctx.Output.Success("Magnification size set to %dpx", val)
	return nil
}

func (d *DockModule) minimizeEffect(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		return module.NewExitError(module.ExitUsage, "Effect required: mac dock minimize-effect <genie|scale>")
	}

	effect := strings.ToLower(ctx.Args[0])
	if effect != "genie" && effect != "scale" {
		return module.NewExitError(module.ExitUsage,
			fmt.Sprintf("Invalid minimize effect: %s (valid: genie, scale)", ctx.Args[0]))
	}

	ctx.Output.Info("Setting minimize effect to %s...", effect)
	if err := ctx.Platform.WriteDefault(domain, "mineffect", effect); err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to set minimize effect: %v", err))
	}

	if err := restartDock(ctx); err != nil {
		return err
	}

	ctx.Output.Success("Minimize effect set to %s", effect)
	return nil
}

func (d *DockModule) status(ctx *module.Context) error {
	ctx.Output.Header("Dock Status")
	fmt.Println()

	// Position
	if pos, err := readString(ctx, "orientation"); err == nil {
		fmt.Printf("  Position:          %s\n", pos)
	} else {
		fmt.Printf("  Position:          bottom (default)\n")
	}

	// Auto-hide
	if val, err := readBool(ctx, "autohide"); err == nil {
		fmt.Printf("  Auto-hide:         %s\n", stateString(val))
	} else {
		fmt.Printf("  Auto-hide:         disabled\n")
	}

	// Size
	if raw, err := readString(ctx, "tilesize"); err == nil {
		px, _ := strconv.Atoi(strings.TrimSpace(raw))
		label := sizeLabel(px)
		fmt.Printf("  Size:              %s (%dpx)\n", label, px)
	} else {
		fmt.Printf("  Size:              medium (64px)\n")
	}

	// Magnification
	if val, err := readBool(ctx, "magnification"); err == nil {
		fmt.Printf("  Magnification:     %s\n", stateString(val))
	} else {
		fmt.Printf("  Magnification:     disabled\n")
	}

	// Magnification size
	if raw, err := readString(ctx, "largesize"); err == nil {
		px, _ := strconv.Atoi(strings.TrimSpace(raw))
		fmt.Printf("  Mag. size:         %dpx\n", px)
	}

	// Minimize effect
	if effect, err := readString(ctx, "mineffect"); err == nil {
		fmt.Printf("  Minimize effect:   %s\n", strings.TrimSpace(effect))
	} else {
		fmt.Printf("  Minimize effect:   genie (default)\n")
	}

	// Show recent apps
	if val, err := readBool(ctx, "show-recents"); err == nil {
		if val {
			fmt.Printf("  Recent apps:       shown\n")
		} else {
			fmt.Printf("  Recent apps:       hidden\n")
		}
	}

	return nil
}

func (d *DockModule) reset(ctx *module.Context) error {
	ctx.Output.Warning("This will reset all dock settings to defaults")
	ctx.Output.Info("Resetting dock to default settings...")

	// Delete all dock preferences
	_ = ctx.Platform.DeleteDefault(domain, "")

	// Set sensible defaults
	_ = ctx.Platform.WriteDefault(domain, "orientation", "bottom")
	_ = ctx.Platform.WriteDefault(domain, "autohide", "-bool false")
	_ = ctx.Platform.WriteDefault(domain, "tilesize", "-int 64")
	_ = ctx.Platform.WriteDefault(domain, "magnification", "-bool false")
	_ = ctx.Platform.WriteDefault(domain, "show-recents", "-bool true")

	if err := restartDock(ctx); err != nil {
		return err
	}

	ctx.Output.Success("Dock reset to default settings")
	ctx.Output.Info("Position: bottom, Size: medium, Auto-hide: off")
	return nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// restartDock sends killall Dock to apply pending changes.
func restartDock(ctx *module.Context) error {
	cmd := exec.Command("killall", "Dock")
	if err := cmd.Run(); err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to restart Dock: %v", err))
	}
	return nil
}

// readString reads a defaults key from the dock domain.
func readString(ctx *module.Context, key string) (string, error) {
	return ctx.Platform.ReadDefault(domain, key)
}

// readBool reads a defaults key and interprets it as a bool.
func readBool(ctx *module.Context, key string) (bool, error) {
	out, err := ctx.Platform.ReadDefault(domain, key)
	if err != nil {
		return false, err
	}
	trimmed := strings.TrimSpace(out)
	return trimmed == "1" || trimmed == "true", nil
}

// normalizePosition maps shorthand or full names to a valid position string.
// Returns empty string if invalid.
func normalizePosition(input string) string {
	switch strings.ToLower(input) {
	case "left", "l":
		return "left"
	case "bottom", "b":
		return "bottom"
	case "right", "r":
		return "right"
	default:
		return ""
	}
}

// parseOnOffToggle resolves "on"/"off"/"toggle" to a boolean.
// For "toggle", it reads the current value from the given key.
func parseOnOffToggle(input string, ctx *module.Context, key string) (bool, error) {
	switch strings.ToLower(input) {
	case "on", "enable", "true", "1":
		return true, nil
	case "off", "disable", "false", "0":
		return false, nil
	case "toggle", "t":
		current, err := readBool(ctx, key)
		if err != nil {
			return true, nil // default to on if unreadable
		}
		return !current, nil
	default:
		return false, module.NewExitError(module.ExitUsage,
			fmt.Sprintf("Invalid state: %s (valid: on, off, toggle)", input))
	}
}

// parseSize maps named sizes or raw pixel counts to a tilesize value.
// Returns (0, "") if invalid.
func parseSize(input string) (int, string) {
	switch strings.ToLower(input) {
	case "small", "s":
		return 32, "small"
	case "medium", "m":
		return 64, "medium"
	case "large", "l":
		return 96, "large"
	default:
		val, err := strconv.Atoi(input)
		if err != nil || val < 16 || val > 256 {
			return 0, ""
		}
		return val, fmt.Sprintf("%dpx", val)
	}
}

// sizeLabel returns a human-readable label for a pixel tile size.
func sizeLabel(px int) string {
	switch {
	case px <= 40:
		return "small"
	case px <= 80:
		return "medium"
	default:
		return "large"
	}
}

// stateString returns "enabled" or "disabled" for a bool.
func stateString(on bool) string {
	if on {
		return "enabled"
	}
	return "disabled"
}
