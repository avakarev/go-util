package slog_test

import (
	"log/slog"
	"testing"

	testslog "github.com/avakarev/go-util/testutil/slog"
)

func TestCapture(t *testing.T) {
	buf := testslog.Capture(t, slog.LevelInfo)

	slog.Info("hello", "key", "val")

	out := buf.String()
	if want := "level=INFO"; buf.Count(want) != 1 {
		t.Errorf("want %q once, got output:\n%s", want, out)
	}
	if want := "msg=hello"; buf.Count(want) != 1 {
		t.Errorf("want %q once, got output:\n%s", want, out)
	}
	if want := "key=val"; buf.Count(want) != 1 {
		t.Errorf("want %q once, got output:\n%s", want, out)
	}
}

func TestCaptureNoTimestamp(t *testing.T) {
	buf := testslog.Capture(t, slog.LevelInfo)

	slog.Info("hello")

	if got := buf.Count("time="); got != 0 {
		t.Errorf("want no timestamp, got %d in output:\n%s", got, buf.String())
	}
}

func TestCaptureLevelFilter(t *testing.T) {
	buf := testslog.Capture(t, slog.LevelWarn)

	slog.Info("skipped")
	slog.Warn("kept")

	if got := buf.Count("skipped"); got != 0 {
		t.Errorf("want debug/info filtered out, got %d", got)
	}
	if got := buf.Count("kept"); got != 1 {
		t.Errorf("want warn kept once, got %d", got)
	}
}

func TestBufferCount(t *testing.T) {
	buf := testslog.Capture(t, slog.LevelInfo)

	slog.Info("ping")
	slog.Info("ping")

	if got := buf.Count("msg=ping"); got != 2 {
		t.Errorf("want 2, got %d", got)
	}
}
