package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	configDir  = ".supermac"
	configFile = "config.yaml"
)

// Config represents the SuperMac configuration.
type Config struct {
	Version int            `yaml:"version"`
	Output  OutputConfig   `yaml:"output"`
	Updates UpdatesConfig  `yaml:"updates"`
	Modules ModulesConfig  `yaml:"modules"`
	Aliases map[string]string `yaml:"aliases"`
}

type OutputConfig struct {
	Color  bool   `yaml:"color"`
	Format string `yaml:"format"` // text, json, quiet
}

type UpdatesConfig struct {
	Check   bool   `yaml:"check"`
	Channel string `yaml:"channel"` // stable, beta
}

type ModulesConfig struct {
	Screenshot ScreenshotConfig `yaml:"screenshot"`
	Audio      AudioConfig      `yaml:"audio"`
	Display    DisplayConfig    `yaml:"display"`
}

type ScreenshotConfig struct {
	Location string `yaml:"location"`
	Format   string `yaml:"format"`
	Shadow   bool   `yaml:"shadow"`
}

type AudioConfig struct {
	VolumeStep int `yaml:"volume_step"`
}

type DisplayConfig struct {
	BrightnessStep int `yaml:"brightness_step"`
}

// Default returns a Config with sensible defaults.
func Default() *Config {
	return &Config{
		Version: 1,
		Output: OutputConfig{
			Color:  true,
			Format: "text",
		},
		Updates: UpdatesConfig{
			Check:   true,
			Channel: "stable",
		},
		Modules: ModulesConfig{
			Screenshot: ScreenshotConfig{
				Location: "Desktop",
				Format:   "PNG",
				Shadow:   false,
			},
			Audio: AudioConfig{
				VolumeStep: 10,
			},
			Display: DisplayConfig{
				BrightnessStep: 10,
			},
		},
		Aliases: map[string]string{
			"kp":   "dev kill-port",
			"dark": "display dark-mode",
		},
	}
}

// ConfigPath returns the full path to the config file.
func ConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, configDir, configFile)
}

// Load reads config from disk, creating defaults if missing.
func Load() (*Config, error) {
	path := ConfigPath()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := Default()
			if saveErr := Save(cfg); saveErr != nil {
				return cfg, saveErr
			}
			return cfg, nil
		}
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Save writes config to disk.
func Save(cfg *Config) error {
	path := ConfigPath()

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
