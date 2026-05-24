package controllers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/constants"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/services"
)

type WatchlistController struct {
	service *services.WatchlistService
}

func NewWatchlistController(service *services.WatchlistService) *WatchlistController {
	return &WatchlistController{service: service}
}

func (c *WatchlistController) CreateWatchlist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidUserIDFormat, err)
		return
	}
	var payload struct {
		Name string `json:"name"`
	}
	if err := parseJSONBody(r, &payload); err != nil || payload.Name == "" {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	wl := &db.UserWatchlist{UserID: userID, Name: payload.Name}
	created, err := c.service.CreateWatchlist(wl)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to create watchlist", err)
		return
	}
	writeJSONResponse(w, http.StatusCreated, SuccessResponse{Data: created})
}

func (c *WatchlistController) GetWatchlists(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidUserIDFormat, err)
		return
	}
	lists, err := c.service.GetWatchlistsByUser(userID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to get watchlists", err)
		return
	}
	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: lists})
}

func (c *WatchlistController) AddStock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	watchlistIDStr := vars["watchlistID"]
	watchlistID, err := strconv.Atoi(watchlistIDStr)
	if err != nil || watchlistID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid watchlist ID", err)
		return
	}
	var payload struct {
		StockID int `json:"stock_id"`
	}
	if err := parseJSONBody(r, &payload); err != nil || payload.StockID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	item := &db.WatchlistItem{WatchlistID: watchlistID, StockID: payload.StockID}
	created, err := c.service.AddStockToWatchlist(item)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to add stock to watchlist", err)
		return
	}
	writeJSONResponse(w, http.StatusCreated, SuccessResponse{Data: created})
}

func (c *WatchlistController) GetItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	watchlistIDStr := vars["watchlistID"]
	watchlistID, err := strconv.Atoi(watchlistIDStr)
	if err != nil || watchlistID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid watchlist ID", err)
		return
	}
	items, err := c.service.GetWatchlistItems(watchlistID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to get watchlist items", err)
		return
	}
	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: items})
}

func (c *WatchlistController) RemoveStock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	watchlistIDStr := vars["watchlistID"]
	stockIDStr := vars["stockID"]
	watchlistID, err := strconv.Atoi(watchlistIDStr)
	stockID, err2 := strconv.Atoi(stockIDStr)
	if err != nil || err2 != nil || watchlistID <= 0 || stockID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid IDs", nil)
		return
	}
	if err := c.service.RemoveStockFromWatchlist(watchlistID, stockID); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to remove stock from watchlist", err)
		return
	}
	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: map[string]string{"message": "Stock removed from watchlist"}})
}

func (c *WatchlistController) DeleteWatchlist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	watchlistIDStr := vars["watchlistID"]
	watchlistID, err := strconv.Atoi(watchlistIDStr)
	if err != nil || watchlistID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid watchlist ID", err)
		return
	}
	if err := c.service.DeleteWatchlist(watchlistID); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to delete watchlist", err)
		return
	}
	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: map[string]string{"message": "Watchlist deleted"}})
}
