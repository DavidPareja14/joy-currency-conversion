package domain

import (
	"time"
)

// Currency represents a currency with its code and country
type Currency struct {
	Code    string `json:"code"`
	Country string `json:"country"`
}

// ConversionRequest represents the request for currency conversion
type ConversionRequest struct {
	Origin      string  `json:"origin" binding:"required"`
	Destination string  `json:"destination" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
}

// ConversionResponse represents the response for currency conversion
type ConversionResponse struct {
	Origin          Currency  `json:"origin"`
	Destination     Currency  `json:"destination"`
	Rate            float64   `json:"rate"`
	Amount          float64   `json:"amount"`
	ConvertedAmount float64   `json:"converted_amount"`
	Timestamp       time.Time `json:"timestamp"`
	RatesSource     string    `json:"rates_source"`
}

// HistoryRequest represents the request for historical data
type HistoryRequest struct {
	Origin     string `json:"origin" binding:"required"`
	Destination string `json:"destination" binding:"required"`
	StartDate  string `json:"start_date" binding:"required"`
	EndDate    string `json:"end_date" binding:"required"`
}

// HistoryRate represents a single historical rate entry
type HistoryRate struct {
	Date string  `json:"date"`
	Rate float64 `json:"rate"`
}

// HistoryResponse represents the response for historical data
type HistoryResponse struct {
	Origin      Currency      `json:"origin"`
	Destination Currency      `json:"destination"`
	StartDate   string        `json:"start_date"`
	EndDate     string        `json:"end_date"`
	Rates       []HistoryRate `json:"rates"`
	Timestamp   time.Time     `json:"timestamp"`
	RatesSource string        `json:"rates_source"`
}

// ForecastRequest represents the request for forecast
type ForecastRequest struct {
	Origin      string `json:"origin" binding:"required"`
	Destination string `json:"destination" binding:"required"`
}

// ForecastResponse represents the response for forecast
type ForecastResponse struct {
	Origin        Currency  `json:"origin"`
	Destination   Currency  `json:"destination"`
	PredictedDate string    `json:"predicted_date"`
	PredictedRate float64   `json:"predicted_rate"`
	Confidence    float64   `json:"confidence"`
	Last30Days    struct {
		Average float64 `json:"average"`
	} `json:"last_30_days"`
	Timestamp   time.Time `json:"timestamp"`
	RatesSource string    `json:"rates_source"`
}

// DestinationsResponse represents the response for available destinations
type DestinationsResponse struct {
	Origin      Currency   `json:"origin"`
	Destinations []Currency `json:"destinations"`
	Timestamp   time.Time  `json:"timestamp"`
	RatesSource string     `json:"rates_source"`
}

// FavoriteRequest represents the request to save a favorite
type FavoriteRequest struct {
	Origin      string  `json:"origin" binding:"required"`
	Destination string  `json:"destination" binding:"required"`
	Threshold   float64 `json:"threshold" binding:"required,gt=0"`
	NotifyEmail string  `json:"notify_email" binding:"required,email"`
}

// Favorite represents a saved favorite conversion
type Favorite struct {
	ID          string    `json:"id"`
	Origin      Currency  `json:"origin"`
	Destination Currency  `json:"destination"`
	Threshold   float64   `json:"threshold"`
	NotifyEmail string    `json:"notify_email"`
	CreatedAt   time.Time `json:"created_at"`
}

// FavoriteCheckResult represents the result of checking a favorite
type FavoriteCheckResult struct {
	FavoriteID        string    `json:"favorite_id"`
	Origin            Currency  `json:"origin"`
	Destination       Currency  `json:"destination"`
	Threshold         float64   `json:"threshold"`
	CurrentRate       float64   `json:"current_rate"`
	Date              string    `json:"date"`
	Exceeded          bool      `json:"exceeded"`
	Notified          bool      `json:"notified"`
	CurrentRateSource string    `json:"current_rate_source"`
}

// FavoriteCheckResponse represents the response for favorite checks
type FavoriteCheckResponse struct {
	Results   []FavoriteCheckResult `json:"results"`
	Timestamp time.Time             `json:"timestamp"`
}

// NotificationRequest represents the request to send a notification
type NotificationRequest struct {
	FavoriteID  string    `json:"favorite_id" binding:"required"`
	Origin      Currency  `json:"origin" binding:"required"`
	Destination Currency  `json:"destination" binding:"required"`
	Threshold   float64   `json:"threshold" binding:"required"`
	CurrentRate float64   `json:"current_rate" binding:"required"`
	Date        string    `json:"date" binding:"required"`
	NotifyEmail string    `json:"notify_email" binding:"required,email"`
}

// NotificationResponse represents the response for notifications
type NotificationResponse struct {
	Message string `json:"message"`
	SentTo  string `json:"sent_to"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string      `json:"error"`
	Code    string      `json:"code"`
	Details interface{} `json:"details,omitempty"`
}
