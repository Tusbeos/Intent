package controller

import (
	"context"
	"encoding/json"
	"log"
	"time"
)

// pushToQueue đẩy message vào Redis queue (ghi log)
func (uc *UserController) pushToLog(userID int, action string) {
	msg, err := json.Marshal(map[string]interface{}{
		"user_id":   userID,
		"action":    action,
		"timestamp": time.Now().Format(time.RFC3339),
	})
	if err != nil {
		log.Println("Failed to marshal message:", err)
		return
	}

	if err := uc.RedisClient.RPush(context.Background(), "user_action_queue", msg).Err(); err != nil {
		log.Println("Failed to push message to Redis:", err)
		return
	}

	log.Printf("Sent log action to Redis for user_id: %d with action: %s", userID, action)
}
