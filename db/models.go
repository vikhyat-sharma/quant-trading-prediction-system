package db

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/vikhyat-sharma/quant-trading-prediction-system/constants"
)

// Stock represents a stock entity
type Stock struct {
	ID       int    `json:"id" db:"id"`
	Symbol   string `json:"symbol" db:"symbol"`
	Name     string `json:"name" db:"name"`
	Exchange string `json:"exchange" db:"exchange"`
}

// Prediction represents a price prediction for a stock
type Prediction struct {
	ID              int       `json:"id" db:"id"`
	StockID         int       `json:"stock_id" db:"stock_id"`
	PredictedPrice  float64   `json:"predicted_price" db:"predicted_price"`
	Algorithm       string    `json:"algorithm" db:"algorithm"`
	ConfidenceScore float64   `json:"confidence_score" db:"confidence_score"`
	UpperBound      float64   `json:"upper_bound" db:"upper_bound"`
	LowerBound      float64   `json:"lower_bound" db:"lower_bound"`
	Date            time.Time `json:"date" db:"date"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

type User struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	Role      string    `json:"role" db:"role"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Portfolio struct {
	ID            int       `json:"id" db:"id"`
	UserID        int       `json:"user_id" db:"user_id"`
	Name          string    `json:"name" db:"name"`
	Description   string    `json:"description" db:"description"`
	TotalValue    float64   `json:"total_value" db:"total_value"`
	CostBasis     float64   `json:"cost_basis" db:"cost_basis"`
	GainLoss      float64   `json:"gain_loss" db:"gain_loss"`
	ReturnPercent float64   `json:"return_percent" db:"return_percent"`
	LastUpdated   time.Time `json:"last_updated" db:"last_updated"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type PortfolioItem struct {
	ID          int       `json:"id" db:"id"`
	PortfolioID int       `json:"portfolio_id" db:"portfolio_id"`
	StockID     int       `json:"stock_id" db:"stock_id"`
	Quantity    float64   `json:"quantity" db:"quantity"`
	AvgCost     float64   `json:"avg_cost" db:"avg_cost"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// PriceHistory represents historical price data for a stock
type PriceHistory struct {
	ID        int       `json:"id" db:"id"`
	StockID   int       `json:"stock_id" db:"stock_id"`
	Price     float64   `json:"price" db:"price"`
	Date      time.Time `json:"date" db:"date"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

const (
	AlertConditionAbove = "above"
	AlertConditionBelow = "below"
)

type Alert struct {
	ID        int       `json:"id" db:"id"`
	StockID   int       `json:"stock_id" db:"stock_id"`
	Threshold float64   `json:"threshold" db:"threshold"`
	Condition string    `json:"condition" db:"condition"`
	Enabled   bool      `json:"enabled" db:"enabled"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Notification struct {
	ID        int       `json:"id" db:"id"`
	AlertID   int       `json:"alert_id" db:"alert_id"`
	StockID   int       `json:"stock_id" db:"stock_id"`
	Price     float64   `json:"price" db:"price"`
	Message   string    `json:"message" db:"message"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// UserWatchlist represents a user's watchlist
type UserWatchlist struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// WatchlistItem represents a stock in a user's watchlist
type WatchlistItem struct {
	ID          int       `json:"id" db:"id"`
	WatchlistID int       `json:"watchlist_id" db:"watchlist_id"`
	StockID     int       `json:"stock_id" db:"stock_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// UserAlertRule represents a user-specific alert rule
type UserAlertRule struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	StockID   int       `json:"stock_id" db:"stock_id"`
	Threshold float64   `json:"threshold" db:"threshold"`
	Condition string    `json:"condition" db:"condition"`
	Enabled   bool      `json:"enabled" db:"enabled"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Validate checks if the stock data is valid
func (s *Stock) Validate() error {
	if s.Symbol == "" {
		return ErrInvalidSymbol
	}

	validSymbol := regexp.MustCompile(`^[A-Z0-9\.]{1,10}$`)
	if !validSymbol.MatchString(strings.ToUpper(s.Symbol)) {
		return ErrInvalidSymbol
	}

	if s.Name == "" {
		return ErrInvalidName
	}

	exchange := strings.ToUpper(strings.TrimSpace(s.Exchange))
	if exchange != constants.ExchangeNSE && exchange != constants.ExchangeBSE {
		return ErrInvalidExchange
	}

	s.Exchange = exchange
	return nil
}

// Validate checks if the prediction data is valid
func (p *Prediction) Validate() error {
	if p.StockID <= 0 {
		return ErrInvalidStockID
	}
	if p.PredictedPrice < 0 {
		return ErrInvalidPrice
	}
	return nil
}

// Validate checks if the price history data is valid
func (ph *PriceHistory) Validate() error {
	if ph.StockID <= 0 {
		return ErrInvalidStockID
	}
	if ph.Price < 0 {
		return ErrInvalidPrice
	}
	return nil
}

// Validate checks if the alert data is valid
func (a *Alert) Validate() error {
	if a.StockID <= 0 {
		return ErrInvalidStockID
	}
	if a.Threshold < 0 {
		return ErrInvalidThreshold
	}

	condition := strings.ToLower(strings.TrimSpace(a.Condition))
	if condition != AlertConditionAbove && condition != AlertConditionBelow {
		return ErrInvalidCondition
	}
	a.Condition = condition
	return nil
}

// Validate checks if the notification data is valid
func (n *Notification) Validate() error {
	if n.AlertID <= 0 || n.StockID <= 0 {
		return ErrInvalidStockID
	}
	if n.Price < 0 {
		return ErrInvalidPrice
	}
	if n.Message == "" {
		return errors.New("invalid notification message")
	}
	return nil
}

// TaxLot represents a single purchase lot of a stock for tax tracking
type TaxLot struct {
	ID              int       `json:"id" db:"id"`
	PortfolioID     int       `json:"portfolio_id" db:"portfolio_id"`
	StockID         int       `json:"stock_id" db:"stock_id"`
	Quantity        float64   `json:"quantity" db:"quantity"`
	CostPerShare    float64   `json:"cost_per_share" db:"cost_per_share"`
	TotalCost       float64   `json:"total_cost" db:"total_cost"`
	AcquisitionDate time.Time `json:"acquisition_date" db:"acquisition_date"`
	QuantitySold    float64   `json:"quantity_sold" db:"quantity_sold"`
	IsComplete      bool      `json:"is_complete" db:"is_complete"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// TaxTransaction represents a buy/sell transaction for tax tracking
type TaxTransaction struct {
	ID              int       `json:"id" db:"id"`
	TaxLotID        int       `json:"tax_lot_id" db:"tax_lot_id"`
	PortfolioID     int       `json:"portfolio_id" db:"portfolio_id"`
	StockID         int       `json:"stock_id" db:"stock_id"`
	Type            string    `json:"type" db:"type"` // BUY or SELL
	Quantity        float64   `json:"quantity" db:"quantity"`
	Price           float64   `json:"price" db:"price"`
	TotalAmount     float64   `json:"total_amount" db:"total_amount"`
	Fees            float64   `json:"fees" db:"fees"`
	TransactionDate time.Time `json:"transaction_date" db:"transaction_date"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// TaxLotGains represents realized and unrealized gains for a tax lot
type TaxLotGains struct {
	TaxLotID        int       `json:"tax_lot_id"`
	StockID         int       `json:"stock_id"`
	Symbol          string    `json:"symbol"`
	AcquisitionDate time.Time `json:"acquisition_date"`
	QuantityHeld    float64   `json:"quantity_held"`
	QuantitySold    float64   `json:"quantity_sold"`
	CostPerShare    float64   `json:"cost_per_share"`
	CurrentPrice    float64   `json:"current_price"`
	CostBasis       float64   `json:"cost_basis"`
	CurrentValue    float64   `json:"current_value"`
	RealizedGain    float64   `json:"realized_gain"`
	UnrealizedGain  float64   `json:"unrealized_gain"`
	TotalGain       float64   `json:"total_gain"`
	HoldingPeriod   string    `json:"holding_period"` // SHORT_TERM or LONG_TERM
	IsLongTerm      bool      `json:"is_long_term"`
}

// Validate checks if the tax lot data is valid
func (tl *TaxLot) Validate() error {
	if tl.PortfolioID <= 0 || tl.StockID <= 0 {
		return ErrInvalidStockID
	}
	if tl.Quantity <= 0 {
		return errors.New("invalid tax lot quantity")
	}
	if tl.CostPerShare < 0 {
		return errors.New("invalid cost per share")
	}
	return nil
}

// Validate checks if the tax transaction data is valid
func (tt *TaxTransaction) Validate() error {
	if tt.PortfolioID <= 0 || tt.StockID <= 0 {
		return ErrInvalidStockID
	}
	if tt.Quantity <= 0 {
		return errors.New("invalid transaction quantity")
	}
	if tt.Price < 0 || tt.Fees < 0 {
		return errors.New("invalid transaction price or fees")
	}
	tt.TotalAmount = tt.Quantity * tt.Price
	if tt.Type == "" {
		tt.Type = "BUY"
	}
	return nil
}
