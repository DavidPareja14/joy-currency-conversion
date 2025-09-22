package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/joy-currency-conversion-private/domain"
	"github.com/joy-currency-conversion-private/infrastructure/response"
)

// CurrencyService implements domain.CurrencyService using AWS services
type CurrencyService struct {
	dynamoDB *dynamodb.DynamoDB
	ExchangeRateAPIKey string
}

// NewCurrencyService creates a new CurrencyService
func NewCurrencyService(dynamoDB *dynamodb.DynamoDB, exchangeRateAPIKey string) *CurrencyService {
	return &CurrencyService{
		dynamoDB: dynamoDB,
		ExchangeRateAPIKey: exchangeRateAPIKey,
	}
}

type exchangeRateResponse struct {
	ConversionRate float64 `json:"conversion_rate"`
}

// GetExchangeRate returns the current exchange rate between two currencies
func (s *CurrencyService) GetExchangeRateGivenAmount(ctx context.Context, origin, destination string, amount float64) (response.ExchangeRateResponse, error) {
	// TODO: Implement actual exchange rate fetching
	// This could integrate with external APIs like:
	// - ExchangeRate-API
	// - Fixer.io
	// - CurrencyLayer
	// - Or store rates in DynamoDB and update them periodically
	
	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/pair/%s/%s/%.3f", s.ExchangeRateAPIKey, origin, destination, amount)

	// Make the GET request
	resp, err := http.Get(url)
	if err != nil {
		return response.ExchangeRateResponse{}, fmt.Errorf("error when get the conversion rate for %s to %s", origin, destination)
	}
	defer resp.Body.Close()

	// Check the HTTP status code
	if resp.StatusCode != http.StatusOK {
		return response.ExchangeRateResponse{}, fmt.Errorf("unexpected status code: %d %s", resp.StatusCode, resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response.ExchangeRateResponse{}, fmt.Errorf("error reading response body: %v", err)
	}

	var exchangeRateResponse response.ExchangeRateResponse
	err = json.Unmarshal(body, &exchangeRateResponse)
	if err != nil {
		return response.ExchangeRateResponse{}, fmt.Errorf("error unmarshalling response body: %v", err)
	}
	exchangeRateResponse.RatesSource = "exchange-rate-api"

	return exchangeRateResponse, nil
}


// GetExchangeRate returns the current exchange rate between two currencies
func (s *CurrencyService) GetExchangeRate(ctx context.Context, origin, destination string) (float64, string, error) {
	// TODO: Implement actual exchange rate fetching
	// This could integrate with external APIs like:
	// - ExchangeRate-API
	// - Fixer.io
	// - CurrencyLayer
	// - Or store rates in DynamoDB and update them periodically
	
	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/pair/%s/%s", s.ExchangeRateAPIKey, origin, destination)

	// Make the GET request
	resp, err := http.Get(url)
	if err != nil {
		return 0, "", fmt.Errorf("error when get the conversion rate for %s to %s", origin, destination)
	}
	defer resp.Body.Close()

	// Check the HTTP status code
	if resp.StatusCode != http.StatusOK {
		return 0, "", fmt.Errorf("unexpected status code: %d %s", resp.StatusCode, resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", fmt.Errorf("error reading response body: %v", err)
	}

	var exchangeRateResponse exchangeRateResponse
	err = json.Unmarshal(body, &exchangeRateResponse)
	if err != nil {
		return 0, "", fmt.Errorf("error unmarshalling response body: %v", err)
	}

	return exchangeRateResponse.ConversionRate, "exchange-rate-api", nil
	
	// For now, return mock data
	/*
	mockRates := map[string]float64{
		"COP-USD": 0.00025,
		"USD-COP": 4000.0,
		"EUR-USD": 1.05,
		"USD-EUR": 0.95,
	}
	
	key := fmt.Sprintf("%s-%s", origin, destination)
	if rate, exists := mockRates[key]; exists {
		return rate, "mock-provider", nil
	}
	
	
	return 0, "", fmt.Errorf("exchange rate not available for %s to %s", origin, destination)
	*/
}

// GetHistoricalRates returns historical exchange rates for a date range
func (s *CurrencyService) GetHistoricalRates(ctx context.Context, origin, destination string, startDate, endDate time.Time) ([]domain.HistoryRate, string, error) {
	// TODO: Implement historical rate fetching
	// This could:
	// - Query DynamoDB for stored historical rates
	// - Call external APIs for historical data
	// - Use AWS Lambda to fetch and cache historical data
	
	// For now, return mock data
	var rates []domain.HistoryRate
	current := startDate
	for current.Before(endDate) || current.Equal(endDate) {
		// Mock rate with some variation
		baseRate := 0.00025
		variation := float64(current.Day()%10) * 0.00001
		rate := baseRate + variation
		
		rates = append(rates, domain.HistoryRate{
			Date: current.Format("2006-01-02"),
			Rate: rate,
		})
		current = current.AddDate(0, 0, 1)
	}
	
	return rates, "mock-provider", nil
}

// GetForecast returns a forecast for the next day's exchange rate
func (s *CurrencyService) GetForecast(ctx context.Context, origin, destination string) (*domain.ForecastResponse, error) {
	// TODO: Implement forecast logic
	// This could:
	// - Use historical data to calculate trends
	// - Integrate with ML services like AWS SageMaker
	// - Use simple statistical methods
	
	// Get currency information
	originCurrency, err := s.GetCurrencyInfo(ctx, origin)
	if err != nil {
		return nil, err
	}
	
	destCurrency, err := s.GetCurrencyInfo(ctx, destination)
	if err != nil {
		return nil, err
	}
	
	// Mock forecast
	tomorrow := time.Now().AddDate(0, 0, 1)
	
	response := &domain.ForecastResponse{
		Origin:        *originCurrency,
		Destination:   *destCurrency,
		PredictedDate: tomorrow.Format("2006-01-02"),
		PredictedRate: 0.00026, // Mock prediction
		Confidence:    0.65,    // Mock confidence
		Last30Days: struct {
			Average float64 `json:"average"`
		}{
			Average: 0.000255, // Mock average
		},
		Timestamp:   time.Now().UTC(),
		RatesSource: "mock-provider",
	}
	
	return response, nil
}

// GetSupportedDestinations returns supported destination currencies for an origin
func (s *CurrencyService) GetSupportedDestinations(ctx context.Context, origin string) ([]domain.Currency, string, error) {
	// TODO: Implement destination lookup
	// This could:
	// - Query DynamoDB for supported currency pairs
	// - Use a predefined list of supported currencies
	// - Call external APIs to get available destinations
	
	// Mock supported destinations
	destinations := map[string][]domain.Currency{
		"COP": {
			{Code: "USD", Country: "United States"},
			{Code: "EUR", Country: "Eurozone"},
		},
		"USD": {
			{Code: "COP", Country: "Colombia"},
			{Code: "EUR", Country: "Eurozone"},
		},
		"EUR": {
			{Code: "USD", Country: "United States"},
			{Code: "COP", Country: "Colombia"},
		},
	}
	
	if dests, exists := destinations[origin]; exists {
		return dests, "mock-provider", nil
	}
	
	return nil, "", fmt.Errorf("origin currency %s not supported", origin)
}

// GetCurrencyInfo returns currency information (code and country)
func (s *CurrencyService) GetCurrencyInfo(ctx context.Context, code string) (*domain.Currency, error) {
	// TODO: Implement currency info lookup
	// This could:
	// - Query DynamoDB for currency information
	// - Use a predefined map of currency codes to countries
	// - Call external APIs for currency metadata
	
	// Mock currency information
	currencies := map[string]domain.Currency{
		"COP": {Code: "COP", Country: "Colombia"},
		"USD": {Code: "USD", Country: "United States"},
		"EUR": {Code: "EUR", Country: "Eurozone"},
	}
	
	if currency, exists := currencies[code]; exists {
		return &currency, nil
	}
	
	return nil, fmt.Errorf("currency code %s not found", code)
}
