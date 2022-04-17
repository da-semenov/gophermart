package storage

import (
	"context"
	"errors"
	"github.com/da-semenov/gophermart/internal/app/config"
	"github.com/da-semenov/gophermart/internal/app/database"
	"github.com/da-semenov/gophermart/internal/app/models"
	"github.com/da-semenov/gophermart/internal/app/storage/basedbhandler"
	"go.uber.org/zap"
)

type OrderRepository struct {
	h basedbhandler.DBHandler
	l *config.Logger
}

func NewOrderRepository(dbHandler basedbhandler.DBHandler, log *config.Logger) (models.OrderRepository, error) {
	var repo OrderRepository
	if dbHandler == nil {
		return nil, errors.New("can't init order repository")
	}
	repo.h = dbHandler
	repo.l = log
	return &repo, nil
}

func (or *OrderRepository) Save(ctx context.Context, order models.Order) error {
	err := or.h.Execute(ctx, database.CreateOrder, order.UserID, order.Number, order.Status, order.UploadAt, order.UpdatedAt)
	return err
}

func (or *OrderRepository) GetByID(ctx context.Context, orderID int) (*models.Order, error) {
	var res models.Order
	row, err := or.h.QueryRow(ctx, database.GetOrderByID, orderID)
	if err != nil {
		or.l.Error("request error", zap.String("query", database.GetOrderByID), zap.Int("user_id", orderID), zap.Error(err))
		return nil, err
	}
	err = row.Scan(&res.ID, &res.UserID, &res.Number, &res.Status, &res.UploadAt, &res.UpdatedAt)
	if err != nil {
		or.l.Error("scan rows error", zap.String("query", database.GetOrderByID), zap.Int("user_id", orderID), zap.Error(err))
	}
	return &res, nil
}

func (or *OrderRepository) GetByNumber(ctx context.Context, num string) (*models.Order, error) {
	var res models.Order
	row, err := or.h.QueryRow(ctx, database.GetOrderByNum, num)
	if err != nil {
		or.l.Error("request error", zap.String("query", database.GetOrderByNum), zap.String("user_id", num), zap.Error(err))
		return nil, err
	}
	err = row.Scan(&res.ID, &res.UserID, &res.Number, &res.Status, &res.UploadAt, &res.UpdatedAt)
	if err != nil {
		or.l.Error("scan rows error", zap.String("query", database.GetOrderByNum), zap.String("user_id", num), zap.Error(err))
	}
	return &res, nil
}

func (or *OrderRepository) UpdateStatus(ctx context.Context, userID int, num string, statusNew string) error {
	err := or.h.Execute(ctx, database.UpdateOrderStatus, userID, num, statusNew)
	return err
}

func (or *OrderRepository) FindByUser(ctx context.Context, userID int) ([]models.Order, error) {
	rows, err := or.h.Query(ctx, database.FindOrdersByUser, userID)
	var resArray []models.Order

	if err != nil {
		or.l.Error("request error", zap.String("query", database.FindOrdersByUser), zap.Int("user_id", userID), zap.Error(err))
	}

	for rows.Next() {
		var o models.Order
		err := rows.Scan(&o.ID, &o.Number, &o.UserID, &o.Status, &o.UploadAt, &o.UpdatedAt)
		if err != nil {
			or.l.Error("scan rows error", zap.String("query", database.FindOrdersByUser), zap.Int("user_id", userID), zap.Error(err))
			break
		}
		resArray = append(resArray, o)
	}
	return resArray, nil
}
