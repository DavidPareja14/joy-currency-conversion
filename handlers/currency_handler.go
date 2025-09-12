package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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
func (h *CurrencyHandler) Convert(c *gin.Context) {
	origin := c.Query("origin")
	destination := c.Query("destination")
	amountStr := c.Query("amount")

	if origin == "" || destination == "" || amountStr == "" {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error: "Missing required parameters: origin, destination, amount",
			Code:  "MISSING_PARAMETERS",
		})
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error: "Invalid amount parameter",
			Code:  "INVALID_AMOUNT",
		})
		return
	}

	// Get exchange rate
	rate, source, err := h.awsServices.CurrencyService.GetExchangeRate(c.Request.Context(), origin, destination)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, domain.ErrorResponse{
			Error: "Unable to get exchange rate",
			Code:  "RATE_UNAVAILABLE",
		})
		return
	}

	// Get currency information
	originCurrency, err := h.awsServices.CurrencyService.GetCurrencyInfo(c.Request.Context(), origin)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error: "Invalid origin currency",
			Code:  "INVALID_ORIGIN",
		})
		return
	}

	destCurrency, err := h.awsServices.CurrencyService.GetCurrencyInfo(c.Request.Context(), destination)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error: "Invalid destination currency",
			Code:  "INVALID_DESTINATION",
		})
		return
	}

	response := domain.ConversionResponse{
		Origin:          *originCurrency,
		Destination:     *destCurrency,
		Rate:            rate,
		Amount:          amount,
		ConvertedAmount: amount * rate,
		Timestamp:       time.Now().UTC(),
		RatesSource:     source,
	}

	c.JSON(http.StatusOK, response)
}

// History handles historical exchange rate requests
// GET /api/v1/history?origin={ORIGIN}&destination={DEST}&start_date={YYYY-MM-DD}&end_date={YYYY-MM-DD}
func (h *CurrencyHandler) History(c *gin.Context) {
	origin := c.Query("origin")
	destination := c.Query("destination")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if origin == "" || destination == "" || startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error: "Missing required parameters: origin, destination, start_date, end_date",
			Code:  "MISSING_PARAMETERS",
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error: "Invalid start_date format. Use YYYY-MM-DD",
			Code:  "INVALID_DATE_FORMAT",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error: "Invalid end_date format. Use YYYY-MM-DD",
			Code:  "INVALID_DATE_FORMAT",
		})
		return
	}

	// Get currency information
	originCurrency, err := h.awsServices.CurrencyService.GetCurrencyInfo(c.Request.Context(), origin)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error: "Invalid origin currency",
			Code:  "INVALID_ORIGIN",
		})
		return
	}

	destCurrency, err := h.awsServices.CurrencyService.GetCurrencyInfo(c.Request.Context(), destination)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error: "Invalid destination currency",
			Code:  "INVALID_DESTINATION",
		})
		return
	}

	// Get historical rates
	rates, source, err := h.awsServices.CurrencyService.GetHistoricalRates(c.Request.Context(), origin, destination, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, domain.ErrorResponse{
			Error: "No historical data available",
			Code:  "NO_DATA_AVAILABLE",
		})
		return
	}

	response := domain.HistoryResponse{
		Origin:      *originCurrency,
		Destination: *destCurrency,
		StartDate:   startDateStr,
		EndDate:     endDateStr,
		Rates:       rates,
		Timestamp:   time.Now().UTC(),
		RatesSource: source,
	}

	c.JSON(http.StatusOK, response)
}

// Forecast handles forecast requests
// GET /api/v1/forecast?origin={ORIGIN}&destination={DEST}
func (h *CurrencyHandler) Forecast(c *gin.Context) {
	origin := c.Query("origin")
	destination := c.Query("destination")

	if origin == "" || destination == "" {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error: "Missing required parameters: origin, destination",
			Code:  "MISSING_PARAMETERS",
		})
		return
	}

	forecast, err := h.awsServices.CurrencyService.GetForecast(c.Request.Context(), origin, destination)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, domain.ErrorResponse{
			Error: "Insufficient data for forecast",
			Code:  "INSUFFICIENT_DATA",
		})
		return
	}

	c.JSON(http.StatusOK, forecast)
}

// GetDestinations handles available destinations requests
// GET /api/v1/origins/{origin}/destinations
func (h *CurrencyHandler) GetDestinations(c *gin.Context) {
	origin := c.Param("origin")

	if origin == "" {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error: "Missing origin parameter",
			Code:  "MISSING_ORIGIN",
		})
		return
	}

	// Get currency information
	originCurrency, err := h.awsServices.CurrencyService.GetCurrencyInfo(c.Request.Context(), origin)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.ErrorResponse{
			Error: "Origin not supported",
			Code:  "ORIGIN_NOT_SUPPORTED",
		})
		return
	}

	// Get supported destinations
	destinations, source, err := h.awsServices.CurrencyService.GetSupportedDestinations(c.Request.Context(), origin)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.ErrorResponse{
			Error: "No destinations available",
			Code:  "NO_DESTINATIONS",
		})
		return
	}

	response := domain.DestinationsResponse{
		Origin:      *originCurrency,
		Destinations: destinations,
		Timestamp:   time.Now().UTC(),
		RatesSource: source,
	}

	c.JSON(http.StatusOK, response)
}

// SaveFavorite handles saving favorite conversions
// POST /api/v1/favorites
func (h *CurrencyHandler) SaveFavorite(c *gin.Context) {
	var req domain.FavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error: "Invalid request body",
			Code:  "INVALID_REQUEST",
		})
		return
	}

	favorite, err := h.awsServices.FavoriteService.SaveFavorite(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusConflict, domain.ErrorResponse{
			Error: "Favorite already exists",
			Code:  "FAVORITE_EXISTS",
		})
		return
	}

	c.JSON(http.StatusCreated, favorite)
}

// CheckFavorites handles daily favorite checks
// POST /api/v1/favorites/check
func (h *CurrencyHandler) CheckFavorites(c *gin.Context) {
	results, err := h.awsServices.FavoriteService.CheckFavorites(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Error: "Failed to check favorites",
			Code:  "CHECK_FAILED",
		})
		return
	}

	c.JSON(http.StatusOK, results)
}

// SendNotification handles email notifications
// POST /api/v1/notifications/email
func (h *CurrencyHandler) SendNotification(c *gin.Context) {
	var req domain.NotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error: "Invalid request body",
			Code:  "INVALID_REQUEST",
		})
		return
	}

	response, err := h.awsServices.NotificationService.SendEmailNotification(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Error: "Failed to send email",
			Code:  "EMAIL_FAILED",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
