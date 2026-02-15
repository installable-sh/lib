package version

import (
	"fmt"
	"regexp"
	"runtime/debug"
)

// V can be set at build time via:
//
//	go build -ldflags="-X github.com/installable-sh/lib/version.V=1.0.0"
var V = ""

// versionRegex extracts major version from module path (e.g., /v1 -> 1)
var versionRegex = regexp.MustCompile(`/v(\d+)(?:/|$)`)

// Get returns a semver-compliant version string.
// If V is set via ldflags, it returns that.
// Otherwise, it infers the major version from the module path and constructs
// a dev version: 1.0.0-dev+sha.abcd123
func Get() string {
	if V != "" {
		return V
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "0.0.0-dev"
	}

	// Extract major version from module path
	major := "0"
	if matches := versionRegex.FindStringSubmatch(info.Main.Path); len(matches) > 1 {
		major = matches[1]
	}

	var revision, modified string
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			if len(setting.Value) >= 7 {
				revision = setting.Value[:7]
			} else {
				revision = setting.Value
			}
		case "vcs.modified":
			if setting.Value == "true" {
				modified = ".dirty"
			}
		}
	}

	if revision == "" {
		return fmt.Sprintf("%s.0.0-dev", major)
	}

	return fmt.Sprintf("%s.0.0-dev+%s%s", major, revision, modified)
}

// Print outputs version information for a command.
func Print(command string) {
	fmt.Printf("%s version %s\n", command, Get())
}
