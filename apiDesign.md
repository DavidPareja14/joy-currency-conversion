# Project Joy — API Design (detailed, v1.2)

This document is the detailed API design for Project Joy, with all 7 endpoints fully specified. It keeps the previous detailed format and applies the requested updates (date format, query param names, removed slope, etc.).

---

## Common conventions
- Base path: `/api/v1`.
- Currency codes: ISO 4217 (three uppercase letters, `^[A-Z]{3}$`).
- Dates: strict `yyyy-mm-dd`, parse in Go with `time.Parse("2006-01-02", s)`.
- Timestamps in responses: ISO 8601 UTC (`YYYY-MM-DDTHH:MM:SSZ`).
- Error format (JSON): `{ "error": "message", "code": "ERROR_CODE", "details": {...} }`.
- Field `rates_source`: included in rate-related responses to identify the external provider used for the data.

---

## Endpoint 1: Currency Conversion
**GET** `/api/v1/convert?origin={ORIGIN}&destination={DEST}&amount={AMOUNT}`

Convert an amount from one currency to another.

### Parameters
- `origin` (query, required): origin currency code (ISO 4217).
- `destination` (query, required): destination currency code (ISO 4217).
- `amount` (query, required): amount to convert (positive number).

### Responses
**200 OK**
```json
{
  "origin": {"code": "COP", "country": "Colombia"},
  "destination": {"code": "USD", "country": "United States"},
  "rate": 0.00025,
  "amount": 100000,
  "converted_amount": 25,
  "timestamp": "2025-09-07T12:00:00Z",
  "rates_source": "example-provider"
}
```
- `rates_source`: provider that supplied the exchange rate for this conversion.

**400 Bad Request** — invalid parameters.
**422 Unprocessable Entity** — no rate available.

---

## Endpoint 2: Daily Historic Values
**GET** `/api/v1/history?origin={ORIGIN}&destination={DEST}&start_date={YYYY-MM-DD}&end_date={YYYY-MM-DD}`

Retrieve historical exchange rates for a currency pair in a date range.

### Parameters
- `origin` (query, required): origin currency code.
- `destination` (query, required): destination currency code.
- `start_date` (query, required): start date (`yyyy-mm-dd`).
- `end_date` (query, required): end date (`yyyy-mm-dd`).

### Responses
**200 OK**
```json
{
  "origin": {"code": "COP", "country": "Colombia"},
  "destination": {"code": "USD", "country": "United States"},
  "start_date": "2025-09-01",
  "end_date": "2025-09-30",
  "rates": [
    {"date": "2025-09-01", "rate": 0.00024},
    {"date": "2025-09-02", "rate": 0.00025},
    {"date": "2025-09-03", "rate": 0.00028}
  ],
  "timestamp": "2025-09-30T12:00:00Z",
  "rates_source": "example-provider"
}
```
- `rates_source`: provider that supplied the historical exchange rates.

**400 Bad Request** — invalid date format.
**422 Unprocessable Entity** — no data available.

---

## Endpoint 3: Probability Forecast (Basic)
**GET** `/api/v1/forecast?origin={ORIGIN}&destination={DEST}`

Forecast the next day’s exchange rate for a currency pair based on the last 30 days.

### Parameters
- `origin` (query, required): origin currency code.
- `destination` (query, required): destination currency code.

### Responses
**200 OK**
```json
{
  "origin": {"code": "COP", "country": "Colombia"},
  "destination": {"code": "USD", "country": "United States"},
  "predicted_date": "2025-09-08",
  "predicted_rate": 0.00026,
  "confidence": 0.65,
  "last_30_days": {
    "average": 0.000255
  },
  "timestamp": "2025-09-07T12:00:00Z",
  "rates_source": "example-provider"
}
```
- `rates_source`: provider that supplied the rates used for the forecast.

**400 Bad Request** — invalid parameters.
**422 Unprocessable Entity** — insufficient data.

---

## Endpoint 4: Available Destination Currencies
**GET** `/api/v1/origins/{ORIGIN}/destinations`

List supported destination currencies for an origin.

### Parameters
- `origin` (path, required): origin currency code.

### Responses
**200 OK**
```json
{
  "origin": {"code": "COP", "country": "Colombia"},
  "destinations": [
    {"code": "USD", "country": "United States"},
    {"code": "EUR", "country": "Eurozone"}
  ],
  "timestamp": "2025-09-07T12:00:00Z",
  "rates_source": "example-provider"
}
```
- `rates_source`: provider that supplied the supported destinations.

**400 Bad Request** — invalid origin.
**404 Not Found** — origin not supported.

---

## Endpoint 5: Save a Favorite Conversion
**POST** `/api/v1/favorites`

Save an origin/destination pair with a threshold and notify email.

### Request Body
```json
{
  "origin": "COP",
  "destination": "USD",
  "threshold": 0.00030,
  "notify_email": "dendrite@example.com"
}
```

### Responses
**201 Created**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "origin": {"code": "COP", "country": "Colombia"},
  "destination": {"code": "USD", "country": "United States"},
  "threshold": 0.00030,
  "notify_email": "dendrite@example.com",
  "created_at": "2025-09-07T12:00:00Z"
}
```

**400 Bad Request** — invalid input.
**409 Conflict** — favorite already exists.

---

## Endpoint 6: Daily Favorite Check
**POST** `/api/v1/favorites/check`

Run the daily check for all saved favorites. Intended to be called by a scheduler.

### Responses
**200 OK**
```json
{
  "results": [
    {
      "favorite_id": "550e8400-e29b-41d4-a716-446655440000",
      "origin": {"code": "COP", "country": "Colombia"},
      "destination": {"code": "USD", "country": "United States"},
      "threshold": 0.00030,
      "current_rate": 0.00031,
      "date": "2025-09-07",
      "exceeded": true,
      "notified": true,
      "current_rate_source": "example-provider"
    }
  ],
  "timestamp": "2025-09-07T12:00:00Z"
}
```

**500 Internal Server Error** — failure during checks.

---

## Endpoint 7: Email Notification on Threshold Exceeded
**POST** `/api/v1/notifications/email`

Send an email notification when a favorite threshold is exceeded.

### Request Body
```json
{
  "favorite_id": "550e8400-e29b-41d4-a716-446655440000",
  "origin": {"code": "COP", "country": "Colombia"},
  "destination": {"code": "USD", "country": "United States"},
  "threshold": 0.00030,
  "current_rate": 0.00031,
  "date": "2025-09-07",
  "notify_email": "dendrite@example.com"
}
```

### Responses
**200 OK**
```json
{
  "message": "Email queued/sent",
  "sent_to": "dendrite@example.com"
}
```

**400 Bad Request** — invalid input.
**500 Internal Server Error** — failure to send email.

---

