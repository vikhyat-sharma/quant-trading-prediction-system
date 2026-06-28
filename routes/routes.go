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
	r.HandleFunc(constants.RouteStockAlerts, alertController.GetAlerts).Methods(constants.MethodGET)
	r.HandleFunc(constants.RouteStockAlerts, alertController.CreateAlert).Methods(constants.MethodPOST)
	r.HandleFunc(constants.RouteStockAlertByID, alertController.DeleteAlert).Methods(constants.MethodDELETE)
	r.HandleFunc(constants.RouteStockAlertsEvaluate, alertController.EvaluateAlerts).Methods(constants.MethodPOST)
	r.HandleFunc(constants.RouteStockNotifications, alertController.GetNotifications).Methods(constants.MethodGET)

	// User routes
	r.HandleFunc(constants.RouteUsers, userController.GetUsers).Methods(constants.MethodGET)
	r.HandleFunc(constants.RouteUsers, userController.CreateUser).Methods(constants.MethodPOST)
	r.HandleFunc(constants.RouteUserByID, userController.GetUser).Methods(constants.MethodGET)
	r.HandleFunc(constants.RouteUserByID, userController.UpdateUser).Methods(constants.MethodPUT)
	r.HandleFunc(constants.RouteUserByID, userController.DeleteUser).Methods(constants.MethodDELETE)

	// Portfolio routes
	r.HandleFunc(constants.RouteUserPortfolios, portfolioController.GetPortfolios).Methods(constants.MethodGET)
	r.HandleFunc(constants.RouteUserPortfolios, portfolioController.CreatePortfolio).Methods(constants.MethodPOST)
	r.HandleFunc(constants.RouteUserPortfolioByID, portfolioController.GetPortfolio).Methods(constants.MethodGET)
	r.HandleFunc(constants.RouteUserPortfolioByID, portfolioController.UpdatePortfolio).Methods(constants.MethodPUT)
	r.HandleFunc(constants.RouteUserPortfolioByID, portfolioController.DeletePortfolio).Methods(constants.MethodDELETE)
	r.HandleFunc(constants.RouteUserPortfolioHoldings, portfolioController.GetHoldings).Methods(constants.MethodGET)
	r.HandleFunc(constants.RouteUserPortfolioHoldings, portfolioController.AddHolding).Methods(constants.MethodPOST)
	r.HandleFunc(constants.RouteUserPortfolioHoldingByID, portfolioController.UpdateHolding).Methods(constants.MethodPUT)
	r.HandleFunc(constants.RouteUserPortfolioHoldingByID, portfolioController.DeleteHolding).Methods(constants.MethodDELETE)
	// Portfolio value
	r.HandleFunc(constants.RouteUserPortfolioValue, portfolioController.GetPortfolioValue).Methods(constants.MethodGET)

	// User Watchlist routes
	r.HandleFunc(constants.RouteUserWatchlists, watchlistController.GetWatchlists).Methods(constants.MethodGET)
	r.HandleFunc(constants.RouteUserWatchlists, watchlistController.CreateWatchlist).Methods(constants.MethodPOST)
	r.HandleFunc(constants.RouteUserWatchlistByID, watchlistController.DeleteWatchlist).Methods(constants.MethodDELETE)
	r.HandleFunc(constants.RouteUserWatchlistItems, watchlistController.GetItems).Methods(constants.MethodGET)
	r.HandleFunc(constants.RouteUserWatchlistItems, watchlistController.AddStock).Methods(constants.MethodPOST)
	r.HandleFunc(constants.RouteUserWatchlistItemByID, watchlistController.RemoveStock).Methods(constants.MethodDELETE)

	// User Alert Rule routes
	r.HandleFunc(constants.RouteUserAlertRules, userAlertRuleController.GetAlertRules).Methods(constants.MethodGET)
	r.HandleFunc(constants.RouteUserAlertRules, userAlertRuleController.CreateAlertRule).Methods(constants.MethodPOST)
	r.HandleFunc(constants.RouteUserAlertRuleByID, userAlertRuleController.DeleteAlertRule).Methods(constants.MethodDELETE)

	// Tax Lot routes
	r.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-lots/buy", taxLotController.RecordBuy).Methods(constants.MethodPOST)
	r.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-lots/sell-fifo", taxLotController.RecordSellFIFO).Methods(constants.MethodPOST)
	r.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-lots/sell-lifo", taxLotController.RecordSellLIFO).Methods(constants.MethodPOST)
	r.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-lots/{taxLotID}/sell", taxLotController.RecordSellSpecificLot).Methods(constants.MethodPOST)
	r.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-lots/{taxLotID}/gains", taxLotController.GetTaxLotGains).Methods(constants.MethodGET)
	r.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-gains", taxLotController.GetPortfolioTaxGains).Methods(constants.MethodGET)
	r.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-report", taxLotController.GetTaxableGains).Methods(constants.MethodGET)
	r.HandleFunc("/users/{userID}/portfolios/{portfolioID}/tax-transactions", taxLotController.GetTaxTransactions).Methods(constants.MethodGET)

	return r
}
