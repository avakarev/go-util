package testutil

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

var homeDir string

func caller() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}

	if homeDir != "" {
		file = strings.Replace(file, homeDir, "~", 1)
	}

	return fmt.Sprintf("\nFailed at %s:%d\n", file, line)
}

func init() {
	homeDir = os.Getenv("HOME")
}
