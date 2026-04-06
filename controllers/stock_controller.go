package controllers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
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
		writeErrorResponse(w, http.StatusBadRequest, "Invalid stock ID format", err)
		return
	}

	if id <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, "Stock ID must be a positive integer", nil)
		return
	}

	stock, err := c.service.GetStock(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, "Stock not found", nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve stock", err)
		return
	}

	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: stock})
}

func (c *StockController) GetAllStocks(w http.ResponseWriter, r *http.Request) {
	stocks, err := c.service.GetAllStocks()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve stocks", err)
		return
	}

	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: stocks})
}
