package controllers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/constants"
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
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidStockIDFormat, err)
		return
	}

	if id <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgStockIDMustBePositive, nil)
		return
	}

	stock, err := c.service.GetStock(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgStockNotFound, nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToRetrieveStock, err)
		return
	}

	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: stock})
}

func (c *StockController) GetAllStocks(w http.ResponseWriter, r *http.Request) {
	stocks, err := c.service.GetAllStocks()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToRetrieveStocks, err)
		return
	}

	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: stocks})
}
