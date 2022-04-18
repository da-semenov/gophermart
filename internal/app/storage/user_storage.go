package storage

import (
	"context"
	"errors"
	"github.com/da-semenov/gophermart/internal/app/config"
	"github.com/da-semenov/gophermart/internal/app/database"
	"github.com/da-semenov/gophermart/internal/app/storage/basedbhandler"
	"go.uber.org/zap"
)

type UserRepository struct {
	h basedbhandler.DBHandler
	l *config.Logger
}

func NewUserRepository(dbHandler basedbhandler.DBHandler, log *config.Logger) *UserRepository {
	var repo UserRepository
	repo.h = dbHandler
	repo.l = log
	return &repo
}

func (ur *UserRepository) Save(ctx context.Context, login string, pass string) error {
	if login == "" {
		ur.l.Info("UserRepository: empty login authorization attempt")
		return errors.New("can't register empty login")
	}
	if pass == "" {
		ur.l.Info("UserRepository: empty password authorization attempt")
		return errors.New("empty password")
	}
	var userID int
	row, err := ur.h.QueryRow(ctx, database.CreateUser, login, pass)
	err = row.Scan(&userID)
	if err != nil {
		ur.l.Error("UserRepository: can't get user_id", zap.Error(err))
		return err
	}
	err = ur.h.Execute(ctx, database.CreateAccount, userID)
	if err != nil {
		ur.l.Error("UserRepository: can't create account for user", zap.Error(err))
		return err
	}
	return nil
}

func (ur *UserRepository) Check(ctx context.Context, login string, pass string) (bool, error) {
	if login == "" {
		ur.l.Info("UserRepository: empty login authorization attempt")
		return false, errors.New("can't register empty login")
	}
	if pass == "" {
		ur.l.Info("UserRepository: empty password authorization attempt")
		return false, errors.New("empty password")
	}
	row, err := ur.h.QueryRow(ctx, database.CheckUser, login, pass)
	if err != nil {
		return false, err
	}

	var res int
	err = row.Scan(&res)
	if err != nil {
		return false, nil
	}
	return true, nil
}
