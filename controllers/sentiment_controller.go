package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/constants"
	"github.com/vikhyat-sharma/quant-trading-prediction-system/services"
)

type SentimentController struct {
	service *services.SentimentService
}

type SentimentRequest struct {
	Text string `json:"text"`
}

func NewSentimentController(service *services.SentimentService) *SentimentController {
	return &SentimentController{service: service}
}

func (c *SentimentController) AnalyzeSentiment(w http.ResponseWriter, r *http.Request) {
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

	var request SentimentRequest
	if err := parseJSONBody(r, &request); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid sentiment request body", err)
		return
	}

	if request.Text == "" {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgSentimentTextRequired, nil)
		return
	}

	result := c.service.Analyze(request.Text)
	if result == nil {
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToAnalyzeSentiment, errors.New("sentiment analysis failed"))
		return
	}

	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: result})
}
