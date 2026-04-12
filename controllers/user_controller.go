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

type UserController struct {
	service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := c.service.GetAllUsers()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToRetrieveUsers, err)
		return
	}

	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: users})
}

func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["id"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidUserIDFormat, err)
		return
	}

	if userID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgUserIDMustBePositive, nil)
		return
	}

	user, err := c.service.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgUserNotFound, nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToRetrieveUser, err)
		return
	}

	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: user})
}

func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := parseJSONBody(r, &payload); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if payload.Name == "" || payload.Email == "" {
		writeErrorResponse(w, http.StatusBadRequest, "Name and email are required", nil)
		return
	}

	user := &db.User{Name: payload.Name, Email: payload.Email}
	createdUser, err := c.service.CreateUser(user)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToCreateUser, err)
		return
	}

	writeJSONResponse(w, http.StatusCreated, SuccessResponse{Data: createdUser})
}

func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["id"]

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
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := parseJSONBody(r, &payload); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if payload.Name == "" || payload.Email == "" {
		writeErrorResponse(w, http.StatusBadRequest, "Name and email are required", nil)
		return
	}

	user := &db.User{ID: userID, Name: payload.Name, Email: payload.Email}
	updatedUser, err := c.service.UpdateUser(userID, user)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgUserNotFound, nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToUpdateUser, err)
		return
	}

	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: updatedUser})
}

func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["id"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidUserIDFormat, err)
		return
	}

	if userID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgUserIDMustBePositive, nil)
		return
	}

	if err := c.service.DeleteUser(userID); err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgUserNotFound, nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToDeleteUser, err)
		return
	}

	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: map[string]string{"message": "User deleted successfully"}})
}
