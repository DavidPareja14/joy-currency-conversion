package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joy-currency-conversion-private/domain"
	"github.com/joy-currency-conversion-private/infrastructure"
)

// CurrencyHandler handles all currency-related HTTP requests
type CurrencyHandler struct {
	awsServices *infrastructure.AWSServices
}

// NewCurrencyHandler creates a new CurrencyHandler
func NewCurrencyHandler(awsServices *infrastructure.AWSServices) *CurrencyHandler {
	return &CurrencyHandler{
		awsServices: awsServices,
	}
}

// Convert handles currency conversion requests
// GET /api/v1/convert?origin={ORIGIN}&destination={DEST}&amount={AMOUNT}
func (h *CurrencyHandler) Convert(w http.ResponseWriter, r *http.Request) {
	origin := r.URL.Query().Get("origin")
	destination := r.URL.Query().Get("destination")
	amountStr := r.URL.Query().Get("amount")

	if origin == "" || destination == "" || amountStr == "" {
		JSONError(w, http.StatusBadRequest, "Missing required parameters: origin, destination, amount", "MISSING_PARAMETERS")
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		JSONError(w, http.StatusBadRequest, "Invalid amount parameter", "INVALID_AMOUNT")
		return
	}

	// Get exchange rate
	rateResponse, err := h.awsServices.CurrencyService.GetExchangeRateGivenAmount(r.Context(), origin, destination, amount)
	if err != nil {
		JSONError(w, http.StatusUnprocessableEntity, "Unable to get exchange rate", "RATE_UNAVAILABLE")
		return
	}

	// Get currency information
	originCurrency, err := h.awsServices.CurrencyService.GetCurrencyInfo(r.Context(), origin)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid origin currency", "INVALID_ORIGIN")
		return
	}

	destCurrency, err := h.awsServices.CurrencyService.GetCurrencyInfo(r.Context(), destination)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid destination currency", "INVALID_DESTINATION")
		return
	}

	response := domain.ConversionResponse{
		Origin:          *originCurrency,
		Destination:     *destCurrency,
		Rate:            rateResponse.ConversionRate,
		Amount:          amount,
		ConvertedAmount: rateResponse.ConversionResult,
		Timestamp:       time.Now().UTC(),
		RatesSource:     rateResponse.RatesSource,
	}

	JSONResponse(w, http.StatusOK, response)
}

// History handles historical exchange rate requests
// GET /api/v1/history?origin={ORIGIN}&destination={DEST}&start_date={YYYY-MM-DD}&end_date={YYYY-MM-DD}
func (h *CurrencyHandler) History(w http.ResponseWriter, r *http.Request) {
	origin := r.URL.Query().Get("origin")
	destination := r.URL.Query().Get("destination")
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	if origin == "" || destination == "" || startDateStr == "" || endDateStr == "" {
		JSONError(w, http.StatusBadRequest, "Missing required parameters: origin, destination, start_date, end_date", "MISSING_PARAMETERS")
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid start_date format. Use YYYY-MM-DD", "INVALID_DATE_FORMAT")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid end_date format. Use YYYY-MM-DD", "INVALID_DATE_FORMAT")
		return
	}

	// Get currency information
	originCurrency, err := h.awsServices.CurrencyService.GetCurrencyInfo(r.Context(), origin)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid origin currency", "INVALID_ORIGIN")
		return
	}

	destCurrency, err := h.awsServices.CurrencyService.GetCurrencyInfo(r.Context(), destination)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid destination currency", "INVALID_DESTINATION")
		return
	}

	// Get historical rates
	rates, source, err := h.awsServices.CurrencyService.GetHistoricalRates(r.Context(), origin, destination, startDate, endDate)
	if err != nil {
		JSONError(w, http.StatusUnprocessableEntity, fmt.Sprintf("No historical data available %s", err.Error()), "NO_DATA_AVAILABLE")
		return
	}
	if origin != "EUR" {
		originCurrency, _ = h.awsServices.CurrencyService.GetCurrencyInfo(r.Context(), "EUR")
	}

	response := domain.HistoryResponse{
		Origin:      *originCurrency,
		Destination: *destCurrency,
		StartDate:   startDateStr,
		EndDate:     endDateStr,
		Rates:       rates,
		Timestamp:   time.Now().UTC(),
		RatesSource: source,
		Message: "I'm sorry, for now the origin must alway be the EUR code",
	}

	JSONResponse(w, http.StatusOK, response)
}

// Forecast handles forecast requests
// GET /api/v1/forecast?origin={ORIGIN}&destination={DEST}
func (h *CurrencyHandler) Forecast(w http.ResponseWriter, r *http.Request) {
	origin := r.URL.Query().Get("origin")
	destination := r.URL.Query().Get("destination")

	if origin == "" || destination == "" {
		JSONError(w, http.StatusBadRequest, "Missing required parameters: origin, destination", "MISSING_PARAMETERS")
		return
	}

	forecast, err := h.awsServices.CurrencyService.GetForecast(r.Context(), origin, destination)
	if err != nil {
		JSONError(w, http.StatusUnprocessableEntity, "Insufficient data for forecast", "INSUFFICIENT_DATA")
		return
	}

	JSONResponse(w, http.StatusOK, forecast)
}

// GetDestinations handles available destinations requests
// GET /api/v1/origins/{origin}/destinations
func (h *CurrencyHandler) GetDestinations(w http.ResponseWriter, r *http.Request) {
	origin := chi.URLParam(r, "origin")

	if origin == "" {
		JSONError(w, http.StatusBadRequest, "Missing origin parameter", "MISSING_ORIGIN")
		return
	}

	// Get currency information
	originCurrency, err := h.awsServices.CurrencyService.GetCurrencyInfo(r.Context(), origin)
	if err != nil {
		JSONError(w, http.StatusNotFound, "Origin not supported", "ORIGIN_NOT_SUPPORTED")
		return
	}

	// Get supported destinations
	destinations, source, err := h.awsServices.CurrencyService.GetSupportedDestinations(r.Context(), origin)
	if err != nil {
		JSONError(w, http.StatusNotFound, "No destinations available, for now, only EUR code is supported", "NO_DESTINATIONS")
		return
	}

	response := domain.DestinationsResponse{
		Origin:       *originCurrency,
		Destinations: destinations,
		Timestamp:    time.Now().UTC(),
		RatesSource:  source,
	}

	JSONResponse(w, http.StatusOK, response)
}

// SaveFavorite handles saving favorite conversions
// POST /api/v1/favorites
func (h *CurrencyHandler) SaveFavorite(w http.ResponseWriter, r *http.Request) {
	var req domain.FavoriteRequest
	if err := BindJSON(r, &req); err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST")
		return
	}

	favorite, err := h.awsServices.FavoriteService.SaveFavorite(r.Context(), &req)
	if err != nil {
		JSONError(w, http.StatusConflict, "Favorite already exists", "FAVORITE_EXISTS")
		return
	}

	JSONResponse(w, http.StatusCreated, favorite)
}

// CheckFavorites handles daily favorite checks
// POST /api/v1/favorites/check
func (h *CurrencyHandler) CheckFavorites(w http.ResponseWriter, r *http.Request) {
	results, err := h.awsServices.FavoriteService.CheckFavorites(r.Context())
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to check favorites", "CHECK_FAILED")
		return
	}

	JSONResponse(w, http.StatusOK, results)
}

// SendNotification handles email notifications
// POST /api/v1/notifications/email
func (h *CurrencyHandler) SendNotification(w http.ResponseWriter, r *http.Request) {
	var req domain.NotificationRequest
	if err := BindJSON(r, &req); err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST")
		return
	}

	response, err := h.awsServices.NotificationService.SendEmailNotification(r.Context(), &req)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to send email", "EMAIL_FAILED")
		return
	}

	JSONResponse(w, http.StatusOK, response)
}
