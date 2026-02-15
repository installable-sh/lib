package log

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	logger := New("test")
	if logger == nil {
		t.Fatal("New() returned nil")
	}
	if logger.prefix != "test" {
		t.Errorf("prefix = %q, want %q", logger.prefix, "test")
	}
	if logger.debug {
		t.Error("debug should be false by default")
	}
}

func TestLogger_SetDebug(t *testing.T) {
	logger := New("test")
	logger.SetDebug(true)
	if !logger.Debug() {
		t.Error("Debug() should return true after SetDebug(true)")
	}
	logger.SetDebug(false)
	if logger.Debug() {
		t.Error("Debug() should return false after SetDebug(false)")
	}
}

func TestLogger_Debugf(t *testing.T) {
	var buf bytes.Buffer
	logger := New("test")
	logger.SetOutput(&buf)

	// Debug disabled - should not output
	logger.Debugf("message %d", 1)
	if buf.Len() != 0 {
		t.Error("Debugf should not output when debug is disabled")
	}

	// Debug enabled - should output
	logger.SetDebug(true)
	logger.Debugf("message %d", 2)
	output := buf.String()
	if !strings.Contains(output, "[test]") {
		t.Errorf("output should contain prefix: %q", output)
	}
	if !strings.Contains(output, "[debug]") {
		t.Errorf("output should contain [debug]: %q", output)
	}
	if !strings.Contains(output, "message 2") {
		t.Errorf("output should contain message: %q", output)
	}
}

func TestLogger_Infof(t *testing.T) {
	var buf bytes.Buffer
	logger := New("test")
	logger.SetOutput(&buf)

	logger.Infof("info %s", "message")
	output := buf.String()
	if !strings.Contains(output, "[test]") {
		t.Errorf("output should contain prefix: %q", output)
	}
	if !strings.Contains(output, "info message") {
		t.Errorf("output should contain message: %q", output)
	}
}

func TestLogger_Errorf(t *testing.T) {
	var buf bytes.Buffer
	logger := New("test")
	logger.SetOutput(&buf)

	err := logger.Errorf("error %s", "occurred")
	output := buf.String()

	// Check log output
	if !strings.Contains(output, "[test]") {
		t.Errorf("output should contain prefix: %q", output)
	}
	if !strings.Contains(output, "error:") {
		t.Errorf("output should contain 'error:': %q", output)
	}
	if !strings.Contains(output, "error occurred") {
		t.Errorf("output should contain message: %q", output)
	}

	// Check returned error
	if err == nil {
		t.Fatal("Errorf should return an error")
	}
	if err.Error() != "error occurred" {
		t.Errorf("error message = %q, want %q", err.Error(), "error occurred")
	}
}

func TestLogger_Errorf_Wrap(t *testing.T) {
	var buf bytes.Buffer
	logger := New("test")
	logger.SetOutput(&buf)

	origErr := errors.New("original error")
	err := logger.Errorf("wrapped: %w", origErr)

	// Check that the error is wrapped
	if !errors.Is(err, origErr) {
		t.Error("Errorf with %w should wrap the original error")
	}

	// Check log output shows error value, not %w
	output := buf.String()
	if strings.Contains(output, "%w") {
		t.Errorf("output should not contain %%w: %q", output)
	}
	if !strings.Contains(output, "original error") {
		t.Errorf("output should contain the error value: %q", output)
	}
}
