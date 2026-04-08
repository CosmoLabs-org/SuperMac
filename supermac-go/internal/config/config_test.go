package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg.Version != 1 {
		t.Errorf("expected version 1, got %d", cfg.Version)
	}
	if cfg.Output.Format != "text" {
		t.Errorf("expected text format, got %s", cfg.Output.Format)
	}
	if cfg.Output.Color != true {
		t.Error("expected color to be true")
	}
	if cfg.Updates.Check != true {
		t.Error("expected updates.check to be true")
	}
	if cfg.Modules.Screenshot.Location != "Desktop" {
		t.Errorf("expected Desktop, got %s", cfg.Modules.Screenshot.Location)
	}
	if cfg.Aliases["kp"] != "dev kill-port" {
		t.Error("expected kp alias")
	}
}

func TestSaveAndLoad(t *testing.T) {
	// Use a temp dir to avoid clobbering real config
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	t.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cfg := Default()
	cfg.Output.Format = "json"

	if err := Save(cfg); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Output.Format != "json" {
		t.Errorf("expected json format, got %s", loaded.Output.Format)
	}

	// Verify file exists
	path := filepath.Join(tmpDir, configDir, configFile)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("config file was not created")
	}
}

func TestLoadCreatesDefault(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	// Load with no existing config should create defaults
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Version != 1 {
		t.Errorf("expected default version 1, got %d", cfg.Version)
	}
}
