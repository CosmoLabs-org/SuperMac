package system

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cosmolabs-org/supermac/internal/module"
)

func init() {
	module.Register(&SystemModule{})
}

// SystemModule implements module.Module for system information and maintenance.
type SystemModule struct{}

func (s *SystemModule) Name() string            { return "system" }
func (s *SystemModule) ShortDescription() string { return "System information and maintenance" }
func (s *SystemModule) Emoji() string            { return "🖥️" }

func (s *SystemModule) Commands() []module.Command {
	return []module.Command{
		{
			Name:        "info",
			Description: "Comprehensive system information",
			Run:         s.info,
		},
		{
			Name:        "cleanup",
			Description: "Deep system cleanup (cache, downloads, trash, logs)",
			Run:         s.cleanup,
		},
		{
			Name:        "battery",
			Description: "Battery status and health",
			Run:         s.battery,
		},
		{
			Name:        "memory",
			Description: "Memory usage statistics",
			Aliases:     []string{"mem"},
			Run:         s.memory,
		},
		{
			Name:        "cpu",
			Description: "CPU usage and information",
			Run:         s.cpu,
		},
		{
			Name:        "hardware",
			Description: "Hardware information and specs",
			Run:         s.hardware,
		},
		{
			Name:        "disk-usage",
			Description: "Disk usage analysis by directory",
			Args: []module.Arg{
				{Name: "path", Required: false, Description: "Directory to analyze (default: home)"},
			},
			Run: s.diskUsage,
		},
		{
			Name:        "processes",
			Description: "Top processes by resource usage (cpu or memory)",
			Args: []module.Arg{
				{Name: "sort", Required: false, Description: "Sort by: cpu (default) or memory"},
			},
			Run: s.processes,
		},
		{
			Name:        "uptime",
			Description: "System uptime with details",
			Run:         s.uptime,
		},
	}
}

func (s *SystemModule) Search(term string) []module.SearchResult {
	var results []module.SearchResult
	for _, cmd := range s.Commands() {
		if matchesSearch(cmd.Name, cmd.Description, cmd.Aliases, term) {
			results = append(results, module.SearchResult{
				Command:     cmd.Name,
				Description: cmd.Description,
				Module:      s.Name(),
			})
		}
	}
	// Additional keyword matches not covered by name/description
	extra := map[string]module.SearchResult{
		"power":     {Command: "battery", Description: "Battery status and health", Module: s.Name()},
		"ram":       {Command: "memory", Description: "Memory usage statistics", Module: s.Name()},
		"processor": {Command: "cpu", Description: "CPU usage and information", Module: s.Name()},
		"clean":     {Command: "cleanup", Description: "Deep system cleanup (cache, downloads, trash, logs)", Module: s.Name()},
		"status":    {Command: "info", Description: "Comprehensive system information", Module: s.Name()},
	}
	if r, ok := extra[term]; ok {
		results = append(results, r)
	}
	return results
}

// matchesSearch checks whether a command matches the search term by name,
// description, or any alias.
func matchesSearch(name, desc string, aliases []string, term string) bool {
	t := strings.ToLower(term)
	if strings.Contains(name, t) {
		return true
	}
	if strings.Contains(strings.ToLower(desc), t) {
		return true
	}
	for _, a := range aliases {
		if strings.Contains(strings.ToLower(a), t) {
			return true
		}
	}
	return false
}

// ============================================================================
// Command implementations
// ============================================================================

func (s *SystemModule) info(ctx *module.Context) error {
	ctx.Output.Header("System Information")

	// macOS version
	productName, _ := runCmd("sw_vers", "-productName")
	productVersion, _ := runCmd("sw_vers", "-productVersion")
	build, _ := runCmd("sw_vers", "-buildVersion")
	fmt.Printf("  OS:          %s %s\n", productName, productVersion)
	fmt.Printf("  Build:       %s\n", build)

	// Architecture
	arch, _ := runCmd("uname", "-m")
	fmt.Printf("  Architecture: %s\n", arch)

	// Hostname
	hostname, _ := runCmd("hostname")
	fmt.Printf("  Hostname:    %s\n", hostname)

	// Uptime
	uptime, _ := runCmd("uptime")
	// Extract just the uptime portion
	if parts := strings.Split(uptime, "up "); len(parts) > 1 {
		if loadParts := strings.Split(parts[1], ", load"); len(loadParts) > 0 {
			fmt.Printf("  Uptime:      %s\n", strings.TrimSpace(strings.TrimRight(loadParts[0], ",")))
		}
	}

	// Shell
	shell := "unknown"
	if sh := os.Getenv("SHELL"); sh != "" {
		shell = filepath.Base(sh)
	}
	fmt.Printf("  Shell:       %s\n", shell)

	// Disk usage summary
	ctx.Output.Info("Storage:")
	dfOut, err := runCmd("df", "-h", "/")
	if err == nil {
		lines := strings.Split(dfOut, "\n")
		if len(lines) >= 2 {
			fields := strings.Fields(lines[1])
			if len(fields) >= 6 {
				fmt.Printf("  Used:        %s of %s (%s)\n", fields[2], fields[1], fields[4])
				fmt.Printf("  Free:        %s\n", fields[3])
			}
		}
	}

	return nil
}

func (s *SystemModule) cleanup(ctx *module.Context) error {
	ctx.Output.Header("System Cleanup")
	ctx.Output.Warning("This will clean caches, logs, and temporary files")

	confirmed, err := ctx.Prompt.Confirm("Continue with system cleanup?")
	if err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Prompt failed: %v", err))
	}
	if !confirmed {
		ctx.Output.Info("Cleanup cancelled")
		return nil
	}

	// User caches
	if out, err := runCmd("test", "-d", "~/Library/Caches"); err == nil && out != "" {
		ctx.Output.Info("Cleaning user caches...")
		if _, err := ctx.Platform.RunCommand("find", "~/Library/Caches", "-type", "f", "-atime", "+7", "-delete"); err == nil {
			ctx.Platform.RunCommand("find", "~/Library/Caches", "-type", "d", "-empty", "-delete")
		}
		ctx.Output.Success("User caches cleaned")
	}

	// Downloads cleanup (files older than 30 days)
	ctx.Output.Info("Cleaning old downloads (30+ days)...")
	if count, err := ctx.Platform.RunCommand("find", "~/Downloads", "-type", "f", "-mtime", "+30"); err == nil {
		files := strings.Count(count, "\n") + 1
		if files > 0 && count != "" {
			ctx.Platform.RunCommand("find", "~/Downloads", "-type", "f", "-mtime", "+30", "-delete")
			ctx.Output.Success("Cleaned %d old download files", files)
		} else {
			ctx.Output.Info("No old downloads to clean")
		}
	}

	// Trash
	ctx.Output.Info("Emptying trash...")
	ctx.Platform.RunOSAScript(`tell application "Finder" to empty trash`)
	ctx.Output.Success("Trash emptied")

	// System logs (best effort)
	ctx.Output.Info("Cleaning system logs...")
	ctx.Platform.RunSudoCommand("find", "/var/log", "-name", "*.log", "-mtime", "+7", "-delete")
	ctx.Platform.RunSudoCommand("find", "/private/var/log", "-name", "*.log", "-mtime", "+7", "-delete")
	ctx.Output.Success("System logs cleaned")

	// Safari cache
	ctx.Output.Info("Cleaning Safari cache...")
	ctx.Platform.RunCommand("rm", "-rf", "~/Library/Caches/com.apple.Safari/*")
	ctx.Output.Success("Safari cache cleaned")

	// Temporary files
	ctx.Output.Info("Cleaning temporary files...")
	ctx.Platform.RunCommand("find", "/tmp", "-type", "f", "-atime", "+7", "-delete")
	ctx.Output.Success("Temporary files cleaned")

	// Font caches
	ctx.Output.Info("Clearing font caches...")
	ctx.Platform.RunSudoCommand("atsutil", "databases", "-remove")
	ctx.Output.Success("Font caches cleared")

	ctx.Output.Success("System cleanup completed!")
	ctx.Output.Info("Consider restarting applications to free additional memory")
	return nil
}

func (s *SystemModule) battery(ctx *module.Context) error {
	info, err := ctx.Platform.GetBatteryInfo()
	if err != nil {
		ctx.Output.Info("No battery information available (desktop Mac)")
		return nil
	}

	ctx.Output.Header("Battery Information")
	fmt.Printf("  Charge:          %d%%\n", info.Percent)

	status := "Discharging"
	if info.Charging {
		status = "Charging"
	}
	fmt.Printf("  Status:          %s\n", status)

	if info.TimeRemaining != "" {
		fmt.Printf("  Time remaining:  %s\n", info.TimeRemaining)
	}

	fmt.Printf("  Cycle count:     %d\n", info.CycleCount)
	fmt.Printf("  Health:          %s\n", info.Health)

	return nil
}

func (s *SystemModule) memory(ctx *module.Context) error {
	info, err := ctx.Platform.GetMemoryInfo()
	if err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to get memory info: %v", err))
	}

	ctx.Output.Header("Memory Usage")

	totalMB := info.Total / 1024 / 1024
	usedMB := info.Used / 1024 / 1024
	freeMB := info.Free / 1024 / 1024
	activeMB := info.Active / 1024 / 1024
	inactiveMB := info.Inactive / 1024 / 1024
	wiredMB := info.Wired / 1024 / 1024

	fmt.Printf("  Total:      %d MB\n", totalMB)
	fmt.Printf("  Used:       %d MB\n", usedMB)
	fmt.Printf("  Free:       %d MB\n", freeMB)
	fmt.Println()
	fmt.Printf("  Active:     %d MB\n", activeMB)
	fmt.Printf("  Inactive:   %d MB\n", inactiveMB)
	fmt.Printf("  Wired:      %d MB\n", wiredMB)

	if info.Compressed > 0 {
		compressedMB := info.Compressed / 1024 / 1024
		fmt.Printf("  Compressed: %d MB\n", compressedMB)
	}

	if info.SwapTotal > 0 {
		swapUsedMB := info.SwapUsed / 1024 / 1024
		swapTotalMB := info.SwapTotal / 1024 / 1024
		fmt.Println()
		fmt.Printf("  Swap used:  %d MB / %d MB\n", swapUsedMB, swapTotalMB)
	}

	// Memory pressure assessment
	freePct := float64(0)
	if totalMB > 0 {
		freePct = float64(freeMB) / float64(totalMB) * 100
	}
	fmt.Println()
	if freePct > 20 {
		ctx.Output.Success("Memory pressure: Low (%.0f%% free)", freePct)
	} else if freePct > 10 {
		ctx.Output.Warning("Memory pressure: Medium (%.0f%% free)", freePct)
	} else {
		ctx.Output.Warning("Memory pressure: High (%.0f%% free)", freePct)
		ctx.Output.Info("Consider closing some applications")
	}

	return nil
}

func (s *SystemModule) cpu(ctx *module.Context) error {
	info, err := ctx.Platform.GetCPUInfo()
	if err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to get CPU info: %v", err))
	}

	ctx.Output.Header("CPU Information")
	fmt.Printf("  Processor:       %s\n", info.Model)
	fmt.Printf("  Physical cores:  %d\n", info.Cores)
	fmt.Printf("  Logical cores:   %d\n", info.Threads)
	fmt.Println()
	fmt.Printf("  Current usage:   %.1f%%\n", info.Usage)

	// Load average from uptime
	uptime, _ := runCmd("uptime")
	if parts := strings.Split(uptime, "load averages: "); len(parts) > 1 {
		fmt.Printf("  Load average:    %s\n", strings.TrimSpace(parts[1]))
	}

	return nil
}

func (s *SystemModule) hardware(ctx *module.Context) error {
	info, err := ctx.Platform.GetHardwareInfo()
	if err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to get hardware info: %v", err))
	}

	ctx.Output.Header("Hardware Information")
	fmt.Printf("  Model:       %s\n", info.Model)
	fmt.Printf("  Chip:        %s\n", info.Chip)
	fmt.Printf("  Memory:      %s\n", info.Memory)
	fmt.Printf("  Serial:      %s\n", info.Serial)
	fmt.Printf("  OS Version:  %s\n", info.OSVersion)
	fmt.Printf("  Build:       %s\n", info.Build)
	fmt.Printf("  Arch:        %s\n", info.Arch)

	return nil
}

// ============================================================================
// Disk Usage
// ============================================================================

func (s *SystemModule) diskUsage(ctx *module.Context) error {
	target := os.Getenv("HOME")
	if len(ctx.Args) > 0 {
		target = ctx.Args[0]
	}

	info, err := os.Stat(target)
	if err != nil || !info.IsDir() {
		return module.NewExitError(module.ExitNotFound, fmt.Sprintf("Directory not found: %s", target))
	}

	ctx.Output.Header("Disk Usage: " + target)

	// Volume info
	dfOut, err := runCmd("df", "-h", target)
	if err == nil {
		lines := strings.Split(dfOut, "\n")
		if len(lines) >= 2 {
			fields := strings.Fields(lines[1])
			if len(fields) >= 6 {
				fmt.Printf("  Volume:     %s\n", fields[0])
				fmt.Printf("  Total:      %s\n", fields[1])
				fmt.Printf("  Used:       %s\n", fields[2])
				fmt.Printf("  Available:  %s\n", fields[3])
				fmt.Printf("  Usage:      %s\n", fields[4])
			}
		}
	}

	fmt.Println()
	ctx.Output.Info("Largest directories:")

	// du -sh on children, sorted by size
	out, err := runCmd("sh", "-c", fmt.Sprintf("du -sh %s/* 2>/dev/null | sort -hr | head -10", target))
	if err == nil {
		for _, line := range strings.Split(out, "\n") {
			if line == "" {
				continue
			}
			parts := strings.SplitN(line, "\t", 2)
			if len(parts) == 2 {
				fmt.Printf("  %-10s %s\n", parts[0], filepath.Base(parts[1]))
			}
		}
	}

	return nil
}

// ============================================================================
// Processes
// ============================================================================

func (s *SystemModule) processes(ctx *module.Context) error {
	sortBy := "cpu"
	if len(ctx.Args) > 0 {
		sortBy = strings.ToLower(ctx.Args[0])
	}

	switch sortBy {
	case "cpu":
		ctx.Output.Header("Top Processes (by CPU)")
	case "memory", "mem":
		ctx.Output.Header("Top Processes (by Memory)")
	default:
		return module.NewExitError(module.ExitUsage, "Unknown sort option. Use: cpu, memory")
	}

	fmt.Println()
	fmt.Printf("  %-20s %8s %8s %8s\n", "PROCESS", "CPU%", "MEM%", "PID")
	fmt.Println(strings.Repeat("-", 50))

	sortField := "-k3"
	if sortBy == "memory" || sortBy == "mem" {
		sortField = "-k4"
	}
	out, err := runCmd("sh", "-c", fmt.Sprintf("ps aux | sort -nr %s | head -15", sortField))
	if err != nil {
		return module.NewExitError(module.ExitGeneral, "Failed to list processes")
	}

	for _, line := range strings.Split(out, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 11 || fields[0] == "USER" {
			continue
		}
		cmdName := filepath.Base(fields[10])
		fmt.Printf("  %-20s %8s %8s %8s\n", cmdName, fields[2], fields[3], fields[1])
	}

	return nil
}

// ============================================================================
// Uptime
// ============================================================================

func (s *SystemModule) uptime(ctx *module.Context) error {
	ctx.Output.Header("System Uptime")

	out, err := runCmd("uptime")
	if err != nil {
		return module.NewExitError(module.ExitGeneral, "Failed to get uptime")
	}

	fmt.Println()
	if parts := strings.Split(out, "up "); len(parts) > 1 {
		uptimePart := parts[1]
		if loadParts := strings.Split(uptimePart, ", load"); len(loadParts) > 0 {
			fmt.Printf("  Up since:    %s\n", strings.TrimSpace(strings.TrimRight(loadParts[0], ",")))
		}
	}

	// Users
	for _, p := range strings.Split(out, ", ") {
		if strings.HasSuffix(p, "users") || strings.HasSuffix(p, "user") {
			fmt.Printf("  Users:       %s\n", strings.TrimSpace(p))
		}
	}

	// Load averages
	if parts := strings.Split(out, "load averages: "); len(parts) > 1 {
		loadFields := strings.Fields(strings.TrimSpace(parts[1]))
		if len(loadFields) >= 3 {
			fmt.Println()
			fmt.Printf("  Load avg:    1m: %s  5m: %s  15m: %s\n",
				loadFields[0], loadFields[1], loadFields[2])
		}
	}

	return nil
}

// ============================================================================
// Helpers
// ============================================================================

// runCmd executes a command and returns its trimmed stdout.
func runCmd(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
