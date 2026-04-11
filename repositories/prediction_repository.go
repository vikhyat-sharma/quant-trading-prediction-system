package repositories

import (
	"database/sql"

	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
)

type PredictionRepository struct {
	db *sql.DB
}

func NewPredictionRepository(db *sql.DB) *PredictionRepository {
	return &PredictionRepository{db: db}
}

func (r *PredictionRepository) GetPredictionsByStockID(stockID int) ([]*db.Prediction, error) {
	rows, err := r.db.Query("SELECT id, stock_id, predicted_price, date FROM predictions WHERE stock_id = $1", stockID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var predictions []*db.Prediction
	for rows.Next() {
		var prediction db.Prediction
		err := rows.Scan(&prediction.ID, &prediction.StockID, &prediction.PredictedPrice, &prediction.Date)
		if err != nil {
			return nil, err
		}
		predictions = append(predictions, &prediction)
	}
	return predictions, nil
}

func (r *PredictionRepository) CreatePrediction(prediction *db.Prediction) (*db.Prediction, error) {
	var id int
	err := r.db.QueryRow(
		"INSERT INTO predictions (stock_id, predicted_price, date) VALUES ($1, $2, $3) RETURNING id",
		prediction.StockID, prediction.PredictedPrice, prediction.Date,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	prediction.ID = id
	return prediction, nil
}
