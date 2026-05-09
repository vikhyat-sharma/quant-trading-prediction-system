package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
)

type StockRepository struct {
	db *sql.DB
}

// StockFilter holds filtering criteria for stocks
type StockFilter struct {
	Search   string // Search by symbol or name
	Exchange string // Filter by exchange
}

// SearchAndFilterStocks searches and filters stocks based on criteria
func (r *StockRepository) SearchAndFilterStocks(filter *StockFilter) ([]*db.Stock, error) {
	query := "SELECT id, symbol, exchange, name FROM stocks WHERE 1=1"
	var args []interface{}
	argCount := 1

	if filter.Search != "" {
		searchTerm := "%" + strings.ToUpper(filter.Search) + "%"
		query += fmt.Sprintf(" AND (UPPER(symbol) LIKE $%d OR UPPER(name) LIKE $%d)", argCount, argCount)
		args = append(args, searchTerm)
		argCount++
	}

	if filter.Exchange != "" {
		query += fmt.Sprintf(" AND UPPER(exchange) = UPPER($%d)", argCount)
		args = append(args, filter.Exchange)
		argCount++
	}

	query += " ORDER BY symbol"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []*db.Stock
	for rows.Next() {
		var stock db.Stock
		if err := rows.Scan(&stock.ID, &stock.Symbol, &stock.Exchange, &stock.Name); err != nil {
			return nil, err
		}
		stocks = append(stocks, &stock)
	}

	return stocks, nil
}

func NewStockRepository(db *sql.DB) *StockRepository {
	return &StockRepository{db: db}
}

func (r *StockRepository) GetStock(id int) (*db.Stock, error) {
	var stock db.Stock
	err := r.db.QueryRow("SELECT id, symbol, exchange, name FROM stocks WHERE id = $1", id).
		Scan(&stock.ID, &stock.Symbol, &stock.Exchange, &stock.Name)
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

func (r *StockRepository) GetAllStocks() ([]*db.Stock, error) {
	rows, err := r.db.Query("SELECT id, symbol, exchange, name FROM stocks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var stocks []*db.Stock
	for rows.Next() {
		var stock db.Stock
		err := rows.Scan(&stock.ID, &stock.Symbol, &stock.Exchange, &stock.Name)
		if err != nil {
			return nil, err
		}
		stocks = append(stocks, &stock)
	}
	return stocks, nil
}

func (r *StockRepository) CreateStock(stock *db.Stock) (*db.Stock, error) {
	var id int
	err := r.db.QueryRow(
		"INSERT INTO stocks (symbol, exchange, name) VALUES ($1, $2, $3) RETURNING id",
		stock.Symbol, stock.Exchange, stock.Name,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	stock.ID = id
	return stock, nil
}

func (r *StockRepository) UpdateStock(id int, stock *db.Stock) (*db.Stock, error) {
	_, err := r.db.Exec(
		"UPDATE stocks SET symbol = $1, exchange = $2, name = $3 WHERE id = $4",
		stock.Symbol, stock.Exchange, stock.Name, id,
	)
	if err != nil {
		return nil, err
	}
	stock.ID = id
	return stock, nil
}

func (r *StockRepository) DeleteStock(id int) error {
	result, err := r.db.Exec("DELETE FROM stocks WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return db.ErrRecordNotFound
	}
	return nil
}
