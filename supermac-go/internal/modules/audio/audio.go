package audio

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/cosmolabs-org/supermac/internal/module"
	"github.com/cosmolabs-org/supermac/internal/platform"
)

func init() {
	module.Register(&AudioModule{})
}

// AudioModule handles audio control and device management.
type AudioModule struct{}

func (a *AudioModule) Name() string            { return "audio" }
func (a *AudioModule) ShortDescription() string { return "Audio control and device management" }
func (a *AudioModule) Emoji() string            { return "🔊" }

func (a *AudioModule) Commands() []module.Command {
	return []module.Command{
		{
			Name:        "volume",
			Description: "Get or set system volume (0-100)",
			Args: []module.Arg{
				{Name: "level", Required: false, Description: "Volume level (0-100). Omit to get current volume."},
			},
			Run: a.volume,
		},
		{
			Name:        "up",
			Description: "Increase volume by step (default: 10)",
			Args: []module.Arg{
				{Name: "step", Required: false, Description: "Step amount (default: 10)"},
			},
			Run: a.volumeUp,
		},
		{
			Name:        "down",
			Description: "Decrease volume by step (default: 10)",
			Args: []module.Arg{
				{Name: "step", Required: false, Description: "Step amount (default: 10)"},
			},
			Run: a.volumeDown,
		},
		{
			Name:        "mute",
			Description: "Mute system audio",
			Run:         a.mute,
		},
		{
			Name:        "unmute",
			Description: "Unmute system audio",
			Run:         a.unmute,
		},
		{
			Name:        "toggle-mute",
			Description: "Toggle mute state",
			Run:         a.toggleMute,
		},
		{
			Name:        "devices",
			Description: "List audio devices (all/input/output)",
			Args: []module.Arg{
				{Name: "type", Required: false, Description: "Device type: all (default), input, output"},
			},
			Run: a.devices,
		},
		{
			Name:        "input-device",
			Description: "Set audio input device",
			Args: []module.Arg{
				{Name: "name", Required: true, Description: "Device name"},
			},
			Run: a.inputDevice,
		},
		{
			Name:        "output-device",
			Description: "Set audio output device",
			Args: []module.Arg{
				{Name: "name", Required: true, Description: "Device name"},
			},
			Run: a.outputDevice,
		},
		{
			Name:        "status",
			Description: "Show audio status",
			Run:         a.status,
		},
	}
}

func (a *AudioModule) Search(term string) []module.SearchResult {
	var results []module.SearchResult
	for _, cmd := range a.Commands() {
		if strings.Contains(cmd.Name, term) || strings.Contains(strings.ToLower(cmd.Description), term) {
			results = append(results, module.SearchResult{
				Command:     cmd.Name,
				Description: cmd.Description,
				Module:      a.Name(),
			})
		}
	}
	return results
}

// ---------------------------------------------------------------------------
// Volume
// ---------------------------------------------------------------------------

func (a *AudioModule) volume(ctx *module.Context) error {
	// No argument: display current volume.
	if len(ctx.Args) == 0 {
		vol, err := getVolume(ctx)
		if err != nil {
			return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to get volume: %v", err))
		}
		muted, _ := getMuteState(ctx)
		ctx.Output.Info("Volume: %d%% (%s)", vol, muted)
		return nil
	}

	target, err := strconv.Atoi(ctx.Args[0])
	if err != nil || target < 0 || target > 100 {
		return module.NewExitError(module.ExitUsage, "Volume must be between 0 and 100")
	}

	current, _ := getVolume(ctx)
	ctx.Output.Info("Setting volume from %d%% to %d%%...", current, target)

	if _, err := ctx.Platform.RunOSAScript(fmt.Sprintf("set volume output volume %d", target)); err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to set volume: %v", err))
	}

	ctx.Output.Success("Volume set to %d%%", target)

	switch {
	case target == 0:
		ctx.Output.Info("Audio is now silent")
	case target <= 25:
		ctx.Output.Info("Low volume")
	case target <= 75:
		ctx.Output.Info("Medium volume")
	default:
		ctx.Output.Info("High volume")
	}
	return nil
}

func (a *AudioModule) volumeUp(ctx *module.Context) error {
	step := 10
	if len(ctx.Args) > 0 {
		if s, err := strconv.Atoi(ctx.Args[0]); err == nil && s > 0 {
			step = s
		}
	}

	current, err := getVolume(ctx)
	if err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to get volume: %v", err))
	}

	newVol := current + step
	if newVol > 100 {
		newVol = 100
	}
	return setVolume(ctx, newVol)
}

func (a *AudioModule) volumeDown(ctx *module.Context) error {
	step := 10
	if len(ctx.Args) > 0 {
		if s, err := strconv.Atoi(ctx.Args[0]); err == nil && s > 0 {
			step = s
		}
	}

	current, err := getVolume(ctx)
	if err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to get volume: %v", err))
	}

	newVol := current - step
	if newVol < 0 {
		newVol = 0
	}
	return setVolume(ctx, newVol)
}

// ---------------------------------------------------------------------------
// Mute / Unmute / Toggle
// ---------------------------------------------------------------------------

func (a *AudioModule) mute(ctx *module.Context) error {
	muted, _ := getMuteState(ctx)
	if muted == "muted" {
		ctx.Output.Info("Audio is already muted")
		return nil
	}

	ctx.Output.Info("Muting audio...")
	if _, err := ctx.Platform.RunOSAScript("set volume with output muted"); err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to mute audio: %v", err))
	}
	ctx.Output.Success("Audio muted")
	ctx.Output.Info("Use 'mac audio unmute' to restore sound")
	return nil
}

func (a *AudioModule) unmute(ctx *module.Context) error {
	muted, _ := getMuteState(ctx)
	if muted == "unmuted" {
		ctx.Output.Info("Audio is already unmuted")
		return nil
	}

	ctx.Output.Info("Unmuting audio...")
	if _, err := ctx.Platform.RunOSAScript("set volume without output muted"); err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to unmute audio: %v", err))
	}
	vol, _ := getVolume(ctx)
	ctx.Output.Success("Audio unmuted (Volume: %d%%)", vol)
	return nil
}

func (a *AudioModule) toggleMute(ctx *module.Context) error {
	muted, _ := getMuteState(ctx)
	ctx.Output.Info("Current state: Audio is %s", muted)
	if muted == "muted" {
		return a.unmute(ctx)
	}
	return a.mute(ctx)
}

// ---------------------------------------------------------------------------
// Devices
// ---------------------------------------------------------------------------

func (a *AudioModule) devices(ctx *module.Context) error {
	deviceType := "all"
	if len(ctx.Args) > 0 {
		deviceType = strings.ToLower(ctx.Args[0])
	}

	devices, err := ctx.Platform.GetAudioDevices()
	if err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to list audio devices: %v", err))
	}

	ctx.Output.Header("Audio Devices")

	var outputDevs, inputDevs []platform.AudioDevice //nolint:prealloc // split from shared slice
	for _, d := range devices {
		switch d.Type {
		case "output":
			outputDevs = append(outputDevs, d)
		case "input":
			inputDevs = append(inputDevs, d)
		}
	}

	switch deviceType {
	case "output", "out":
		printDeviceList(ctx, "Output Devices", outputDevs)
	case "input", "in":
		printDeviceList(ctx, "Input Devices", inputDevs)
	default:
		printDeviceList(ctx, "Output Devices", outputDevs)
		fmt.Println()
		printDeviceList(ctx, "Input Devices", inputDevs)
	}

	// Hint about SwitchAudioSource for enhanced switching.
	if _, err := exec.LookPath("SwitchAudioSource"); err != nil {
		fmt.Println()
		ctx.Output.Info("Install SwitchAudioSource for enhanced device switching:")
		ctx.Output.Info("    brew install switchaudio-osx")
	}
	return nil
}

func printDeviceList(ctx *module.Context, header string, devices []platform.AudioDevice) {
	ctx.Output.Info("%s:", header)
	if len(devices) == 0 {
		fmt.Println("  (none detected)")
		return
	}
	for _, d := range devices {
		marker := "  "
		if d.Active {
			marker = "* "
		}
		fmt.Printf("  %s%s\n", marker, d.Name)
	}
}

// ---------------------------------------------------------------------------
// Input / Output device switching
// ---------------------------------------------------------------------------

func (a *AudioModule) inputDevice(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		return module.NewExitError(module.ExitUsage, "Device name required: mac audio input-device <name>")
	}
	return switchDevice(ctx, "input", ctx.Args[0])
}

func (a *AudioModule) outputDevice(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		return module.NewExitError(module.ExitUsage, "Device name required: mac audio output-device <name>")
	}
	return switchDevice(ctx, "output", ctx.Args[0])
}

func switchDevice(ctx *module.Context, deviceType, name string) error {
	// Prefer SwitchAudioSource when available.
	if _, err := exec.LookPath("SwitchAudioSource"); err == nil {
		ctx.Output.Info("Setting %s device to: %s", deviceType, name)
		out, err := exec.Command("SwitchAudioSource", "-s", name, "-t", deviceType).CombinedOutput() //nolint:gosec // device name is user-provided, not a security risk
		if err != nil {
			return module.NewExitError(module.ExitGeneral, fmt.Sprintf(
				"Failed to set %s device: %s\nMake sure the device name is correct. Use 'mac audio devices %s' to see available devices.",
				deviceType, strings.TrimSpace(string(out)), deviceType))
		}
		ctx.Output.Success("%s device set to: %s", strings.Title(deviceType), name)
		return nil
	}

	// Graceful fallback: try via platform RunCommand if SwitchAudioSource is absent.
	ctx.Output.Info("SwitchAudioSource not found, attempting fallback...")
	ctx.Output.Info("Install SwitchAudioSource for reliable device switching: brew install switchaudio-osx")
	return module.NewExitError(module.ExitNotFound,
		"SwitchAudioSource is required for device switching. Install with: brew install switchaudio-osx")
}

// ---------------------------------------------------------------------------
// Status
// ---------------------------------------------------------------------------

func (a *AudioModule) status(ctx *module.Context) error {
	ctx.Output.Header("Audio Status")
	fmt.Println()

	vol, err := getVolume(ctx)
	if err != nil {
		ctx.Output.Warning("Could not read volume")
		vol = -1
	}
	muted, _ := getMuteState(ctx)

	fmt.Printf("  Volume:        %d%%\n", vol)
	fmt.Printf("  Status:        %s\n", muted)

	// Current devices via platform.
	devices, _ := ctx.Platform.GetAudioDevices()
	for _, d := range devices {
		if d.Active {
			fmt.Printf("  %s device:   %s\n", strings.Title(d.Type), d.Name)
		}
	}

	// Sound effects.
	soundEffects, err := ctx.Platform.ReadDefault("NSGlobalDomain", "com.apple.sound.uiaudio.enabled")
	if err == nil {
		if soundEffects == "1" {
			fmt.Println("  Sound effects: enabled")
		} else {
			fmt.Println("  Sound effects: disabled")
		}
	}

	// Alert volume.
	alertOut, err := ctx.Platform.RunOSAScript("alert volume of (get volume settings)")
	if err == nil {
		fmt.Printf("  Alert volume:  %s%%\n", strings.TrimSpace(alertOut))
	}
	return nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func getVolume(ctx *module.Context) (int, error) {
	out, err := ctx.Platform.RunOSAScript("output volume of (get volume settings)")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(out))
}

func getMuteState(ctx *module.Context) (string, error) {
	out, err := ctx.Platform.RunOSAScript("output muted of (get volume settings)")
	if err != nil {
		return "unknown", err
	}
	if strings.TrimSpace(out) == "true" {
		return "muted", nil
	}
	return "unmuted", nil
}

func setVolume(ctx *module.Context, level int) error {
	ctx.Output.Info("Setting volume to %d%%...", level)
	if _, err := ctx.Platform.RunOSAScript(fmt.Sprintf("set volume output volume %d", level)); err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to set volume: %v", err))
	}
	ctx.Output.Success("Volume set to %d%%", level)
	return nil
}
