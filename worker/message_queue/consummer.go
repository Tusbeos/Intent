package message_queue

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"intent/models"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	Reader *kafka.Reader
}

func NewKafkaConsumer(brokers []string, topic string, groupID string) *KafkaConsumer {
	return &KafkaConsumer{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MaxBytes: 10e6,
		}),
	}
}

func (kc *KafkaConsumer) Start() {
	log.Println("[KafkaConsumer] Starting consumer...")

	for {
		msg, err := kc.Reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("[KafkaConsumer] Error reading message:", err)
			continue
		}

		log.Println("[KafkaConsumer] Received message:", string(msg.Value))

		var action models.UserActionMessage
		if err := json.Unmarshal(msg.Value, &action); err != nil {
			log.Println("[KafkaConsumer] JSON parse error:", err)
			continue
		}

		log.Printf("[KafkaConsumer] Processing message: UserID=%d, Action=%s, Timestamp=%s\n",
			action.UserID, action.Action, time.Now().Format(time.RFC3339))

		// Giả lập thành công
		success := true
		if success {
			log.Println("[KafkaConsumer] Successfully processed action:", action)
		} else {
			log.Println("[KafkaConsumer] Failed to process action:", action)
		}
	}
}
