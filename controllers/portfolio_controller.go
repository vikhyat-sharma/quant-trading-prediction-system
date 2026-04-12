package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/constants"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/services"
)

type PortfolioController struct {
	service *services.PortfolioService
}

func NewPortfolioController(service *services.PortfolioService) *PortfolioController {
	return &PortfolioController{service: service}
}

func (c *PortfolioController) GetPortfolios(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidUserIDFormat, err)
		return
	}

	if userID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgUserIDMustBePositive, nil)
		return
	}

	portfolios, err := c.service.GetPortfoliosByUserID(userID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgPortfolioNotFound, nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToRetrievePortfolios, err)
		return
	}

	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: portfolios})
}

func (c *PortfolioController) CreatePortfolio(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidUserIDFormat, err)
		return
	}

	if userID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgUserIDMustBePositive, nil)
		return
	}

	var payload struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := parseJSONBody(r, &payload); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if payload.Name == "" {
		writeErrorResponse(w, http.StatusBadRequest, "Portfolio name is required", nil)
		return
	}

	portfolio := &db.Portfolio{UserID: userID, Name: payload.Name, Description: payload.Description}
	createdPortfolio, err := c.service.CreatePortfolio(portfolio)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToCreatePortfolio, err)
		return
	}

	writeJSONResponse(w, http.StatusCreated, SuccessResponse{Data: createdPortfolio})
}

func (c *PortfolioController) GetPortfolio(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	portfolioIDStr := vars["portfolioID"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidUserIDFormat, err)
		return
	}

	portfolioID, err := strconv.Atoi(portfolioIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidPortfolioIDFormat, err)
		return
	}

	if userID <= 0 || portfolioID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgUserIDMustBePositive, nil)
		return
	}

	portfolio, err := c.service.GetPortfolioByID(userID, portfolioID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgPortfolioNotFound, nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToRetrievePortfolio, err)
		return
	}

	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: portfolio})
}

func (c *PortfolioController) UpdatePortfolio(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	portfolioIDStr := vars["portfolioID"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidUserIDFormat, err)
		return
	}

	portfolioID, err := strconv.Atoi(portfolioIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidPortfolioIDFormat, err)
		return
	}

	if userID <= 0 || portfolioID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgUserIDMustBePositive, nil)
		return
	}

	var payload struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := parseJSONBody(r, &payload); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if payload.Name == "" {
		writeErrorResponse(w, http.StatusBadRequest, "Portfolio name is required", nil)
		return
	}

	portfolio := &db.Portfolio{ID: portfolioID, UserID: userID, Name: payload.Name, Description: payload.Description}
	updatedPortfolio, err := c.service.UpdatePortfolio(portfolio)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgPortfolioNotFound, nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToUpdatePortfolio, err)
		return
	}

	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: updatedPortfolio})
}

func (c *PortfolioController) DeletePortfolio(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	portfolioIDStr := vars["portfolioID"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidUserIDFormat, err)
		return
	}

	portfolioID, err := strconv.Atoi(portfolioIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidPortfolioIDFormat, err)
		return
	}

	if userID <= 0 || portfolioID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgUserIDMustBePositive, nil)
		return
	}

	if err := c.service.DeletePortfolio(userID, portfolioID); err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgPortfolioNotFound, nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToDeletePortfolio, err)
		return
	}

	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: map[string]string{"message": "Portfolio deleted successfully"}})
}

func (c *PortfolioController) GetHoldings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	portfolioIDStr := vars["portfolioID"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidUserIDFormat, err)
		return
	}

	portfolioID, err := strconv.Atoi(portfolioIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidPortfolioIDFormat, err)
		return
	}

	if userID <= 0 || portfolioID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgUserIDMustBePositive, nil)
		return
	}

	if _, err := c.service.GetPortfolioByID(userID, portfolioID); err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgPortfolioNotFound, nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToRetrievePortfolio, err)
		return
	}

	holdings, err := c.service.GetHoldings(portfolioID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: []*db.PortfolioItem{}})
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToRetrieveHoldings, err)
		return
	}

	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: holdings})
}

func (c *PortfolioController) AddHolding(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	portfolioIDStr := vars["portfolioID"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidUserIDFormat, err)
		return
	}

	portfolioID, err := strconv.Atoi(portfolioIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidPortfolioIDFormat, err)
		return
	}

	if userID <= 0 || portfolioID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgUserIDMustBePositive, nil)
		return
	}

	if _, err := c.service.GetPortfolioByID(userID, portfolioID); err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgPortfolioNotFound, nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToRetrievePortfolio, err)
		return
	}

	var payload struct {
		StockID  int     `json:"stock_id"`
		Quantity float64 `json:"quantity"`
		AvgCost  float64 `json:"avg_cost"`
	}

	if err := parseJSONBody(r, &payload); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if payload.StockID <= 0 || payload.Quantity <= 0 || payload.AvgCost < 0 {
		writeErrorResponse(w, http.StatusBadRequest, "Stock ID, quantity and avg_cost are required and must be valid", nil)
		return
	}

	item := &db.PortfolioItem{
		PortfolioID: portfolioID,
		StockID:     payload.StockID,
		Quantity:    payload.Quantity,
		AvgCost:     payload.AvgCost,
	}

	createdHolding, err := c.service.CreateHolding(item)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToCreateHolding, err)
		return
	}

	writeJSONResponse(w, http.StatusCreated, SuccessResponse{Data: createdHolding})
}

func (c *PortfolioController) UpdateHolding(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	portfolioIDStr := vars["portfolioID"]
	holdingIDStr := vars["holdingID"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidUserIDFormat, err)
		return
	}

	portfolioID, err := strconv.Atoi(portfolioIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidPortfolioIDFormat, err)
		return
	}

	holdingID, err := strconv.Atoi(holdingIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidHoldingIDFormat, err)
		return
	}

	if userID <= 0 || portfolioID <= 0 || holdingID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgUserIDMustBePositive, nil)
		return
	}

	if _, err := c.service.GetPortfolioByID(userID, portfolioID); err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgPortfolioNotFound, nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToRetrievePortfolio, err)
		return
	}

	var payload struct {
		Quantity float64 `json:"quantity"`
		AvgCost  float64 `json:"avg_cost"`
	}

	if err := parseJSONBody(r, &payload); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if payload.Quantity <= 0 || payload.AvgCost < 0 {
		writeErrorResponse(w, http.StatusBadRequest, "Quantity must be greater than zero and avg_cost must be non-negative", nil)
		return
	}

	item := &db.PortfolioItem{ID: holdingID, PortfolioID: portfolioID, Quantity: payload.Quantity, AvgCost: payload.AvgCost}
	updatedHolding, err := c.service.UpdateHolding(item)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgHoldingNotFound, nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToUpdateHolding, err)
		return
	}

	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: updatedHolding})
}

func (c *PortfolioController) DeleteHolding(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	portfolioIDStr := vars["portfolioID"]
	holdingIDStr := vars["holdingID"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidUserIDFormat, err)
		return
	}

	portfolioID, err := strconv.Atoi(portfolioIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidPortfolioIDFormat, err)
		return
	}

	holdingID, err := strconv.Atoi(holdingIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidHoldingIDFormat, err)
		return
	}

	if userID <= 0 || portfolioID <= 0 || holdingID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgUserIDMustBePositive, nil)
		return
	}

	if _, err := c.service.GetPortfolioByID(userID, portfolioID); err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgPortfolioNotFound, nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToRetrievePortfolio, err)
		return
	}

	if err := c.service.DeleteHolding(portfolioID, holdingID); err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgHoldingNotFound, nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToDeleteHolding, err)
		return
	}

	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: map[string]string{"message": "Holding deleted successfully"}})
}
