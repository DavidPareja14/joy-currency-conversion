package infrastructure

import (
	"context"
	"fmt"
	"time"

	/*
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	*/
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	"github.com/joy-currency-conversion-private/domain"
)

// FavoriteService implements domain.FavoriteService using DynamoDB
type FavoriteService struct {
	dynamoDB *dynamodb.DynamoDB
	ExchangeRateAPIKey string
	ExchangeRatesAPIKey string
}

// NewFavoriteService creates a new FavoriteService
func NewFavoriteService(dynamoDB *dynamodb.DynamoDB, exchangeRateAPIKey, exchangeRatesAPIKey string) *FavoriteService {
	return &FavoriteService{
		dynamoDB: dynamoDB,
		ExchangeRateAPIKey: exchangeRateAPIKey,
		ExchangeRatesAPIKey: exchangeRatesAPIKey,
	}
}

// SaveFavorite saves a new favorite conversion
func (s *FavoriteService) SaveFavorite(ctx context.Context, req *domain.FavoriteRequest) (*domain.Favorite, error) {
	// Generate UUID for the favorite
	id := uuid.New().String()
	
	// Get currency information
	currencyService := NewCurrencyService(s.dynamoDB, s.ExchangeRateAPIKey, s.ExchangeRatesAPIKey)
	originCurrency, err := currencyService.GetCurrencyInfo(ctx, req.Origin)
	if err != nil {
		return nil, fmt.Errorf("invalid origin currency: %w", err)
	}
	
	destCurrency, err := currencyService.GetCurrencyInfo(ctx, req.Destination)
	if err != nil {
		return nil, fmt.Errorf("invalid destination currency: %w", err)
	}
	
	// Create favorite object
	favorite := &domain.Favorite{
		ID:          id,
		Origin:      *originCurrency,
		Destination: *destCurrency,
		Threshold:   req.Threshold,
		NotifyEmail: req.NotifyEmail,
		CreatedAt:   time.Now().UTC(),
	}
	
	// TODO: Save to DynamoDB
	// This would involve:
	// 1. Creating a DynamoDB item with the favorite data
	// 2. Using PutItem or UpdateItem to store it
	// 3. Handling conditional writes to prevent duplicates
	
	// For now, just return the favorite (mock implementation)
	return favorite, nil
}

// GetAllFavorites returns all saved favorites
func (s *FavoriteService) GetAllFavorites(ctx context.Context) ([]domain.Favorite, error) {
	// TODO: Implement DynamoDB scan or query
	// This would involve:
	// 1. Using Scan or Query operation to get all favorites
	// 2. Unmarshaling the results into domain.Favorite objects
	// 3. Handling pagination if there are many favorites
	
	favorites := make([]domain.Favorite, 0)
	favorites = append(favorites, domain.Favorite{
		ID: "abcd",
		Origin: domain.Currency{Code: "USD", Country: "EE.UU"},
		Destination: domain.Currency{Code: "COP", Country: "Colombia"},
		Threshold: 3000,
		NotifyEmail: "test@gmail.com",
	})
	return favorites, nil
}

// CheckFavorites checks all favorites against current rates
func (s *FavoriteService) CheckFavorites(ctx context.Context) (*domain.FavoriteCheckResponse, error) {
	// Get all favorites
	favorites, err := s.GetAllFavorites(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get favorites: %w", err)
	}
	
	// Initialize currency service for rate checking
	currencyService := NewCurrencyService(s.dynamoDB, s.ExchangeRateAPIKey, s.ExchangeRatesAPIKey)
	
	var results []domain.FavoriteCheckResult
	today := time.Now().Format("2006-01-02")
	
	// Check each favorite
	for _, favorite := range favorites {
		// Get current rate
		currentRate, source, err := currencyService.GetExchangeRate(
			ctx, 
			favorite.Origin.Code, 
			favorite.Destination.Code,
		)
		if err != nil {
			// Skip this favorite if we can't get the rate
			continue
		}
		
		// Check if threshold is exceeded
		exceeded := currentRate >= favorite.Threshold
		
		// TODO: Send notification if threshold is exceeded
		// This would involve calling the notification service
		notified := false
		if exceeded {
			// Mock notification sending
			notified = true
		}
		
		result := domain.FavoriteCheckResult{
			FavoriteID:        favorite.ID,
			Origin:            favorite.Origin,
			Destination:       favorite.Destination,
			Threshold:         favorite.Threshold,
			CurrentRate:       currentRate,
			Date:              today,
			Exceeded:          exceeded,
			Notified:          notified,
			CurrentRateSource: source,
		}
		
		results = append(results, result)
	}
	
	response := &domain.FavoriteCheckResponse{
		Results:   results,
		Timestamp: time.Now().UTC(),
	}
	
	return response, nil
}
