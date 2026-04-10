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

type AlertController struct {
	service *services.AlertService
}

func NewAlertController(service *services.AlertService) *AlertController {
	return &AlertController{service: service}
}

func (c *AlertController) CreateAlert(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stockIDStr := vars["stockID"]

	stockID, err := strconv.Atoi(stockIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidStockIDFormat, err)
		return
	}

	if stockID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgStockIDMustBePositive, nil)
		return
	}

	var payload struct {
		Threshold float64 `json:"threshold"`
		Condition string  `json:"condition"`
	}

	if err := parseJSONBody(r, &payload); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	alert := &db.Alert{
		StockID:   stockID,
		Threshold: payload.Threshold,
		Condition: payload.Condition,
		Enabled:   true,
	}

	createdAlert, err := c.service.CreateAlert(alert)
	if err != nil {
		if errors.Is(err, db.ErrInvalidThreshold) || errors.Is(err, db.ErrInvalidCondition) {
			writeErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to create alert", err)
		return
	}

	writeJSONResponse(w, http.StatusCreated, SuccessResponse{Data: createdAlert})
}

func (c *AlertController) GetAlerts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stockIDStr := vars["stockID"]

	stockID, err := strconv.Atoi(stockIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidStockIDFormat, err)
		return
	}

	if stockID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgStockIDMustBePositive, nil)
		return
	}

	alerts, err := c.service.GetAlerts(stockID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve alerts", err)
		return
	}

	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: alerts})
}

func (c *AlertController) DeleteAlert(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stockIDStr := vars["stockID"]
	alertIDStr := vars["alertID"]

	stockID, err := strconv.Atoi(stockIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidStockIDFormat, err)
		return
	}

	alertID, err := strconv.Atoi(alertIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidPredictionIDFormat, err)
		return
	}

	if stockID <= 0 || alertID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgStockIDMustBePositive, nil)
		return
	}

	if err := c.service.DeleteAlert(stockID, alertID); err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, "Alert not found", nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to delete alert", err)
		return
	}

	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: map[string]string{"message": "Alert deleted"}})
}

func (c *AlertController) EvaluateAlerts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stockIDStr := vars["stockID"]

	stockID, err := strconv.Atoi(stockIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidStockIDFormat, err)
		return
	}

	if stockID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgStockIDMustBePositive, nil)
		return
	}

	notifications, err := c.service.EvaluateAlerts(stockID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, "No price history available for this stock", nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to evaluate alerts", err)
		return
	}

	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: notifications})
}

func (c *AlertController) GetNotifications(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stockIDStr := vars["stockID"]

	stockID, err := strconv.Atoi(stockIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidStockIDFormat, err)
		return
	}

	if stockID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgStockIDMustBePositive, nil)
		return
	}

	notifications, err := c.service.GetNotifications(stockID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, "No notifications found for this stock", nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve notifications", err)
		return
	}

	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: notifications})
}
