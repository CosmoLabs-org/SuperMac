package platform

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// DarwinPlatform implements Interface using real macOS system commands.
type DarwinPlatform struct{}

// ---------------------------------------------------------------------------
// osascript
// ---------------------------------------------------------------------------

func (d *DarwinPlatform) RunOSAScript(script string) (string, error) {
	out, err := exec.Command("osascript", "-e", script).Output() //nolint:gosec // script is hardcoded in modules
	if err != nil {
		return "", fmt.Errorf("osascript: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// ---------------------------------------------------------------------------
// defaults read/write/delete
// ---------------------------------------------------------------------------

func (d *DarwinPlatform) ReadDefault(domain, key string) (string, error) {
	out, err := exec.Command("defaults", "read", domain, key).Output()
	if err != nil {
		return "", fmt.Errorf("defaults read %s %s: %w", domain, key, err)
	}
	return strings.TrimSpace(string(out)), nil
}

func (d *DarwinPlatform) WriteDefault(domain, key, value string) error {
	args := []string{"write", domain, key}
	args = append(args, strings.Fields(value)...)
	return exec.Command("defaults", args...).Run()
}

func (d *DarwinPlatform) DeleteDefault(domain, key string) error {
	return exec.Command("defaults", "delete", domain, key).Run()
}

// ---------------------------------------------------------------------------
// Network
// ---------------------------------------------------------------------------

func (d *DarwinPlatform) SetWiFi(on bool) error {
	action := "on"
	if !on {
		action = "off"
	}
	out, err := exec.Command("networksetup", "-setairportpower", "en0", action).CombinedOutput()
	if err != nil {
		return fmt.Errorf("set wifi %s: %s: %w", action, strings.TrimSpace(string(out)), err)
	}
	return nil
}

func (d *DarwinPlatform) GetWiFiStatus() (*WiFiInfo, error) {
	out, err := exec.Command(
		"/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport",
		"-I",
	).Output()
	if err != nil {
		return nil, fmt.Errorf("airport: %w", err)
	}
	info := &WiFiInfo{}
	for _, line := range strings.Split(string(out), "\n") {
		parts := strings.SplitN(strings.TrimSpace(line), ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		switch key {
		case "SSID":
			info.SSID = val
			info.Connected = val != ""
		case "BSSID":
			info.BSSID = val
		case "rssi":
			if v, err := strconv.Atoi(val); err == nil {
				info.Signal = v
			}
		case "channel":
			info.Channel = val
		}
	}
	return info, nil
}

func (d *DarwinPlatform) ScanWiFiNetworks() ([]Network, error) {
	out, err := exec.Command(
		"/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport",
		"-s",
	).Output()
	if err != nil {
		return nil, fmt.Errorf("airport scan: %w", err)
	}
	var networks []Network
	lines := strings.Split(string(out), "\n")
	for _, line := range lines[1:] { // skip header
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		networks = append(networks, Network{
			SSID:     fields[1],
			Signal:   0, // parsing RSSI from airport scan output is complex
			Security: fields[len(fields)-1],
		})
	}
	return networks, nil
}

func (d *DarwinPlatform) FlushDNS() error {
	return exec.Command("dscacheutil", "-flushcache").Run()
}

func (d *DarwinPlatform) ResetNetwork() error {
	_ = d.FlushDNS()
	_ = exec.Command("sudo", "-n", "ifconfig", "en0", "down").Run()
	_ = exec.Command("sudo", "-n", "ifconfig", "en0", "up").Run()
	return nil
}

// ---------------------------------------------------------------------------
// System
// ---------------------------------------------------------------------------

func (d *DarwinPlatform) GetMemoryInfo() (*MemoryInfo, error) {
	out, err := exec.Command("sysctl", "-n", "hw.memsize").Output()
	if err != nil {
		return nil, err
	}
	total, _ := strconv.ParseUint(strings.TrimSpace(string(out)), 10, 64)
	return &MemoryInfo{Total: total}, nil
}

func (d *DarwinPlatform) GetCPUInfo() (*CPUInfo, error) {
	model, _ := exec.Command("sysctl", "-n", "machdep.cpu.brand_string").Output()
	cores, _ := exec.Command("sysctl", "-n", "hw.ncpu").Output()
	info := &CPUInfo{
		Model: strings.TrimSpace(string(model)),
	}
	info.Cores, _ = strconv.Atoi(strings.TrimSpace(string(cores)))
	return info, nil
}

func (d *DarwinPlatform) GetBatteryInfo() (*BatteryInfo, error) {
	out, err := exec.Command("pmset", "-g", "batt").Output()
	if err != nil {
		return nil, err
	}
	info := &BatteryInfo{}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "%") {
			parts := strings.Split(line, "\t")
			if len(parts) >= 2 {
				pctStr := strings.TrimSuffix(strings.TrimSpace(parts[1]), "%;")
				fields := strings.Fields(pctStr)
				if len(fields) > 0 {
					info.Percent, _ = strconv.Atoi(fields[0])
				}
				info.Charging = strings.Contains(line, "AC Power") || strings.Contains(line, "charging")
			}
		}
	}
	return info, nil
}

func (d *DarwinPlatform) GetHardwareInfo() (*HardwareInfo, error) {
	model, _ := exec.Command("sysctl", "-n", "hw.model").Output()
	osVer, _ := exec.Command("sw_vers", "-productVersion").Output()
	build, _ := exec.Command("sw_vers", "-buildVersion").Output()
	arch, _ := exec.Command("uname", "-m").Output()
	chip, _ := exec.Command("sysctl", "-n", "machdep.cpu.brand_string").Output()
	mem, _ := exec.Command("sysctl", "-n", "hw.memsize").Output()

	memGB := ""
	if memBytes, err := strconv.ParseUint(strings.TrimSpace(string(mem)), 10, 64); err == nil {
		memGB = fmt.Sprintf("%d GB", memBytes/(1024*1024*1024))
	}

	return &HardwareInfo{
		Model:     strings.TrimSpace(string(model)),
		Chip:      strings.TrimSpace(string(chip)),
		Memory:    memGB,
		OSVersion: strings.TrimSpace(string(osVer)),
		Build:     strings.TrimSpace(string(build)),
		Arch:      strings.TrimSpace(string(arch)),
	}, nil
}

func (d *DarwinPlatform) GetPageSize() (int, error) {
	out, err := exec.Command("pagesize").Output()
	if err != nil {
		return 16384, nil // default for Apple Silicon
	}
	ps, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return 16384, nil
	}
	return ps, nil
}

// ---------------------------------------------------------------------------
// Display
// ---------------------------------------------------------------------------

func (d *DarwinPlatform) SetBrightness(level float64) error {
	// Requires external tool or CoreBrightness private framework
	return fmt.Errorf("SetBrightness requires external tool")
}

func (d *DarwinPlatform) GetDarkMode() (bool, error) {
	out, err := d.ReadDefault("NSGlobalDomain", "AppleInterfaceStyle")
	if err != nil {
		// Key doesn't exist = light mode
		return false, nil
	}
	return strings.TrimSpace(out) == "Dark", nil
}

func (d *DarwinPlatform) SetDarkMode(on bool) error {
	if on {
		return d.WriteDefault("NSGlobalDomain", "AppleInterfaceStyle", "-string Dark")
	}
	return d.DeleteDefault("NSGlobalDomain", "AppleInterfaceStyle")
}

// ---------------------------------------------------------------------------
// Audio
// ---------------------------------------------------------------------------

func (d *DarwinPlatform) GetVolume() (int, error) {
	out, err := d.RunOSAScript("output volume of (get volume settings)")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(out))
}

func (d *DarwinPlatform) SetVolume(level int) error {
	_, err := d.RunOSAScript(fmt.Sprintf("set volume output volume %d", level))
	return err
}

func (d *DarwinPlatform) GetAudioDevices() ([]AudioDevice, error) {
	// Try system_profiler first
	out, err := exec.Command("system_profiler", "SPAudioDataType", "-json").Output()
	if err != nil {
		return nil, err
	}
	// Basic parsing — the full JSON structure is complex.
	// For now, return an empty list and rely on SwitchAudioSource when available.
	_ = out
	return []AudioDevice{}, nil
}

// ---------------------------------------------------------------------------
// Process management
// ---------------------------------------------------------------------------

func (d *DarwinPlatform) ListProcesses(filter string) ([]Process, error) {
	args := []string{"-c", "ps aux"}
	if filter != "" {
		args = []string{"-c", fmt.Sprintf("ps aux | %s", filter)}
	}
	out, err := exec.Command("sh", args...).Output()
	if err != nil {
		return nil, err
	}
	var processes []Process
	lines := strings.Split(string(out), "\n")
	for _, line := range lines[1:] { // skip header
		fields := strings.Fields(line)
		if len(fields) < 11 {
			continue
		}
		cpu, _ := strconv.ParseFloat(fields[2], 64)
		mem, _ := strconv.ParseFloat(fields[3], 64)
		pid, _ := strconv.Atoi(fields[1])
		processes = append(processes, Process{
			PID:     pid,
			User:    fields[0],
			CPU:     cpu,
			Memory:  mem,
			Command: strings.Join(fields[10:], " "),
		})
	}
	return processes, nil
}

func (d *DarwinPlatform) KillPort(port int) error {
	return exec.Command("lsof", "-ti", fmt.Sprintf(":%d", port)).Run()
}

func (d *DarwinPlatform) GetPortUser(port int) (string, error) {
	out, err := exec.Command("lsof", "-ti", fmt.Sprintf(":%d", port)).Output()
	if err != nil || len(out) == 0 {
		return "", fmt.Errorf("no process on port %d", port)
	}
	return strings.TrimSpace(string(out)), nil
}

// ---------------------------------------------------------------------------
// General command execution
// ---------------------------------------------------------------------------

func (d *DarwinPlatform) RunCommand(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		return string(out), err
	}
	return strings.TrimSpace(string(out)), nil
}

func (d *DarwinPlatform) RunSudoCommand(name string, args ...string) (string, error) {
	allArgs := append([]string{"-n", name}, args...)
	out, err := exec.Command("sudo", allArgs...).CombinedOutput()
	if err != nil {
		return string(out), err
	}
	return strings.TrimSpace(string(out)), nil
}
