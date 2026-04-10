package db

import "errors"

// Common database errors
var (
	ErrInvalidSymbol    = errors.New("invalid stock symbol")
	ErrInvalidName      = errors.New("invalid stock name")
	ErrInvalidExchange  = errors.New("invalid stock exchange")
	ErrInvalidStockID   = errors.New("invalid stock ID")
	ErrInvalidPrice     = errors.New("invalid predicted price")
	ErrInvalidThreshold = errors.New("invalid alert threshold")
	ErrInvalidCondition = errors.New("invalid alert condition")
	ErrRecordNotFound   = errors.New("record not found")
	ErrDatabaseError    = errors.New("database error")
)
