package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	RateValueBitSize int    `env:"RATE_VALUE_BIT_SIZE"        envDefault:"64"`
	FromCurrency     string `env:"FROM_CURRENCY"              envDefault:"BTC"`
	ToCurrency       string `env:"TO_CURRENCY"                envDefault:"UAH"`
	CurrencyProvider string `env:"CURRENCY_PROVIDER"`
	CachingPeriodMin int    `env:"CACHING_PERIOD_MIN"         envDefault:"5"`
	ServerHost       string `env:"SERVER_HOST"                envDefault:"localhost"`
	ServerPort       string `env:"SERVER_PORT"                envDefault:"80"`
	LogLevel         int    `env:"LOG_LEVEL"                  envDefault:"5"`
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
