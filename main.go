package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joy-currency-conversion-private/handlers"
	"github.com/joy-currency-conversion-private/infrastructure"
)

func main() {
	// Initialize AWS services
	awsServices := infrastructure.NewAWSServices()

	// Initialize handlers
	currencyHandler := handlers.NewCurrencyHandler(awsServices)

	// Setup Gin router
	router := gin.Default()

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Endpoint 1: Currency Conversion
		v1.GET("/convert", currencyHandler.Convert)

		// Endpoint 2: Daily Historic Values
		v1.GET("/history", currencyHandler.History)

		// Endpoint 3: Probability Forecast
		v1.GET("/forecast", currencyHandler.Forecast)

		// Endpoint 4: Available Destination Currencies
		v1.GET("/origins/:origin/destinations", currencyHandler.GetDestinations)

		// Endpoint 5: Save a Favorite Conversion
		v1.POST("/favorites", currencyHandler.SaveFavorite)

		// Endpoint 6: Daily Favorite Check
		v1.POST("/favorites/check", currencyHandler.CheckFavorites)

		// Endpoint 7: Email Notification
		v1.POST("/notifications/email", currencyHandler.SendNotification)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	log.Println("Starting Project Joy API server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
