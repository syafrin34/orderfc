package service

import "orderfc/cmd/order/repository"

type OrderService struct {
	OrderRepository repository.OrderRepository
}

func NewOrderService(op repository.OrderRepository) *OrderService {
	return &OrderService{
		OrderRepository: op,
	}
}
