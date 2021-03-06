package service

import (
	"context"
	"errors"
	"github.com/da-semenov/gophermart/internal/app/domain"
	"github.com/da-semenov/gophermart/internal/app/infrastructure"
	"github.com/da-semenov/gophermart/internal/app/models"
	"go.uber.org/zap"
	"time"
)

type AccrualClient interface {
	GetAccrual(ctx context.Context, orderNum string) (*domain.Accrual, error)
}

type GophermartClient interface {
	ProcessRequest(ctx context.Context, orderNum string) bool
}

type AccrualService struct {
	dbOrder          models.OrderRepository
	dbBalance        models.BalanceRepository
	accrualClient    AccrualClient
	gophermartClient GophermartClient
	log              *infrastructure.Logger
	enable           bool
}

func NewAccrualService(
	orderRepo models.OrderRepository,
	balanceRepo models.BalanceRepository,
	accrualClient AccrualClient,
	gophermartClient GophermartClient,
	log *infrastructure.Logger,
	enable bool,
) *AccrualService {
	var target AccrualService
	target.dbOrder = orderRepo
	target.dbBalance = balanceRepo
	target.log = log
	target.accrualClient = accrualClient
	target.gophermartClient = gophermartClient
	target.enable = enable
	return &target
}

func (s *AccrualService) StartProcessJob(latency time.Duration) {
	if !s.enable {
		return
	}
	t := time.NewTicker(latency * time.Second)
	defer t.Stop()
	for {
		<-t.C
		ctx := context.Background()
		s.process(ctx)
	}
}

func (s *AccrualService) process(ctx context.Context) {
	orderList, err := s.dbOrder.FindNotProcessed(ctx)
	if err != nil {
		s.log.Error("AccrualService: process. Can't get order list", zap.Error(err))
		return
	}
	for _, order := range orderList {
		s.gophermartClient.ProcessRequest(ctx, order.Num)
	}
}

func (s *AccrualService) ProcessOrder(ctx context.Context, orderNum string) error {
	s.log.Debug("AccrualService: processOrder. Request")
	accrual, err := s.accrualClient.GetAccrual(context.Background(), orderNum)
	if err != nil {
		s.log.Error("AccrualService: processOrder. Can't get accruals from remote service", zap.Error(err))
		return err
	}
	order, err := s.dbOrder.LockOrder(ctx, orderNum)
	if err != nil {
		s.log.Error("AccrualService: processOrder. Can't lock order", zap.Error(err))
		return err
	}
	if accrual.Status == models.OrderStatusProcessed && order.Status != models.OrderStatusProcessed {
		account, err := s.dbBalance.LockAccount(ctx, order.UserID)
		if err != nil {
			s.log.Error("AccrualService: processOrder. Can't lock account", zap.Error(err))
			return err
		}

		operation := models.Operation{
			AccountID:     account.ID,
			Amount:        accrual.Accrual,
			OrderID:       order.ID,
			OrderNum:      order.Num,
			OperationType: models.OperationCredit,
			ProcessedAt:   time.Now().Truncate(time.Second),
		}
		account.Balance += accrual.Accrual
		account.Credit += accrual.Accrual

		order.Status = accrual.Status
		order.UpdatedAt = time.Now().Truncate(time.Second)
		err = s.dbBalance.CreateOperation(ctx, &operation)
		if err != nil {
			s.log.Error("AccrualService: processOrder. Can't create operation", zap.Error(err))
			return err
		}
		err = s.dbBalance.SaveAccount(ctx, account)
		if err != nil {
			s.log.Error("AccrualService: processOrder. Can't save account", zap.Error(err))
			return err
		}
		err = s.dbOrder.UpdateStatus(ctx, order)
		if err != nil {
			s.log.Error("AccrualService: processOrder. Can't save order", zap.Error(err))
			return err
		}
	} else if accrual.Status == models.OrderStatusProcessing || accrual.Status == models.OrderStatusRegistered || accrual.Status == models.OrderStatusInvalid {
		order.Status = accrual.Status
		order.UpdatedAt = time.Now().Truncate(time.Second)
		s.dbOrder.Save(ctx, order)
	} else {
		s.log.Error("AccrualService: processOrder. Received unexpected status", zap.String("OrderNum", order.Num), zap.String("Status", accrual.Status))
		return errors.New("received unexpected status")
	}
	s.log.Debug("AccrualService: processOrder. Success")
	return nil
}
