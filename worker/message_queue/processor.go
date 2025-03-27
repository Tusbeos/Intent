package message_queue

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"intent/models"
)

type KafkaProcessor struct{}

func NewKafkaProcessor() *KafkaProcessor {
	return &KafkaProcessor{}
}

func (p *KafkaProcessor) ProcessMessage(data []byte) error {
	var action models.UserActionMessage
	if err := json.Unmarshal(data, &action); err != nil {
		log.Println("[KafkaProcessor] JSON parse error:", err)
		return err
	}

	// Kiểm tra dữ liệu hợp lệ
	if action.UserID == 0 || action.Action == "" {
		log.Println("[KafkaProcessor] Invalid action data:", action)
		return errors.New("invalid action data")
	}

	// Giả lập xử lý thành công
	log.Printf("[KafkaProcessor] Processing message: UserID=%d, Action=%s, Timestamp=%s\n",
		action.UserID, action.Action, time.Now().Format(time.RFC3339))

	log.Println("[KafkaProcessor] Successfully processed action:", action)
	return nil
}
