package db

import (
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

// PriceHistory represents historical price data for a stock
type PriceHistory struct {
	ID        int       `json:"id" db:"id"`
	StockID   int       `json:"stock_id" db:"stock_id"`
	Price     float64   `json:"price" db:"price"`
	Date      time.Time `json:"date" db:"date"`
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
