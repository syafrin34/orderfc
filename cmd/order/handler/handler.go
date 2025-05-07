package handler

import (
	"orderfc/cmd/order/usecase"
)

type OrderHandler struct {
	OrderUseCase usecase.OrderUseCase
}

func NewOrderHandler(oh usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		OrderUseCase: oh,
	}
}
