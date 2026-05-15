package repositories

import (
	"database/sql"
	"time"

	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
)

type PredictionMetricsRepository struct {
	db *sql.DB
}

func NewPredictionMetricsRepository(database *sql.DB) *PredictionMetricsRepository {
	return &PredictionMetricsRepository{db: database}
}

// CreateMetric creates a new prediction metric record
func (r *PredictionMetricsRepository) CreateMetric(metric *db.PredictionMetric) (*db.PredictionMetric, error) {
	var id int
	err := r.db.QueryRow(
		`INSERT INTO prediction_metrics 
		(prediction_id, stock_id, algorithm, predicted_price, actual_price, absolute_error, percent_error, is_accurate, accuracy_threshold, prediction_date, actual_date) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
		RETURNING id, created_at, updated_at`,
		metric.PredictionID, metric.StockID, metric.Algorithm, metric.PredictedPrice, metric.ActualPrice,
		metric.AbsoluteError, metric.PercentError, metric.IsAccurate, metric.AccuracyThreshold,
		metric.PredictionDate, metric.ActualDate,
	).Scan(&id, &metric.CreatedAt, &metric.UpdatedAt)

	if err != nil {
		return nil, err
	}

	metric.ID = id
	return metric, nil
}

// GetMetricByPredictionID gets metrics for a specific prediction
func (r *PredictionMetricsRepository) GetMetricByPredictionID(predictionID int) (*db.PredictionMetric, error) {
	var metric db.PredictionMetric
	err := r.db.QueryRow(
		`SELECT id, prediction_id, stock_id, algorithm, predicted_price, actual_price, absolute_error, percent_error, is_accurate, accuracy_threshold, prediction_date, actual_date, created_at, updated_at 
		FROM prediction_metrics 
		WHERE prediction_id = $1`,
		predictionID,
	).Scan(&metric.ID, &metric.PredictionID, &metric.StockID, &metric.Algorithm, &metric.PredictedPrice,
		&metric.ActualPrice, &metric.AbsoluteError, &metric.PercentError, &metric.IsAccurate,
		&metric.AccuracyThreshold, &metric.PredictionDate, &metric.ActualDate, &metric.CreatedAt, &metric.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, db.ErrRecordNotFound
		}
		return nil, err
	}

	return &metric, nil
}

// GetPerformanceByAlgorithm gets aggregated performance metrics for an algorithm
func (r *PredictionMetricsRepository) GetPerformanceByAlgorithm(algorithm string) (*db.AlgorithmPerformance, error) {
	var performance db.AlgorithmPerformance

	err := r.db.QueryRow(
		`SELECT 
		algorithm, 
		COUNT(*) as total_predictions,
		SUM(CASE WHEN is_accurate = true THEN 1 ELSE 0 END) as accurate_predictions,
		AVG(CASE WHEN is_accurate IS NOT NULL THEN CAST(is_accurate AS INT) ELSE NULL END) * 100 as accuracy_rate,
		AVG(absolute_error) as average_absolute_error,
		AVG(percent_error) as average_percent_error,
		MAX(created_at) as last_updated
		FROM prediction_metrics 
		WHERE algorithm = $1 AND actual_price IS NOT NULL
		GROUP BY algorithm`,
		algorithm,
	).Scan(&performance.Algorithm, &performance.TotalPredictions, &performance.AccuratePredictions,
		&performance.AccuracyRate, &performance.AverageAbsoluteError, &performance.AveragePercentError,
		&performance.LastUpdated)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, db.ErrRecordNotFound
		}
		return nil, err
	}

	return &performance, nil
}

// GetMetricsForDateRange gets metrics for a specific date range
func (r *PredictionMetricsRepository) GetMetricsForDateRange(algorithm string, startDate, endDate time.Time) ([]*db.PredictionMetric, error) {
	rows, err := r.db.Query(
		`SELECT id, prediction_id, stock_id, algorithm, predicted_price, actual_price, absolute_error, percent_error, is_accurate, accuracy_threshold, prediction_date, actual_date, created_at, updated_at 
		FROM prediction_metrics 
		WHERE algorithm = $1 AND created_at BETWEEN $2 AND $3 AND actual_price IS NOT NULL
		ORDER BY created_at DESC`,
		algorithm, startDate, endDate,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []*db.PredictionMetric
	for rows.Next() {
		var metric db.PredictionMetric
		err := rows.Scan(&metric.ID, &metric.PredictionID, &metric.StockID, &metric.Algorithm, &metric.PredictedPrice,
			&metric.ActualPrice, &metric.AbsoluteError, &metric.PercentError, &metric.IsAccurate,
			&metric.AccuracyThreshold, &metric.PredictionDate, &metric.ActualDate, &metric.CreatedAt, &metric.UpdatedAt)

		if err != nil {
			return nil, err
		}
		metrics = append(metrics, &metric)
	}

	return metrics, nil
}

// UpdateMetricWithActualPrice updates a metric with actual price and calculates error
func (r *PredictionMetricsRepository) UpdateMetricWithActualPrice(metricID int, actualPrice float64, actualDate time.Time) (*db.PredictionMetric, error) {
	// Get the metric first to calculate errors
	metric := &db.PredictionMetric{}
	err := r.db.QueryRow(
		`SELECT id, prediction_id, stock_id, algorithm, predicted_price, actual_price, absolute_error, percent_error, is_accurate, accuracy_threshold, prediction_date, actual_date, created_at, updated_at 
		FROM prediction_metrics 
		WHERE id = $1`,
		metricID,
	).Scan(&metric.ID, &metric.PredictionID, &metric.StockID, &metric.Algorithm, &metric.PredictedPrice,
		&metric.ActualPrice, &metric.AbsoluteError, &metric.PercentError, &metric.IsAccurate,
		&metric.AccuracyThreshold, &metric.PredictionDate, &metric.ActualDate, &metric.CreatedAt, &metric.UpdatedAt)

	if err != nil {
		return nil, err
	}

	// Calculate errors
	absError := metric.PredictedPrice - actualPrice
	if absError < 0 {
		absError = -absError
	}

	percentError := (absError / actualPrice) * 100
	isAccurate := absError <= metric.AccuracyThreshold

	// Update the metric
	err = r.db.QueryRow(
		`UPDATE prediction_metrics 
		SET actual_price = $1, absolute_error = $2, percent_error = $3, is_accurate = $4, actual_date = $5, updated_at = $6
		WHERE id = $7
		RETURNING id, prediction_id, stock_id, algorithm, predicted_price, actual_price, absolute_error, percent_error, is_accurate, accuracy_threshold, prediction_date, actual_date, created_at, updated_at`,
		actualPrice, absError, percentError, isAccurate, actualDate, time.Now(), metricID,
	).Scan(&metric.ID, &metric.PredictionID, &metric.StockID, &metric.Algorithm, &metric.PredictedPrice,
		&metric.ActualPrice, &metric.AbsoluteError, &metric.PercentError, &metric.IsAccurate,
		&metric.AccuracyThreshold, &metric.PredictionDate, &metric.ActualDate, &metric.CreatedAt, &metric.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return metric, nil
}

// GetAllAlgorithmPerformance gets performance stats for all algorithms
func (r *PredictionMetricsRepository) GetAllAlgorithmPerformance() ([]*db.AlgorithmPerformance, error) {
	rows, err := r.db.Query(
		`SELECT 
		algorithm, 
		COUNT(*) as total_predictions,
		SUM(CASE WHEN is_accurate = true THEN 1 ELSE 0 END) as accurate_predictions,
		AVG(CASE WHEN is_accurate IS NOT NULL THEN CAST(is_accurate AS INT) ELSE NULL END) * 100 as accuracy_rate,
		AVG(absolute_error) as average_absolute_error,
		AVG(percent_error) as average_percent_error,
		MAX(created_at) as last_updated
		FROM prediction_metrics 
		WHERE actual_price IS NOT NULL
		GROUP BY algorithm
		ORDER BY accuracy_rate DESC`,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var performances []*db.AlgorithmPerformance
	for rows.Next() {
		var performance db.AlgorithmPerformance
		err := rows.Scan(&performance.Algorithm, &performance.TotalPredictions, &performance.AccuratePredictions,
			&performance.AccuracyRate, &performance.AverageAbsoluteError, &performance.AveragePercentError,
			&performance.LastUpdated)

		if err != nil {
			return nil, err
		}
		performances = append(performances, &performance)
	}

	return performances, nil
}

// GetMetricsForStock gets all metrics for a specific stock
func (r *PredictionMetricsRepository) GetMetricsForStock(stockID int) ([]*db.PredictionMetric, error) {
	rows, err := r.db.Query(
		`SELECT id, prediction_id, stock_id, algorithm, predicted_price, actual_price, absolute_error, percent_error, is_accurate, accuracy_threshold, prediction_date, actual_date, created_at, updated_at 
		FROM prediction_metrics 
		WHERE stock_id = $1 AND actual_price IS NOT NULL
		ORDER BY created_at DESC`,
		stockID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []*db.PredictionMetric
	for rows.Next() {
		var metric db.PredictionMetric
		err := rows.Scan(&metric.ID, &metric.PredictionID, &metric.StockID, &metric.Algorithm, &metric.PredictedPrice,
			&metric.ActualPrice, &metric.AbsoluteError, &metric.PercentError, &metric.IsAccurate,
			&metric.AccuracyThreshold, &metric.PredictionDate, &metric.ActualDate, &metric.CreatedAt, &metric.UpdatedAt)

		if err != nil {
			return nil, err
		}
		metrics = append(metrics, &metric)
	}

	return metrics, nil
}
