package response

type ExchangeRateResponse struct {
	ConversionRate float64 `json:"conversion_rate"`
	ConversionResult float64 `json:"conversion_result"`
	RatesSource string
}