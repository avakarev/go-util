package gormutil

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm/logger"
)

// Logger is a zerolog-backed gorm logger
var Logger logger.Interface

// NewLogger returns new logger value
func NewLogger(lvl logger.LogLevel) logger.Interface {
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

func init() {
	if lvl := os.Getenv("LOG_LEVEL"); lvl == zerolog.DebugLevel.String() {
		Logger = NewLogger(logger.Info)
	} else {
		Logger = NewLogger(logger.Error)
	}
}
