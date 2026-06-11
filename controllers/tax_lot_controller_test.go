package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
)

// mock service implementing necessary methods
type mockTaxLotService struct{}

func (m *mockTaxLotService) RecordBuy(portfolioID, stockID int, quantity, price, fees float64, buyDate time.Time) (*db.TaxLot, error) {
	tl := &db.TaxLot{ID: 1, PortfolioID: portfolioID, StockID: stockID, Quantity: quantity, CostPerShare: price, TotalCost: quantity * price, AcquisitionDate: buyDate}
	return tl, nil
}
func (m *mockTaxLotService) RecordSellFIFO(portfolioID, stockID int, quantity, price, fees float64, sellDate time.Time) (float64, error) {
	return 123.45, nil
}
func (m *mockTaxLotService) RecordSellLIFO(portfolioID, stockID int, quantity, price, fees float64, sellDate time.Time) (float64, error) {
	return 67.89, nil
}
func (m *mockTaxLotService) RecordSellSpecificLot(taxLotID int, quantity, price, fees float64, sellDate time.Time) (float64, error) {
	return 11.22, nil
}
func (m *mockTaxLotService) GetTaxLotGains(taxLotID int, currentPrice float64) (*db.TaxLotGains, error) {
	return &db.TaxLotGains{TaxLotID: taxLotID, StockID: 1, Symbol: "FOO", QuantityHeld: 10, QuantitySold: 5, CostPerShare: 10.0, CurrentPrice: currentPrice, CostBasis: 150.0, CurrentValue: 200.0}, nil
}
func (m *mockTaxLotService) GetPortfolioTaxGains(portfolioID int, currentPrices map[int]float64) (map[string]interface{}, error) {
	return map[string]interface{}{"total_realized_gain": 100.0}, nil
}
func (m *mockTaxLotService) CalculateTaxableGainsBySellDate(portfolioID int) (map[string]interface{}, error) {
	return map[string]interface{}{"short_term_gains": 10.0, "long_term_gains": 20.0}, nil
}
func (m *mockTaxLotService) GetTaxTransactionsByPortfolio(portfolioID int) ([]db.TaxTransaction, error) {
	return []db.TaxTransaction{{ID: 1, PortfolioID: portfolioID, StockID: 1, Type: "BUY", Quantity: 100, Price: 10}}, nil
}

func setupRouterWithMock() *mux.Router {
	svc := &mockTaxLotService{}
	ctrl := NewTaxLotController(svc)
	r := mux.NewRouter()
	r.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-lots/buy", ctrl.RecordBuy).Methods("POST")
	r.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-lots/sell-fifo", ctrl.RecordSellFIFO).Methods("POST")
	r.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-lots/sell-lifo", ctrl.RecordSellLIFO).Methods("POST")
	r.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-lots/{taxLotID}/sell", ctrl.RecordSellSpecificLot).Methods("POST")
	r.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-lots/{taxLotID}/gains", ctrl.GetTaxLotGains).Methods("GET")
	r.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-gains", ctrl.GetPortfolioTaxGains).Methods("POST")
	r.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-report", ctrl.GetTaxableGains).Methods("GET")
	r.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-transactions", ctrl.GetTaxTransactions).Methods("GET")
	return r
}

func TestRecordBuy(t *testing.T) {
	r := setupRouterWithMock()
	body := map[string]interface{}{"stock_id": 1, "quantity": 10, "price": 5.0, "fees": 0.5}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/users/1/portfolios/1/tax-lots/buy", bytes.NewReader(b))
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d: %s", rw.Code, rw.Body.String())
	}
}

func TestRecordSellFIFO(t *testing.T) {
	r := setupRouterWithMock()
	body := map[string]interface{}{"stock_id": 1, "quantity": 5, "price": 6.0, "fees": 0.2}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/users/1/portfolios/1/tax-lots/sell-fifo", bytes.NewReader(b))
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d: %s", rw.Code, rw.Body.String())
	}
}

func TestRecordSellSpecific(t *testing.T) {
	r := setupRouterWithMock()
	body := map[string]interface{}{"quantity": 2, "price": 6.0, "fees": 0.1}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/users/1/portfolios/1/tax-lots/1/sell", bytes.NewReader(b))
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d: %s", rw.Code, rw.Body.String())
	}
}

func TestGetTaxLotGains(t *testing.T) {
	r := setupRouterWithMock()
	req := httptest.NewRequest("GET", "/users/1/portfolios/1/tax-lots/1/gains?current_price=12.0", nil)
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d: %s", rw.Code, rw.Body.String())
	}
}

func TestGetPortfolioTaxGains(t *testing.T) {
	r := setupRouterWithMock()
	body := map[string]interface{}{"current_prices": map[string]float64{"1": 12.0}}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/users/1/portfolios/1/tax-gains", bytes.NewReader(b))
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d: %s", rw.Code, rw.Body.String())
	}
}

func TestGetTaxableGains(t *testing.T) {
	r := setupRouterWithMock()
	req := httptest.NewRequest("GET", "/users/1/portfolios/1/tax-report", nil)
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d: %s", rw.Code, rw.Body.String())
	}
}

func TestGetTaxTransactions(t *testing.T) {
	r := setupRouterWithMock()
	req := httptest.NewRequest("GET", "/users/1/portfolios/1/tax-transactions", nil)
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d: %s", rw.Code, rw.Body.String())
	}
}
