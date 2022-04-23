package domain

import "time"

type Withdrawal struct {
	OrderNum    string    `json:"order"`
	Amount      int       `json:"sum"`
	Status      string    `json:"status"`
	ProcessedAt time.Time `json:"processed_at"`
}

type Withdraw struct {
	OrderNum string `json:"order"`
	Amount   int    `json:"sum"`
}
