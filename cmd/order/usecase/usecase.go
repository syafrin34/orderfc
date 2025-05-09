package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"orderfc/cmd/order/service"
	"orderfc/infrastructure/constant"
	"orderfc/infrastructure/logger"
	"orderfc/kafka"
	"orderfc/models"
	"time"

	"github.com/sirupsen/logrus"
)

type OrderUsecase struct {
	OrderService  service.OrderService
	KafkaProducer kafka.KafkaProducer
}

func NewOrderUsecase(ou service.OrderService, kafkaProduceer kafka.KafkaProducer) *OrderUsecase {
	return &OrderUsecase{
		OrderService:  ou,
		KafkaProducer: kafkaProduceer,
	}
}

func (uc *OrderUsecase) CheckOutOrder(ctx context.Context, param *models.CheckoutRequest) (int64, error) {
	var orderID int64

	//check idempotency
	if param.IdempotencyToken != "" {
		isExistst, err := uc.OrderService.CheckIdempotency(ctx, param.IdempotencyToken)
		if err != nil {
			return 0, err
		}
		if isExistst {
			return 0, errors.New("order already created, please check again")
		}
	}

	//validate product
	err := uc.validateProducts(param.Items)
	if err != nil {
		return 0, err
	}

	//calculate product amount

	totalQty, totalAmout := uc.calculateOrderSummary(param.Items)

	//construct order detail
	products, orderHistory := uc.constructOrderDetail(param.Items)

	// save order and order detail
	orderDetail := models.OrderDetail{
		Products:     products,     // stringfy json dari param.Items ([]models.CheckoutItem)
		OrderHistory: orderHistory, //
	}

	order := models.Order{
		UserID:          param.UserID,
		Amount:          totalAmout,
		TotalQty:        totalQty,
		Status:          constant.OrderStatusCreated,
		PaymentMethod:   param.PaymentMethod,
		ShippingAddress: param.ShippingAddress,
	}

	orderID, err = uc.OrderService.SaveOrderAndOrderDetail(ctx, &order, &orderDetail)
	if err != nil {
		return 0, err
	}

	// save idempotency
	if param.IdempotencyToken != "" {
		err = uc.OrderService.SaveIdempotency(ctx, param.IdempotencyToken)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{
				"err":   err.Error(),
				"token": param.IdempotencyToken,
			}).Info("uc.OrderService.Saveidempotency got error")
		}
	}

	//todo
	// connect to payment service -> infoin payment request
	// kafka
	// publish order created to kafka
	orderCreatedEvent := models.OrderCreatedEvent{
		OrderID:         orderID,
		UserID:          param.UserID,
		TotalAmount:     order.Amount,
		PaymentMethod:   param.PaymentMethod,
		ShippingAddress: param.ShippingAddress,
	}
	err = uc.KafkaProducer.PublishOrderCreated(ctx, orderCreatedEvent)
	if err != nil {
		return 0, err
	}
	return orderID, nil
}

func (uc *OrderUsecase) validateProducts(ctx context.Context, items []models.CheckoutItem) error {
	seen := map[int64]bool{}
	for _, item := range items {
		productInfo, err := uc.OrderService.GerProductInfo(ctx, item.ProductID)
		if err != nil {
			return fmt.Errorf("failed get product info: %d", item.ProductID)
		}
		//product price
		if item.Price != productInfo.Price {
			return fmt.Errorf("invalid price for product  %d", item.ProductID)
		}
		// dupliacated product
		// if seen[item.ProductID] {
		// 	return fmt.Errorf("duplicate product: %d", item.ProductID)
		// }

		//qty product
		if item.Price <= 0 {
			return fmt.Errorf("invalid price for product %d", item.ProductID)
		}
		if item.Quantity > 1000 {
			return fmt.Errorf("invalid quantity  product %d maximum quantity is 1000", item.ProductID)
		}
		//validate stock

		if item.Quantity > productInfo.Stock {
			return fmt.Errorf("invalid product %d, product stock is only %d", item.ProductID, item.Quantity)
		}

	}

	return nil
}

func (uc *OrderUsecase) calculateOrderSummary(items []models.CheckoutItem) (int, float64) {
	var TotalQty int
	var totalAmount float64

	for _, item := range items {
		TotalQty += item.Quantity
		totalAmount += float64(item.Quantity) * item.Price
	}
	return TotalQty, totalAmount
}

func (uc *OrderUsecase) constructOrderDetail(items []models.CheckoutItem) (string, string) {
	//products, order history
	productJson, _ := json.Marshal(items)
	history := []map[string]interface{}{
		{"status": "created", "timestamp": time.Now()},
	}

	historyJSON, _ := json.Marshal(history)
	return string(productJson), string(historyJSON)

}

func (uc *OrderUsecase) GetOrderHistoryByUserID(ctx context.Context, param models.OrderHistoryParam) ([]models.OrderHistoryResponse, error) {
	orderHistories, err := uc.OrderService.GetOrderHistoryByUserID(ctx, param)
	if err != nil {
		return nil, err
	}
	return orderHistories, nil
}
