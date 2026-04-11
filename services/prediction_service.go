package services

import (
	"time"

	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/repositories"
)

type PredictionService struct {
	repo *repositories.PredictionRepository
}

func NewPredictionService(repo *repositories.PredictionRepository) *PredictionService {
	return &PredictionService{repo: repo}
}

func (s *PredictionService) GetPredictionsByStockID(stockID int) ([]*db.Prediction, error) {
	return s.repo.GetPredictionsByStockID(stockID)
}

// Mock prediction generation
func (s *PredictionService) GeneratePrediction(stockID int) (*db.Prediction, error) {
	// Simple mock: predict price as 100.0
	prediction := &db.Prediction{
		StockID:        stockID,
		PredictedPrice: 100.0,
		Date:           time.Now(),
	}
	// Save to database
	return s.repo.CreatePrediction(prediction)
}
