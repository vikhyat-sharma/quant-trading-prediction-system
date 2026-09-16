package routes

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/constants"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/controllers"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/middleware"
)

// SetupRoutes wires all routes and middleware.
func SetupRoutes(
	stockController *controllers.StockController,
	predictionController *controllers.PredictionController,
	priceHistoryController *controllers.PriceHistoryController,
	alertController *controllers.AlertController,
	userController *controllers.UserController,
	portfolioController *controllers.PortfolioController,
	sentimentController *controllers.SentimentController,
	watchlistController *controllers.WatchlistController,
	userAlertRuleController *controllers.UserAlertRuleController,
	taxLotController *controllers.TaxLotController,
	rateLimiter *middleware.RateLimiter,
	database *sql.DB,
) *mux.Router {
	r := mux.NewRouter()

	// Global middleware (applied in order)
	r.Use(middleware.RecoveryMiddleware)
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.CORSMiddleware)
	r.Use(middleware.ContentTypeMiddleware)
	r.Use(middleware.BodyLimitMiddleware)
	r.Use(rateLimiter.Middleware)

	// Health and readiness probes (no auth, no rate limit bypass needed)
	r.HandleFunc("/health", middleware.LivenessHandler).Methods(constants.MethodGET)
	r.HandleFunc("/ready", middleware.ReadinessHandlerFunc(database)).Methods(constants.MethodGET)

	// ── Public routes ──────────────────────────────────────────────────────────

	// Auth
	r.HandleFunc("/auth/login", userController.Login).Methods(constants.MethodPOST)

	// Stock routes (read-only public; mutations require auth)
	r.HandleFunc(constants.RouteStocks, stockController.GetAllStocks).Methods(constants.MethodGET)
	r.HandleFunc(constants.RouteStockByID, stockController.GetStock).Methods(constants.MethodGET)

	// Prediction routes (read-only public)
	r.HandleFunc(constants.RouteStockPredictions, predictionController.GetPredictions).Methods(constants.MethodGET)
	r.HandleFunc(constants.RouteStockBacktest, predictionController.BacktestStrategy).Methods(constants.MethodGET)

	// Price History routes (read-only public)
	r.HandleFunc(constants.RouteStockPriceHistory, priceHistoryController.GetPriceHistory).Methods(constants.MethodGET)
	r.HandleFunc(constants.RouteStockPriceHistoryRange, priceHistoryController.GetPriceHistoryByDateRange).Methods(constants.MethodGET)
	r.HandleFunc(constants.RouteStockPriceStats, priceHistoryController.GetPriceStats).Methods(constants.MethodGET)
	r.HandleFunc(constants.RouteStockLatestPrice, priceHistoryController.GetLatestPrice).Methods(constants.MethodGET)

	// User registration (public)
	r.HandleFunc(constants.RouteUsers, userController.CreateUser).Methods(constants.MethodPOST)

	// ── Authenticated routes ───────────────────────────────────────────────────

	auth := r.NewRoute().Subrouter()
	auth.Use(middleware.AuthMiddleware)

	// Stock mutations (authenticated)
	auth.HandleFunc(constants.RouteStocks, stockController.CreateStock).Methods(constants.MethodPOST)
	auth.HandleFunc(constants.RouteStockByID, stockController.UpdateStock).Methods(constants.MethodPUT)
	auth.HandleFunc(constants.RouteStockByID, stockController.DeleteStock).Methods(constants.MethodDELETE)

	// Prediction generation (authenticated)
	auth.HandleFunc(constants.RouteStockPredictionsGenerate, predictionController.GeneratePrediction).Methods(constants.MethodPOST)
	auth.HandleFunc(constants.RouteStockSentiment, sentimentController.AnalyzeSentiment).Methods(constants.MethodPOST)

	// Price history recording (authenticated)
	auth.HandleFunc(constants.RouteStockPriceHistory, priceHistoryController.RecordPrice).Methods(constants.MethodPOST)

	// Alert routes (authenticated)
	auth.HandleFunc(constants.RouteStockAlerts, alertController.GetAlerts).Methods(constants.MethodGET)
	auth.HandleFunc(constants.RouteStockAlerts, alertController.CreateAlert).Methods(constants.MethodPOST)
	auth.HandleFunc(constants.RouteStockAlertByID, alertController.DeleteAlert).Methods(constants.MethodDELETE)
	auth.HandleFunc(constants.RouteStockAlertsEvaluate, alertController.EvaluateAlerts).Methods(constants.MethodPOST)
	auth.HandleFunc(constants.RouteStockNotifications, alertController.GetNotifications).Methods(constants.MethodGET)

	// User routes (authenticated; ownership enforced in handlers)
	auth.HandleFunc(constants.RouteUsers, userController.GetUsers).Methods(constants.MethodGET)
	auth.HandleFunc(constants.RouteUserByID, userController.GetUser).Methods(constants.MethodGET)
	auth.HandleFunc(constants.RouteUserByID, userController.UpdateUser).Methods(constants.MethodPUT)
	auth.HandleFunc(constants.RouteUserByID, userController.DeleteUser).Methods(constants.MethodDELETE)

	// Portfolio routes (authenticated; ownership enforced in handlers)
	auth.HandleFunc(constants.RouteUserPortfolios, portfolioController.GetPortfolios).Methods(constants.MethodGET)
	auth.HandleFunc(constants.RouteUserPortfolios, portfolioController.CreatePortfolio).Methods(constants.MethodPOST)
	auth.HandleFunc(constants.RouteUserPortfolioByID, portfolioController.GetPortfolio).Methods(constants.MethodGET)
	auth.HandleFunc(constants.RouteUserPortfolioByID, portfolioController.UpdatePortfolio).Methods(constants.MethodPUT)
	auth.HandleFunc(constants.RouteUserPortfolioByID, portfolioController.DeletePortfolio).Methods(constants.MethodDELETE)
	auth.HandleFunc(constants.RouteUserPortfolioHoldings, portfolioController.GetHoldings).Methods(constants.MethodGET)
	auth.HandleFunc(constants.RouteUserPortfolioHoldings, portfolioController.AddHolding).Methods(constants.MethodPOST)
	auth.HandleFunc(constants.RouteUserPortfolioHoldingByID, portfolioController.UpdateHolding).Methods(constants.MethodPUT)
	auth.HandleFunc(constants.RouteUserPortfolioHoldingByID, portfolioController.DeleteHolding).Methods(constants.MethodDELETE)
	auth.HandleFunc(constants.RouteUserPortfolioValue, portfolioController.GetPortfolioValue).Methods(constants.MethodGET)

	// Watchlist routes (authenticated)
	auth.HandleFunc(constants.RouteUserWatchlists, watchlistController.GetWatchlists).Methods(constants.MethodGET)
	auth.HandleFunc(constants.RouteUserWatchlists, watchlistController.CreateWatchlist).Methods(constants.MethodPOST)
	auth.HandleFunc(constants.RouteUserWatchlistByID, watchlistController.DeleteWatchlist).Methods(constants.MethodDELETE)
	auth.HandleFunc(constants.RouteUserWatchlistItems, watchlistController.GetItems).Methods(constants.MethodGET)
	auth.HandleFunc(constants.RouteUserWatchlistItems, watchlistController.AddStock).Methods(constants.MethodPOST)
	auth.HandleFunc(constants.RouteUserWatchlistItemByID, watchlistController.RemoveStock).Methods(constants.MethodDELETE)

	// User Alert Rule routes (authenticated)
	auth.HandleFunc(constants.RouteUserAlertRules, userAlertRuleController.GetAlertRules).Methods(constants.MethodGET)
	auth.HandleFunc(constants.RouteUserAlertRules, userAlertRuleController.CreateAlertRule).Methods(constants.MethodPOST)
	auth.HandleFunc(constants.RouteUserAlertRuleByID, userAlertRuleController.DeleteAlertRule).Methods(constants.MethodDELETE)

	// Tax Lot routes (authenticated; ownership enforced in handlers)
	auth.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-lots/buy", taxLotController.RecordBuy).Methods(constants.MethodPOST)
	auth.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-lots/sell-fifo", taxLotController.RecordSellFIFO).Methods(constants.MethodPOST)
	auth.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-lots/sell-lifo", taxLotController.RecordSellLIFO).Methods(constants.MethodPOST)
	auth.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-lots/{taxLotID}/sell", taxLotController.RecordSellSpecificLot).Methods(constants.MethodPOST)
	auth.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-lots/{taxLotID}/gains", taxLotController.GetTaxLotGains).Methods(constants.MethodGET)
	auth.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-gains", taxLotController.GetPortfolioTaxGains).Methods(constants.MethodGET)
	auth.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-report", taxLotController.GetTaxableGains).Methods(constants.MethodGET)
	auth.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-transactions", taxLotController.GetTaxTransactions).Methods(constants.MethodGET)

	return r
}
