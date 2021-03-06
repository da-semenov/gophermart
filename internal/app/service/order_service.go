package service

import (
	"context"
	"github.com/da-semenov/gophermart/internal/app/domain"
	"github.com/da-semenov/gophermart/internal/app/infrastructure"
	"github.com/da-semenov/gophermart/internal/app/models"
	"go.uber.org/zap"
	"time"
)

type OrderService struct {
	dbOrder          models.OrderRepository
	log              *infrastructure.Logger
	EnableValidation bool
}

func NewOrderService(orderRepo models.OrderRepository, log *infrastructure.Logger, enableValidation bool) *OrderService {
	var target OrderService
	target.dbOrder = orderRepo
	target.log = log
	target.EnableValidation = enableValidation
	return &target
}

func (s *OrderService) mapOrderDomainToModel(src *domain.Order) *models.Order {
	return &models.Order{
		UserID:   src.UserID,
		Num:      src.Num,
		Status:   src.Status,
		UploadAt: src.UploadAt,
	}
}

func (s *OrderService) mapOrderModelToDomain(src *models.Order) *domain.Order {
	return &domain.Order{
		UserID:   src.UserID,
		Num:      src.Num,
		Status:   src.Status,
		Accrual:  src.Accrual,
		UploadAt: src.UploadAt.Truncate(time.Second),
	}
}

func (s *OrderService) mapOrderListModelToDomain(src []models.Order) (resList []domain.Order) {
	for _, o := range src {
		o := o
		resList = append(resList, *s.mapOrderModelToDomain(&o))
	}
	return resList
}

func (s *OrderService) Save(ctx context.Context, order *domain.Order) error {
	if order == nil {
		s.log.Debug("OrderService: Save. Got nil order")
		return domain.ErrBadParam
	}
	if (order.UserID == 0) || (order.Num == "") {
		s.log.Debug("OrderService: Save. Validation error")
		return domain.ErrBadParam
	}
	if s.EnableValidation && !CheckOrderNum(order.Num) {
		s.log.Debug("OrderService: Save. Order num validation error")
		return domain.ErrBadOrderNum
	}

	exOrder, err := s.dbOrder.GetByNum(ctx, order.Num)
	if err == nil {
		if exOrder.UserID == order.UserID {
			return domain.ErrOrderRegistered
		}

		return domain.ErrOrderRegisteredByAnotherUser

	} else if err != &models.NoRowFound {
		s.log.Error("OrderService: Save. Unexpected error", zap.Error(err))
		return err
	}
	modelOrder := s.mapOrderDomainToModel(order)
	modelOrder.Status = "NEW"
	modelOrder.UploadAt = time.Now().Truncate(time.Microsecond)
	modelOrder.UpdatedAt = time.Now().Truncate(time.Microsecond)

	err = s.dbOrder.Save(ctx, modelOrder)
	if err != nil {
		s.log.Error("OrderService: Save. Can't save order",
			zap.Int("userID", order.UserID),
			zap.String("num", order.Num),
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (s *OrderService) GetOrderList(ctx context.Context, userID int) ([]domain.Order, error) {
	if userID == 0 {
		s.log.Debug("OrderService: GetOrderList. Got nil userID")
		return nil, domain.ErrBadParam
	}

	orderList, err := s.dbOrder.FindByUser(ctx, userID)
	if err != nil {
		s.log.Error("OrderService: GetOrderList. Can't get order list",
			zap.Int("userID", userID),
			zap.Error(err),
		)
		return nil, err
	}

	resList := s.mapOrderListModelToDomain(orderList)
	return resList, nil
}
