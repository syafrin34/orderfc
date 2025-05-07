package usecase

import "orderfc/cmd/order/service"

type OrderUseCase struct {
	OrderService service.OrderService
}

func NewOrderUsecase(ou service.OrderService) *OrderUseCase {
	return &OrderUseCase{
		OrderService: ou,
	}
}
