package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func PublishMessage(broker, topic string, message interface{}) error {
	// Mở kết nối tới Kafka
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	defer w.Close()

	// Mã hóa message thành JSON
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Gửi message vào Kafka
	err = w.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(time.Now().Format(time.RFC3339)),
		Value: data,
	})
	if err != nil {
		log.Println("[Kafka] Error sending message:", err)
		return err
	}

	log.Println("[Kafka] Message sent successfully:", string(data))
	return nil
}
