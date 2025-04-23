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

// ShouldStr returns environment variable and error if it's not set
func ShouldStr(name string) (string, error) {
	value, ok := os.LookupEnv(name)
	if !ok {
		return "", fmt.Errorf("environment variable %q is required, but not set", name)
	}
	if value == "" {
		return "", fmt.Errorf("environment variable %q is required, but empty", name)
	}
	return value, nil
}

// MustStr is like ShouldStr but panics in case of error
func MustStr(name string) string {
	value, err := ShouldStr(name)
	if err != nil {
		panic(err)
	}
	return value
}

// ShouldStrSlice returns environment variable split by given separator and error if it's not set
func ShouldStrSlice(name string, sep string) ([]string, error) {
	str, err := ShouldStr(name)
	if err != nil {
		return nil, err
	}
	return strings.Split(str, sep), nil
}

// MustStrSlice is like ShouldStrSlice but panics in case of error
func MustStrSlice(name string, sep string) []string {
	value, err := ShouldStrSlice(name, sep)
	if err != nil {
		panic(err)
	}
	return value
}

// ShouldInt returns environment variable as int64 and conversion error if any
func ShouldInt(name string) (int64, error) {
	str := Str(name)
	if str == "" {
		return 0, nil
	}
	return strconv.ParseInt(str, 10, 64)
}

// MustInt is like ShouldInt but panics in case of error
func MustInt(name string) int64 {
	i, err := ShouldInt(name)
	if err != nil {
		panic(err)
	}
	return i
}
