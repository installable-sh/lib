package shell

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

// Script represents a shell script to execute.
type Script struct {
	Content string
	Name    string
}

// Run executes a shell script with the given arguments.
func Run(ctx context.Context, script Script, args []string, sigCh <-chan os.Signal) error {
	return RunWithIO(ctx, script, args, os.Stdin, os.Stdout, os.Stderr, sigCh)
}

// RunWithIO executes a shell script with custom I/O streams.
// sigCh receives signals that should be forwarded to the running script.
func RunWithIO(ctx context.Context, script Script, args []string, stdin io.Reader, stdout, stderr io.Writer, sigCh <-chan os.Signal) error {
	parser := syntax.NewParser()
	prog, err := parser.Parse(strings.NewReader(script.Content), script.Name)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	// Prepend "--" to args to prevent them from being interpreted as shell options
	params := append([]string{"--"}, args...)
	runner, err := interp.New(
		interp.StdIO(stdin, stdout, stderr),
		interp.Env(expand.ListEnviron(os.Environ()...)),
		interp.Params(params...),
	)
	if err != nil {
		return fmt.Errorf("interpreter error: %w", err)
	}

	// Run the script in a goroutine so we can handle signals
	errCh := make(chan error, 1)
	go func() {
		errCh <- runner.Run(ctx, prog)
	}()

	// Wait for completion or forward signals
	for {
		select {
		case err := <-errCh:
			return err
		case sig := <-sigCh:
			// Forward signal to the process group
			// Using negative PID sends to all processes in the group
			if signum, ok := sig.(syscall.Signal); ok {
				_ = syscall.Kill(-os.Getpid(), signum)
			}
		}
	}
}
