package repository

import (
	"context"
	"errors"
	"github.com/da-semenov/gophermart/internal/app/db-queries"
	"github.com/da-semenov/gophermart/internal/app/infrastructure"
	"github.com/da-semenov/gophermart/internal/app/models"
	"github.com/da-semenov/gophermart/internal/app/repository/basedbhandler"
	"go.uber.org/zap"
)

type BalanceRepository struct {
	h basedbhandler.DBHandler
	l *infrastructure.Logger
}

func NewBalanceRepository(dbHandler basedbhandler.DBHandler, log *infrastructure.Logger) (models.BalanceRepository, error) {
	var target BalanceRepository
	if dbHandler == nil {
		return nil, errors.New("can't init balance repository")
	}
	target.h = dbHandler
	target.l = log
	return &target, nil
}

func (r *BalanceRepository) FindWithdrawalByUser(ctx context.Context, userID int) ([]models.Withdrawal, error) {
	rows, err := r.h.Query(ctx, db_queries.GetWithdrawalByUser, userID)
	var resArray []models.Withdrawal
	if err != nil {
		r.l.Error("BalanceRepository: request error", zap.String("query", db_queries.GetWithdrawalByUser), zap.Int("userID", userID), zap.Error(err))
		return nil, err
	}
	for rows.Next() {
		var o models.Withdrawal
		err := rows.Scan(&o.OrderNum, &o.Amount, &o.Status, &o.ProcessedAt)
		if err != nil {
			r.l.Error("BalanceRepository: scan rows error", zap.String("query", db_queries.GetWithdrawalByUser), zap.Int("userID", userID), zap.Error(err))
			break
		}
		resArray = append(resArray, o)
	}
	return resArray, nil
}

func (r *BalanceRepository) LockAccount(ctx context.Context, userID int) (*models.Account, error) {
	row, err := r.h.QueryRow(ctx, db_queries.GetAccountForUpdate, userID)
	if err != nil {
		r.l.Error("BalanceRepository: can't get account for update", zap.Error(err))
		return nil, err
	}
	account := models.Account{}
	err = row.Scan(&account.ID, &account.UserID, &account.Balance, &account.Debit, &account.Credit)
	if err != nil {
		r.l.Error("BalanceRepository: can't get account for update", zap.Error(err))
		if err.Error() == "no rows in result set" {
			return nil, &models.NoRowFound
		} else {
			return nil, err
		}
	}
	return &account, nil
}

func (r *BalanceRepository) SaveAccount(ctx context.Context, account *models.Account) error {
	err := r.h.Execute(ctx, db_queries.UpdateAccountForUser, account.UserID, account.Balance, account.Debit, account.Credit)
	if err != nil {
		r.l.Error("UserRepository: can't create user", zap.Error(err))
		return err
	}
	return nil
}

func (r *BalanceRepository) CreateOperation(ctx context.Context, operation *models.Operation) error {
	err := r.h.Execute(ctx, db_queries.CreateOperation,
		operation.AccountID,
		operation.OrderID,
		operation.OrderNum,
		operation.OperationType,
		operation.Amount,
		operation.ProcessedAt)
	if err != nil {
		r.l.Error("BalanceRepository: can't create operation", zap.Error(err))
		return err
	}
	return nil
}

func (r *BalanceRepository) GetAccount(ctx context.Context, userID int) (*models.Account, error) {
	row, err := r.h.QueryRow(ctx, db_queries.GetAccount, userID)
	if err != nil {
		r.l.Error("BalanceRepository: can't get account for update", zap.Error(err))
		return nil, err
	}
	account := models.Account{}
	err = row.Scan(&account.ID, &account.UserID, &account.Balance, &account.Debit, &account.Credit)
	if err != nil {
		r.l.Error("BalanceRepository: can't get account for update", zap.Error(err))
		if err.Error() == "no rows in result set" {
			return nil, &models.NoRowFound
		} else {
			return nil, err
		}
	}
	return &account, nil
}
