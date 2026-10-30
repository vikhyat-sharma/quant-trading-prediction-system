# Implementation & Migration Guide

## Quick Start

This guide walks you through implementing all the advanced features in your existing system.

---

## Step 1: Update Dependencies

Add new required packages to `go.mod`:

```bash
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto
```

Or update `go.mod` directly:
```go
require (
    github.com/golang-jwt/jwt/v5 v5.2.0
    golang.org/x/crypto v0.24.0
)
```

Then run:
```bash
go mod tidy
```

---

## Step 2: Database Migration

### Option A: Fresh Database
If starting fresh, simply run `init.sql`:
```bash
psql -U postgres -d quant_trading < init.sql
```

### Option B: Existing Database
Run these migration commands to update existing schema:

```sql
-- 1. Update users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS password VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(50) DEFAULT 'user';
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT true;
ALTER TABLE users ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- 2. Update predictions table
ALTER TABLE predictions ADD COLUMN IF NOT EXISTS algorithm VARCHAR(50) DEFAULT 'ENSEMBLE';
ALTER TABLE predictions ADD COLUMN IF NOT EXISTS confidence_score DECIMAL(3,2) DEFAULT 0.5;
ALTER TABLE predictions ADD COLUMN IF NOT EXISTS upper_bound DECIMAL(10,2);
ALTER TABLE predictions ADD COLUMN IF NOT EXISTS lower_bound DECIMAL(10,2);
ALTER TABLE predictions ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- 3. Update portfolios table
ALTER TABLE portfolios ADD COLUMN IF NOT EXISTS total_value DECIMAL(15,2) DEFAULT 0;
ALTER TABLE portfolios ADD COLUMN IF NOT EXISTS cost_basis DECIMAL(15,2) DEFAULT 0;
ALTER TABLE portfolios ADD COLUMN IF NOT EXISTS gain_loss DECIMAL(15,2) DEFAULT 0;
ALTER TABLE portfolios ADD COLUMN IF NOT EXISTS return_percent DECIMAL(5,2) DEFAULT 0;
ALTER TABLE portfolios ADD COLUMN IF NOT EXISTS last_updated TIMESTAMP;

-- 4. Create new tables
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

-- 5. Create indexes
CREATE INDEX IF NOT EXISTS idx_predictions_algorithm ON predictions(algorithm);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_prediction_metrics_stock_id ON prediction_metrics(stock_id);
CREATE INDEX IF NOT EXISTS idx_prediction_metrics_algorithm ON prediction_metrics(algorithm);
CREATE INDEX IF NOT EXISTS idx_prediction_metrics_prediction_id ON prediction_metrics(prediction_id);
CREATE INDEX IF NOT EXISTS idx_portfolio_performance_portfolio_id ON portfolio_performance(portfolio_id);
CREATE INDEX IF NOT EXISTS idx_portfolio_performance_record_date ON portfolio_performance(record_date);
```

---

## Step 3: Update Services

### 3.1 Update PredictionService

Replace the old simple implementation with the new one:

```go
// Old service - remove
type PredictionService struct {
    repo *repositories.PredictionRepository
}

// New service - use this
type PredictionService struct {
    repo                 *repositories.PredictionRepository
    priceHistoryRepo     *repositories.PriceHistoryRepository
    metricsRepo          *repositories.PredictionMetricsRepository
    defaultAlgorithm     string
    lookbackPeriod       int
}

func NewPredictionServiceWithMetrics(
    repo *repositories.PredictionRepository,
    priceHistoryRepo *repositories.PriceHistoryRepository,
    metricsRepo *repositories.PredictionMetricsRepository,
) *PredictionService {
    return &PredictionService{
        repo:             repo,
        priceHistoryRepo: priceHistoryRepo,
        metricsRepo:      metricsRepo,
        defaultAlgorithm: "ENSEMBLE",
        lookbackPeriod:   100,
    }
}
```

### 3.2 Update Portfolio Service

Add analytics capabilities:

```go
// Add to main.go initialization
portfolioAnalyticsService := services.NewPortfolioAnalyticsService(portfolioRepo, priceHistoryRepo)
```

### 3.3 Create Metrics Service

Add metrics tracking:

```go
// New service for tracking metrics
type MetricsService struct {
    metricsRepo *repositories.PredictionMetricsRepository
}

func NewMetricsService(repo *repositories.PredictionMetricsRepository) *MetricsService {
    return &MetricsService{metricsRepo: repo}
}
```

---

## Step 4: Update Models

### 4.1 Update User Model

Add password and role fields:

```go
type User struct {
    ID        int       `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    Email     string    `json:"email" db:"email"`
    Password  string    `json:"-" db:"password"`    // NEW: Hashed password
    Role      string    `json:"role" db:"role"`      // NEW: user/admin/superadmin
    IsActive  bool      `json:"is_active" db:"is_active"` // NEW
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"` // NEW
}
```

### 4.2 Update Prediction Model

Add algorithm and confidence fields:

```go
type Prediction struct {
    ID               int       `json:"id" db:"id"`
    StockID          int       `json:"stock_id" db:"stock_id"`
    PredictedPrice   float64   `json:"predicted_price" db:"predicted_price"`
    Algorithm        string    `json:"algorithm" db:"algorithm"`                    // NEW
    ConfidenceScore  float64   `json:"confidence_score" db:"confidence_score"`      // NEW
    UpperBound       float64   `json:"upper_bound" db:"upper_bound"`                // NEW
    LowerBound       float64   `json:"lower_bound" db:"lower_bound"`                // NEW
    Date             time.Time `json:"date" db:"date"`
    CreatedAt        time.Time `json:"created_at" db:"created_at"`                  // NEW
}
```

### 4.3 Create New Models

Create `db/metrics_models.go` with:

```go
type PredictionMetric struct { ... }
type AlgorithmPerformance struct { ... }
type PortfolioPerformance struct { ... }
```

---

## Step 5: Add Repositories

### 5.1 Add Metrics Repository

Create `repositories/prediction_metrics_repository.go` with methods:
- `CreateMetric()`
- `GetMetricByPredictionID()`
- `GetPerformanceByAlgorithm()`
- `UpdateMetricWithActualPrice()`
- `GetAllAlgorithmPerformance()`

### 5.2 Update Price History Repository

Add method:
- `GetHistoricalPrices(stockID, limit)` - For algorithm calculations

---

## Step 6: Add Utilities

### 6.1 Create Authentication Utils

Create `util/auth.go` with:
```go
func HashPassword(password string) (string, error)
func VerifyPassword(hashedPassword, password string) bool
func GenerateJWT(userID int, email, role string) (string, error)
func VerifyJWT(tokenString string) (*Claims, error)
func ValidateEmail(email string) bool
func GenerateAPIKey() (string, error)
func IsAdminRole(role string) bool
func IsUserRole(role string) bool
```

---

## Step 7: Create Algorithm Package

### 7.1 Create Prediction Algorithms

Create `services/algorithms/prediction_algorithms.go` with:

**Algorithms**:
- `SimpleMovingAveragePrediction()`
- `ExponentialMovingAveragePrediction()`
- `MomentumPrediction()`
- `MeanReversionPrediction()`
- `EnsemblePrediction()`

**Indicators**:
- `CalculateSMA()`
- `CalculateEMA()`
- `CalculateRSI()`
- `CalculateMACD()`
- `CalculateBollingerBands()`
- `CalculateTechnicalIndicators()`

---

## Step 8: Update Middleware

### 8.1 Add Authentication Middleware

Create `middleware/auth_middleware.go` with:
- `AuthMiddleware()` - Requires JWT
- `OptionalAuthMiddleware()` - Optional JWT
- `AdminMiddleware()` - Requires admin role
- `RateLimitMiddleware()` - Rate limiting
- `RequestIDMiddleware()` - Request tracking
- `RecoveryMiddleware()` - Panic recovery

---

## Step 9: Update Main Application

### 9.1 Initialize New Services in main.go

```go
// Add metrics repository
metricsRepo := repositories.NewPredictionMetricsRepository(database)

// Update prediction service
predictionService := services.NewPredictionServiceWithMetrics(
    predictionRepo,
    priceHistoryRepo,
    metricsRepo,
)

// Add analytics service
portfolioAnalyticsService := services.NewPortfolioAnalyticsService(
    portfolioRepo,
    priceHistoryRepo,
)
```

### 9.2 Add Middleware to Router

```go
router.Use(middleware.RecoveryMiddleware)
router.Use(middleware.RequestIDMiddleware)
router.Use(middleware.LoggingMiddleware)
router.Use(middleware.CORSMiddleware)
router.Use(middleware.ContentTypeMiddleware)

// Optional: Add rate limiting
rateLimiter := middleware.NewRateLimiter(100, time.Minute)
router.Use(rateLimiter.Middleware)
```

---

## Step 10: Create Controllers

### 10.1 Create Auth Controller

```go
// POST /auth/register
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request)

// POST /auth/login
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request)

// POST /auth/verify-token
func (c *AuthController) VerifyToken(w http.ResponseWriter, r *http.Request)
```

### 10.2 Create Metrics Controller

```go
// GET /metrics/predictions/{predictionID}
func (c *MetricsController) GetPredictionMetric(w http.ResponseWriter, r *http.Request)

// GET /metrics/algorithms/{algorithmName}
func (c *MetricsController) GetAlgorithmPerformance(w http.ResponseWriter, r *http.Request)

// GET /metrics/algorithms
func (c *MetricsController) GetAllAlgorithmPerformance(w http.ResponseWriter, r *http.Request)

// PUT /metrics/predictions/{predictionID}
func (c *MetricsController) UpdateMetricWithActualPrice(w http.ResponseWriter, r *http.Request)
```

### 10.3 Create Analytics Controller

```go
// GET /users/{userID}/portfolios/{portfolioID}/metrics
func (c *AnalyticsController) GetPortfolioMetrics(w http.ResponseWriter, r *http.Request)

// GET /users/{userID}/portfolios/{portfolioID}/value
func (c *AnalyticsController) GetPortfolioValue(w http.ResponseWriter, r *http.Request)

// GET /users/{userID}/portfolios/{portfolioID}/volatility
func (c *AnalyticsController) GetPortfolioVolatility(w http.ResponseWriter, r *http.Request)

// GET /users/{userID}/portfolios/{portfolioID}/sharpe-ratio
func (c *AnalyticsController) GetSharpeRatio(w http.ResponseWriter, r *http.Request)

// GET /users/{userID}/portfolios/{portfolioID}/top-holdings
func (c *AnalyticsController) GetTopHoldings(w http.ResponseWriter, r *http.Request)

// GET /users/{userID}/portfolios/{portfolioID}/diversification
func (c *AnalyticsController) GetDiversificationScore(w http.ResponseWriter, r *http.Request)
```

---

## Step 11: Update Routes

### 11.1 Add New Routes

```go
// Authentication routes (no auth required)
router.HandleFunc("/auth/register", authController.Register).Methods("POST")
router.HandleFunc("/auth/login", authController.Login).Methods("POST")
router.HandleFunc("/auth/verify-token", authController.VerifyToken).Methods("POST")

// Prediction routes with algorithm support
router.HandleFunc("/stocks/{stockID}/predictions", 
    authMiddleware(predictionController.GetPredictions)).Methods("GET")
router.HandleFunc("/stocks/{stockID}/predictions", 
    authMiddleware(predictionController.GeneratePrediction)).Methods("POST")

// Metrics routes (requires auth)
router.HandleFunc("/metrics/predictions/{predictionID}", 
    authMiddleware(metricsController.GetPredictionMetric)).Methods("GET")
router.HandleFunc("/metrics/algorithms/{algorithmName}", 
    authMiddleware(metricsController.GetAlgorithmPerformance)).Methods("GET")
router.HandleFunc("/metrics/algorithms", 
    authMiddleware(metricsController.GetAllAlgorithmPerformance)).Methods("GET")

// Analytics routes (requires auth)
router.HandleFunc("/users/{userID}/portfolios/{portfolioID}/metrics", 
    authMiddleware(analyticsController.GetPortfolioMetrics)).Methods("GET")
router.HandleFunc("/users/{userID}/portfolios/{portfolioID}/value", 
    authMiddleware(analyticsController.GetPortfolioValue)).Methods("GET")
router.HandleFunc("/users/{userID}/portfolios/{portfolioID}/volatility", 
    authMiddleware(analyticsController.GetPortfolioVolatility)).Methods("GET")
```

---

## Step 12: Testing

### 12.1 Test Authentication

```bash
# Register
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Trader",
    "email": "john@example.com",
    "password": "SecurePassword123!"
  }'

# Login
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePassword123!"
  }'
```

### 12.2 Test Predictions

```bash
# Generate prediction with ENSEMBLE algorithm
curl -X POST "http://localhost:8080/stocks/1/predictions?algorithm=ENSEMBLE" \
  -H "Authorization: Bearer <token>"

# Get predictions with filters
curl -X GET "http://localhost:8080/stocks/1/predictions?start_date=2024-01-01&min_price=100" \
  -H "Authorization: Bearer <token>"
```

### 12.3 Test Metrics

```bash
# Get algorithm performance
curl -X GET "http://localhost:8080/metrics/algorithms/ENSEMBLE" \
  -H "Authorization: Bearer <token>"

# Get all algorithm performance
curl -X GET "http://localhost:8080/metrics/algorithms" \
  -H "Authorization: Bearer <token>"
```

### 12.4 Test Portfolio Analytics

```bash
# Get portfolio metrics
curl -X GET "http://localhost:8080/users/1/portfolios/1/metrics" \
  -H "Authorization: Bearer <token>"

# Get portfolio volatility
curl -X GET "http://localhost:8080/users/1/portfolios/1/volatility?days=30" \
  -H "Authorization: Bearer <token>"

# Get top holdings
curl -X GET "http://localhost:8080/users/1/portfolios/1/top-holdings?limit=5" \
  -H "Authorization: Bearer <token>"
```

---

## Step 13: Verify Implementation

### 13.1 Compile the Project

```bash
go build -o app .
```

### 13.2 Run the Application

```bash
./app
# or
go run main.go
```

### 13.3 Check Logs

Verify middleware and routes are working:
```
[INFO] Server starting on port 8080
[DEBUG] Initialized 5 repositories
[DEBUG] Initialized 3 services
[DEBUG] Registered 15 routes
```

---

## Troubleshooting

### Issue: "table prediction_metrics does not exist"
**Solution**: Run database migration script

### Issue: "jwt: signature verification failed"
**Solution**: Ensure JWT_SECRET_KEY is correct in auth.go

### Issue: "password must be at least 8 characters"
**Solution**: Ensure password >= 8 characters in registration

### Issue: "algorithm not found"
**Solution**: Use one of: SMA, EMA, MOMENTUM, MEAN_REVERSION, ENSEMBLE

### Issue: "rate limit exceeded"
**Solution**: Reduce request frequency or adjust RateLimiter configuration

---

## Performance Tuning

### Optimize Prediction Generation
```go
// Increase lookback period for more accurate predictions
service.lookbackPeriod = 250  // 1 year of data

// Decrease for faster calculation
service.lookbackPeriod = 30   // 1 month of data
```

### Database Optimization
```sql
-- Analyze query performance
EXPLAIN ANALYZE
SELECT * FROM predictions 
WHERE stock_id = 1 
  AND algorithm = 'ENSEMBLE' 
  AND created_at > NOW() - INTERVAL '30 days';

-- Vacuum and analyze
VACUUM ANALYZE predictions;
VACUUM ANALYZE prediction_metrics;
```

### Cache Frequent Queries
```go
// Consider caching frequently accessed data:
// - Latest prices per stock
// - Algorithm performance stats
// - Portfolio metrics
```

---

## Next Steps

1. ✅ Update dependencies
2. ✅ Migrate database
3. ✅ Update services and models
4. ✅ Add repositories
5. ✅ Create utilities and algorithms
6. ✅ Add middleware and controllers
7. ✅ Update routes
8. ✅ Test all features
9. ⬜ Deploy to production
10. ⬜ Monitor performance
11. ⬜ Iterate based on user feedback

---

## Additional Resources

- JWT Documentation: https://github.com/golang-jwt/jwt
- Bcrypt Documentation: https://golang.org/x/crypto
- PostgreSQL JSON: https://www.postgresql.org/docs/current/datatype-json.html
- Technical Analysis: https://en.wikipedia.org/wiki/Technical_analysis
