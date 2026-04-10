package dev

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cosmolabs-org/supermac/internal/dep"
	"github.com/cosmolabs-org/supermac/internal/module"
)

func init() {
	module.Register(&DevModule{})
}

type DevModule struct{}

func (d *DevModule) Name() string            { return "dev" }
func (d *DevModule) ShortDescription() string { return "Developer tools and utilities" }
func (d *DevModule) Emoji() string            { return "💻" }

func (d *DevModule) Commands() []module.Command {
	return []module.Command{
		{
			Name:        "kill-port",
			Description: "Kill process on specific port",
			Aliases:     []string{"kp"},
			Args: []module.Arg{
				{Name: "port", Required: true, Description: "Port number"},
			},
			Run: d.killPort,
		},
		{
			Name:        "ports",
			Description: "Show all processes using network ports",
			Aliases:     []string{"list-ports"},
			Run:         d.listPorts,
		},
		{
			Name:        "servers",
			Description: "List running development servers",
			Run:         d.servers,
		},
		{
			Name:        "localhost",
			Description: "Open localhost in browser",
			Args: []module.Arg{
				{Name: "port", Required: true, Description: "Port number"},
			},
			Flags: []module.Flag{
				{Name: "protocol", Shorthand: "p", DefaultValue: "http", Description: "Protocol to use (http/https)"},
			},
			Run: d.localhost,
		},
		{
			Name:        "serve",
			Description: "Start HTTP server in directory",
			Args: []module.Arg{
				{Name: "dir", Required: false, Description: "Directory to serve (default: current)"},
			},
			Flags: []module.Flag{
				{Name: "port", Shorthand: "p", DefaultValue: "8000", Description: "Port to listen on"},
			},
			Run: d.serve,
		},
		{
			Name:        "processes",
			Description: "Enhanced process viewer",
			Flags: []module.Flag{
				{Name: "sort", Shorthand: "s", DefaultValue: "cpu", Description: "Sort by: cpu, memory"},
				{Name: "count", Shorthand: "n", DefaultValue: "15", Description: "Number of processes to show"},
			},
			Run: d.processes,
		},
		{
			Name:        "cpu-hogs",
			Description: "Show CPU-intensive processes",
			Run:         d.cpuHogs,
		},
		{
			Name:        "memory-hogs",
			Description: "Show memory-intensive processes",
			Run:         d.memoryHogs,
		},
		{
			Name:        "uuid",
			Description: "Generate a UUID and copy to clipboard",
			Run:         d.uuid,
		},
		{
			Name:        "env",
			Description: "Show development environment info",
			Run:         d.env,
		},
		{
			Name:        "json-format",
			Description: "Pretty-print a JSON file in-place",
			Aliases:     []string{"jf"},
			Args: []module.Arg{
				{Name: "file", Required: true, Description: "JSON file to format"},
			},
			Run: d.jsonFormat,
		},
		{
			Name:        "base64-encode",
			Description: "Base64 encode text (copies to clipboard)",
			Aliases:     []string{"b64e"},
			Args: []module.Arg{
				{Name: "text", Required: true, Description: "Text to encode"},
			},
			Run: d.base64Encode,
		},
		{
			Name:        "base64-decode",
			Description: "Base64 decode text (copies to clipboard)",
			Aliases:     []string{"b64d"},
			Args: []module.Arg{
				{Name: "text", Required: true, Description: "Base64 string to decode"},
			},
			Run: d.base64Decode,
		},
		{
			Name:        "password",
			Description: "Generate a secure random password (copies to clipboard)",
			Aliases:     []string{"pw"},
			Args: []module.Arg{
				{Name: "length", Required: false, Description: "Password length (default: 20)"},
			},
			Run: d.password,
		},
		{
			Name:        "hash",
			Description: "Compute file hash (default: sha256)",
			Args: []module.Arg{
				{Name: "file", Required: true, Description: "File to hash"},
				{Name: "algorithm", Required: false, Description: "Algorithm: md5, sha1, sha256 (default), sha512"},
			},
			Run: d.hashFile,
		},
		{
			Name:        "timestamp",
			Description: "Convert between unix timestamps and human dates",
			Args: []module.Arg{
				{Name: "value", Required: false, Description: "Unix timestamp or date string (default: now)"},
			},
			Run: d.timestamp,
		},
	}
}

func (d *DevModule) Search(term string) []module.SearchResult {
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

func (d *DevModule) Dependencies() []dep.Dependency { return nil }

// ---------------------------------------------------------------------------
// Command implementations
// ---------------------------------------------------------------------------

func (d *DevModule) killPort(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		return module.NewExitError(module.ExitUsage, "Port number required: mac dev kill-port <port>")
	}

	port, err := strconv.Atoi(ctx.Args[0])
	if err != nil || port < 1 || port > 65535 {
		return module.NewExitError(module.ExitUsage, fmt.Sprintf("Invalid port number: %s", ctx.Args[0]))
	}

	ctx.Output.Info("Looking for processes on port %d...", port)

	out, err := exec.Command("lsof", "-ti:"+strconv.Itoa(port)).Output()
	if err != nil || len(out) == 0 {
		ctx.Output.Warning("No process found running on port %d", port)
		return nil
	}

	pids := strings.Fields(strings.TrimSpace(string(out)))
	for _, pidStr := range pids {
		// Get process name for reporting
		name := getProcessName(pidStr)

		ctx.Output.Info("Found process: %s (PID: %s)", name, pidStr)

		if err := exec.Command("kill", "-9", pidStr).Run(); err != nil {
			ctx.Output.Error("Failed to kill process %s", pidStr)
			continue
		}
		ctx.Output.Success("Killed process %s (PID: %s) on port %d", name, pidStr, port)
	}

	return nil
}

func (d *DevModule) listPorts(ctx *module.Context) error {
	ctx.Output.Header("Active Network Ports")

	// Listening ports
	out, err := exec.Command("lsof", "-i", "-P", "-n", "-sTCP:LISTEN").Output()
	if err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to list ports: %v", err))
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) <= 1 {
		ctx.Output.Info("No listening ports found")
		return nil
	}

	fmt.Println()
	fmt.Println("  Listening Ports:")
	var rows [][]string
	for _, line := range lines[1:] { // skip header
		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}
		rows = append(rows, []string{fields[0], fields[1], fields[8]})
	}
	ctx.Output.Table([]string{"Process", "PID", "Address"}, rows)

	// Common dev ports check
	fmt.Println()
	fmt.Println("  Common Development Ports:")
	devPorts := map[int]string{
		3000: "React/Next.js",
		3001: "React (alt)",
		4000: "Gatsby/Express",
		5000: "Flask/Express",
		5173: "Vite",
		8000: "Django/Python",
		8080: "Webpack/Tomcat",
		8888: "Jupyter",
		9000: "PHP/Node",
		9001: "SvelteKit",
	}

	// Sort ports for consistent output
	var sortedPorts []int
	for p := range devPorts {
		sortedPorts = append(sortedPorts, p)
	}
	sort.Ints(sortedPorts)

	found := false
	for _, p := range sortedPorts {
		out, err := exec.Command("lsof", "-ti:"+strconv.Itoa(p)).Output()
		if err != nil || len(out) == 0 {
			continue
		}
		pid := strings.TrimSpace(string(out))
		name := getProcessName(pid)
		fmt.Printf("  %-6d %-20s %s (PID: %s)\n", p, devPorts[p], name, pid)
		found = true
	}

	if !found {
		ctx.Output.Info("No processes found on common development ports")
	}

	return nil
}

func (d *DevModule) servers(ctx *module.Context) error {
	ctx.Output.Header("Running Development Servers")

	commonPorts := map[int]string{
		3000: "React/Next.js",
		3001: "React (alt)",
		4000: "Gatsby/Express",
		5000: "Flask/Express",
		5173: "Vite",
		8000: "Django/Python",
		8080: "Webpack/Tomcat",
		8888: "Jupyter",
		9000: "PHP/Node",
		9001: "SvelteKit",
	}

	var sortedPorts []int
	for p := range commonPorts {
		sortedPorts = append(sortedPorts, p)
	}
	sort.Ints(sortedPorts)

	found := false
	for _, port := range sortedPorts {
		out, err := exec.Command("lsof", "-ti:"+strconv.Itoa(port)).Output()
		if err != nil || len(out) == 0 {
			continue
		}

		pid := strings.Fields(strings.TrimSpace(string(out)))[0]
		name := getProcessName(pid)
		cmdLine := getProcessCmdLine(pid)

		fmt.Printf("  %-6d %-15s %-12s PID: %s\n", port, commonPorts[port], name, pid)
		fmt.Printf("         Command: %.50s\n", cmdLine)
		fmt.Println()
		found = true
	}

	if !found {
		ctx.Output.Info("No development servers found on common ports")
		ctx.Output.Info("Try 'mac dev ports' to see all active ports")
	}

	fmt.Println()
	fmt.Println("  Quick Actions:")
	fmt.Println("    mac dev kill-port <port>   Kill server on port")
	fmt.Println("    mac dev localhost <port>   Open in browser")
	fmt.Println("    mac dev serve <dir>        Start HTTP server")

	return nil
}

func (d *DevModule) localhost(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		return module.NewExitError(module.ExitUsage, "Port number required: mac dev localhost <port> [protocol]")
	}

	port, err := strconv.Atoi(ctx.Args[0])
	if err != nil || port < 1 || port > 65535 {
		return module.NewExitError(module.ExitUsage, fmt.Sprintf("Invalid port number: %s", ctx.Args[0]))
	}

	protocol := ctx.Flags["protocol"]
	if protocol == "" {
		protocol = "http"
	}

	url := fmt.Sprintf("%s://localhost:%d", protocol, port)

	// Check if anything is listening
	out, err := exec.Command("lsof", "-ti:"+strconv.Itoa(port)).Output()
	if err != nil || len(out) == 0 {
		ctx.Output.Warning("No service detected on port %d", port)
		confirmed, confirmErr := ctx.Prompt.Confirm(fmt.Sprintf("Open %s anyway?", url))
		if confirmErr != nil || !confirmed {
			return nil
		}
	}

	ctx.Output.Info("Opening %s in default browser...", url)
	if err := exec.Command("open", url).Run(); err != nil {
		ctx.Output.Error("Failed to open browser")
		ctx.Output.Info("URL: %s", url)
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("failed to open browser: %v", err))
	}

	ctx.Output.Success("Browser opened!")
	return nil
}

func (d *DevModule) serve(ctx *module.Context) error {
	dir := "."
	if len(ctx.Args) > 0 {
		dir = ctx.Args[0]
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Invalid directory: %s", dir))
	}

	info, err := os.Stat(absDir)
	if err != nil || !info.IsDir() {
		return module.NewExitError(module.ExitNotFound, fmt.Sprintf("Directory not found: %s", dir))
	}

	portStr := ctx.Flags["port"]
	if portStr == "" {
		portStr = "8000"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		return module.NewExitError(module.ExitUsage, fmt.Sprintf("Invalid port: %s", portStr))
	}

	// Check if port is already in use
	out, err := exec.Command("lsof", "-ti:"+strconv.Itoa(port)).Output()
	if err == nil && len(out) > 0 {
		ctx.Output.Error("Port %d is already in use", port)
		ctx.Output.Info("Use 'mac dev kill-port %d' to free it", port)
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("port %d already in use", port))
	}

	ctx.Output.Info("Starting HTTP server on port %d...", port)
	ctx.Output.Info("Serving directory: %s", absDir)
	ctx.Output.Info("URL: http://localhost:%d", port)
	ctx.Output.Info("Press Ctrl+C to stop")
	fmt.Println()

	// Use Go stdlib HTTP server — no Python dependency
	handler := http.FileServer(http.Dir(absDir))
	return http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
}

func (d *DevModule) processes(ctx *module.Context) error {
	sortBy := ctx.Flags["sort"]
	if sortBy == "" {
		sortBy = "cpu"
	}

	countStr := ctx.Flags["count"]
	if countStr == "" {
		countStr = "15"
	}
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 {
		count = 15
	}

	procs, err := getProcesses()
	if err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to list processes: %v", err))
	}

	switch sortBy {
	case "memory", "mem":
		sort.Slice(procs, func(i, j int) bool {
			return procs[i].MemPercent > procs[j].MemPercent
		})
		ctx.Output.Header(fmt.Sprintf("Top %d Processes by Memory Usage", min(count, len(procs))))
	default:
		sort.Slice(procs, func(i, j int) bool {
			return procs[i].CPUPercent > procs[j].CPUPercent
		})
		ctx.Output.Header(fmt.Sprintf("Top %d Processes by CPU Usage", min(count, len(procs))))
	}

	fmt.Println()
	var rows [][]string
	limit := min(count, len(procs))
	for i := 0; i < limit; i++ {
		p := procs[i]
		rows = append(rows, []string{
			truncate(p.Command, 20),
			fmt.Sprintf("%.1f%%", p.CPUPercent),
			fmt.Sprintf("%.1f%%", p.MemPercent),
			p.PID,
			p.User,
		})
	}
	ctx.Output.Table([]string{"Command", "CPU", "Memory", "PID", "User"}, rows)

	fmt.Println()
	ctx.Output.Info("Use 'mac dev cpu-hogs' or 'mac dev memory-hogs' for focused views")
	return nil
}

func (d *DevModule) cpuHogs(ctx *module.Context) error {
	ctx.Output.Header("CPU-Intensive Processes")

	procs, err := getProcesses()
	if err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to list processes: %v", err))
	}

	// Filter to > 1% CPU, sort descending, take top 10
	var hogs []processInfo
	for _, p := range procs {
		if p.CPUPercent > 1.0 {
			hogs = append(hogs, p)
		}
	}
	sort.Slice(hogs, func(i, j int) bool {
		return hogs[i].CPUPercent > hogs[j].CPUPercent
	})

	if len(hogs) > 10 {
		hogs = hogs[:10]
	}

	fmt.Println()
	if len(hogs) == 0 {
		ctx.Output.Info("No CPU-intensive processes found (> 1%% CPU)")
		return nil
	}

	var rows [][]string
	for _, p := range hogs {
		rows = append(rows, []string{
			truncate(p.Command, 20),
			fmt.Sprintf("%.1f%%", p.CPUPercent),
			fmt.Sprintf("%.1f%%", p.MemPercent),
			p.PID,
		})
	}
	ctx.Output.Table([]string{"Command", "CPU", "Memory", "PID"}, rows)

	fmt.Println()
	ctx.Output.Info("High CPU usage may indicate runaway processes")
	ctx.Output.Info("Use 'mac dev kill-port <port>' to stop development servers")
	return nil
}

func (d *DevModule) memoryHogs(ctx *module.Context) error {
	ctx.Output.Header("Memory-Intensive Processes")

	procs, err := getProcesses()
	if err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to list processes: %v", err))
	}

	// Filter to > 1% memory, sort descending, take top 10
	var hogs []processInfo
	for _, p := range procs {
		if p.MemPercent > 1.0 {
			hogs = append(hogs, p)
		}
	}
	sort.Slice(hogs, func(i, j int) bool {
		return hogs[i].MemPercent > hogs[j].MemPercent
	})

	if len(hogs) > 10 {
		hogs = hogs[:10]
	}

	fmt.Println()
	if len(hogs) == 0 {
		ctx.Output.Info("No memory-intensive processes found (> 1%% memory)")
		return nil
	}

	var rows [][]string
	for _, p := range hogs {
		rows = append(rows, []string{
			truncate(p.Command, 20),
			fmt.Sprintf("%.1f%%", p.CPUPercent),
			fmt.Sprintf("%.1f%%", p.MemPercent),
			p.PID,
		})
	}
	ctx.Output.Table([]string{"Command", "CPU", "Memory", "PID"}, rows)

	fmt.Println()
	ctx.Output.Info("High memory usage may slow down your system")
	ctx.Output.Info("Consider closing unused applications")
	return nil
}

func (d *DevModule) uuid(ctx *module.Context) error {
	// Generate UUID v4 using crypto/rand
	id := generateUUID()
	ctx.Output.Success("Generated UUID:")
	fmt.Printf("  %s\n", id)

	// Copy to clipboard via pbcopy
	if err := exec.Command("pbcopy").Run(); err == nil {
		// pbcopy is available — pipe the UUID to it
		cmd := exec.Command("pbcopy")
		cmd.Stdin = strings.NewReader(id)
		if cmd.Run() == nil {
			ctx.Output.Info("Copied to clipboard")
		}
	}

	return nil
}

func (d *DevModule) env(ctx *module.Context) error {
	ctx.Output.Header("Development Environment")

	fmt.Println()

	// Language versions
	langs := []struct {
		name string
		cmd  string
		args []string
	}{
		{"Go", "go", []string{"version"}},
		{"Node.js", "node", []string{"--version"}},
		{"Python3", "python3", []string{"--version"}},
		{"Rust", "rustc", []string{"--version"}},
		{"Java", "java", []string{"-version"}},
		{"Bun", "bun", []string{"--version"}},
	}

	fmt.Println("  Languages:")
	for _, lang := range langs {
		out, err := exec.Command(lang.cmd, lang.args...).CombinedOutput()
		if err != nil {
			fmt.Printf("    %-12s not installed\n", lang.name)
			continue
		}
		version := strings.TrimSpace(string(out))
		// Take only first line for messy output (like java -version going to stderr)
		if idx := strings.Index(version, "\n"); idx > 0 {
			version = version[:idx]
		}
		fmt.Printf("    %-12s %s\n", lang.name, version)
	}

	fmt.Println()

	// Package managers
	pkgs := []struct {
		name string
		cmd  string
		args []string
	}{
		{"Homebrew", "brew", []string{"--version"}},
		{"Git", "git", []string{"--version"}},
		{"Docker", "docker", []string{"--version"}},
	}

	fmt.Println("  Tools:")
	for _, tool := range pkgs {
		out, err := exec.Command(tool.cmd, tool.args...).CombinedOutput()
		if err != nil {
			fmt.Printf("    %-12s not installed\n", tool.name)
			continue
		}
		version := strings.TrimSpace(string(out))
		if idx := strings.Index(version, "\n"); idx > 0 {
			version = version[:idx]
		}
		fmt.Printf("    %-12s %s\n", tool.name, version)
	}

	fmt.Println()

	// Shell info
	fmt.Printf("  Shell:        %s\n", os.Getenv("SHELL"))
	fmt.Printf("  Terminal:     %s\n", os.Getenv("TERM_PROGRAM"))
	if term := os.Getenv("TERM"); term != "" {
		fmt.Printf("  TERM:         %s\n", term)
	}

	// Editor
	if editor := os.Getenv("EDITOR"); editor != "" {
		fmt.Printf("  Editor:       %s\n", editor)
	}

	return nil
}

// ---------------------------------------------------------------------------
// Helper types and functions
// ---------------------------------------------------------------------------

// processInfo holds parsed ps output for a single process.
type processInfo struct {
	User        string
	PID         string
	CPUPercent  float64
	MemPercent  float64
	Command     string
}

// getProcesses runs ps aux and returns parsed process info sorted by CPU.
func getProcesses() ([]processInfo, error) {
	out, err := exec.Command("ps", "aux").Output()
	if err != nil {
		return nil, fmt.Errorf("ps aux failed: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) <= 1 {
		return nil, nil
	}

	var procs []processInfo
	for _, line := range lines[1:] { // skip header
		fields := strings.Fields(line)
		if len(fields) < 11 {
			continue
		}

		cpu, _ := strconv.ParseFloat(fields[2], 64)
		mem, _ := strconv.ParseFloat(fields[3], 64)

		// Command is fields[10] onwards, take just the basename of the first part
		cmdParts := strings.Join(fields[10:], " ")
		cmdName := filepath.Base(strings.Fields(cmdParts)[0])

		procs = append(procs, processInfo{
			User:       fields[0],
			PID:        fields[1],
			CPUPercent: cpu,
			MemPercent: mem,
			Command:    cmdName,
		})
	}

	return procs, nil
}

// getProcessName returns the command name for a PID.
func getProcessName(pid string) string {
	out, err := exec.Command("ps", "-p", pid, "-o", "comm=").Output()
	if err != nil {
		return "unknown"
	}
	return filepath.Base(strings.TrimSpace(string(out)))
}

// getProcessCmdLine returns the full command line for a PID (truncated).
func getProcessCmdLine(pid string) string {
	out, err := exec.Command("ps", "-p", pid, "-o", "args=").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// generateUUID creates a UUID v4 using crypto/rand, falling back to uuidgen.
func generateUUID() string {
	// Try uuidgen first (always available on macOS)
	out, err := exec.Command("uuidgen").Output()
	if err == nil {
		return strings.ToLower(strings.TrimSpace(string(out)))
	}

	// Fallback: crypto/rand
	var uuid [16]byte
	if _, err := rand.Read(uuid[:]); err != nil {
		// Should never happen on macOS, but return something
		return "00000000-0000-4000-8000-000000000000"
	}

	// Set version 4 and variant bits per RFC 4122
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // variant 10

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16])
}

// truncate shortens a string to maxLen, appending "..." if needed.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// --- JSON Format ---

func (d *DevModule) jsonFormat(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		return module.NewExitError(module.ExitUsage, "File required: mac dev json-format <file.json>")
	}
	path := ctx.Args[0]
	if _, err := os.Stat(path); err != nil {
		return module.NewExitError(module.ExitNotFound, fmt.Sprintf("File not found: %s", path))
	}

	out, err := exec.Command("python3", "-m", "json.tool", path).Output()
	if err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Invalid JSON in %s", path))
	}

	if err := os.WriteFile(path, out, 0644); err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to write: %v", err))
	}

	ctx.Output.Success("Formatted %s", path)
	return nil
}

// --- Base64 Encode ---

func (d *DevModule) base64Encode(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		return module.NewExitError(module.ExitUsage, "Text required: mac dev base64-encode <text>")
	}
	input := strings.Join(ctx.Args, " ")
	encoded := base64.StdEncoding.EncodeToString([]byte(input))

	fmt.Println(encoded)
	// Copy to clipboard
	_ = exec.Command("pbcopy").Run()
	pbcopy := exec.Command("pbcopy")
	pbcopy.Stdin = strings.NewReader(encoded)
	_ = pbcopy.Run()

	ctx.Output.Success("Copied to clipboard")
	return nil
}

// --- Base64 Decode ---

func (d *DevModule) base64Decode(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		return module.NewExitError(module.ExitUsage, "Base64 string required: mac dev base64-decode <string>")
	}
	input := ctx.Args[0]

	decoded, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Invalid base64: %v", err))
	}

	fmt.Println(string(decoded))
	pbcopy := exec.Command("pbcopy")
	pbcopy.Stdin = strings.NewReader(string(decoded))
	_ = pbcopy.Run()

	ctx.Output.Success("Copied to clipboard")
	return nil
}

// --- Hash ---

func (d *DevModule) hashFile(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		return module.NewExitError(module.ExitUsage, "File required: mac dev hash <file> [algorithm]")
	}

	filePath := ctx.Args[0]
	if _, err := os.Stat(filePath); err != nil {
		return module.NewExitError(module.ExitNotFound, fmt.Sprintf("File not found: %s", filePath))
	}

	algorithm := "sha256"
	if len(ctx.Args) > 1 {
		algorithm = strings.ToLower(ctx.Args[1])
	}

	var h hash.Hash
	switch algorithm {
	case "md5":
		h = md5.New()
	case "sha1":
		h = sha1.New()
	case "sha256":
		h = sha256.New()
	case "sha512":
		h = sha512.New()
	default:
		return module.NewExitError(module.ExitUsage,
			fmt.Sprintf("Unknown algorithm: %s (valid: md5, sha1, sha256, sha512)", algorithm))
	}

	f, err := os.Open(filePath)
	if err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to open file: %v", err))
	}
	defer f.Close()

	if _, err := io.Copy(h, f); err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to read file: %v", err))
	}

	checksum := hex.EncodeToString(h.Sum(nil))

	fmt.Printf("%s  %s\n", checksum, filePath)
	ctx.Output.Info("Algorithm: %s", strings.ToUpper(algorithm))

	// Copy to clipboard
	pbcopy := exec.Command("pbcopy")
	pbcopy.Stdin = strings.NewReader(checksum)
	_ = pbcopy.Run()

	ctx.Output.Success("Hash copied to clipboard")
	return nil
}

// --- Timestamp ---

func (d *DevModule) timestamp(ctx *module.Context) error {
	now := time.Now()

	if len(ctx.Args) == 0 {
		// No args: show current unix timestamp
		fmt.Printf("  Unix timestamp: %d\n", now.Unix())
		fmt.Printf("  UTC:            %s\n", now.UTC().Format(time.RFC3339))
		fmt.Printf("  Local:          %s\n", now.Format("Mon Jan 2 15:04:05 MST 2006"))
		return nil
	}

	input := ctx.Args[0]

	// Try parsing as unix timestamp (integer)
	if ts, err := strconv.ParseInt(input, 10, 64); err == nil {
		t := time.Unix(ts, 0)
		fmt.Printf("  Unix timestamp: %d\n", ts)
		fmt.Printf("  UTC:            %s\n", t.UTC().Format(time.RFC3339))
		fmt.Printf("  Local:          %s\n", t.Format("Mon Jan 2 15:04:05 MST 2006"))
		fmt.Printf("  Relative:       %s\n", formatRelativeTime(now.Sub(t)))
		return nil
	}

	// Try parsing as a date string
	formats := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		time.RFC3339,
		"Jan 2, 2006",
		"January 2, 2006",
		"01/02/2006",
		"01/02/2006 15:04:05",
	}

	var parsed time.Time
	for _, layout := range formats {
		if t, err := time.Parse(layout, input); err == nil {
			parsed = t
			break
		}
	}

	if parsed.IsZero() {
		return module.NewExitError(module.ExitUsage,
			fmt.Sprintf("Cannot parse '%s' as timestamp or date. Try formats like: 1700000000, 2024-01-15, 'Jan 15, 2024'", input))
	}

	fmt.Printf("  Date:           %s\n", parsed.Format("Mon Jan 2 15:04:05 MST 2006"))
	fmt.Printf("  Unix timestamp: %d\n", parsed.Unix())
	fmt.Printf("  UTC:            %s\n", parsed.UTC().Format(time.RFC3339))

	// Copy unix timestamp to clipboard
	tsStr := strconv.FormatInt(parsed.Unix(), 10)
	pbcopy := exec.Command("pbcopy")
	pbcopy.Stdin = strings.NewReader(tsStr)
	_ = pbcopy.Run()
	ctx.Output.Success("Unix timestamp copied to clipboard")

	return nil
}

// formatRelativeTime returns a human-readable relative time string.
func formatRelativeTime(d time.Duration) string {
	if d < 0 {
		d = -d
	}
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(d.Hours()))
	case d < 365*24*time.Hour:
		return fmt.Sprintf("%d days ago", int(d.Hours()/24))
	default:
		return fmt.Sprintf("%.1f years ago", d.Hours()/(24*365))
	}
}

// --- Password Generator ---

func (d *DevModule) password(ctx *module.Context) error {
	length := 20
	if len(ctx.Args) > 0 {
		if l, err := strconv.Atoi(ctx.Args[0]); err == nil && l > 0 {
			length = l
		}
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:,.<>?"
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		return module.NewExitError(module.ExitGeneral, "Failed to generate password")
	}

	var password strings.Builder
	for i := 0; i < length; i++ {
		password.WriteByte(charset[int(buf[i])%len(charset)])
	}

	result := password.String()
	fmt.Println(result)

	pbcopy := exec.Command("pbcopy")
	pbcopy.Stdin = strings.NewReader(result)
	_ = pbcopy.Run()

	ctx.Output.Success("Password (%d chars) copied to clipboard", length)
	return nil
}
