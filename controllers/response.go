package controllers

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Data interface{} `json:"data"`
}

// writeJSONResponse writes a JSON response with proper headers
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// writeErrorResponse writes an error response
func writeErrorResponse(w http.ResponseWriter, statusCode int, message string, err error) {
	response := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
		Code:    statusCode,
	}

	if err != nil && message == "" {
		response.Message = err.Error()
	}

	writeJSONResponse(w, statusCode, response)
}
