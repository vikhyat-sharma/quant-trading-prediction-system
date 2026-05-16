package db

import "time"

// PredictionMetric tracks accuracy and performance of predictions
type PredictionMetric struct {
	ID                int        `json:"id" db:"id"`
	PredictionID      int        `json:"prediction_id" db:"prediction_id"`
	StockID           int        `json:"stock_id" db:"stock_id"`
	Algorithm         string     `json:"algorithm" db:"algorithm"`
	PredictedPrice    float64    `json:"predicted_price" db:"predicted_price"`
	ActualPrice       *float64   `json:"actual_price" db:"actual_price"`
	AbsoluteError     *float64   `json:"absolute_error" db:"absolute_error"`
	PercentError      *float64   `json:"percent_error" db:"percent_error"`
	IsAccurate        *bool      `json:"is_accurate" db:"is_accurate"`
	AccuracyThreshold float64    `json:"accuracy_threshold" db:"accuracy_threshold"`
	PredictionDate    time.Time  `json:"prediction_date" db:"prediction_date"`
	ActualDate        *time.Time `json:"actual_date" db:"actual_date"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

// AlgorithmPerformance summarizes algorithm accuracy across all predictions
type AlgorithmPerformance struct {
	Algorithm              string    `json:"algorithm" db:"algorithm"`
	TotalPredictions       int       `json:"total_predictions" db:"total_predictions"`
	AccuratePredictions    int       `json:"accurate_predictions" db:"accurate_predictions"`
	AccuracyRate           float64   `json:"accuracy_rate" db:"accuracy_rate"`
	AverageAbsoluteError   float64   `json:"average_absolute_error" db:"average_absolute_error"`
	AveragePercentError    float64   `json:"average_percent_error" db:"average_percent_error"`
	AverageConfidenceScore float64   `json:"average_confidence_score" db:"average_confidence_score"`
	LastUpdated            time.Time `json:"last_updated" db:"last_updated"`
}

// PortfolioPerformance tracks portfolio metrics over time
type PortfolioPerformance struct {
	ID            int       `json:"id" db:"id"`
	PortfolioID   int       `json:"portfolio_id" db:"portfolio_id"`
	TotalValue    float64   `json:"total_value" db:"total_value"`
	CostBasis     float64   `json:"cost_basis" db:"cost_basis"`
	GainLoss      float64   `json:"gain_loss" db:"gain_loss"`
	ReturnPercent float64   `json:"return_percent" db:"return_percent"`
	DailyReturn   float64   `json:"daily_return" db:"daily_return"`
	Volatility    float64   `json:"volatility" db:"volatility"`
	Sharpe        float64   `json:"sharpe" db:"sharpe"`
	MaxDrawdown   float64   `json:"max_drawdown" db:"max_drawdown"`
	RecordDate    time.Time `json:"record_date" db:"record_date"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}
