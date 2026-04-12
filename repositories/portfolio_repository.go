package repositories

import (
	"database/sql"

	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
)

type PortfolioRepository struct {
	db *sql.DB
}

func NewPortfolioRepository(database *sql.DB) *PortfolioRepository {
	return &PortfolioRepository{db: database}
}

func (r *PortfolioRepository) GetPortfoliosByUserID(userID int) ([]*db.Portfolio, error) {
	rows, err := r.db.Query(
		"SELECT id, user_id, name, description, created_at FROM portfolios WHERE user_id = $1 ORDER BY id",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var portfolios []*db.Portfolio
	for rows.Next() {
		var portfolio db.Portfolio
		if err := rows.Scan(&portfolio.ID, &portfolio.UserID, &portfolio.Name, &portfolio.Description, &portfolio.CreatedAt); err != nil {
			return nil, err
		}
		portfolios = append(portfolios, &portfolio)
	}

	if len(portfolios) == 0 {
		return nil, db.ErrRecordNotFound
	}

	return portfolios, nil
}

func (r *PortfolioRepository) GetPortfolioByID(userID, portfolioID int) (*db.Portfolio, error) {
	var portfolio db.Portfolio
	if err := r.db.QueryRow(
		"SELECT id, user_id, name, description, created_at FROM portfolios WHERE id = $1 AND user_id = $2",
		portfolioID, userID,
	).Scan(&portfolio.ID, &portfolio.UserID, &portfolio.Name, &portfolio.Description, &portfolio.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, db.ErrRecordNotFound
		}
		return nil, err
	}
	return &portfolio, nil
}

func (r *PortfolioRepository) CreatePortfolio(portfolio *db.Portfolio) (*db.Portfolio, error) {
	var id int
	if err := r.db.QueryRow(
		"INSERT INTO portfolios (user_id, name, description) VALUES ($1, $2, $3) RETURNING id, created_at",
		portfolio.UserID, portfolio.Name, portfolio.Description,
	).Scan(&id, &portfolio.CreatedAt); err != nil {
		return nil, err
	}
	portfolio.ID = id
	return portfolio, nil
}

func (r *PortfolioRepository) UpdatePortfolio(portfolio *db.Portfolio) (*db.Portfolio, error) {
	result, err := r.db.Exec(
		"UPDATE portfolios SET name = $1, description = $2 WHERE id = $3 AND user_id = $4",
		portfolio.Name, portfolio.Description, portfolio.ID, portfolio.UserID,
	)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, db.ErrRecordNotFound
	}
	return portfolio, nil
}

func (r *PortfolioRepository) DeletePortfolio(userID, portfolioID int) error {
	result, err := r.db.Exec("DELETE FROM portfolios WHERE id = $1 AND user_id = $2", portfolioID, userID)
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

func (r *PortfolioRepository) GetPortfolioHoldings(portfolioID int) ([]*db.PortfolioItem, error) {
	rows, err := r.db.Query(
		"SELECT id, portfolio_id, stock_id, quantity, avg_cost, created_at FROM portfolio_items WHERE portfolio_id = $1 ORDER BY id",
		portfolioID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var holdings []*db.PortfolioItem
	for rows.Next() {
		var item db.PortfolioItem
		if err := rows.Scan(&item.ID, &item.PortfolioID, &item.StockID, &item.Quantity, &item.AvgCost, &item.CreatedAt); err != nil {
			return nil, err
		}
		holdings = append(holdings, &item)
	}

	if len(holdings) == 0 {
		return nil, db.ErrRecordNotFound
	}

	return holdings, nil
}

func (r *PortfolioRepository) GetHoldingByID(portfolioID, holdingID int) (*db.PortfolioItem, error) {
	var item db.PortfolioItem
	if err := r.db.QueryRow(
		"SELECT id, portfolio_id, stock_id, quantity, avg_cost, created_at FROM portfolio_items WHERE id = $1 AND portfolio_id = $2",
		holdingID, portfolioID,
	).Scan(&item.ID, &item.PortfolioID, &item.StockID, &item.Quantity, &item.AvgCost, &item.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, db.ErrRecordNotFound
		}
		return nil, err
	}
	return &item, nil
}

func (r *PortfolioRepository) CreateHolding(item *db.PortfolioItem) (*db.PortfolioItem, error) {
	var id int
	if err := r.db.QueryRow(
		"INSERT INTO portfolio_items (portfolio_id, stock_id, quantity, avg_cost) VALUES ($1, $2, $3, $4) RETURNING id, created_at",
		item.PortfolioID, item.StockID, item.Quantity, item.AvgCost,
	).Scan(&id, &item.CreatedAt); err != nil {
		return nil, err
	}
	item.ID = id
	return item, nil
}

func (r *PortfolioRepository) UpdateHolding(item *db.PortfolioItem) (*db.PortfolioItem, error) {
	result, err := r.db.Exec(
		"UPDATE portfolio_items SET quantity = $1, avg_cost = $2 WHERE id = $3 AND portfolio_id = $4",
		item.Quantity, item.AvgCost, item.ID, item.PortfolioID,
	)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, db.ErrRecordNotFound
	}
	return item, nil
}

func (r *PortfolioRepository) DeleteHolding(portfolioID, holdingID int) error {
	result, err := r.db.Exec("DELETE FROM portfolio_items WHERE id = $1 AND portfolio_id = $2", holdingID, portfolioID)
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
