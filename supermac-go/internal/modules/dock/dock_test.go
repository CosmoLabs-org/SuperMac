package dock

import (
	"testing"

	"github.com/cosmolabs-org/supermac/internal/module"
)

func TestModuleRegistration(t *testing.T) {
	mod := &DockModule{}
	if mod.Name() != "dock" {
		t.Errorf("expected name 'dock', got %q", mod.Name())
	}
	if mod.Emoji() != "🚢" {
		t.Errorf("expected emoji '🚢', got %q", mod.Emoji())
	}
	if mod.ShortDescription() == "" {
		t.Error("ShortDescription should not be empty")
	}
}

func TestCommandsCount(t *testing.T) {
	mod := &DockModule{}
	cmds := mod.Commands()
	if len(cmds) != 11 {
		t.Errorf("expected 11 commands, got %d", len(cmds))
	}
}

func TestCommandNames(t *testing.T) {
	mod := &DockModule{}
	cmds := mod.Commands()
	expected := []string{"position", "autohide", "size", "magnification", "magnification-size",
		"minimize-effect", "status", "reset", "add", "list", "remove"}
	for _, name := range expected {
		found := false
		for _, cmd := range cmds {
			if cmd.Name == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected command %q not found", name)
		}
	}
}

func TestAllCommandsHaveRun(t *testing.T) {
	mod := &DockModule{}
	for _, cmd := range mod.Commands() {
		if cmd.Run == nil {
			t.Errorf("command %q missing Run function", cmd.Name)
		}
	}
}

func TestNormalizePosition(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"left", "left"},
		{"l", "left"},
		{"LEFT", "left"},
		{"bottom", "bottom"},
		{"b", "bottom"},
		{"right", "right"},
		{"r", "right"},
		{"top", ""},
		{"", ""},
		{"invalid", ""},
	}
	for _, tt := range tests {
		got := normalizePosition(tt.input)
		if got != tt.want {
			t.Errorf("normalizePosition(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestParseSize(t *testing.T) {
	tests := []struct {
		input string
		want  int
		label string
	}{
		{"small", 32, "small"},
		{"s", 32, "small"},
		{"medium", 64, "medium"},
		{"m", 64, "medium"},
		{"large", 96, "large"},
		{"l", 96, "large"},
		{"48", 48, "48px"},
		{"256", 256, "256px"},
		{"15", 0, ""},  // too small
		{"300", 0, ""}, // too large
		{"abc", 0, ""}, // not a number
	}
	for _, tt := range tests {
		val, label := parseSize(tt.input)
		if val != tt.want || label != tt.label {
			t.Errorf("parseSize(%q) = (%d, %q), want (%d, %q)", tt.input, val, label, tt.want, tt.label)
		}
	}
}

func TestSizeLabel(t *testing.T) {
	tests := []struct {
		input int
		want  string
	}{
		{32, "small"},
		{40, "small"},
		{64, "medium"},
		{80, "medium"},
		{96, "large"},
		{128, "large"},
	}
	for _, tt := range tests {
		got := sizeLabel(tt.input)
		if got != tt.want {
			t.Errorf("sizeLabel(%d) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestStateString(t *testing.T) {
	if stateString(true) != "enabled" {
		t.Error("stateString(true) should be 'enabled'")
	}
	if stateString(false) != "disabled" {
		t.Error("stateString(false) should be 'disabled'")
	}
}

func TestSearch(t *testing.T) {
	mod := &DockModule{}
	results := mod.Search("position")
	if len(results) == 0 {
		t.Error("expected results for 'position'")
	}
	results = mod.Search("nonexistent")
	if len(results) != 0 {
		t.Errorf("expected 0 results for 'nonexistent', got %d", len(results))
	}
}

// parseOnOffToggle tests need a module.Context with a mock platform.
// These test the pure parsing logic (on/off paths only).
func TestParseOnOffToggle_DirectValues(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"on", true},
		{"enable", true},
		{"true", true},
		{"1", true},
		{"ON", true},
		{"off", false},
		{"disable", false},
		{"false", false},
		{"0", false},
	}
	for _, tt := range tests {
		got, err := parseOnOffToggle(tt.input, nil, "")
		if err != nil {
			t.Errorf("parseOnOffToggle(%q) unexpected error: %v", tt.input, err)
		}
		if got != tt.want {
			t.Errorf("parseOnOffToggle(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestParseOnOffToggle_Invalid(t *testing.T) {
	_, err := parseOnOffToggle("invalid", nil, "")
	if err == nil {
		t.Error("expected error for invalid input")
	}
	// Should be an ExitError
	if exitErr, ok := err.(*module.ExitError); !ok {
		t.Errorf("expected ExitError, got %T", err)
	} else if exitErr.Code != module.ExitUsage {
		t.Errorf("expected ExitUsage code, got %d", exitErr.Code)
	}
}
