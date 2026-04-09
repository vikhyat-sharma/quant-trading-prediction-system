package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/vikhyat-sharma/quant-trading-prediction-system/config"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/constants"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/controllers"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/repositories"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/routes"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/services"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(constants.LogMsgFailedToLoadConfig+":", err)
	}

	database, err := db.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(constants.LogMsgFailedToConnectDB+":", err)
	}
	defer database.Close()

	// Ensure database schema exists
	if err := db.EnsureSchema(database); err != nil {
		log.Fatal(constants.LogMsgFailedToEnsureSchema+":", err)
	}

	stockRepo := repositories.NewStockRepository(database)
	predictionRepo := repositories.NewPredictionRepository(database)
	priceHistoryRepo := repositories.NewPriceHistoryRepository(database)

	stockService := services.NewStockService(stockRepo)
	predictionService := services.NewPredictionService(predictionRepo)
	priceHistoryService := services.NewPriceHistoryService(priceHistoryRepo)

	stockController := controllers.NewStockController(stockService)
	predictionController := controllers.NewPredictionController(predictionService)
	priceHistoryController := controllers.NewPriceHistoryController(priceHistoryService)

	router := routes.SetupRoutes(stockController, predictionController, priceHistoryController)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Channel to listen for interrupt signal
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	// Register interrupt signals
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Printf(constants.LogMsgServerStarting, cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(constants.LogMsgFailedToStartServer+":", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Println(constants.LogMsgServerShuttingDown)

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultServerShutdownTimeout)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(constants.LogMsgServerForcedShutdown+":", err)
	}

	close(done)
	log.Println(constants.LogMsgServerExited)
}
