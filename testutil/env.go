package testutil

import "os"

// Env is a helper type to declare map of environment variables key/value pairs
type Env map[string]string

// SetEnv set given Env value entries to the process's environment and returns reset function
func SetEnv(env Env) (reset func()) {
	originalEnv := Env{}

	for key, val := range env {
		if origVal, ok := os.LookupEnv(key); ok {
			originalEnv[key] = origVal
		}
		_ = os.Setenv(key, val)
	}

	return func() {
		for key := range env {
			origVal, ok := originalEnv[key]
			if ok {
				_ = os.Setenv(key, origVal)
				continue
			}
			_ = os.Unsetenv(key)
		}
	}
}
