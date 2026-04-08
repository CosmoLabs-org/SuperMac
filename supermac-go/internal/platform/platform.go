package platform

// Interface is the test seam for all macOS system calls.
// Modules never call exec.Command directly — they use this interface.
type Interface interface {
	// osascript
	RunOSAScript(script string) (string, error)

	// defaults read/write/delete
	ReadDefault(domain, key string) (string, error)
	WriteDefault(domain, key, value string) error
	DeleteDefault(domain, key string) error

	// Network
	SetWiFi(on bool) error
	GetWiFiStatus() (*WiFiInfo, error)
	ScanWiFiNetworks() ([]Network, error)
	FlushDNS() error
	ResetNetwork() error

	// System
	GetMemoryInfo() (*MemoryInfo, error)
	GetCPUInfo() (*CPUInfo, error)
	GetBatteryInfo() (*BatteryInfo, error)
	GetHardwareInfo() (*HardwareInfo, error)
	GetPageSize() (int, error)

	// Display
	SetBrightness(level float64) error
	GetDarkMode() (bool, error)
	SetDarkMode(on bool) error

	// Audio
	GetVolume() (int, error)
	SetVolume(level int) error
	GetAudioDevices() ([]AudioDevice, error)

	// Process management
	ListProcesses(filter string) ([]Process, error)
	KillPort(port int) error
	GetPortUser(port int) (string, error)

	// General command execution
	RunCommand(name string, args ...string) (string, error)
	RunSudoCommand(name string, args ...string) (string, error)
}

// Data types returned by platform methods.

type WiFiInfo struct {
	SSID     string
	BSSID    string
	Signal   int
	Channel  string
	Connected bool
}

type Network struct {
	SSID   string
	Signal int
	Security string
}

type MemoryInfo struct {
	Total     uint64
	Used      uint64
	Free      uint64
	Active    uint64
	Inactive  uint64
	Wired     uint64
	Compressed uint64
	SwapTotal uint64
	SwapUsed  uint64
}

type CPUInfo struct {
	Model     string
	Cores     int
	Threads   int
	Usage     float64
}

type BatteryInfo struct {
	Percent      int
	Charging     bool
	Health       string
	CycleCount   int
	TimeRemaining string
}

type HardwareInfo struct {
	Model       string
	Chip        string
	Memory      string
	Serial      string
	OSVersion   string
	Build       string
	Arch        string
}

type AudioDevice struct {
	ID    string
	Name  string
	Type  string // input, output
	Active bool
}

type Process struct {
	PID     int
	User    string
	CPU     float64
	Memory  float64
	Command string
}
