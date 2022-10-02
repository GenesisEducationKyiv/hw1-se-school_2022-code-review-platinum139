package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	StorageFilename  string `env:"STORAGE_FILENAME"           envDefault:"emails.storage"`
	RateValueBitSize int    `env:"RATE_VALUE_BIT_SIZE"        envDefault:"64"`
	FromCurrency     string `env:"FROM_CURRENCY"              envDefault:"BTC"`
	ToCurrency       string `env:"TO_CURRENCY"                envDefault:"UAH"`
	CurrencySvcHost  string `env:"CURRENCY_SERVICE_HOST"      envDefault:"localhost"`
	CurrencySvcPort  string `env:"CURRENCY_SERVICE_PORT"      envDefault:"8000"`
	SMTPHost         string `env:"SMTP_HOST"                  envDefault:"smtp.gmail.com"`
	SMTPPort         int    `env:"SMTP_PORT"                  envDefault:"587"`
	SMTPUsername     string `env:"SMTP_USERNAME"`
	SMTPPassword     string `env:"SMTP_PASSWORD"`
	ServerHost       string `env:"SERVER_HOST"                envDefault:"localhost"`
	ServerPort       string `env:"SERVER_PORT"                envDefault:"80"`
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
