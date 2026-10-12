package controllers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/constants"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/middleware"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/services"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/util"
)

type WatchlistController struct {
	service *services.WatchlistService
}

func NewWatchlistController(service *services.WatchlistService) *WatchlistController {
	return &WatchlistController{service: service}
}

// ownsUser returns true if the caller owns the given userID or is an admin.
func ownsUser(r *http.Request, userID int) bool {
	callerID := middleware.ContextUserID(r.Context())
	return callerID == userID || util.IsAdminRole(middleware.ContextUserRole(r.Context()))
}

func (c *WatchlistController) CreateWatchlist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userID"])
	if err != nil || userID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidUserIDFormat, nil)
		return
	}
	if !ownsUser(r, userID) {
		writeErrorResponse(w, http.StatusForbidden, "Insufficient permissions", nil)
		return
	}

	var payload struct {
		Name string `json:"name"`
	}
	if err := parseJSONBody(r, &payload); err != nil || payload.Name == "" {
		writeErrorResponse(w, http.StatusBadRequest, "name is required", nil)
		return
	}

	wl := &db.UserWatchlist{UserID: userID, Name: payload.Name}
	created, err := c.service.CreateWatchlist(wl)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to create watchlist", nil)
		return
	}
	writeJSONResponse(w, http.StatusCreated, SuccessResponse{Data: created})
}

func (c *WatchlistController) GetWatchlists(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userID"])
	if err != nil || userID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidUserIDFormat, nil)
		return
	}
	if !ownsUser(r, userID) {
		writeErrorResponse(w, http.StatusForbidden, "Insufficient permissions", nil)
		return
	}

	lists, err := c.service.GetWatchlistsByUser(userID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to get watchlists", nil)
		return
	}
	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: lists})
}

func (c *WatchlistController) AddStock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userID"])
	if err != nil || userID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidUserIDFormat, nil)
		return
	}
	if !ownsUser(r, userID) {
		writeErrorResponse(w, http.StatusForbidden, "Insufficient permissions", nil)
		return
	}

	watchlistID, err := strconv.Atoi(vars["watchlistID"])
	if err != nil || watchlistID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid watchlist ID", nil)
		return
	}

	var payload struct {
		StockID int `json:"stock_id"`
	}
	if err := parseJSONBody(r, &payload); err != nil || payload.StockID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, "stock_id is required", nil)
		return
	}

	item := &db.WatchlistItem{WatchlistID: watchlistID, StockID: payload.StockID}
	created, err := c.service.AddStockToWatchlist(item)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to add stock to watchlist", nil)
		return
	}
	writeJSONResponse(w, http.StatusCreated, SuccessResponse{Data: created})
}

func (c *WatchlistController) GetItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userID"])
	if err != nil || userID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidUserIDFormat, nil)
		return
	}
	if !ownsUser(r, userID) {
		writeErrorResponse(w, http.StatusForbidden, "Insufficient permissions", nil)
		return
	}

	watchlistID, err := strconv.Atoi(vars["watchlistID"])
	if err != nil || watchlistID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid watchlist ID", nil)
		return
	}

	items, err := c.service.GetWatchlistItems(watchlistID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to get watchlist items", nil)
		return
	}
	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: items})
}

func (c *WatchlistController) RemoveStock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userID"])
	if err != nil || userID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidUserIDFormat, nil)
		return
	}
	if !ownsUser(r, userID) {
		writeErrorResponse(w, http.StatusForbidden, "Insufficient permissions", nil)
		return
	}

	watchlistID, err := strconv.Atoi(vars["watchlistID"])
	stockID, err2 := strconv.Atoi(vars["stockID"])
	if err != nil || err2 != nil || watchlistID <= 0 || stockID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid IDs", nil)
		return
	}

	if err := c.service.RemoveStockFromWatchlist(watchlistID, stockID); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to remove stock from watchlist", nil)
		return
	}
	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: map[string]string{"message": "Stock removed from watchlist"}})
}

func (c *WatchlistController) DeleteWatchlist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userID"])
	if err != nil || userID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidUserIDFormat, nil)
		return
	}
	if !ownsUser(r, userID) {
		writeErrorResponse(w, http.StatusForbidden, "Insufficient permissions", nil)
		return
	}

	watchlistID, err := strconv.Atoi(vars["watchlistID"])
	if err != nil || watchlistID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid watchlist ID", nil)
		return
	}

	if err := c.service.DeleteWatchlist(watchlistID); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to delete watchlist", nil)
		return
	}
	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: map[string]string{"message": "Watchlist deleted"}})
}
