// Package zerologger implements gorm zerolog adapter
package zerologger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm/logger"
)

// New returns new logger value
func New(lvl logger.LogLevel) logger.Interface {
	zlog := log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC822}).With().Caller().Logger()
	return logger.New(
		&zlog, // io writer
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold
			LogLevel:                  lvl,                    // Log level
			IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,                  // Disable color
		},
	)
}

// NewDefault returns new logger value with default level
func NewDefault() logger.Interface {
	if lvl := os.Getenv("LOG_LEVEL"); lvl == zerolog.DebugLevel.String() {
		return New(logger.Info)
	}
	return New(logger.Error)
}
