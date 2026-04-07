package db

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/constants"
)

// DBConfig holds database configuration
type DBConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// NewDB creates a new database connection with optimized settings
func NewDB(url string) (*sql.DB, error) {
	db, err := sql.Open(constants.DatabaseDriverPostgres, url)
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	config := DBConfig{
		MaxOpenConns:    constants.DefaultMaxOpenConns,
		MaxIdleConns:    constants.DefaultMaxIdleConns,
		ConnMaxLifetime: constants.DefaultConnMaxLifetime,
	}

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// NewDBWithConfig creates a database connection with custom configuration
func NewDBWithConfig(url string, config DBConfig) (*sql.DB, error) {
	db, err := sql.Open(constants.DatabaseDriverPostgres, url)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// EnsureSchema creates the required database schema if it doesn't exist
func EnsureSchema(db *sql.DB) error {
	schema := []string{
		`CREATE TABLE IF NOT EXISTS stocks (
		    id SERIAL PRIMARY KEY,
		    symbol VARCHAR(10) NOT NULL UNIQUE,
		    name VARCHAR(255) NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS predictions (
		    id SERIAL PRIMARY KEY,
		    stock_id INTEGER NOT NULL REFERENCES stocks(id) ON DELETE CASCADE,
		    predicted_price DECIMAL(10,2) NOT NULL,
		    date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE INDEX IF NOT EXISTS idx_predictions_stock_id ON predictions(stock_id);`,
		`CREATE INDEX IF NOT EXISTS idx_predictions_date ON predictions(date);`,
	}

	for _, stmt := range schema {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}

	return nil
}
