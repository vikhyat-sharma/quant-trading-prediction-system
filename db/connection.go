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
		`CREATE TABLE IF NOT EXISTS users (
		    id SERIAL PRIMARY KEY,
		    name VARCHAR(255) NOT NULL,
		    email VARCHAR(255) NOT NULL UNIQUE,
		    password VARCHAR(255) NOT NULL,
		    role VARCHAR(50) NOT NULL DEFAULT 'user',
		    is_active BOOLEAN NOT NULL DEFAULT true,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS predictions (
		    id SERIAL PRIMARY KEY,
		    stock_id INTEGER NOT NULL REFERENCES stocks(id) ON DELETE CASCADE,
		    predicted_price DECIMAL(10,2) NOT NULL,
		    algorithm VARCHAR(50) NOT NULL DEFAULT 'ENSEMBLE',
		    confidence_score DECIMAL(3,2) NOT NULL DEFAULT 0.5,
		    upper_bound DECIMAL(10,2),
		    lower_bound DECIMAL(10,2),
		    date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS price_history (
		    id SERIAL PRIMARY KEY,
		    stock_id INTEGER NOT NULL REFERENCES stocks(id) ON DELETE CASCADE,
		    price DECIMAL(10,2) NOT NULL,
		    date TIMESTAMP NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS portfolios (
		    id SERIAL PRIMARY KEY,
		    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		    name VARCHAR(255) NOT NULL,
		    description TEXT DEFAULT '',
		    total_value DECIMAL(15,2) DEFAULT 0,
		    cost_basis DECIMAL(15,2) DEFAULT 0,
		    gain_loss DECIMAL(15,2) DEFAULT 0,
		    return_percent DECIMAL(5,2) DEFAULT 0,
		    last_updated TIMESTAMP,
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
		`CREATE TABLE IF NOT EXISTS alerts (
		    id SERIAL PRIMARY KEY,
		    stock_id INTEGER NOT NULL REFERENCES stocks(id) ON DELETE CASCADE,
		    threshold DECIMAL(10,2) NOT NULL,
		    condition VARCHAR(50) NOT NULL,
		    enabled BOOLEAN NOT NULL DEFAULT TRUE,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS notifications (
		    id SERIAL PRIMARY KEY,
		    alert_id INTEGER NOT NULL REFERENCES alerts(id) ON DELETE CASCADE,
		    stock_id INTEGER NOT NULL REFERENCES stocks(id) ON DELETE CASCADE,
		    price DECIMAL(10,2) NOT NULL,
		    message TEXT NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS user_watchlists (
		    id SERIAL PRIMARY KEY,
		    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		    name VARCHAR(255) NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS watchlist_items (
		    id SERIAL PRIMARY KEY,
		    watchlist_id INTEGER NOT NULL REFERENCES user_watchlists(id) ON DELETE CASCADE,
		    stock_id INTEGER NOT NULL REFERENCES stocks(id) ON DELETE CASCADE,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS user_alert_rules (
		    id SERIAL PRIMARY KEY,
		    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		    stock_id INTEGER NOT NULL REFERENCES stocks(id) ON DELETE CASCADE,
		    threshold DECIMAL(10,2) NOT NULL,
		    condition VARCHAR(50) NOT NULL,
		    enabled BOOLEAN NOT NULL DEFAULT TRUE,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS prediction_metrics (
		    id SERIAL PRIMARY KEY,
		    prediction_id INTEGER NOT NULL REFERENCES predictions(id) ON DELETE CASCADE,
		    stock_id INTEGER NOT NULL REFERENCES stocks(id) ON DELETE CASCADE,
		    algorithm VARCHAR(50) NOT NULL,
		    predicted_price DECIMAL(10,2) NOT NULL,
		    actual_price DECIMAL(10,2),
		    absolute_error DECIMAL(10,2),
		    percent_error DECIMAL(5,2),
		    is_accurate BOOLEAN,
		    accuracy_threshold DECIMAL(5,2) DEFAULT 5.0,
		    prediction_date TIMESTAMP NOT NULL,
		    actual_date TIMESTAMP,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS portfolio_performance (
		    id SERIAL PRIMARY KEY,
		    portfolio_id INTEGER NOT NULL REFERENCES portfolios(id) ON DELETE CASCADE,
		    total_value DECIMAL(15,2) NOT NULL,
		    cost_basis DECIMAL(15,2) NOT NULL,
		    gain_loss DECIMAL(15,2) NOT NULL,
		    return_percent DECIMAL(5,2) NOT NULL,
		    daily_return DECIMAL(5,3),
		    volatility DECIMAL(5,2),
		    sharpe DECIMAL(5,2),
		    max_drawdown DECIMAL(5,2),
		    record_date DATE NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS tax_lots (
		    id SERIAL PRIMARY KEY,
		    portfolio_id INTEGER NOT NULL REFERENCES portfolios(id) ON DELETE CASCADE,
		    stock_id INTEGER NOT NULL REFERENCES stocks(id) ON DELETE CASCADE,
		    quantity DECIMAL(18,4) NOT NULL,
		    cost_per_share DECIMAL(12,2) NOT NULL,
		    total_cost DECIMAL(15,2) NOT NULL,
		    acquisition_date TIMESTAMP NOT NULL,
		    quantity_sold DECIMAL(18,4) NOT NULL DEFAULT 0,
		    is_complete BOOLEAN NOT NULL DEFAULT false,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS tax_transactions (
		    id SERIAL PRIMARY KEY,
		    tax_lot_id INTEGER REFERENCES tax_lots(id) ON DELETE SET NULL,
		    portfolio_id INTEGER NOT NULL REFERENCES portfolios(id) ON DELETE CASCADE,
		    stock_id INTEGER NOT NULL REFERENCES stocks(id) ON DELETE CASCADE,
		    type VARCHAR(10) NOT NULL,
		    quantity DECIMAL(18,4) NOT NULL,
		    price DECIMAL(12,2) NOT NULL,
		    total_amount DECIMAL(15,2) NOT NULL,
		    fees DECIMAL(12,2) DEFAULT 0,
		    transaction_date TIMESTAMP NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE INDEX IF NOT EXISTS idx_predictions_stock_id ON predictions(stock_id);`,
		`CREATE INDEX IF NOT EXISTS idx_predictions_date ON predictions(date);`,
		`CREATE INDEX IF NOT EXISTS idx_predictions_algorithm ON predictions(algorithm);`,
		`CREATE INDEX IF NOT EXISTS idx_price_history_stock_id ON price_history(stock_id);`,
		`CREATE INDEX IF NOT EXISTS idx_price_history_date ON price_history(date);`,
		`CREATE INDEX IF NOT EXISTS idx_stocks_symbol_exchange ON stocks(symbol, exchange);`,
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);`,
		`CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);`,
		`CREATE INDEX IF NOT EXISTS idx_portfolios_user_id ON portfolios(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_portfolio_items_portfolio_id ON portfolio_items(portfolio_id);`,
		`CREATE INDEX IF NOT EXISTS idx_portfolio_items_stock_id ON portfolio_items(stock_id);`,
		`CREATE INDEX IF NOT EXISTS idx_alerts_stock_id ON alerts(stock_id);`,
		`CREATE INDEX IF NOT EXISTS idx_notifications_alert_id ON notifications(alert_id);`,
		`CREATE INDEX IF NOT EXISTS idx_prediction_metrics_stock_id ON prediction_metrics(stock_id);`,
		`CREATE INDEX IF NOT EXISTS idx_prediction_metrics_algorithm ON prediction_metrics(algorithm);`,
		`CREATE INDEX IF NOT EXISTS idx_prediction_metrics_prediction_id ON prediction_metrics(prediction_id);`,
		`CREATE INDEX IF NOT EXISTS idx_portfolio_performance_portfolio_id ON portfolio_performance(portfolio_id);`,
		`CREATE INDEX IF NOT EXISTS idx_portfolio_performance_record_date ON portfolio_performance(record_date);`,
		`CREATE INDEX IF NOT EXISTS idx_user_watchlists_user_id ON user_watchlists(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_watchlist_items_watchlist_id ON watchlist_items(watchlist_id);`,
		`CREATE INDEX IF NOT EXISTS idx_watchlist_items_stock_id ON watchlist_items(stock_id);`,
		`CREATE INDEX IF NOT EXISTS idx_user_alert_rules_user_id ON user_alert_rules(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_user_alert_rules_stock_id ON user_alert_rules(stock_id);`,
		`CREATE INDEX IF NOT EXISTS idx_tax_lots_portfolio_id ON tax_lots(portfolio_id);`,
		`CREATE INDEX IF NOT EXISTS idx_tax_lots_stock_id ON tax_lots(stock_id);`,
		`CREATE INDEX IF NOT EXISTS idx_tax_lots_acquisition_date ON tax_lots(acquisition_date);`,
		`CREATE INDEX IF NOT EXISTS idx_tax_transactions_portfolio_id ON tax_transactions(portfolio_id);`,
		`CREATE INDEX IF NOT EXISTS idx_tax_transactions_stock_id ON tax_transactions(stock_id);`,
		`CREATE INDEX IF NOT EXISTS idx_tax_transactions_tax_lot_id ON tax_transactions(tax_lot_id);`,
		`CREATE INDEX IF NOT EXISTS idx_tax_transactions_type ON tax_transactions(type);`,
		`CREATE INDEX IF NOT EXISTS idx_tax_transactions_transaction_date ON tax_transactions(transaction_date);`,
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
