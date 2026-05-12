# Search and Filtering Capabilities

This document describes the comprehensive search and filtering features added to the Quant Trading Prediction System.

## Overview

The API now supports flexible search and filtering across multiple endpoints, allowing clients to efficiently retrieve specific data based on various criteria.

## Implemented Features

### 1. Stock Search and Filtering

**Endpoint**: `GET /stocks`

**Available Parameters**:
- `search` - Search by stock symbol or name (case-insensitive, partial matching)
- `exchange` - Filter by exchange name (case-insensitive)

**Examples**:
```bash
# Search for stocks with "apple" in symbol or name
curl "http://localhost:8080/stocks?search=apple"

# Filter stocks by exchange
curl "http://localhost:8080/stocks?exchange=NYSE"

# Combine filters
curl "http://localhost:8080/stocks?search=tech&exchange=NASDAQ"
```

**Implementation**:
- Repository: `repositories/stock_repository.go` - `SearchAndFilterStocks()`
- Service: `services/stock_service.go` - `SearchAndFilterStocks()`
- Controller: `controllers/stock_controller.go` - `GetAllStocks()`

---

### 2. User Search and Filtering

**Endpoint**: `GET /users`

**Available Parameters**:
- `search` - Search by name or email (case-insensitive, partial matching)

**Examples**:
```bash
# Search for users by name
curl "http://localhost:8080/users?search=john"

# Search by email
curl "http://localhost:8080/users?search=john@example.com"
```

**Implementation**:
- Repository: `repositories/user_repository.go` - `SearchAndFilterUsers()`
- Service: `services/user_service.go` - `SearchAndFilterUsers()`
- Controller: `controllers/user_controller.go` - `GetUsers()`

---

### 3. Portfolio Search and Filtering

**Endpoint**: `GET /users/{userID}/portfolios`

**Available Parameters**:
- `search` - Search by portfolio name (case-insensitive, partial matching)

**Examples**:
```bash
# Search for portfolios with specific name
curl "http://localhost:8080/users/1/portfolios?search=conservative"

# Filter user's portfolios by name
curl "http://localhost:8080/users/1/portfolios?search=tech"
```

**Implementation**:
- Repository: `repositories/portfolio_repository.go` - `SearchAndFilterPortfolios()`
- Service: `services/portfolio_service.go` - `SearchAndFilterPortfolios()`
- Controller: `controllers/portfolio_controller.go` - `GetPortfolios()`

---

### 4. Price History Search and Filtering

**Endpoint**: `GET /stocks/{stockID}/price-history`

**Available Parameters**:
- `start_date` - Filter prices from this date (format: YYYY-MM-DD)
- `end_date` - Filter prices until this date (format: YYYY-MM-DD)
- `min_price` - Filter prices with minimum value
- `max_price` - Filter prices with maximum value

**Examples**:
```bash
# Get price history for a date range
curl "http://localhost:8080/stocks/1/price-history?start_date=2024-01-01&end_date=2024-12-31"

# Get prices within a price range
curl "http://localhost:8080/stocks/1/price-history?min_price=50&max_price=150"

# Get prices from specific period and price range
curl "http://localhost:8080/stocks/1/price-history?start_date=2024-06-01&min_price=100&max_price=200"
```

**Implementation**:
- Repository: `repositories/price_history_repository.go` - `SearchAndFilterPriceHistory()`
- Service: `services/price_history_service.go` - `SearchAndFilterPriceHistory()`
- Controller: `controllers/price_history_controller.go` - `GetPriceHistory()`

---

### 5. Prediction Search and Filtering

**Endpoint**: `GET /stocks/{stockID}/predictions`

**Available Parameters**:
- `start_date` - Filter predictions from this date (format: YYYY-MM-DD)
- `end_date` - Filter predictions until this date (format: YYYY-MM-DD)
- `min_price` - Filter predictions with minimum predicted price
- `max_price` - Filter predictions with maximum predicted price

**Examples**:
```bash
# Get predictions between two dates
curl "http://localhost:8080/stocks/1/predictions?start_date=2024-01-01&end_date=2024-12-31"

# Get predictions within a price range
curl "http://localhost:8080/stocks/1/predictions?min_price=100&max_price=200"

# Get predictions from specific period and price range
curl "http://localhost:8080/stocks/1/predictions?start_date=2024-06-01&min_price=150&max_price=300"
```

**Implementation**:
- Repository: `repositories/prediction_repository.go` - `SearchAndFilterPredictions()`
- Service: `services/prediction_service.go` - `SearchAndFilterPredictions()`
- Controller: `controllers/prediction_controller.go` - `GetPredictions()`

---

## Filter Types

### Text Search Filters
Used for `search` parameters:
- **Case-Insensitive**: Matches both uppercase and lowercase variations
- **Partial Matching**: Finds results containing the search term anywhere in the field
- **Multi-field**: Searches across relevant fields (e.g., symbol and name for stocks)

### Date Range Filters
Used for `start_date` and `end_date` parameters:
- **Format**: YYYY-MM-DD (ISO 8601)
- **Inclusive**: Both start and end dates are included in the range
- **Optional**: Either can be used independently
- **Validation**: Returns 400 Bad Request for invalid date formats

### Price Range Filters
Used for `min_price` and `max_price` parameters:
- **Format**: Floating-point numbers
- **Inclusive**: Both boundaries are included in the range
- **Optional**: Either can be used independently
- **Validation**: Returns 400 Bad Request for invalid number formats

---

## Database Query Patterns

All filters use parameterized queries to prevent SQL injection:

### Example SQL Generation
```sql
-- Stock search with exchange filter
SELECT id, symbol, exchange, name FROM stocks 
WHERE 1=1 
  AND (UPPER(symbol) LIKE $1 OR UPPER(name) LIKE $1)
  AND UPPER(exchange) = $2
ORDER BY symbol

-- Price history with date and price range
SELECT id, stock_id, price, date, created_at FROM price_history 
WHERE 1=1 
  AND stock_id = $1
  AND date BETWEEN $2 AND $3
  AND price BETWEEN $4 AND $5
ORDER BY date DESC
```

---

## Error Handling

All filter endpoints include comprehensive error handling:

- **400 Bad Request**: Invalid parameter formats (e.g., malformed dates)
- **200 OK**: Empty result set is valid (returns empty array)
- **500 Internal Server Error**: Database or server errors

**Example Error Response**:
```json
{
  "error": "Invalid start_date format (use YYYY-MM-DD)",
  "details": "parsing time \"2024/01/01\": ..."
}
```

---

## Performance Considerations

### Indexed Fields
The following fields should be indexed in the database for optimal performance:
- `stocks.symbol` (text search)
- `stocks.name` (text search)
- `stocks.exchange` (filtering)
- `users.name` (text search)
- `users.email` (text search)
- `portfolios.name` (text search)
- `portfolios.user_id` (filtering)
- `price_history.stock_id` (filtering)
- `price_history.date` (range filtering)
- `price_history.price` (range filtering)
- `predictions.stock_id` (filtering)
- `predictions.date` (range filtering)
- `predictions.predicted_price` (range filtering)

---

## Usage Examples

### Search for Tech Stocks
```bash
GET /stocks?search=tech
```

### Find User by Email
```bash
GET /users?search=user@example.com
```

### Get Recent Price Data
```bash
GET /stocks/1/price-history?start_date=2024-12-01&end_date=2024-12-31
```

### Find Stable Predictions
```bash
GET /stocks/1/predictions?min_price=100&max_price=110
```

### Complex Filter: Recent Predictions Within Price Range
```bash
GET /stocks/1/predictions?start_date=2024-06-01&end_date=2024-12-31&min_price=100&max_price=200
```

---

## Testing

All search and filtering functionality can be tested with curl:

```bash
# Test stock search
curl -X GET "http://localhost:8080/stocks?search=apple&exchange=NYSE"

# Test user search
curl -X GET "http://localhost:8080/users?search=john"

# Test price history filtering
curl -X GET "http://localhost:8080/stocks/1/price-history?start_date=2024-01-01&end_date=2024-12-31&min_price=50&max_price=150"

# Test prediction filtering
curl -X GET "http://localhost:8080/stocks/1/predictions?min_price=100&max_price=200"
```

---

## Future Enhancements

Potential improvements to search and filtering:
1. **Pagination**: Add limit/offset for large result sets
2. **Sorting**: Add sort parameters for results
3. **Full-Text Search**: Implement PostgreSQL full-text search for better text matching
4. **Advanced Operators**: Support AND/OR/NOT logic in searches
5. **Fuzzy Matching**: Implement approximate string matching for typo tolerance
6. **Caching**: Cache frequently used search results
7. **Search Analytics**: Track popular searches for insights

---

## Implementation Details

### Type Definitions

All filter types are defined in their respective repositories:

- **StockFilter**: `repositories/stock_repository.go`
- **UserFilter**: `repositories/user_repository.go`
- **PortfolioFilter**: `repositories/portfolio_repository.go`
- **PriceHistoryFilter**: `repositories/price_history_repository.go`
- **PredictionFilter**: `repositories/prediction_repository.go`

### Architecture

The implementation follows the clean architecture pattern:
1. **Controllers**: Parse query parameters and validate input
2. **Services**: Call repository methods and handle business logic
3. **Repositories**: Execute parameterized SQL queries with filters

This separation ensures maintainability, testability, and adherence to SOLID principles.
