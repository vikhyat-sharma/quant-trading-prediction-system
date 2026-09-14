package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/constants"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/middleware"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/repositories"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/services"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/util"
)

type PortfolioController struct {
	service   *services.PortfolioService
	priceRepo *repositories.PriceHistoryRepository
}

func NewPortfolioController(service *services.PortfolioService, priceRepo *repositories.PriceHistoryRepository) *PortfolioController {
	return &PortfolioController{service: service, priceRepo: priceRepo}
}

// ownsPortfolio returns true if the caller is the owner of the userID resource or is an admin.
func ownsPortfolio(r *http.Request, userID int) bool {
	callerID := middleware.ContextUserID(r.Context())
	return callerID == userID || util.IsAdminRole(middleware.ContextUserRole(r.Context()))
}

func parsePortfolioIDs(w http.ResponseWriter, r *http.Request) (userID, portfolioID int, ok bool) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userID"])
	if err != nil || userID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidUserIDFormat, nil)
		return 0, 0, false
	}
	portfolioID, err = strconv.Atoi(vars["portfolioID"])
	if err != nil || portfolioID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidPortfolioIDFormat, nil)
		return 0, 0, false
	}
	return userID, portfolioID, true
}

func (c *PortfolioController) GetPortfolios(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userID"])
	if err != nil || userID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidUserIDFormat, nil)
		return
	}
	if !ownsPortfolio(r, userID) {
		writeErrorResponse(w, http.StatusForbidden, "Insufficient permissions", nil)
		return
	}

	search := r.URL.Query().Get("search")
	if search == "" {
		portfolios, err := c.service.GetPortfoliosByUserID(userID)
		if err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToRetrievePortfolios, nil)
			return
		}
		writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: portfolios})
		return
	}

	portfolios, err := c.service.SearchAndFilterPortfolios(&repositories.PortfolioFilter{Search: search, UserID: userID})
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToRetrievePortfolios, nil)
		return
	}
	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: portfolios})
}

func (c *PortfolioController) CreatePortfolio(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userID"])
	if err != nil || userID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidUserIDFormat, nil)
		return
	}
	if !ownsPortfolio(r, userID) {
		writeErrorResponse(w, http.StatusForbidden, "Insufficient permissions", nil)
		return
	}

	var payload struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := parseJSONBody(r, &payload); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}
	if payload.Name == "" {
		writeErrorResponse(w, http.StatusBadRequest, "Portfolio name is required", nil)
		return
	}

	portfolio := &db.Portfolio{UserID: userID, Name: payload.Name, Description: payload.Description}
	created, err := c.service.CreatePortfolio(portfolio)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToCreatePortfolio, nil)
		return
	}
	writeJSONResponse(w, http.StatusCreated, SuccessResponse{Data: created})
}

func (c *PortfolioController) GetPortfolio(w http.ResponseWriter, r *http.Request) {
	userID, portfolioID, ok := parsePortfolioIDs(w, r)
	if !ok {
		return
	}
	if !ownsPortfolio(r, userID) {
		writeErrorResponse(w, http.StatusForbidden, "Insufficient permissions", nil)
		return
	}

	portfolio, err := c.service.GetPortfolioByID(userID, portfolioID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgPortfolioNotFound, nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToRetrievePortfolio, nil)
		return
	}
	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: portfolio})
}

func (c *PortfolioController) UpdatePortfolio(w http.ResponseWriter, r *http.Request) {
	userID, portfolioID, ok := parsePortfolioIDs(w, r)
	if !ok {
		return
	}
	if !ownsPortfolio(r, userID) {
		writeErrorResponse(w, http.StatusForbidden, "Insufficient permissions", nil)
		return
	}

	var payload struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := parseJSONBody(r, &payload); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}
	if payload.Name == "" {
		writeErrorResponse(w, http.StatusBadRequest, "Portfolio name is required", nil)
		return
	}

	portfolio := &db.Portfolio{ID: portfolioID, UserID: userID, Name: payload.Name, Description: payload.Description}
	updated, err := c.service.UpdatePortfolio(portfolio)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgPortfolioNotFound, nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToUpdatePortfolio, nil)
		return
	}
	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: updated})
}

func (c *PortfolioController) DeletePortfolio(w http.ResponseWriter, r *http.Request) {
	userID, portfolioID, ok := parsePortfolioIDs(w, r)
	if !ok {
		return
	}
	if !ownsPortfolio(r, userID) {
		writeErrorResponse(w, http.StatusForbidden, "Insufficient permissions", nil)
		return
	}

	if err := c.service.DeletePortfolio(userID, portfolioID); err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgPortfolioNotFound, nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToDeletePortfolio, nil)
		return
	}
	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: map[string]string{"message": "Portfolio deleted successfully"}})
}

func (c *PortfolioController) GetHoldings(w http.ResponseWriter, r *http.Request) {
	userID, portfolioID, ok := parsePortfolioIDs(w, r)
	if !ok {
		return
	}
	if !ownsPortfolio(r, userID) {
		writeErrorResponse(w, http.StatusForbidden, "Insufficient permissions", nil)
		return
	}

	if _, err := c.service.GetPortfolioByID(userID, portfolioID); err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgPortfolioNotFound, nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToRetrievePortfolio, nil)
		return
	}

	holdings, err := c.service.GetHoldings(portfolioID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: []*db.PortfolioItem{}})
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToRetrieveHoldings, nil)
		return
	}
	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: holdings})
}

func (c *PortfolioController) AddHolding(w http.ResponseWriter, r *http.Request) {
	userID, portfolioID, ok := parsePortfolioIDs(w, r)
	if !ok {
		return
	}
	if !ownsPortfolio(r, userID) {
		writeErrorResponse(w, http.StatusForbidden, "Insufficient permissions", nil)
		return
	}

	if _, err := c.service.GetPortfolioByID(userID, portfolioID); err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgPortfolioNotFound, nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToRetrievePortfolio, nil)
		return
	}

	var payload struct {
		StockID  int     `json:"stock_id"`
		Quantity float64 `json:"quantity"`
		AvgCost  float64 `json:"avg_cost"`
	}
	if err := parseJSONBody(r, &payload); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}
	if payload.StockID <= 0 || payload.Quantity <= 0 || payload.AvgCost < 0 {
		writeErrorResponse(w, http.StatusBadRequest, "stock_id, quantity (>0), and avg_cost (>=0) are required", nil)
		return
	}

	item := &db.PortfolioItem{PortfolioID: portfolioID, StockID: payload.StockID, Quantity: payload.Quantity, AvgCost: payload.AvgCost}
	created, err := c.service.CreateHolding(item)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToCreateHolding, nil)
		return
	}
	writeJSONResponse(w, http.StatusCreated, SuccessResponse{Data: created})
}

func (c *PortfolioController) UpdateHolding(w http.ResponseWriter, r *http.Request) {
	userID, portfolioID, ok := parsePortfolioIDs(w, r)
	if !ok {
		return
	}
	if !ownsPortfolio(r, userID) {
		writeErrorResponse(w, http.StatusForbidden, "Insufficient permissions", nil)
		return
	}

	vars := mux.Vars(r)
	holdingID, err := strconv.Atoi(vars["holdingID"])
	if err != nil || holdingID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidHoldingIDFormat, nil)
		return
	}

	if _, err := c.service.GetPortfolioByID(userID, portfolioID); err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgPortfolioNotFound, nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToRetrievePortfolio, nil)
		return
	}

	var payload struct {
		Quantity float64 `json:"quantity"`
		AvgCost  float64 `json:"avg_cost"`
	}
	if err := parseJSONBody(r, &payload); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}
	if payload.Quantity <= 0 || payload.AvgCost < 0 {
		writeErrorResponse(w, http.StatusBadRequest, "quantity must be > 0 and avg_cost >= 0", nil)
		return
	}

	item := &db.PortfolioItem{ID: holdingID, PortfolioID: portfolioID, Quantity: payload.Quantity, AvgCost: payload.AvgCost}
	updated, err := c.service.UpdateHolding(item)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgHoldingNotFound, nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToUpdateHolding, nil)
		return
	}
	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: updated})
}

func (c *PortfolioController) DeleteHolding(w http.ResponseWriter, r *http.Request) {
	userID, portfolioID, ok := parsePortfolioIDs(w, r)
	if !ok {
		return
	}
	if !ownsPortfolio(r, userID) {
		writeErrorResponse(w, http.StatusForbidden, "Insufficient permissions", nil)
		return
	}

	vars := mux.Vars(r)
	holdingID, err := strconv.Atoi(vars["holdingID"])
	if err != nil || holdingID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidHoldingIDFormat, nil)
		return
	}

	if _, err := c.service.GetPortfolioByID(userID, portfolioID); err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgPortfolioNotFound, nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToRetrievePortfolio, nil)
		return
	}

	if err := c.service.DeleteHolding(portfolioID, holdingID); err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgHoldingNotFound, nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToDeleteHolding, nil)
		return
	}
	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: map[string]string{"message": "Holding deleted successfully"}})
}

func (c *PortfolioController) GetPortfolioValue(w http.ResponseWriter, r *http.Request) {
	userID, portfolioID, ok := parsePortfolioIDs(w, r)
	if !ok {
		return
	}
	if !ownsPortfolio(r, userID) {
		writeErrorResponse(w, http.StatusForbidden, "Insufficient permissions", nil)
		return
	}

	if _, err := c.service.GetPortfolioByID(userID, portfolioID); err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgPortfolioNotFound, nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToRetrievePortfolio, nil)
		return
	}

	holdings, err := c.service.GetHoldings(portfolioID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: map[string]float64{"total_value": 0, "cost_basis": 0}})
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToRetrieveHoldings, nil)
		return
	}

	totalValue := 0.0
	totalCostBasis := 0.0
	for _, h := range holdings {
		price := h.AvgCost
		if latest, err := c.priceRepo.GetLatestPrice(h.StockID); err == nil && latest != nil {
			price = latest.Price
		}
		totalValue += h.Quantity * price
		totalCostBasis += h.Quantity * h.AvgCost
	}
	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: map[string]float64{"total_value": totalValue, "cost_basis": totalCostBasis}})
}
