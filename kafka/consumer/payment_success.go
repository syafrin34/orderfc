package consumer

import (
	"context"
	"encoding/json"
	"orderfc/cmd/order/service"
	"orderfc/infrastructure/constant"
	"orderfc/infrastructure/logger"
	kafkaFC "orderfc/kafka"
	"orderfc/models"

	"github.com/segmentio/kafka-go"
)

type PaymentConsumer struct {
	Reader       *kafka.Reader
	Producer     kafkaFC.KafkaProducer
	OrderService service.OrderService
}

// func NewPaymentFailedConsumer(brokers []string, topic string, orderService service.OrderService, kafkaProducer kafkaFC.KafkaProducer) *PaymentConsumer {
// 	reader := kafka.NewReader(kafka.ReaderConfig{
// 		Brokers: brokers,
// 		Topic:   topic,
// 		GroupID: "orderfc",
// 	})
// 	return &PaymentConsumer{
// 		Reader:       reader,
// 		OrderService: orderService,
// 		Producer:     kafkaProducer,
// 	}
// }

func NewPaymentSuccessConsumer(brokers []string, topic string, orderService service.OrderService, kafkaProducer kafkaFC.KafkaProducer) *PaymentConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: "orderfc",
	})
	return &PaymentConsumer{
		Reader:       reader,
		OrderService: orderService,
		Producer:     kafkaProducer,
	}
}

func (c *PaymentConsumer) StartPaymentSuccessConsumer(ctx context.Context) {
	logger.Logger.Println("[KAFKA] Listening to topic: payment.success")
	for {
		message, err := c.Reader.ReadMessage(ctx)
		if err != nil {
			logger.Logger.Println("[KAFKA] error rad message: ", err)
			continue
		}
		var event models.PaymentUpdateStatusEvent
		err = json.Unmarshal(message.Value, &event)
		if err != nil {
			logger.Logger.Println("[KAFKA] error unmarshal event message value: ", err)
			continue
		}
		logger.Logger.Printf("[KAFKA] received payment.success event for order id #%d", event.OrderID)

		// update db

		err = c.OrderService.UpdateOrderStatus(ctx, event.OrderID, constant.OrderStatusCompleted)
		if err != nil {
			logger.Logger.Println("[KAFKA] error update order status", err)
			continue
		}

		// get order info from db
		// ordeInfo, err := c.OrderService.GetOrderInfoByOrderID(ctx, event.OrderID)
		// if err != nil {
		// 	logger.Logger.Println("[KAFKA] error get order info", err)
		// 	continue
		// }

		// get order detail
		// orderDetail, err := c.OrderService.GetOrderDetailByOrderDetailID(ctx, ordeInfo.OrderDetailID)
		// if err != nil {
		// 	logger.Logger.Println("[KAFKA] error get order detail", err)
		// 	continue
		// }

		// get product list from order detail

		// var products []models.CheckoutItem
		// err = json.Unmarshal([]byte(orderDetail.Products), &products)
		// if err != nil {
		// 	logger.Logger.Println("[KAFKA] error get product list from order detail", err)
		// 	continue
		// }

		// publish event product service
		// err = c.Producer.PublishProductStockUpdate(ctx, models.ProductStockUpdateEvent{
		// 	OrderID:   event.OrderID,
		// 	Products:  convertCheckOutItemsToProductItems(products),
		// 	EventTime: time.Now(),
		// })
		// if err != nil {
		// 	logger.Logger.Println("[KAFKA] error publish product stock info", err)
		// 	continue
		// }

	}
}

// func convertCheckOutItemsToProductItems(source []models.CheckoutItem) []models.ProductItem {
// 	result := make([]models.ProductItem, len(source))

// 	for index, item := range source {
// 		result[index] = models.ProductItem{
// 			ProductID: item.ProductID,
// 			Qty:       item.Quantity,
// 		}
// 	}
// 	return result
// }
