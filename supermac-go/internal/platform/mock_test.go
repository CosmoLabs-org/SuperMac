package platform

import (
	"testing"
)

func TestMockWiFiStatus(t *testing.T) {
	m := &MockPlatform{
		WiFiStatus: &WiFiInfo{SSID: "TestNet", Connected: true},
	}
	info, err := m.GetWiFiStatus()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.SSID != "TestNet" {
		t.Errorf("expected TestNet, got %s", info.SSID)
	}
	if !info.Connected {
		t.Error("expected connected")
	}
}

func TestMockMemory(t *testing.T) {
	m := &MockPlatform{
		Memory: &MemoryInfo{Total: 16 * 1024 * 1024 * 1024, Free: 8 * 1024 * 1024 * 1024},
	}
	info, err := m.GetMemoryInfo()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Total != 16*1024*1024*1024 {
		t.Errorf("expected 16GB, got %d", info.Total)
	}
}

func TestMockDefaultPageSize(t *testing.T) {
	m := &MockPlatform{}
	size, err := m.GetPageSize()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if size != 16384 {
		t.Errorf("expected default 16384, got %d", size)
	}
}

func TestMockCustomPageSize(t *testing.T) {
	m := &MockPlatform{PageSize: 4096}
	size, err := m.GetPageSize()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if size != 4096 {
		t.Errorf("expected 4096, got %d", size)
	}
}

func TestMockVolume(t *testing.T) {
	m := &MockPlatform{}
	m.SetVolume(75)
	vol, err := m.GetVolume()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vol != 75 {
		t.Errorf("expected 75, got %d", vol)
	}
}

func TestMockDarkMode(t *testing.T) {
	m := &MockPlatform{}
	m.SetDarkMode(true)
	dark, err := m.GetDarkMode()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !dark {
		t.Error("expected dark mode to be true")
	}
}

func TestMockRecordsSudoCalls(t *testing.T) {
	m := &MockPlatform{}
	m.RunSudoCommand("networksetup", "-setdnsservers", "Wi-Fi", "empty")
	if len(m.SudoCalls) != 1 {
		t.Fatalf("expected 1 sudo call, got %d", len(m.SudoCalls))
	}
	if m.SudoCalls[0][0] != "networksetup" {
		t.Errorf("expected networksetup, got %s", m.SudoCalls[0][0])
	}
}

func TestMockRecordsOSAScript(t *testing.T) {
	m := &MockPlatform{}
	m.RunOSAScript(`tell application "Finder" to quit`)
	if len(m.OSAScriptCalls) != 1 {
		t.Fatalf("expected 1 osascript call, got %d", len(m.OSAScriptCalls))
	}
}

func TestMockRecordsDefaultWrites(t *testing.T) {
	m := &MockPlatform{}
	m.WriteDefault("com.apple.dock", "autohide", "-bool true")
	if len(m.DefaultWrites) != 1 {
		t.Fatalf("expected 1 default write, got %d", len(m.DefaultWrites))
	}
	dw := m.DefaultWrites[0]
	if dw.Domain != "com.apple.dock" || dw.Key != "autohide" {
		t.Errorf("unexpected default write: %+v", dw)
	}
}
