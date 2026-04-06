package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
		log.Printf("Server starting on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Println("Server is shutting down...")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	close(done)
	log.Println("Server exited")
}
