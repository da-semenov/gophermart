package domain

import (
	"errors"
)

type Error struct {
	Msg string `json:"msg"`
}

var ErrDuplicateKey = errors.New("duplicate key")
var ErrBadParam = errors.New("bad param occurred")
var ErrTooManyRequest = errors.New("too many request to remote service")
var ErrRemoteServiceError = errors.New("remote service error")

var ErrOrderRegistered = errors.New("order registered early")
var ErrOrderRegisteredByAnotherUser = errors.New("order registered early by another user")
var ErrBadOrderNum = errors.New("bad order num")
var ErrNotEnoughFunds = errors.New("not enough funds")
