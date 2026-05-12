package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
)

type PredictionRepository struct {
	db *sql.DB
}

// PredictionFilter holds filtering criteria for predictions
type PredictionFilter struct {
	StockID   int       // Filter by stock ID
	StartDate time.Time // Filter by start date (inclusive)
	EndDate   time.Time // Filter by end date (inclusive)
	MinPrice  float64   // Filter by minimum predicted price
	MaxPrice  float64   // Filter by maximum predicted price
}

// SearchAndFilterPredictions searches and filters predictions based on criteria
func (r *PredictionRepository) SearchAndFilterPredictions(filter *PredictionFilter) ([]*db.Prediction, error) {
	query := "SELECT id, stock_id, predicted_price, algorithm, confidence_score, upper_bound, lower_bound, date, created_at FROM predictions WHERE 1=1"
	var args []interface{}
	argCount := 1

	if filter.StockID > 0 {
		query += fmt.Sprintf(" AND stock_id = $%d", argCount)
		args = append(args, filter.StockID)
		argCount++
	}

	if !filter.StartDate.IsZero() && !filter.EndDate.IsZero() {
		query += fmt.Sprintf(" AND date BETWEEN $%d AND $%d", argCount, argCount+1)
		args = append(args, filter.StartDate, filter.EndDate)
		argCount += 2
	} else if !filter.StartDate.IsZero() {
		query += fmt.Sprintf(" AND date >= $%d", argCount)
		args = append(args, filter.StartDate)
		argCount++
	} else if !filter.EndDate.IsZero() {
		query += fmt.Sprintf(" AND date <= $%d", argCount)
		args = append(args, filter.EndDate)
		argCount++
	}

	if filter.MinPrice > 0 && filter.MaxPrice > 0 {
		query += fmt.Sprintf(" AND predicted_price BETWEEN $%d AND $%d", argCount, argCount+1)
		args = append(args, filter.MinPrice, filter.MaxPrice)
		argCount += 2
	} else if filter.MinPrice > 0 {
		query += fmt.Sprintf(" AND predicted_price >= $%d", argCount)
		args = append(args, filter.MinPrice)
		argCount++
	} else if filter.MaxPrice > 0 {
		query += fmt.Sprintf(" AND predicted_price <= $%d", argCount)
		args = append(args, filter.MaxPrice)
		argCount++
	}

	query += " ORDER BY date DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var predictions []*db.Prediction
	for rows.Next() {
		var prediction db.Prediction
		if err := rows.Scan(&prediction.ID, &prediction.StockID, &prediction.PredictedPrice, &prediction.Algorithm, &prediction.ConfidenceScore, &prediction.UpperBound, &prediction.LowerBound, &prediction.Date, &prediction.CreatedAt); err != nil {
			return nil, err
		}
		predictions = append(predictions, &prediction)
	}

	return predictions, nil
}

func NewPredictionRepository(db *sql.DB) *PredictionRepository {
	return &PredictionRepository{db: db}
}

func (r *PredictionRepository) GetPredictionsByStockID(stockID int) ([]*db.Prediction, error) {
	rows, err := r.db.Query("SELECT id, stock_id, predicted_price, algorithm, confidence_score, upper_bound, lower_bound, date, created_at FROM predictions WHERE stock_id = $1 ORDER BY date DESC", stockID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var predictions []*db.Prediction
	for rows.Next() {
		var prediction db.Prediction
		err := rows.Scan(&prediction.ID, &prediction.StockID, &prediction.PredictedPrice, &prediction.Algorithm, &prediction.ConfidenceScore, &prediction.UpperBound, &prediction.LowerBound, &prediction.Date, &prediction.CreatedAt)
		if err != nil {
			return nil, err
		}
		predictions = append(predictions, &prediction)
	}
	return predictions, nil
}

func (r *PredictionRepository) CreatePrediction(prediction *db.Prediction) (*db.Prediction, error) {
	var id int
	var createdAt time.Time
	err := r.db.QueryRow(
		"INSERT INTO predictions (stock_id, predicted_price, algorithm, confidence_score, upper_bound, lower_bound, date) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at",
		prediction.StockID, prediction.PredictedPrice, prediction.Algorithm, prediction.ConfidenceScore, prediction.UpperBound, prediction.LowerBound, prediction.Date,
	).Scan(&id, &createdAt)
	if err != nil {
		return nil, err
	}
	prediction.ID = id
	prediction.CreatedAt = createdAt
	return prediction, nil
}
