package finder

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/cosmolabs-org/supermac/internal/module"
)

func init() {
	module.Register(&FinderModule{})
}

type FinderModule struct{}

func (f *FinderModule) Name() string                { return "finder" }
func (f *FinderModule) ShortDescription() string     { return "File visibility and Finder management" }
func (f *FinderModule) Emoji() string                { return "📁" }

func (f *FinderModule) Commands() []module.Command {
	return []module.Command{
		{
			Name:        "restart",
			Description: "Restart Finder",
			Run:         f.restart,
		},
		{
			Name:        "show-hidden",
			Description: "Show hidden files in Finder",
			Run:         f.showHidden,
		},
		{
			Name:        "hide-hidden",
			Description: "Hide hidden files in Finder",
			Run:         f.hideHidden,
		},
		{
			Name:        "toggle-hidden",
			Description: "Toggle hidden files visibility",
			Run:         f.toggleHidden,
		},
		{
			Name:        "reveal",
			Description: "Reveal file in Finder",
			Args: []module.Arg{
				{Name: "path", Required: true, Description: "Path to reveal"},
			},
			Run: f.reveal,
		},
		{
			Name:        "status",
			Description: "Show Finder status and settings",
			Run:         f.status,
		},
	}
}

func (f *FinderModule) Search(term string) []module.SearchResult {
	var results []module.SearchResult
	for _, cmd := range f.Commands() {
		if strings.Contains(cmd.Name, term) || strings.Contains(strings.ToLower(cmd.Description), term) {
			results = append(results, module.SearchResult{
				Command:     cmd.Name,
				Description: cmd.Description,
				Module:      f.Name(),
			})
		}
	}
	return results
}

func (f *FinderModule) restart(ctx *module.Context) error {
	ctx.Output.Info("Restarting Finder...")
	cmd := exec.Command("killall", "Finder")
	if err := cmd.Run(); err != nil {
		return module.NewExitError(module.ExitGeneral, fmt.Sprintf("Failed to restart Finder: %v", err))
	}
	ctx.Output.Success("Finder restarted")
	return nil
}

func (f *FinderModule) showHidden(ctx *module.Context) error {
	ctx.Output.Info("Showing hidden files...")
	if err := setAppleShowAllFiles(true, ctx); err != nil {
		return err
	}
	restartFinder(ctx)
	ctx.Output.Success("Hidden files are now visible")
	return nil
}

func (f *FinderModule) hideHidden(ctx *module.Context) error {
	ctx.Output.Info("Hiding hidden files...")
	if err := setAppleShowAllFiles(false, ctx); err != nil {
		return err
	}
	restartFinder(ctx)
	ctx.Output.Success("Hidden files are now hidden")
	return nil
}

func (f *FinderModule) toggleHidden(ctx *module.Context) error {
	visible, err := getAppleShowAllFiles(ctx)
	if err != nil {
		return err
	}
	if visible {
		return f.hideHidden(ctx)
	}
	return f.showHidden(ctx)
}

func (f *FinderModule) reveal(ctx *module.Context) error {
	if len(ctx.Args) == 0 {
		return module.NewExitError(module.ExitUsage, "Path required: mac finder reveal <path>")
	}
	path := ctx.Args[0]
	ctx.Output.Info("Revealing %s in Finder...", path)
	cmd := exec.Command("open", "-R", path)
	if err := cmd.Run(); err != nil {
		return module.NewExitError(module.ExitNotFound, fmt.Sprintf("Path not found: %s", path))
	}
	ctx.Output.Success("Revealed in Finder")
	return nil
}

func (f *FinderModule) status(ctx *module.Context) error {
	ctx.Output.Header("Finder Status")

	visible, err := getAppleShowAllFiles(ctx)
	if err != nil {
		ctx.Output.Warning("Could not read hidden files state")
	} else {
		state := "hidden"
		if visible {
			state = "visible"
		}
		fmt.Printf("  Hidden files:    %s\n", state)
	}

	// Get Finder version
	cmd := exec.Command("defaults", "read", "/System/Library/CoreServices/Finder.app/Contents/Info", "CFBundleShortVersionString")
	if out, err := cmd.Output(); err == nil {
		fmt.Printf("  Finder version:  %s\n", strings.TrimSpace(string(out)))
	}

	return nil
}

// Helper functions

func getAppleShowAllFiles(ctx *module.Context) (bool, error) {
	out, err := ctx.Platform.ReadDefault("com.apple.finder", "AppleShowAllFiles")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(out) == "true" || strings.TrimSpace(out) == "1", nil
}

func setAppleShowAllFiles(show bool, ctx *module.Context) error {
	val := "false"
	if show {
		val = "true"
	}
	return ctx.Platform.WriteDefault("com.apple.finder", "AppleShowAllFiles", "-bool "+val)
}

func restartFinder(ctx *module.Context) {
	cmd := exec.Command("killall", "Finder")
	cmd.Run() // Best effort
}
