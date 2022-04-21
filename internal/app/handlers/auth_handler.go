package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/da-semenov/gophermart/internal/app/domain/model"
	"github.com/da-semenov/gophermart/internal/app/infrastructure"
	"go.uber.org/zap"
	"net/http"
)

type AuthService interface {
	Register(ctx context.Context, user *model.User) (*model.User, error)
	Check(ctx context.Context, user *model.User) (*model.User, error)
}

type AuthHandler struct {
	authService AuthService
	auth        *Auth
	log         *infrastructure.Logger
}

func NewAuthHandler(as AuthService, auth *Auth, l *infrastructure.Logger) *AuthHandler {
	var target AuthHandler
	target.log = l
	target.authService = as
	target.auth = auth
	return &target
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user model.User
	b, err := getRequestBody(r)
	if err != nil {
		h.log.Error("AuthHandler:can't get request body", zap.Error(err))
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("внутренняя ошибка сервера")); err != nil {
			h.log.Error("AuthHandler: can't write response", zap.Error(err))
		}
		return
	}
	if err := json.Unmarshal(b, &user); err != nil {
		h.log.Error("AuthHandler:can't unmarshal body", zap.Error(err))
		if err = WriteResponse(w, http.StatusBadRequest, ErrMessage("неверный формат запроса")); err != nil {
			h.log.Error("AuthHandler: can't write response", zap.Error(err))
		}
		return
	}
	ctx := r.Context()
	u, err := h.authService.Register(ctx, &user)
	if err != nil {
		h.log.Error("AuthHandler:received an error", zap.Error(err))
		if errors.Is(err, model.ErrDuplicateKey) {
			if err = WriteResponse(w, http.StatusConflict, ErrMessage("логин уже занят")); err != nil {
				h.log.Error("AuthHandler: can't write response", zap.Error(err))
			}
			return
		} else if errors.Is(err, model.ErrBadParam) {
			if err = WriteResponse(w, http.StatusBadRequest, ErrMessage("неверный формат запроса")); err != nil {
				h.log.Error("AuthHandler: can't write response", zap.Error(err))
			}
			return
		} else {
			if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("внутренняя ошибка сервера")); err != nil {
				h.log.Error("AuthHandler: can't write response", zap.Error(err))
			}
			return
		}
	}
	token, err := h.auth.GetNewToken(u.ID, u.Login)
	if err != nil {
		h.log.Error("AuthHandler: can't make token", zap.Error(err))
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("внутренняя ошибка сервера")); err != nil {
			h.log.Error("AuthHandler: can't write response", zap.Error(err))
		}
		return
	}
	jwtCookie, err := bakeCookie(token)
	if err != nil {
		h.log.Error("AuthHandler: can't set cookie", zap.Error(err))
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("внутренняя ошибка сервера")); err != nil {
			h.log.Error("AuthHandler: can't write response", zap.Error(err))
		}
		return
	}
	http.SetCookie(w, jwtCookie)
	if err = WriteResponse(w, http.StatusOK, ""); err != nil {
		h.log.Error("AuthHandler: can't write response", zap.Error(err))
		return
	}
	h.log.Info(fmt.Sprintf("User %s successfully registered", user.Login))

}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user model.User
	b, err := getRequestBody(r)
	if err != nil {
		h.log.Error("AuthHandler:can't get request body", zap.Error(err))
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("внутренняя ошибка сервера")); err != nil {
			h.log.Error("AuthHandler: can't write response", zap.Error(err))
		}
		return
	}
	if err := json.Unmarshal(b, &user); err != nil {
		h.log.Error("AuthHandler:can't unmarshal body", zap.Error(err))
		if err = WriteResponse(w, http.StatusBadRequest, ErrMessage("неверный формат запроса")); err != nil {
			h.log.Error("AuthHandler: can't write response", zap.Error(err))
		}
		return
	}
	h.log.Info(fmt.Sprintf("User %s login attempt", user.Login))
	if user.Login == "" {
		h.log.Error("AuthHandler: login can't be empty", zap.Error(err))
		if err = WriteResponse(w, http.StatusBadRequest, ErrMessage("неверный формат запроса")); err != nil {
			h.log.Error("AuthHandler: can't write response", zap.Error(err))
		}
		return
	}
	ctx := r.Context()
	u, err := h.authService.Check(ctx, &user)
	if err != nil {
		h.log.Error("AuthHandler: can't check user credential", zap.Error(err))
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("внутренняя ошибка сервера")); err != nil {
			h.log.Error("AuthHandler: can't write response", zap.Error(err))
		}
		return
	}
	if u == nil {
		if err = WriteResponse(w, http.StatusUnauthorized, ErrMessage("неверная пара логин/пароль")); err != nil {
			h.log.Error("AuthHandler: can't write response", zap.Error(err))
			return
		}
		return
	}
	token, err := h.auth.GetNewToken(u.ID, u.Login)
	if err != nil {
		h.log.Error("AuthHandler: can't make token", zap.Error(err))
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("внутренняя ошибка сервера")); err != nil {
			h.log.Error("AuthHandler: can't write response", zap.Error(err))
		}
		return
	}
	jwtCookie, err := bakeCookie(token)
	if err != nil {
		h.log.Error("AuthHandler: can't set cookie", zap.Error(err))
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("внутренняя ошибка сервера")); err != nil {
			h.log.Error("AuthHandler: can't write response", zap.Error(err))
		}
		return
	}
	http.SetCookie(w, jwtCookie)
	if err = WriteResponse(w, http.StatusOK, ""); err != nil {
		h.log.Error("AuthHandler: can't write response", zap.Error(err))
		return
	}
	h.log.Info(fmt.Sprintf("User %s successfully logined", user.Login))
}
