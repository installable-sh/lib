package shell

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name       string
		script     Script
		args       []string
		wantStdout string
		wantStderr string
		wantErr    bool
	}{
		{
			name: "echo hello",
			script: Script{
				Content: "echo hello",
				Name:    "test.sh",
			},
			wantStdout: "hello\n",
		},
		{
			name: "echo with args",
			script: Script{
				Content: `echo "arg1=$1 arg2=$2"`,
				Name:    "test.sh",
			},
			args:       []string{"foo", "bar"},
			wantStdout: "arg1=foo arg2=bar\n",
		},
		{
			name: "multiline script",
			script: Script{
				Content: "echo line1\necho line2",
				Name:    "test.sh",
			},
			wantStdout: "line1\nline2\n",
		},
		{
			name: "stderr output",
			script: Script{
				Content: "echo error >&2",
				Name:    "test.sh",
			},
			wantStderr: "error\n",
		},
		{
			name: "exit code",
			script: Script{
				Content: "exit 1",
				Name:    "test.sh",
			},
			wantErr: true,
		},
		{
			name: "syntax error",
			script: Script{
				Content: "if then",
				Name:    "test.sh",
			},
			wantErr: true,
		},
		{
			name: "variable expansion",
			script: Script{
				Content: "FOO=bar; echo $FOO",
				Name:    "test.sh",
			},
			wantStdout: "bar\n",
		},
		{
			name: "command substitution",
			script: Script{
				Content: `echo "result: $(echo nested)"`,
				Name:    "test.sh",
			},
			wantStdout: "result: nested\n",
		},
		{
			name: "conditional",
			script: Script{
				Content: `if [ "a" = "a" ]; then echo yes; else echo no; fi`,
				Name:    "test.sh",
			},
			wantStdout: "yes\n",
		},
		{
			name: "loop",
			script: Script{
				Content: "for i in 1 2 3; do echo $i; done",
				Name:    "test.sh",
			},
			wantStdout: "1\n2\n3\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer

			err := RunWithIO(context.Background(), tt.script, tt.args, strings.NewReader(""), &stdout, &stderr, nil)

			if (err != nil) != tt.wantErr {
				t.Errorf("RunWithIO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if got := stdout.String(); got != tt.wantStdout {
				t.Errorf("RunWithIO() stdout = %q, want %q", got, tt.wantStdout)
			}

			if got := stderr.String(); got != tt.wantStderr {
				t.Errorf("RunWithIO() stderr = %q, want %q", got, tt.wantStderr)
			}
		})
	}
}

func TestRunWithIO_Stdin(t *testing.T) {
	script := Script{
		Content: "read line; echo got: $line",
		Name:    "test.sh",
	}

	var stdout bytes.Buffer
	stdin := strings.NewReader("hello\n")

	err := RunWithIO(context.Background(), script, nil, stdin, &stdout, &bytes.Buffer{}, nil)
	if err != nil {
		t.Fatalf("RunWithIO() error: %v", err)
	}

	want := "got: hello\n"
	if got := stdout.String(); got != want {
		t.Errorf("RunWithIO() stdout = %q, want %q", got, want)
	}
}
