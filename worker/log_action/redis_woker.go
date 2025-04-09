package log_action

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"intent/models"
	"intent/repository"
)

type RedisWorker struct {
	RedisClient *redis.Client
	UserRepo    *repository.UserRepository
}

func NewRedisWorker(redisClient *redis.Client, userRepo *repository.UserRepository) *RedisWorker {
	return &RedisWorker{RedisClient: redisClient, UserRepo: userRepo}
}

func (w *RedisWorker) Start() {
	ctx := context.Background()
	log.Println("Log Action Worker is running...")

	for {
		msg, err := w.RedisClient.RPop(ctx, "user_action_queue").Result()
		if err != nil {
			if err == redis.Nil {
				time.Sleep(2 * time.Second)
				continue
			}
			log.Println("Redis queue error:", err)
			continue
		}
		w.processMessage(msg)
	}
}

func (w *RedisWorker) processMessage(msg string) {
	var userAction models.UserActionMessage
	if err := json.Unmarshal([]byte(msg), &userAction); err != nil {
		log.Println("JSON parse error:", err)
		return
	}

	parsedTime, err := time.Parse(time.RFC3339, userAction.Timestamp)
	if err != nil {
		log.Println("Timestamp parse error:", err)
		return
	}

	err = w.UserRepo.SaveLogAction(models.LogAction{
		UserID:    userAction.UserID,
		Action:    userAction.Action,
		Timestamp: parsedTime,
	})
	if err != nil {
		log.Println("DB save error:", err)
	}
}
