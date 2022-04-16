package app

import (
	"context"
	"fmt"
	conf "github.com/da-semenov/gophermart/internal/app/config"
	"github.com/da-semenov/gophermart/internal/app/storage"
	"log"
)

func RunApp() {
	config := conf.NewConfig()
	err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	postgresHandler, err := storage.NewPostgresHandler(context.Background(), config.DatabaseDSN)
	if err != nil {
		fmt.Println("can't init postgres handler", err)
	}

	if config.ReInit {
		err = storage.ClearDatabase(context.Background(), postgresHandler)
		if err != nil {
			fmt.Println("can't clear database structure", err)
			return
		}
	}
	err = storage.InitDatabase(context.Background(), postgresHandler)
	if err != nil {
		fmt.Println("can't init database structure", err)
		return
	}
}
