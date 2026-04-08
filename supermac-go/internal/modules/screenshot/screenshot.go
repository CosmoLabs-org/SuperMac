package screenshot

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/cosmolabs-org/supermac/internal/module"
)

const domain = "com.apple.screencapture"

var validFormats = map[string]string{
	"png":  "png",
	"jpg":  "jpg",
	"jpeg": "jpg",
	"tiff": "tiff",
	"tif":  "tiff",
	"gif":  "gif",
}

func init() {
	module.Register(&ScreenshotModule{})
}

type ScreenshotModule struct{}

func (s *ScreenshotModule) Name() string                { return "screenshot" }
func (s *ScreenshotModule) ShortDescription() string     { return "Screenshot settings and management" }
func (s *ScreenshotModule) Emoji() string                { return "📸" }

func (s *ScreenshotModule) Commands() []module.Command {
	return []module.Command{
		{
			Name:        "location",
			Description: "Get or set screenshot save location",
			Aliases:     []string{"loc"},
			Args: []module.Arg{
				{Name: "destination", Required: false, Description: "desktop, downloads, clipboard, documents, pictures, or a custom path"},
			},
			Run: s.location,
		},
		{
			Name:        "format",
			Description: "Get or set screenshot file format",
			Aliases:     []string{"type"},
			Args: []module.Arg{
				{Name: "format", Required: false, Description: "png, jpg, tiff, or gif"},
			},
			Run: s.format,
		},
		{
			Name:        "shadow",
			Description: "Get or set window shadow in screenshots",
			Aliases:     []string{"shadows"},
			Args: []module.Arg{
				{Name: "action", Required: false, Description: "on, off, or toggle"},
			},
			Run: s.shadow,
		},
		{
			Name:        "cursor",
			Description: "Get or set cursor visibility in screenshots",
			Aliases:     []string{"show-cursor", "showsCursor"},
			Args: []module.Arg{
				{Name: "action", Required: false, Description: "show, hide, or toggle"},
			},
			Run: s.cursor,
		},
		{
			Name:        "naming",
			Description: "Get or set screenshot naming mode",
			Aliases:     []string{"name-format", "name"},
			Args: []module.Arg{
				{Name: "mode", Required: false, Description: "sequential or timestamp"},
			},
			Run: s.naming,
		},
		{
			Name:        "status",
			Description: "Show current screenshot settings",
			Run:         s.status,
		},
		{
			Name:        "thumbnail",
			Description: "Toggle screenshot thumbnail preview",
			Args: []module.Arg{
				{Name: "action", Required: false, Description: "on, off, or toggle"},
			},
			Run: s.thumbnail,
		},
		{
			Name:        "sound",
			Description: "Toggle screenshot camera shutter sound",
			Args: []module.Arg{
				{Name: "action", Required: false, Description: "on, off, or toggle"},
			},
			Run: s.sound,
		},
		{
			Name:        "take",
			Description: "Take a screenshot now (area, window, or screen)",
			Args: []module.Arg{
				{Name: "type", Required: false, Description: "area (default), window, or screen"},
			},
			Run: s.take,
		},
		{
			Name:        "reset",
			Description: "Reset all screenshot settings to defaults",
			Run:         s.reset,
		},
	}
}

func (s *ScreenshotModule) Search(term string) []module.SearchResult {
	var results []module.SearchResult
	for _, cmd := range s.Commands() {
		match := strings.Contains(cmd.Name, term) || strings.Contains(strings.ToLower(cmd.Description), term)
		if !match {
			for _, alias := range cmd.Aliases {
				if strings.Contains(alias, term) {
					match = true
					break
				}
			}
		}
		if match {
			results = append(results, module.SearchResult{
				Command:     cmd.Name,
				Description: cmd.Description,
				Module:      s.Name(),
			})
		}
	}
	return results
}

// --- Command implementations ---

func (s *ScreenshotModule) location(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		loc, err := readSetting(ctx, "location")
		if err != nil {
			ctx.Output.Info("Desktop (default)")
			return nil
		}
		loc = strings.TrimSpace(loc)
		if loc == "" {
			ctx.Output.Info("Desktop (default)")
		} else {
			ctx.Output.Info("%s", loc)
		}
		return nil
	}

	dest := ctx.Args[0]
	targetPath := ""

	switch strings.ToLower(dest) {
	case "desktop":
		targetPath = homeDir() + "/Desktop"
		ctx.Output.Info("Setting screenshot location to Desktop...")
	case "downloads":
		targetPath = homeDir() + "/Downloads"
		ctx.Output.Info("Setting screenshot location to Downloads...")
	case "clipboard":
		ctx.Output.Info("Setting screenshots to save to clipboard only...")
		if err := writeSetting(ctx, "location", "clipboard"); err != nil {
			return err
		}
		restartSystemUI(ctx)
		ctx.Output.Success("Screenshots will now be saved to clipboard only")
		return nil
	case "documents":
		targetPath = homeDir() + "/Documents"
		ctx.Output.Info("Setting screenshot location to Documents...")
	case "pictures":
		targetPath = homeDir() + "/Pictures"
		ctx.Output.Info("Setting screenshot location to Pictures...")
	default:
		// Custom path
		abs, err := filepath.Abs(dest)
		if err != nil {
			return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Invalid path: %s", dest))
		}
		info, err := os.Stat(abs)
		if err != nil || !info.IsDir() {
			return module.NewExitError(module.ExitNotFound, fmt.Sprintf("Directory does not exist: %s", dest))
		}
		targetPath = abs
		ctx.Output.Info("Setting screenshot location to: %s", dest)
	}

	if err := writeSetting(ctx, "location", targetPath); err != nil {
		return err
	}
	restartSystemUI(ctx)
	ctx.Output.Success("%s", fmt.Sprintf("Screenshot location set to: %s", targetPath))
	return nil
}

func (s *ScreenshotModule) format(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		val, err := readSetting(ctx, "type")
		if err != nil || val == "" {
			ctx.Output.Info("PNG (default)")
		} else {
			ctx.Output.Info("%s", strings.ToUpper(strings.TrimSpace(val)))
		}
		return nil
	}

	input := strings.ToLower(ctx.Args[0])
	canonical, ok := validFormats[input]
	if !ok {
		return module.NewExitError(module.ExitUsage, fmt.Sprintf("Invalid format: %s. Available: png, jpg, tiff, gif", input))
	}

	ctx.Output.Info("Setting screenshot format to %s...", strings.ToUpper(canonical))
	if err := writeSetting(ctx, "type", canonical); err != nil {
		return err
	}
	restartSystemUI(ctx)

	descriptions := map[string]string{
		"png":  "lossless, best quality",
		"jpg":  "smaller files, some compression",
		"tiff": "lossless, larger files",
		"gif":  "limited colors, smaller files",
	}
	desc := descriptions[canonical]
	ctx.Output.Success("%s", fmt.Sprintf("Format set to %s (%s)", strings.ToUpper(canonical), desc))
	return nil
}

func (s *ScreenshotModule) shadow(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		val, err := readSetting(ctx, "disable-shadow")
		if err != nil {
			ctx.Output.Info("enabled (default)")
			return nil
		}
		if isTrue(val) {
			ctx.Output.Info("disabled")
		} else {
			ctx.Output.Info("enabled")
		}
		return nil
	}

	action := strings.ToLower(ctx.Args[0])
	switch action {
	case "on", "enable", "true":
		ctx.Output.Info("Enabling window shadows in screenshots...")
		if err := writeSetting(ctx, "disable-shadow", "-bool false"); err != nil {
			return err
		}
		restartSystemUI(ctx)
		ctx.Output.Success("Window shadows enabled")
	case "off", "disable", "false":
		ctx.Output.Info("Disabling window shadows in screenshots...")
		if err := writeSetting(ctx, "disable-shadow", "-bool true"); err != nil {
			return err
		}
		restartSystemUI(ctx)
		ctx.Output.Success("Window shadows disabled")
	case "toggle":
		val, err := readSetting(ctx, "disable-shadow")
		if err != nil {
			val = "false"
		}
		if isTrue(val) {
			return s.shadowWithArg(ctx, "on")
		}
		return s.shadowWithArg(ctx, "off")
	default:
		return module.NewExitError(module.ExitUsage, fmt.Sprintf("Invalid action: %s. Use: on, off, toggle", action))
	}
	return nil
}

func (s *ScreenshotModule) shadowWithArg(ctx *module.Context, action string) error {
	return s.shadow(&module.Context{
		Config:   ctx.Config,
		Output:   ctx.Output,
		Platform: ctx.Platform,
		Prompt:   ctx.Prompt,
		Args:     []string{action},
		Flags:    ctx.Flags,
		Verbose:  ctx.Verbose,
		DryRun:   ctx.DryRun,
	})
}

func (s *ScreenshotModule) cursor(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		val, err := readSetting(ctx, "showsCursor")
		if err != nil {
			ctx.Output.Info("hidden (default)")
			return nil
		}
		if isTrue(val) {
			ctx.Output.Info("visible")
		} else {
			ctx.Output.Info("hidden")
		}
		return nil
	}

	action := strings.ToLower(ctx.Args[0])
	switch action {
	case "show", "on", "enable", "true":
		ctx.Output.Info("Enabling cursor in screenshots...")
		if err := writeSetting(ctx, "showsCursor", "-bool true"); err != nil {
			return err
		}
		restartSystemUI(ctx)
		ctx.Output.Success("Cursor will appear in screenshots")
	case "hide", "off", "disable", "false":
		ctx.Output.Info("Disabling cursor in screenshots...")
		if err := writeSetting(ctx, "showsCursor", "-bool false"); err != nil {
			return err
		}
		restartSystemUI(ctx)
		ctx.Output.Success("Cursor will not appear in screenshots")
	case "toggle":
		val, err := readSetting(ctx, "showsCursor")
		if err != nil {
			val = "false"
		}
		if isTrue(val) {
			return s.cursorWithArg(ctx, "hide")
		}
		return s.cursorWithArg(ctx, "show")
	default:
		return module.NewExitError(module.ExitUsage, fmt.Sprintf("Invalid action: %s. Use: show, hide, toggle", action))
	}
	return nil
}

func (s *ScreenshotModule) cursorWithArg(ctx *module.Context, action string) error {
	return s.cursor(&module.Context{
		Config:   ctx.Config,
		Output:   ctx.Output,
		Platform: ctx.Platform,
		Prompt:   ctx.Prompt,
		Args:     []string{action},
		Flags:    ctx.Flags,
		Verbose:  ctx.Verbose,
		DryRun:   ctx.DryRun,
	})
}

func (s *ScreenshotModule) naming(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		val, err := readSetting(ctx, "name")
		if err != nil || val == "" {
			ctx.Output.Info("Default (Screenshot date-time)")
		} else {
			ctx.Output.Info("%s", strings.TrimSpace(val))
		}
		return nil
	}

	input := ctx.Args[0]
	switch strings.ToLower(input) {
	case "sequential":
		ctx.Output.Info("Setting naming to sequential mode...")
		if err := writeSetting(ctx, "name", "Screenshot"); err != nil {
			return err
		}
		restartSystemUI(ctx)
		ctx.Output.Success("Screenshots will use sequential naming (Screenshot, Screenshot 1, Screenshot 2...)")
	case "timestamp":
		ctx.Output.Info("Setting naming to timestamp mode...")
		if err := writeSetting(ctx, "name", "Screenshot %Y-%m-%d at %H.%M.%S"); err != nil {
			return err
		}
		restartSystemUI(ctx)
		ctx.Output.Success("Screenshots will use timestamp naming")
	default:
		// Allow custom format strings
		ctx.Output.Info("Setting screenshot name format...")
		if err := writeSetting(ctx, "name", input); err != nil {
			return err
		}
		restartSystemUI(ctx)
		ctx.Output.Success("%s", fmt.Sprintf("Screenshot name format set to: %s", input))
	}
	return nil
}

func (s *ScreenshotModule) status(ctx *module.Context) error {
	ctx.Output.Header("Screenshot Settings")

	// Location
	loc, err := readSetting(ctx, "location")
	if err != nil || strings.TrimSpace(loc) == "" {
		fmt.Println("  Location:         Desktop (default)")
	} else if strings.TrimSpace(loc) == "clipboard" {
		fmt.Println("  Location:         Clipboard only")
	} else {
		fmt.Printf("  Location:         %s\n", strings.TrimSpace(loc))
	}

	// Format
	fmtVal, err := readSetting(ctx, "type")
	if err != nil || strings.TrimSpace(fmtVal) == "" {
		fmt.Println("  Format:           PNG (default)")
	} else {
		fmt.Printf("  Format:           %s\n", strings.ToUpper(strings.TrimSpace(fmtVal)))
	}

	// Shadow
	shadowVal, err := readSetting(ctx, "disable-shadow")
	if err != nil {
		fmt.Println("  Window shadows:   enabled (default)")
	} else if isTrue(shadowVal) {
		fmt.Println("  Window shadows:   disabled")
	} else {
		fmt.Println("  Window shadows:   enabled")
	}

	// Cursor
	cursorVal, err := readSetting(ctx, "showsCursor")
	if err != nil {
		fmt.Println("  Show cursor:      hidden (default)")
	} else if isTrue(cursorVal) {
		fmt.Println("  Show cursor:      visible")
	} else {
		fmt.Println("  Show cursor:      hidden")
	}

	// Naming
	nameVal, err := readSetting(ctx, "name")
	if err != nil || strings.TrimSpace(nameVal) == "" {
		fmt.Println("  Name format:      Default (Screenshot date-time)")
	} else {
		fmt.Printf("  Name format:      %s\n", strings.TrimSpace(nameVal))
	}

	fmt.Println("")
	fmt.Println("  Keyboard shortcuts:")
	fmt.Println("    Cmd+Shift+3    Full screen screenshot")
	fmt.Println("    Cmd+Shift+4    Select area screenshot")
	fmt.Println("    Cmd+Shift+5    Screenshot options menu")

	return nil
}

func (s *ScreenshotModule) reset(ctx *module.Context) error {
	ctx.Output.Info("Resetting screenshot settings to defaults...")

	cmd := exec.Command("defaults", "delete", domain)
	_ = cmd.Run() // Best effort — may fail if no overrides exist

	restartSystemUI(ctx)
	ctx.Output.Success("Screenshot settings reset to defaults")
	ctx.Output.Info("Screenshots will now save to Desktop in PNG format")
	return nil
}

// --- Thumbnail ---

func (s *ScreenshotModule) thumbnail(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		val, err := readSetting(ctx, "include-thumbnail")
		if err != nil || val == "" {
			ctx.Output.Info("Thumbnail preview: enabled (default)")
		} else if isTrue(val) {
			ctx.Output.Info("Thumbnail preview: enabled")
		} else {
			ctx.Output.Info("Thumbnail preview: disabled")
		}
		return nil
	}

	action := strings.ToLower(ctx.Args[0])
	current, _ := readSetting(ctx, "include-thumbnail")
	currentOn := isTrue(current)

	var newVal string
	switch action {
	case "on":
		newVal = "1"
	case "off":
		newVal = "0"
	case "toggle":
		if currentOn {
			newVal = "0"
		} else {
			newVal = "1"
		}
	default:
		return module.NewExitError(module.ExitUsage, "Usage: mac screenshot thumbnail [on|off|toggle]")
	}

	ctx.Output.Info("Setting thumbnail preview to %s...", newVal)
	if err := writeSetting(ctx, "include-thumbnail", "-bool "+newVal); err != nil {
		return err
	}
	restartSystemUI(ctx)
	if newVal == "1" {
		ctx.Output.Success("Thumbnail preview enabled")
	} else {
		ctx.Output.Success("Thumbnail preview disabled")
	}
	return nil
}

// --- Sound ---

func (s *ScreenshotModule) sound(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		soundVal, soundErr := ctx.Platform.ReadDefault(domain, "disable-sound")
		if soundErr != nil || soundVal == "" {
			ctx.Output.Info("Camera sound: enabled (default)")
		} else if isTrue(soundVal) {
			ctx.Output.Info("Camera sound: disabled")
		} else {
			ctx.Output.Info("Camera sound: enabled")
		}
		_ = soundErr
		return nil
	}

	action := strings.ToLower(ctx.Args[0])
	soundVal, _ := ctx.Platform.ReadDefault(domain, "disable-sound")
	currentDisabled := isTrue(soundVal)

	var newVal string
	switch action {
	case "on":
		newVal = "0" // disable-sound = 0 means sound ON
	case "off":
		newVal = "1"
	case "toggle":
		if currentDisabled {
			newVal = "0"
		} else {
			newVal = "1"
		}
	default:
		return module.NewExitError(module.ExitUsage, "Usage: mac screenshot sound [on|off|toggle]")
	}

	ctx.Output.Info("Setting camera sound...")
	if err := ctx.Platform.WriteDefault(domain, "disable-sound", "-bool "+newVal); err != nil {
		return err
	}
	restartSystemUI(ctx)
	if newVal == "0" {
		ctx.Output.Success("Camera shutter sound enabled")
	} else {
		ctx.Output.Success("Camera shutter sound disabled")
	}
	return nil
}

// --- Take ---

func (s *ScreenshotModule) take(ctx *module.Context) error {
	captureType := "area"
	if len(ctx.Args) > 0 {
		captureType = strings.ToLower(ctx.Args[0])
	}

	loc, _ := readSetting(ctx, "location")
	if loc == "" {
		loc, _ = os.UserHomeDir()
		loc = filepath.Join(loc, "Desktop")
	}

	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	filename := filepath.Join(loc, fmt.Sprintf("Screenshot %s.png", timestamp))

	switch captureType {
	case "screen", "fullscreen":
		ctx.Output.Info("Taking full screen screenshot...")
		if err := exec.Command("screencapture", "-x", filename).Run(); err != nil {
			return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to take screenshot: %v", err))
		}
		ctx.Output.Success("Full screen screenshot taken!")

	case "area", "selection":
		ctx.Output.Info("Select area for screenshot (press Space for window mode)...")
		if err := exec.Command("screencapture", "-i", filename).Run(); err != nil {
			return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to take screenshot: %v", err))
		}
		ctx.Output.Success("Area screenshot taken!")

	case "window":
		ctx.Output.Info("Click on a window to capture...")
		if err := exec.Command("screencapture", "-i", "-w", filename).Run(); err != nil {
			return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to take screenshot: %v", err))
		}
		ctx.Output.Success("Window screenshot taken!")

	default:
		return module.NewExitError(module.ExitUsage, "Invalid type. Use: area, window, or screen")
	}

	ctx.Output.Info("Saved to: %s", filename)
	return nil
}

// --- Helpers ---

func readSetting(ctx *module.Context, key string) (string, error) {
	return ctx.Platform.ReadDefault(domain, key)
}

func writeSetting(ctx *module.Context, key, value string) error {
	return ctx.Platform.WriteDefault(domain, key, value)
}

func restartSystemUI(ctx *module.Context) {
	cmd := exec.Command("killall", "SystemUIServer")
	_ = cmd.Run()
}

func isTrue(val string) bool {
	v := strings.TrimSpace(val)
	return v == "true" || v == "1" || v == "YES"
}

func homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "/tmp"
	}
	return home
}
