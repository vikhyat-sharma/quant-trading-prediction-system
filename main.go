package main

import (
	"log"
	"net/http"

	"github.com/vikhyat-sharma/quant-trading-prediction-system/config"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/controllers"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/repositories"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/routes"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/services"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	database, err := db.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	stockRepo := repositories.NewStockRepository(database)
	predictionRepo := repositories.NewPredictionRepository(database)

	stockService := services.NewStockService(stockRepo)
	predictionService := services.NewPredictionService(predictionRepo)

	stockController := controllers.NewStockController(stockService)
	predictionController := controllers.NewPredictionController(predictionService)

	router := routes.SetupRoutes(stockController, predictionController)

	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
