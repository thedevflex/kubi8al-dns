package main

import (
	"os"
	"strings"
	"time"
)

type Config struct {
	Port              string
	BaseDomain        string
	DefaultEnv        string
	AllowedNamespaces []string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

func LoadConfig() *Config {
	config := &Config{
		Port:         getEnv("PORT", "8080"),
		BaseDomain:   getEnv("BASE_DOMAIN", "code-craft.co.in"),
		DefaultEnv:   getEnv("DEFAULT_ENV", "true"),
		ReadTimeout:  getDurationEnv("READ_TIMEOUT", 15*time.Second),
		WriteTimeout: getDurationEnv("WRITE_TIMEOUT", 15*time.Second),
		IdleTimeout:  getDurationEnv("IDLE_TIMEOUT", 60*time.Second),
	}

	if allowedNS := os.Getenv("ALLOWED_NAMESPACES"); allowedNS != "" {
		config.AllowedNamespaces = strings.Split(allowedNS, ",")
		for i, ns := range config.AllowedNamespaces {
			config.AllowedNamespaces[i] = strings.TrimSpace(ns)
		}
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
