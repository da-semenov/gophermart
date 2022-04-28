package service

import (
	"context"
	"github.com/da-semenov/gophermart/internal/app/domain"
	"github.com/da-semenov/gophermart/internal/app/infrastructure"
	"github.com/da-semenov/gophermart/internal/app/models"
	"go.uber.org/zap"
	"time"
)

type BalanceService struct {
	dbBalance models.BalanceRepository
	log       *infrastructure.Logger
}

func NewBalanceService(balanceRepo models.BalanceRepository, log *infrastructure.Logger) *BalanceService {
	var target BalanceService
	target.dbBalance = balanceRepo
	target.log = log
	return &target
}

func (s *BalanceService) mapWithdrawalModelToDomain(src models.Withdrawal) domain.Withdrawal {
	return domain.Withdrawal{
		OrderNum:    src.OrderNum,
		Amount:      src.Amount,
		Status:      src.Status,
		ProcessedAt: src.ProcessedAt,
	}
}

func (s *BalanceService) mapWithdrawalListModelToDomain(src []models.Withdrawal) (resList []domain.Withdrawal) {
	for _, o := range src {
		resList = append(resList, s.mapWithdrawalModelToDomain(o))
	}
	return resList
}

func (s *BalanceService) GetCurrentBalance(ctx context.Context, userID int) (*domain.Balance, error) {
	if userID == 0 {
		s.log.Debug("BalanceService: GetCurrentBalance. Got nil userID")
		return nil, domain.ErrBadParam
	}

	account, err := s.dbBalance.GetAccount(ctx, userID)
	if err != nil {
		s.log.Debug("BalanceService: GetCurrentBalance. Can't get current balance")
		return nil, err
	}

	return &domain.Balance{
		Current:   account.Balance,
		Withdrawn: account.Credit,
	}, nil

}

func (s *BalanceService) Withdraw(ctx context.Context, obj *domain.Withdraw, userID int) error {
	if userID == 0 {
		s.log.Debug("BalanceService: Withdraw. Got nil userID")
		return domain.ErrBadParam
	}
	if obj == nil {
		s.log.Debug("BalanceService: Withdraw. Got nil order")
		return domain.ErrBadParam
	}
	if !CheckOrderNum(obj.OrderNum) {
		s.log.Debug("BalanceService: Withdraw. Order num validation error", zap.String("orderNum", obj.OrderNum))
		return domain.ErrBadOrderNum
	}
	account, err := s.dbBalance.LockAccount(ctx, userID)
	if err != nil {
		s.log.Error("BalanceService: Withdraw. Unexpected error", zap.Error(err))
		return err
	}

	if account.Balance < obj.Amount {
		s.log.Debug("BalanceService: Withdraw. In account not enough funds")
		return domain.ErrNotEnoughFunds
	}

	operation := models.Operation{
		AccountID:     account.ID,
		Amount:        obj.Amount,
		OrderNum:      obj.OrderNum,
		OperationType: models.OperationDebit,
		ProcessedAt:   time.Now().Truncate(time.Second),
	}
	err = s.dbBalance.CreateOperation(ctx, &operation)

	if err != nil {
		s.log.Error("BalanceService: Withdraw. Can't save operation", zap.Error(err))
		return err
	}
	account.Balance -= obj.Amount
	account.Debit += obj.Amount
	err = s.dbBalance.SaveAccount(ctx, account)
	if err != nil {
		s.log.Error("BalanceService: Withdraw. Can't save account", zap.Error(err))
		return err
	}
	return nil
}

func (s *BalanceService) GetWithdrawalsList(ctx context.Context, userID int) ([]domain.Withdrawal, error) {
	if userID == 0 {
		s.log.Debug("BalanceService: GetWithdrawalsList. Got nil userID")
		return nil, domain.ErrBadParam
	}

	withdrawalList, err := s.dbBalance.FindWithdrawalByUser(ctx, userID)
	if err != nil {
		s.log.Error("BalanceService: GetWithdrawalsList. Can't get withdrawal list",
			zap.Int("userID", userID),
			zap.Error(err),
		)
		return nil, err
	}
	resList := s.mapWithdrawalListModelToDomain(withdrawalList)
	return resList, nil
}
