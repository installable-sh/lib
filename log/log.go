package log

import (
	"fmt"
	"io"
	"os"
)

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
		fmt.Fprintf(l.output, "[%s] [debug] %s\n", l.prefix, fmt.Sprintf(format, args...))
	}
}

// Infof prints an info message.
func (l *Logger) Infof(format string, args ...any) {
	fmt.Fprintf(l.output, "[%s] %s\n", l.prefix, fmt.Sprintf(format, args...))
}

// Errorf prints an error message.
func (l *Logger) Errorf(format string, args ...any) {
	fmt.Fprintf(l.output, "[%s] error: %s\n", l.prefix, fmt.Sprintf(format, args...))
}
