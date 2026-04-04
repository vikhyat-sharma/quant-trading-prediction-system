package routes

import (
	"github.com/gorilla/mux"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/controllers"
)

func SetupRoutes(stockController *controllers.StockController, predictionController *controllers.PredictionController) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/stocks", stockController.GetAllStocks).Methods("GET")
	r.HandleFunc("/stocks/{id}", stockController.GetStock).Methods("GET")
	r.HandleFunc("/stocks/{stockID}/predictions", predictionController.GetPredictions).Methods("GET")
	r.HandleFunc("/stocks/{stockID}/predictions/generate", predictionController.GeneratePrediction).Methods("POST")
	return r
}
