package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func StartConsumer(broker, topic, groupID string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		GroupID:  groupID,
		MaxBytes: 10e6,
	})
	defer func() {
		if err := reader.Close(); err != nil {
			log.Printf("error closing Kafka consumer: %v", err)
		}
	}()

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("error reading message from Kafka: %v", err)
			return
		}
		fmt.Printf("message received: %s = %s\n", string(msg.Key), string(msg.Value))
	}
}
