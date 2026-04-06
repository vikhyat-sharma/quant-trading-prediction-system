package db

import "errors"

// Common database errors
var (
	ErrInvalidSymbol  = errors.New("invalid stock symbol")
	ErrInvalidName    = errors.New("invalid stock name")
	ErrInvalidStockID = errors.New("invalid stock ID")
	ErrInvalidPrice   = errors.New("invalid predicted price")
	ErrRecordNotFound = errors.New("record not found")
	ErrDatabaseError  = errors.New("database error")
)
