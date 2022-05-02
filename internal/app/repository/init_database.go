package repository

import (
	"context"
	"fmt"
	"github.com/da-semenov/gophermart/internal/app/db-queries"
	"github.com/da-semenov/gophermart/internal/app/repository/basedbhandler"
)

func InitDatabase(ctx context.Context, h basedbhandler.DBHandler) error {
	err := h.Execute(ctx, db_queries.CreateDatabaseStructure)
	if err != nil {
		return err
	}
	fmt.Println("db-queries structure created successfully")
	return nil
}

func ClearDatabase(ctx context.Context, h basedbhandler.DBHandler) error {
	err := h.Execute(ctx, db_queries.ClearDatabaseStructure)
	if err != nil {
		return err
	}
	return nil
}
