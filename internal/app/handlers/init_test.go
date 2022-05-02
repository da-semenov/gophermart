package handlers

import (
	"context"
	"fmt"
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
	log, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println("init zap log failed")
		panic(err)
	}
	postgresHandler, err := datastore.NewPostgresHandlerTX(context.Background(), Datasource, log)
	if err != nil {
		fmt.Println("can't init PostgresHandler")
		panic(err)
	}
	repository.ClearDatabase(context.Background(), postgresHandler)
	repository.InitDatabase(context.Background(), postgresHandler)
	repo, err := repository.NewUserRepository(postgresHandler, log)
	if err != nil {
		fmt.Println("can't init UserRepo")
		panic(err)
	}
	authService = service.NewAuthService(repo, log)
	auth = NewAuth("secret")
	os.Exit(m.Run())
}
