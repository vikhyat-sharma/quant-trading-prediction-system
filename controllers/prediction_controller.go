package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
		http.Error(w, "Invalid Stock ID", http.StatusBadRequest)
		return
	}
	predictions, err := c.service.GetPredictionsByStockID(stockID)
	if err != nil {
		http.Error(w, "Error fetching predictions", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(predictions)
}

func (c *PredictionController) GeneratePrediction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stockIDStr := vars["stockID"]
	stockID, err := strconv.Atoi(stockIDStr)
	if err != nil {
		http.Error(w, "Invalid Stock ID", http.StatusBadRequest)
		return
	}
	prediction, err := c.service.GeneratePrediction(stockID)
	if err != nil {
		http.Error(w, "Error generating prediction", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(prediction)
}
