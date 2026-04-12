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
	ID             int       `json:"id" db:"id"`
	StockID        int       `json:"stock_id" db:"stock_id"`
	PredictedPrice float64   `json:"predicted_price" db:"predicted_price"`
	Date           time.Time `json:"date" db:"date"`
}

type User struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Portfolio struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
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
