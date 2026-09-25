package slogutil_test

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/avakarev/go-util/slogutil"
	"github.com/avakarev/go-util/testutil"
)

func TestConfigure(t *testing.T) {
	// unset LOG_LEVEL uses the default and enables info
	testutil.MustNoErr(slogutil.Configure(), t)
	testutil.Diff(true, slog.Default().Enabled(t.Context(), slog.LevelInfo), t)

	// each supported level configures without error
	for _, level := range []string{
		slogutil.LevelTrace,
		slogutil.LevelDebug,
		slogutil.LevelInfo,
		slogutil.LevelWarn,
		slogutil.LevelError,
		slogutil.LevelFatal,
		slogutil.LevelPanic,
	} {
		resetEnv := testutil.SetEnv(testutil.Env{"LOG_LEVEL": level})
		testutil.MustNoErr(slogutil.Configure(), t)
		resetEnv()
	}

	// debug level enables debug logging
	resetEnv := testutil.SetEnv(testutil.Env{"LOG_LEVEL": slogutil.LevelDebug})
	testutil.MustNoErr(slogutil.Configure(), t)
	testutil.Diff(true, slog.Default().Enabled(t.Context(), slog.LevelDebug), t)
	resetEnv()

	// error level disables info logging
	resetEnv = testutil.SetEnv(testutil.Env{"LOG_LEVEL": slogutil.LevelError})
	testutil.MustNoErr(slogutil.Configure(), t)
	testutil.Diff(false, slog.Default().Enabled(t.Context(), slog.LevelInfo), t)
	resetEnv()

	// invalid LOG_LEVEL yields an error
	resetEnv = testutil.SetEnv(testutil.Env{"LOG_LEVEL": "nope"})
	defer resetEnv()
	testutil.MustErr(
		errors.New("slogutil: invalid LOG_LEVEL \"nope\", want one of: trace, debug, info, warn, error, fatal, panic"),
		slogutil.Configure(),
		t,
	)
}

func TestConfigureWithLevel(t *testing.T) {
	// WithLevel overrides an unset env var
	testutil.MustNoErr(slogutil.Configure(slogutil.WithLevel(slogutil.LevelError)), t)
	testutil.Diff(false, slog.Default().Enabled(t.Context(), slog.LevelInfo), t)

	// WithLevel overrides the LOG_LEVEL env var
	resetEnv := testutil.SetEnv(testutil.Env{"LOG_LEVEL": slogutil.LevelError})
	testutil.MustNoErr(slogutil.Configure(slogutil.WithLevel(slogutil.LevelDebug)), t)
	testutil.Diff(true, slog.Default().Enabled(t.Context(), slog.LevelDebug), t)
	resetEnv()

	// invalid WithLevel value yields an error
	testutil.MustErr(
		errors.New("slogutil: invalid LOG_LEVEL \"nope\", want one of: trace, debug, info, warn, error, fatal, panic"),
		slogutil.Configure(slogutil.WithLevel("nope")),
		t,
	)
}

func TestConfigureFormat(t *testing.T) {
	// JSON (default) configures without error
	testutil.MustNoErr(slogutil.Configure(), t)
	testutil.MustNoErr(slogutil.Configure(slogutil.WithJSON(true)), t)

	// console configures without error
	testutil.MustNoErr(slogutil.Configure(slogutil.WithJSON(false)), t)

	// combined with WithLevel
	testutil.MustNoErr(slogutil.Configure(slogutil.WithLevel(slogutil.LevelDebug), slogutil.WithJSON(false)), t)
	testutil.Diff(true, slog.Default().Enabled(t.Context(), slog.LevelDebug), t)
}

func TestConfigureWithTime(t *testing.T) {
	// default time field and format configure without error
	testutil.MustNoErr(slogutil.Configure(), t)

	// custom field and layout configure without error
	testutil.MustNoErr(slogutil.Configure(slogutil.WithTime("ts", time.RFC3339Nano, nil)), t)

	// empty field and fmt fall back to defaults
	testutil.MustNoErr(slogutil.Configure(slogutil.WithTime("", "", nil)), t)

	// combined with other options
	testutil.MustNoErr(slogutil.Configure(
		slogutil.WithLevel(slogutil.LevelDebug),
		slogutil.WithJSON(false),
		slogutil.WithTime("ts", time.RFC3339, nil),
	), t)
	testutil.Diff(true, slog.Default().Enabled(t.Context(), slog.LevelDebug), t)
}

func TestConfigureWithMsgField(t *testing.T) {
	// default message field configures without error
	testutil.MustNoErr(slogutil.Configure(), t)

	// custom message field configures without error
	testutil.MustNoErr(slogutil.Configure(slogutil.WithMsgField("message")), t)

	// combined with other options
	testutil.MustNoErr(slogutil.Configure(
		slogutil.WithLevel(slogutil.LevelDebug),
		slogutil.WithMsgField("message"),
	), t)
	testutil.Diff(true, slog.Default().Enabled(t.Context(), slog.LevelDebug), t)
}

func TestConfigureWithCaller(t *testing.T) {
	// caller disabled (default) configures without error
	testutil.MustNoErr(slogutil.Configure(), t)

	// caller enabled configures without error
	testutil.MustNoErr(slogutil.Configure(slogutil.WithCaller(true)), t)

	// combined with other options
	testutil.MustNoErr(slogutil.Configure(
		slogutil.WithLevel(slogutil.LevelDebug),
		slogutil.WithCaller(true),
	), t)
	testutil.Diff(true, slog.Default().Enabled(t.Context(), slog.LevelDebug), t)
}

// captureStderr configures the logger with the given options, runs fn, and
// returns everything written to os.Stderr (where phuslu/log writes by default).
func captureStderr(t *testing.T, fn func(), opts ...slogutil.Option) string {
	t.Helper()

	r, w, err := os.Pipe()
	testutil.MustNoErr(err, t)

	orig := os.Stderr
	os.Stderr = w
	defer func() { os.Stderr = orig }()

	testutil.MustNoErr(slogutil.Configure(opts...), t)
	fn()

	testutil.MustNoErr(w.Close(), t)
	out, err := io.ReadAll(r)
	testutil.MustNoErr(err, t)
	return string(out)
}

// timeValue parses the RFC3339 "time" field out of a single JSON log line.
func timeValue(t *testing.T, line string) time.Time {
	t.Helper()
	var entry struct {
		Time time.Time `json:"time"`
	}
	testutil.MustNoErr(json.Unmarshal([]byte(strings.TrimSpace(line)), &entry), t)
	return entry.Time
}

func TestConfigureOutput(t *testing.T) {
	// default output uses the "msg" and "time" field names
	line := captureStderr(t, func() { slog.Info("hello", "k", "v") })
	testutil.Diff(true, strings.Contains(line, `"msg":"hello"`), t)
	testutil.Diff(true, strings.Contains(line, `"time":`), t)
	testutil.Diff(true, strings.Contains(line, `"k":"v"`), t)
	// caller info is off by default
	testutil.Diff(false, strings.Contains(line, `"caller":`), t)
	testutil.Diff(false, strings.Contains(line, `"goid":`), t)

	// custom field names are honored in the output
	line = captureStderr(t, func() { slog.Info("hey") },
		slogutil.WithMsgField("message"),
		slogutil.WithTime("ts", time.RFC3339, nil),
	)
	testutil.Diff(true, strings.Contains(line, `"message":"hey"`), t)
	testutil.Diff(true, strings.Contains(line, `"ts":`), t)

	// enabling caller adds caller, callerfunc and goid
	line = captureStderr(t, func() { slog.Info("hi") }, slogutil.WithCaller(true))
	testutil.Diff(true, strings.Contains(line, `"caller":`), t)
	testutil.Diff(true, strings.Contains(line, `"callerfunc":`), t)
	testutil.Diff(true, strings.Contains(line, `"goid":`), t)

	// console output is not JSON but still carries the message. Color is auto-
	// detected from the terminal; under test os.Stderr is a pipe (not a TTY),
	// so no ANSI codes are emitted.
	line = captureStderr(t, func() { slog.Info("hello") }, slogutil.WithJSON(false))
	testutil.Diff(true, strings.Contains(line, "hello"), t)
	testutil.Diff(false, strings.Contains(line, `"msg":"hello"`), t)
	testutil.Diff(false, strings.Contains(line, "\x1b["), t)
}

func TestConfigureFieldCollision(t *testing.T) {
	// time field colliding with the message field
	testutil.MustErr(
		errors.New("slogutil: message field \"time\" collides with time field"),
		slogutil.Configure(slogutil.WithMsgField("time")),
		t,
	)

	// time field colliding with the reserved level field
	testutil.MustErr(
		errors.New("slogutil: level field \"level\" collides with time field"),
		slogutil.Configure(slogutil.WithTime("level", time.RFC3339, nil)),
		t,
	)

	// distinct field names configure without error
	testutil.MustNoErr(slogutil.Configure(
		slogutil.WithTime("ts", time.RFC3339, nil),
		slogutil.WithMsgField("message"),
	), t)
}

func TestConfigureWithTimeLocation(t *testing.T) {
	// UTC renders an RFC3339 timestamp with a "Z" zone suffix
	line := captureStderr(t, func() { slog.Info("hello") }, slogutil.WithTime("", "", time.UTC))
	testutil.Diff(true, strings.Contains(line, `"msg":"hello"`), t)
	testutil.Diff(true, timeValue(t, line).Location() == time.UTC, t)

	// a fixed non-UTC location is reflected in the timestamp's offset
	kolkata, err := time.LoadLocation("Asia/Kolkata") // +05:30, has no DST
	testutil.MustNoErr(err, t)
	line = captureStderr(t, func() { slog.Info("hello") }, slogutil.WithTime("", "", kolkata))
	_, offset := timeValue(t, line).Zone()
	testutil.Diff(int((5*time.Hour+30*time.Minute)/time.Second), offset, t)
}

func TestConfigureTZEnv(t *testing.T) {
	// a valid TZ is respected when no location is passed
	resetEnv := testutil.SetEnv(testutil.Env{"TZ": "Asia/Kolkata"})
	line := captureStderr(t, func() { slog.Info("hello") })
	_, offset := timeValue(t, line).Zone()
	testutil.Diff(int((5*time.Hour+30*time.Minute)/time.Second), offset, t)
	resetEnv()

	// an explicit location overrides the TZ env var
	resetEnv = testutil.SetEnv(testutil.Env{"TZ": "Asia/Kolkata"})
	line = captureStderr(t, func() { slog.Info("hello") }, slogutil.WithTime("", "", time.UTC))
	testutil.Diff(true, timeValue(t, line).Location() == time.UTC, t)
	resetEnv()

	// an invalid TZ yields an error
	resetEnv = testutil.SetEnv(testutil.Env{"TZ": "not/a/zone"})
	defer resetEnv()
	_, err := time.LoadLocation("not/a/zone")
	testutil.MustErr(
		errors.New("slogutil: invalid TZ \"not/a/zone\": "+err.Error()),
		slogutil.Configure(),
		t,
	)
}

func TestMustConfigure(t *testing.T) {
	// valid level does not panic
	resetEnv := testutil.SetEnv(testutil.Env{"LOG_LEVEL": slogutil.LevelInfo})
	slogutil.MustConfigure()
	resetEnv()

	// invalid level panics
	resetEnv = testutil.SetEnv(testutil.Env{"LOG_LEVEL": "nope"})
	defer resetEnv()
	func() {
		defer func() {
			testutil.Diff(true, recover() != nil, t)
		}()
		slogutil.MustConfigure()
	}()
}
