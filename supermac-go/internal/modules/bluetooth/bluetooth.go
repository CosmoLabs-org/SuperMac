package bluetooth

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/cosmolabs-org/supermac/internal/module"
)

func init() {
	module.Register(&BluetoothModule{})
}

// BluetoothModule handles Bluetooth device management via blueutil.
type BluetoothModule struct{}

func (b *BluetoothModule) Name() string            { return "bluetooth" }
func (b *BluetoothModule) ShortDescription() string { return "Bluetooth device management" }
func (b *BluetoothModule) Emoji() string            { return "📡" }

func (b *BluetoothModule) Commands() []module.Command {
	return []module.Command{
		{
			Name:        "status",
			Description: "Show Bluetooth power state and connected devices",
			Run:         b.status,
		},
		{
			Name:        "devices",
			Description: "List all paired Bluetooth devices",
			Run:         b.devices,
		},
		{
			Name:        "connect",
			Description: "Connect to a paired device by MAC address",
			Args: []module.Arg{
				{Name: "mac", Required: true, Description: "MAC address of the device"},
			},
			Run: b.connect,
		},
		{
			Name:        "disconnect",
			Description: "Disconnect a Bluetooth device by MAC address",
			Args: []module.Arg{
				{Name: "mac", Required: true, Description: "MAC address of the device"},
			},
			Run: b.disconnect,
		},
		{
			Name:        "power",
			Description: "Toggle Bluetooth power (on, off, toggle). No args = show current state.",
			Args: []module.Arg{
				{Name: "state", Required: false, Description: "on, off, or toggle"},
			},
			Run: b.power,
		},
		{
			Name:        "discoverable",
			Description: "Set Bluetooth discoverable mode (on or off)",
			Args: []module.Arg{
				{Name: "state", Required: true, Description: "on or off"},
			},
			Run: b.discoverable,
		},
	}
}

func (b *BluetoothModule) Search(term string) []module.SearchResult {
	var results []module.SearchResult
	for _, cmd := range b.Commands() {
		if strings.Contains(cmd.Name, term) || strings.Contains(strings.ToLower(cmd.Description), term) {
			results = append(results, module.SearchResult{
				Command:     cmd.Name,
				Description: cmd.Description,
				Module:      b.Name(),
			})
		}
	}
	return results
}

// checkBlueutil verifies blueutil is installed and returns an error if not.
func checkBlueutil(ctx *module.Context) error {
	if _, err := exec.LookPath("blueutil"); err != nil {
		ctx.Output.Warning("blueutil is not installed")
		ctx.Output.Info("Install it with: brew install blueutil")
		return module.NewExitError(module.ExitGeneral,
			"blueutil is required for Bluetooth management. Install with: brew install blueutil")
	}
	return nil
}

// ---------------------------------------------------------------------------
// Status
// ---------------------------------------------------------------------------

func (b *BluetoothModule) status(ctx *module.Context) error {
	if err := checkBlueutil(ctx); err != nil {
		return err
	}

	// Power state
	powerOut, err := exec.Command("blueutil", "--power").Output()
	if err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to read Bluetooth power state: %v", err))
	}
	powerState := strings.TrimSpace(string(powerOut))

	ctx.Output.Header("Bluetooth Status")
	fmt.Println()

	if powerState == "1" {
		fmt.Println("  Power:         on")
	} else {
		fmt.Println("  Power:         off")
	}

	// Connected devices
	connectedOut, err := exec.Command("blueutil", "--connected").Output()
	if err != nil {
		// Not fatal — may simply have no connected devices
		fmt.Println("  Connected:     (none)")
		return nil
	}

	lines := strings.Split(strings.TrimSpace(string(connectedOut)), "\n")
	if len(lines) == 0 || (len(lines) == 1 && strings.TrimSpace(lines[0]) == "") {
		fmt.Println("  Connected:     (none)")
		return nil
	}

	fmt.Printf("  Connected:     %d device(s)\n", len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		name, mac := parseBlueutilDeviceLine(line)
		if name != "" {
			fmt.Printf("    - %s (%s)\n", name, mac)
		} else {
			fmt.Printf("    - %s\n", mac)
		}
	}

	return nil
}

// ---------------------------------------------------------------------------
// Devices
// ---------------------------------------------------------------------------

func (b *BluetoothModule) devices(ctx *module.Context) error {
	if err := checkBlueutil(ctx); err != nil {
		return err
	}

	out, err := exec.Command("blueutil", "--paired").Output()
	if err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to list paired devices: %v", err))
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")

	ctx.Output.Header("Paired Bluetooth Devices")
	fmt.Println()

	if len(lines) == 0 || (len(lines) == 1 && strings.TrimSpace(lines[0]) == "") {
		fmt.Println("  No paired devices found")
		return nil
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		name, mac := parseBlueutilDeviceLine(line)
		connected := strings.Contains(line, "connected")
		state := "disconnected"
		if connected {
			state = "connected"
		}

		if name != "" {
			fmt.Printf("  %-30s %s   %s\n", name, mac, state)
		} else {
			fmt.Printf("  %-30s %s   %s\n", "(unknown)", mac, state)
		}
	}

	return nil
}

// ---------------------------------------------------------------------------
// Connect / Disconnect
// ---------------------------------------------------------------------------

func (b *BluetoothModule) connect(ctx *module.Context) error {
	if err := checkBlueutil(ctx); err != nil {
		return err
	}

	if len(ctx.Args) == 0 {
		return module.NewExitError(module.ExitUsage, "MAC address required: mac bluetooth connect <mac>")
	}

	mac := ctx.Args[0]
	ctx.Output.Info("Connecting to %s...", mac)

	out, err := exec.Command("blueutil", "--connect", mac).CombinedOutput()
	if err != nil {
		return module.NewExitError(module.ExitGeneral,
			fmt.Sprintf("Failed to connect to %s: %s", mac, strings.TrimSpace(string(out))))
	}

	ctx.Output.Success("Connected to %s", mac)
	return nil
}

func (b *BluetoothModule) disconnect(ctx *module.Context) error {
	if err := checkBlueutil(ctx); err != nil {
		return err
	}

	if len(ctx.Args) == 0 {
		return module.NewExitError(module.ExitUsage, "MAC address required: mac bluetooth disconnect <mac>")
	}

	mac := ctx.Args[0]
	ctx.Output.Info("Disconnecting %s...", mac)

	out, err := exec.Command("blueutil", "--disconnect", mac).CombinedOutput()
	if err != nil {
		return module.NewExitError(module.ExitGeneral,
			fmt.Sprintf("Failed to disconnect %s: %s", mac, strings.TrimSpace(string(out))))
	}

	ctx.Output.Success("Disconnected %s", mac)
	return nil
}

// ---------------------------------------------------------------------------
// Power
// ---------------------------------------------------------------------------

func (b *BluetoothModule) power(ctx *module.Context) error {
	if err := checkBlueutil(ctx); err != nil {
		return err
	}

	// No args: show current state
	if len(ctx.Args) == 0 {
		out, err := exec.Command("blueutil", "--power").Output()
		if err != nil {
			return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to read Bluetooth power state: %v", err))
		}
		state := strings.TrimSpace(string(out))
		if state == "1" {
			ctx.Output.Info("Bluetooth power: on")
		} else {
			ctx.Output.Info("Bluetooth power: off")
		}
		return nil
	}

	arg := strings.ToLower(ctx.Args[0])
	switch arg {
	case "on":
		return setPower(ctx, true)
	case "off":
		return setPower(ctx, false)
	case "toggle":
		out, err := exec.Command("blueutil", "--power").Output()
		if err != nil {
			return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to read Bluetooth power state: %v", err))
		}
		current := strings.TrimSpace(string(out))
		return setPower(ctx, current != "1")
	default:
		return module.NewExitError(module.ExitUsage, "Argument must be on, off, or toggle")
	}
}

func setPower(ctx *module.Context, on bool) error {
	val := "0"
	label := "off"
	if on {
		val = "1"
		label = "on"
	}

	ctx.Output.Info("Setting Bluetooth power %s...", label)

	out, err := exec.Command("blueutil", "--power", val).CombinedOutput()
	if err != nil {
		return module.NewExitError(module.ExitGeneral,
			fmt.Sprintf("Failed to set Bluetooth power: %s", strings.TrimSpace(string(out))))
	}

	ctx.Output.Success("Bluetooth power: %s", label)
	return nil
}

// ---------------------------------------------------------------------------
// Discoverable
// ---------------------------------------------------------------------------

func (b *BluetoothModule) discoverable(ctx *module.Context) error {
	if err := checkBlueutil(ctx); err != nil {
		return err
	}

	if len(ctx.Args) == 0 {
		return module.NewExitError(module.ExitUsage, "Argument required: mac bluetooth discoverable <on|off>")
	}

	arg := strings.ToLower(ctx.Args[0])
	switch arg {
	case "on":
		return setDiscoverable(ctx, true)
	case "off":
		return setDiscoverable(ctx, false)
	default:
		return module.NewExitError(module.ExitUsage, "Argument must be on or off")
	}
}

func setDiscoverable(ctx *module.Context, on bool) error {
	val := "0"
	label := "off"
	if on {
		val = "1"
		label = "on"
	}

	ctx.Output.Info("Setting Bluetooth discoverable %s...", label)

	out, err := exec.Command("blueutil", "--discoverable", val).CombinedOutput()
	if err != nil {
		return module.NewExitError(module.ExitGeneral,
			fmt.Sprintf("Failed to set discoverable mode: %s", strings.TrimSpace(string(out))))
	}

	ctx.Output.Success("Bluetooth discoverable: %s", label)
	return nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// parseBlueutilDeviceLine extracts name and MAC from a blueutil output line.
// Typical format: "address: XX:XX:XX:XX:XX:XX, name: Device Name, ..."
// Fallback: just the MAC if parsing fails.
func parseBlueutilDeviceLine(line string) (name, mac string) {
	parts := strings.Split(line, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "address:") {
			mac = strings.TrimSpace(strings.TrimPrefix(part, "address:"))
		} else if strings.HasPrefix(part, "name:") {
			name = strings.TrimSpace(strings.TrimPrefix(part, "name:"))
		}
	}
	return name, mac
}
