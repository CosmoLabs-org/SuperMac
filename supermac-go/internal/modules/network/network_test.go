package network

import (
	"testing"
)

func TestModuleRegistration(t *testing.T) {
	mod := &NetworkModule{}
	if mod.Name() != "network" {
		t.Errorf("expected name 'network', got %q", mod.Name())
	}
	if mod.Emoji() != "📡" {
		t.Errorf("expected emoji '📡', got %q", mod.Emoji())
	}
	if mod.ShortDescription() == "" {
		t.Error("ShortDescription should not be empty")
	}
}

func TestCommandsCount(t *testing.T) {
	mod := &NetworkModule{}
	cmds := mod.Commands()
	if len(cmds) != 12 {
		t.Errorf("expected 12 commands, got %d", len(cmds))
	}
}

func TestAllCommandsHaveRun(t *testing.T) {
	mod := &NetworkModule{}
	for _, cmd := range mod.Commands() {
		if cmd.Run == nil {
			t.Errorf("command %q missing Run function", cmd.Name)
		}
	}
}

func TestSearch(t *testing.T) {
	mod := &NetworkModule{}
	results := mod.Search("flush")
	if len(results) == 0 {
		t.Error("expected results for 'flush'")
	}
	results = mod.Search("speed")
	if len(results) == 0 {
		t.Error("expected results for 'speed'")
	}
	results = mod.Search("nonexistent")
	if len(results) != 0 {
		t.Errorf("expected 0 results for 'nonexistent', got %d", len(results))
	}
}
