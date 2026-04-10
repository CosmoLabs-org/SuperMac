package dep

import (
	"testing"
)

func TestIsInstalled_KnownBinary(t *testing.T) {
	// "ls" exists on every macOS
	d := Dependency{Name: "ls", Check: "ls"}
	if !d.IsInstalled() {
		t.Error("expected ls to be installed")
	}
}

func TestIsInstalled_UnknownBinary(t *testing.T) {
	d := Dependency{Name: "nonexistent-tool-xyz", Check: "nonexistent-tool-xyz"}
	if d.IsInstalled() {
		t.Error("expected nonexistent binary to not be installed")
	}
}

func TestAffectsCommand_NilCommands(t *testing.T) {
	d := Dependency{Name: "test", Commands: nil}
	if !d.AffectsCommand("anything") {
		t.Error("nil Commands should affect all commands")
	}
}

func TestAffectsCommand_SpecificCommands(t *testing.T) {
	d := Dependency{Name: "test", Commands: []string{"add", "remove"}}
	if !d.AffectsCommand("add") {
		t.Error("should affect 'add'")
	}
	if !d.AffectsCommand("remove") {
		t.Error("should affect 'remove'")
	}
	if d.AffectsCommand("list") {
		t.Error("should not affect 'list'")
	}
}

func TestEnsure_NonInteractive(t *testing.T) {
	d := Dependency{Name: "nonexistent-tool-xyz", Brew: "nonexistent", Check: "nonexistent-tool-xyz"}
	err := d.Ensure(false)
	if err == nil {
		t.Error("expected error for missing dep in non-interactive mode")
	}
}

func TestCheckBrew(t *testing.T) {
	path, ok := CheckBrew()
	// brew may or may not be installed in CI, just verify no panic
	_ = path
	_ = ok
}
