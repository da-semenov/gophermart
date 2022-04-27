package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/da-semenov/gophermart/internal/app/domain"
	"github.com/da-semenov/gophermart/internal/app/infrastructure"
	"go.uber.org/zap"
	"net/http"
)

type BalanceService interface {
	GetCurrentBalance(ctx context.Context, userID int) (*domain.Balance, error)
	Withdraw(ctx context.Context, obj *domain.Withdraw, userID int) error
	GetWithdrawalsList(ctx context.Context, userID int) ([]domain.Withdrawal, error)
}

type BalanceHandler struct {
	balanceService BalanceService
	auth           *Auth
	log            *infrastructure.Logger
}

func NewBalanceHandler(bs BalanceService, auth *Auth, l *infrastructure.Logger) *BalanceHandler {
	var target BalanceHandler
	target.log = l
	target.balanceService = bs
	target.auth = auth
	return &target
}

func (h *BalanceHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _, err := h.auth.GetFromContext(ctx)
	if err != nil {
		h.log.Error("BalanceHandler:can't get params from the token", zap.Error(err))
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("внутренняя ошибка сервера")); err != nil {
			h.log.Error("AuthHandler: can't write response", zap.Error(err))
		}
		return
	}
	balance, err := h.balanceService.GetCurrentBalance(ctx, userID)
	if err != nil {
		h.log.Error("BalanceHandler:can't get params from the token", zap.Error(err))
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("внутренняя ошибка сервера")); err != nil {
			h.log.Error("AuthHandler: can't write response", zap.Error(err))
		}
		return
	}
	responseBody, err := json.Marshal(balance)
	if err != nil {
		h.log.Error("BalanceHandler: can't serialize response", zap.Error(err))
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("внутренняя ошибка сервера")); err != nil {
			h.log.Error("BalanceHandler: can't write response", zap.Error(err))
		}
		return
	}
	if err = WriteResponse(w, http.StatusOK, responseBody); err != nil {
		h.log.Error("BalanceHandler: can't write response", zap.Error(err))
	}
}

func (h *BalanceHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	b, err := getRequestBody(r)
	if err != nil {
		h.log.Error("BalanceHandler:can't get request body", zap.Error(err))
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("внутренняя ошибка сервера")); err != nil {
			h.log.Error("BalanceHandler: can't write response", zap.Error(err))
		}
		return
	}
	if len(b) == 0 || r.Header.Get("Content-Type") != "application/json" {
		h.log.Info("BalanceHandler:empty request body")
		if err = WriteResponse(w, http.StatusBadRequest, ErrMessage("неверный формат запроса")); err != nil {
			h.log.Error("BalanceHandler: can't write response", zap.Error(err))
		}
		return
	}

	var (
		withdraw domain.Withdraw
	)
	err = json.Unmarshal(b, &withdraw)
	if err != nil {
		h.log.Error("BalanceHandler:can't unmarshal request body", zap.Error(err))
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("внутренняя ошибка сервера")); err != nil {
			h.log.Error("BalanceHandler: can't write response", zap.Error(err))
		}
		return
	}
	ctx := r.Context()
	userID, _, err := h.auth.GetFromContext(ctx)
	h.log.Info("Try to withdraw funds", zap.String("OrderNum", withdraw.OrderNum), zap.Int("userID", userID))
	if err != nil {
		h.log.Error("BalanceHandler:can't get params from the token", zap.Error(err))
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("внутренняя ошибка сервера")); err != nil {
			h.log.Error("BalanceHandler: can't write response", zap.Error(err))
		}
		return
	}
	err = h.balanceService.Withdraw(ctx, &withdraw, userID)
	if err != nil {
		h.log.Error("BalanceHandler:Withdraw error", zap.Error(err))
		if errors.Is(err, domain.ErrNotEnoughFunds) {
			if err = WriteResponse(w, http.StatusPaymentRequired, ErrMessage("на счету недостаточно средств")); err != nil {
				h.log.Error("BalanceHandler: can't write response", zap.Error(err))
			}
		} else if errors.Is(err, domain.ErrBadOrderNum) {
			if err = WriteResponse(w, http.StatusUnprocessableEntity, ErrMessage("неверный номер заказа")); err != nil {
				h.log.Error("BalanceHandler: can't write response", zap.Error(err))
			}
		} else {
			if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("внутренняя ошибка сервера")); err != nil {
				h.log.Error("BalanceHandler: can't write response", zap.Error(err))
			}
		}
		return
	}
	if err = WriteResponse(w, http.StatusOK, ""); err != nil {
		h.log.Error("BalanceHandler: can't write response", zap.Error(err))
	}
	h.log.Info("Withdraw success", zap.String("OrderNum", withdraw.OrderNum), zap.Int("userID", userID))
}

func (h *BalanceHandler) GetWithdrawalsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _, err := h.auth.GetFromContext(ctx)
	if err != nil {
		h.log.Error("BalanceHandler:can't get params from the token", zap.Error(err))
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("внутренняя ошибка сервера")); err != nil {
			h.log.Error("AuthHandler: can't write response", zap.Error(err))
		}
		return
	}
	res, err := h.balanceService.GetWithdrawalsList(ctx, userID)
	if err != nil {
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("внутренняя ошибка сервера")); err != nil {
			h.log.Error("BalanceHandler: can't write response", zap.Error(err))
		}
	} else {
		if err = WriteResponse(w, http.StatusOK, res); err != nil {
			h.log.Error("BalanceHandler: can't write response", zap.Error(err))
		}
	}
}