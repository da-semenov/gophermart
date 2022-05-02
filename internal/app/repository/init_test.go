package repository

import (
	"context"
	"fmt"
	"github.com/da-semenov/gophermart/internal/app/infrastructure/datastore"
	"go.uber.org/zap"
	"os"
	"testing"
)

const Datasource = "postgresql://practicum_test:practicum_test@127.0.0.1:5432/practicum_test"

var postgresHandler *datastore.PostgresHandlerTX
var Log *zap.Logger

func initDatabase(ctx context.Context, h *datastore.PostgresHandlerTX) {
	err := ClearDatabase(ctx, h)
	if err != nil {
		fmt.Println("can't clear db-queries")
		panic(err)
	}
	err = InitDatabase(ctx, h)
	if err != nil {
		fmt.Println("can't init db-queries")
		panic(err)
	}
}

func TestMain(m *testing.M) {
	var err error
	Log, err = zap.NewDevelopment()
	if err != nil {
		fmt.Println("init zap log failed")
		panic(err)
	}
	postgresHandler, err = datastore.NewPostgresHandlerTX(context.Background(), Datasource, Log)
	if err != nil {
		fmt.Println("can't init PostgresHandler")
		panic(err)
	}
	os.Exit(m.Run())
}
