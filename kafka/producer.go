package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
	topic  string
}

func NewProducer(broker, topic string) *Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(broker),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		WriteTimeout: 10 * time.Second,
	}

	return &Producer{
		writer: writer,
		topic:  topic,
	}
}

func (p *Producer) Send(data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(time.Now().String()),
		Value: payload,
	}

	err = p.writer.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Printf("Failed to send Kafka message: %v\n", err)
		return err
	}

	log.Println("Sent Kafka message successfully.")
	return nil
}
