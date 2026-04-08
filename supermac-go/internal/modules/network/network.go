package network

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/cosmolabs-org/supermac/internal/module"
)

func init() {
	module.Register(&NetworkModule{})
}

type NetworkModule struct{}

func (n *NetworkModule) Name() string            { return "network" }
func (n *NetworkModule) ShortDescription() string { return "Network information and troubleshooting" }
func (n *NetworkModule) Emoji() string            { return "📡" }

func (n *NetworkModule) Commands() []module.Command {
	return []module.Command{
		{
			Name:        "ip",
			Description: "Show local IP address and interface",
			Run:         n.localIP,
		},
		{
			Name:        "public-ip",
			Description: "Show public IP address with geolocation",
			Run:         n.publicIP,
		},
		{
			Name:        "dns-flush",
			Description: "Clear DNS cache (requires sudo)",
			Aliases:     []string{"flush-dns"},
			Run:         n.dnsFlush,
		},
		{
			Name:        "ping",
			Description: "Ping a host with enhanced output",
			Args: []module.Arg{
				{Name: "host", Required: true, Description: "Hostname or IP to ping"},
			},
			Flags: []module.Flag{
				{Name: "count", Shorthand: "c", DefaultValue: "5", Description: "Number of packets"},
			},
			Run: n.ping,
		},
		{
			Name:        "ports",
			Description: "Show listening ports and processes",
			Run:         n.ports,
		},
		{
			Name:        "reset",
			Description: "Reset network settings to defaults (requires sudo)",
			Run:         n.reset,
		},
		{
			Name:        "interfaces",
			Description: "List network interfaces and status",
			Run:         n.interfaces,
		},
		{
									Name:        "status",
			Aliases:     []string{"info"},
			Run:         n.status,
		},
		{
			Name:        "speed-test",
			Description: "Quick network speed test (download via curl)",
			Run:         n.speedTest,
		},
		{
			Name:        "renew-dhcp",
			Description: "Renew DHCP lease on primary interface",
			Run:         n.renewDHCP,
		},
		{
			Name:        "locations",
			Description: "List network locations and show current",
			Aliases:     []string{"loc"},
			Run:         n.locations,
		},
		{
			Name:        "connections",
			Description: "Show all active network connections",
			Run:         n.connections,
		},
	}
}

func (n *NetworkModule) Search(term string) []module.SearchResult {
	var results []module.SearchResult
	for _, cmd := range n.Commands() {
		if strings.Contains(cmd.Name, term) || strings.Contains(strings.ToLower(cmd.Description), term) {
			results = append(results, module.SearchResult{
				Command:     cmd.Name,
				Description: cmd.Description,
				Module:      n.Name(),
			})
		}
	}
	return results
}

// ---------------------------------------------------------------------------
// Command implementations
// ---------------------------------------------------------------------------

func (n *NetworkModule) localIP(ctx *module.Context) error {
	ip, iface, err := getLocalIP()
	if err != nil || ip == "" {
		ctx.Output.Warning("No active network connection found")
		ctx.Output.Info("Make sure you're connected to WiFi or Ethernet")
		return module.NewExitError(module.ExitNetwork, "no active network connection")
	}
	ctx.Output.Success("Local IP address: %s", ip)
	ctx.Output.Info("Interface: %s", iface)
	return nil
}

func (n *NetworkModule) publicIP(ctx *module.Context) error {
	ctx.Output.Info("Fetching public IP address...")

	ip, err := fetchPublicIP()
	if err != nil || ip == "" {
		ctx.Output.Error("Failed to retrieve public IP address")
		ctx.Output.Info("Check your internet connection")
		return module.NewExitError(module.ExitNetwork, fmt.Sprintf("failed to retrieve public IP: %v", err))
	}

	ctx.Output.Success("Public IP address: %s", ip)

	// Show geolocation if available
	if city, country, isp, locErr := fetchIPLocation(ip); locErr == nil {
		if city != "" && country != "" {
			fmt.Printf("  Location: %s, %s\n", city, country)
		}
		if isp != "" {
			fmt.Printf("  ISP:      %s\n", isp)
		}
	}
	return nil
}

func (n *NetworkModule) dnsFlush(ctx *module.Context) error {
	ctx.Output.Info("Flushing DNS cache...")
	if err := ctx.Platform.FlushDNS(); err != nil {
		return module.NewExitError(module.ExitPermission, fmt.Sprintf("failed to flush DNS: %v", err))
	}
	ctx.Output.Success("DNS cache cleared successfully")
	ctx.Output.Info("This can resolve DNS-related connectivity issues")
	return nil
}

func (n *NetworkModule) ping(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		return module.NewExitError(module.ExitUsage, "Host required: mac network ping <host>")
	}

	host := ctx.Args[0]
	count := ctx.Flags["count"]
	if count == "" {
		count = "5"
	}

	ctx.Output.Info("Pinging %s (%s packets)...", host, count)
	fmt.Println()

	cmd := exec.Command("ping", "-c", count, host)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println()
		ctx.Output.Error("Ping failed - host may be unreachable")
		return module.NewExitError(module.ExitNetwork, fmt.Sprintf("ping to %s failed", host))
	}

	fmt.Println()
	ctx.Output.Success("Ping completed successfully")
	return nil
}

func (n *NetworkModule) ports(ctx *module.Context) error {
	ctx.Output.Header("Listening Ports")

	cmd := exec.Command("lsof", "-i", "-P", "-n", "-sTCP:LISTEN")
	out, err := cmd.Output()
	if err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("failed to list ports: %v", err))
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 0 {
		ctx.Output.Info("No listening ports found")
		return nil
	}

	// Parse and display as table
	var rows [][]string
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}
		rows = append(rows, []string{fields[0], fields[1], fields[8]})
	}

	if len(rows) > 0 {
		ctx.Output.Table([]string{"Process", "PID", "Address"}, rows)
	}
	return nil
}

func (n *NetworkModule) reset(ctx *module.Context) error {
	ctx.Output.Warning("This will reset all network settings to defaults")

	confirmed, err := ctx.Prompt.Confirm("Are you sure you want to reset network settings?")
	if err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("prompt failed: %v", err))
	}
	if !confirmed {
		ctx.Output.Info("Network reset cancelled")
		return nil
	}

	ctx.Output.Info("Resetting network settings...")
	if err := ctx.Platform.ResetNetwork(); err != nil {
		return module.NewExitError(module.ExitPermission, fmt.Sprintf("failed to reset network: %v", err))
	}

	ctx.Output.Success("Network settings reset")
	ctx.Output.Warning("You may need to reconfigure WiFi networks and other settings")
	ctx.Output.Info("Consider restarting your Mac for a complete reset")
	return nil
}

func (n *NetworkModule) interfaces(ctx *module.Context) error {
	ctx.Output.Header("Network Interfaces")

	interfaces, err := net.Interfaces()
	if err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("failed to list interfaces: %v", err))
	}

	var rows [][]string
	for _, iface := range interfaces {
		status := "down"
		if iface.Flags&net.FlagUp != 0 {
			status = "up"
		}

		addrs, _ := iface.Addrs()
		addrStr := "-"
		if len(addrs) > 0 {
			addrStr = addrs[0].String()
		}

		rows = append(rows, []string{iface.Name, status, fmt.Sprintf("%d", iface.MTU), addrStr})
	}

	ctx.Output.Table([]string{"Interface", "Status", "MTU", "Address"}, rows)
	return nil
}

func (n *NetworkModule) status(ctx *module.Context) error {
	ctx.Output.Header("Network Status")
	fmt.Println()

	// Local IP and interface
	ip, iface, err := getLocalIP()
	if err != nil || ip == "" {
		fmt.Println("  Local IP:   Not connected")
	} else {
		fmt.Printf("  Local IP:   %s\n", ip)
		fmt.Printf("  Interface:  %s\n", iface)
	}

	// Gateway
	if gw := getGateway(); gw != "" {
		fmt.Printf("  Gateway:    %s\n", gw)
	}

	// DNS servers
	if dns := getDNSServers(); len(dns) > 0 {
		fmt.Printf("  DNS:        %s\n", strings.Join(dns, ", "))
	}

	// Public IP (best effort)
	ctx.Output.Info("Fetching public IP...")
	if pubIP, pubErr := fetchPublicIP(); pubErr == nil && pubIP != "" {
		fmt.Printf("  Public IP:  %s\n", pubIP)
	}

	fmt.Println()
	ctx.Output.Info("Use 'mac network interfaces' for detailed interface list")
	ctx.Output.Info("Use 'mac network ports' for listening port details")
	return nil
}

// ---------------------------------------------------------------------------
// Helper functions
// ---------------------------------------------------------------------------

// interfaces to probe for a local IP address.
var probeInterfaces = []string{"en0", "en1", "en2", "en3"}

// getLocalIP returns the first active IPv4 address and the interface it was found on.
func getLocalIP() (ipAddr, iface string, err error) {
	for _, name := range probeInterfaces {
		out, err := exec.Command("ipconfig", "getifaddr", name).Output()
		if err != nil {
			continue
		}
		candidate := strings.TrimSpace(string(out))
		if candidate != "" && net.ParseIP(candidate) != nil {
			return candidate, name, nil
		}
	}
	return "", "", fmt.Errorf("no active interface")
}

// fetchPublicIP tries multiple public IP services and returns the first valid result.
func fetchPublicIP() (string, error) {
	services := []string{
		"https://ifconfig.me",
		"https://ipinfo.io/ip",
		"https://api.ipify.org",
		"https://checkip.amazonaws.com",
	}

	ipPattern := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)

	for _, svc := range services {
		out, err := exec.Command("curl", "-s", "--connect-timeout", "10", svc).Output()
		if err != nil {
			continue
		}
		candidate := strings.TrimSpace(string(out))
		if ipPattern.MatchString(candidate) {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("all services failed")
}

// fetchIPLocation returns city, country, ISP for the given IP using ipinfo.io.
func fetchIPLocation(ip string) (city, country, isp string, err error) {
	out, err := exec.Command("curl", "-s", "--connect-timeout", "5",
		fmt.Sprintf("https://ipinfo.io/%s/json", ip)).Output()
	if err != nil {
		return "", "", "", err
	}

	body := string(out)
	city = extractJSONValue(body, "city")
	country = extractJSONValue(body, "country")
	isp = extractJSONValue(body, "org")
	return city, country, isp, nil
}

// extractJSONValue is a minimal JSON field extractor (avoids importing encoding/json
// for a simple two-field lookup from a known API shape).
func extractJSONValue(body, key string) string {
	pattern := regexp.MustCompile(`"` + key + `"\s*:\s*"([^"]*)"`)
	matches := pattern.FindStringSubmatch(body)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

// getGateway returns the default route gateway.
func getGateway() string {
	out, err := exec.Command("route", "-n", "get", "default").Output()
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "gateway:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "gateway:"))
		}
	}
	return ""
}

// getDNSServers returns up to 3 DNS nameserver addresses.
func getDNSServers() []string {
	out, err := exec.Command("scutil", "--dns").Output()
	if err != nil {
		return nil
	}

	seen := map[string]bool{}
	var servers []string
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "nameserver") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			addr := fields[2]
			if !seen[addr] {
				seen[addr] = true
				servers = append(servers, addr)
				if len(servers) >= 3 {
					break
				}
			}
		}
	}
	return servers
}

func (n *NetworkModule) speedTest(ctx *module.Context) error {
	ctx.Output.Header("Network Speed Test")
	fmt.Println()
	ctx.Output.Info("Downloading test file from Cloudflare...")

	start := time.Now()
	out, err := exec.Command("curl", "-o", "/dev/null", "-s", "-w",
		"%{speed_download} %{size_download} %{time_total}",
		"https://speed.cloudflare.com/__down?bytes=10000000").Output()

	if err != nil {
		ctx.Output.Info("Cloudflare test failed, trying basic connectivity...")
		err = exec.Command("curl", "-o", "/dev/null", "-s", "https://google.com").Run()
		elapsed := time.Since(start)
		if err != nil {
			return module.NewExitError(module.ExitGeneral, "Network appears to be down")
		}
		ctx.Output.Success("Network is up (response in %s)", elapsed.Round(time.Millisecond))
		return nil
	}

	parts := strings.Fields(string(out))
	if len(parts) >= 1 {
		var speedBytes float64
		fmt.Sscanf(parts[0], "%f", &speedBytes)
		speedMbps := speedBytes * 8 / 1000000
		elapsed := time.Since(start)

		fmt.Printf("  Download:    %.1f Mbps\n", speedMbps)
		fmt.Printf("  Time:        %s\n", elapsed.Round(time.Millisecond))
		fmt.Println()
		switch {
		case speedMbps > 100:
			ctx.Output.Success("Fast connection (%.0f Mbps)", speedMbps)
		case speedMbps > 30:
			ctx.Output.Success("Good connection (%.0f Mbps)", speedMbps)
		case speedMbps > 10:
			ctx.Output.Info("Moderate speed (%.0f Mbps)", speedMbps)
		default:
			ctx.Output.Warning("Slow connection (%.0f Mbps)", speedMbps)
		}
	}
	return nil
}

func (n *NetworkModule) renewDHCP(ctx *module.Context) error {
	iface := "en0"
	if len(ctx.Args) > 0 {
		iface = ctx.Args[0]
	}

	ctx.Output.Info("Renewing DHCP lease on %s...", iface)
	out, err := exec.Command("sudo", "-n", "ipconfig", "set", iface, "DHCP").CombinedOutput()
	if err != nil {
		if strings.Contains(string(out), "sudo") {
			return module.NewExitError(module.ExitGeneral,
				"Admin privileges required. Run: sudo mac network renew-dhcp")
		}
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed: %s", strings.TrimSpace(string(out))))
	}

	ip, _, err := getLocalIP()
	if err == nil {
		ctx.Output.Success("DHCP renewed. New IP: %s", ip)
	} else {
		ctx.Output.Success("DHCP lease renewed on %s", iface)
	}
	return nil
}

func (n *NetworkModule) locations(ctx *module.Context) error {
	ctx.Output.Header("Network Locations")

	out, err := exec.Command("networksetup", "-listlocations").Output()
	if err != nil {
		return module.NewExitError(module.ExitGeneral, "Failed to list network locations")
	}

	current, _ := exec.Command("networksetup", "-getcurrentlocation").Output()
	currentLoc := strings.TrimSpace(string(current))

	fmt.Println()
	for _, loc := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		loc = strings.TrimSpace(loc)
		if loc == "" {
			continue
		}
		marker := "  "
		if loc == currentLoc {
			marker = "* "
		}
		fmt.Printf("  %s%s\n", marker, loc)
	}

	if currentLoc != "" {
		fmt.Println()
		ctx.Output.Info("Current: %s", currentLoc)
	}

	fmt.Println()
	ctx.Output.Info("Switch with: sudo networksetup -switchtolocation <name>")
	return nil
}

func (n *NetworkModule) connections(ctx *module.Context) error {
	ctx.Output.Header("Active Network Connections")

	out, err := exec.Command("lsof", "-i", "-n", "-P").Output()
	if err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to list connections: %v", err))
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) <= 1 {
		ctx.Output.Info("No active network connections found")
		return nil
	}

	fmt.Println()
	fmt.Printf("  %-20s %-8s %-10s %-6s %-24s %-24s %s\n",
		"PROCESS", "PID", "USER", "PROTO", "LOCAL", "FOREIGN", "STATE")
	fmt.Println(strings.Repeat(" ", 2) + strings.Repeat("-", 110))

	for _, line := range lines[1:] { // skip header
		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}

		process := fields[0]
		pid := fields[1]
		user := fields[2]
		// fields[3] is "fd", fields[4] may be the protocol descriptor
		// The node field (local addr) and name field (foreign addr) positions vary
		// lsof -i -n -P output: COMMAND PID USER FD TYPE DEVICE SIZE/OFF NODE NAME
		// We need to find the connection info after NODE

		// Find the protocol and addresses from the end of the line
		// Re-join everything from field 8 onwards (NAME column and state)
		rest := strings.Join(fields[8:], " ")

		// Parse: (TCPUDP) (local->foreign) (state)
		// or: (TCPUDP) (local->foreign)
		var proto, localAddr, foreignAddr, state string

		// Extract protocol from the TYPE field area - look at fields[4]
		nodeType := ""
		if len(fields) > 4 {
			nodeType = fields[4]
		}
		if strings.HasPrefix(nodeType, "IPv") {
			// Look for TCP/UDP in the name portion
			// The name field contains: protocol local->foreign (state)
			parts := strings.SplitN(rest, " ", 2)
			if len(parts) >= 1 {
				connInfo := parts[0]
				stateParts := ""
				if len(parts) >= 2 {
					stateParts = parts[1]
				}

				// Parse local->foreign
				connParts := strings.SplitN(connInfo, "->", 2)
				if len(connParts) == 2 {
					localAddr = connParts[0]
					foreignAddr = connParts[1]
				} else {
					localAddr = connInfo
				}

				// Determine protocol from the * in FD or from connection type
				if strings.Contains(fields[3], "u") || strings.Contains(rest, "UDP") {
					proto = "UDP"
				} else {
					proto = "TCP"
				}

				// Extract state from parentheses
				if stateParts != "" {
					state = strings.Trim(stateParts, "()")
				}
				if state == "" && proto == "UDP" {
					state = "-"
				}
			}
		}

		if localAddr == "" {
			continue
		}

		// Truncate long process names
		if len(process) > 20 {
			process = process[:17] + "..."
		}

		fmt.Printf("  %-20s %-8s %-10s %-6s %-24s %-24s %s\n",
			process, pid, user, proto, localAddr, foreignAddr, state)
	}

	return nil
}
