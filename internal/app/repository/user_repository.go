package repository

import (
	"context"
	"errors"
	"github.com/da-semenov/gophermart/internal/app/database"
	"github.com/da-semenov/gophermart/internal/app/infrastructure"
	"github.com/da-semenov/gophermart/internal/app/models"
	"github.com/da-semenov/gophermart/internal/app/repository/basedbhandler"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"go.uber.org/zap"
)

type UserRepository struct {
	h basedbhandler.DBHandler
	l *infrastructure.Logger
}

func NewUserRepository(dbHandler basedbhandler.DBHandler, log *infrastructure.Logger) (models.UserRepository, error) {
	var repo UserRepository
	if dbHandler == nil {
		return nil, errors.New("can't init user repository")
	}
	repo.h = dbHandler
	repo.l = log
	return &repo, nil
}

func (ur *UserRepository) Save(ctx context.Context, login string, pass string) (int, error) {
	var userID int
	row, err := ur.h.QueryRow(ctx, database.GetNextUserID)
	if err != nil {
		ur.l.Error("UserRepository: can't get userID")
		return 0, err
	}
	err = row.Scan(&userID)
	if err != nil {
		ur.l.Error("UserRepository: can't get userID")
		return 0, err
	}

	err = ur.h.Execute(ctx, database.CreateUser, userID, login, pass)
	if err != nil {
		ur.l.Error("UserRepository: can't create user", zap.Error(err))
		return 0, err
	}
	err = ur.h.Execute(ctx, database.CreateAccount, userID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return 0, &models.UniqueViolation
		}
	}
	if err != nil {
		ur.l.Error("UserRepository: can't create account for user", zap.Error(err))
		return 0, err
	}
	return userID, nil
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

func (ur *UserRepository) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	row, err := ur.h.QueryRow(ctx, database.GetUserByLogin, login)
	if err != nil {
		return nil, err
	}
	var res models.User
	err = row.Scan(&res.ID, &res.Login, &res.Pass)
	if err != nil && err.Error() == "no rows in result set" {
		return nil, &models.NoRowFound
	}

	if err != nil {
		return nil, nil
	}
	return &res, nil
}
