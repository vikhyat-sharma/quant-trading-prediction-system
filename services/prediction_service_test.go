package services

import (
	"testing"
)

// Test that PredictionService can be instantiated
func TestPredictionService_NewPredictionService(t *testing.T) {
	// This test verifies that the PredictionService constructor works
	// In a real scenario with dependency injection, you'd test with mocked repositories
	service := NewPredictionService(nil)

	if service == nil {
		t.Errorf("expected PredictionService, got nil")
	}
}

// Test that GeneratePrediction returns a prediction with expected fields
func TestPredictionService_GeneratePrediction_Fields(t *testing.T) {
	service := NewPredictionService(nil)

	prediction, err := service.GeneratePrediction(1)

	if err != nil {
		t.Errorf("GeneratePrediction failed: %v", err)
	}

	if prediction == nil {
		t.Errorf("expected prediction, got nil")
	}

	if prediction.StockID != 1 {
		t.Errorf("expected StockID 1, got %d", prediction.StockID)
	}

	if prediction.PredictedPrice != 100.0 {
		t.Errorf("expected PredictedPrice 100.0, got %f", prediction.PredictedPrice)
	}

	if prediction.Date.IsZero() {
		t.Errorf("expected non-zero Date")
	}
}
