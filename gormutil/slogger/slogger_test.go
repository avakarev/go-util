package slogger_test

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/avakarev/go-util/gormutil/slogger"
	"github.com/avakarev/go-util/testutil"
	"gorm.io/gorm/logger"
)

// compile-time assertion that New satisfies gorm's logger contract
var _ logger.Interface = slogger.New(nil, logger.Info)

func newTestLogger(lvl slog.Level) (*slog.Logger, *bytes.Buffer) {
	var buf bytes.Buffer
	h := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: lvl})
	return slog.New(h), &buf
}

func TestTraceEmitsSQL(t *testing.T) {
	l, buf := newTestLogger(slog.LevelInfo)
	gl := slogger.New(l, logger.Info)

	begin := time.Now().Add(-5 * time.Millisecond)
	gl.Trace(context.Background(), begin, func() (string, int64) {
		return "SELECT 1", 1
	}, nil)

	out := buf.String()
	testutil.Diff(true, strings.Contains(out, "SELECT 1"), t)
	testutil.Diff(true, strings.Contains(out, `"rows":1`), t)
}

func TestTraceIgnoresRecordNotFound(t *testing.T) {
	l, buf := newTestLogger(slog.LevelInfo)
	gl := slogger.New(l, logger.Info)

	gl.Trace(context.Background(), time.Now(), func() (string, int64) {
		return "SELECT 1", 0
	}, logger.ErrRecordNotFound)

	// record-not-found is not surfaced as an error line
	testutil.Diff(false, strings.Contains(buf.String(), `"level":"ERROR"`), t)
}

func TestTraceSurfacesError(t *testing.T) {
	l, buf := newTestLogger(slog.LevelInfo)
	gl := slogger.New(l, logger.Info)

	gl.Trace(context.Background(), time.Now(), func() (string, int64) {
		return "SELECT 1", 0
	}, errors.New("boom"))

	out := buf.String()
	testutil.Diff(true, strings.Contains(out, `"level":"ERROR"`), t)
	testutil.Diff(true, strings.Contains(out, "boom"), t)
}

func TestNewNilFallsBackToDefault(t *testing.T) {
	// New(nil, ...) must not panic and must satisfy the interface
	gl := slogger.New(nil, logger.Warn)
	testutil.Diff(true, gl != nil, t)
}

func TestNewDefaultLevelTracksSlog(t *testing.T) {
	orig := slog.Default()
	defer slog.SetDefault(orig)

	// debug enabled -> gorm Info level (info logs are emitted)
	dbg, buf := newTestLogger(slog.LevelDebug)
	slog.SetDefault(dbg)
	slogger.NewDefault().Info(context.Background(), "hello")
	testutil.Diff(true, strings.Contains(buf.String(), "hello"), t)

	// debug disabled -> gorm Error level (info logs are suppressed)
	warn, buf2 := newTestLogger(slog.LevelWarn)
	slog.SetDefault(warn)
	slogger.NewDefault().Info(context.Background(), "hello")
	testutil.Diff(false, strings.Contains(buf2.String(), "hello"), t)
}
