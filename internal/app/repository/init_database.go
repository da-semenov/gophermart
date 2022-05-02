package repository

import (
	"context"
	"fmt"
	"github.com/da-semenov/gophermart/internal/app/dbqueries"
	"github.com/da-semenov/gophermart/internal/app/repository/basedbhandler"
)

func InitDatabase(ctx context.Context, h basedbhandler.DBHandler) error {
	err := h.Execute(ctx, dbqueries.CreateDatabaseStructure)
	if err != nil {
		return err
	}
	fmt.Println("dbqueries structure created successfully")
	return nil
}

func ClearDatabase(ctx context.Context, h basedbhandler.DBHandler) error {
	err := h.Execute(ctx, dbqueries.ClearDatabaseStructure)
	if err != nil {
		return err
	}
	return nil
}
