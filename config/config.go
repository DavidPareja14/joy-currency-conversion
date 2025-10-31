package config

import (
	// "context"
	"fmt"
	"os"

	// "github.com/aws/aws-sdk-go-v2/aws"
    // "github.com/aws/aws-sdk-go-v2/config"
    // "github.com/aws/aws-sdk-go-v2/service/ssm"
)

const (
	name = "/exchange-rate/api-key"
)

type Config struct {
	KyeEchangeRateAPI string
	KyeEchangeRatesAPI string
}

func LoadConfig() (*Config, error) {
	exchangeKey := os.Getenv("EXCHANGE_RATE_API_KEY")
	if exchangeKey == "" {
		return &Config{}, fmt.Errorf("EXCHANGE_RATE_API_KEY is not set")
	}

	// Fetching secrets directly from AWS Parameter Store
	// ctx := context.TODO()

    // // Load AWS default configuration (uses EC2 IAM role automatically)
    // cfg, err := config.LoadDefaultConfig(ctx)
    // if err != nil {
    //     return &Config{}, fmt.Errorf("failed to load AWS config: %w", err)
    // }

    // client := ssm.NewFromConfig(cfg)
    // param, err := client.GetParameter(ctx, &ssm.GetParameterInput{
    //     Name:           aws.String(name),
    //     WithDecryption: aws.Bool(true),
    // })
    // if err != nil {
    //     return &Config{}, fmt.Errorf("EXCHANGE_RATE_API_KEY is not set, Error: %w", err)
    // }

    // exchangeKey := *param.Parameter.Value

	// Fetching secrets from .env file
	exchangeRatesKey := os.Getenv("EXCHANGE_RATES_API_KEY")
	if exchangeRatesKey == "" {
		return &Config{}, fmt.Errorf("EXCHANGE_RATES_API_KEY is not set")
	}

	return &Config{
		KyeEchangeRateAPI: exchangeKey,
		KyeEchangeRatesAPI: exchangeRatesKey,
	}, nil
}