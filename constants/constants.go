package constants

import "time"

// HTTP Headers
const (
	HeaderContentType               = "Content-Type"
	HeaderAccessControlAllowOrigin  = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders = "Access-Control-Allow-Headers"
	HeaderAuthorization             = "Authorization"
)

// HTTP Methods
const (
	MethodGET     = "GET"
	MethodPOST    = "POST"
	MethodPUT     = "PUT"
	MethodDELETE  = "DELETE"
	MethodOPTIONS = "OPTIONS"
)

// Content Types
const (
	ContentTypeJSON = "application/json"
	ContentTypeText = "text/plain"
)

// CORS Settings
const (
	CORSAllowOrigin  = "*"
	CORSAllowMethods = "GET, POST, PUT, DELETE, OPTIONS"
	CORSAllowHeaders = "Content-Type, Authorization"
)

// Error Messages
const (
	// Validation errors
	ErrMsgInvalidStockIDFormat       = "Invalid stock ID format"
	ErrMsgStockIDMustBePositive      = "Stock ID must be a positive integer"
	ErrMsgInvalidPredictionIDFormat  = "Invalid prediction ID format"
	ErrMsgPredictionIDMustBePositive = "Prediction ID must be a positive integer"
	ErrMsgInvalidUserIDFormat        = "Invalid user ID format"
	ErrMsgUserIDMustBePositive       = "User ID must be a positive integer"
	ErrMsgInvalidPortfolioIDFormat   = "Invalid portfolio ID format"
	ErrMsgPortfolioIDMustBePositive  = "Portfolio ID must be a positive integer"
	ErrMsgInvalidHoldingIDFormat     = "Invalid holding ID format"
	ErrMsgHoldingIDMustBePositive    = "Holding ID must be a positive integer"

	// Not found errors
	ErrMsgStockNotFound         = "Stock not found"
	ErrMsgPredictionNotFound    = "Prediction not found"
	ErrMsgNoPredictionsForStock = "No predictions found for this stock"
	ErrMsgUserNotFound          = "User not found"
	ErrMsgPortfolioNotFound     = "Portfolio not found"
	ErrMsgHoldingNotFound       = "Holding not found"

	// Operation errors
	ErrMsgFailedToRetrieveStock       = "Failed to retrieve stock"
	ErrMsgFailedToRetrieveStocks      = "Failed to retrieve stocks"
	ErrMsgFailedToRetrievePredictions = "Failed to retrieve predictions"
	ErrMsgFailedToGeneratePrediction  = "Failed to generate prediction"
	ErrMsgFailedToCreateStock         = "Failed to create stock"
	ErrMsgFailedToUpdateStock         = "Failed to update stock"
	ErrMsgFailedToDeleteStock         = "Failed to delete stock"
	ErrMsgFailedToRetrieveUsers       = "Failed to retrieve users"
	ErrMsgFailedToRetrieveUser        = "Failed to retrieve user"
	ErrMsgFailedToCreateUser          = "Failed to create user"
	ErrMsgFailedToUpdateUser          = "Failed to update user"
	ErrMsgFailedToDeleteUser          = "Failed to delete user"
	ErrMsgFailedToRetrievePortfolios  = "Failed to retrieve portfolios"
	ErrMsgFailedToRetrievePortfolio   = "Failed to retrieve portfolio"
	ErrMsgFailedToCreatePortfolio     = "Failed to create portfolio"
	ErrMsgFailedToUpdatePortfolio     = "Failed to update portfolio"
	ErrMsgFailedToDeletePortfolio     = "Failed to delete portfolio"
	ErrMsgFailedToRetrieveHoldings    = "Failed to retrieve holdings"
	ErrMsgFailedToCreateHolding       = "Failed to create holding"
	ErrMsgFailedToUpdateHolding       = "Failed to update holding"
	ErrMsgFailedToDeleteHolding       = "Failed to delete holding"

	// Configuration errors
	ErrMsgPortCannotBeEmpty        = "PORT cannot be empty"
	ErrMsgPortMustBeValidNumber    = "PORT must be a valid number"
	ErrMsgDatabaseURLCannotBeEmpty = "DATABASE_URL cannot be empty"
	ErrMsgEnvironmentInvalid       = "ENVIRONMENT must be one of: development, staging, production"
	ErrMsgLogLevelInvalid          = "LOG_LEVEL must be one of: debug, info, warn, error"
)

// Stock market errors
const (
	ErrMsgInvalidExchange = "Exchange must be NSE or BSE"
)

// Configuration Keys and Defaults
const (
	// Environment variable keys
	EnvKeyPort        = "PORT"
	EnvKeyDatabaseURL = "DATABASE_URL"
	EnvKeyEnvironment = "ENVIRONMENT"
	EnvKeyLogLevel    = "LOG_LEVEL"

	// Default values
	DefaultPort        = "8080"
	DefaultDatabaseURL = "postgres://user:password@localhost/quant_trading?sslmode=disable"
	DefaultEnvironment = "development"
	DefaultLogLevel    = "info"
)

// Database Settings
const (
	DatabaseDriverPostgres = "postgres"

	// Connection pool defaults
	DefaultMaxOpenConns    = 25
	DefaultMaxIdleConns    = 5
	DefaultConnMaxLifetime = 5 * time.Minute
)

// Stock exchange constants
const (
	ExchangeNSE     = "NSE"
	ExchangeBSE     = "BSE"
	DefaultExchange = ExchangeNSE
)

// Route Paths
const (
	RouteStocks                   = "/stocks"
	RouteStockByID                = "/stocks/{id}"
	RouteStockPredictions         = "/stocks/{stockID}/predictions"
	RouteStockPredictionsGenerate = "/stocks/{stockID}/predictions/generate"
	RouteStockPriceHistory        = "/stocks/{stockID}/price-history"
	RouteStockPriceHistoryRange   = "/stocks/{stockID}/price-history/range"
	RouteStockPriceStats          = "/stocks/{stockID}/price-stats"
	RouteStockLatestPrice         = "/stocks/{stockID}/price-latest"
	RouteStockAlerts              = "/stocks/{stockID}/alerts"
	RouteStockAlertByID           = "/stocks/{stockID}/alerts/{alertID}"
	RouteStockAlertsEvaluate      = "/stocks/{stockID}/alerts/evaluate"
	RouteStockNotifications       = "/stocks/{stockID}/notifications"
	RouteUsers                    = "/users"
	RouteUserByID                 = "/users/{id}"
	RouteUserPortfolios           = "/users/{userID}/portfolios"
	RouteUserPortfolioByID        = "/users/{userID}/portfolios/{portfolioID}"
	RouteUserPortfolioHoldings    = "/users/{userID}/portfolios/{portfolioID}/holdings"
	RouteUserPortfolioHoldingByID = "/users/{userID}/portfolios/{portfolioID}/holdings/{holdingID}"
)

// Server Settings
const (
	DefaultServerShutdownTimeout = 30 * time.Second
)

// Log Messages
const (
	LogMsgServerStarting       = "Server starting on port %s"
	LogMsgServerShuttingDown   = "Server is shutting down..."
	LogMsgServerExited         = "Server exited"
	LogMsgFailedToLoadConfig   = "Failed to load config"
	LogMsgFailedToConnectDB    = "Failed to connect to database"
	LogMsgFailedToStartServer  = "Failed to start server"
	LogMsgServerForcedShutdown = "Server forced to shutdown"
	LogMsgFailedToEnsureSchema = "Failed to ensure database schema"
)
