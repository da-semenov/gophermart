package repository

import (
	"context"
	"errors"
	"github.com/da-semenov/gophermart/internal/app/database"
	"github.com/da-semenov/gophermart/internal/app/infrastructure"
	"github.com/da-semenov/gophermart/internal/app/models"
	"github.com/da-semenov/gophermart/internal/app/repository/basedbhandler"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"go.uber.org/zap"
)

type OrderRepository struct {
	h basedbhandler.DBHandler
	l *infrastructure.Logger
}

func NewOrderRepository(dbHandler basedbhandler.DBHandler, log *infrastructure.Logger) (models.OrderRepository, error) {
	var target OrderRepository
	if dbHandler == nil {
		return nil, errors.New("can't init order repository")
	}
	target.h = dbHandler
	target.l = log
	return &target, nil
}

func (or *OrderRepository) Save(ctx context.Context, order *models.Order) error {
	err := or.h.Execute(ctx, database.CreateOrder, order.UserID, order.Num, order.Status, order.UploadAt, order.UpdatedAt)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return &models.UniqueViolation
		}
	}
	return err
}

func (or *OrderRepository) GetByID(ctx context.Context, orderID int) (*models.Order, error) {
	var res models.Order
	row, err := or.h.QueryRow(ctx, database.GetOrderByID, orderID)
	if err != nil {
		or.l.Error("OrderRepository: request error", zap.String("query", database.GetOrderByID), zap.Int("orderID", orderID), zap.Error(err))
		return nil, err
	}
	err = row.Scan(&res.ID, &res.UserID, &res.Num, &res.Status, &res.UploadAt, &res.UpdatedAt)
	if err != nil && err.Error() == "no rows in result set" {
		return nil, &models.NoRowFound
	}
	if err != nil {
		or.l.Error("OrderRepository: scan rows error", zap.String("query", database.GetOrderByID), zap.Int("orderID", orderID), zap.Error(err))
		return nil, err
	}
	return &res, nil
}

func (or *OrderRepository) GetByNum(ctx context.Context, num string) (*models.Order, error) {
	var res models.Order
	row, err := or.h.QueryRow(ctx, database.GetOrderByNum, num)
	if err != nil {
		or.l.Error("OrderRepository: request error", zap.String("query", database.GetOrderByNum), zap.String("Num", num), zap.Error(err))
		return nil, err
	}
	err = row.Scan(&res.ID, &res.UserID, &res.Num, &res.Status, &res.UploadAt, &res.UpdatedAt)
	if err != nil && err.Error() == "no rows in result set" {
		return nil, &models.NoRowFound
	}
	if err != nil {
		or.l.Error("OrderRepository: scan rows error", zap.String("query", database.GetOrderByNum), zap.String("Num", num), zap.Error(err))
		return nil, err
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
		or.l.Error("OrderRepository: request error", zap.String("query", database.FindOrdersByUser), zap.Int("userID", userID), zap.Error(err))
		return nil, err
	}

	for rows.Next() {
		var o models.Order
		err := rows.Scan(&o.ID, &o.Num, &o.UserID, &o.Status, &o.UploadAt, &o.UpdatedAt)
		if err != nil {
			or.l.Error("OrderRepository: scan rows error", zap.String("query", database.FindOrdersByUser), zap.Int("userID", userID), zap.Error(err))
			break
		}
		resArray = append(resArray, o)
	}
	return resArray, nil
}

func (or *OrderRepository) LockOrder(ctx context.Context, OrderNum string) (*models.Order, error) {
	var res models.Order
	row, err := or.h.QueryRow(ctx, database.GetOrderByNumForUpdate, OrderNum)
	if err != nil {
		or.l.Error("OrderRepository: can't get order for update", zap.Error(err))
		return nil, err
	}
	err = row.Scan(&res.ID, &res.UserID, &res.Num, &res.Status, &res.UploadAt, &res.UpdatedAt)

	if err != nil {
		or.l.Error("OrderRepository: can't get account for update", zap.Error(err))
		if err.Error() == "no rows in result set" {
			return nil, &models.NoRowFound
		} else {
			return nil, err
		}
	}
	return &res, nil
}

func (or *OrderRepository) FindNotProcessed(ctx context.Context) ([]models.Order, error) {
	const rowCountLimit = 20
	rows, err := or.h.Query(ctx, database.FindOrderByStatuses, models.OrderStatusProcessing, models.OrderStatusNew, models.OrderStatusRegistered, "", "", rowCountLimit)
	var resArray []models.Order
	if err != nil {
		or.l.Error("OrderRepository: request error", zap.String("query", database.FindOrderByStatuses), zap.Error(err))
		return nil, err
	}
	for rows.Next() {
		var o models.Order
		err := rows.Scan(&o.ID, &o.UserID, &o.Num, &o.Status, &o.UploadAt, &o.UpdatedAt)
		if err != nil {
			or.l.Error("OrderRepository: scan rows error", zap.String("query", database.FindOrderByStatuses), zap.Error(err))
			break
		}
		resArray = append(resArray, o)
	}
	return resArray, nil
}
