package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
)

type PriceHistoryRepository struct {
	db *sql.DB
}

// PriceHistoryFilter holds filtering criteria for price history
type PriceHistoryFilter struct {
	StockID   int       // Filter by stock ID
	StartDate time.Time // Filter by start date (inclusive)
	EndDate   time.Time // Filter by end date (inclusive)
	MinPrice  float64   // Filter by minimum price
	MaxPrice  float64   // Filter by maximum price
}

// SearchAndFilterPriceHistory searches and filters price history based on criteria
func (r *PriceHistoryRepository) SearchAndFilterPriceHistory(filter *PriceHistoryFilter) ([]*db.PriceHistory, error) {
	query := "SELECT id, stock_id, price, date, created_at FROM price_history WHERE 1=1"
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
		query += fmt.Sprintf(" AND price BETWEEN $%d AND $%d", argCount, argCount+1)
		args = append(args, filter.MinPrice, filter.MaxPrice)
		argCount += 2
	} else if filter.MinPrice > 0 {
		query += fmt.Sprintf(" AND price >= $%d", argCount)
		args = append(args, filter.MinPrice)
		argCount++
	} else if filter.MaxPrice > 0 {
		query += fmt.Sprintf(" AND price <= $%d", argCount)
		args = append(args, filter.MaxPrice)
		argCount++
	}

	query += " ORDER BY date DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var priceHistories []*db.PriceHistory
	for rows.Next() {
		var price db.PriceHistory
		if err := rows.Scan(&price.ID, &price.StockID, &price.Price, &price.Date, &price.CreatedAt); err != nil {
			return nil, err
		}
		priceHistories = append(priceHistories, &price)
	}

	return priceHistories, nil
}

func NewPriceHistoryRepository(database *sql.DB) *PriceHistoryRepository {
	return &PriceHistoryRepository{db: database}
}

func (r *PriceHistoryRepository) GetPriceHistoryByStockID(stockID int) ([]*db.PriceHistory, error) {
	rows, err := r.db.Query(
		"SELECT id, stock_id, price, date, created_at FROM price_history WHERE stock_id = $1 ORDER BY date DESC",
		stockID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var priceHistories []*db.PriceHistory
	for rows.Next() {
		var price db.PriceHistory
		err := rows.Scan(&price.ID, &price.StockID, &price.Price, &price.Date, &price.CreatedAt)
		if err != nil {
			return nil, err
		}
		priceHistories = append(priceHistories, &price)
	}

	if len(priceHistories) == 0 {
		return nil, db.ErrRecordNotFound
	}

	return priceHistories, nil
}

func (r *PriceHistoryRepository) GetPriceHistoryByStockIDAndDateRange(stockID int, startDate, endDate time.Time) ([]*db.PriceHistory, error) {
	rows, err := r.db.Query(
		"SELECT id, stock_id, price, date, created_at FROM price_history WHERE stock_id = $1 AND date BETWEEN $2 AND $3 ORDER BY date DESC",
		stockID, startDate, endDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var priceHistories []*db.PriceHistory
	for rows.Next() {
		var price db.PriceHistory
		err := rows.Scan(&price.ID, &price.StockID, &price.Price, &price.Date, &price.CreatedAt)
		if err != nil {
			return nil, err
		}
		priceHistories = append(priceHistories, &price)
	}

	if len(priceHistories) == 0 {
		return nil, db.ErrRecordNotFound
	}

	return priceHistories, nil
}

func (r *PriceHistoryRepository) RecordPrice(priceHistory *db.PriceHistory) (*db.PriceHistory, error) {
	var id int
	err := r.db.QueryRow(
		"INSERT INTO price_history (stock_id, price, date) VALUES ($1, $2, $3) RETURNING id, created_at",
		priceHistory.StockID, priceHistory.Price, priceHistory.Date,
	).Scan(&id, &priceHistory.CreatedAt)
	if err != nil {
		return nil, err
	}
	priceHistory.ID = id
	return priceHistory, nil
}

func (r *PriceHistoryRepository) DeletePriceHistoryByStockID(stockID int) error {
	_, err := r.db.Exec("DELETE FROM price_history WHERE stock_id = $1", stockID)
	return err
}

func (r *PriceHistoryRepository) GetLatestPrice(stockID int) (*db.PriceHistory, error) {
	var price db.PriceHistory
	err := r.db.QueryRow(
		"SELECT id, stock_id, price, date, created_at FROM price_history WHERE stock_id = $1 ORDER BY date DESC LIMIT 1",
		stockID,
	).Scan(&price.ID, &price.StockID, &price.Price, &price.Date, &price.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, db.ErrRecordNotFound
		}
		return nil, err
	}
	return &price, nil
}
