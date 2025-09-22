package config

import (
	"fmt"
	"os"
)

type Config struct {
	KyeEchangeRateAPI string
}

func LoadConfig() (*Config, error) {
	exchangeKey := os.Getenv("EXCHANGE_RATE_API_KEY")
	if exchangeKey == "" {
		return &Config{}, fmt.Errorf("EXCHANGE_RATE_API_KEY is not set")
	}

	return &Config{
		KyeEchangeRateAPI: exchangeKey,
	}, nil
}