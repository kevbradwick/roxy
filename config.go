package roxy

import (
	"fmt"
	"os"
)

var config Config

type Config struct {
	Port            string
	Host            string
	Target          string
	HealthCheckPath string
}

func envDefault(name, defaultValue string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}

	return defaultValue
}

func envRequired(name string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}

	fmt.Printf("ERROR: Missing environment variable '%s'", name)
	os.Exit(1)
	return ""
}

func init() {
	config = Config{}
	config.Port = envDefault("PORT", "9000")
	config.Host = envDefault("HOST", "0.0.0.0")
	config.Target = envRequired("TARGET")
	config.HealthCheckPath = envDefault("HEALTHCHECK_PATH", "/__health__")
}

func EnvConfig() *Config {
	return &config
}
