package apps

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cosmolabs-org/supermac/internal/dep"
	"github.com/cosmolabs-org/supermac/internal/module"
)

func init() {
	module.Register(&AppsModule{})
}

// AppsModule handles application listing, info, cache management, and lifecycle.
type AppsModule struct{}

func (a *AppsModule) Name() string            { return "apps" }
func (a *AppsModule) ShortDescription() string { return "Application management" }
func (a *AppsModule) Emoji() string            { return "📱" }

func (a *AppsModule) Commands() []module.Command {
	return []module.Command{
		{
			Name:        "list",
			Description: "List installed applications",
			Args: []module.Arg{
				{Name: "filter", Required: false, Description: "Optional name filter (case-insensitive grep)"},
			},
			Run: a.list,
		},
		{
			Name:        "info",
			Description: "Show detailed info about an application",
			Args: []module.Arg{
				{Name: "appname", Required: true, Description: "Application name to look up"},
			},
			Run: a.info,
		},
		{
			Name:        "cache-clear",
			Description: "Clear an application's cache, support, and preference files",
			Args: []module.Arg{
				{Name: "appname", Required: true, Description: "Application name or bundle ID"},
			},
			Run: a.cacheClear,
		},
		{
			Name:        "recent",
			Description: "Show recently used applications",
			Run:         a.recent,
		},
		{
			Name:        "kill",
			Description: "Kill/force-quit an application by name",
			Args: []module.Arg{
				{Name: "appname", Required: true, Description: "Application process name to kill"},
			},
			Run: a.kill,
		},
		{
			Name:        "open",
			Description: "Open an application by name",
			Args: []module.Arg{
				{Name: "appname", Required: true, Description: "Application name to open"},
			},
			Run: a.open,
		},
	}
}

func (a *AppsModule) Search(term string) []module.SearchResult {
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

func (a *AppsModule) Dependencies() []dep.Dependency { return nil }

// ---------------------------------------------------------------------------
// list — List installed applications
// ---------------------------------------------------------------------------

// spApp represents a single app entry from system_profiler JSON output.
type spApp struct {
	Name           string `json:"_name"`
	Version        string `json:"version"`
	ObtainedFrom   string `json:"obtained_from"`
	Path           string `json:"path"`
	LastModified   string `json:"lastModified"`
	Kind           string `json:"kind"`
}

// spAppsOutput represents the top-level system_profiler JSON structure.
type spAppsOutput struct {
	SPApplicationsDataType []spApp `json:"SPApplicationsDataType"`
}

func (a *AppsModule) list(ctx *module.Context) error {
	filter := ""
	if len(ctx.Args) > 0 {
		filter = strings.ToLower(ctx.Args[0])
	}

	ctx.Output.Info("Scanning installed applications...")

	apps, err := getAppsList()
	if err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to list applications: %v", err))
	}

	// Apply filter if provided.
	if filter != "" {
		var filtered []spApp
		for _, app := range apps {
			if strings.Contains(strings.ToLower(app.Name), filter) {
				filtered = append(filtered, app)
			}
		}
		apps = filtered
	}

	if len(apps) == 0 {
		if filter != "" {
			ctx.Output.Warning("No applications found matching '%s'", filter)
		} else {
			ctx.Output.Warning("No applications found")
		}
		return nil
	}

	ctx.Output.Header(fmt.Sprintf("Installed Applications (%d)", len(apps)))
	fmt.Println()

	headers := []string{"Name", "Version", "Path"}
	rows := make([][]string, 0, len(apps))
	for _, app := range apps {
		name := app.Name
		if name == "" {
			name = "(unknown)"
		}
		version := app.Version
		if version == "" {
			version = "-"
		}
		path := app.Path
		if path == "" {
			path = "-"
		}
		rows = append(rows, []string{name, version, path})
	}
	ctx.Output.Table(headers, rows)

	ctx.Output.Info("Showing %d applications", len(apps))
	return nil
}

// getAppsList tries system_profiler first, falls back to /Applications listing.
func getAppsList() ([]spApp, error) {
	// Try system_profiler for rich data.
	out, err := exec.Command("system_profiler", "SPApplicationsDataType", "-json").Output()
	if err == nil {
		var result spAppsOutput
		if jsonErr := json.Unmarshal(out, &result); jsonErr == nil && len(result.SPApplicationsDataType) > 0 {
			return result.SPApplicationsDataType, nil
		}
	}

	// Fallback: list /Applications directory.
	return listAppsFallback()
}

// listAppsFallback scans /Applications for .app bundles.
func listAppsFallback() ([]spApp, error) {
	entries, err := os.ReadDir("/Applications")
	if err != nil {
		return nil, fmt.Errorf("cannot read /Applications: %w", err)
	}

	apps := make([]spApp, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasSuffix(name, ".app") {
			continue
		}
		appName := strings.TrimSuffix(name, ".app")
		fullPath := filepath.Join("/Applications", name)
		apps = append(apps, spApp{
			Name: appName,
			Path: fullPath,
		})
	}
	return apps, nil
}

// ---------------------------------------------------------------------------
// info — Show detailed info about an application
// ---------------------------------------------------------------------------

func (a *AppsModule) info(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		return module.NewExitError(module.ExitUsage, "Application name required: mac apps info <appname>")
	}

	appName := ctx.Args[0]

	// Find the application using mdfind.
	ctx.Output.Info("Searching for '%s'...", appName)
	out, err := exec.Command("mdfind", "kMDItemKind == 'Application'", "-name", appName).Output()
	if err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Search failed: %v", err))
	}

	paths := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(paths) == 0 || paths[0] == "" {
		return module.NewExitError(module.ExitNotFound, fmt.Sprintf("Application '%s' not found", appName))
	}

	// Use the first match.
	appPath := paths[0]
	ctx.Output.Header(fmt.Sprintf("Application Info: %s", appName))
	fmt.Println()

	// Name.
	fmt.Printf("  Name:            %s\n", filepath.Base(strings.TrimSuffix(appPath, ".app")))

	// Version from Info.plist.
	version := readPlistValue(appPath, "CFBundleShortVersionString")
	if version == "" {
		version = readPlistValue(appPath, "CFBundleVersion")
	}
	if version != "" {
		fmt.Printf("  Version:         %s\n", version)
	} else {
		fmt.Printf("  Version:         (unknown)\n")
	}

	// Path.
	fmt.Printf("  Path:            %s\n", appPath)

	// Size.
	sizeOut, err := exec.Command("du", "-sh", appPath).Output()
	if err == nil {
		parts := strings.Fields(strings.TrimSpace(string(sizeOut)))
		if len(parts) > 0 {
			fmt.Printf("  Size:            %s\n", parts[0])
		}
	}

	// Last modified.
	stat, err := os.Stat(appPath)
	if err == nil {
		fmt.Printf("  Last Modified:   %s\n", stat.ModTime().Format("2006-01-02 15:04:05"))
	}

	// Architecture — inspect the main binary.
	binaryPath := findMainBinary(appPath)
	if binaryPath != "" {
		fileOut, err := exec.Command("file", "-b", binaryPath).Output()
		if err == nil {
			arch := strings.TrimSpace(string(fileOut))
			fmt.Printf("  Architecture:    %s\n", arch)
		}
	}

	// Show other matches if multiple were found.
	if len(paths) > 1 {
		fmt.Println()
		ctx.Output.Info("Other matches:")
		for i, p := range paths {
			if i == 0 {
				continue
			}
			fmt.Printf("    %s\n", p)
		}
	}

	return nil
}

// readPlistValue reads a key from the app's Info.plist using defaults/PlistBuddy.
func readPlistValue(appPath, key string) string {
	plistPath := filepath.Join(appPath, "Contents", "Info.plist")
	if _, err := os.Stat(plistPath); err != nil {
		return ""
	}
	out, err := exec.Command("/usr/libexec/PlistBuddy", "-c", "Print :"+key, plistPath).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// findMainBinary locates the executable inside a .app bundle.
func findMainBinary(appPath string) string {
	// Try to read CFBundleExecutable from plist.
	exeName := readPlistValue(appPath, "CFBundleExecutable")
	if exeName != "" {
		candidate := filepath.Join(appPath, "Contents", "MacOS", exeName)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	// Fallback: list MacOS/ dir and pick the first file.
	macosDir := filepath.Join(appPath, "Contents", "MacOS")
	entries, err := os.ReadDir(macosDir)
	if err != nil {
		return ""
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			return filepath.Join(macosDir, entry.Name())
		}
	}
	return ""
}

// ---------------------------------------------------------------------------
// cache-clear — Clear an app's cache files
// ---------------------------------------------------------------------------

func (a *AppsModule) cacheClear(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		return module.NewExitError(module.ExitUsage, "Application name required: mac apps cache-clear <appname>")
	}

	appName := ctx.Args[0]

	// Try to resolve bundle ID from the app.
	bundleID := resolveBundleID(appName)

	if bundleID == "" {
		// Use the app name as a heuristic.
		bundleID = appName
	}

	// Build list of cache locations to check.
	home, _ := os.UserHomeDir()
	locations := []struct {
		label string
		path  string
	}{
		{"Cache", filepath.Join(home, "Library", "Caches", bundleID)},
		{"Application Support", filepath.Join(home, "Library", "Application Support", bundleID)},
		{"Preferences", filepath.Join(home, "Library", "Preferences", bundleID+".plist")},
	}

	ctx.Output.Header(fmt.Sprintf("Cache Locations for %s (bundle: %s)", appName, bundleID))
	fmt.Println()

	found := false
	var toDelete []string
	for _, loc := range locations {
		if _, err := os.Stat(loc.path); err == nil {
			found = true
			// Get size.
			sizeStr := "-"
			sizeOut, err := exec.Command("du", "-sh", loc.path).Output()
			if err == nil {
				parts := strings.Fields(strings.TrimSpace(string(sizeOut)))
				if len(parts) > 0 {
					sizeStr = parts[0]
				}
			}
			fmt.Printf("  %-22s %s  (%s)\n", loc.label+":", loc.path, sizeStr)
			toDelete = append(toDelete, loc.path)
		} else {
			fmt.Printf("  %-22s %s  (not found)\n", loc.label+":", loc.path)
		}
	}

	if !found {
		ctx.Output.Warning("No cache files found for '%s'", appName)
		ctx.Output.Info("Try using the full bundle ID (e.g., com.apple.Safari)")
		return nil
	}

	fmt.Println()
	if ctx.DryRun {
		ctx.Output.Info("Dry run — would delete %d location(s):", len(toDelete))
		for _, p := range toDelete {
			fmt.Printf("    rm -rf %s\n", p)
		}
		return nil
	}

	confirmed, err := ctx.Prompt.Confirm(fmt.Sprintf("Delete %d cache location(s) for %s?", len(toDelete), appName))
	if err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Prompt failed: %v", err))
	}
	if !confirmed {
		ctx.Output.Info("Cancelled")
		return nil
	}

	for _, p := range toDelete {
		ctx.Output.Info("Removing %s...", p)
		if err := os.RemoveAll(p); err != nil {
			ctx.Output.Warning("Failed to remove %s: %v", p, err)
		} else {
			ctx.Output.Success("Removed %s", p)
		}
	}

	ctx.Output.Success("Cache cleared for %s", appName)
	return nil
}

// resolveBundleID tries to find the bundle ID for an application name.
func resolveBundleID(appName string) string {
	// Search for the app using mdfind.
	out, err := exec.Command("mdfind", "kMDItemKind == 'Application'", "-name", appName).Output()
	if err != nil {
		return ""
	}

	paths := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(paths) == 0 || paths[0] == "" {
		return ""
	}

	// Read CFBundleIdentifier from the first match.
	return readPlistValue(paths[0], "CFBundleIdentifier")
}

// ---------------------------------------------------------------------------
// recent — Show recently used applications
// ---------------------------------------------------------------------------

func (a *AppsModule) recent(ctx *module.Context) error {
	ctx.Output.Header("Recent Applications")
	fmt.Println()

	// Try reading from Dock's recent-apps.
	out, err := exec.Command("defaults", "read", "com.apple.dock", "recent-apps").Output()
	if err != nil {
		// No recent apps or domain not found.
		ctx.Output.Info("No recent applications recorded in Dock")
		return nil
	}

	// Parse the output — it's a plist-style array.
	// Format: (
	//     { "bundle-id" = "..."; "tile-type" = "..."; },
	//     ...
	// )
	text := string(out)
	entries := parseRecentApps(text)

	if len(entries) == 0 {
		ctx.Output.Info("No recent applications found")
		return nil
	}

	headers := []string{"#", "Bundle ID"}
	rows := make([][]string, 0, len(entries))
	for i, entry := range entries {
		rows = append(rows, []string{fmt.Sprintf("%d", i+1), entry})
	}
	ctx.Output.Table(headers, rows)

	return nil
}

// parseRecentApps extracts bundle IDs from the dock recent-apps defaults output.
func parseRecentApps(text string) []string {
	var results []string
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, `"bundle-id"`) {
			// Extract the value after the = sign.
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				val := strings.TrimSpace(parts[1])
				val = strings.Trim(val, ` ";`)
				if val != "" {
					results = append(results, val)
				}
			}
		}
	}
	return results
}

// ---------------------------------------------------------------------------
// kill — Kill/force-quit an application
// ---------------------------------------------------------------------------

func (a *AppsModule) kill(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		return module.NewExitError(module.ExitUsage, "Application name required: mac apps kill <appname>")
	}

	appName := ctx.Args[0]

	// Strip .app suffix if the user included it.
	appName = strings.TrimSuffix(appName, ".app")

	ctx.Output.Info("Killing '%s'...", appName)

	// Try killall first.
	err := exec.Command("killall", appName).Run()
	if err == nil {
		ctx.Output.Success("'%s' has been terminated", appName)
		return nil
	}

	// Fallback to pkill -x for exact match.
	ctx.Output.Info("killall failed, trying pkill...")
	err = exec.Command("pkill", "-x", appName).Run()
	if err == nil {
		ctx.Output.Success("'%s' has been terminated (via pkill)", appName)
		return nil
	}

	return module.NewExitError(module.ExitNotFound,
		fmt.Sprintf("Could not kill '%s' — it may not be running", appName))
}

// ---------------------------------------------------------------------------
// open — Open an application
// ---------------------------------------------------------------------------

func (a *AppsModule) open(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		return module.NewExitError(module.ExitUsage, "Application name required: mac apps open <appname>")
	}

	appName := ctx.Args[0]

	ctx.Output.Info("Opening '%s'...", appName)

	err := exec.Command("open", "-a", appName).Run()
	if err != nil {
		return module.NewExitError(module.ExitGeneral,
			fmt.Sprintf("Failed to open '%s': %v", appName, err))
	}

	ctx.Output.Success("'%s' opened", appName)
	return nil
}
