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

func (c *StockController) CreateStock(w http.ResponseWriter, r *http.Request) {
	var stock db.Stock
	if err := parseJSONBody(r, &stock); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	createdStock, err := c.service.CreateStock(&stock)
	if err != nil {
		if errors.Is(err, db.ErrInvalidSymbol) || errors.Is(err, db.ErrInvalidName) || errors.Is(err, db.ErrInvalidExchange) {
			writeErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to create stock", err)
		return
	}

	writeJSONResponse(w, http.StatusCreated, SuccessResponse{Data: createdStock})
}

func (c *StockController) UpdateStock(w http.ResponseWriter, r *http.Request) {
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

	var stock db.Stock
	if err := parseJSONBody(r, &stock); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	updatedStock, err := c.service.UpdateStock(id, &stock)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) || errors.Is(err, sql.ErrNoRows) {
			writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgStockNotFound, nil)
			return
		}
		if errors.Is(err, db.ErrInvalidSymbol) || errors.Is(err, db.ErrInvalidName) || errors.Is(err, db.ErrInvalidExchange) {
			writeErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to update stock", err)
		return
	}

	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: updatedStock})
}

func (c *StockController) DeleteStock(w http.ResponseWriter, r *http.Request) {
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

	if err := c.service.DeleteStock(id); err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgStockNotFound, nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to delete stock", err)
		return
	}

	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: map[string]string{"message": "Stock deleted successfully"}})
}
