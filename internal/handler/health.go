package handler

import (
	"database/sql"
	"net/http"

	"github.com/davidgrcias/digital-wallet/pkg/response"
)

type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	// Actually ping the database
	if err := h.db.Ping(); err != nil {
		response.JSON(w, http.StatusServiceUnavailable, response.APIResponse{
			Success: false,
			Message: "Database connection failed",
		})
		return
	}

	response.JSON(w, http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Service is healthy",
		Data: map[string]string{
			"status":   "healthy",
			"database": "connected",
		},
	})
}
