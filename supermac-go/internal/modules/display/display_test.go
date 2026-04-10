package display

import (
	"testing"
)

func TestModuleRegistration(t *testing.T) {
	mod := &DisplayModule{}
	if mod.Name() != "display" {
		t.Errorf("expected name 'display', got %q", mod.Name())
	}
	if mod.Emoji() != "🖥️" {
		t.Errorf("expected emoji '🖥️', got %q", mod.Emoji())
	}
	if mod.ShortDescription() == "" {
		t.Error("ShortDescription should not be empty")
	}
}

func TestCommandsCount(t *testing.T) {
	mod := &DisplayModule{}
	cmds := mod.Commands()
	if len(cmds) != 8 {
		t.Errorf("expected 8 commands, got %d", len(cmds))
	}
}

func TestAllCommandsHaveRun(t *testing.T) {
	mod := &DisplayModule{}
	for _, cmd := range mod.Commands() {
		if cmd.Run == nil {
			t.Errorf("command %q missing Run function", cmd.Name)
		}
	}
}

func TestSearch(t *testing.T) {
	mod := &DisplayModule{}
	results := mod.Search("brightness")
	if len(results) == 0 {
		t.Error("expected results for 'brightness'")
	}
	results = mod.Search("dark")
	if len(results) == 0 {
		t.Error("expected results for 'dark'")
	}
	results = mod.Search("nonexistent")
	if len(results) != 0 {
		t.Errorf("expected 0 results for 'nonexistent', got %d", len(results))
	}
}
