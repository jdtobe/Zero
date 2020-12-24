package main

import (
	"os"

	"github.com/spf13/cast"
)

func getEnvDefault(name, def string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}

	return def
}

func getEnvDefaultFloat64(name string, def float64) float64 {
	if v := cast.ToFloat64(os.Getenv(name)); v != 0.0 {
		return v
	}

	return def
}

func getEnvDefaultInt(name string, def int) int {
	if v := cast.ToInt(os.Getenv(name)); v != 0 {
		return v
	}

	return def
}
