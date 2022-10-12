package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	PostgresURL      string `env:"POSTGRES_URL"`
	ServerHost       string `env:"SERVER_HOST"                envDefault:"localhost"`
	ServerPort       string `env:"SERVER_PORT"                envDefault:"80"`
	LogLevel         int    `env:"LOG_LEVEL"                  envDefault:"5"`
	RateValueBitSize int    `env:"RATE_VALUE_BIT_SIZE"        envDefault:"64"`
}

func NewAppConfig(envFilename string) (*AppConfig, error) {
	if err := godotenv.Load(envFilename); err != nil {
		return nil, LoadError{Message: err.Error()}
	}

	cfg := &AppConfig{}
	if err := env.Parse(cfg); err != nil {
		return nil, LoadError{Message: err.Error()}
	}

	return cfg, nil
}
