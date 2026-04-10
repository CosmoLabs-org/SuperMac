package dev

import (
	"testing"
)

func TestModuleRegistration(t *testing.T) {
	mod := &DevModule{}
	if mod.Name() != "dev" {
		t.Errorf("expected name 'dev', got %q", mod.Name())
	}
	if mod.Emoji() != "💻" {
		t.Errorf("expected emoji '💻', got %q", mod.Emoji())
	}
	if mod.ShortDescription() == "" {
		t.Error("ShortDescription should not be empty")
	}
}

func TestCommandsCount(t *testing.T) {
	mod := &DevModule{}
	cmds := mod.Commands()
	if len(cmds) != 16 {
		t.Errorf("expected 16 commands, got %d", len(cmds))
	}
}

func TestAllCommandsHaveRun(t *testing.T) {
	mod := &DevModule{}
	for _, cmd := range mod.Commands() {
		if cmd.Run == nil {
			t.Errorf("command %q missing Run function", cmd.Name)
		}
	}
}

func TestSearch(t *testing.T) {
	mod := &DevModule{}
	results := mod.Search("kill")
	if len(results) == 0 {
		t.Error("expected results for 'kill'")
	}
	results = mod.Search("uuid")
	if len(results) == 0 {
		t.Error("expected results for 'uuid'")
	}
	results = mod.Search("nonexistent")
	if len(results) != 0 {
		t.Errorf("expected 0 results for 'nonexistent', got %d", len(results))
	}
}
