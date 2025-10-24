package handlers

import (
	"encoding/json"
	"net/http"

	"go-starter/internal/logger"
	"go-starter/pkg/database"

	"go.uber.org/zap"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db *database.DB
}

// NewHealthHandler creates a new health check handler
func NewHealthHandler(db *database.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

// Healthz godoc
// @Summary Health check
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 503 {object} HealthResponse
// @Router /healthz [get]
func (h *HealthHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:   "ok",
		Database: "ok",
	}

	statusCode := http.StatusOK

	// Check database health
	if err := h.db.Health(r.Context()); err != nil {
		logger.FromContext(r.Context()).Error("database health check failed", zap.Error(err))
		response.Database = "unhealthy"
		response.Status = "unhealthy"
		statusCode = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// Ready godoc
// @Summary Readiness check
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 503 {object} HealthResponse
// @Router /ready [get]
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	// Same as healthz for now, but can be extended for more complex readiness checks
	h.Healthz(w, r)
}
