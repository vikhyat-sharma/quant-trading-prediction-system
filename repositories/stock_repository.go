package repositories

import (
	"database/sql"

	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
)

type StockRepository struct {
	db *sql.DB
}

func NewStockRepository(db *sql.DB) *StockRepository {
	return &StockRepository{db: db}
}

func (r *StockRepository) GetStock(id int) (*db.Stock, error) {
	var stock db.Stock
	err := r.db.QueryRow("SELECT id, symbol, name FROM stocks WHERE id = $1", id).Scan(&stock.ID, &stock.Symbol, &stock.Name)
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

func (r *StockRepository) GetAllStocks() ([]*db.Stock, error) {
	rows, err := r.db.Query("SELECT id, symbol, name FROM stocks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var stocks []*db.Stock
	for rows.Next() {
		var stock db.Stock
		err := rows.Scan(&stock.ID, &stock.Symbol, &stock.Name)
		if err != nil {
			return nil, err
		}
		stocks = append(stocks, &stock)
	}
	return stocks, nil
}
