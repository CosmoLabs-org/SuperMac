package audio

import (
	"testing"
)

func TestModuleRegistration(t *testing.T) {
	mod := &AudioModule{}
	if mod.Name() != "audio" {
		t.Errorf("expected name 'audio', got %q", mod.Name())
	}
	if mod.Emoji() != "🔊" {
		t.Errorf("expected emoji '🔊', got %q", mod.Emoji())
	}
	if mod.ShortDescription() == "" {
		t.Error("ShortDescription should not be empty")
	}
}

func TestCommandsCount(t *testing.T) {
	mod := &AudioModule{}
	cmds := mod.Commands()
	if len(cmds) != 11 {
		t.Errorf("expected 11 commands, got %d", len(cmds))
	}
}

func TestAllCommandsHaveRun(t *testing.T) {
	mod := &AudioModule{}
	for _, cmd := range mod.Commands() {
		if cmd.Run == nil {
			t.Errorf("command %q missing Run function", cmd.Name)
		}
	}
}

func TestSearch(t *testing.T) {
	mod := &AudioModule{}
	results := mod.Search("volume")
	if len(results) == 0 {
		t.Error("expected results for 'volume'")
	}
	results = mod.Search("balance")
	if len(results) == 0 {
		t.Error("expected results for 'balance'")
	}
	results = mod.Search("nonexistent")
	if len(results) != 0 {
		t.Errorf("expected 0 results for 'nonexistent', got %d", len(results))
	}
}
