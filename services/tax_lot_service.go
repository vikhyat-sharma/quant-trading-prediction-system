package services

import (
	"errors"
	"math"
	"time"

	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
)

// TaxLotRepository defines the methods required by the tax lot service.
type TaxLotRepository interface {
	CreateTaxLot(*db.TaxLot) error
	GetTaxLotByID(int) (*db.TaxLot, error)
	GetTaxLotsByPortfolioID(int) ([]db.TaxLot, error)
	GetActiveTaxLotsByStockID(int, int) ([]db.TaxLot, error)
	UpdateTaxLot(*db.TaxLot) error
	CreateTaxTransaction(*db.TaxTransaction) error
	GetTaxTransactionsByPortfolioID(int) ([]db.TaxTransaction, error)
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

// TaxLotService handles tax lot and gains calculations
type TaxLotService struct {
	taxLotRepo TaxLotRepository
	stockRepo  StockRepository
}

// NewTaxLotService creates a new tax lot service
func NewTaxLotService(
	taxLotRepo TaxLotRepository,
	stockRepo StockRepository,
) *TaxLotService {
	return &TaxLotService{
		taxLotRepo: taxLotRepo,
		stockRepo:  stockRepo,
	}
}

// RecordBuy records a buy transaction and creates a tax lot
func (s *TaxLotService) RecordBuy(portfolioID, stockID int, quantity, price, fees float64, buyDate time.Time) (*db.TaxLot, error) {
	if quantity <= 0 || price < 0 || fees < 0 {
		return nil, errors.New("invalid buy parameters")
	}

	totalCost := (quantity * price) + fees

	taxLot := &db.TaxLot{
		PortfolioID:     portfolioID,
		StockID:         stockID,
		Quantity:        quantity,
		CostPerShare:    price + (fees / quantity),
		TotalCost:       totalCost,
		AcquisitionDate: buyDate,
		QuantitySold:    0,
		IsComplete:      false,
	}

	if err := taxLot.Validate(); err != nil {
		return nil, err
	}

	if err := s.taxLotRepo.CreateTaxLot(taxLot); err != nil {
		return nil, err
	}

	// Record the transaction
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

	if err := s.taxLotRepo.CreateTaxTransaction(transaction); err != nil {
		return nil, err
	}

	return taxLot, nil
}

// RecordSellFIFO sells using First-In-First-Out method (oldest lots first)
func (s *TaxLotService) RecordSellFIFO(portfolioID, stockID int, quantity, price, fees float64, sellDate time.Time) (float64, error) {
	return s.recordSell(portfolioID, stockID, quantity, price, fees, sellDate, "FIFO")
}

// RecordSellLIFO sells using Last-In-First-Out method (newest lots first)
func (s *TaxLotService) RecordSellLIFO(portfolioID, stockID int, quantity, price, fees float64, sellDate time.Time) (float64, error) {
	return s.recordSell(portfolioID, stockID, quantity, price, fees, sellDate, "LIFO")
}

// RecordSellSpecificLot sells from a specific tax lot
func (s *TaxLotService) RecordSellSpecificLot(taxLotID int, quantity, price, fees float64, sellDate time.Time) (float64, error) {
	taxLot, err := s.taxLotRepo.GetTaxLotByID(taxLotID)
	if err != nil {
		return 0, err
	}
	if taxLot == nil {
		return 0, errors.New("tax lot not found")
	}

	quantityAvailable := taxLot.Quantity - taxLot.QuantitySold
	if quantity > quantityAvailable {
		return 0, errors.New("insufficient quantity in tax lot")
	}

	taxLot.QuantitySold += quantity
	if taxLot.QuantitySold >= taxLot.Quantity {
		taxLot.IsComplete = true
	}

	if err := s.taxLotRepo.UpdateTaxLot(taxLot); err != nil {
		return 0, err
	}

	// Record the transaction
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

	if err := s.taxLotRepo.CreateTaxTransaction(transaction); err != nil {
		return 0, err
	}

	// Calculate realized gain
	costBasis := quantity * taxLot.CostPerShare
	proceeds := (quantity * price) - fees
	realizedGain := proceeds - costBasis

	return realizedGain, nil
}

// recordSell is the internal method for selling with FIFO or LIFO strategy
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
		for i := len(activeLots)/2 - 1; i >= 0; i-- {
			opp := len(activeLots) - 1 - i
			activeLots[i], activeLots[opp] = activeLots[opp], activeLots[i]
		}
	}

	totalRealizedGain := 0.0
	remainingQuantity := quantity
	totalFees := fees

	for _, lot := range activeLots {
		if remainingQuantity <= 0 {
			break
		}

		quantityAvailable := lot.Quantity - lot.QuantitySold
		if quantityAvailable <= 0 {
			continue
		}

		quantityToSell := math.Min(remainingQuantity, quantityAvailable)
		lot.QuantitySold += quantityToSell

		if lot.QuantitySold >= lot.Quantity {
			lot.IsComplete = true
		}

		if err := s.taxLotRepo.UpdateTaxLot(&lot); err != nil {
			return 0, err
		}

		// Record the transaction
		proportionalFees := (quantityToSell / quantity) * totalFees
		transaction := &db.TaxTransaction{
			TaxLotID:        lot.ID,
			PortfolioID:     portfolioID,
			StockID:         stockID,
			Type:            "SELL",
			Quantity:        quantityToSell,
			Price:           price,
			TotalAmount:     quantityToSell * price,
			Fees:            proportionalFees,
			TransactionDate: sellDate,
		}

		if err := s.taxLotRepo.CreateTaxTransaction(transaction); err != nil {
			return 0, err
		}

		// Calculate realized gain for this portion
		costBasis := quantityToSell * lot.CostPerShare
		proceeds := (quantityToSell * price) - proportionalFees
		realizedGain := proceeds - costBasis
		totalRealizedGain += realizedGain

		remainingQuantity -= quantityToSell
	}

	if remainingQuantity > 0 {
		return 0, errors.New("insufficient quantity across all tax lots")
	}

	return totalRealizedGain, nil
}

// GetTaxLotGains calculates realized and unrealized gains for a tax lot
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

	// Calculate realized gain (from sold quantity)
	costBasisSold := taxLot.QuantitySold * taxLot.CostPerShare
	proceedsSold := taxLot.QuantitySold * currentPrice
	realizedGain := proceedsSold - costBasisSold

	// Calculate unrealized gain (from held quantity)
	costBasisHeld := quantityHeld * taxLot.CostPerShare
	currentValue := quantityHeld * currentPrice
	unrealizedGain := currentValue - costBasisHeld

	// Determine holding period (long-term if > 1 year)
	holdingDays := time.Since(taxLot.AcquisitionDate).Hours() / 24
	holdingPeriod := "SHORT_TERM"
	isLongTerm := false
	if holdingDays > 365 {
		holdingPeriod = "LONG_TERM"
		isLongTerm = true
	}

	gains := &db.TaxLotGains{
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
	}

	return gains, nil
}

// GetPortfolioTaxGains calculates total realized and unrealized gains for a portfolio
func (s *TaxLotService) GetPortfolioTaxGains(portfolioID int, currentPrices map[int]float64) (map[string]interface{}, error) {
	taxLots, err := s.taxLotRepo.GetTaxLotsByPortfolioID(portfolioID)
	if err != nil {
		return nil, err
	}

	totalRealizedGain := 0.0
	totalUnrealizedGain := 0.0
	totalCostBasis := 0.0
	totalCurrentValue := 0.0
	longTermGain := 0.0
	shortTermGain := 0.0
	var allTaxLotGains []db.TaxLotGains

	for _, lot := range taxLots {
		price := currentPrices[lot.StockID]
		if price == 0 {
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

		allTaxLotGains = append(allTaxLotGains, *gains)
	}

	return map[string]interface{}{
		"total_realized_gain":   totalRealizedGain,
		"total_unrealized_gain": totalUnrealizedGain,
		"total_gain":            totalRealizedGain + totalUnrealizedGain,
		"total_cost_basis":      totalCostBasis,
		"total_current_value":   totalCurrentValue,
		"long_term_gain":        longTermGain,
		"short_term_gain":       shortTermGain,
		"tax_lot_count":         len(allTaxLotGains),
		"tax_lots":              allTaxLotGains,
	}, nil
}

// GetTaxTransactionsByPortfolio gets all transactions for a portfolio
func (s *TaxLotService) GetTaxTransactionsByPortfolio(portfolioID int) ([]db.TaxTransaction, error) {
	return s.taxLotRepo.GetTaxTransactionsByPortfolioID(portfolioID)
}

// CalculateTaxableGainsBySellDate calculates tax consequences by sell date
func (s *TaxLotService) CalculateTaxableGainsBySellDate(portfolioID int) (map[string]interface{}, error) {
	transactions, err := s.taxLotRepo.GetTaxTransactionsByPortfolioID(portfolioID)
	if err != nil {
		return nil, err
	}

	sellTransactions := make([]db.TaxTransaction, 0)
	for _, t := range transactions {
		if t.Type == "SELL" {
			sellTransactions = append(sellTransactions, t)
		}
	}

	shortTermGains := 0.0
	longTermGains := 0.0

	for _, sellTx := range sellTransactions {
		if sellTx.TaxLotID == 0 {
			continue
		}

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
		"sell_transactions": len(sellTransactions),
	}, nil
}
