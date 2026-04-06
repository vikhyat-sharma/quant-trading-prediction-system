package controllers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
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
		writeErrorResponse(w, http.StatusBadRequest, "Invalid stock ID format", err)
		return
	}

	if stockID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, "Stock ID must be a positive integer", nil)
		return
	}

	predictions, err := c.service.GetPredictionsByStockID(stockID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, "No predictions found for this stock", nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve predictions", err)
		return
	}

	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: predictions})
}

func (c *PredictionController) GeneratePrediction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stockIDStr := vars["stockID"]

	stockID, err := strconv.Atoi(stockIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid stock ID format", err)
		return
	}

	if stockID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, "Stock ID must be a positive integer", nil)
		return
	}

	prediction, err := c.service.GeneratePrediction(stockID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, "Stock not found", nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to generate prediction", err)
		return
	}

	writeJSONResponse(w, http.StatusCreated, SuccessResponse{Data: prediction})
}
