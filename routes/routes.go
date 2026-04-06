package routes

import (
	"github.com/gorilla/mux"
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
	r.HandleFunc("/stocks", stockController.GetAllStocks).Methods("GET")
	r.HandleFunc("/stocks/{id}", stockController.GetStock).Methods("GET")

	// Prediction routes
	r.HandleFunc("/stocks/{stockID}/predictions", predictionController.GetPredictions).Methods("GET")
	r.HandleFunc("/stocks/{stockID}/predictions/generate", predictionController.GeneratePrediction).Methods("POST")

	return r
}
