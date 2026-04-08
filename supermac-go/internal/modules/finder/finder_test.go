package finder

import (
	"bytes"
	"strings"
	"testing"

	"github.com/cosmolabs-org/supermac/internal/module"
	"github.com/cosmolabs-org/supermac/internal/output"
	"github.com/cosmolabs-org/supermac/internal/platform"
)

func testContext(defaults map[string]string, args ...string) *module.Context {
	buf := &bytes.Buffer{}
	return &module.Context{
		Output:   output.NewWriter("plain", buf),
		Platform: &stubPlatform{defaults: defaults},
		Args:     args,
	}
}

// stubPlatform records defaults reads/writes for assertions.
type stubPlatform struct {
	defaults       map[string]string
	writes         []platform.DefaultWrite
	osaScriptCalls []string
}

func (s *stubPlatform) RunOSAScript(script string) (string, error) {
	s.osaScriptCalls = append(s.osaScriptCalls, script)
	return "", nil
}

func (s *stubPlatform) ReadDefault(domain, key string) (string, error) {
	if v, ok := s.defaults[domain+"."+key]; ok {
		return v, nil
	}
	return "", module.NewExitError(module.ExitGeneral, "not found")
}

func (s *stubPlatform) WriteDefault(domain, key, value string) error {
	s.writes = append(s.writes, platform.DefaultWrite{Domain: domain, Key: key, Value: value})
	return nil
}

func (s *stubPlatform) DeleteDefault(domain, key string) error { return nil }
func (s *stubPlatform) SetWiFi(bool) error                     { return nil }
func (s *stubPlatform) GetWiFiStatus() (*platform.WiFiInfo, error) {
	return &platform.WiFiInfo{}, nil
}
func (s *stubPlatform) ScanWiFiNetworks() ([]platform.Network, error) {
	return nil, nil
}
func (s *stubPlatform) FlushDNS() error            { return nil }
func (s *stubPlatform) ResetNetwork() error         { return nil }
func (s *stubPlatform) GetMemoryInfo() (*platform.MemoryInfo, error) {
	return &platform.MemoryInfo{}, nil
}
func (s *stubPlatform) GetCPUInfo() (*platform.CPUInfo, error) {
	return &platform.CPUInfo{}, nil
}
func (s *stubPlatform) GetBatteryInfo() (*platform.BatteryInfo, error) {
	return &platform.BatteryInfo{}, nil
}
func (s *stubPlatform) GetHardwareInfo() (*platform.HardwareInfo, error) {
	return &platform.HardwareInfo{}, nil
}
func (s *stubPlatform) GetPageSize() (int, error) { return 16384, nil }
func (s *stubPlatform) SetBrightness(float64) error { return nil }
func (s *stubPlatform) GetDarkMode() (bool, error)   { return false, nil }
func (s *stubPlatform) SetDarkMode(bool) error        { return nil }
func (s *stubPlatform) GetVolume() (int, error)       { return 50, nil }
func (s *stubPlatform) SetVolume(int) error            { return nil }
func (s *stubPlatform) GetAudioDevices() ([]platform.AudioDevice, error) {
	return nil, nil
}
func (s *stubPlatform) ListProcesses(string) ([]platform.Process, error) { return nil, nil }
func (s *stubPlatform) KillPort(int) error                                { return nil }
func (s *stubPlatform) GetPortUser(int) (string, error) {
	return "", nil
}
func (s *stubPlatform) RunCommand(string, ...string) (string, error) { return "", nil }
func (s *stubPlatform) RunSudoCommand(string, ...string) (string, error) {
	return "", nil
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestModuleRegistration(t *testing.T) {
	mod := &FinderModule{}
	if mod.Name() != "finder" {
		t.Errorf("expected name 'finder', got %q", mod.Name())
	}
	if mod.Emoji() != "📁" {
		t.Errorf("expected emoji '📁', got %q", mod.Emoji())
	}
	cmds := mod.Commands()
	if len(cmds) != 6 {
		t.Errorf("expected 6 commands, got %d", len(cmds))
	}
}

func TestShowHidden(t *testing.T) {
	ctx := testContext(map[string]string{
		"com.apple.finder.AppleShowAllFiles": "false",
	})
	mod := &FinderModule{}
	if err := mod.showHidden(ctx); err != nil {
		t.Fatalf("showHidden failed: %v", err)
	}
	stub := ctx.Platform.(*stubPlatform)
	if len(stub.writes) != 1 {
		t.Fatalf("expected 1 write, got %d", len(stub.writes))
	}
	w := stub.writes[0]
	if w.Domain != "com.apple.finder" || w.Key != "AppleShowAllFiles" {
		t.Errorf("wrong write: %+v", w)
	}
	if !strings.Contains(w.Value, "true") {
		t.Errorf("expected value containing 'true', got %q", w.Value)
	}
}

func TestHideHidden(t *testing.T) {
	ctx := testContext(map[string]string{
		"com.apple.finder.AppleShowAllFiles": "true",
	})
	mod := &FinderModule{}
	if err := mod.hideHidden(ctx); err != nil {
		t.Fatalf("hideHidden failed: %v", err)
	}
	stub := ctx.Platform.(*stubPlatform)
	if len(stub.writes) != 1 {
		t.Fatalf("expected 1 write, got %d", len(stub.writes))
	}
	w := stub.writes[0]
	if !strings.Contains(w.Value, "false") {
		t.Errorf("expected value containing 'false', got %q", w.Value)
	}
}

func TestToggleHidden_FromVisible(t *testing.T) {
	ctx := testContext(map[string]string{
		"com.apple.finder.AppleShowAllFiles": "true",
	})
	mod := &FinderModule{}
	if err := mod.toggleHidden(ctx); err != nil {
		t.Fatalf("toggleHidden failed: %v", err)
	}
	stub := ctx.Platform.(*stubPlatform)
	if len(stub.writes) != 1 {
		t.Fatalf("expected 1 write, got %d", len(stub.writes))
	}
	if !strings.Contains(stub.writes[0].Value, "false") {
		t.Errorf("toggle from visible should write false, got %q", stub.writes[0].Value)
	}
}

func TestToggleHidden_FromHidden(t *testing.T) {
	ctx := testContext(map[string]string{
		"com.apple.finder.AppleShowAllFiles": "false",
	})
	mod := &FinderModule{}
	if err := mod.toggleHidden(ctx); err != nil {
		t.Fatalf("toggleHidden failed: %v", err)
	}
	stub := ctx.Platform.(*stubPlatform)
	if !strings.Contains(stub.writes[0].Value, "true") {
		t.Errorf("toggle from hidden should write true, got %q", stub.writes[0].Value)
	}
}

func TestStatus_Visible(t *testing.T) {
	ctx := testContext(map[string]string{
		"com.apple.finder.AppleShowAllFiles": "true",
	})
	mod := &FinderModule{}
	if err := mod.status(ctx); err != nil {
		t.Fatalf("status failed: %v", err)
	}
	// status uses fmt.Printf directly, so we can't easily capture it.
	// The test verifies no error occurs.
}

func TestStatus_Hidden(t *testing.T) {
	ctx := testContext(map[string]string{
		"com.apple.finder.AppleShowAllFiles": "false",
	})
	mod := &FinderModule{}
	if err := mod.status(ctx); err != nil {
		t.Fatalf("status failed: %v", err)
	}
}

func TestReveal_NoArgs(t *testing.T) {
	ctx := testContext(nil)
	mod := &FinderModule{}
	err := mod.reveal(ctx)
	if err == nil {
		t.Fatal("expected error when no path provided")
	}
}

func TestSearch(t *testing.T) {
	mod := &FinderModule{}
	results := mod.Search("hidden")
	if len(results) < 2 {
		t.Errorf("expected at least 2 results for 'hidden', got %d", len(results))
	}
	results = mod.Search("nonexistent")
	if len(results) != 0 {
		t.Errorf("expected 0 results for 'nonexistent', got %d", len(results))
	}
}
