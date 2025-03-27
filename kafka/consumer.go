package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"intent/config"
	"intent/request"
	"intent/service"

	"github.com/segmentio/kafka-go"
)

// StartKafkaConsumer lắng nghe message từ Kafka
func StartKafkaConsumer(userService *service.UserService) {
	cfg := config.LoadConfig()
	kafkaBroker, topic := config.GetKafkaConfig(cfg)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaBroker},
		Topic:    topic,
		GroupID:  "user-group",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	log.Println("Kafka Consumer started...")

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error reading message:", err)
			continue
		}

		var req request.UserCreateRequest
		if err := json.Unmarshal(msg.Value, &req); err != nil {
			log.Println("Failed to unmarshal Kafka message:", err)
			continue
		}

		// Gọi service xử lý
		_, err = userService.CreateUser(req)
		if err != nil {
			log.Println("Error processing user creation:", err)
		} else {
			fmt.Println("User created successfully from Kafka message")
		}
	}
}
