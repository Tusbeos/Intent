package message_queue

import (
	"context"
	"encoding/json"
	"log"

	"intent/repository"

	"github.com/segmentio/kafka-go"
)

type KafkaWorker struct {
	broker   string
	topic    string
	userRepo *repository.UserRepository
	retries  int
}

func NewKafkaWorker(broker, topic string, userRepo *repository.UserRepository, retries int) *KafkaWorker {
	return &KafkaWorker{
		broker:   broker,
		topic:    topic,
		userRepo: userRepo,
		retries:  retries,
	}
}

func (w *KafkaWorker) Start() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{w.broker},
		Topic:   w.topic,
		GroupID: "user-actions",
	})
	defer r.Close()

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("[Kafka] Error reading message:", err)
			continue
		}

		// Parse message
		var data map[string]interface{}
		err = json.Unmarshal(msg.Value, &data)
		if err != nil {
			log.Println("[Kafka] Failed to unmarshal message:", err)
			continue
		}

		// Lấy request_id để logging
		requestID, _ := data["request_id"].(string)
		log.Printf("[Kafka] Processing message with request_id: %s", requestID)

		// Giả lập xử lý dữ liệu
		if err := w.processData(data); err != nil {
			log.Printf("[Kafka] Failed to process message with request_id: %s, error: %v", requestID, err)
			continue
		}

		// Xử lý thành công
		log.Printf("[Kafka] Successfully processed message with request_id: %s", requestID)
	}
}

// processData giả lập xử lý dữ liệu (sau này có thể thay bằng lưu DB hoặc xử lý khác)
func (w *KafkaWorker) processData(data map[string]interface{}) error {
	// Giả lập xử lý logic (ở đây chỉ return nil, sau này có thể thêm logic khác)
	return nil
}
