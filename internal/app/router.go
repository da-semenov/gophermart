package app

import (
	"github.com/da-semenov/gophermart/internal/app/handlers"
	"github.com/da-semenov/gophermart/internal/app/infrastructure"
	"github.com/da-semenov/gophermart/internal/app/infrastructure/datastore"
	"github.com/da-semenov/gophermart/internal/app/mymiddleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func publicRoutes(
	r chi.Router,
	handler *handlers.AuthHandler,
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
	})
}
