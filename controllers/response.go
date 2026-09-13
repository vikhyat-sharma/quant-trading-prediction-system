package controllers

import (
	"encoding/json"
	"io"
	"net/http"
)

// apiError is the canonical error envelope returned to clients.
// Internal error details are never included in the response body.
type apiError struct {
	Error apiErrorBody `json:"error"`
}

type apiErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// SuccessResponse wraps successful response data.
type SuccessResponse struct {
	Data interface{} `json:"data"`
}

// writeJSONResponse writes a JSON response with the given status code.
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}

// writeErrorResponse writes a structured error response.
// The internal err is intentionally not forwarded to the client.
func writeErrorResponse(w http.ResponseWriter, statusCode int, message string, _ error) {
	code := httpStatusToCode(statusCode)
	writeJSONResponse(w, statusCode, apiError{
		Error: apiErrorBody{Code: code, Message: message},
	})
}

// parseJSONBody decodes the request body into v.
// Returns an error if the body is missing, malformed, or too large.
func parseJSONBody(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil && err != io.EOF {
		return err
	}
	return nil
}

func httpStatusToCode(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "VALIDATION_ERROR"
	case http.StatusUnauthorized:
		return "UNAUTHORIZED"
	case http.StatusForbidden:
		return "FORBIDDEN"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusConflict:
		return "CONFLICT"
	case http.StatusTooManyRequests:
		return "RATE_LIMITED"
	case http.StatusInternalServerError:
		return "INTERNAL_ERROR"
	default:
		return "ERROR"
	}
}
