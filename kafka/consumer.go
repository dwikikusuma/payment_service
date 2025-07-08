package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"payment_service/models"
)

func StartKafkaConsumer(broker string, topic string, handler func(event models.OrderCreatedEvent)) {
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: "PaymentFC",
	})

	go func(r *kafka.Reader) {
		for {
			message, err := r.ReadMessage(context.Background())
			if err != nil {
				fmt.Println("error occurred while reading message on kafka consumer ", err.Error())
				continue
			}

			var event models.OrderCreatedEvent
			err = json.Unmarshal(message.Value, &event)
			if err != nil {
				fmt.Println("error unmarshal event ", err.Error())
				continue
			}

			fmt.Println("message received")
			handler(event)
		}
	}(consumer)

}
