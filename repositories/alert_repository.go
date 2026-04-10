package repositories

import (
	"database/sql"

	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
)

type AlertRepository struct {
	db *sql.DB
}

func NewAlertRepository(database *sql.DB) *AlertRepository {
	return &AlertRepository{db: database}
}

func (r *AlertRepository) CreateAlert(alert *db.Alert) (*db.Alert, error) {
	var id int
	err := r.db.QueryRow(
		"INSERT INTO price_alerts (stock_id, threshold, condition, enabled) VALUES ($1, $2, $3, $4) RETURNING id",
		alert.StockID, alert.Threshold, alert.Condition, alert.Enabled,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	alert.ID = id
	return alert, nil
}

func (r *AlertRepository) GetAlertsByStockID(stockID int) ([]*db.Alert, error) {
	rows, err := r.db.Query(
		"SELECT id, stock_id, threshold, condition, enabled, created_at FROM price_alerts WHERE stock_id = $1 ORDER BY created_at DESC",
		stockID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []*db.Alert
	for rows.Next() {
		var alert db.Alert
		err := rows.Scan(&alert.ID, &alert.StockID, &alert.Threshold, &alert.Condition, &alert.Enabled, &alert.CreatedAt)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, &alert)
	}

	return alerts, nil
}

func (r *AlertRepository) GetEnabledAlertsByStockID(stockID int) ([]*db.Alert, error) {
	rows, err := r.db.Query(
		"SELECT id, stock_id, threshold, condition, enabled, created_at FROM price_alerts WHERE stock_id = $1 AND enabled = TRUE ORDER BY created_at DESC",
		stockID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []*db.Alert
	for rows.Next() {
		var alert db.Alert
		err := rows.Scan(&alert.ID, &alert.StockID, &alert.Threshold, &alert.Condition, &alert.Enabled, &alert.CreatedAt)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, &alert)
	}

	return alerts, nil
}

func (r *AlertRepository) DeleteAlert(stockID int, alertID int) error {
	result, err := r.db.Exec(
		"DELETE FROM price_alerts WHERE id = $1 AND stock_id = $2",
		alertID, stockID,
	)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return db.ErrRecordNotFound
	}

	return nil
}
