package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	StorageFilename        string `env:"STORAGE_FILENAME"           envDefault:"emails.storage"`
	CurrencyServiceBaseURL string `env:"CURRENCY_SERVICE_BASE_URL"`
	FromCurrency           string `env:"FROM_CURRENCY"              envDefault:"BTC"`
	ToCurrency             string `env:"TO_CURRENCY"                envDefault:"UAH"`
	SMTPHost               string `env:"SMTP_HOST"                  envDefault:"smtp.gmail.com"`
	SMTPPort               int    `env:"SMTP_PORT"                  envDefault:"587"`
	SMTPUsername           string `env:"SMTP_USERNAME"`
	SMTPPassword           string `env:"SMTP_PASSWORD"`
	ServerHost             string `env:"SERVER_HOST"                envDefault:"localhost"`
	ServerPort             string `env:"SERVER_PORT"                envDefault:"80"`
	LogLevel               int    `env:"LOG_LEVEL"                  envDefault:"5"`
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
