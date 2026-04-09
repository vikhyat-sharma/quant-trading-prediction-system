package routes

import (
	"github.com/gorilla/mux"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/constants"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/controllers"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/middleware"
)

func SetupRoutes(stockController *controllers.StockController, predictionController *controllers.PredictionController, priceHistoryController *controllers.PriceHistoryController) *mux.Router {
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

	// Price History routes
	r.HandleFunc(constants.RouteStockPriceHistory, priceHistoryController.GetPriceHistory).Methods(constants.MethodGET)
	r.HandleFunc(constants.RouteStockPriceHistory, priceHistoryController.RecordPrice).Methods(constants.MethodPOST)
	r.HandleFunc(constants.RouteStockPriceHistoryRange, priceHistoryController.GetPriceHistoryByDateRange).Methods(constants.MethodGET)
	r.HandleFunc(constants.RouteStockPriceStats, priceHistoryController.GetPriceStats).Methods(constants.MethodGET)
	r.HandleFunc(constants.RouteStockLatestPrice, priceHistoryController.GetLatestPrice).Methods(constants.MethodGET)

	return r
}
