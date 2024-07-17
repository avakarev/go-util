// Package envutil implements environment variable helpers
package envutil

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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

// StrSlice returns environment variable split by given separator
func StrSlice(name string, sep string) []string {
	s := Str(name)
	if s == "" {
		return nil
	}
	return strings.Split(s, sep)
}

// MustStrSlice returns environment variable split by given separator and errors if it's not set
func MustStrSlice(name string, sep string) ([]string, error) {
	s, err := MustStr(name)
	if err != nil {
		return nil, err
	}
	return strings.Split(s, sep), nil
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
		return 0, err
	}
	return strconv.ParseInt(s, 10, 64)
}
