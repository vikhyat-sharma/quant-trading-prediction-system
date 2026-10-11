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
	StockID   int
	StartDate time.Time
	EndDate   time.Time
	MinPrice  float64
	MaxPrice  float64
}

func NewPredictionRepository(database *sql.DB) *PredictionRepository {
	return &PredictionRepository{db: database}
}

const predictionColumns = `id, stock_id, predicted_price, algorithm, algorithm_version,
	confidence_score, upper_bound, lower_bound, date, created_at`

func scanPrediction(row interface {
	Scan(...interface{}) error
}) (*db.Prediction, error) {
	var p db.Prediction
	err := row.Scan(
		&p.ID, &p.StockID, &p.PredictedPrice, &p.Algorithm, &p.AlgorithmVersion,
		&p.ConfidenceScore, &p.UpperBound, &p.LowerBound, &p.Date, &p.CreatedAt,
	)
	return &p, err
}

func (r *PredictionRepository) GetPredictionsByStockID(stockID int) ([]*db.Prediction, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM predictions WHERE stock_id = $1 ORDER BY date DESC",
		predictionColumns,
	)
	rows, err := r.db.Query(query, stockID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var predictions []*db.Prediction
	for rows.Next() {
		p, err := scanPrediction(rows)
		if err != nil {
			return nil, err
		}
		predictions = append(predictions, p)
	}
	return predictions, rows.Err()
}

// GetPredictionsByStockIDPaginated returns a page of predictions for a stock.
func (r *PredictionRepository) GetPredictionsByStockIDPaginated(stockID, limit, offset int) ([]*db.Prediction, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM predictions WHERE stock_id = $1 ORDER BY date DESC LIMIT $2 OFFSET $3",
		predictionColumns,
	)
	rows, err := r.db.Query(query, stockID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var predictions []*db.Prediction
	for rows.Next() {
		p, err := scanPrediction(rows)
		if err != nil {
			return nil, err
		}
		predictions = append(predictions, p)
	}
	return predictions, rows.Err()
}

func (r *PredictionRepository) SearchAndFilterPredictions(filter *PredictionFilter) ([]*db.Prediction, error) {
	query := fmt.Sprintf("SELECT %s FROM predictions WHERE 1=1", predictionColumns)
	var args []interface{}
	n := 1

	if filter.StockID > 0 {
		query += fmt.Sprintf(" AND stock_id = $%d", n)
		args = append(args, filter.StockID)
		n++
	}
	if !filter.StartDate.IsZero() {
		query += fmt.Sprintf(" AND date >= $%d", n)
		args = append(args, filter.StartDate)
		n++
	}
	if !filter.EndDate.IsZero() {
		query += fmt.Sprintf(" AND date <= $%d", n)
		args = append(args, filter.EndDate)
		n++
	}
	if filter.MinPrice > 0 {
		query += fmt.Sprintf(" AND predicted_price >= $%d", n)
		args = append(args, filter.MinPrice)
		n++
	}
	if filter.MaxPrice > 0 {
		query += fmt.Sprintf(" AND predicted_price <= $%d", n)
		args = append(args, filter.MaxPrice)
		n++
	}
	query += " ORDER BY date DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var predictions []*db.Prediction
	for rows.Next() {
		p, err := scanPrediction(rows)
		if err != nil {
			return nil, err
		}
		predictions = append(predictions, p)
	}
	return predictions, rows.Err()
}

func (r *PredictionRepository) CreatePrediction(prediction *db.Prediction) (*db.Prediction, error) {
	err := r.db.QueryRow(
		`INSERT INTO predictions
			(stock_id, predicted_price, algorithm, algorithm_version, confidence_score, upper_bound, lower_bound, date)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id, created_at`,
		prediction.StockID, prediction.PredictedPrice, prediction.Algorithm, prediction.AlgorithmVersion,
		prediction.ConfidenceScore, prediction.UpperBound, prediction.LowerBound, prediction.Date,
	).Scan(&prediction.ID, &prediction.CreatedAt)
	if err != nil {
		return nil, err
	}
	return prediction, nil
}
