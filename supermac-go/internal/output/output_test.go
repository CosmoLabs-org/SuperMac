package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestColoredWriterInfo(t *testing.T) {
	var buf bytes.Buffer
	w := &ColoredWriter{w: &buf, color: false}
	w.Info("test %s", "message")

	if !strings.Contains(buf.String(), "test message") {
		t.Errorf("expected output to contain 'test message', got: %s", buf.String())
	}
}

func TestColoredWriterError(t *testing.T) {
	var buf bytes.Buffer
	w := &ColoredWriter{w: &buf, color: false}
	w.Error("failed: %d", 42)

	if !strings.Contains(buf.String(), "failed: 42") {
		t.Errorf("expected output to contain 'failed: 42', got: %s", buf.String())
	}
}

func TestQuietWriterSuppressesInfo(t *testing.T) {
	var buf bytes.Buffer
	w := &QuietWriter{w: &buf}
	w.Info("should not appear")
	w.Success("should not appear")
	w.Warning("should not appear")

	if buf.Len() > 0 {
		t.Errorf("quiet writer should suppress info/success/warning, got: %s", buf.String())
	}
}

func TestQuietWriterShowsErrors(t *testing.T) {
	var buf bytes.Buffer
	w := &QuietWriter{w: &buf}
	w.Error("actual error")

	if !strings.Contains(buf.String(), "actual error") {
		t.Errorf("quiet writer should show errors, got: %s", buf.String())
	}
}

func TestJSONWriterInfo(t *testing.T) {
	var buf bytes.Buffer
	w := &JSONWriter{w: &buf}
	w.Info("test message")

	output := buf.String()
	if !strings.Contains(output, `"type"`) || !strings.Contains(output, "info") {
		t.Errorf("expected JSON type info, got: %s", output)
	}
	if !strings.Contains(output, `"message"`) || !strings.Contains(output, "test message") {
		t.Errorf("expected JSON message, got: %s", output)
	}
}

func TestNewWriterText(t *testing.T) {
	w := NewWriter("text", nil)
	if _, ok := w.(*ColoredWriter); !ok {
		t.Error("expected ColoredWriter for text format")
	}
}

func TestNewWriterJSON(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter("json", &buf)
	if _, ok := w.(*JSONWriter); !ok {
		t.Error("expected JSONWriter for json format")
	}
}

func TestNewWriterQuiet(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter("quiet", &buf)
	if _, ok := w.(*QuietWriter); !ok {
		t.Error("expected QuietWriter for quiet format")
	}
}
