package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"orderfc/models"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string) *KafkaProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.LeastBytes{},
	}
	return &KafkaProducer{writer: writer}
}

func (p *KafkaProducer) PublishOrderCreated(ctx context.Context, event models.OrderCreatedEvent) error {
	value, err := json.Marshal(event)
	if err != nil {
		return err
	}
	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("order-%d", event.OrderID)),
		Value: value,
		Topic: "order.created",
	}
	return p.writer.WriteMessages(ctx, msg)
}

func (p *KafkaProducer) PublishProductStockUpdate(ctx context.Context, event models.ProductStockUpdateEvent) error {
	value, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("order-%d", event.OrderID)),
		Value: value,
		Topic: "stock.update",
	}
	return p.writer.WriteMessages(ctx, msg)

}

func (p *KafkaProducer) PublishProductStockRollback(ctx context.Context, event models.ProductStockUpdateEvent) error {
	value, err := json.Marshal(event)
	if err != nil {
		return err
	}
	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("order-%d", event.OrderID)),
		Value: value,
		Topic: "stock.rollback",
	}
	return p.writer.WriteMessages(ctx, msg)
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
