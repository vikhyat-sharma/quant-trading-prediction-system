package db

import "time"

type Stock struct {
	ID     int    `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

type Prediction struct {
	ID             int       `json:"id"`
	StockID        int       `json:"stock_id"`
	PredictedPrice float64   `json:"predicted_price"`
	Date           time.Time `json:"date"`
}
