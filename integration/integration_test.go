package integration_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/vikhyat-sharma/quant-trading-prediction-system/constants"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/controllers"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/repositories"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/routes"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/services"
)

type stockResponse struct {
	ID     int    `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

type predictionResponse struct {
	ID             int       `json:"id"`
	StockID        int       `json:"stock_id"`
	PredictedPrice float64   `json:"predicted_price"`
	Date           time.Time `json:"date"`
}

func getTestDatabaseURL() string {
	if url := os.Getenv("TEST_DATABASE_URL"); url != "" {
		return url
	}
	return os.Getenv("DATABASE_URL")
}

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	url := getTestDatabaseURL()
	if url == "" {
		t.Skip("TEST_DATABASE_URL or DATABASE_URL is required for integration tests")
	}

	database, err := db.NewDB(url)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	// Ensure schema exists
	if err := db.EnsureSchema(database); err != nil {
		database.Close()
		t.Fatalf("failed to ensure schema: %v", err)
	}

	return database
}

func cleanupTestData(t *testing.T, database *sql.DB, symbol string) {
	t.Helper()

	_, err := database.Exec(`DELETE FROM price_history WHERE stock_id IN (SELECT id FROM stocks WHERE symbol = $1)`, symbol)
	if err != nil {
		t.Fatalf("failed to clean price_history: %v", err)
	}

	_, err = database.Exec(`DELETE FROM predictions WHERE stock_id IN (SELECT id FROM stocks WHERE symbol = $1)`, symbol)
	if err != nil {
		t.Fatalf("failed to clean predictions: %v", err)
	}

	_, err = database.Exec(`DELETE FROM stocks WHERE symbol = $1`, symbol)
	if err != nil {
		t.Fatalf("failed to clean stocks: %v", err)
	}
}

func createSampleStock(t *testing.T, database *sql.DB, symbol, name string) int {
	t.Helper()

	var id int
	row := database.QueryRow(`INSERT INTO stocks (symbol, name) VALUES ($1, $2) ON CONFLICT (symbol) DO UPDATE SET name = EXCLUDED.name RETURNING id`, symbol, name)
	if err := row.Scan(&id); err != nil {
		t.Fatalf("failed to insert sample stock: %v", err)
	}

	return id
}

func createSamplePrediction(t *testing.T, database *sql.DB, stockID int, price float64) int {
	t.Helper()

	var id int
	row := database.QueryRow(`INSERT INTO predictions (stock_id, predicted_price, date) VALUES ($1, $2, $3) RETURNING id`, stockID, price, time.Now())
	if err := row.Scan(&id); err != nil {
		t.Fatalf("failed to insert sample prediction: %v", err)
	}

	return id
}

func buildRouter(database *sql.DB) http.Handler {
	stockRepo := repositories.NewStockRepository(database)
	predictionRepo := repositories.NewPredictionRepository(database)
	priceHistoryRepo := repositories.NewPriceHistoryRepository(database)
	alertRepo := repositories.NewAlertRepository(database)
	notificationRepo := repositories.NewNotificationRepository(database)

	stockService := services.NewStockService(stockRepo)
	predictionService := services.NewPredictionService(predictionRepo)
	priceHistoryService := services.NewPriceHistoryService(priceHistoryRepo)
	alertService := services.NewAlertService(alertRepo, notificationRepo, priceHistoryRepo, stockRepo)

	stockController := controllers.NewStockController(stockService)
	predictionController := controllers.NewPredictionController(predictionService)
	priceHistoryController := controllers.NewPriceHistoryController(priceHistoryService)
	alertController := controllers.NewAlertController(alertService)

	return routes.SetupRoutes(stockController, predictionController, priceHistoryController, alertController)
}

func TestIntegration_GetAllStocks(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	symbol := fmt.Sprintf("TEST_%d", time.Now().UnixNano())
	cleanupTestData(t, db, symbol)
	defer cleanupTestData(t, db, symbol)

	stockID := createSampleStock(t, db, symbol, "Integration Test Stock")

	router := buildRouter(db)
	req := httptest.NewRequest(constants.MethodGET, constants.RouteStocks, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var response struct {
		Data []stockResponse `json:"data"`
	}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response.Data) == 0 {
		t.Fatal("expected at least one stock in response")
	}

	found := false
	for _, stock := range response.Data {
		if stock.ID == stockID && stock.Symbol == symbol {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("expected inserted stock %s to be present", symbol)
	}
}

func TestIntegration_GetStockByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	symbol := fmt.Sprintf("TEST_%d", time.Now().UnixNano())
	cleanupTestData(t, db, symbol)
	defer cleanupTestData(t, db, symbol)

	stockID := createSampleStock(t, db, symbol, "Integration Test Stock")

	router := buildRouter(db)
	url := fmt.Sprintf("/stocks/%d", stockID)
	req := httptest.NewRequest(constants.MethodGET, url, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var response struct {
		Data stockResponse `json:"data"`
	}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Data.ID != stockID || response.Data.Symbol != symbol {
		t.Fatalf("unexpected stock response: %+v", response.Data)
	}
}

func TestIntegration_GetPredictionsByStockID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	symbol := fmt.Sprintf("TEST_%d", time.Now().UnixNano())
	cleanupTestData(t, db, symbol)
	defer cleanupTestData(t, db, symbol)

	stockID := createSampleStock(t, db, symbol, "Integration Test Stock")
	createSamplePrediction(t, db, stockID, 123.45)

	router := buildRouter(db)
	url := fmt.Sprintf("/stocks/%d/predictions", stockID)
	req := httptest.NewRequest(constants.MethodGET, url, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var response struct {
		Data []predictionResponse `json:"data"`
	}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response.Data) == 0 {
		t.Fatal("expected at least one prediction in response")
	}

	if response.Data[0].StockID != stockID {
		t.Fatalf("expected stock_id %d, got %d", stockID, response.Data[0].StockID)
	}
}

func TestIntegration_GeneratePrediction(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	symbol := fmt.Sprintf("TEST_%d", time.Now().UnixNano())
	cleanupTestData(t, db, symbol)
	defer cleanupTestData(t, db, symbol)

	stockID := createSampleStock(t, db, symbol, "Integration Test Stock")

	router := buildRouter(db)
	url := fmt.Sprintf("/stocks/%d/predictions/generate", stockID)
	req := httptest.NewRequest(constants.MethodPOST, url, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	var response struct {
		Data predictionResponse `json:"data"`
	}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Data.StockID != stockID {
		t.Fatalf("expected stock_id %d, got %d", stockID, response.Data.StockID)
	}

	if response.Data.PredictedPrice != 100.0 {
		t.Fatalf("expected predicted_price 100.0, got %f", response.Data.PredictedPrice)
	}
}
