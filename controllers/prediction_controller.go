package controllers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/constants"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/repositories"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/services"
)

type PredictionController struct {
	service *services.PredictionService
}

func NewPredictionController(service *services.PredictionService) *PredictionController {
	return &PredictionController{service: service}
}

func (c *PredictionController) GetPredictions(w http.ResponseWriter, r *http.Request) {
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

	// Check for filter query parameters
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")
	minPriceStr := r.URL.Query().Get("min_price")
	maxPriceStr := r.URL.Query().Get("max_price")

	// If no filters provided, get all predictions for the stock
	if startDateStr == "" && endDateStr == "" && minPriceStr == "" && maxPriceStr == "" {
		predictions, err := c.service.GetPredictionsByStockID(stockID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) || errors.Is(err, db.ErrRecordNotFound) {
				writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgNoPredictionsForStock, nil)
				return
			}
			writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToRetrievePredictions, err)
			return
		}
		writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: predictions})
		return
	}

	// Build filter from query parameters
	filter := &repositories.PredictionFilter{
		StockID: stockID,
	}

	if startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			writeErrorResponse(w, http.StatusBadRequest, "Invalid start_date format (use YYYY-MM-DD)", err)
			return
		}
		filter.StartDate = startDate
	}

	if endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			writeErrorResponse(w, http.StatusBadRequest, "Invalid end_date format (use YYYY-MM-DD)", err)
			return
		}
		filter.EndDate = endDate.Add(time.Hour * 24)
	}

	if minPriceStr != "" {
		minPrice, err := strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			writeErrorResponse(w, http.StatusBadRequest, "Invalid min_price format", err)
			return
		}
		filter.MinPrice = minPrice
	}

	if maxPriceStr != "" {
		maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			writeErrorResponse(w, http.StatusBadRequest, "Invalid max_price format", err)
			return
		}
		filter.MaxPrice = maxPrice
	}

	predictions, err := c.service.SearchAndFilterPredictions(filter)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToRetrievePredictions, err)
		return
	}

	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: predictions})
}

func (c *PredictionController) GeneratePrediction(w http.ResponseWriter, r *http.Request) {
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

	prediction, err := c.service.GeneratePrediction(stockID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgStockNotFound, nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToGeneratePrediction, err)
		return
	}

	writeJSONResponse(w, http.StatusCreated, SuccessResponse{Data: prediction})
}
