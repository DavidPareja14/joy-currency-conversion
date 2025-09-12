# Project Joy — API Design (clean rebuild, updated)

This document is a clean, updated rebuild of the Project Joy API design. It preserves endpoints 1–7 and applies your requested changes.

## Endpoints (summary)

1. `GET /api/v1/convert?origin={ORIGIN}&destination={DEST}&amount={AMOUNT}` — Currency conversion.
2. `GET /api/v1/history?origin={ORIGIN}&destination={DEST}&start_date={YYYY-MM-DD}&end_date={YYYY-MM-DD}` — Daily historic values.
3. `GET /api/v1/forecast?origin={ORIGIN}&destination={DEST}` — Forecast for tomorrow.
4. `GET /api/v1/origins/{ORIGIN}/destinations` — Supported destinations for an origin.
5. `POST /api/v1/favorites` — Save a favorite conversion (origin, destination, threshold, notify_email).
6. `POST /api/v1/favorites/check` — Run the daily check for all favorites (intended for scheduler).
7. `POST /api/v1/notifications/email` — Send an email notification (internal / test use).

---

## Common conventions
- Base path: `/api/v1`.
- Currency codes: ISO 4217 (three uppercase letters, `^[A-Z]{3}$`).
- Dates: strict `yyyy-mm-dd`, parse in Go with `time.Parse("2006-01-02", s)`.
- Timestamps in responses: ISO 8601 UTC (`YYYY-MM-DDTHH:MM:SSZ`).
- Error format (JSON): `{ "error": "message", "code": "ERROR_CODE", "details": {...} }`.
- Field `rates_source` is included in rate-related responses to identify the external provider (brief explanation added in each endpoint schema).

---

(Full endpoint reference, request/response shapes, validation rules and examples are included in the accompanying OpenAPI YAML and the embedded Swagger UI.)

