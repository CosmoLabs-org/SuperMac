package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/cosmolabs-org/supermac/internal/config"
	"github.com/cosmolabs-org/supermac/internal/module"
	"github.com/cosmolabs-org/supermac/internal/output"
	"github.com/cosmolabs-org/supermac/internal/platform"
	"github.com/cosmolabs-org/supermac/internal/version"

	// Module registrations
	_ "github.com/cosmolabs-org/supermac/internal/modules/finder"
	_ "github.com/cosmolabs-org/supermac/internal/modules/dock"
	_ "github.com/cosmolabs-org/supermac/internal/modules/system"
	_ "github.com/cosmolabs-org/supermac/internal/modules/wifi"
	_ "github.com/cosmolabs-org/supermac/internal/modules/network"
	_ "github.com/cosmolabs-org/supermac/internal/modules/display"
	_ "github.com/cosmolabs-org/supermac/internal/modules/dev"
	_ "github.com/cosmolabs-org/supermac/internal/modules/audio"
	_ "github.com/cosmolabs-org/supermac/internal/modules/screenshot"

	"github.com/spf13/cobra"
)

var (
	jsonFlag   bool
	quietFlag  bool
	noColor    bool
	verbose    bool
	dryRun     bool
	yesFlag    bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "mac",
		Short: "SuperMac — macOS power tools for the CLI",
		Long:  "Professional command-line tool for macOS with organized, powerful shortcuts.\nBuilt by CosmoLabs — https://cosmolabs.org",
		Run:   runHelp,
	}

	rootCmd.PersistentFlags().BoolVar(&jsonFlag, "json", false, "Output in JSON format")
	rootCmd.PersistentFlags().BoolVar(&quietFlag, "quiet", false, "Suppress all output except errors")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable color output")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Show what would be done without executing")
	rootCmd.PersistentFlags().BoolVarP(&yesFlag, "yes", "y", false, "Skip confirmation prompts")

	// Built-in commands
	rootCmd.AddCommand(versionCmd())
	rootCmd.AddCommand(configCmd())

	// Register all modules as Cobra subcommands
	registerModules(rootCmd)

	// Global shortcuts (convenience aliases matching the original Bash CLI)
	addShortcuts(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			w := getOutput()
			w.Header(fmt.Sprintf("SuperMac v%s", version.Version))
			fmt.Printf("  Version:    %s\n", version.Version)
			fmt.Printf("  Build:      %s\n", version.BuildDate)
			fmt.Println()

			modules := module.All()
			names := make([]string, 0, len(modules))
			for name := range modules {
				names = append(names, name)
			}
			sort.Strings(names)
			fmt.Printf("  Modules:    %s\n", strings.Join(names, ", "))
		},
	}
}

func configCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage SuperMac configuration",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "Show current configuration",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Load()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
				os.Exit(1)
			}
			w := getOutput()
			if jsonFlag {
				w.JSON(cfg)
				return
			}
			w.Header("SuperMac Configuration")
			fmt.Printf("  Format:   %s\n", cfg.Output.Format)
			fmt.Printf("  Color:    %v\n", cfg.Output.Color)
			fmt.Printf("  Updates:  %v (%s)\n", cfg.Output.Format, cfg.Updates.Channel)
			if len(cfg.Aliases) > 0 {
				fmt.Println("  Aliases:")
				for k, v := range cfg.Aliases {
					fmt.Printf("    %-15s → %s\n", k, v)
				}
			}
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "get <key>",
		Short: "Get a config value",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Load()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			// Simple dot-path lookup
			switch args[0] {
			case "output.format":
				fmt.Println(cfg.Output.Format)
			case "output.color":
				fmt.Println(cfg.Output.Color)
			case "updates.check":
				fmt.Println(cfg.Updates.Check)
			case "updates.channel":
				fmt.Println(cfg.Updates.Channel)
			default:
				fmt.Fprintf(os.Stderr, "Unknown config key: %s\n", args[0])
				os.Exit(1)
			}
		},
	})

	return cmd
}

func registerModules(rootCmd *cobra.Command) {
	modules := module.All()

	for _, mod := range modules {
		modCmd := &cobra.Command{
			Use:   mod.Name(),
			Short: fmt.Sprintf("%s %s", mod.Emoji(), mod.ShortDescription()),
		}

		commands := mod.Commands()
		for _, cmd := range commands {
			cmd := cmd // capture
			subCmd := &cobra.Command{
				Use:   cmd.Name,
				Short: cmd.Description,
				Aliases: cmd.Aliases,
				RunE: func(subCmd *cobra.Command, args []string) error {
					ctx := &module.Context{
						Config:   loadConfig(),
						Output:   getOutput(),
						Platform: getPlatform(),
						Prompt:   getPrompt(),
						Args:     args,
						Flags:    make(map[string]string),
						Verbose:  verbose,
						DryRun:   dryRun,
					}
					return cmd.Run(ctx)
				},
			}
			// Register per-command flags
			for _, flag := range cmd.Flags {
				subCmd.Flags().StringP(flag.Name, flag.Shorthand, flag.DefaultValue, flag.Description)
				if flag.Required {
					subCmd.MarkFlagRequired(flag.Name)
				}
			}
			modCmd.AddCommand(subCmd)
		}

		rootCmd.AddCommand(modCmd)
	}
}

func runHelp(cmd *cobra.Command, args []string) {
	w := getOutput()
	w.Header(fmt.Sprintf("SuperMac v%s", version.Version))
	fmt.Println()
	fmt.Println("Usage: mac <category> <action> [arguments]")
	fmt.Println()

	modules := module.All()
	if len(modules) > 0 {
		names := make([]string, 0, len(modules))
		for name := range modules {
			names = append(names, name)
		}
		sort.Strings(names)

		fmt.Println("Available Categories:")
		for _, name := range names {
			mod := modules[name]
			fmt.Printf("  %-12s %s %s\n", name, mod.Emoji(), mod.ShortDescription())
		}
		fmt.Println()
	}

	fmt.Println("Built-in Commands:")
	fmt.Println("  help           Show this help")
	fmt.Println("  version        Show version information")
	fmt.Println("  config list    Show current configuration")
	fmt.Println()
	fmt.Println("Global Flags:")
	fmt.Println("  --json         Output in JSON format")
	fmt.Println("  --quiet        Suppress output except errors")
	fmt.Println("  --no-color     Disable color output")
	fmt.Println("  --verbose      Verbose output")
	fmt.Println("  --dry-run      Show what would be done")
	fmt.Println("  --yes          Skip confirmation prompts")
	fmt.Println()
	fmt.Println("Built by CosmoLabs — https://cosmolabs.org")
}

func loadConfig() *config.Config {
	cfg, err := config.Load()
	if err != nil {
		return config.Default()
	}
	return cfg
}

func getOutput() output.Writer {
	format := "text"
	if jsonFlag {
		format = "json"
	} else if quietFlag {
		format = "quiet"
	}
	return output.NewWriter(format, os.Stdout)
}

func getPlatform() platform.Interface {
	return &platform.DarwinPlatform{}
}

func getPrompt() module.PromptInterface {
	if yesFlag {
		return &module.AutoYesPrompt{}
	}
	return &module.TerminalPrompt{}
}

// addShortcuts creates top-level convenience commands that delegate to module subcommands.
// These mirror the original Bash CLI's global shortcuts.
func addShortcuts(rootCmd *cobra.Command) {
	shortcuts := map[string]string{
		"ip":             "network ip",
		"cleanup":        "system cleanup",
		"restart-finder": "finder restart",
		"kp":             "dev kill-port",
		"vol":            "audio volume",
		"dark":           "display dark-mode on",
		"light":          "display dark-mode off",
		"search":         "", // special: built-in search
	}

	for name, target := range shortcuts {
		if name == "search" {
			// Search command: find commands by keyword
			rootCmd.AddCommand(&cobra.Command{
				Use:   "search <term>",
				Short: "Search for a command by keyword",
				Args:  cobra.MinimumNArgs(1),
				Run: func(cmd *cobra.Command, args []string) {
					term := strings.ToLower(args[0])
					w := getOutput()
					w.Header("Search Results")
					fmt.Println()
					found := false
					for _, mod := range module.All() {
						for _, r := range mod.Search(term) {
							fmt.Printf("  %-25s %s (%s)\n", r.Module+" "+r.Command, r.Description, r.Module)
							found = true
						}
					}
					if !found {
						fmt.Printf("  No commands found matching '%s'\n", args[0])
					}
				},
			})
			continue
		}

		parts := strings.SplitN(target, " ", 2)
		moduleName := parts[0]
		subCmd := ""
		if len(parts) > 1 {
			subCmd = parts[1]
		}
		desc := fmt.Sprintf("Shortcut for: mac %s", target)

		shortcutName := name
		shortcutModule := moduleName
		shortcutSubCmd := subCmd

		rootCmd.AddCommand(&cobra.Command{
			Use:   shortcutName + " [args...]",
			Short: desc,
			Run: func(cmd *cobra.Command, args []string) {
				// Find the module command and its subcommand, then execute
				for _, c := range rootCmd.Commands() {
					if c.Name() == shortcutModule {
						for _, sc := range c.Commands() {
							if sc.Name() == shortcutSubCmd {
								sc.RunE(sc, args)
								return
							}
						}
					}
				}
				fmt.Fprintf(os.Stderr, "Error: shortcut target 'mac %s' not found\n", shortcutModule+" "+shortcutSubCmd)
				os.Exit(1)
			},
		})
	}
}
