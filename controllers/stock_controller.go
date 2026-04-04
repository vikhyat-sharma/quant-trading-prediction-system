package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/services"
)

type StockController struct {
	service *services.StockService
}

func NewStockController(service *services.StockService) *StockController {
	return &StockController{service: service}
}

func (c *StockController) GetStock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	stock, err := c.service.GetStock(id)
	if err != nil {
		http.Error(w, "Stock not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(stock)
}

func (c *StockController) GetAllStocks(w http.ResponseWriter, r *http.Request) {
	stocks, err := c.service.GetAllStocks()
	if err != nil {
		http.Error(w, "Error fetching stocks", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stocks)
}
