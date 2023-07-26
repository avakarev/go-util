// Package zerologutil initializes and sets defaults for zerolog
package zerologutil

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// DefaultLevel is default log level
const DefaultLevel = "info"

// SetLevel sets `zerolog`'s global level
func SetLevel(name string) error {
	level, err := zerolog.ParseLevel(name)
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(level)
	log.Debug().Msgf("log level is set to %q", level)
	return nil
}

// IsDebug checks whether log level is `debug`
func IsDebug() bool {
	return zerolog.GlobalLevel() == zerolog.DebugLevel
}

// Init initializes logger
func Init() error {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = DefaultLevel
	}

	zerolog.MessageFieldName = "msg"

	return SetLevel(level)
}

// MustInit is like Init but panics in case of error
func MustInit() {
	if err := Init(); err != nil {
		panic("log.MustInit(): " + err.Error())
	}
}
