package version

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	// Save original V and restore after test
	origV := V
	defer func() { V = origV }()

	t.Run("with V set", func(t *testing.T) {
		V = "2.0.0"
		if got := Get(); got != "2.0.0" {
			t.Errorf("Get() = %q, want %q", got, "2.0.0")
		}
	})

	t.Run("without V set", func(t *testing.T) {
		V = ""
		got := Get()
		// Should return a semver-compliant version
		// Either "X.0.0-dev" or "X.0.0-dev+sha" format
		if !strings.Contains(got, ".0.0-dev") {
			t.Errorf("Get() = %q, want semver dev format", got)
		}
	})
}

func TestPrint(t *testing.T) {
	// Save original V and restore after test
	origV := V
	defer func() { V = origV }()

	V = "1.2.3"

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Print("TEST")

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	got := buf.String()

	want := "TEST version 1.2.3\n"
	if got != want {
		t.Errorf("Print() output = %q, want %q", got, want)
	}
}

func TestVersionRegex(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"github.com/foo/bar/v1", "1"},
		{"github.com/foo/bar/v2", "2"},
		{"github.com/foo/bar/v10", "10"},
		{"github.com/foo/bar/v1/internal", "1"},
		{"github.com/foo/bar", ""},
		{"github.com/foo/bar/pkg", ""},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			matches := versionRegex.FindStringSubmatch(tt.path)
			got := ""
			if len(matches) > 1 {
				got = matches[1]
			}
			if got != tt.want {
				t.Errorf("versionRegex on %q = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}
