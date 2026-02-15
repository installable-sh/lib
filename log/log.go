package log

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// DebugLogger is an interface for logging used by lib packages.
// Pass nil to disable logging in functions that accept this interface.
type DebugLogger interface {
	Debugf(format string, args ...any)
	Errorf(format string, args ...any) error
}

// Logger provides leveled logging with debug support.
type Logger struct {
	debug  bool
	output io.Writer
	prefix string
}

// New creates a new Logger. Debug output is disabled by default.
func New(prefix string) *Logger {
	return &Logger{
		debug:  false,
		output: os.Stderr,
		prefix: prefix,
	}
}

// SetDebug enables or disables debug output.
func (l *Logger) SetDebug(enabled bool) {
	l.debug = enabled
}

// SetOutput sets the output writer for log messages.
func (l *Logger) SetOutput(w io.Writer) {
	l.output = w
}

// Debug returns true if debug mode is enabled.
func (l *Logger) Debug() bool {
	return l.debug
}

// Debugf prints a debug message if debug mode is enabled.
func (l *Logger) Debugf(format string, args ...any) {
	if l.debug {
		_, _ = fmt.Fprintf(l.output, "[%s] [debug] %s\n", l.prefix, fmt.Sprintf(format, args...))
	}
}

// Infof prints an info message.
func (l *Logger) Infof(format string, args ...any) {
	_, _ = fmt.Fprintf(l.output, "[%s] %s\n", l.prefix, fmt.Sprintf(format, args...))
}

// Errorf prints an error message and returns an error with the same message.
// Supports %w for error wrapping.
func (l *Logger) Errorf(format string, args ...any) error {
	// Replace %w with %v for logging (fmt.Sprintf doesn't handle %w)
	logFormat := strings.ReplaceAll(format, "%w", "%v")
	_, _ = fmt.Fprintf(l.output, "[%s] error: %s\n", l.prefix, fmt.Sprintf(logFormat, args...))
	//nolint:govet // format string comes from caller, supports %w wrapping
	return fmt.Errorf(format, args...)
}
