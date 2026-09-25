// Package slogger implements a gorm logger adapter backed by log/slog
package slogger

import (
	"context"
	"log/slog"
	"time"

	"gorm.io/gorm/logger"
)

// New returns a gorm logger that writes through the given *slog.Logger
// A nil logger falls back to slog.Default()
func New(l *slog.Logger, lvl logger.LogLevel) logger.Interface {
	if l == nil {
		l = slog.Default()
	}
	return logger.NewSlogLogger(l, logger.Config{
		SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold
		LogLevel:                  lvl,                    // Log level
		IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound error for logger
	})
}

// NewDefault returns a gorm logger backed by slog.Default()
// The gorm level tracks slog.Default(): Info when debug logging is enabled,
// otherwise Error, so it stays in sync with slogutil.Configure
func NewDefault() logger.Interface {
	lvl := logger.Error
	if slog.Default().Enabled(context.Background(), slog.LevelDebug) {
		lvl = logger.Info
	}
	return New(slog.Default(), lvl)
}
