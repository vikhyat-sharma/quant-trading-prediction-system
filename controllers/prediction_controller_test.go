package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/services"
)

func TestPredictionController_GetPredictions_InvalidStockID(t *testing.T) {
	controller := &PredictionController{service: &services.PredictionService{}}

	req := httptest.NewRequest("GET", "/stocks/invalid/predictions", nil)
	w := httptest.NewRecorder()

	vars := map[string]string{"stockID": "invalid"}
	req = mux.SetURLVars(req, vars)

	controller.GetPredictions(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestPredictionController_GeneratePrediction_InvalidStockID(t *testing.T) {
	controller := &PredictionController{service: &services.PredictionService{}}

	req := httptest.NewRequest("POST", "/stocks/invalid/predictions/generate", nil)
	w := httptest.NewRecorder()

	vars := map[string]string{"stockID": "invalid"}
	req = mux.SetURLVars(req, vars)

	controller.GeneratePrediction(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}
