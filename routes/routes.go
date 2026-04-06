package routes

import (
	"github.com/gorilla/mux"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/constants"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/controllers"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/middleware"
)

func SetupRoutes(stockController *controllers.StockController, predictionController *controllers.PredictionController) *mux.Router {
	r := mux.NewRouter()

	// Apply middleware to all routes
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.CORSMiddleware)
	r.Use(middleware.ContentTypeMiddleware)

	// Stock routes
	r.HandleFunc(constants.RouteStocks, stockController.GetAllStocks).Methods(constants.MethodGET)
	r.HandleFunc(constants.RouteStockByID, stockController.GetStock).Methods(constants.MethodGET)

	// Prediction routes
	r.HandleFunc(constants.RouteStockPredictions, predictionController.GetPredictions).Methods(constants.MethodGET)
	r.HandleFunc(constants.RouteStockPredictionsGenerate, predictionController.GeneratePrediction).Methods(constants.MethodPOST)

	return r
}
