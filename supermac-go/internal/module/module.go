package module

import (
	"github.com/cosmolabs-org/supermac/internal/config"
	"github.com/cosmolabs-org/supermac/internal/output"
	"github.com/cosmolabs-org/supermac/internal/platform"
)

// Module is the interface every SuperMac module must implement.
type Module interface {
	Name() string
	ShortDescription() string
	Emoji() string
	Commands() []Command
	Search(term string) []SearchResult
}

// Command represents a single action within a module.
type Command struct {
	Name        string
	Description string
	Aliases     []string
	Args        []Arg
	Flags       []Flag
	Run         func(ctx *Context) error
}

// Flag defines a per-command flag (e.g., --sort, --force).
type Flag struct {
	Name         string
	Shorthand    string
	DefaultValue string
	Description  string
	Required     bool
}

// Arg defines a positional argument for a command.
type Arg struct {
	Name        string
	Required    bool
	Description string
}

// Context is passed to every command's Run function.
type Context struct {
	Config   *config.Config
	Output   output.Writer
	Platform platform.Interface
	Prompt   PromptInterface
	Args     []string
	Flags    map[string]string
	Verbose  bool
	DryRun   bool
}

// SearchResult represents a matched command from a module search.
type SearchResult struct {
	Command     string
	Description string
	Module      string
}
