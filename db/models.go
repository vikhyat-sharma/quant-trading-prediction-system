package db

import (
	"time"
)

// Stock represents a stock entity
type Stock struct {
	ID     int    `json:"id" db:"id"`
	Symbol string `json:"symbol" db:"symbol"`
	Name   string `json:"name" db:"name"`
}

// Prediction represents a price prediction for a stock
type Prediction struct {
	ID             int       `json:"id" db:"id"`
	StockID        int       `json:"stock_id" db:"stock_id"`
	PredictedPrice float64   `json:"predicted_price" db:"predicted_price"`
	Date           time.Time `json:"date" db:"date"`
}

// Validate checks if the stock data is valid
func (s *Stock) Validate() error {
	if s.Symbol == "" {
		return ErrInvalidSymbol
	}
	if s.Name == "" {
		return ErrInvalidName
	}
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
