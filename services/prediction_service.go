package services

import (
	"fmt"
	"time"

	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/repositories"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/services/algorithms"
)

type PredictionService struct {
	repo             *repositories.PredictionRepository
	priceHistoryRepo *repositories.PriceHistoryRepository
	metricsRepo      *repositories.PredictionMetricsRepository
	defaultAlgorithm string
	lookbackPeriod   int
}

func NewPredictionService(
	repo *repositories.PredictionRepository,
	priceHistoryRepo *repositories.PriceHistoryRepository,
) *PredictionService {
	return &PredictionService{
		repo:             repo,
		priceHistoryRepo: priceHistoryRepo,
		defaultAlgorithm: "ENSEMBLE",
		lookbackPeriod:   100,
	}
}

func NewPredictionServiceWithMetrics(
	repo *repositories.PredictionRepository,
	priceHistoryRepo *repositories.PriceHistoryRepository,
	metricsRepo *repositories.PredictionMetricsRepository,
) *PredictionService {
	return &PredictionService{
		repo:             repo,
		priceHistoryRepo: priceHistoryRepo,
		metricsRepo:      metricsRepo,
		defaultAlgorithm: "ENSEMBLE",
		lookbackPeriod:   100,
	}
}

func (s *PredictionService) GetPredictionsByStockID(stockID int) ([]*db.Prediction, error) {
	return s.repo.GetPredictionsByStockID(stockID)
}

func (s *PredictionService) SearchAndFilterPredictions(filter *repositories.PredictionFilter) ([]*db.Prediction, error) {
	return s.repo.SearchAndFilterPredictions(filter)
}

// GeneratePrediction generates a prediction for a stock using the default algorithm
func (s *PredictionService) GeneratePrediction(stockID int) (*db.Prediction, error) {
	return s.GeneratePredictionWithAlgorithm(stockID, s.defaultAlgorithm)
}

// GeneratePredictionWithAlgorithm generates a prediction using a specific algorithm
func (s *PredictionService) GeneratePredictionWithAlgorithm(stockID int, algorithmType string) (*db.Prediction, error) {
	// Get historical price data
	prices, err := s.priceHistoryRepo.GetHistoricalPrices(stockID, s.lookbackPeriod)
	if err != nil {
		return nil, fmt.Errorf("failed to get price history: %w", err)
	}

	if len(prices) < 5 {
		return nil, fmt.Errorf("insufficient price data for prediction (need at least 5, got %d)", len(prices))
	}

	// Extract prices in chronological order
	priceValues := make([]float64, len(prices))
	for i, p := range prices {
		priceValues[i] = p.Price
	}

	// Generate prediction based on algorithm type
	var result *algorithms.PredictionResult
	switch algorithmType {
	case "SMA":
		result = algorithms.SimpleMovingAveragePrediction(priceValues)
	case "EMA":
		result = algorithms.ExponentialMovingAveragePrediction(priceValues)
	case "MOMENTUM":
		result = algorithms.MomentumPrediction(priceValues)
	case "MEAN_REVERSION":
		result = algorithms.MeanReversionPrediction(priceValues)
	case "ENSEMBLE":
		result = algorithms.EnsemblePrediction(priceValues)
	default:
		result = algorithms.EnsemblePrediction(priceValues)
	}

	prediction := &db.Prediction{
		StockID:         stockID,
		PredictedPrice:  result.PredictedPrice,
		Algorithm:       result.Algorithm,
		ConfidenceScore: result.ConfidenceScore,
		UpperBound:      result.UpperBound,
		LowerBound:      result.LowerBound,
		Date:            time.Now().AddDate(0, 0, 1), // Predict for next day
	}

	if s.repo == nil {
		return prediction, nil
	}

	return s.repo.CreatePrediction(prediction)
}

// GetPredictionMetrics retrieves metrics for a specific prediction
func (s *PredictionService) GetPredictionMetrics(predictionID int) (*db.PredictionMetric, error) {
	if s.metricsRepo == nil {
		return nil, fmt.Errorf("metrics repository not initialized")
	}
	return s.metricsRepo.GetMetricByPredictionID(predictionID)
}

// GetAlgorithmPerformance retrieves performance stats for an algorithm
func (s *PredictionService) GetAlgorithmPerformance(algorithm string) (*db.AlgorithmPerformance, error) {
	if s.metricsRepo == nil {
		return nil, fmt.Errorf("metrics repository not initialized")
	}
	return s.metricsRepo.GetPerformanceByAlgorithm(algorithm)
}
