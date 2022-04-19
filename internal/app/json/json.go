package json

import (
	"time"
)

type User struct {
	ID    int
	Login string `json:"login"`
	Pass  string `json:"password"`
}

type Order struct {
	Num      string    `json:"number"`
	UserID   int       `json:"-"`
	Status   string    `json:"status"`
	Accrual  int       `json:"accrual"`
	UploadAt time.Time `json:"upload_at"`
}

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

type Balance struct {
	Current   int `json:"current"`
	Withdrawn int `json:"withdrawn"`
}

type Error struct {
	Msg string `json:"msg"`
}
