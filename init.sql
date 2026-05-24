-- Create the stocks table
CREATE TABLE IF NOT EXISTS stocks (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL,
    exchange VARCHAR(10) NOT NULL DEFAULT 'NSE',
    name VARCHAR(255) NOT NULL,
    UNIQUE(symbol, exchange)
);

-- Create the users table with authentication fields
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create the predictions table with algorithm and confidence
CREATE TABLE IF NOT EXISTS predictions (
    id SERIAL PRIMARY KEY,
    stock_id INTEGER NOT NULL REFERENCES stocks(id) ON DELETE CASCADE,
    predicted_price DECIMAL(10,2) NOT NULL,
    algorithm VARCHAR(50) NOT NULL DEFAULT 'ENSEMBLE',
    confidence_score DECIMAL(3,2) NOT NULL DEFAULT 0.5,
    upper_bound DECIMAL(10,2),
    lower_bound DECIMAL(10,2),
    date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create the price history table
CREATE TABLE IF NOT EXISTS price_history (
    id SERIAL PRIMARY KEY,
    stock_id INTEGER NOT NULL REFERENCES stocks(id) ON DELETE CASCADE,
    price DECIMAL(10,2) NOT NULL,
    date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create the portfolios table with performance metrics
CREATE TABLE IF NOT EXISTS portfolios (
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
);

-- Create portfolio holdings table
CREATE TABLE IF NOT EXISTS portfolio_items (
    id SERIAL PRIMARY KEY,
    portfolio_id INTEGER NOT NULL REFERENCES portfolios(id) ON DELETE CASCADE,
    stock_id INTEGER NOT NULL REFERENCES stocks(id) ON DELETE CASCADE,
    quantity DECIMAL(18,4) NOT NULL,
    avg_cost DECIMAL(12,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create alerts table
CREATE TABLE IF NOT EXISTS alerts (
    id SERIAL PRIMARY KEY,
    stock_id INTEGER NOT NULL REFERENCES stocks(id) ON DELETE CASCADE,
    threshold DECIMAL(10,2) NOT NULL,
    condition VARCHAR(50) NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create notifications table
CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    alert_id INTEGER NOT NULL REFERENCES alerts(id) ON DELETE CASCADE,
    stock_id INTEGER NOT NULL REFERENCES stocks(id) ON DELETE CASCADE,
    price DECIMAL(10,2) NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create user watchlists table
CREATE TABLE IF NOT EXISTS user_watchlists (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create watchlist items table
CREATE TABLE IF NOT EXISTS watchlist_items (
    id SERIAL PRIMARY KEY,
    watchlist_id INTEGER NOT NULL REFERENCES user_watchlists(id) ON DELETE CASCADE,
    stock_id INTEGER NOT NULL REFERENCES stocks(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create user alert rules table
CREATE TABLE IF NOT EXISTS user_alert_rules (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    stock_id INTEGER NOT NULL REFERENCES stocks(id) ON DELETE CASCADE,
    threshold DECIMAL(10,2) NOT NULL,
    condition VARCHAR(50) NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for watchlists and alert rules
CREATE INDEX IF NOT EXISTS idx_user_watchlists_user_id ON user_watchlists(user_id);
CREATE INDEX IF NOT EXISTS idx_watchlist_items_watchlist_id ON watchlist_items(watchlist_id);
CREATE INDEX IF NOT EXISTS idx_watchlist_items_stock_id ON watchlist_items(stock_id);
CREATE INDEX IF NOT EXISTS idx_user_alert_rules_user_id ON user_alert_rules(user_id);
CREATE INDEX IF NOT EXISTS idx_user_alert_rules_stock_id ON user_alert_rules(stock_id);

-- Create prediction metrics table for accuracy tracking
CREATE TABLE IF NOT EXISTS prediction_metrics (
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
);

-- Create portfolio performance history table
CREATE TABLE IF NOT EXISTS portfolio_performance (
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
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_predictions_stock_id ON predictions(stock_id);
CREATE INDEX IF NOT EXISTS idx_predictions_date ON predictions(date);
CREATE INDEX IF NOT EXISTS idx_predictions_algorithm ON predictions(algorithm);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_portfolios_user_id ON portfolios(user_id);
CREATE INDEX IF NOT EXISTS idx_portfolio_items_portfolio_id ON portfolio_items(portfolio_id);
CREATE INDEX IF NOT EXISTS idx_portfolio_items_stock_id ON portfolio_items(stock_id);
CREATE INDEX IF NOT EXISTS idx_price_history_stock_id ON price_history(stock_id);
CREATE INDEX IF NOT EXISTS idx_price_history_date ON price_history(date);
CREATE INDEX IF NOT EXISTS idx_alerts_stock_id ON alerts(stock_id);
CREATE INDEX IF NOT EXISTS idx_notifications_alert_id ON notifications(alert_id);
CREATE INDEX IF NOT EXISTS idx_prediction_metrics_stock_id ON prediction_metrics(stock_id);
CREATE INDEX IF NOT EXISTS idx_prediction_metrics_algorithm ON prediction_metrics(algorithm);
CREATE INDEX IF NOT EXISTS idx_prediction_metrics_prediction_id ON prediction_metrics(prediction_id);
CREATE INDEX IF NOT EXISTS idx_portfolio_performance_portfolio_id ON portfolio_performance(portfolio_id);
CREATE INDEX IF NOT EXISTS idx_portfolio_performance_record_date ON portfolio_performance(record_date);
