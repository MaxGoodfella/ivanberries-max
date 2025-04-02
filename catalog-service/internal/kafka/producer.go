package kafka

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(broker, topic string) *Producer {
	createTopic(broker, topic)

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(broker),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: false,
	}

	return &Producer{writer: writer}
}

func createTopic(broker, topic string) {
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		log.Fatalf("error connecting to Kafka: %v", err)
	}
	defer conn.Close()

	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     3,
		ReplicationFactor: 1,
	})
	if err != nil {
		log.Printf("error creating topic %s: %v", topic, err)
		return
	}
}

func (p *Producer) SendMessage(key, value string) error {
	msg := kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("error sending message to Kafka: %v", err)
		return err
	}

	log.Printf("message sent to Kafka, key: %s, value: %s", key, value)
	return nil
}

func (p *Producer) Close() {
	if err := p.writer.Close(); err != nil {
		log.Printf("error closing Kafka-producer: %v", err)
	}
}
