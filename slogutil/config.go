package slogutil

import (
	"fmt"
	"time"

	"github.com/avakarev/go-util/envutil"
	"github.com/phuslu/log"
)

const (
	// LevelTrace is the "trace" log level
	LevelTrace = "trace"
	// LevelDebug is the "debug" log level
	LevelDebug = "debug"
	// LevelInfo is the "info" log level
	LevelInfo = "info"
	// LevelWarn is the "warn" log level
	LevelWarn = "warn"
	// LevelError is the "error" log level
	LevelError = "error"
	// LevelFatal is the "fatal" log level
	LevelFatal = "fatal"
	// LevelPanic is the "panic" log level
	LevelPanic = "panic"
)

type config struct {
	// level is the log level, it falls back to the `LOG_LEVEL` env var, then LevelInfo
	level string
	// timeField is the time field name, "time" by default
	timeField string
	// timeFmt is the Go time layout for the time field, time.RFC3339 by default
	// See github.com/phuslu/log for the special TimeFormatUnix* values, which
	// produce a UNIX timestamp instead
	timeFmt string
	// timeLoc is the time location for the time field, it falls back to the `TZ`
	// env var when set, otherwise nil so phuslu/log uses time.Local
	timeLoc *time.Location
	// msgField is the message field name, "msg" by default
	msgField string
	// caller enables the caller, callerfunc and goid fields, disabled by default.
	// caller is the "file:line" of the log call site, callerfunc is the enclosing
	// function name, and goid is the id of the goroutine that emitted the entry.
	// phuslu/log emits all three together; they cannot be toggled individually.
	caller bool
	// json sets the output format to structured JSON, true by default. When
	// false, colorized console output is used instead (color only on a terminal)
	json bool
}

// defaults fills unset fields with their default values
// It returns an error if `TZ` is set but not a loadable location.
func (cfg *config) defaults() error {
	if cfg.level == "" {
		cfg.level = envutil.StrDef("LOG_LEVEL", LevelInfo)
	}
	if cfg.timeField == "" {
		cfg.timeField = "time"
	}
	if cfg.timeFmt == "" {
		cfg.timeFmt = time.RFC3339
	}
	if cfg.timeLoc == nil {
		if tz := envutil.Str("TZ"); tz != "" {
			loc, err := time.LoadLocation(tz)
			if err != nil {
				return fmt.Errorf("slogutil: invalid TZ %q: %w", tz, err)
			}
			cfg.timeLoc = loc
		}
	}
	if cfg.msgField == "" {
		cfg.msgField = "msg"
	}
	return nil
}

// validate ensures the configured field names do not collide with each other
// or with the reserved level field, which would emit duplicate keys.
func (cfg *config) validate() error {
	fields := []struct {
		kind string
		name string
	}{
		{"time", cfg.timeField},
		{"message", cfg.msgField},
		{"level", log.LevelKey},
	}
	seen := make(map[string]string, len(fields))
	for _, f := range fields {
		if other, ok := seen[f.name]; ok {
			return fmt.Errorf("slogutil: %s field %q collides with %s field", f.kind, f.name, other)
		}
		seen[f.name] = f.kind
	}
	return nil
}

// parseLevel resolves a log level string into a phuslu/log level type
// unsupported names return an error
func (cfg *config) parseLevel() (log.Level, error) {
	switch cfg.level {
	case LevelTrace:
		return log.TraceLevel, nil
	case LevelDebug:
		return log.DebugLevel, nil
	case LevelInfo:
		return log.InfoLevel, nil
	case LevelWarn:
		return log.WarnLevel, nil
	case LevelError:
		return log.ErrorLevel, nil
	case LevelFatal:
		return log.FatalLevel, nil
	case LevelPanic:
		return log.PanicLevel, nil
	default:
		return 0, fmt.Errorf("slogutil: invalid LOG_LEVEL %q, want one of: trace, debug, info, warn, error, fatal, panic", cfg.level)
	}
}
