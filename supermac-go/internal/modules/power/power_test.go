package power

import (
	"testing"
)

func TestModuleRegistration(t *testing.T) {
	mod := &PowerModule{}
	if mod.Name() != "power" {
		t.Errorf("expected name 'power', got %q", mod.Name())
	}
	if mod.Emoji() != "⚡" {
		t.Errorf("expected emoji '⚡', got %q", mod.Emoji())
	}
}

func TestCommandsCount(t *testing.T) {
	mod := &PowerModule{}
	cmds := mod.Commands()
	// 20 toggles + 1 status = 21
	if len(cmds) != 21 {
		t.Errorf("expected 21 commands, got %d", len(cmds))
	}
}

func TestStatusCommand(t *testing.T) {
	mod := &PowerModule{}
	cmds := mod.Commands()
	found := false
	for _, cmd := range cmds {
		if cmd.Name == "status" {
			found = true
			break
		}
	}
	if !found {
		t.Error("status command not found")
	}
}

func TestAllTogglesHaveNames(t *testing.T) {
	toggles := allToggles()
	if len(toggles) != 20 {
		t.Errorf("expected 20 toggles, got %d", len(toggles))
	}
	for _, tg := range toggles {
		if tg.name == "" {
			t.Error("toggle missing name")
		}
		if tg.desc == "" {
			t.Errorf("toggle %q missing description", tg.name)
		}
		if tg.getState == nil {
			t.Errorf("toggle %q missing getState", tg.name)
		}
		if tg.setState == nil {
			t.Errorf("toggle %q missing setState", tg.name)
		}
	}
}

func TestIsOn(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"1", true},
		{"true", true},
		{"TRUE", true},
		{"yes", true},
		{"on", true},
		{"enabled", true},
		{"active", true},
		{"0", false},
		{"false", false},
		{"off", false},
		{"disabled", false},
		{"", false},
		{"random", false},
	}
	for _, tt := range tests {
		got := isOn(tt.input)
		if got != tt.want {
			t.Errorf("isOn(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestFormatState(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"1", "enabled"},
		{"true", "enabled"},
		{"0", "disabled"},
		{"false", "disabled"},
		{"", "disabled"},
		{"custom-value", "custom-value"},
	}
	for _, tt := range tests {
		got := formatState(tt.input, nil)
		if got != tt.want {
			t.Errorf("formatState(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestSearch(t *testing.T) {
	mod := &PowerModule{}
	results := mod.Search("caffeinate")
	if len(results) != 1 {
		t.Errorf("expected 1 result for 'caffeinate', got %d", len(results))
	}
	results = mod.Search("nonexistent")
	if len(results) != 0 {
		t.Errorf("expected 0 results for 'nonexistent', got %d", len(results))
	}
}
