package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/services"
)

// TaxLotController handles tax lot operations
type TaxLotController struct {
	taxLotService *services.TaxLotService
}

// NewTaxLotController creates a new tax lot controller
func NewTaxLotController(taxLotService *services.TaxLotService) *TaxLotController {
	return &TaxLotController{
		taxLotService: taxLotService,
	}
}

// RecordBuy records a new buy transaction and creates a tax lot
// POST /users/{userID}/portfolios/{portfolioID}/tax-lots/buy
func (c *TaxLotController) RecordBuy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	portfolioID, err := strconv.Atoi(vars["portfolioID"])
	if err != nil {
		http.Error(w, "Invalid portfolio ID", http.StatusBadRequest)
		return
	}

	var req struct {
		StockID  int       `json:"stock_id"`
		Quantity float64   `json:"quantity"`
		Price    float64   `json:"price"`
		Fees     float64   `json:"fees"`
		BuyDate  time.Time `json:"buy_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.BuyDate.IsZero() {
		req.BuyDate = time.Now()
	}

	taxLot, err := c.taxLotService.RecordBuy(portfolioID, req.StockID, req.Quantity, req.Price, req.Fees, req.BuyDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"tax_lot": taxLot,
		"message": "Buy transaction recorded successfully",
	})
}

// RecordSellFIFO records a sell transaction using FIFO method
// POST /users/{userID}/portfolios/{portfolioID}/tax-lots/sell-fifo
func (c *TaxLotController) RecordSellFIFO(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	portfolioID, err := strconv.Atoi(vars["portfolioID"])
	if err != nil {
		http.Error(w, "Invalid portfolio ID", http.StatusBadRequest)
		return
	}

	var req struct {
		StockID  int       `json:"stock_id"`
		Quantity float64   `json:"quantity"`
		Price    float64   `json:"price"`
		Fees     float64   `json:"fees"`
		SellDate time.Time `json:"sell_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.SellDate.IsZero() {
		req.SellDate = time.Now()
	}

	realizedGain, err := c.taxLotService.RecordSellFIFO(portfolioID, req.StockID, req.Quantity, req.Price, req.Fees, req.SellDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":        "success",
		"realized_gain": realizedGain,
		"method":        "FIFO",
		"message":       "Sell transaction recorded using FIFO method",
	})
}

// RecordSellLIFO records a sell transaction using LIFO method
// POST /users/{userID}/portfolios/{portfolioID}/tax-lots/sell-lifo
func (c *TaxLotController) RecordSellLIFO(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	portfolioID, err := strconv.Atoi(vars["portfolioID"])
	if err != nil {
		http.Error(w, "Invalid portfolio ID", http.StatusBadRequest)
		return
	}

	var req struct {
		StockID  int       `json:"stock_id"`
		Quantity float64   `json:"quantity"`
		Price    float64   `json:"price"`
		Fees     float64   `json:"fees"`
		SellDate time.Time `json:"sell_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.SellDate.IsZero() {
		req.SellDate = time.Now()
	}

	realizedGain, err := c.taxLotService.RecordSellLIFO(portfolioID, req.StockID, req.Quantity, req.Price, req.Fees, req.SellDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":        "success",
		"realized_gain": realizedGain,
		"method":        "LIFO",
		"message":       "Sell transaction recorded using LIFO method",
	})
}

// RecordSellSpecificLot records a sell transaction from a specific tax lot
// POST /users/{userID}/portfolios/{portfolioID}/tax-lots/{taxLotID}/sell
func (c *TaxLotController) RecordSellSpecificLot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taxLotID, err := strconv.Atoi(vars["taxLotID"])
	if err != nil {
		http.Error(w, "Invalid tax lot ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Quantity float64   `json:"quantity"`
		Price    float64   `json:"price"`
		Fees     float64   `json:"fees"`
		SellDate time.Time `json:"sell_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.SellDate.IsZero() {
		req.SellDate = time.Now()
	}

	realizedGain, err := c.taxLotService.RecordSellSpecificLot(taxLotID, req.Quantity, req.Price, req.Fees, req.SellDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":        "success",
		"realized_gain": realizedGain,
		"method":        "SPECIFIC_LOT",
		"message":       "Sell transaction recorded from specific tax lot",
	})
}

// GetTaxLotGains gets realized and unrealized gains for a tax lot
// GET /users/{userID}/portfolios/{portfolioID}/tax-lots/{taxLotID}/gains
func (c *TaxLotController) GetTaxLotGains(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taxLotID, err := strconv.Atoi(vars["taxLotID"])
	if err != nil {
		http.Error(w, "Invalid tax lot ID", http.StatusBadRequest)
		return
	}

	currentPrice := r.URL.Query().Get("current_price")
	if currentPrice == "" {
		http.Error(w, "current_price query parameter is required", http.StatusBadRequest)
		return
	}

	price, err := strconv.ParseFloat(currentPrice, 64)
	if err != nil {
		http.Error(w, "Invalid current_price value", http.StatusBadRequest)
		return
	}

	gains, err := c.taxLotService.GetTaxLotGains(taxLotID, price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"gains":  gains,
	})
}

// GetPortfolioTaxGains gets total tax gains for a portfolio
// GET /users/{userID}/portfolios/{portfolioID}/tax-gains
func (c *TaxLotController) GetPortfolioTaxGains(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	portfolioID, err := strconv.Atoi(vars["portfolioID"])
	if err != nil {
		http.Error(w, "Invalid portfolio ID", http.StatusBadRequest)
		return
	}

	// Parse current prices from query parameters
	// Format: ?prices=1:100,2:150,3:200 (stockID:price pairs)
	pricesQuery := r.URL.Query().Get("prices")
	currentPrices := make(map[int]float64)

	if pricesQuery != "" {
		// Parse the prices query parameter
		// This is a simple implementation; can be enhanced
		// You can also accept JSON body for complex scenarios
		http.Error(w, "prices parameter format not specified. Use JSON body instead.", http.StatusBadRequest)
		return
	}

	// For now, accept prices in JSON body
	var req struct {
		CurrentPrices map[int]float64 `json:"current_prices"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	gains, err := c.taxLotService.GetPortfolioTaxGains(portfolioID, req.CurrentPrices)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   gains,
	})
}

// GetTaxableGains gets tax consequences by holding period
// GET /users/{userID}/portfolios/{portfolioID}/tax-report
func (c *TaxLotController) GetTaxableGains(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	portfolioID, err := strconv.Atoi(vars["portfolioID"])
	if err != nil {
		http.Error(w, "Invalid portfolio ID", http.StatusBadRequest)
		return
	}

	gains, err := c.taxLotService.CalculateTaxableGainsBySellDate(portfolioID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   gains,
	})
}

// GetTaxTransactions gets all tax transactions for a portfolio
// GET /users/{userID}/portfolios/{portfolioID}/tax-transactions
func (c *TaxLotController) GetTaxTransactions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	portfolioID, err := strconv.Atoi(vars["portfolioID"])
	if err != nil {
		http.Error(w, "Invalid portfolio ID", http.StatusBadRequest)
		return
	}

	transactions, err := c.taxLotService.GetTaxTransactionsByPortfolio(portfolioID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":       "success",
		"transactions": transactions,
		"total_count":  len(transactions),
	})
}
