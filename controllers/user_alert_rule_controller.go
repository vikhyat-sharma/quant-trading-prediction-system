package controllers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/constants"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/db"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/services"
)

type UserAlertRuleController struct {
	service *services.UserAlertRuleService
}

func NewUserAlertRuleController(service *services.UserAlertRuleService) *UserAlertRuleController {
	return &UserAlertRuleController{service: service}
}

func (c *UserAlertRuleController) CreateAlertRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidUserIDFormat, err)
		return
	}
	var payload struct {
		StockID   int     `json:"stock_id"`
		Threshold float64 `json:"threshold"`
		Condition string  `json:"condition"`
	}
	if err := parseJSONBody(r, &payload); err != nil || payload.StockID <= 0 || payload.Condition == "" {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	rule := &db.UserAlertRule{UserID: userID, StockID: payload.StockID, Threshold: payload.Threshold, Condition: payload.Condition, Enabled: true}
	created, err := c.service.CreateAlertRule(rule)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to create alert rule", err)
		return
	}
	writeJSONResponse(w, http.StatusCreated, SuccessResponse{Data: created})
}

func (c *UserAlertRuleController) GetAlertRules(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidUserIDFormat, err)
		return
	}
	rules, err := c.service.GetAlertRulesByUser(userID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to get alert rules", err)
		return
	}
	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: rules})
}

func (c *UserAlertRuleController) DeleteAlertRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ruleIDStr := vars["ruleID"]
	ruleID, err := strconv.Atoi(ruleIDStr)
	if err != nil || ruleID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid alert rule ID", err)
		return
	}
	if err := c.service.DeleteAlertRule(ruleID); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to delete alert rule", err)
		return
	}
	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: map[string]string{"message": "Alert rule deleted"}})
}
