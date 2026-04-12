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
func EnsureSchema(database *sql.DB) error {
	schema := []string{
		`CREATE TABLE IF NOT EXISTS stocks (
		    id SERIAL PRIMARY KEY,
		    symbol VARCHAR(10) NOT NULL,
		    exchange VARCHAR(10) NOT NULL DEFAULT 'NSE',
		    name VARCHAR(255) NOT NULL,
		    UNIQUE(symbol, exchange)
		);`,
		`CREATE TABLE IF NOT EXISTS predictions (
		    id SERIAL PRIMARY KEY,
		    stock_id INTEGER NOT NULL REFERENCES stocks(id) ON DELETE CASCADE,
		    predicted_price DECIMAL(10,2) NOT NULL,
		    date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS price_history (
		    id SERIAL PRIMARY KEY,
		    stock_id INTEGER NOT NULL REFERENCES stocks(id) ON DELETE CASCADE,
		    price DECIMAL(10,2) NOT NULL,
		    date TIMESTAMP NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS price_alerts (
		    id SERIAL PRIMARY KEY,
		    stock_id INTEGER NOT NULL REFERENCES stocks(id) ON DELETE CASCADE,
		    threshold DECIMAL(10,2) NOT NULL,
		    condition VARCHAR(10) NOT NULL,
		    enabled BOOLEAN NOT NULL DEFAULT TRUE,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS users (
		    id SERIAL PRIMARY KEY,
		    name VARCHAR(255) NOT NULL,
		    email VARCHAR(255) NOT NULL UNIQUE,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS portfolios (
		    id SERIAL PRIMARY KEY,
		    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		    name VARCHAR(255) NOT NULL,
		    description TEXT DEFAULT '',
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS portfolio_items (
		    id SERIAL PRIMARY KEY,
		    portfolio_id INTEGER NOT NULL REFERENCES portfolios(id) ON DELETE CASCADE,
		    stock_id INTEGER NOT NULL REFERENCES stocks(id) ON DELETE CASCADE,
		    quantity DECIMAL(18,4) NOT NULL,
		    avg_cost DECIMAL(12,2) NOT NULL DEFAULT 0,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS notifications (
		    id SERIAL PRIMARY KEY,
		    alert_id INTEGER NOT NULL REFERENCES price_alerts(id) ON DELETE CASCADE,
		    stock_id INTEGER NOT NULL REFERENCES stocks(id) ON DELETE CASCADE,
		    price DECIMAL(10,2) NOT NULL,
		    message TEXT NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE INDEX IF NOT EXISTS idx_predictions_stock_id ON predictions(stock_id);`,
		`CREATE INDEX IF NOT EXISTS idx_predictions_date ON predictions(date);`,
		`CREATE INDEX IF NOT EXISTS idx_price_history_stock_id ON price_history(stock_id);`,
		`CREATE INDEX IF NOT EXISTS idx_price_history_date ON price_history(date);`,
		`CREATE INDEX IF NOT EXISTS idx_stocks_symbol_exchange ON stocks(symbol, exchange);`,
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);`,
		`CREATE INDEX IF NOT EXISTS idx_portfolios_user_id ON portfolios(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_portfolio_items_portfolio_id ON portfolio_items(portfolio_id);`,
	}

	for _, stmt := range schema {
		if _, err := database.Exec(stmt); err != nil {
			return err
		}
	}

	alterations := []string{
		`ALTER TABLE stocks ADD COLUMN IF NOT EXISTS exchange VARCHAR(10) NOT NULL DEFAULT 'NSE';`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_stocks_symbol_exchange ON stocks(symbol, exchange);`,
	}

	for _, stmt := range alterations {
		if _, err := database.Exec(stmt); err != nil {
			return err
		}
	}

	return nil
}
