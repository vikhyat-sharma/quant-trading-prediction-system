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

type UserAlertRuleController struct {
	service *services.UserAlertRuleService
}

func NewUserAlertRuleController(service *services.UserAlertRuleService) *UserAlertRuleController {
	return &UserAlertRuleController{service: service}
}

func (c *UserAlertRuleController) CreateAlertRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userID"])
	if err != nil || userID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidUserIDFormat, nil)
		return
	}
	// Enforce ownership: caller can only create rules for themselves
	callerID := middleware.ContextUserID(r.Context())
	if callerID != userID && !util.IsAdminRole(middleware.ContextUserRole(r.Context())) {
		writeErrorResponse(w, http.StatusForbidden, "Insufficient permissions", nil)
		return
	}

	var payload struct {
		StockID   int     `json:"stock_id"`
		Threshold float64 `json:"threshold"`
		Condition string  `json:"condition"`
	}
	if err := parseJSONBody(r, &payload); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}
	if payload.StockID <= 0 || payload.Condition == "" {
		writeErrorResponse(w, http.StatusBadRequest, "stock_id and condition are required", nil)
		return
	}
	if payload.Condition != "above" && payload.Condition != "below" {
		writeErrorResponse(w, http.StatusBadRequest, "condition must be 'above' or 'below'", nil)
		return
	}

	rule := &db.UserAlertRule{
		UserID:    userID,
		StockID:   payload.StockID,
		Threshold: payload.Threshold,
		Condition: payload.Condition,
		Enabled:   true,
	}
	created, err := c.service.CreateAlertRule(rule)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to create alert rule", nil)
		return
	}
	writeJSONResponse(w, http.StatusCreated, SuccessResponse{Data: created})
}

func (c *UserAlertRuleController) GetAlertRules(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userID"])
	if err != nil || userID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidUserIDFormat, nil)
		return
	}
	// Enforce ownership
	callerID := middleware.ContextUserID(r.Context())
	if callerID != userID && !util.IsAdminRole(middleware.ContextUserRole(r.Context())) {
		writeErrorResponse(w, http.StatusForbidden, "Insufficient permissions", nil)
		return
	}

	rules, err := c.service.GetAlertRulesByUser(userID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to get alert rules", nil)
		return
	}
	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: rules})
}

func (c *UserAlertRuleController) DeleteAlertRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userID"])
	if err != nil || userID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidUserIDFormat, nil)
		return
	}
	// Enforce ownership
	callerID := middleware.ContextUserID(r.Context())
	if callerID != userID && !util.IsAdminRole(middleware.ContextUserRole(r.Context())) {
		writeErrorResponse(w, http.StatusForbidden, "Insufficient permissions", nil)
		return
	}

	ruleID, err := strconv.Atoi(vars["ruleID"])
	if err != nil || ruleID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid alert rule ID", nil)
		return
	}

	if err := c.service.DeleteAlertRule(ruleID); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to delete alert rule", nil)
		return
	}
	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: map[string]string{"message": "Alert rule deleted"}})
}
