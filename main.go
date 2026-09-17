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
	"github.com/vikhyat-sharma/quant-trading-prediction-system/middleware"
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

	if err := db.EnsureSchema(database); err != nil {
		log.Fatal(constants.LogMsgFailedToEnsureSchema+":", err)
	}

	// Repositories
	stockRepo := repositories.NewStockRepository(database)
	predictionRepo := repositories.NewPredictionRepository(database)
	priceHistoryRepo := repositories.NewPriceHistoryRepository(database)
	alertRepo := repositories.NewAlertRepository(database)
	notificationRepo := repositories.NewNotificationRepository(database)
	userRepo := repositories.NewUserRepository(database)
	portfolioRepo := repositories.NewPortfolioRepository(database)
	taxLotRepo := repositories.NewTaxLotRepository(database)
	watchlistRepo := repositories.NewWatchlistRepository(database)
	userAlertRuleRepo := repositories.NewUserAlertRuleRepository(database)

	// Services
	stockService := services.NewStockService(stockRepo)
	predictionService := services.NewPredictionService(predictionRepo, priceHistoryRepo)
	priceHistoryService := services.NewPriceHistoryService(priceHistoryRepo)
	alertService := services.NewAlertService(alertRepo, notificationRepo, priceHistoryRepo, stockRepo)
	userService := services.NewUserService(userRepo)
	portfolioService := services.NewPortfolioService(portfolioRepo)
	sentimentService := services.NewSentimentService()
	taxLotService := services.NewTaxLotService(taxLotRepo, stockRepo)
	watchlistService := services.NewWatchlistService(watchlistRepo)
	userAlertRuleService := services.NewUserAlertRuleService(userAlertRuleRepo)

	// Controllers
	stockController := controllers.NewStockController(stockService)
	predictionController := controllers.NewPredictionController(predictionService)
	priceHistoryController := controllers.NewPriceHistoryController(priceHistoryService)
	alertController := controllers.NewAlertController(alertService)
	userController := controllers.NewUserController(userService)
	portfolioController := controllers.NewPortfolioController(portfolioService, priceHistoryRepo)
	sentimentController := controllers.NewSentimentController(sentimentService)
	taxLotController := controllers.NewTaxLotController(taxLotService)
	watchlistController := controllers.NewWatchlistController(watchlistService)
	userAlertRuleController := controllers.NewUserAlertRuleController(userAlertRuleService)

	// Rate limiter (100 req/min per IP); stopped on shutdown
	rateLimiter := middleware.NewRateLimiter(100, 60)

	router := routes.SetupRoutes(
		stockController,
		predictionController,
		priceHistoryController,
		alertController,
		userController,
		portfolioController,
		sentimentController,
		watchlistController,
		userAlertRuleController,
		taxLotController,
		rateLimiter,
		database,
	)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf(constants.LogMsgServerStarting, cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(constants.LogMsgFailedToStartServer+":", err)
		}
	}()

	<-quit
	log.Println(constants.LogMsgServerShuttingDown)

	rateLimiter.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultServerShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(constants.LogMsgServerForcedShutdown+":", err)
	}

	log.Println(constants.LogMsgServerExited)
}
