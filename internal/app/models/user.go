package models

import "context"

type UserRepository interface {
	Save(ctx context.Context, login string, pass string) (userID int, err error)
	Check(ctx context.Context, login string, pass string) (bool, error)
}

type User struct {
	ID    int
	Login string
	Pass  string
}
