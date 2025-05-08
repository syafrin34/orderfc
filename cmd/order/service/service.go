package service

import (
	"context"
	"orderfc/cmd/order/repository"
	"orderfc/models"

	"gorm.io/gorm"
)

type OrderService struct {
	OrderRepository repository.OrderRepository
}

func NewOrderService(op repository.OrderRepository) *OrderService {
	return &OrderService{
		OrderRepository: op,
	}
}

// check idempotency
func (s *OrderService) CheckIdempotency(ctx context.Context, token string) (bool, error) {
	isExistst, err := s.OrderRepository.CheckIdempotency(ctx, token)
	if err != nil {
		return false, err
	}
	return isExistst, nil
}

//save idmpotency

func (s *OrderService) SaveIdempotency(ctx context.Context, token string) error {
	err := s.OrderRepository.SaveIdempotency(ctx, token)
	if err != nil {
		return err
	}
	return nil
}

// save order dan order detail
func (s *OrderService) SaveOrderAndOrderDetail(ctx context.Context, order *models.Order, orderDetail *models.OrderDetail) (int64, error) {
	var orderID int64
	err := s.OrderRepository.WithTransaction(ctx, func(tx *gorm.DB) error {
		err := s.OrderRepository.InsertorderDetailTx(ctx, tx, orderDetail)
		if err != nil {
			return err
		}

		//set order detail id to order table
		order.OrderDetailID = orderDetail.ID
		err = s.OrderRepository.InsertorderTx(ctx, tx, order)
		if err != nil {
			return err
		}

		orderID = order.ID
		return nil
	})

	if err != nil {
		return 0, err
	}
	return orderID, nil
}

func (s *OrderService) GetOrderHistoryByUserID(ctx context.Context, param models.OrderHistoryParam) ([]models.OrderHistoryResponse, error) {
	orderHistories, err := s.OrderRepository.GetOrderHistoryByUserID(ctx, param)
	if err != nil {
		return nil, err
	}
	return orderHistories, nil
}
