package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	StorageFilename string `env:"STORAGE_FILENAME" envDefault:"emails.storage"`
	FromCurrency    string `env:"FROM_CURRENCY"    envDefault:"BTC"`
	ToCurrency      string `env:"TO_CURRENCY"      envDefault:"UAH"`
	SmtpHost        string `env:"SMTP_HOST"        envDefault:"smtp.gmail.com"`
	SmtpPort        int    `env:"SMTP_PORT"        envDefault:"587"`
	SmtpUsername    string `env:"SMTP_USERNAME"`
	SmtpPassword    string `env:"SMTP_PASSWORD"`
	ServerHost      string `env:"SERVER_HOST"      envDefault:"localhost"`
	ServerPort      string `env:"SERVER_PORT"      envDefault:"80"`
	LogLevel        int    `env:"LOG_LEVEL"        envDefault:"5"`
}

func NewAppConfig() (*AppConfig, error) {
	if err := godotenv.Load(); err != nil {
		return nil, LoadError{Message: err.Error()}
	}

	cfg := &AppConfig{}
	if err := env.Parse(cfg); err != nil {
		return nil, LoadError{Message: err.Error()}
	}

	return cfg, nil
}
