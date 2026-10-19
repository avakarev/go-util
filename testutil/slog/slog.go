// Package slog captures slog output
package slog

import (
	"bytes"
	"log/slog"
	"strings"
	"sync"
	"testing"
)

// Buffer holds captured log output. It is safe for concurrent writes (from
// paho, gocron or worker goroutines) while the test reads it.
type Buffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

// Write implements io.Writer.
func (b *Buffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

// String returns everything captured so far.
func (b *Buffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

// Count returns how many times s occurs in the captured output.
func (b *Buffer) Count(s string) int {
	return strings.Count(b.String(), s)
}

// Capture routes the default slog logger to a text handler at level, without
// timestamps, until the test ends, and returns what it writes. The default
// logger is process-wide, so a test that calls Capture must not run in
// parallel.
func Capture(t testing.TB, level slog.Level) *Buffer {
	t.Helper()
	buf := &Buffer{}
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	})))
	t.Cleanup(func() { slog.SetDefault(prev) })
	return buf
}
