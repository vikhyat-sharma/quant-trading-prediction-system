package repositories

import (
	"database/sql"
	"time"

	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
)

// TaxLotRepository handles tax lot data operations
type TaxLotRepository struct {
	db *sql.DB
}

// NewTaxLotRepository creates a new tax lot repository
func NewTaxLotRepository(database *sql.DB) *TaxLotRepository {
	return &TaxLotRepository{db: database}
}

// CreateTaxLot creates a new tax lot
func (r *TaxLotRepository) CreateTaxLot(taxLot *db.TaxLot) error {
	query := `INSERT INTO tax_lots 
		(portfolio_id, stock_id, quantity, cost_per_share, total_cost, acquisition_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRow(query,
		taxLot.PortfolioID,
		taxLot.StockID,
		taxLot.Quantity,
		taxLot.CostPerShare,
		taxLot.TotalCost,
		taxLot.AcquisitionDate,
		now,
		now,
	).Scan(&taxLot.ID, &taxLot.CreatedAt, &taxLot.UpdatedAt)

	return err
}

// GetTaxLotByID gets a tax lot by ID
func (r *TaxLotRepository) GetTaxLotByID(id int) (*db.TaxLot, error) {
	query := `SELECT id, portfolio_id, stock_id, quantity, cost_per_share, total_cost, 
		acquisition_date, quantity_sold, is_complete, created_at, updated_at
		FROM tax_lots WHERE id = $1`

	taxLot := &db.TaxLot{}
	err := r.db.QueryRow(query, id).Scan(
		&taxLot.ID,
		&taxLot.PortfolioID,
		&taxLot.StockID,
		&taxLot.Quantity,
		&taxLot.CostPerShare,
		&taxLot.TotalCost,
		&taxLot.AcquisitionDate,
		&taxLot.QuantitySold,
		&taxLot.IsComplete,
		&taxLot.CreatedAt,
		&taxLot.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return taxLot, err
}

// GetTaxLotsByPortfolioID gets all tax lots for a portfolio
func (r *TaxLotRepository) GetTaxLotsByPortfolioID(portfolioID int) ([]db.TaxLot, error) {
	query := `SELECT id, portfolio_id, stock_id, quantity, cost_per_share, total_cost, 
		acquisition_date, quantity_sold, is_complete, created_at, updated_at
		FROM tax_lots WHERE portfolio_id = $1 ORDER BY acquisition_date DESC`

	rows, err := r.db.Query(query, portfolioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var taxLots []db.TaxLot
	for rows.Next() {
		taxLot := db.TaxLot{}
		err := rows.Scan(
			&taxLot.ID,
			&taxLot.PortfolioID,
			&taxLot.StockID,
			&taxLot.Quantity,
			&taxLot.CostPerShare,
			&taxLot.TotalCost,
			&taxLot.AcquisitionDate,
			&taxLot.QuantitySold,
			&taxLot.IsComplete,
			&taxLot.CreatedAt,
			&taxLot.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		taxLots = append(taxLots, taxLot)
	}

	return taxLots, rows.Err()
}

// GetActiveTaxLotsByStockID gets active (not fully sold) tax lots for a stock
func (r *TaxLotRepository) GetActiveTaxLotsByStockID(portfolioID, stockID int) ([]db.TaxLot, error) {
	query := `SELECT id, portfolio_id, stock_id, quantity, cost_per_share, total_cost, 
		acquisition_date, quantity_sold, is_complete, created_at, updated_at
		FROM tax_lots 
		WHERE portfolio_id = $1 AND stock_id = $2 AND is_complete = false
		ORDER BY acquisition_date ASC`

	rows, err := r.db.Query(query, portfolioID, stockID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var taxLots []db.TaxLot
	for rows.Next() {
		taxLot := db.TaxLot{}
		err := rows.Scan(
			&taxLot.ID,
			&taxLot.PortfolioID,
			&taxLot.StockID,
			&taxLot.Quantity,
			&taxLot.CostPerShare,
			&taxLot.TotalCost,
			&taxLot.AcquisitionDate,
			&taxLot.QuantitySold,
			&taxLot.IsComplete,
			&taxLot.CreatedAt,
			&taxLot.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		taxLots = append(taxLots, taxLot)
	}

	return taxLots, rows.Err()
}

// UpdateTaxLot updates a tax lot
func (r *TaxLotRepository) UpdateTaxLot(taxLot *db.TaxLot) error {
	query := `UPDATE tax_lots 
		SET quantity_sold = $1, is_complete = $2, updated_at = $3
		WHERE id = $4`

	_, err := r.db.Exec(query,
		taxLot.QuantitySold,
		taxLot.IsComplete,
		time.Now(),
		taxLot.ID,
	)

	return err
}

// DeleteTaxLot deletes a tax lot
func (r *TaxLotRepository) DeleteTaxLot(id int) error {
	query := `DELETE FROM tax_lots WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// CreateTaxTransaction creates a new tax transaction
func (r *TaxLotRepository) CreateTaxTransaction(transaction *db.TaxTransaction) error {
	query := `INSERT INTO tax_transactions 
		(tax_lot_id, portfolio_id, stock_id, type, quantity, price, total_amount, fees, transaction_date, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at`

	now := time.Now()
	err := r.db.QueryRow(query,
		transaction.TaxLotID,
		transaction.PortfolioID,
		transaction.StockID,
		transaction.Type,
		transaction.Quantity,
		transaction.Price,
		transaction.TotalAmount,
		transaction.Fees,
		transaction.TransactionDate,
		now,
	).Scan(&transaction.ID, &transaction.CreatedAt)

	return err
}

// GetTaxTransactionsByPortfolioID gets all tax transactions for a portfolio
func (r *TaxLotRepository) GetTaxTransactionsByPortfolioID(portfolioID int) ([]db.TaxTransaction, error) {
	query := `SELECT id, tax_lot_id, portfolio_id, stock_id, type, quantity, price, 
		total_amount, fees, transaction_date, created_at
		FROM tax_transactions WHERE portfolio_id = $1 ORDER BY transaction_date DESC`

	rows, err := r.db.Query(query, portfolioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []db.TaxTransaction
	for rows.Next() {
		transaction := db.TaxTransaction{}
		var taxLotID sql.NullInt64
		err := rows.Scan(
			&transaction.ID,
			&taxLotID,
			&transaction.PortfolioID,
			&transaction.StockID,
			&transaction.Type,
			&transaction.Quantity,
			&transaction.Price,
			&transaction.TotalAmount,
			&transaction.Fees,
			&transaction.TransactionDate,
			&transaction.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if taxLotID.Valid {
			transaction.TaxLotID = int(taxLotID.Int64)
		}
		transactions = append(transactions, transaction)
	}

	return transactions, rows.Err()
}

// GetTaxTransactionsByTaxLotID gets all transactions for a specific tax lot
func (r *TaxLotRepository) GetTaxTransactionsByTaxLotID(taxLotID int) ([]db.TaxTransaction, error) {
	query := `SELECT id, tax_lot_id, portfolio_id, stock_id, type, quantity, price, 
		total_amount, fees, transaction_date, created_at
		FROM tax_transactions WHERE tax_lot_id = $1 ORDER BY transaction_date ASC`

	rows, err := r.db.Query(query, taxLotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []db.TaxTransaction
	for rows.Next() {
		transaction := db.TaxTransaction{}
		var taxLotID sql.NullInt64
		err := rows.Scan(
			&transaction.ID,
			&taxLotID,
			&transaction.PortfolioID,
			&transaction.StockID,
			&transaction.Type,
			&transaction.Quantity,
			&transaction.Price,
			&transaction.TotalAmount,
			&transaction.Fees,
			&transaction.TransactionDate,
			&transaction.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if taxLotID.Valid {
			transaction.TaxLotID = int(taxLotID.Int64)
		}
		transactions = append(transactions, transaction)
	}

	return transactions, rows.Err()
}
