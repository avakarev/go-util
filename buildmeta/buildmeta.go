// Package buildmeta holds build and runtime information
package buildmeta

import (
	"runtime"
	"time"

	"github.com/avakarev/go-util/timeutil"
)

var (
	// Commit is git commit sha
	Commit string

	// Ref is git branch or tag ref name
	Ref string

	// BuildTimeUTC is a build datetime in UTC
	BuildTimeUTC string
)

var (
	// buildTimeLocal is a build datetime in local timezone
	buildTimeLocal string

	// uptime is the application's uptime
	uptime string
)

// Compiler returns Go compiler version
func Compiler() string {
	return runtime.Version()
}

// OS returns running program's operating system target
func OS() string {
	return runtime.GOOS
}

// Arch returns running program's architecture target
func Arch() string {
	return runtime.GOARCH
}

// Uptime return server's uptime
func Uptime() string {
	return uptime
}

// Time returns server's local time
func Time() string {
	return timeutil.Local(time.Now()).Format(time.RFC3339)
}

// Timezone returns server's local timezone
func Timezone() string {
	return timeutil.Location.String()
}

// Fields returns build meta as map
func Fields() map[string]any {
	return map[string]any{
		"compiler": Compiler(),
		"os":       OS(),
		"arch":     Arch(),

		"commit":         Commit,
		"ref":            Ref,
		"buildTimeUTC":   BuildTimeUTC,
		"buildTimeLocal": buildTimeLocal,

		"serverUptime":   uptime,
		"serverTime":     Time(),
		"serverTimezone": Timezone(),
	}
}

func init() {
	if t, err := time.Parse(time.RFC3339, BuildTimeUTC); err == nil && !t.IsZero() {
		buildTimeLocal = timeutil.Local(t).Format(time.RFC3339)
	}
	uptime = timeutil.Local(time.Now()).Format(time.RFC3339)
}
