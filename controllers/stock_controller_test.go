package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/services"
)

// MockStockService is a mock implementation for testing
type MockStockService struct {
	stock       *db.Stock
	stocks      []*db.Stock
	err         error
	shouldError bool
}

func (m *MockStockService) GetStock(id int) (*db.Stock, error) {
	if m.shouldError {
		return nil, m.err
	}
	return m.stock, nil
}

func (m *MockStockService) GetAllStocks() ([]*db.Stock, error) {
	if m.shouldError {
		return nil, m.err
	}
	return m.stocks, nil
}

func TestStockController_GetStock_InvalidID(t *testing.T) {
	controller := &StockController{service: &services.StockService{}}

	req := httptest.NewRequest("GET", "/stocks/invalid", nil)
	w := httptest.NewRecorder()

	vars := map[string]string{"id": "invalid"}
	req = mux.SetURLVars(req, vars)

	controller.GetStock(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}
