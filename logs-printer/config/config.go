package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	LogLevel         string `env:"LOG_LEVEL"                  envDefault:"error"`
	RabbitMqHost     string `env:"RABBIT_MQ_HOST"             envDefault:"localhost"`
	RabbitMqPort     string `env:"RABBIT_MQ_PORT"             envDefault:"5672"`
	RabbitMqUserName string `env:"RABBIT_MQ_USERNAME"         envDefault:"guest"`
	RabbitMqPassword string `env:"RABBIT_MQ_PASSWORD"         envDefault:"guest"`
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
