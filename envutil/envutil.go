// Package envutil implements environment variable helpers
package envutil

import (
	"fmt"
	"os"
	"strconv"
)

// Str returns environment variable as string
func Str(name string) string {
	return os.Getenv(name)
}

// MustStr returns environment variable and error if it's not set
func MustStr(name string) (string, error) {
	value, ok := os.LookupEnv(name)
	if !ok {
		return "", fmt.Errorf("environment variable %q is not set", name)
	}
	if value == "" {
		return "", fmt.Errorf("environment variable %q is empty", name)
	}
	return value, nil
}

// Int returns environment variable as int64 and conversion error if any
func Int(name string) (int64, error) {
	s := Str(name)
	if s == "" {
		return 0, nil
	}
	return strconv.ParseInt(s, 10, 64)
}

// MustInt returns environment variable and error if it's not set or conversion error if any
func MustInt(name string) (int64, error) {
	s, err := MustStr(name)
	if err != nil {
		return 0, nil
	}
	return strconv.ParseInt(s, 10, 64)
}
