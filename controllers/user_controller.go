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

type UserController struct {
	service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{service: service}
}

// Login authenticates a user and returns a JWT token.
// POST /auth/login
func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := parseJSONBody(r, &payload); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}
	if payload.Email == "" || payload.Password == "" {
		writeErrorResponse(w, http.StatusBadRequest, "Email and password are required", nil)
		return
	}

	user, err := c.service.GetUserByEmail(payload.Email)
	if err != nil || user == nil {
		// Constant-time response to prevent user enumeration
		writeErrorResponse(w, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	if !util.VerifyPassword(user.Password, payload.Password) {
		writeErrorResponse(w, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	token, err := util.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to generate token", nil)
		return
	}

	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: map[string]string{"token": token}})
}

func (c *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Only admins may list all users
	if !util.IsAdminRole(middleware.ContextUserRole(r.Context())) {
		writeErrorResponse(w, http.StatusForbidden, "Insufficient permissions", nil)
		return
	}

	search := r.URL.Query().Get("search")
	if search == "" {
		users, err := c.service.GetAllUsers()
		if err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToRetrieveUsers, nil)
			return
		}
		writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: users})
		return
	}

	users, err := c.service.SearchAndFilterUsers(&repositories.UserFilter{Search: search})
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToRetrieveUsers, nil)
		return
	}
	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: users})
}

func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := parseUserID(w, r)
	if !ok {
		return
	}

	// Users may only fetch their own profile; admins may fetch any
	callerID := middleware.ContextUserID(r.Context())
	if callerID != userID && !util.IsAdminRole(middleware.ContextUserRole(r.Context())) {
		writeErrorResponse(w, http.StatusForbidden, "Insufficient permissions", nil)
		return
	}

	user, err := c.service.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgUserNotFound, nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToRetrieveUser, nil)
		return
	}
	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: user})
}

func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := parseJSONBody(r, &payload); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}
	if payload.Name == "" || payload.Email == "" || payload.Password == "" {
		writeErrorResponse(w, http.StatusBadRequest, "Name, email, and password are required", nil)
		return
	}
	if !util.ValidateEmail(payload.Email) {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid email format", nil)
		return
	}

	hashedPassword, err := util.HashPassword(payload.Password)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid password: "+err.Error(), nil)
		return
	}

	// Role is always "user" on self-registration; admins must be promoted separately
	user := &db.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: hashedPassword,
		Role:     "user",
		IsActive: true,
	}
	createdUser, err := c.service.CreateUser(user)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToCreateUser, nil)
		return
	}
	writeJSONResponse(w, http.StatusCreated, SuccessResponse{Data: createdUser})
}

func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := parseUserID(w, r)
	if !ok {
		return
	}

	// Users may only update their own profile; admins may update any
	callerID := middleware.ContextUserID(r.Context())
	if callerID != userID && !util.IsAdminRole(middleware.ContextUserRole(r.Context())) {
		writeErrorResponse(w, http.StatusForbidden, "Insufficient permissions", nil)
		return
	}

	var payload struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := parseJSONBody(r, &payload); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}
	if payload.Name == "" || payload.Email == "" {
		writeErrorResponse(w, http.StatusBadRequest, "Name and email are required", nil)
		return
	}
	if !util.ValidateEmail(payload.Email) {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid email format", nil)
		return
	}

	user := &db.User{ID: userID, Name: payload.Name, Email: payload.Email}
	updatedUser, err := c.service.UpdateUser(userID, user)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgUserNotFound, nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToUpdateUser, nil)
		return
	}
	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: updatedUser})
}

func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := parseUserID(w, r)
	if !ok {
		return
	}

	// Users may only delete their own account; admins may delete any
	callerID := middleware.ContextUserID(r.Context())
	if callerID != userID && !util.IsAdminRole(middleware.ContextUserRole(r.Context())) {
		writeErrorResponse(w, http.StatusForbidden, "Insufficient permissions", nil)
		return
	}

	if err := c.service.DeleteUser(userID); err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			writeErrorResponse(w, http.StatusNotFound, constants.ErrMsgUserNotFound, nil)
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, constants.ErrMsgFailedToDeleteUser, nil)
		return
	}
	writeJSONResponse(w, http.StatusOK, SuccessResponse{Data: map[string]string{"message": "User deleted successfully"}})
}

// parseUserID extracts and validates the {id} path variable.
func parseUserID(w http.ResponseWriter, r *http.Request) (int, bool) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgInvalidUserIDFormat, nil)
		return 0, false
	}
	if userID <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, constants.ErrMsgUserIDMustBePositive, nil)
		return 0, false
	}
	return userID, true
}
