package models

import (
	"context"
	"time"
)

type BalanceRepository interface {
	FindWithdrawalByUser(ctx context.Context, userID int) ([]Withdrawal, error)
	LockAccount(ctx context.Context, userID int) (*Account, error)
	SaveAccount(ctx context.Context, account *Account) error
	CreateOperation(ctx context.Context, operation *Operation) error
	GetAccount(ctx context.Context, userID int) (*Account, error)
}

type Withdrawal struct {
	OrderNum    string
	Amount      int
	Status      string
	ProcessedAt time.Time
}
