package platform

import "fmt"

// MockPlatform implements Interface for testing.
// All methods return configurable responses. Calls are recorded for assertion.
type MockPlatform struct {
	WiFiStatus    *WiFiInfo
	WiFIDevices   []Network
	Memory        *MemoryInfo
	CPU           *CPUInfo
	Battery       *BatteryInfo
	Hardware      *HardwareInfo
	Volume        int
	DarkMode      bool
	AudioDevices  []AudioDevice
	Processes     []Process
	PageSize      int

	// Records of calls made
	OSAScriptCalls    []string
	SudoCalls         [][]string
	CommandCalls      [][]string
	DefaultWrites     []DefaultWrite
}

type DefaultWrite struct {
	Domain string
	Key    string
	Value  string
}

func (m *MockPlatform) RunOSAScript(script string) (string, error) {
	m.OSAScriptCalls = append(m.OSAScriptCalls, script)
	return "", nil
}

func (m *MockPlatform) ReadDefault(domain, key string) (string, error) {
	return "", nil
}

func (m *MockPlatform) WriteDefault(domain, key, value string) error {
	m.DefaultWrites = append(m.DefaultWrites, DefaultWrite{domain, key, value})
	return nil
}

func (m *MockPlatform) DeleteDefault(domain, key string) error {
	return nil
}

func (m *MockPlatform) SetWiFi(on bool) error {
	return nil
}

func (m *MockPlatform) GetWiFiStatus() (*WiFiInfo, error) {
	if m.WiFiStatus != nil {
		return m.WiFiStatus, nil
	}
	return &WiFiInfo{}, nil
}

func (m *MockPlatform) ScanWiFiNetworks() ([]Network, error) {
	if m.WiFIDevices != nil {
		return m.WiFIDevices, nil
	}
	return []Network{}, nil
}

func (m *MockPlatform) FlushDNS() error {
	return nil
}

func (m *MockPlatform) ResetNetwork() error {
	return nil
}

func (m *MockPlatform) GetMemoryInfo() (*MemoryInfo, error) {
	if m.Memory != nil {
		return m.Memory, nil
	}
	return &MemoryInfo{}, nil
}

func (m *MockPlatform) GetCPUInfo() (*CPUInfo, error) {
	if m.CPU != nil {
		return m.CPU, nil
	}
	return &CPUInfo{}, nil
}

func (m *MockPlatform) GetBatteryInfo() (*BatteryInfo, error) {
	if m.Battery != nil {
		return m.Battery, nil
	}
	return &BatteryInfo{}, nil
}

func (m *MockPlatform) GetHardwareInfo() (*HardwareInfo, error) {
	if m.Hardware != nil {
		return m.Hardware, nil
	}
	return &HardwareInfo{}, nil
}

func (m *MockPlatform) GetPageSize() (int, error) {
	if m.PageSize > 0 {
		return m.PageSize, nil
	}
	return 16384, nil
}

func (m *MockPlatform) SetBrightness(level float64) error {
	return nil
}

func (m *MockPlatform) GetDarkMode() (bool, error) {
	return m.DarkMode, nil
}

func (m *MockPlatform) SetDarkMode(on bool) error {
	m.DarkMode = on
	return nil
}

func (m *MockPlatform) GetVolume() (int, error) {
	return m.Volume, nil
}

func (m *MockPlatform) SetVolume(level int) error {
	m.Volume = level
	return nil
}

func (m *MockPlatform) GetAudioDevices() ([]AudioDevice, error) {
	if m.AudioDevices != nil {
		return m.AudioDevices, nil
	}
	return []AudioDevice{}, nil
}

func (m *MockPlatform) ListProcesses(filter string) ([]Process, error) {
	if m.Processes != nil {
		return m.Processes, nil
	}
	return []Process{}, nil
}

func (m *MockPlatform) KillPort(port int) error {
	return nil
}

func (m *MockPlatform) GetPortUser(port int) (string, error) {
	return "", fmt.Errorf("no process on port %d", port)
}

func (m *MockPlatform) RunCommand(name string, args ...string) (string, error) {
	call := append([]string{name}, args...)
	m.CommandCalls = append(m.CommandCalls, call)
	return "", nil
}

func (m *MockPlatform) RunSudoCommand(name string, args ...string) (string, error) {
	call := append([]string{name}, args...)
	m.SudoCalls = append(m.SudoCalls, call)
	return "", nil
}
