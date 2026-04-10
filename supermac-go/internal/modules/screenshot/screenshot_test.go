package screenshot

import (
	"testing"
)

func TestModuleRegistration(t *testing.T) {
	mod := &ScreenshotModule{}
	if mod.Name() != "screenshot" {
		t.Errorf("expected name 'screenshot', got %q", mod.Name())
	}
	if mod.Emoji() != "📸" {
		t.Errorf("expected emoji '📸', got %q", mod.Emoji())
	}
	if mod.ShortDescription() == "" {
		t.Error("ShortDescription should not be empty")
	}
}

func TestCommandsCount(t *testing.T) {
	mod := &ScreenshotModule{}
	cmds := mod.Commands()
	if len(cmds) != 11 {
		t.Errorf("expected 11 commands, got %d", len(cmds))
	}
}

func TestAllCommandsHaveRun(t *testing.T) {
	mod := &ScreenshotModule{}
	for _, cmd := range mod.Commands() {
		if cmd.Run == nil {
			t.Errorf("command %q missing Run function", cmd.Name)
		}
	}
}

func TestIsTrue(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"true", true},
		{"1", true},
		{"YES", true},
		{"false", false},
		{"0", false},
		{"no", false},
		{"", false},
		{"maybe", false},
	}
	for _, tt := range tests {
		got := isTrue(tt.input)
		if got != tt.want {
			t.Errorf("isTrue(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestValidFormats(t *testing.T) {
	tests := []struct {
		input string
		want  string
		ok    bool
	}{
		{"png", "png", true},
		{"jpg", "jpg", true},
		{"jpeg", "jpg", true},
		{"tiff", "tiff", true},
		{"tif", "tiff", true},
		{"gif", "gif", true},
		{"bmp", "", false},
		{"", "", false},
	}
	for _, tt := range tests {
		got, ok := validFormats[tt.input]
		if ok != tt.ok {
			t.Errorf("validFormats[%q] exists=%v, want %v", tt.input, ok, tt.ok)
		}
		if ok && got != tt.want {
			t.Errorf("validFormats[%q] = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestSearch(t *testing.T) {
	mod := &ScreenshotModule{}
	results := mod.Search("format")
	if len(results) == 0 {
		t.Error("expected results for 'format'")
	}
	results = mod.Search("location")
	if len(results) == 0 {
		t.Error("expected results for 'location'")
	}
	results = mod.Search("nonexistent")
	if len(results) != 0 {
		t.Errorf("expected 0 results for 'nonexistent', got %d", len(results))
	}
}
