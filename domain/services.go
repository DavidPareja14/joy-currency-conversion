package domain

import (
	"context"
	"time"

	"github.com/joy-currency-conversion-private/infrastructure/response"
)

// CurrencyService defines the interface for currency-related operations
type CurrencyService interface {
	// GetExchangeRate returns the current exchange rate between two currencies
	GetExchangeRate(ctx context.Context, origin, destination string) (float64, string, error)

	// GetExchangeRateGivenAmount returns the current exchange rate between two currencies given an amount
	GetExchangeRateGivenAmount(ctx context.Context, origin, destination string, amount float64) (response.ExchangeRateResponse, error)
	
	// GetHistoricalRates returns historical exchange rates for a date range
	GetHistoricalRates(ctx context.Context, origin, destination string, startDate, endDate time.Time) ([]HistoryRate, string, error)
	
	// GetForecast returns a forecast for the next day's exchange rate
	GetForecast(ctx context.Context, origin, destination string) (*ForecastResponse, error)
	
	// GetSupportedDestinations returns supported destination currencies for an origin
	GetSupportedDestinations(ctx context.Context, origin string) ([]Currency, string, error)
	
	// GetCurrencyInfo returns currency information (code and country)
	GetCurrencyInfo(ctx context.Context, code string) (*Currency, error)
}

// FavoriteService defines the interface for favorite-related operations
type FavoriteService interface {
	// SaveFavorite saves a new favorite conversion
	SaveFavorite(ctx context.Context, req *FavoriteRequest) (*Favorite, error)
	
	// GetAllFavorites returns all saved favorites
	GetAllFavorites(ctx context.Context) ([]Favorite, error)
	
	// CheckFavorites checks all favorites against current rates
	CheckFavorites(ctx context.Context) (*FavoriteCheckResponse, error)
}

// NotificationService defines the interface for notification operations
type NotificationService interface {
	// SendEmailNotification sends an email notification
	SendEmailNotification(ctx context.Context, req *NotificationRequest) (*NotificationResponse, error)
}
