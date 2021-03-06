package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/da-semenov/gophermart/internal/app/domain"
	"github.com/da-semenov/gophermart/internal/app/infrastructure"
	"go.uber.org/zap"
	"net/http"
)

type OrderService interface {
	Save(ctx context.Context, order *domain.Order) error
	GetOrderList(ctx context.Context, userID int) ([]domain.Order, error)
}

type OrderHandler struct {
	orderService OrderService
	auth         *Auth
	log          *infrastructure.Logger
}

func NewOrderHandler(os OrderService, auth *Auth, l *infrastructure.Logger) *OrderHandler {
	var target OrderHandler
	target.log = l
	target.orderService = os
	target.auth = auth
	return &target
}

func (h *OrderHandler) RegisterNewOrder(w http.ResponseWriter, r *http.Request) {
	b, err := getRequestBody(r)
	if err != nil {
		h.log.Error("OrderHandler:can't get request body", zap.Error(err))
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("внутренняя ошибка сервера")); err != nil {
			h.log.Error("OrderHandler: can't write response", zap.Error(err))
		}
		return
	}

	if len(b) == 0 || r.Header.Get("Content-Type") != "text/plain" {
		h.log.Info("OrderHandler:empty request body")
		if err = WriteResponse(w, http.StatusBadRequest, ErrMessage("неверный формат запроса")); err != nil {
			h.log.Error("OrderHandler: can't write response", zap.Error(err))
		}
		return
	}
	var (
		order domain.Order
	)

	order.Num = string(b)
	ctx := r.Context()
	u, _, err := h.auth.GetFromContext(ctx)
	if err != nil {
		h.log.Error("OrderHandler:can't get params from the token", zap.Error(err))
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("внутренняя ошибка сервера")); err != nil {
			h.log.Error("OrderHandler: can't write response", zap.Error(err))
		}
		return
	}
	order.UserID = u
	err = h.orderService.Save(ctx, &order)
	var (
		statusCode int
		msg        string
	)
	if err != nil {
		h.log.Error("OrderHandler:received an error", zap.Error(err))
		switch err {
		case domain.ErrOrderRegistered:
			statusCode = http.StatusOK
			msg = "номер заказа уже был загружен этим пользователем"
		case domain.ErrOrderRegisteredByAnotherUser:
			statusCode = http.StatusConflict
			msg = "номер заказа уже был загружен другим пользователем"
		case domain.ErrBadParam:
			statusCode = http.StatusBadRequest
			msg = "неверный формат запроса"
		case domain.ErrBadOrderNum:
			statusCode = http.StatusUnprocessableEntity
			msg = "неверный формат номера заказа"
		default:
			statusCode = http.StatusInternalServerError
			msg = "внутренняя ошибка сервера"
		}
		if err = WriteResponse(w, statusCode, ErrMessage(msg)); err != nil {
			h.log.Error("OrderHandler: can't write response", zap.Error(err))
		}
		return
	}

	if err = WriteResponse(w, http.StatusAccepted, nil); err != nil {
		h.log.Error("OrderHandler: can't write response", zap.Error(err))
	}
	h.log.Info(fmt.Sprintf("Order %s successfully registered", order.Num))
}

func (h *OrderHandler) GetOrderList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _, err := h.auth.GetFromContext(ctx)
	if err != nil {
		h.log.Error("OrderHandler:can't get params from the token", zap.Error(err))
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("внутренняя ошибка сервера")); err != nil {
			h.log.Error("OrderHandler: can't write response", zap.Error(err))
		}
		return
	}

	res, err := h.orderService.GetOrderList(ctx, userID)
	if err != nil {
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("внутренняя ошибка сервера")); err != nil {
			h.log.Error("OrderHandler: can't write response", zap.Error(err))
		}
	} else if len(res) == 0 {
		if err = WriteResponse(w, http.StatusNoContent, ErrMessage("нет данных для ответа")); err != nil {
			h.log.Error("OrderHandler: can't write response", zap.Error(err))
		}
	} else {
		responseBody, err := json.Marshal(res)
		h.log.Debug("All orders:", zap.String("domain", string(responseBody)))
		if err != nil {
			h.log.Error("OrderHandler: can't serialize response", zap.Error(err))
			return
		}
		if err = WriteResponse(w, http.StatusOK, responseBody); err != nil {
			h.log.Error("OrderHandler: can't write response", zap.Error(err))
		}
	}
}
