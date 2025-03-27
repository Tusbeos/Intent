package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"intent/config"
	"intent/models"
)

func PublishUserAction(userID int, action string) {
	msg := models.UserActionMessage{
		UserID:    userID,
		Action:    action,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("JSON marshal error:", err)
		return
	}

	if err := config.RedisClient.LPush(context.Background(), "user_action_queue", data).Err(); err != nil {
		log.Println("Redis push error:", err)
	}
}
