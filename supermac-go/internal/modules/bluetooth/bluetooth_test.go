package bluetooth

import (
	"testing"
)

func TestModuleRegistration(t *testing.T) {
	mod := &BluetoothModule{}
	if mod.Name() != "bluetooth" {
		t.Errorf("expected name 'bluetooth', got %q", mod.Name())
	}
	if mod.Emoji() != "📡" {
		t.Errorf("expected emoji '📡', got %q", mod.Emoji())
	}
	if mod.ShortDescription() == "" {
		t.Error("ShortDescription should not be empty")
	}
}

func TestCommandsCount(t *testing.T) {
	mod := &BluetoothModule{}
	cmds := mod.Commands()
	if len(cmds) != 6 {
		t.Errorf("expected 6 commands, got %d", len(cmds))
	}
}

func TestAllCommandsHaveRun(t *testing.T) {
	mod := &BluetoothModule{}
	for _, cmd := range mod.Commands() {
		if cmd.Run == nil {
			t.Errorf("command %q missing Run function", cmd.Name)
		}
	}
}

func TestParseBlueutilDeviceLine(t *testing.T) {
	tests := []struct {
		line    string
		wantName string
		wantMAC  string
	}{
		{
			"address: AA:BB:CC:DD:EE:FF, name: AirPods Pro",
			"AirPods Pro",
			"AA:BB:CC:DD:EE:FF",
		},
		{
			"address: 11:22:33:44:55:66, name: Magic Mouse",
			"Magic Mouse",
			"11:22:33:44:55:66",
		},
		{
			"name: Keyboard, address: AA:11:BB:22:CC:33",
			"Keyboard",
			"AA:11:BB:22:CC:33",
		},
		{
			"address: FF:FF:FF:FF:FF:FF",
			"",
			"FF:FF:FF:FF:FF:FF",
		},
		{
			"",
			"",
			"",
		},
	}
	for _, tt := range tests {
		name, mac := parseBlueutilDeviceLine(tt.line)
		if name != tt.wantName {
			t.Errorf("parseBlueutilDeviceLine(%q) name = %q, want %q", tt.line, name, tt.wantName)
		}
		if mac != tt.wantMAC {
			t.Errorf("parseBlueutilDeviceLine(%q) mac = %q, want %q", tt.line, mac, tt.wantMAC)
		}
	}
}

func TestSearch(t *testing.T) {
	mod := &BluetoothModule{}
	results := mod.Search("connect")
	if len(results) == 0 {
		t.Error("expected results for 'connect'")
	}
	results = mod.Search("power")
	if len(results) == 0 {
		t.Error("expected results for 'power'")
	}
	results = mod.Search("nonexistent")
	if len(results) != 0 {
		t.Errorf("expected 0 results for 'nonexistent', got %d", len(results))
	}
}
