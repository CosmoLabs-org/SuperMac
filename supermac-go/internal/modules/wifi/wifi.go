package wifi

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/cosmolabs-org/supermac/internal/module"
	"github.com/cosmolabs-org/supermac/internal/platform"
)

const airportBin = "/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport"

func init() {
	module.Register(&WiFiModule{})
}

type WiFiModule struct{}

func (w *WiFiModule) Name() string            { return "wifi" }
func (w *WiFiModule) ShortDescription() string { return "WiFi control and management" }
func (w *WiFiModule) Emoji() string            { return "\U0001F310" }

func (w *WiFiModule) Commands() []module.Command {
	return []module.Command{
		{
			Name:        "on",
			Description: "Turn WiFi on",
			Run:         w.on,
		},
		{
			Name:        "off",
			Description: "Turn WiFi off",
			Run:         w.off,
		},
		{
			Name:        "toggle",
			Description: "Toggle WiFi state",
			Run:         w.toggle,
		},
		{
			Name:        "status",
			Description: "Show WiFi connection status",
			Run:         w.status,
		},
		{
			Name:        "scan",
			Description: "Scan for available networks",
			Run:         w.scan,
		},
		{
			Name:        "connect",
			Description: "Connect to a WiFi network",
			Aliases:     []string{"join"},
			Args: []module.Arg{
				{Name: "network", Required: true, Description: "Network name to connect to"},
				{Name: "password", Required: false, Description: "Network password"},
			},
			Run: w.connect,
		},
		{
			Name:        "forget",
			Description: "Forget a saved WiFi network",
			Args: []module.Arg{
				{Name: "network", Required: true, Description: "Network name to forget"},
			},
			Run: w.forget,
		},
		{
			Name:        "info",
			Description: "Detailed connection information",
			Run:         w.info,
		},
		{
			Name:        "list-saved",
			Description: "Show saved WiFi networks",
			Aliases:     []string{"saved"},
			Run:         w.listSaved,
		},
	}
}

func (w *WiFiModule) Search(term string) []module.SearchResult {
	var results []module.SearchResult
	for _, cmd := range w.Commands() {
		if strings.Contains(cmd.Name, term) || strings.Contains(strings.ToLower(cmd.Description), term) {
			results = append(results, module.SearchResult{
				Command:     cmd.Name,
				Description: cmd.Description,
				Module:      w.Name(),
			})
		}
	}
	// Also match common synonyms not covered by command names/descriptions.
	synonyms := map[string][]module.SearchResult{
		"wireless": {{Command: "status", Description: "Show WiFi connection status", Module: w.Name()}},
		"network": {
			{Command: "scan", Description: "Scan for available networks", Module: w.Name()},
			{Command: "connect", Description: "Connect to a WiFi network", Module: w.Name()},
		},
		"remove": {{Command: "forget", Description: "Forget a saved WiFi network", Module: w.Name()}},
	}
	if matches, ok := synonyms[term]; ok {
		results = append(results, matches...)
	}
	return results
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// getInterface returns the WiFi interface name (typically en0).
func getInterface() (string, error) {
	out, err := exec.Command("networksetup", "-listallhardwareports").Output()
	if err != nil {
		return "", module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to list hardware ports: %v", err))
	}
	lines := strings.Split(string(out), "\n")
	for i, line := range lines {
		if strings.Contains(line, "Wi-Fi") || strings.Contains(line, "AirPort") {
			if i+1 < len(lines) {
				fields := strings.Fields(lines[i+1])
				if len(fields) >= 2 {
					return fields[1], nil
				}
			}
		}
	}
	return "", module.NewExitError(module.ExitNotFound, "WiFi interface not found")
}

// getPowerState returns the current WiFi power state ("On", "Off", or error).
func getPowerState() (string, error) {
	iface, err := getInterface()
	if err != nil {
		return "unknown", err
	}
	out, err := exec.Command("networksetup", "-getairportpower", iface).Output()
	if err != nil {
		return "unknown", err
	}
	parts := strings.Fields(string(out))
	if len(parts) == 0 {
		return "unknown", nil
	}
	return parts[len(parts)-1], nil
}

// getCurrentNetwork returns the SSID of the currently connected network,
// or "Not connected" if none.
func getCurrentNetwork() (string, error) {
	iface, err := getInterface()
	if err != nil {
		return "", err
	}
	out, err := exec.Command("networksetup", "-getairportnetwork", iface).Output()
	if err != nil {
		return "Not connected", nil
	}
	trimmed := strings.TrimSpace(string(out))
	if strings.Contains(trimmed, "You are not associated") {
		return "Not connected", nil
	}
	return strings.TrimPrefix(trimmed, "Current Wi-Fi Network: "), nil
}

// parseAirportInfo parses airport -I output into a map of key-value pairs.
func parseAirportInfo(raw string) map[string]string {
	result := make(map[string]string)
	for _, line := range strings.Split(raw, "\n") {
		idx := strings.Index(line, ":")
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		result[key] = val
	}
	return result
}

// waitForPower polls the WiFi power state until it matches want or timeout.
func waitForPower(want string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		time.Sleep(500 * time.Millisecond)
		state, err := getPowerState()
		if err == nil && state == want {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// Commands
// ---------------------------------------------------------------------------

func (w *WiFiModule) on(ctx *module.Context) error {
	iface, err := getInterface()
	if err != nil {
		return err
	}

	state, _ := getPowerState()
	if state == "On" {
		ctx.Output.Info("WiFi is already on")
		return nil
	}

	ctx.Output.Info("Turning WiFi on...")
	if err := ctx.Platform.SetWiFi(true); err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to turn WiFi on: %v", err))
	}

	if !waitForPower("On", 10*time.Second) {
		return module.NewExitError(module.ExitGeneral, "Failed to turn WiFi on (timeout)")
	}

	ctx.Output.Success("WiFi is now on")

	// Brief pause to let the interface scan, then show current connection.
	time.Sleep(2 * time.Second)
	showCurrentConnection(ctx, iface)
	return nil
}

func (w *WiFiModule) off(ctx *module.Context) error {
	if _, err := getInterface(); err != nil {
		return err
	}

	state, _ := getPowerState()
	if state == "Off" {
		ctx.Output.Info("WiFi is already off")
		return nil
	}

	ctx.Output.Info("Turning WiFi off...")
	if err := ctx.Platform.SetWiFi(false); err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to turn WiFi off: %v", err))
	}

	time.Sleep(2 * time.Second)
	state, _ = getPowerState()
	if state == "Off" {
		ctx.Output.Success("WiFi is now off")
		ctx.Output.Info("Use 'mac wifi on' to turn it back on")
	} else {
		return module.NewExitError(module.ExitGeneral, "Failed to turn WiFi off")
	}
	return nil
}

func (w *WiFiModule) toggle(ctx *module.Context) error {
	state, err := getPowerState()
	if err != nil {
		return err
	}
	ctx.Output.Info("Current WiFi state: %s", state)
	switch state {
	case "On":
		return w.off(ctx)
	case "Off":
		return w.on(ctx)
	default:
		return module.NewExitError(module.ExitGeneral, "Unable to determine WiFi state")
	}
}

func (w *WiFiModule) status(ctx *module.Context) error {
	ctx.Output.Header("\U0001F310 WiFi Status")

	iface, err := getInterface()
	if err != nil {
		return err
	}

	fmt.Printf("  Interface:      %s\n", iface)

	state, _ := getPowerState()
	fmt.Printf("  Power:          %s\n", state)

	if state != "On" {
		return nil
	}

	network, _ := getCurrentNetwork()
	if network == "Not connected" {
		ctx.Output.Warning("Not connected to any WiFi network")
		return nil
	}
	fmt.Printf("  Connected to:   %s\n", network)

	// Signal strength via airport -I
	info, err := exec.Command(airportBin, "-I").Output()
	if err == nil {
		parsed := parseAirportInfo(string(info))
		if rssi, ok := parsed["agrCtlRSSI"]; ok {
			fmt.Printf("  Signal:         %s dBm\n", rssi)
		}
	}

	// IP address
	ip, err := exec.Command("ipconfig", "getifaddr", iface).Output()
	if err == nil {
		fmt.Printf("  IP Address:     %s\n", strings.TrimSpace(string(ip)))
	}

	return nil
}

func (w *WiFiModule) scan(ctx *module.Context) error {
	iface, err := getInterface()
	if err != nil {
		return err
	}

	state, _ := getPowerState()
	if state != "On" {
		return module.NewExitError(module.ExitGeneral, "WiFi is turned off. Turn on WiFi first: mac wifi on")
	}

	ctx.Output.Info("Scanning for available WiFi networks...")

	// Try airport scan first via platform interface.
	networks, err := ctx.Platform.ScanWiFiNetworks()
	if err == nil && len(networks) > 0 {
		for _, n := range networks {
			fmt.Printf("  %-32s %d dBm  %s\n", n.SSID, n.Signal, n.Security)
		}
		return nil
	}

	// Fallback: show saved networks.
	ctx.Output.Info("Airport scan unavailable, listing saved networks...")
	out, err := exec.Command("networksetup", "-listpreferredwirelessnetworks", iface).Output()
	if err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to list networks: %v", err))
	}
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "Preferred networks") {
			continue
		}
		fmt.Printf("  %s\n", line)
	}
	return nil
}

func (w *WiFiModule) connect(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		return module.NewExitError(module.ExitUsage, "Network name required: mac wifi connect <network> [password]")
	}
	networkName := ctx.Args[0]

	iface, err := getInterface()
	if err != nil {
		return err
	}

	state, _ := getPowerState()
	if state != "On" {
		return module.NewExitError(module.ExitGeneral, "WiFi is turned off. Turn on WiFi first: mac wifi on")
	}

	ctx.Output.Info("Connecting to network: %s", networkName)

	var cmd *exec.Cmd
	if len(ctx.Args) >= 2 {
		// Connect with password.
		password := ctx.Args[1]
		cmd = exec.Command("networksetup", "-setairportnetwork", iface, networkName, password)
	} else {
		cmd = exec.Command("networksetup", "-setairportnetwork", iface, networkName)
	}

	if err := cmd.Run(); err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf(
			"Failed to connect to %s. Make sure the network name is correct and you have the password", networkName))
	}

	ctx.Output.Success("Connected to %s", networkName)
	return nil
}

func (w *WiFiModule) forget(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		return module.NewExitError(module.ExitUsage, "Network name required: mac wifi forget <network>")
	}
	networkName := ctx.Args[0]

	iface, err := getInterface()
	if err != nil {
		return err
	}

	ctx.Output.Info("Forgetting network: %s", networkName)

	if err := exec.Command("networksetup", "-removepreferredwirelessnetwork", iface, networkName).Run(); err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf(
			"Failed to forget network %s. Make sure the network name is correct", networkName))
	}

	ctx.Output.Success("Forgot network: %s", networkName)
	ctx.Output.Info("You will need to re-enter the password to reconnect")
	return nil
}

func (w *WiFiModule) info(ctx *module.Context) error {
	ctx.Output.Header("\U0001F310 Detailed WiFi Information")

	iface, err := getInterface()
	if err != nil {
		return err
	}

	// Basic status first.
	_ = w.status(ctx)
	fmt.Println()

	network, _ := getCurrentNetwork()
	if network == "Not connected" {
		return nil
	}

	// Detailed airport -I information.
	ctx.Output.Header("Connection Details:")
	out, err := exec.Command(airportBin, "-I").Output()
	if err == nil {
		parsed := parseAirportInfo(string(out))
		if v, ok := parsed["SSID"]; ok {
			fmt.Printf("  Network:   %s\n", v)
		}
		if v, ok := parsed["BSSID"]; ok {
			fmt.Printf("  Router:    %s\n", v)
		}
		if v, ok := parsed["channel"]; ok {
			fmt.Printf("  Channel:   %s\n", v)
		}
		if v, ok := parsed["CC"]; ok {
			fmt.Printf("  Country:   %s\n", v)
		}
	}

	// Network configuration (gateway, DNS).
	fmt.Println()
	ctx.Output.Header("Network Configuration:")

	gw, err := exec.Command("route", "-n", "get", "default").Output()
	if err == nil {
		for _, line := range strings.Split(string(gw), "\n") {
			if strings.Contains(line, "gateway") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					fmt.Printf("  Gateway:   %s\n", fields[len(fields)-1])
				}
			}
		}
	}

	dnsOut, err := exec.Command("scutil", "--dns").Output()
	if err == nil {
		var dnsServers []string
		for _, line := range strings.Split(string(dnsOut), "\n") {
			if strings.Contains(line, "nameserver") {
				fields := strings.Fields(line)
				if len(fields) >= 3 {
					dnsServers = append(dnsServers, fields[len(fields)-1])
				}
			}
			if len(dnsServers) >= 3 {
				break
			}
		}
		if len(dnsServers) > 0 {
			fmt.Printf("  DNS:       %s\n", strings.Join(dnsServers, " "))
		}
	}

	// Unused import guard (iface is used above via getInterface).
	_ = iface
	return nil
}

func (w *WiFiModule) listSaved(ctx *module.Context) error {
	iface, err := getInterface()
	if err != nil {
		return err
	}

	ctx.Output.Header("\U0001F4BE Saved WiFi Networks")

	out, err := exec.Command("networksetup", "-listpreferredwirelessnetworks", iface).Output()
	if err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to list saved networks: %v", err))
	}

	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "Preferred networks") {
			continue
		}
		fmt.Printf("  %s\n", line)
	}
	return nil
}

// showCurrentConnection prints the current network and signal strength.
func showCurrentConnection(ctx *module.Context, iface string) {
	network, err := getCurrentNetwork()
	if err != nil || network == "Not connected" {
		return
	}
	fmt.Printf("  Connected to:   %s\n", network)

	info, err := exec.Command(airportBin, "-I").Output()
	if err == nil {
		parsed := parseAirportInfo(string(info))
		if rssi, ok := parsed["agrCtlRSSI"]; ok {
			fmt.Printf("  Signal:         %s dBm\n", rssi)
		}
	}
}

// Ensure platform types are referenced (avoids unused import).
var _ platform.WiFiInfo
