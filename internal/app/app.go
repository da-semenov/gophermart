package app

import (
	"context"
	conf "github.com/da-semenov/gophermart/internal/app/config"
	"github.com/da-semenov/gophermart/internal/app/storage"
	"go.uber.org/zap"
	"log"
)

func RunApp() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	config := conf.NewConfig()
	err = config.Init()
	if err != nil {
		logger.Fatal("can't init configuration", zap.Error(err))
	}

	postgresHandler, err := storage.NewPostgresHandler(context.Background(), config.DatabaseDSN)
	if err != nil {
		logger.Fatal("can't init postgres handler", zap.Error(err))
	}

	if config.ReInit {
		err = storage.ClearDatabase(context.Background(), postgresHandler)
		if err != nil {
			logger.Fatal("can't clear database structure", zap.Error(err))
			return
		}
	}
	err = storage.InitDatabase(context.Background(), postgresHandler)
	if err != nil {
		logger.Fatal("can't init database structure", zap.Error(err))
		return
	}
}
