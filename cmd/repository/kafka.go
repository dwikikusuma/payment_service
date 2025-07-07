package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
)

type PaymentEventPublisher interface {
	PublishPaymentSuccess(orderID int64) error
}

type KafkaEventPublisher struct {
	writer *kafka.Writer
}

func NewKafkaEventPublisher(writer *kafka.Writer) PaymentEventPublisher {
	return &KafkaEventPublisher{writer: writer}
}

func (k *KafkaEventPublisher) PublishPaymentSuccess(orderID int64) error {
	payload := map[string]interface{}{
		"order_id": orderID,
		"status":   "paid",
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return k.writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(fmt.Sprintf("order-%d", orderID)),
		Value: data,
	})
}
