package controllers

import (
	"net/http"
	"strconv"
)

const (
	defaultPageSize = 50
	maxPageSize     = 500
)

// Pagination holds validated pagination parameters.
type Pagination struct {
	Page     int
	PageSize int
	Offset   int
}

// parsePagination extracts and validates ?page and ?page_size query parameters.
// Returns safe defaults if parameters are absent or invalid.
func parsePagination(r *http.Request) Pagination {
	page := 1
	pageSize := defaultPageSize

	if v := r.URL.Query().Get("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			page = n
		}
	}
	if v := r.URL.Query().Get("page_size"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			if n > maxPageSize {
				n = maxPageSize
			}
			pageSize = n
		}
	}

	return Pagination{
		Page:     page,
		PageSize: pageSize,
		Offset:   (page - 1) * pageSize,
	}
}

// PaginatedResponse wraps a paginated collection response.
type PaginatedResponse struct {
	Data     interface{} `json:"data"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	Total    int         `json:"total,omitempty"`
}
