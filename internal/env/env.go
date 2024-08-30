package env

import "os"

// The `MustGet` function retrieves the value of an environment variable or panics if it is missing.
func MustGet(key string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	panic("missing required environment variable: " + key)
}
