package app

import (
	"github.com/da-semenov/gophermart/internal/app/handlers"
	"github.com/da-semenov/gophermart/internal/app/infrastructure"
	"github.com/da-semenov/gophermart/internal/app/infrastructure/datastore"
	"github.com/da-semenov/gophermart/internal/app/mymiddleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

func publicRoutes(
	r chi.Router,
	handler *handlers.AuthHandler,
	accrual *handlers.AccrualHandler,
	postgresHandlerTx *datastore.PostgresHandlerTX,
	log *infrastructure.Logger,
) {
	r.Group(func(router chi.Router) {
		router.Use(middleware.CleanPath)
		router.Use(middleware.Logger)
		router.Use(middleware.Recoverer)
		router.Use(mymiddleware.Transactional(postgresHandlerTx, log))
		router.Post("/api/user/register", handler.Register)
		router.Post("/api/user/login", handler.Login)
		router.Post("/api/accrual/process/{orderNum}", accrual.ProcessOrder)
	})
}

func protectedOrderRoutes(
	r chi.Router,
	tokenAuth *jwtauth.JWTAuth,
	postgresHandlerTx *datastore.PostgresHandlerTX,
	handler *handlers.OrderHandler,
	log *infrastructure.Logger,
) {
	r.Group(func(router chi.Router) {
		router.Use(middleware.CleanPath)
		router.Use(middleware.Logger)
		router.Use(middleware.Recoverer)
		router.Use(jwtauth.Verifier(tokenAuth))
		router.Use(jwtauth.Authenticator)
		router.Use(mymiddleware.Transactional(postgresHandlerTx, log))
		router.Post("/api/user/orders", handler.RegisterNewOrder)
		router.Get("/api/user/orders", handler.GetOrderList)
	})
}

func protectedBalanceRoutes(
	r chi.Router,
	tokenAuth *jwtauth.JWTAuth,
	postgresHandlerTx *datastore.PostgresHandlerTX,
	handler *handlers.BalanceHandler,
	log *infrastructure.Logger,
) {
	r.Group(func(router chi.Router) {
		router.Use(middleware.CleanPath)
		router.Use(middleware.Logger)
		router.Use(middleware.Recoverer)
		router.Use(jwtauth.Verifier(tokenAuth))
		router.Use(jwtauth.Authenticator)
		router.Use(mymiddleware.Transactional(postgresHandlerTx, log))
		router.Get("/api/user/balance", handler.GetBalance)
		router.Post("/api/user/balance/withdraw", handler.Withdraw)
		router.Get("/api/user/balance/withdrawals", handler.GetWithdrawalsList)
	})
}
