package handlers

import (
	"encoding/json"
	"net/http"

	"go-starter/internal/logger"
	"go-starter/internal/models"
	"go-starter/internal/services"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService *services.AuthService
	validate    *validator.Validate
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(),
	}
}

// Register godoc
// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration credentials"
// @Success 201 {object} models.AuthResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, r, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		respondWithError(w, r, http.StatusBadRequest, "validation failed", err)
		return
	}

	// Register user
	response, err := h.authService.Register(r.Context(), &req)
	if err != nil {
		if err == services.ErrUserExists {
			respondWithError(w, r, http.StatusConflict, "user already exists", err)
		} else {
			respondWithError(w, r, http.StatusInternalServerError, "failed to register user", err)
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, response)
}

// Login godoc
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, r, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		respondWithError(w, r, http.StatusBadRequest, "validation failed", err)
		return
	}

	// Login user
	response, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			respondWithError(w, r, http.StatusUnauthorized, "invalid credentials", err)
		} else {
			respondWithError(w, r, http.StatusInternalServerError, "failed to login", err)
		}
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

// respondWithJSON sends a JSON response
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		// Log encoding error but don't send another response
		logger.Error("failed to encode JSON response", zap.Error(err))
	}
}

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, r *http.Request, code int, message string, err error) {
	logger.FromContext(r.Context()).Error(message,
		zap.Error(err),
		zap.Int("status_code", code),
	)

	response := models.ErrorResponse{
		Error:   message,
		Message: err.Error(),
	}

	respondWithJSON(w, code, response)
}
