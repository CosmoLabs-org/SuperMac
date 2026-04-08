package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// Writer is the interface for all output formatting.
type Writer interface {
	Info(msg string, args ...interface{})
	Success(msg string, args ...interface{})
	Warning(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Header(title string)
	Table(headers []string, rows [][]string)
	JSON(v interface{}) error
}

// NewWriter selects the appropriate writer based on flags and environment.
func NewWriter(format string, w io.Writer) Writer {
	if w == nil {
		w = os.Stdout
	}

	switch format {
	case "json":
		return &JSONWriter{w: w}
	case "quiet":
		return &QuietWriter{w: w}
	default:
		return NewColoredWriter(w)
	}
}

// ColoredWriter produces ANSI-colored terminal output.
type ColoredWriter struct {
	w     io.Writer
	color bool
}

func NewColoredWriter(w io.Writer) *ColoredWriter {
	color := true
	if os.Getenv("NO_COLOR") != "" {
		color = false
	}
	// Disable color when not a TTY
	if f, ok := w.(*os.File); ok {
		if !isTerminal(f) {
			color = false
		}
	}
	return &ColoredWriter{w: w, color: color}
}

func isTerminal(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func (c *ColoredWriter) colorize(code, msg string) string {
	if !c.color {
		return msg
	}
	return fmt.Sprintf("\033[%sm%s\033[0m", code, msg)
}

func (c *ColoredWriter) Info(msg string, args ...interface{}) {
	fmt.Fprintln(c.w, c.colorize("36", fmt.Sprintf("  ℹ️  "+msg, args...)))
}

func (c *ColoredWriter) Success(msg string, args ...interface{}) {
	fmt.Fprintln(c.w, c.colorize("32", fmt.Sprintf("  ✅ "+msg, args...)))
}

func (c *ColoredWriter) Warning(msg string, args ...interface{}) {
	fmt.Fprintln(c.w, c.colorize("33", fmt.Sprintf("  ⚠️  "+msg, args...)))
}

func (c *ColoredWriter) Error(msg string, args ...interface{}) {
	fmt.Fprintln(c.w, c.colorize("31", fmt.Sprintf("  ❌ "+msg, args...)))
}

func (c *ColoredWriter) Header(title string) {
	line := strings.Repeat("─", len(title)+4)
	fmt.Fprintln(c.w, c.colorize("1", fmt.Sprintf("\n┌%s┐", line)))
	fmt.Fprintln(c.w, c.colorize("1", fmt.Sprintf("│ %s │", title)))
	fmt.Fprintln(c.w, c.colorize("1", fmt.Sprintf("└%s┘", line)))
}

func (c *ColoredWriter) Table(headers []string, rows [][]string) {
	// Simple table for now — will improve later
	for _, h := range headers {
		fmt.Fprintf(c.w, "%-20s", c.colorize("1", h))
	}
	fmt.Fprintln(c.w)
	for _, row := range rows {
		for _, cell := range row {
			fmt.Fprintf(c.w, "%-20s", cell)
		}
		fmt.Fprintln(c.w)
	}
}

func (c *ColoredWriter) JSON(v interface{}) error {
	enc := json.NewEncoder(c.w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// JSONWriter produces structured JSON output for all messages.
type JSONWriter struct {
	w io.Writer
}

func (j *JSONWriter) writeJSON(v interface{}) {
	enc := json.NewEncoder(j.w)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}

func (j *JSONWriter) Info(msg string, args ...interface{}) {
	j.writeJSON(map[string]string{"type": "info", "message": fmt.Sprintf(msg, args...)})
}

func (j *JSONWriter) Success(msg string, args ...interface{}) {
	j.writeJSON(map[string]string{"type": "success", "message": fmt.Sprintf(msg, args...)})
}

func (j *JSONWriter) Warning(msg string, args ...interface{}) {
	j.writeJSON(map[string]string{"type": "warning", "message": fmt.Sprintf(msg, args...)})
}

func (j *JSONWriter) Error(msg string, args ...interface{}) {
	j.writeJSON(map[string]string{"type": "error", "message": fmt.Sprintf(msg, args...)})
}

func (j *JSONWriter) Header(title string) {
	j.writeJSON(map[string]string{"type": "header", "title": title})
}

func (j *JSONWriter) Table(headers []string, rows [][]string) {
	j.writeJSON(map[string]interface{}{"type": "table", "headers": headers, "rows": rows})
}

func (j *JSONWriter) JSON(v interface{}) error {
	return json.NewEncoder(j.w).Encode(v)
}

// QuietWriter suppresses all output except errors.
type QuietWriter struct {
	w io.Writer
}

func (q *QuietWriter) Info(string, ...interface{})    {}
func (q *QuietWriter) Success(string, ...interface{}) {}
func (q *QuietWriter) Warning(string, ...interface{}) {}
func (q *QuietWriter) Error(msg string, args ...interface{}) {
	fmt.Fprintln(q.w, fmt.Sprintf(msg, args...))
}
func (q *QuietWriter) Header(string)                                {}
func (q *QuietWriter) Table([]string, [][]string)                   {}
func (q *QuietWriter) JSON(v interface{}) error {
	return json.NewEncoder(q.w).Encode(v)
}
