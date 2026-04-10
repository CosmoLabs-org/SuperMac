package system

import (
	"testing"
)

func TestModuleRegistration(t *testing.T) {
	mod := &SystemModule{}
	if mod.Name() != "system" {
		t.Errorf("expected name 'system', got %q", mod.Name())
	}
	if mod.Emoji() != "🖥️" {
		t.Errorf("expected emoji '🖥️', got %q", mod.Emoji())
	}
	if mod.ShortDescription() == "" {
		t.Error("ShortDescription should not be empty")
	}
}

func TestCommandsCount(t *testing.T) {
	mod := &SystemModule{}
	cmds := mod.Commands()
	if len(cmds) != 11 {
		t.Errorf("expected 11 commands, got %d", len(cmds))
	}
}

func TestCommandNames(t *testing.T) {
	mod := &SystemModule{}
	cmds := mod.Commands()
	expected := []string{"info", "cleanup", "battery", "memory", "cpu",
		"hardware", "disk-usage", "processes", "uptime", "updates", "temperature"}
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
	mod := &SystemModule{}
	for _, cmd := range mod.Commands() {
		if cmd.Run == nil {
			t.Errorf("command %q missing Run function", cmd.Name)
		}
	}
}

func TestSearch(t *testing.T) {
	mod := &SystemModule{}
	results := mod.Search("battery")
	if len(results) == 0 {
		t.Error("expected results for 'battery'")
	}
	results = mod.Search("cleanup")
	if len(results) == 0 {
		t.Error("expected results for 'cleanup'")
	}
	results = mod.Search("nonexistent")
	if len(results) != 0 {
		t.Errorf("expected 0 results for 'nonexistent', got %d", len(results))
	}
}
