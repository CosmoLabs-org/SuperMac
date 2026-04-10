package apps

import (
	"testing"
)

func TestModuleRegistration(t *testing.T) {
	mod := &AppsModule{}
	if mod.Name() != "apps" {
		t.Errorf("expected name 'apps', got %q", mod.Name())
	}
	if mod.Emoji() != "📱" {
		t.Errorf("expected emoji '📱', got %q", mod.Emoji())
	}
	if mod.ShortDescription() == "" {
		t.Error("ShortDescription should not be empty")
	}
}

func TestCommandsCount(t *testing.T) {
	mod := &AppsModule{}
	cmds := mod.Commands()
	if len(cmds) != 6 {
		t.Errorf("expected 6 commands, got %d", len(cmds))
	}
}

func TestCommandNames(t *testing.T) {
	mod := &AppsModule{}
	cmds := mod.Commands()
	expected := []string{"list", "info", "cache-clear", "recent", "kill", "open"}
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
	mod := &AppsModule{}
	for _, cmd := range mod.Commands() {
		if cmd.Run == nil {
			t.Errorf("command %q missing Run function", cmd.Name)
		}
	}
}

func TestSearch(t *testing.T) {
	mod := &AppsModule{}
	results := mod.Search("list")
	if len(results) == 0 {
		t.Error("expected results for 'list'")
	}
	results = mod.Search("cache")
	if len(results) == 0 {
		t.Error("expected results for 'cache'")
	}
	results = mod.Search("nonexistent")
	if len(results) != 0 {
		t.Errorf("expected 0 results for 'nonexistent', got %d", len(results))
	}
}
