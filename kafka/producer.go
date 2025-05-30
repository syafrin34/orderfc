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

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
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
	}
	return p.writer.WriteMessages(ctx, msg)

}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
