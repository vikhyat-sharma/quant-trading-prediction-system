# Advanced Features Implementation Guide

## Overview

This document outlines all the advanced features added to the Quant Trading Prediction System, including prediction algorithms, authentication, analytics, and metrics.

---

## 1. Advanced Prediction Algorithms

### Implemented Algorithms

#### 1.1 Simple Moving Average (SMA)
- **Package**: `services/algorithms/prediction_algorithms.go`
- **Function**: `SimpleMovingAveragePrediction()`
- **Description**: Uses 20-day and 10-day SMAs to predict trends
- **Confidence Score**: 0.5-0.6
- **Use Case**: Identifying short-term price trends

#### 1.2 Exponential Moving Average (EMA)
- **Function**: `ExponentialMovingAveragePrediction()`
- **Description**: Uses 12-day and 5-day EMAs for more responsive predictions
- **Confidence Score**: 0.55-0.65
- **Use Case**: Faster response to price changes

#### 1.3 Momentum Trading
- **Function**: `MomentumPrediction()`
- **Description**: Uses 14-day RSI and 5-day momentum
- **Confidence Score**: 0.55-0.6
- **Features**:
  - Detects overbought (RSI > 70) and oversold (RSI < 30) conditions
  - Predicts reversals in extreme conditions
  - More aggressive in volatility detection

#### 1.4 Mean Reversion
- **Function**: `MeanReversionPrediction()`
- **Description**: Reverts prices toward the long-term mean
- **Confidence Score**: 0.5-0.65
- **Features**:
  - Identifies overvalued/undervalued conditions
  - Predicts gradual price correction
  - Good for ranging markets

#### 1.5 Ensemble (Recommended)
- **Function**: `EnsemblePrediction()`
- **Description**: Combines all four algorithms with weighted averaging
- **Weights**:
  - EMA: 30% (most responsive)
  - SMA: 25%
  - Momentum: 25%
  - Mean Reversion: 20%
- **Confidence Score**: Averaged from all models
- **Advantage**: Reduces model bias, more robust predictions

### Technical Indicators

The system calculates the following technical indicators:

```go
type TechnicalIndicators struct {
    SMA20      float64 // 20-day Simple Moving Average
    EMA12      float64 // 12-day Exponential Moving Average
    RSI        float64 // Relative Strength Index (0-100)
    MACD       float64 // MACD line
    MACDSignal float64 // MACD signal line
    BBUpper    float64 // Bollinger Band Upper
    BBLower    float64 // Bollinger Band Lower
}
```

**Usage**:
```go
indicators := algorithms.CalculateTechnicalIndicators(prices)
```

### API Usage

**Generate prediction with default algorithm (ENSEMBLE)**:
```bash
POST /stocks/{stockID}/predictions
```

**Generate prediction with specific algorithm**:
```bash
POST /stocks/{stockID}/predictions?algorithm=SMA
```

Supported algorithms: `SMA`, `EMA`, `MOMENTUM`, `MEAN_REVERSION`, `ENSEMBLE`

---

## 2. Authentication & Authorization

### User Model Updates

```go
type User struct {
    ID        int
    Name      string
    Email     string
    Password  string    // Hashed with bcrypt
    Role      string    // "user", "admin", "superadmin"
    IsActive  bool
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### Available Roles

- **user**: Regular trader with full trading access
- **admin**: Can manage users, view analytics, create alerts
- **superadmin**: Full system access

### JWT Token

**Generated using**: HS256 algorithm
**Expiry**: 7 days (configurable via `TokenExpiry`)
**Claims**:
```go
type Claims struct {
    UserID int
    Email  string
    Role   string
    jwt.RegisteredClaims
}
```

### Utility Functions (`util/auth.go`)

#### Password Management
```go
// Hash password (min 8 characters)
hash, err := util.HashPassword("mypassword123")

// Verify password
isValid := util.VerifyPassword(hash, "mypassword123")
```

#### JWT Token
```go
// Generate JWT
token, err := util.GenerateJWT(userID, email, role)

// Verify JWT
claims, err := util.VerifyJWT(tokenString)
```

#### Validation
```go
// Validate email format
if !util.ValidateEmail("user@example.com") {
    // Invalid email
}
```

### Middleware

#### AuthMiddleware
Protects routes - requires valid JWT token:
```bash
curl -H "Authorization: Bearer <token>" http://localhost:8080/portfolios
```

#### OptionalAuthMiddleware
Allows both authenticated and unauthenticated requests

#### AdminMiddleware
Restricts access to admin/superadmin users only

#### RateLimitMiddleware
Prevents abuse:
```go
rateLimiter := middleware.NewRateLimiter(100, time.Minute)
// Max 100 requests per minute per IP
```

---

## 3. Prediction Accuracy Metrics

### PredictionMetric Model

```go
type PredictionMetric struct {
    ID                int
    PredictionID      int
    StockID           int
    Algorithm         string
    PredictedPrice    float64
    ActualPrice       *float64
    AbsoluteError     *float64
    PercentError      *float64
    IsAccurate        *bool
    AccuracyThreshold float64  // Default: 5%
    PredictionDate    time.Time
    ActualDate        *time.Time
    CreatedAt         time.Time
    UpdatedAt         time.Time
}
```

### Tracking Predictions

1. **Create Metric** (when prediction is made):
   ```go
   metric := &db.PredictionMetric{
       PredictionID:      pred.ID,
       StockID:           pred.StockID,
       Algorithm:         pred.Algorithm,
       PredictedPrice:    pred.PredictedPrice,
       AccuracyThreshold: 5.0, // 5% error tolerance
       PredictionDate:    pred.Date,
   }
   createdMetric, err := metricsRepo.CreateMetric(metric)
   ```

2. **Update with Actual Price**:
   ```go
   updatedMetric, err := metricsRepo.UpdateMetricWithActualPrice(
       metricID, 
       actualPrice, 
       time.Now(),
   )
   ```

### Retrieving Metrics

**Get individual prediction metrics**:
```bash
GET /metrics/predictions/{predictionID}
```

**Get algorithm performance**:
```bash
GET /metrics/algorithms/{algorithmName}
```

Response:
```json
{
  "algorithm": "ENSEMBLE",
  "total_predictions": 150,
  "accurate_predictions": 127,
  "accuracy_rate": 84.67,
  "average_absolute_error": 2.34,
  "average_percent_error": 2.1,
  "average_confidence_score": 0.72,
  "last_updated": "2024-05-12T10:30:00Z"
}
```

**Get all algorithm performance**:
```bash
GET /metrics/algorithms
```

### Performance Tracking

The system tracks:
- **Total Predictions**: Number of predictions made
- **Accuracy Rate**: % of predictions within threshold
- **Absolute Error**: Average difference in price
- **Percent Error**: Average % difference
- **Confidence Score**: Average confidence of predictions

---

## 4. Portfolio Analytics

### PortfolioAnalyticsService

Provides comprehensive portfolio analysis:


#### Metrics Calculated

```go
type PortfolioAnalytics struct {
    TotalValue           float64 // Current portfolio value
    CostBasis            float64 // Total amount invested
    GainLoss             float64 // Profit/Loss in absolute terms
    ReturnPercent        float64 // Return as percentage
    Volatility           float64 // Annualized volatility
    SharpeRatio          float64 // Risk-adjusted return
    MaxDrawdown          float64 // Largest peak-to-trough decline
    DiversificationScore float64 // 0-1 (0 = concentrated, 1 = diversified)
    ValueAtRisk          float64 // Value at Risk (VaR)
    ExpectedShortfall    float64 // Expected Shortfall (ES/CVaR)
}
```

##### Value at Risk (VaR) & Expected Shortfall (ES)

- **Value at Risk (VaR)**: Maximum expected portfolio loss over a given period at a specified confidence level (e.g., 95%).
- **Expected Shortfall (ES/CVaR)**: Average loss in the worst-case (tail) scenarios beyond the VaR threshold.

**Calculation:**
Implemented in `PortfolioAnalyticsService`:
```go
VaR, _ := analyticsService.CalculateValueAtRisk(portfolioID, 100, 0.95) // 100-day, 95% VaR
ES, _ := analyticsService.CalculateExpectedShortfall(portfolioID, 100, 0.95)
```

**API Usage Example:**
```bash
# Get 95% 30-day Value at Risk
GET /users/{userID}/portfolios/{portfolioID}/value-at-risk?days=30&confidence=0.95

# Get 95% 30-day Expected Shortfall
GET /users/{userID}/portfolios/{portfolioID}/expected-shortfall?days=30&confidence=0.95
```

**Response:**
```json
{
  "value_at_risk": 0.032, // 3.2% daily VaR
  "expected_shortfall": 0.045 // 4.5% daily ES
}
```

**Interpretation:**
- VaR: There is a 95% chance the portfolio will not lose more than 3.2% in a day.
- ES: If the loss exceeds VaR, the average loss is 4.5%.

#### API Endpoints

**Get portfolio value**:
```bash
GET /users/{userID}/portfolios/{portfolioID}/value
```

**Get portfolio metrics**:
```bash
GET /users/{userID}/portfolios/{portfolioID}/metrics
```

Response:
```json
{
  "total_value": 125000.50,
  "cost_basis": 100000.00,
  "gain_loss": 25000.50,
  "return_percent": 25.0,
  "volatility": 0.18,
  "sharpe_ratio": 1.42,
  "max_drawdown": -8.5,
  "diversification_score": 0.72
}
```

**Get portfolio volatility** (30-day annualized):
```bash
GET /users/{userID}/portfolios/{portfolioID}/volatility?days=30
```

**Get Sharpe ratio**:
```bash
GET /users/{userID}/portfolios/{portfolioID}/sharpe-ratio?days=30&risk_free_rate=0.02
```

**Get top holdings**:
```bash
GET /users/{userID}/portfolios/{portfolioID}/top-holdings?limit=5
```

Response:
```json
{
  "holdings": [
    {
      "stock_id": 1,
      "quantity": 100,
      "avg_cost": 150.00,
      "value": 18500.00,
      "weight": 14.8
    }
  ]
}
```

**Get diversification score**:
```bash
GET /users/{userID}/portfolios/{portfolioID}/diversification
```

### Portfolio Performance Tracking

The system records daily portfolio snapshots:

```go
type PortfolioPerformance struct {
    PortfolioID   int
    TotalValue    float64
    CostBasis     float64
    GainLoss      float64
    ReturnPercent float64
    DailyReturn   float64
    Volatility    float64
    Sharpe        float64
    MaxDrawdown   float64
    RecordDate    time.Time
}
```

---

## 5. Alerts & Notifications System

### Enhanced Features

**Existing**: Price alert triggers
**New**: 
- Multiple alert conditions
- Historical notification tracking
- Real-time evaluation

### Alert Creation

```bash
POST /stocks/{stockID}/alerts
{
    "threshold": 150.50,
    "condition": "above",  // or "below"
    "enabled": true
}
```

### Notification Evaluation

Alerts are evaluated periodically:
```bash
POST /stocks/{stockID}/evaluate-alerts
```

### Get Notifications

```bash
GET /stocks/{stockID}/notifications
```

---

## 6. Database Schema Updates

### New Tables

#### prediction_metrics
Tracks accuracy and performance of each prediction

#### portfolio_performance
Records daily portfolio metrics for historical analysis

### Updated Tables

#### users
- Added: `password`, `role`, `is_active`, `updated_at`

#### predictions
- Added: `algorithm`, `confidence_score`, `upper_bound`, `lower_bound`, `created_at`

#### portfolios
- Added: `total_value`, `cost_basis`, `gain_loss`, `return_percent`, `last_updated`

### New Indexes

For optimal performance:
- `idx_predictions_algorithm`
- `idx_users_role`
- `idx_prediction_metrics_stock_id`
- `idx_prediction_metrics_algorithm`
- `idx_prediction_metrics_prediction_id`
- `idx_portfolio_performance_portfolio_id`
- `idx_portfolio_performance_record_date`

---

## 7. Configuration

### Environment Variables

```bash
# Database
DATABASE_URL=postgresql://user:password@localhost:5432/quant_trading

# Server
PORT=8080

# JWT
JWT_SECRET_KEY=QUANT_TRADING_SECRET_KEY (set in code)
JWT_EXPIRY=168h (7 days)

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60 (seconds)
```

---

## 8. Usage Examples

### Example 1: Generate Ensemble Prediction

```bash
# Generate prediction using ENSEMBLE algorithm
POST /stocks/1/predictions?algorithm=ENSEMBLE

# Response:
{
  "id": 456,
  "stock_id": 1,
  "predicted_price": 155.25,
  "algorithm": "ENSEMBLE",
  "confidence_score": 0.68,
  "upper_bound": 163.02,
  "lower_bound": 147.48,
  "date": "2024-05-13T00:00:00Z",
  "created_at": "2024-05-12T10:30:00Z"
}
```

### Example 2: Track Prediction Accuracy

```bash
# Create metric when prediction is made (automatic)
# Then update when actual price is known

PUT /metrics/predictions/456
{
  "actual_price": 156.80,
  "actual_date": "2024-05-13T15:30:00Z"
}

# Response shows accuracy calculation:
{
  "id": 789,
  "prediction_id": 456,
  "predicted_price": 155.25,
  "actual_price": 156.80,
  "absolute_error": 1.55,
  "percent_error": 0.99,
  "is_accurate": true,  // Within 5% threshold
  "accuracy_threshold": 5.0
}
```

### Example 3: Portfolio Analysis

```bash
# Get comprehensive portfolio metrics
GET /users/1/portfolios/5/metrics

# Response:
{
  "total_value": 125000.50,
  "cost_basis": 100000.00,
  "gain_loss": 25000.50,
  "return_percent": 25.0,
  "volatility": 0.18,
  "sharpe_ratio": 1.42,
  "diversification_score": 0.75
}

# Get top holdings by value
GET /users/1/portfolios/5/top-holdings?limit=3

# Get 30-day volatility
GET /users/1/portfolios/5/volatility?days=30
```

### Example 4: Authentication

```bash
# Register/Create User
POST /users
{
  "name": "John Trader",
  "email": "john@example.com",
  "password": "SecurePassword123!"
}

# Login
POST /auth/login
{
  "email": "john@example.com",
  "password": "SecurePassword123!"
}

# Response:
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_id": 1,
  "role": "user"
}

# Use token for authenticated requests
GET /users/1/portfolios \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

## 9. Performance Considerations

### Prediction Generation
- **Time Complexity**: O(n) where n = lookback period (default 100)
- **Lookback Period**: Configurable, default 100 days
- **Cache Opportunities**: Cache price data for frequently analyzed stocks

### Portfolio Metrics
- **Calculation Time**: ~50-200ms per portfolio (varies with holdings count)
- **Optimization**: Calculate daily snapshots once per day
- **Indexing**: Ensure proper database indexes on portfolio_id

### Metrics Query
- **Aggregation**: Perform SUM/AVG at database level
- **Retention**: Consider archiving old metrics (>1 year)
- **Batch Operations**: Update multiple metrics in single transaction

---

## 10. Future Enhancements

1. **Machine Learning Models**: Neural networks, LSTM for sequence prediction
2. **Risk Management**: Value at Risk (VaR), Expected Shortfall
3. **Backtesting**: Historical simulation of strategies
4. **Portfolio Optimization**: Efficient frontier calculation
5. **Sentiment Analysis**: News/social media sentiment incorporation
6. **API Rate Limiting**: Per-user/API-key limits
7. **Prediction Confidence Intervals**: Probabilistic bounds
8. **Multi-Asset Strategies**: Cross-asset correlation analysis
9. **Mobile App**: React Native mobile trading app
10. **WebSocket Integration**: Real-time data streaming

---

## 11. Support & Troubleshooting

### Common Issues

**Issue**: Password too short
**Solution**: Ensure password >= 8 characters

**Issue**: Token expired
**Solution**: Refresh token (implement refresh token endpoint)

**Issue**: High portfolio calculation time
**Solution**: Limit holdings per portfolio or cache recent prices

**Issue**: Inaccurate predictions
**Solution**: Use ENSEMBLE algorithm, adjust lookback period

---

## Database Migration

Run these SQL commands to update existing database:

```sql
-- Add user authentication
ALTER TABLE users ADD COLUMN IF NOT EXISTS password VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(50) DEFAULT 'user';
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT true;
ALTER TABLE users ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Update predictions
ALTER TABLE predictions ADD COLUMN IF NOT EXISTS algorithm VARCHAR(50) DEFAULT 'ENSEMBLE';
ALTER TABLE predictions ADD COLUMN IF NOT EXISTS confidence_score DECIMAL(3,2) DEFAULT 0.5;
ALTER TABLE predictions ADD COLUMN IF NOT EXISTS upper_bound DECIMAL(10,2);
ALTER TABLE predictions ADD COLUMN IF NOT EXISTS lower_bound DECIMAL(10,2);
ALTER TABLE predictions ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Update portfolios
ALTER TABLE portfolios ADD COLUMN IF NOT EXISTS total_value DECIMAL(15,2) DEFAULT 0;
ALTER TABLE portfolios ADD COLUMN IF NOT EXISTS cost_basis DECIMAL(15,2) DEFAULT 0;
ALTER TABLE portfolios ADD COLUMN IF NOT EXISTS gain_loss DECIMAL(15,2) DEFAULT 0;
ALTER TABLE portfolios ADD COLUMN IF NOT EXISTS return_percent DECIMAL(5,2) DEFAULT 0;
ALTER TABLE portfolios ADD COLUMN IF NOT EXISTS last_updated TIMESTAMP;

-- Run init.sql to create new tables (prediction_metrics, portfolio_performance)
```

---

## API Documentation Links

- Authentication: See `middleware/auth_middleware.go`
- Algorithms: See `services/algorithms/prediction_algorithms.go`
- Analytics: See `services/portfolio_analytics_service.go`
- Metrics: See `repositories/prediction_metrics_repository.go`
- Auth Utilities: See `util/auth.go`
