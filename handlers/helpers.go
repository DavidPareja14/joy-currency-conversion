package handlers

import (
	"encoding/json"
	"net/http"
)

// JSONResponse writes a JSON response to the ResponseWriter
func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// JSONError writes a JSON error response to the ResponseWriter
func JSONError(w http.ResponseWriter, statusCode int, message string, code string) {
	errorResponse := map[string]string{
		"error": message,
		"code":  code,
	}
	JSONResponse(w, statusCode, errorResponse)
}

// BindJSON decodes JSON from request body into the provided interface
func BindJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
