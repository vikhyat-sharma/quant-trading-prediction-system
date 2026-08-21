package routes

import (
	"github.com/gorilla/mux"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/constants"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/controllers"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/middleware"
)

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
) *mux.Router {
	r := mux.NewRouter()

	// Apply middleware to all routes
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.CORSMiddleware)
	r.Use(middleware.ContentTypeMiddleware)

	r.HandleFunc("/health", middleware.HealthHandler).Methods(constants.MethodGET)
	r.HandleFunc("/ready", middleware.HealthHandler).Methods(constants.MethodGET)

	protectedUserRoutes := r.PathPrefix("/users").Subrouter()
	protectedUserRoutes.Use(middleware.AuthMiddleware)

	protectedPortfolioRoutes := r.PathPrefix("/users/{userID}/portfolios").Subrouter()
	protectedPortfolioRoutes.Use(middleware.AuthMiddleware)

	protectedAlertRoutes := r.PathPrefix("/stocks/{stockID}/alerts").Subrouter()
	protectedAlertRoutes.Use(middleware.AuthMiddleware)

	// Stock routes
	r.HandleFunc(constants.RouteStocks, stockController.GetAllStocks).Methods(constants.MethodGET)
	r.HandleFunc(constants.RouteStocks, stockController.CreateStock).Methods(constants.MethodPOST)
	r.HandleFunc(constants.RouteStockByID, stockController.GetStock).Methods(constants.MethodGET)
	r.HandleFunc(constants.RouteStockByID, stockController.UpdateStock).Methods(constants.MethodPUT)
	r.HandleFunc(constants.RouteStockByID, stockController.DeleteStock).Methods(constants.MethodDELETE)

	// Prediction routes
	r.HandleFunc(constants.RouteStockPredictions, predictionController.GetPredictions).Methods(constants.MethodGET)
	r.HandleFunc(constants.RouteStockPredictionsGenerate, predictionController.GeneratePrediction).Methods(constants.MethodPOST)
	r.HandleFunc(constants.RouteStockBacktest, predictionController.BacktestStrategy).Methods(constants.MethodGET)
	r.HandleFunc(constants.RouteStockSentiment, sentimentController.AnalyzeSentiment).Methods(constants.MethodPOST)

	// Price History routes
	r.HandleFunc(constants.RouteStockPriceHistory, priceHistoryController.GetPriceHistory).Methods(constants.MethodGET)
	r.HandleFunc(constants.RouteStockPriceHistory, priceHistoryController.RecordPrice).Methods(constants.MethodPOST)
	r.HandleFunc(constants.RouteStockPriceHistoryRange, priceHistoryController.GetPriceHistoryByDateRange).Methods(constants.MethodGET)
	r.HandleFunc(constants.RouteStockPriceStats, priceHistoryController.GetPriceStats).Methods(constants.MethodGET)
	r.HandleFunc(constants.RouteStockLatestPrice, priceHistoryController.GetLatestPrice).Methods(constants.MethodGET)

	// Alert routes
	protectedAlertRoutes.HandleFunc(constants.RouteStockAlerts, alertController.GetAlerts).Methods(constants.MethodGET)
	protectedAlertRoutes.HandleFunc(constants.RouteStockAlerts, alertController.CreateAlert).Methods(constants.MethodPOST)
	protectedAlertRoutes.HandleFunc(constants.RouteStockAlertByID, alertController.DeleteAlert).Methods(constants.MethodDELETE)
	protectedAlertRoutes.HandleFunc(constants.RouteStockAlertsEvaluate, alertController.EvaluateAlerts).Methods(constants.MethodPOST)
	protectedAlertRoutes.HandleFunc(constants.RouteStockNotifications, alertController.GetNotifications).Methods(constants.MethodGET)

	// User routes
	r.HandleFunc(constants.RouteUsers, userController.CreateUser).Methods(constants.MethodPOST)
	protectedUserRoutes.HandleFunc(constants.RouteUsers, userController.GetUsers).Methods(constants.MethodGET)
	protectedUserRoutes.HandleFunc(constants.RouteUserByID, userController.GetUser).Methods(constants.MethodGET)
	protectedUserRoutes.HandleFunc(constants.RouteUserByID, userController.UpdateUser).Methods(constants.MethodPUT)
	protectedUserRoutes.HandleFunc(constants.RouteUserByID, userController.DeleteUser).Methods(constants.MethodDELETE)

	// Portfolio routes
	protectedPortfolioRoutes.HandleFunc(constants.RouteUserPortfolios, portfolioController.GetPortfolios).Methods(constants.MethodGET)
	protectedPortfolioRoutes.HandleFunc(constants.RouteUserPortfolios, portfolioController.CreatePortfolio).Methods(constants.MethodPOST)
	protectedPortfolioRoutes.HandleFunc(constants.RouteUserPortfolioByID, portfolioController.GetPortfolio).Methods(constants.MethodGET)
	protectedPortfolioRoutes.HandleFunc(constants.RouteUserPortfolioByID, portfolioController.UpdatePortfolio).Methods(constants.MethodPUT)
	protectedPortfolioRoutes.HandleFunc(constants.RouteUserPortfolioByID, portfolioController.DeletePortfolio).Methods(constants.MethodDELETE)
	protectedPortfolioRoutes.HandleFunc(constants.RouteUserPortfolioHoldings, portfolioController.GetHoldings).Methods(constants.MethodGET)
	protectedPortfolioRoutes.HandleFunc(constants.RouteUserPortfolioHoldings, portfolioController.AddHolding).Methods(constants.MethodPOST)
	protectedPortfolioRoutes.HandleFunc(constants.RouteUserPortfolioHoldingByID, portfolioController.UpdateHolding).Methods(constants.MethodPUT)
	protectedPortfolioRoutes.HandleFunc(constants.RouteUserPortfolioHoldingByID, portfolioController.DeleteHolding).Methods(constants.MethodDELETE)
	protectedPortfolioRoutes.HandleFunc(constants.RouteUserPortfolioValue, portfolioController.GetPortfolioValue).Methods(constants.MethodGET)

	// User Watchlist routes
	protectedPortfolioRoutes.HandleFunc(constants.RouteUserWatchlists, watchlistController.GetWatchlists).Methods(constants.MethodGET)
	protectedPortfolioRoutes.HandleFunc(constants.RouteUserWatchlists, watchlistController.CreateWatchlist).Methods(constants.MethodPOST)
	protectedPortfolioRoutes.HandleFunc(constants.RouteUserWatchlistByID, watchlistController.DeleteWatchlist).Methods(constants.MethodDELETE)
	protectedPortfolioRoutes.HandleFunc(constants.RouteUserWatchlistItems, watchlistController.GetItems).Methods(constants.MethodGET)
	protectedPortfolioRoutes.HandleFunc(constants.RouteUserWatchlistItems, watchlistController.AddStock).Methods(constants.MethodPOST)
	protectedPortfolioRoutes.HandleFunc(constants.RouteUserWatchlistItemByID, watchlistController.RemoveStock).Methods(constants.MethodDELETE)

	// User Alert Rule routes
	protectedUserRoutes.HandleFunc(constants.RouteUserAlertRules, userAlertRuleController.GetAlertRules).Methods(constants.MethodGET)
	protectedUserRoutes.HandleFunc(constants.RouteUserAlertRules, userAlertRuleController.CreateAlertRule).Methods(constants.MethodPOST)
	protectedUserRoutes.HandleFunc(constants.RouteUserAlertRuleByID, userAlertRuleController.DeleteAlertRule).Methods(constants.MethodDELETE)

	// Tax Lot routes
	protectedPortfolioRoutes.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-lots/buy", taxLotController.RecordBuy).Methods(constants.MethodPOST)
	protectedPortfolioRoutes.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-lots/sell-fifo", taxLotController.RecordSellFIFO).Methods(constants.MethodPOST)
	protectedPortfolioRoutes.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-lots/sell-lifo", taxLotController.RecordSellLIFO).Methods(constants.MethodPOST)
	protectedPortfolioRoutes.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-lots/{taxLotID}/sell", taxLotController.RecordSellSpecificLot).Methods(constants.MethodPOST)
	protectedPortfolioRoutes.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-lots/{taxLotID}/gains", taxLotController.GetTaxLotGains).Methods(constants.MethodGET)
	protectedPortfolioRoutes.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-gains", taxLotController.GetPortfolioTaxGains).Methods(constants.MethodGET)
	protectedPortfolioRoutes.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-report", taxLotController.GetTaxableGains).Methods(constants.MethodGET)
	protectedPortfolioRoutes.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-transactions", taxLotController.GetTaxTransactions).Methods(constants.MethodGET)

	return r
}
