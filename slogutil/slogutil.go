// Package slogutil configures a global slog logger backed by phuslu/log.
package slogutil

import (
	"log/slog"
	"os"
	"time"

	"github.com/phuslu/log"
)

// Option configures Configure's behavior.
type Option func(*config)

// WithLevel sets the log level, overriding the `LOG_LEVEL` env var
// name may be one of: debug, info, warn, error
func WithLevel(name string) Option {
	return func(c *config) {
		c.level = name
	}
}

// WithJSON sets structured JSON output when true, or colorized console output
// when false, JSON by default. In console mode, color is enabled only on a terminal
func WithJSON(json bool) Option {
	return func(c *config) {
		c.json = json
	}
}

// WithTime sets the time field name, format and location
// An empty field or fmt falls back to its default ("time" and time.RFC3339),
// a nil loc falls back to the `TZ` env var, then time.Local. fmt is a Go time
// layout, or one of phuslu/log's special TimeFormatUnix* values for a UNIX timestamp
func WithTime(field string, fmt string, loc *time.Location) Option {
	return func(c *config) {
		c.timeField = field
		c.timeFmt = fmt
		c.timeLoc = loc
	}
}

// WithMsgField sets the message field name, applied in JSON mode only
// In JSON mode this mutates a phuslu/log package-global shared by all loggers
// in the process, console mode renders the message positionally and ignores it
func WithMsgField(name string) Option {
	return func(c *config) {
		c.msgField = name
	}
}

// WithCaller enables or disables the caller, callerfunc and goid fields
// phuslu/log emits all three together, they cannot be toggled individually
func WithCaller(caller bool) Option {
	return func(c *config) {
		c.caller = caller
	}
}

// Configure sets the global slog default logger, backed by phuslu/log.
func Configure(opts ...Option) error {
	// json defaults to true and is seeded here, before options run, so that an
	// explicit WithJSON(false) is distinguishable from the zero value
	cfg := config{json: true}
	for _, opt := range opts {
		opt(&cfg)
	}
	if err := cfg.defaults(); err != nil {
		return err
	}
	if err := cfg.validate(); err != nil {
		return err
	}

	level, err := cfg.parseLevel()
	if err != nil {
		return err
	}

	if cfg.json {
		// log.MessageKey is a phuslu/log package-global shared by all loggers in
		// the process. It only affects JSON output (console renders the message
		// positionally), so it is set only in JSON mode. It defaults to "message",
		// here it follows the configured field.
		log.MessageKey = cfg.msgField
	}

	var caller int
	if cfg.caller {
		caller = 1
	}

	logger := &log.Logger{
		Level:        level,
		TimeField:    cfg.timeField,
		TimeFormat:   cfg.timeFmt,
		TimeLocation: cfg.timeLoc,
		Caller:       caller,
	}
	if !cfg.json {
		// Color only when stderr, phuslu/log's default output, is a terminal.
		logger.Writer = &log.ConsoleWriter{ColorOutput: log.IsTerminal(os.Stderr.Fd())}
	}

	slog.SetDefault(logger.Slog())
	return nil
}

// MustConfigure is like Configure but panics in case of error.
func MustConfigure(opts ...Option) {
	if err := Configure(opts...); err != nil {
		panic(err)
	}
}
