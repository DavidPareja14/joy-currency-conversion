package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/joy-currency-conversion-private/domain"
	"github.com/joy-currency-conversion-private/infrastructure/response"
)

// CurrencyService implements domain.CurrencyService using AWS services
type CurrencyService struct {
	dynamoDB            *dynamodb.DynamoDB
	ExchangeRateAPIKey  string
	ExchangeRatesAPIKey string
}

// NewCurrencyService creates a new CurrencyService
func NewCurrencyService(dynamoDB *dynamodb.DynamoDB, exchangeRateAPIKey, echangeRatesAPIKey string) *CurrencyService {
	return &CurrencyService{
		dynamoDB:            dynamoDB,
		ExchangeRateAPIKey:  exchangeRateAPIKey,
		ExchangeRatesAPIKey: echangeRatesAPIKey,
	}
}

type exchangeRateResponse struct {
	ConversionRate float64 `json:"conversion_rate"`
}

type exchangeRatesResponse struct {
	Historical bool               `json:"historical"`
	Date       string             `json:"date"`
	Base       string             `json:"base"`
	Rates      map[string]float64 `json:"rates"`
}

// GetExchangeRate returns the current exchange rate between two currencies
func (s *CurrencyService) GetExchangeRateGivenAmount(ctx context.Context, origin, destination string, amount float64) (response.ExchangeRateResponse, error) {
	// TODO: Implement actual exchange rate fetching
	// This could integrate with external APIs like:
	// - ExchangeRate-API
	// - Fixer.io
	// - CurrencyLayer
	// - Or store rates in DynamoDB and update them periodically

	// This was built with https://app.exchangerate-api.com/ api
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

	// This was built with https://app.exchangerate-api.com/ api
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

	// Other client used https://manage.exchangeratesapi.io/

	// For now, return mock data
	/*
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
	*/
	fmt.Printf("*** End date: %s and Start date: %s", startDate, endDate)
	if startDate.After(time.Now()) || endDate.After(time.Now()) {
		return []domain.HistoryRate{}, "", fmt.Errorf("the start date and end date must not be greater than the current date")
	}

	if startDate.After(endDate) {
		return []domain.HistoryRate{}, "", fmt.Errorf("the start date must not be greater than end date")
	}

	if endDate.Sub(startDate).Hours()/24 > 5 {
		return []domain.HistoryRate{}, "", fmt.Errorf("the difference between start date and end date must not be greater than 5")
	}
	var rates []domain.HistoryRate
	current := startDate
	for current.Before(endDate) || current.Equal(endDate) {

		// No wokr the query param base and symbols, by the fault, the base is EUR, i think i can't use the endpoint with another base
		url := fmt.Sprintf("https://api.exchangeratesapi.io/v1/%s?access_key=%s&base=%s&symbols=%s", current.Format("2006-01-02"), s.ExchangeRatesAPIKey, "EUR", destination)

		// Make the GET request
		resp, err := http.Get(url)
		if err != nil {
			return []domain.HistoryRate{}, "", fmt.Errorf("error when get the historical data for %s to %s, date: %s", origin, destination, current)
		}
		defer resp.Body.Close()

		// Check the HTTP status code
		if resp.StatusCode != http.StatusOK {
			return []domain.HistoryRate{}, "", fmt.Errorf("unexpected status code: %d %s", resp.StatusCode, resp.Status)
		}

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return []domain.HistoryRate{}, "", fmt.Errorf("error reading response body: %v", err)
		}

		var exchangeRatesResponse exchangeRatesResponse
		err = json.Unmarshal(body, &exchangeRatesResponse)
		if err != nil {
			return []domain.HistoryRate{}, "", fmt.Errorf("error unmarshalling response body: %v", err)
		}

		rates = append(rates, domain.HistoryRate{
			Date: current.Format("2006-01-02"),
			Rate: exchangeRatesResponse.Rates[destination],
		})
		current = current.AddDate(0, 0, 1)
		time.Sleep(1 * time.Second)
	}

	return rates, "api.exchangeratesapi.io", nil
}

// GetForecast returns a forecast for the next day's exchange rate
func (s *CurrencyService) GetForecast(ctx context.Context, origin, destination string) (*domain.ForecastResponse, error) {
	// Get currency information
	originCurrency, err := s.GetCurrencyInfo(ctx, origin)
	if err != nil {
		return nil, err
	}

	// Temporal code because currently the API  only allows EUR origin
	if origin != "EUR" {
		originCurrency, _ = s.GetCurrencyInfo(ctx, "EUR")
	}

	destCurrency, err := s.GetCurrencyInfo(ctx, destination)
	if err != nil {
		return nil, err
	}

	// Calculate date range for last 5 days (to respect API limitations)
	endDate := time.Now().AddDate(0, 0, -1) // Yesterday (since today's data might not be available)
	startDate := endDate.AddDate(0, 0, -4)  // 5 days ago

	// Get historical data for the last 5 days
	historicalRates, source, err := s.GetHistoricalRates(ctx, origin, destination, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("unable to get historical data for forecast: %w", err)
	}

	// Check if we have enough data (at least 3 days)
	if len(historicalRates) < 3 {
		return nil, fmt.Errorf("insufficient historical data for forecast (need at least 3 days, got %d)", len(historicalRates))
	}

	// Calculate statistics from historical data
	var sum float64
	var rates []float64

	for _, rate := range historicalRates {
		sum += rate.Rate
		rates = append(rates, rate.Rate)
	}

	average := sum / float64(len(rates))

	// Calculate standard deviation for confidence estimation
	var variance float64
	for _, rate := range rates {
		diff := rate - average
		variance += diff * diff
	}
	variance = variance / float64(len(rates))
	stdDev := math.Sqrt(variance)

	// Simple trend analysis: compare first half vs second half
	midPoint := len(rates) / 2
	firstHalfAvg := 0.0
	secondHalfAvg := 0.0

	for i := 0; i < midPoint; i++ {
		firstHalfAvg += rates[i]
	}
	firstHalfAvg = firstHalfAvg / float64(midPoint)

	for i := midPoint; i < len(rates); i++ {
		secondHalfAvg += rates[i]
	}
	secondHalfAvg = secondHalfAvg / float64(len(rates)-midPoint)

	// Calculate trend (positive = increasing, negative = decreasing)
	trend := (secondHalfAvg - firstHalfAvg) / firstHalfAvg

	// Predict next day's rate using simple linear trend
	// Use a conservative approach: 50% of the trend + 50% of the average
	trendAdjustment := average * trend * 0.5
	predictedRate := average + trendAdjustment

	// Ensure predicted rate is positive
	if predictedRate <= 0 {
		predictedRate = average
	}

	// Calculate confidence based on data consistency
	// Lower standard deviation = higher confidence
	// More data points = higher confidence
	confidence := 0.5 // Base confidence

	// Adjust confidence based on standard deviation (lower std dev = higher confidence)
	if stdDev > 0 {
		coefficientOfVariation := stdDev / average
		if coefficientOfVariation < 0.05 { // Less than 5% variation
			confidence += 0.3
		} else if coefficientOfVariation < 0.1 { // Less than 10% variation
			confidence += 0.2
		} else if coefficientOfVariation < 0.2 { // Less than 20% variation
			confidence += 0.1
		}
	}

	// Adjust confidence based on data points (adjusted for 5-day period)
	if len(rates) >= 5 {
		confidence += 0.1
	} else if len(rates) >= 4 {
		confidence += 0.05
	}

	// Cap confidence at 0.9 (90%)
	if confidence > 0.9 {
		confidence = 0.9
	}

	// Ensure minimum confidence
	if confidence < 0.3 {
		confidence = 0.3
	}

	// Calculate predicted date (tomorrow)
	tomorrow := time.Now().AddDate(0, 0, 1)

	response := &domain.ForecastResponse{
		Origin:        *originCurrency,
		Destination:   *destCurrency,
		PredictedDate: tomorrow.Format("2006-01-02"),
		PredictedRate: predictedRate,
		Confidence:    confidence,
		Last30Days: struct {
			Average float64 `json:"average"`
		}{
			Average: average,
		},
		Timestamp:   time.Now().UTC(),
		RatesSource: source,
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
