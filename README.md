# joy-currency-conversion-private
 This repository contains the API design and artifacts for **Project Joy**.

## Contents
- `joy_openapi_detailed_v1.2.yaml` — Full OpenAPI 3.0 specification (YAML).
- `project_joy_postman_detailed_v1.2.json` — Postman collection for quick testing.
- `project_joy_swagger_ui_detailed_v1.2.html` — Standalone Swagger UI (open in a browser).

## How to Use

### 1. Swagger UI (Interactive Docs)
1. Download `project_joy_swagger_ui_detailed_v1.2.html`.
2. Open it in your browser (double click or drag+drop).
3. You’ll see a full interactive Swagger UI with all endpoints documented.
4. Change the `baseUrl` server URL if needed to point to your API deployment.

### 2. OpenAPI YAML
- Import `joy_openapi_detailed_v1.2.yaml` into tools like **Stoplight**, **Insomnia**, or **Postman**.
- Use it as the single source of truth for schemas, validation, and code generation.

### 3. Postman Collection
1. Open **Postman**.
2. Go to **File → Import**.
3. Select `project_joy_postman_detailed_v1.2.json`.
4. The collection will be added, with pre-configured requests for all 7 endpoints.
5. Set the collection variable `baseUrl` to your API URL (default is `http://localhost:8080`).

### 4. Run locally (Go example)
If you scaffold the API server in Go, run it on `http://localhost:8080`.
- Then you can test requests using Postman or Swagger UI directly.

## Notes
- All dates use strict format `yyyy-mm-dd`.
- `rates_source` fields indicate which external provider was used for exchange rates data.
- Error responses follow the JSON format: 
  ```json
  { "error": "message", "code": "ERROR_CODE", "details": {...} }
