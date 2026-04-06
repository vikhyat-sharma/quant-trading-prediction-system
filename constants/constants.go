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

	// Not found errors
	ErrMsgStockNotFound         = "Stock not found"
	ErrMsgPredictionNotFound    = "Prediction not found"
	ErrMsgNoPredictionsForStock = "No predictions found for this stock"

	// Operation errors
	ErrMsgFailedToRetrieveStock       = "Failed to retrieve stock"
	ErrMsgFailedToRetrieveStocks      = "Failed to retrieve stocks"
	ErrMsgFailedToRetrievePredictions = "Failed to retrieve predictions"
	ErrMsgFailedToGeneratePrediction  = "Failed to generate prediction"

	// Configuration errors
	ErrMsgPortCannotBeEmpty        = "PORT cannot be empty"
	ErrMsgPortMustBeValidNumber    = "PORT must be a valid number"
	ErrMsgDatabaseURLCannotBeEmpty = "DATABASE_URL cannot be empty"
	ErrMsgEnvironmentInvalid       = "ENVIRONMENT must be one of: development, staging, production"
	ErrMsgLogLevelInvalid          = "LOG_LEVEL must be one of: debug, info, warn, error"
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

// Route Paths
const (
	RouteStocks                   = "/stocks"
	RouteStockByID                = "/stocks/{id}"
	RouteStockPredictions         = "/stocks/{stockID}/predictions"
	RouteStockPredictionsGenerate = "/stocks/{stockID}/predictions/generate"
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
)
