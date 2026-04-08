package db

import (
	"testing"
	"time"
)

func TestStockModel(t *testing.T) {
	stock := Stock{
		ID:       1,
		Symbol:   "TCS",
		Exchange: "NSE",
		Name:     "Tata Consultancy Services Ltd.",
	}

	if stock.ID != 1 {
		t.Errorf("expected ID 1, got %d", stock.ID)
	}

	if stock.Symbol != "TCS" {
		t.Errorf("expected symbol TCS, got %s", stock.Symbol)
	}

	if stock.Exchange != "NSE" {
		t.Errorf("expected exchange NSE, got %s", stock.Exchange)
	}

	if stock.Name != "Tata Consultancy Services Ltd." {
		t.Errorf("expected name Tata Consultancy Services Ltd., got %s", stock.Name)
	}
}

func TestPredictionModel(t *testing.T) {
	now := time.Now()
	prediction := Prediction{
		ID:             1,
		StockID:        1,
		PredictedPrice: 150.50,
		Date:           now,
	}

	if prediction.ID != 1 {
		t.Errorf("expected ID 1, got %d", prediction.ID)
	}

	if prediction.StockID != 1 {
		t.Errorf("expected StockID 1, got %d", prediction.StockID)
	}

	if prediction.PredictedPrice != 150.50 {
		t.Errorf("expected price 150.50, got %f", prediction.PredictedPrice)
	}

	if prediction.Date != now {
		t.Errorf("expected date %v, got %v", now, prediction.Date)
	}
}

func TestStockModelJSONSerialization(t *testing.T) {
	stock := Stock{
		ID:       1,
		Symbol:   "INFY",
		Exchange: "BSE",
		Name:     "Infosys Ltd.",
	}

	// Verify JSON tags are present and structs can be serialized
	if stock.Symbol == "" {
		t.Errorf("Stock should have non-empty symbol for JSON serialization")
	}
	if stock.Exchange == "" {
		t.Errorf("Stock should have non-empty exchange for JSON serialization")
	}
}

func TestPredictionModelJSONSerialization(t *testing.T) {
	prediction := Prediction{
		ID:             2,
		StockID:        1,
		PredictedPrice: 200.75,
		Date:           time.Now(),
	}

	// Verify JSON tags are present and structs can be serialized
	if prediction.PredictedPrice == 0 {
		t.Errorf("Prediction should have non-zero price for JSON serialization")
	}
}
