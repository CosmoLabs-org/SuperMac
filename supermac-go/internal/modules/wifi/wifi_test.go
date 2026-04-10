package wifi

import (
	"testing"
)

func TestModuleRegistration(t *testing.T) {
	mod := &WiFiModule{}
	if mod.Name() != "wifi" {
		t.Errorf("expected name 'wifi', got %q", mod.Name())
	}
	if mod.Emoji() != "🌐" {
		t.Errorf("expected emoji '🌐', got %q", mod.Emoji())
	}
	if mod.ShortDescription() == "" {
		t.Error("ShortDescription should not be empty")
	}
}

func TestCommandsCount(t *testing.T) {
	mod := &WiFiModule{}
	cmds := mod.Commands()
	if len(cmds) != 9 {
		t.Errorf("expected 9 commands, got %d", len(cmds))
	}
}

func TestAllCommandsHaveRun(t *testing.T) {
	mod := &WiFiModule{}
	for _, cmd := range mod.Commands() {
		if cmd.Run == nil {
			t.Errorf("command %q missing Run function", cmd.Name)
		}
	}
}

func TestSearch(t *testing.T) {
	mod := &WiFiModule{}
	results := mod.Search("scan")
	if len(results) == 0 {
		t.Error("expected results for 'scan'")
	}
	results = mod.Search("connect")
	if len(results) == 0 {
		t.Error("expected results for 'connect'")
	}
	results = mod.Search("nonexistent")
	if len(results) != 0 {
		t.Errorf("expected 0 results for 'nonexistent', got %d", len(results))
	}
}
