package shell

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/installable-sh/lib/log"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

// Script represents a shell script to execute.
type Script struct {
	Content string
	Name    string
}

// Run executes a shell script with custom I/O streams.
// Debug output is controlled by the logger's debug level.
func Run(ctx context.Context, script Script, args []string, stdin io.Reader, stdout, stderr io.Writer, logger log.DebugLogger) error {
	logger.Debugf("Parsing script: %s (%d bytes)", script.Name, len(script.Content))
	parser := syntax.NewParser()
	prog, err := parser.Parse(strings.NewReader(script.Content), script.Name)
	if err != nil {
		logger.Errorf("Parse error: %v", err)
		return fmt.Errorf("parse error: %w", err)
	}
	logger.Debugf("Parsed %d statements", len(prog.Stmts))

	// Prepend "--" to args to prevent them from being interpreted as shell options
	params := append([]string{"--"}, args...)
	logger.Debugf("Script arguments: %v", args)

	logger.Debugf("Creating shell interpreter")
	runner, err := interp.New(
		interp.StdIO(stdin, stdout, stderr),
		interp.Env(expand.ListEnviron(os.Environ()...)),
		interp.Params(params...),
	)
	if err != nil {
		logger.Errorf("Interpreter error: %v", err)
		return fmt.Errorf("interpreter error: %w", err)
	}

	logger.Debugf("Executing script")
	err = runner.Run(ctx, prog)
	if err != nil {
		logger.Errorf("Execution error: %v", err)
	} else {
		logger.Debugf("Script completed successfully")
	}

	return err
}
