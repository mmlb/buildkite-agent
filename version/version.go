// Package version provides the agent version strings.
package version

import (
	_ "embed"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
)

// You can overridden buildVersion at compile time by using:
//
//  go run -ldflags "-X github.com/buildkite/agent/v3/agent.buildVersion=abc" . --version
//
// On CI, the binaries are always build with the buildVersion variable set.
//
// Pre-release builds' versions must be in the format `x.y-beta`, `x.y-beta.z` or `x.y-beta.z.a`

//go:embed VERSION
var baseVersion string
var buildNumber string

func Version() string {
	return strings.TrimSpace(baseVersion)
}

func BuildNumber() string {
	if buildNumber == "" {
		return "x"
	}
	return buildNumber
}

func commitInfo() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "x"
	}

	dirty := ".dirty"
	var commit string
	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			commit = setting.Value
		} else if setting.Key == "vcs.modified" && setting.Value == "false" {
			dirty = ""
		}
	}

	return commit + dirty
}

// FullVersion is a SemVer 2.0 compliant version string that includes
// [build metadata](https://semver.org/#spec-item-10) consisting of the build
// number (if any), the commit hash, and whether the build was dirty.
func FullVersion() string {
	return fmt.Sprintf("%s+%s.%s", Version(), BuildNumber(), commitInfo())
}

func UserAgent() string {
	return fmt.Sprintf(
		"buildkite-agent/%s.%s (%s; %s)",
		Version(),
		BuildNumber(),
		runtime.GOOS,
		runtime.GOARCH,
	)
}
