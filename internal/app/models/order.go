package models

import (
	"context"
	"time"
)

type OrderRepository interface {
	Save(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, orderID int) (*Order, error)
	GetByNum(ctx context.Context, num string) (*Order, error)
	UpdateStatus(ctx context.Context, userID int, num string, statusNew string) error
	FindByUser(ctx context.Context, userID int) ([]Order, error)
}

type Order struct {
	ID        int
	UserID    int
	Num       string
	Status    string
	UploadAt  time.Time
	UpdatedAt time.Time
}

const (
	OrderStatusNew        = "NEW"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusInvalid    = "INVALID"
	OrderStatusProcessed  = "PROCESSED"
)
