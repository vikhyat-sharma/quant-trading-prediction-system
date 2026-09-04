package services

import (
	"database/sql"
	"errors"
	"math"
	"time"

	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
)

// TaxLotRepository defines the methods required by the tax lot service.
type TaxLotRepository interface {
	DB() *sql.DB
	CreateTaxLot(*db.TaxLot) error
	GetTaxLotByID(int) (*db.TaxLot, error)
	GetTaxLotsByPortfolioID(int) ([]db.TaxLot, error)
	GetActiveTaxLotsByStockID(int, int) ([]db.TaxLot, error)
	UpdateTaxLot(*db.TaxLot) error
	CreateTaxTransaction(*db.TaxTransaction) error
	GetTaxTransactionsByPortfolioID(int) ([]db.TaxTransaction, error)
	GetTaxTransactionsByTaxLotID(int) ([]db.TaxTransaction, error)
}

// StockRepository defines the methods required by the tax lot service.
type StockRepository interface {
	GetStock(int) (*db.Stock, error)
}

// TaxLotServiceInterface defines the public methods used by controllers and tests.
type TaxLotServiceInterface interface {
	RecordBuy(portfolioID, stockID int, quantity, price, fees float64, buyDate time.Time) (*db.TaxLot, error)
	RecordSellFIFO(portfolioID, stockID int, quantity, price, fees float64, sellDate time.Time) (float64, error)
	RecordSellLIFO(portfolioID, stockID int, quantity, price, fees float64, sellDate time.Time) (float64, error)
	RecordSellSpecificLot(taxLotID int, quantity, price, fees float64, sellDate time.Time) (float64, error)
	GetTaxLotGains(taxLotID int, currentPrice float64) (*db.TaxLotGains, error)
	GetPortfolioTaxGains(portfolioID int, currentPrices map[int]float64) (map[string]interface{}, error)
	CalculateTaxableGainsBySellDate(portfolioID int) (map[string]interface{}, error)
	GetTaxTransactionsByPortfolio(portfolioID int) ([]db.TaxTransaction, error)
}

// TaxLotService handles tax lot and gains calculations.
type TaxLotService struct {
	taxLotRepo TaxLotRepository
	stockRepo  StockRepository
}

// NewTaxLotService creates a new tax lot service.
func NewTaxLotService(taxLotRepo TaxLotRepository, stockRepo StockRepository) *TaxLotService {
	return &TaxLotService{taxLotRepo: taxLotRepo, stockRepo: stockRepo}
}

// RecordBuy records a buy transaction and creates a tax lot atomically.
func (s *TaxLotService) RecordBuy(portfolioID, stockID int, quantity, price, fees float64, buyDate time.Time) (*db.TaxLot, error) {
	if quantity <= 0 || price < 0 || fees < 0 {
		return nil, errors.New("invalid buy parameters")
	}

	taxLot := &db.TaxLot{
		PortfolioID:     portfolioID,
		StockID:         stockID,
		Quantity:        quantity,
		CostPerShare:    price + (fees / quantity),
		TotalCost:       (quantity * price) + fees,
		AcquisitionDate: buyDate,
	}
	if err := taxLot.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.taxLotRepo.DB().Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = createTaxLotTx(tx, taxLot); err != nil {
		return nil, err
	}

	transaction := &db.TaxTransaction{
		TaxLotID:        taxLot.ID,
		PortfolioID:     portfolioID,
		StockID:         stockID,
		Type:            "BUY",
		Quantity:        quantity,
		Price:           price,
		TotalAmount:     quantity * price,
		Fees:            fees,
		TransactionDate: buyDate,
	}
	if err = createTaxTransactionTx(tx, transaction); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return taxLot, nil
}

// RecordSellFIFO sells using First-In-First-Out method (oldest lots first).
func (s *TaxLotService) RecordSellFIFO(portfolioID, stockID int, quantity, price, fees float64, sellDate time.Time) (float64, error) {
	return s.recordSell(portfolioID, stockID, quantity, price, fees, sellDate, "FIFO")
}

// RecordSellLIFO sells using Last-In-First-Out method (newest lots first).
func (s *TaxLotService) RecordSellLIFO(portfolioID, stockID int, quantity, price, fees float64, sellDate time.Time) (float64, error) {
	return s.recordSell(portfolioID, stockID, quantity, price, fees, sellDate, "LIFO")
}

// RecordSellSpecificLot sells from a specific tax lot atomically.
func (s *TaxLotService) RecordSellSpecificLot(taxLotID int, quantity, price, fees float64, sellDate time.Time) (float64, error) {
	if quantity <= 0 || price < 0 || fees < 0 {
		return 0, errors.New("invalid sell parameters")
	}

	taxLot, err := s.taxLotRepo.GetTaxLotByID(taxLotID)
	if err != nil {
		return 0, err
	}
	if taxLot == nil {
		return 0, errors.New("tax lot not found")
	}

	available := taxLot.Quantity - taxLot.QuantitySold
	if quantity > available {
		return 0, errors.New("insufficient quantity in tax lot")
	}

	tx, err := s.taxLotRepo.DB().Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	taxLot.QuantitySold += quantity
	if taxLot.QuantitySold >= taxLot.Quantity {
		taxLot.IsComplete = true
	}
	if err = updateTaxLotTx(tx, taxLot); err != nil {
		return 0, err
	}

	transaction := &db.TaxTransaction{
		TaxLotID:        taxLotID,
		PortfolioID:     taxLot.PortfolioID,
		StockID:         taxLot.StockID,
		Type:            "SELL",
		Quantity:        quantity,
		Price:           price,
		TotalAmount:     quantity * price,
		Fees:            fees,
		TransactionDate: sellDate,
	}
	if err = createTaxTransactionTx(tx, transaction); err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	costBasis := quantity * taxLot.CostPerShare
	proceeds := (quantity * price) - fees
	return proceeds - costBasis, nil
}

// recordSell is the internal method for FIFO/LIFO sells, fully transactional.
func (s *TaxLotService) recordSell(portfolioID, stockID int, quantity, price, fees float64, sellDate time.Time, method string) (float64, error) {
	if quantity <= 0 || price < 0 || fees < 0 {
		return 0, errors.New("invalid sell parameters")
	}

	activeLots, err := s.taxLotRepo.GetActiveTaxLotsByStockID(portfolioID, stockID)
	if err != nil {
		return 0, err
	}
	if len(activeLots) == 0 {
		return 0, errors.New("no active tax lots found for this stock")
	}

	// Reverse order for LIFO
	if method == "LIFO" {
		for i, j := 0, len(activeLots)-1; i < j; i, j = i+1, j-1 {
			activeLots[i], activeLots[j] = activeLots[j], activeLots[i]
		}
	}

	// Validate sufficient total quantity before opening a transaction
	totalAvailable := 0.0
	for _, lot := range activeLots {
		totalAvailable += lot.Quantity - lot.QuantitySold
	}
	if quantity > totalAvailable {
		return 0, errors.New("insufficient quantity across all tax lots")
	}

	tx, err := s.taxLotRepo.DB().Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	totalRealizedGain := 0.0
	remaining := quantity

	for i := range activeLots {
		if remaining <= 0 {
			break
		}
		lot := &activeLots[i]
		available := lot.Quantity - lot.QuantitySold
		if available <= 0 {
			continue
		}

		toSell := math.Min(remaining, available)
		lot.QuantitySold += toSell
		if lot.QuantitySold >= lot.Quantity {
			lot.IsComplete = true
		}

		if err = updateTaxLotTx(tx, lot); err != nil {
			return 0, err
		}

		proportionalFees := (toSell / quantity) * fees
		transaction := &db.TaxTransaction{
			TaxLotID:        lot.ID,
			PortfolioID:     portfolioID,
			StockID:         stockID,
			Type:            "SELL",
			Quantity:        toSell,
			Price:           price,
			TotalAmount:     toSell * price,
			Fees:            proportionalFees,
			TransactionDate: sellDate,
		}
		if err = createTaxTransactionTx(tx, transaction); err != nil {
			return 0, err
		}

		costBasis := toSell * lot.CostPerShare
		proceeds := (toSell * price) - proportionalFees
		totalRealizedGain += proceeds - costBasis
		remaining -= toSell
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return totalRealizedGain, nil
}

// GetTaxLotGains calculates realized and unrealized gains for a tax lot.
func (s *TaxLotService) GetTaxLotGains(taxLotID int, currentPrice float64) (*db.TaxLotGains, error) {
	taxLot, err := s.taxLotRepo.GetTaxLotByID(taxLotID)
	if err != nil {
		return nil, err
	}
	if taxLot == nil {
		return nil, errors.New("tax lot not found")
	}

	stock, err := s.stockRepo.GetStock(taxLot.StockID)
	if err != nil {
		return nil, err
	}
	if stock == nil {
		return nil, errors.New("stock not found")
	}

	quantityHeld := taxLot.Quantity - taxLot.QuantitySold

	// Calculate realized gain from actual sell transactions
	realizedGain := 0.0
	if taxLot.QuantitySold > 0 {
		transactions, err := s.taxLotRepo.GetTaxTransactionsByTaxLotID(taxLot.ID)
		if err != nil {
			return nil, err
		}
		var proceedsSold, feesPaid float64
		for _, tx := range transactions {
			if tx.Type != "SELL" {
				continue
			}
			proceeds := tx.TotalAmount
			if proceeds <= 0 {
				proceeds = tx.Quantity * tx.Price
			}
			proceedsSold += proceeds
			feesPaid += tx.Fees
		}
		costBasisSold := taxLot.QuantitySold * taxLot.CostPerShare
		realizedGain = (proceedsSold - feesPaid) - costBasisSold
	}

	costBasisHeld := quantityHeld * taxLot.CostPerShare
	currentValue := quantityHeld * currentPrice
	unrealizedGain := currentValue - costBasisHeld

	holdingDays := time.Since(taxLot.AcquisitionDate).Hours() / 24
	holdingPeriod := "SHORT_TERM"
	isLongTerm := false
	if holdingDays > 365 {
		holdingPeriod = "LONG_TERM"
		isLongTerm = true
	}

	return &db.TaxLotGains{
		TaxLotID:        taxLot.ID,
		StockID:         taxLot.StockID,
		Symbol:          stock.Symbol,
		AcquisitionDate: taxLot.AcquisitionDate,
		QuantityHeld:    quantityHeld,
		QuantitySold:    taxLot.QuantitySold,
		CostPerShare:    taxLot.CostPerShare,
		CurrentPrice:    currentPrice,
		CostBasis:       taxLot.TotalCost,
		CurrentValue:    currentValue,
		RealizedGain:    realizedGain,
		UnrealizedGain:  unrealizedGain,
		TotalGain:       realizedGain + unrealizedGain,
		HoldingPeriod:   holdingPeriod,
		IsLongTerm:      isLongTerm,
	}, nil
}

// GetPortfolioTaxGains calculates total realized and unrealized gains for a portfolio.
func (s *TaxLotService) GetPortfolioTaxGains(portfolioID int, currentPrices map[int]float64) (map[string]interface{}, error) {
	taxLots, err := s.taxLotRepo.GetTaxLotsByPortfolioID(portfolioID)
	if err != nil {
		return nil, err
	}

	var (
		totalRealizedGain   float64
		totalUnrealizedGain float64
		totalCostBasis      float64
		totalCurrentValue   float64
		longTermGain        float64
		shortTermGain       float64
		allGains            []db.TaxLotGains
	)

	for _, lot := range taxLots {
		price, ok := currentPrices[lot.StockID]
		if !ok || price <= 0 {
			// Skip lots with no price rather than silently corrupting totals;
			// callers should be aware that totals are partial if prices are missing.
			continue
		}

		gains, err := s.GetTaxLotGains(lot.ID, price)
		if err != nil {
			continue
		}

		totalRealizedGain += gains.RealizedGain
		totalUnrealizedGain += gains.UnrealizedGain
		totalCostBasis += gains.CostBasis
		totalCurrentValue += gains.CurrentValue

		if gains.IsLongTerm {
			longTermGain += gains.TotalGain
		} else {
			shortTermGain += gains.TotalGain
		}
		allGains = append(allGains, *gains)
	}

	return map[string]interface{}{
		"total_realized_gain":   totalRealizedGain,
		"total_unrealized_gain": totalUnrealizedGain,
		"total_gain":            totalRealizedGain + totalUnrealizedGain,
		"total_cost_basis":      totalCostBasis,
		"total_current_value":   totalCurrentValue,
		"long_term_gain":        longTermGain,
		"short_term_gain":       shortTermGain,
		"tax_lot_count":         len(allGains),
		"tax_lots":              allGains,
	}, nil
}

// GetTaxTransactionsByPortfolio gets all transactions for a portfolio.
func (s *TaxLotService) GetTaxTransactionsByPortfolio(portfolioID int) ([]db.TaxTransaction, error) {
	return s.taxLotRepo.GetTaxTransactionsByPortfolioID(portfolioID)
}

// CalculateTaxableGainsBySellDate calculates tax consequences by sell date.
func (s *TaxLotService) CalculateTaxableGainsBySellDate(portfolioID int) (map[string]interface{}, error) {
	transactions, err := s.taxLotRepo.GetTaxTransactionsByPortfolioID(portfolioID)
	if err != nil {
		return nil, err
	}

	var shortTermGains, longTermGains float64
	var sellCount int

	for _, sellTx := range transactions {
		if sellTx.Type != "SELL" {
			continue
		}
		if sellTx.TaxLotID == 0 {
			// Cannot determine holding period without a lot reference; skip.
			continue
		}
		sellCount++

		taxLot, err := s.taxLotRepo.GetTaxLotByID(sellTx.TaxLotID)
		if err != nil || taxLot == nil {
			continue
		}

		costBasis := sellTx.Quantity * taxLot.CostPerShare
		proceeds := (sellTx.Quantity * sellTx.Price) - sellTx.Fees
		gain := proceeds - costBasis

		holdingDays := sellTx.TransactionDate.Sub(taxLot.AcquisitionDate).Hours() / 24
		if holdingDays > 365 {
			longTermGains += gain
		} else {
			shortTermGains += gain
		}
	}

	return map[string]interface{}{
		"short_term_gains":  shortTermGains,
		"long_term_gains":   longTermGains,
		"total_gains":       shortTermGains + longTermGains,
		"sell_transactions": sellCount,
	}, nil
}

// ── Transactional helpers ─────────────────────────────────────────────────────

func createTaxLotTx(tx *sql.Tx, taxLot *db.TaxLot) error {
	query := `INSERT INTO tax_lots
		(portfolio_id, stock_id, quantity, cost_per_share, total_cost, acquisition_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, created_at, updated_at`
	return tx.QueryRow(query,
		taxLot.PortfolioID, taxLot.StockID, taxLot.Quantity,
		taxLot.CostPerShare, taxLot.TotalCost, taxLot.AcquisitionDate,
	).Scan(&taxLot.ID, &taxLot.CreatedAt, &taxLot.UpdatedAt)
}

func updateTaxLotTx(tx *sql.Tx, taxLot *db.TaxLot) error {
	_, err := tx.Exec(
		`UPDATE tax_lots SET quantity_sold = $1, is_complete = $2, updated_at = NOW() WHERE id = $3`,
		taxLot.QuantitySold, taxLot.IsComplete, taxLot.ID,
	)
	return err
}

func createTaxTransactionTx(tx *sql.Tx, t *db.TaxTransaction) error {
	query := `INSERT INTO tax_transactions
		(tax_lot_id, portfolio_id, stock_id, type, quantity, price, total_amount, fees, transaction_date, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
		RETURNING id, created_at`
	return tx.QueryRow(query,
		t.TaxLotID, t.PortfolioID, t.StockID, t.Type,
		t.Quantity, t.Price, t.TotalAmount, t.Fees, t.TransactionDate,
	).Scan(&t.ID, &t.CreatedAt)
}
