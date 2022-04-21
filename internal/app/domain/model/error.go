package model

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
