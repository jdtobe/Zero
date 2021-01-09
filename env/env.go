package env

import (
	"os"

	"github.com/spf13/cast"
)

// GetDefault returns an environment variable by name as a string, if found, and def if not.
func GetDefault(name, def string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}

	return def
}

// GetDefaultFloat64 returns an environment variable by name as a float64, if found, and def if not.
func GetDefaultFloat64(name string, def float64) float64 {
	if v := cast.ToFloat64(os.Getenv(name)); v != 0.0 {
		return v
	}

	return def
}

// GetDefaultInt returns an environment variable by name as an int, if found, and def if not.
func GetDefaultInt(name string, def int) int {
	if v := cast.ToInt(os.Getenv(name)); v != 0 {
		return v
	}

	return def
}
