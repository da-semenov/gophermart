package models

import "time"

type Operation struct {
	ID            int
	AccountID     int
	OrderID       int
	OrderNum      string
	OperationType string
	Amount        int
	ProcessedAt   time.Time
}
