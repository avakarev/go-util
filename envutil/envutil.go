// Package envutil implements environment variable helpers
package envutil

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Str returns environment variable as string
func Str(key string) string {
	return os.Getenv(key)
}

// StrDef returns environment variable as string, or def if the variable is unset or empty
func StrDef(key string, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

// ShouldStr returns env variable and error if it's not set
func ShouldStr(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("env variable %q is required, but not set", key)
	}
	if value == "" {
		return "", fmt.Errorf("env variable %q is required, but empty", key)
	}
	return value, nil
}

// MustStr is like ShouldStr but panics in case of error
func MustStr(key string) string {
	value, err := ShouldStr(key)
	if err != nil {
		panic(err)
	}
	return value
}

// ShouldStrSlice returns env variable split by given separator and error if it's not set
func ShouldStrSlice(key string, sep string) ([]string, error) {
	str, err := ShouldStr(key)
	if err != nil {
		return nil, err
	}
	return strings.Split(str, sep), nil
}

// MustStrSlice is like ShouldStrSlice but panics in case of error
func MustStrSlice(key string, sep string) []string {
	value, err := ShouldStrSlice(key, sep)
	if err != nil {
		panic(err)
	}
	return value
}

// ShouldInt returns env variable parsed as int64, or a conversion error if any.
// An unset or empty variable yields 0 without error.
func ShouldInt(key string) (int64, error) {
	str := Str(key)
	if str == "" {
		return 0, nil
	}
	return strconv.ParseInt(str, 10, 64)
}

// MustInt is like ShouldInt but panics in case of error
func MustInt(key string) int64 {
	i, err := ShouldInt(key)
	if err != nil {
		panic(err)
	}
	return i
}

// DurDef returns environment variable parsed as time.Duration, or def if the
// variable is unset or empty. It returns an error if the value cannot be parsed.
func DurDef(key string, def time.Duration) (time.Duration, error) {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return def, nil
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, fmt.Errorf("env variable %q is invalid duration: %w", key, err)
	}
	return d, nil
}
