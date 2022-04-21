package service

import (
	"context"
	"errors"
	"github.com/da-semenov/gophermart/internal/app/domain/model"
	"github.com/da-semenov/gophermart/internal/app/infrastructure"
	"github.com/da-semenov/gophermart/internal/app/models"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	dbUser models.UserRepository
	log    *infrastructure.Logger
}

func NewAuthService(userRepo models.UserRepository, log *infrastructure.Logger) *AuthService {
	var target AuthService
	target.dbUser = userRepo
	target.log = log
	return &target
}

func (s *AuthService) hashPassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 4)
	return string(bytes), err
}

func (s *AuthService) checkPasswordHash(pass string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}

func (s *AuthService) Register(ctx context.Context, user *model.User) (*model.User, error) {
	if user == nil {
		s.log.Debug("AuthService: Register. Got nil user")
		return nil, model.ErrBadParam
	}
	if (user.Login == "") || (user.Pass == "") {
		s.log.Warn("AuthService: Register. Validation error", zap.String("user", user.Login))
		return nil, model.ErrBadParam
	}

	hp, err := s.hashPassword(user.Pass)
	if err != nil {
		s.log.Error("AuthService: Register. Can't calculate hash", zap.String("login", user.Login), zap.Error(err))
		return nil, err
	}
	user.ID, err = s.dbUser.Save(ctx, user.Login, hp)
	if err != nil {
		s.log.Error("AuthService: Register. Can't register user", zap.String("login", user.Login), zap.Error(err))
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Check(ctx context.Context, user *model.User) (*model.User, error) {
	if user == nil {
		s.log.Debug("AuthService: Check. Got nil user")
		return nil, model.ErrBadParam
	}
	if user.Login == "" {
		s.log.Warn("AuthService: Check. Validation error", zap.String("user", user.Login))
		return nil, model.ErrBadParam
	}
	modelUser, err := s.dbUser.GetUserByLogin(ctx, user.Login)
	if err != nil {
		if errors.Is(err, &models.NoRowFound) {
			s.log.Debug("AuthService: Check. User not found in database", zap.String("login", user.Login))
			return nil, nil
		} else {
			s.log.Error("AuthService: Check.", zap.Error(err))
			return nil, err
		}
	}
	if s.checkPasswordHash(user.Pass, modelUser.Pass) {
		user.ID = modelUser.ID
		return user, nil
	} else {
		return nil, nil
	}
}
