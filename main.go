package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joy-currency-conversion-private/config"
	"github.com/joy-currency-conversion-private/handlers"
	"github.com/joy-currency-conversion-private/infrastructure"
	"github.com/joy-currency-conversion-private/infrastructure/db"
)

func main() {
	configuratios, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	if err := db.Connect(); err != nil {
		log.Fatalf("db connect: %v", err)
	}

	// Initialize AWS services
	awsServices := infrastructure.NewAWSServices(configuratios.KyeEchangeRateAPI, configuratios.KyeEchangeRatesAPI)

	// Initialize handlers
	currencyHandler := handlers.NewCurrencyHandler(awsServices)

	// Setup Chi router
	router := chi.NewRouter()

	// Add middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)

	// Health check
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	// API v1 routes
	router.Route("/api/v1", func(r chi.Router) {
		// Endpoint 1: Currency Conversion
		r.Get("/convert", currencyHandler.Convert)

		// Endpoint 2: Daily Historic Values
		r.Get("/history", currencyHandler.History)

		// Endpoint 3: Probability Forecast
		r.Get("/forecast", currencyHandler.Forecast)

		// Endpoint 4: Available Destination Currencies
		r.Get("/origins/{origin}/destinations", currencyHandler.GetDestinations)

		// Endpoint 5: Save a Favorite Conversion
		r.Post("/favorites", currencyHandler.SaveFavorite)

		// Endpoint 6: Daily Favorite Check
		r.Post("/favorites/check", currencyHandler.CheckFavorites)

		// Endpoint 7: Email Notification
		r.Post("/notifications/email", currencyHandler.SendNotification)
	})

	log.Println("Starting Project Joy API server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
