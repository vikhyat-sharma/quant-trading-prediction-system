# Testing Guide

This document provides an overview of the test suite for the Quant Trading Prediction System.

## Test Structure

The project includes comprehensive unit tests for all major components:

### Controllers Tests

- **[stock_controller_test.go](stock_controller_test.go)**: Tests for stock HTTP handlers
  - `TestStockController_GetStock_Success`: Successful stock retrieval
  - `TestStockController_GetStock_InvalidID`: Invalid ID handling
  - `TestStockController_GetAllStocks_Success`: Fetching multiple stocks
  - `TestStockController_GetAllStocks_Error`: Error handling

- **[prediction_controller_test.go](prediction_controller_test.go)**: Tests for prediction HTTP handlers
  - `TestPredictionController_GetPredictions_Success`: Successful prediction retrieval
  - `TestPredictionController_GetPredictions_InvalidStockID`: Invalid stock ID handling
  - `TestPredictionController_GeneratePrediction_Success`: Prediction generation
  - `TestPredictionController_GeneratePrediction_InvalidStockID`: Error handling

### Services Tests

- **[stock_service_test.go](../services/stock_service_test.go)**: Tests for stock business logic
  - `TestStockService_GetStock`: Single stock retrieval with mocked repository
  - `TestStockService_GetAllStocks`: Multiple stocks retrieval with error scenarios

- **[prediction_service_test.go](../services/prediction_service_test.go)**: Tests for prediction business logic
  - `TestPredictionService_GetPredictionsByStockID`: Fetching predictions for a stock
  - `TestPredictionService_GeneratePrediction`: Prediction generation logic

### Repositories Tests

- **[stock_repository_test.go](../repositories/stock_repository_test.go)**: Tests for stock data access
  - Uses `go-sqlmock` for database mocking
  - `TestStockRepository_GetStock_Success`: Successful DB query
  - `TestStockRepository_GetStock_NotFound`: Handling missing records
  - `TestStockRepository_GetAllStocks_Success`: Batch retrieval
  - `TestStockRepository_GetAllStocks_Empty`: Empty result handling

- **[prediction_repository_test.go](../repositories/prediction_repository_test.go)**: Tests for prediction data access
  - Uses `go-sqlmock` for database mocking
  - `TestPredictionRepository_GetPredictionsByStockID_Success`: Successful retrieval
  - `TestPredictionRepository_GetPredictionsByStockID_Empty`: Handling no results
  - `TestPredictionRepository_GetPredictionsByStockID_Error`: DB error handling

### Models Tests

- **[models_test.go](../db/models_test.go)**: Tests for data models
  - `TestStockModel`: Stock struct validation
  - `TestPredictionModel`: Prediction struct validation
  - JSON serialization tests

### Configuration Tests

- **[env_test.go](../config/env_test.go)**: Tests for configuration loading
  - `TestLoadConfig_Defaults`: Default configuration values
  - `TestLoadConfig_WithEnvVars`: Environment variable overrides
  - Configuration struct field validation

## Running Tests

### Run all tests
```bash
go test ./...
```

### Run tests with verbose output
```bash
go test -v ./...
```

### Run specific test package
```bash
go test ./controllers
go test ./services
go test ./repositories
```

### Run specific test
```bash
go test ./controllers -run TestStockController_GetStock_Success
```

### Run tests with coverage
```bash
go test -cover ./...
```

### Generate detailed coverage report
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Using Makefile
```bash
make test
```

## Test Dependencies

The tests use the following packages for mocking and assertions:

- **go-sqlmock**: Mock SQL database connections for repository testing
  - Allows testing database queries without a real database
  - Validates query execution and parameter binding

## Mocking Strategy

### Repository Layer
- Uses `go-sqlmock` to mock database connections
- Validates SQL queries and parameters
- Tests both success and error scenarios

### Service Layer
- Mock repositories are implemented directly in test files
- Services are tested in isolation from the database layer
- Allows testing business logic without database dependencies

### Controller Layer
- Mock services are used to test HTTP handlers
- Uses Go's standard `httptest` package
- Tests request validation and response formatting

## Test Coverage Goals

The test suite aims to cover:

1. **Happy Path**: Normal operation scenarios
2. **Error Handling**: Database errors, invalid inputs, missing data
3. **Edge Cases**: Empty results, boundary conditions
4. **Input Validation**: Invalid IDs, malformed requests
5. **Response Format**: JSON encoding, status codes

## Adding New Tests

When adding new functionality:

1. Add unit tests in the appropriate `*_test.go` file
2. Use table-driven tests for multiple scenarios
3. Mock external dependencies (database, services)
4. Test both success and error cases
5. Ensure coverage remains high with `go test -cover`

## Continuous Integration

These tests can be run in CI/CD pipelines:

```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

## Troubleshooting

### Import errors with go-sqlmock
Ensure the dependency is installed:
```bash
go get -u github.com/DATA-DOG/go-sqlmock
```

Then run:
```bash
go mod tidy
```

### Test failures
- Check that environment variables are properly set or cleared
- Verify mock expectations match actual query calls
- Review test logs for detailed error messages
