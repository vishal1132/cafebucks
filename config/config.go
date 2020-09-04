package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
)

// Environment is the current runtime environment.
type Environment string

const (
	// Development is for when it's the development environment
	Development Environment = "development"

	// Testing is WISOTT
	Testing Environment = "testing"

	// Staging is WISOTT
	Staging Environment = "staging"

	// Production is WISOTT
	Production Environment = "production"
)

func strToEnv(s string) Environment {
	switch strings.ToLower(s) {
	case "production":
		return Production
	case "staging":
		return Staging
	case "testing":
		return Testing
	default:
		return Development
	}
}

// C is the config struct
type C struct {
	// Loglevel any valid zerolog logging level
	// ENV: LOG_LEVEL
	LogLevel zerolog.Level

	// Env is current environment
	// ENV: Env
	Env Environment

	// PORT is the TCP port for web workers
	// ENV: PORT
	Port uint16

	// Kafka is the kafka configuration
	// Not loaded from environment variables directly
	Kafka K
}

// K struct contains the configuration for kafka brokers
type K struct {
	// BootstrapServers for kafka brokers
	BootstrapServers []string
}

// LoadEnv loads the configuration from appropriate environment variables
func LoadEnv() (C, error) {
	var c C

	if p := os.Getenv("PORT"); len(p) > 0 {
		u, err := strconv.ParseUint(p, 10, 16)
		if err != nil {
			return C{}, fmt.Errorf("failed to parse PORT: %w", err)
		}

		c.Port = uint16(u)
	}

	ll := os.Getenv("LOG_LEVEL")
	if len(ll) == 0 {
		ll = "info"
	}

	l, err := zerolog.ParseLevel(ll)
	if err != nil {
		return C{}, fmt.Errorf("failed to parse GOPHER_LOG_LEVEL: %w", err)
	}

	c.LogLevel = l
	c.Env = strToEnv(os.Getenv("ENV"))
	bootstrapServers := os.Getenv("BOOTSTRAP_SERVERS")
	if bootstrapServers != "" {
		c.Kafka.BootstrapServers = strings.Split(bootstrapServers, ",")
	}
	return c, nil
}

// DefaultLogger returns a zerolog.Logger using settings from our config struct.
func DefaultLogger(config C) zerolog.Logger {
	// set up zerolog
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.SetGlobalLevel(config.LogLevel)

	// set up logging
	return zerolog.New(os.Stdout).
		With().Timestamp().Logger()
}
