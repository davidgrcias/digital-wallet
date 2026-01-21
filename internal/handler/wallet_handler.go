package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/davidgarcia/digital-wallet/internal/domain"
	"github.com/davidgarcia/digital-wallet/internal/usecase"
	"github.com/davidgarcia/digital-wallet/pkg/response"
)

type WalletHandler struct {
	walletUsecase usecase.WalletUsecase
}

func NewWalletHandler(walletUsecase usecase.WalletUsecase) *WalletHandler {
	return &WalletHandler{
		walletUsecase: walletUsecase,
	}
}

func (h *WalletHandler) RegisterRoutes(router *mux.Router) {
	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/users/{user_id}/balance", h.GetBalance).Methods(http.MethodGet)
	api.HandleFunc("/users/{user_id}/withdraw", h.Withdraw).Methods(http.MethodPost)
}

func (h *WalletHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	if userID == "" {
		response.BadRequest(w, "user_id is required", nil)
		return
	}

	result, err := h.walletUsecase.GetBalance(r.Context(), userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.Success(w, "Balance retrieved successfully", result)
}

func (h *WalletHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	if userID == "" {
		response.BadRequest(w, "user_id is required", nil)
		return
	}

	var req domain.WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", nil)
		return
	}

	result, err := h.walletUsecase.Withdraw(r.Context(), userID, &req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.Success(w, "Withdrawal successful", result)
}

func (h *WalletHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		response.NotFound(w, "User not found")
	case errors.Is(err, domain.ErrInsufficientBalance):
		response.BadRequest(w, "Insufficient balance", nil)
	case errors.Is(err, domain.ErrInvalidAmount):
		response.BadRequest(w, "Amount must be greater than 0", nil)
	default:
		response.InternalError(w, "Internal server error")
	}
}
