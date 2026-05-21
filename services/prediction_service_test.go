package services

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/repositories"
)

// Test that PredictionService can be instantiated
func TestPredictionService_NewPredictionService(t *testing.T) {
	service := NewPredictionService(nil, nil)

	if service == nil {
		t.Errorf("expected PredictionService, got nil")
	}
}

// Test that GeneratePrediction returns a prediction when price history data is available
func TestPredictionService_GeneratePrediction_WithPriceHistory(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mockDB.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "stock_id", "price", "date", "created_at"}).
		AddRow(1, 1, 100.0, now.AddDate(0, 0, -4), now).
		AddRow(2, 1, 102.0, now.AddDate(0, 0, -3), now).
		AddRow(3, 1, 104.0, now.AddDate(0, 0, -2), now).
		AddRow(4, 1, 103.0, now.AddDate(0, 0, -1), now).
		AddRow(5, 1, 105.0, now, now)

	mock.ExpectQuery("SELECT id, stock_id, price, date, created_at FROM price_history WHERE stock_id = \\$1 ORDER BY date DESC LIMIT \\$2").
		WithArgs(1, 100).
		WillReturnRows(rows)

	priceHistoryRepo := repositories.NewPriceHistoryRepository(mockDB)
	service := NewPredictionService(nil, priceHistoryRepo)

	prediction, err := service.GeneratePrediction(1)
	if err != nil {
		t.Fatalf("GeneratePrediction failed: %v", err)
	}

	if prediction == nil {
		t.Fatal("expected prediction, got nil")
	}

	if prediction.StockID != 1 {
		t.Errorf("expected StockID 1, got %d", prediction.StockID)
	}

	if prediction.Date.IsZero() {
		t.Errorf("expected non-zero Date")
	}

	if prediction.Algorithm == "" {
		t.Errorf("expected algorithm name, got empty string")
	}

	if prediction.ConfidenceScore < 0 || prediction.ConfidenceScore > 1 {
		t.Errorf("expected confidence score between 0 and 1, got %f", prediction.ConfidenceScore)
	}
}
