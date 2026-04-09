package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/constants"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/services"
)

type PriceHistoryController struct {
	service *services.PriceHistoryService
}

func NewPriceHistoryController(service *services.PriceHistoryService) *PriceHistoryController {
	return &PriceHistoryController{service: service}
}

// GetPriceHistory retrieves price history for a stock
func (c *PriceHistoryController) GetPriceHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stockIDStr := vars["stockID"]

	stockID, err := strconv.Atoi(stockIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidStockIDFormat, err)
		return
	}

	if stockID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgStockIDMustBePositive, nil)
		return
	}

	priceHistory, err := c.service.GetPriceHistoryByStockID(stockID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, "No price history found for this stock", nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve price history", err)
		return
	}

	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: priceHistory})
}

// GetPriceHistoryByDateRange retrieves price history for a stock within a date range
func (c *PriceHistoryController) GetPriceHistoryByDateRange(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stockIDStr := vars["stockID"]

	stockID, err := strconv.Atoi(stockIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidStockIDFormat, err)
		return
	}

	if stockID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgStockIDMustBePositive, nil)
		return
	}

	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	if startDate == "" || endDate == "" {
		writeErrorResponse(w, http.StatusBadRequest, "start_date and end_date query parameters are required", nil)
		return
	}

	priceHistory, err := c.service.GetPriceHistoryByDateRange(stockID, startDate, endDate)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, "No price history found for the given date range", nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve price history", err)
		return
	}

	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: priceHistory})
}

// RecordPrice records a new price for a stock
func (c *PriceHistoryController) RecordPrice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stockIDStr := vars["stockID"]

	stockID, err := strconv.Atoi(stockIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidStockIDFormat, err)
		return
	}

	if stockID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgStockIDMustBePositive, nil)
		return
	}

	var payload struct {
		Price float64 `json:"price"`
		Date  string  `json:"date"`
	}

	if err := parseJSONBody(r, &payload); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if payload.Price < 0 {
		writeErrorResponse(w, http.StatusBadRequest, "Price must be non-negative", nil)
		return
	}

	// Parse date if provided, otherwise use current time
	var priceDate time.Time
	if payload.Date != "" {
		priceDate, err = time.Parse("2006-01-02 15:04:05", payload.Date)
		if err != nil {
			writeErrorResponse(w, http.StatusBadRequest, "Invalid date format (use YYYY-MM-DD HH:MM:SS)", nil)
			return
		}
	} else {
		priceDate = time.Now()
	}

	priceHistory, err := c.service.RecordPrice(stockID, payload.Price, priceDate)
	if err != nil {
		if errors.Is(err, db.ErrInvalidStockID) || errors.Is(err, db.ErrInvalidPrice) {
			writeErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to record price", err)
		return
	}

	writeJSONResponse(w, http.StatusCreated, SuccessResponse{Data: priceHistory})
}

// GetPriceStats retrieves price statistics for a stock
func (c *PriceHistoryController) GetPriceStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stockIDStr := vars["stockID"]

	stockID, err := strconv.Atoi(stockIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidStockIDFormat, err)
		return
	}

	if stockID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgStockIDMustBePositive, nil)
		return
	}

	stats, err := c.service.CalculatePriceStats(stockID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, "No price history found for this stock", nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to calculate statistics", err)
		return
	}

	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: stats})
}

// GetLatestPrice retrieves the latest price for a stock
func (c *PriceHistoryController) GetLatestPrice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stockIDStr := vars["stockID"]

	stockID, err := strconv.Atoi(stockIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidStockIDFormat, err)
		return
	}

	if stockID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgStockIDMustBePositive, nil)
		return
	}

	priceHistory, err := c.service.GetLatestPrice(stockID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, "No price history found for this stock", nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve latest price", err)
		return
	}

	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: priceHistory})
}
