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
func (s *OrderService) GerProductInfo(ctx context.Context, productID int64) (models.Product, error) {
	productInfo, err := s.OrderRepository.GetProductInfo(ctx, productID)
	if err != nil {
		return models.Product{}, err
	}

	return productInfo, nil
}

// update orde status
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID int64, status int) error {
	err := s.OrderRepository.UpdateOrderStatus(ctx, orderID, status)
	if err != nil {
		return err
	}
	return nil
}

func (s *OrderService) GetOrderInfoByOrderID(ctx context.Context, orderID int64) (models.Order, error) {
	orderInfo, err := s.OrderRepository.GetOrderInfoByOrderID(ctx, orderID)
	if err != nil {
		return models.Order{}, err
	}
	return orderInfo, nil
}
func (s *OrderService) GetOrderDetailByOrderDetailID(ctx context.Context, orderDetailID int64) (models.OrderDetail, error) {
	orderDetail, err := s.OrderRepository.GetOrderDetailByOrderDetailID(ctx, orderDetailID)
	if err != nil {
		return models.OrderDetail{}, err
	}
	return orderDetail, nil
}
