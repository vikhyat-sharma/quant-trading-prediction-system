package services

import (
	"testing"
	"time"

	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
)

type mockTaxLotRepoForService struct {
	taxLot *db.TaxLot
	err    error
}

func (m *mockTaxLotRepoForService) GetTaxLotByID(id int) (*db.TaxLot, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.taxLot == nil {
		return nil, nil
	}
	return m.taxLot, nil
}

func (m *mockTaxLotRepoForService) GetTaxLotsByPortfolioID(portfolioID int) ([]db.TaxLot, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.taxLot == nil {
		return nil, nil
	}
	return []db.TaxLot{*m.taxLot}, nil
}

func (m *mockTaxLotRepoForService) GetActiveTaxLotsByStockID(portfolioID, stockID int) ([]db.TaxLot, error) {
	return nil, nil
}

func (m *mockTaxLotRepoForService) UpdateTaxLot(taxLot *db.TaxLot) error {
	return nil
}

func (m *mockTaxLotRepoForService) CreateTaxTransaction(transaction *db.TaxTransaction) error {
	return nil
}

func (m *mockTaxLotRepoForService) CreateTaxLot(taxLot *db.TaxLot) error {
	return nil
}

func (m *mockTaxLotRepoForService) GetTaxTransactionsByPortfolioID(portfolioID int) ([]db.TaxTransaction, error) {
	return nil, nil
}

type mockStockRepoForService struct {
	stock *db.Stock
	err   error
}

func (m *mockStockRepoForService) GetStock(id int) (*db.Stock, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.stock, nil
}

func TestGetTaxLotGains_Success(t *testing.T) {
	taxLot := &db.TaxLot{
		ID:              1,
		StockID:         1,
		Quantity:        100,
		QuantitySold:    40,
		CostPerShare:    10.0,
		TotalCost:       1000.0,
		AcquisitionDate: time.Now().AddDate(-1, -1, 0),
	}

	stockRepo := &mockStockRepoForService{stock: &db.Stock{ID: 1, Symbol: "TST", Name: "Test", Exchange: "NSE"}}
	taxLotRepo := &mockTaxLotRepoForService{taxLot: taxLot}
	service := &TaxLotService{taxLotRepo: taxLotRepo, stockRepo: stockRepo}

	gains, err := service.GetTaxLotGains(1, 12.5)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if gains.TaxLotID != 1 {
		t.Fatalf("expected tax lot id 1, got %d", gains.TaxLotID)
	}
	if gains.CurrentPrice != 12.5 {
		t.Fatalf("expected current price 12.5, got %f", gains.CurrentPrice)
	}
	if gains.UnrealizedGain <= 0 {
		t.Fatalf("expected unrealized gain > 0, got %f", gains.UnrealizedGain)
	}
}

func TestGetTaxLotGains_TaxLotNotFound(t *testing.T) {
	taxLotRepo := &mockTaxLotRepoForService{taxLot: nil}
	stockRepo := &mockStockRepoForService{stock: &db.Stock{ID: 1}}
	service := &TaxLotService{taxLotRepo: taxLotRepo, stockRepo: stockRepo}

	_, err := service.GetTaxLotGains(1, 10.0)
	if err == nil || err.Error() != "tax lot not found" {
		t.Fatalf("expected tax lot not found error, got %v", err)
	}
}

func TestGetTaxLotGains_StockNotFound(t *testing.T) {
	taxLot := &db.TaxLot{
		ID:              1,
		StockID:         1,
		Quantity:        50,
		QuantitySold:    10,
		CostPerShare:    10.0,
		TotalCost:       500.0,
		AcquisitionDate: time.Now().AddDate(-2, 0, 0),
	}

	taxLotRepo := &mockTaxLotRepoForService{taxLot: taxLot}
	stockRepo := &mockStockRepoForService{stock: nil}
	service := &TaxLotService{taxLotRepo: taxLotRepo, stockRepo: stockRepo}

	_, err := service.GetTaxLotGains(1, 11.0)
	if err == nil || err.Error() != "stock not found" {
		t.Fatalf("expected stock not found error, got %v", err)
	}
}
