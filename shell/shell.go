package shell

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

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
func Run(ctx context.Context, script Script, args []string) error {
	return RunWithIO(ctx, script, args, os.Stdin, os.Stdout, os.Stderr)
}

// RunWithIO executes a shell script with custom I/O streams.
func RunWithIO(ctx context.Context, script Script, args []string, stdin io.Reader, stdout, stderr io.Writer) error {
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

	return runner.Run(ctx, prog)
}
