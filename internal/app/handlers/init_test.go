package handlers

import (
	"context"
	"github.com/da-semenov/gophermart/internal/app/infrastructure/datastore"
	"github.com/da-semenov/gophermart/internal/app/repository"
	"github.com/da-semenov/gophermart/internal/app/service"
	"go.uber.org/zap"
	"os"
	"testing"
)

var (
	log         *zap.Logger
	authService *service.AuthService
	auth        *Auth
)

const Datasource = "postgresql://practicum_test:practicum_test@127.0.0.1:5432/practicum_test"

func TestMain(m *testing.M) {
	log, _ = zap.NewDevelopment()
	postgresHandler, _ := datastore.NewPostgresHandler(context.Background(), Datasource)
	repository.ClearDatabase(context.Background(), postgresHandler)
	repository.InitDatabase(context.Background(), postgresHandler)
	repo, _ := repository.NewUserRepository(postgresHandler, log)
	authService = service.NewAuthService(repo, log)
	auth = NewAuth("secret")
	os.Exit(m.Run())
}
