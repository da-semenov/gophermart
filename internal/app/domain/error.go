package domain

import (
	"errors"
)

type Error struct {
	Msg string `json:"msg"`
}

var ErrDuplicateKey = errors.New("duplicate key")
var ErrNotFound = errors.New("no rows in result set")
var ErrBadParam = errors.New("bad param occurred")
var ErrUnauthorized = errors.New("user unauthorized")

var ErrOrderRegistered = errors.New("order registered early")
var ErrOrderRegisteredByAnotherUser = errors.New("order registered early by another user")
var ErrBadOrderNum = errors.New("bad order num")
var ErrNotEnoughFunds = errors.New("not enough funds")
