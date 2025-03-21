package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"intent/config"
	"intent/models"
	"intent/repository"
)

func main() {
	// Load cấu hình & kết nối DB
	cfg := config.LoadConfig()
	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Kết nối Redis
	redisClient := config.ConnectRedis(cfg)
	if redisClient == nil {
		log.Fatal("Failed to connect to Redis")
	}

	// Tạo instance repository
	userRepo := repository.NewUserRepository(db, redisClient)

	ctx := context.Background()
	log.Println("Worker is running...")

	for {
		// Lấy message từ Redis queue
		msg, err := redisClient.RPop(ctx, "user_action_queue").Result()
		if err != nil {
			if err == redis.Nil {
				time.Sleep(2 * time.Second) // Đợi nếu queue rỗng
				continue
			}
			log.Println("Redis queue error:", err)
			continue
		}

		// Parse JSON
		var userAction models.UserActionMessage
		if err := json.Unmarshal([]byte(msg), &userAction); err != nil {
			log.Println("JSON parse error:", err)
			continue
		}

		// Chuyển `Timestamp` từ string sang `time.Time`
		parsedTime, err := time.Parse(time.RFC3339, userAction.Timestamp)
		if err != nil {
			log.Println("Timestamp parse error:", err)
			continue
		}

		// Lưu log vào DB
		err = userRepo.SaveLogAction(models.LogAction{
			UserID:    userAction.UserID,
			Action:    userAction.Action,
			Timestamp: parsedTime,
		})
		if err != nil {
			log.Println("DB save error:", err)
		}
	}
}
