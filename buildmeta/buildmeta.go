// Package buildmeta holds build and runtime information
package buildmeta

import (
	"runtime"
	"time"

	"github.com/avakarev/go-util/timeutil"
)

var (
	// Commit is a git commit sha
	Commit string

	// Ref is a git branch or tag ref name
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

// Meta defines build meta
type Meta struct {
	// Compiler is Go compiler version
	Compiler string `json:"compiler"`

	// OS is a running program's operating system target
	OS string `json:"os"`

	// Arch is a running program's architecture target
	Arch string `json:"arch"`

	// Commit is a git commit sha
	Commit string `json:"commit"`

	// Ref is a git branch or tag ref name
	Ref string `json:"ref"`

	// BuildTimeUTC is a build datetime in UTC
	BuildTimeUTC string `json:"buildTimeUTC"`

	// BuildTimeLocal is a build datetime in local timezone
	BuildTimeLocal string `json:"buildTimeLocal"`

	// ServerUptime return server's uptime
	ServerUptime string `json:"serverUptime"`

	// ServerTime returns server's local time
	ServerTime string `json:"serverTime"`

	// Timezone returns server's local timezone
	ServerTimezone string `json:"serverTimezone"`
}

// New returns new Meta value
func New() Meta {
	return Meta{
		Compiler:       runtime.Version(),
		OS:             runtime.GOOS,
		Arch:           runtime.GOARCH,
		Commit:         Commit,
		Ref:            Ref,
		BuildTimeUTC:   BuildTimeUTC,
		BuildTimeLocal: buildTimeLocal,
		ServerUptime:   uptime,
		ServerTime:     timeutil.Local(time.Now()).Format(time.RFC3339),
		ServerTimezone: timeutil.Location.String(),
	}
}

func init() {
	if t, err := time.Parse(time.RFC3339, BuildTimeUTC); err == nil && !t.IsZero() {
		buildTimeLocal = timeutil.Local(t).Format(time.RFC3339)
	}
	uptime = timeutil.Local(time.Now()).Format(time.RFC3339)
}
