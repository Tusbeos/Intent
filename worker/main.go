package main

import (
	"log"

	"intent/config"
	"intent/repository"
	"intent/worker/log_action"
	"intent/worker/message_queue"
)

func main() {
	log.Println("-------------------------------++-------------------------------")
	log.Println("[Worker] Loading configuration...")

	// Load config
	cfg := config.LoadConfig()

	// Kết nối DB
	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("[Worker] Database connection failed: %v", err)
	}

	// Kết nối Redis
	redisClient := config.ConnectRedis(cfg)
	if redisClient == nil {
		log.Fatal("[Worker] Failed to connect to Redis")
	}

	// Tạo instance repository
	userRepo := repository.NewUserRepository(db, redisClient)

	// Khởi động worker Redis (lưu log action)
	redisWorker := log_action.NewRedisWorker(redisClient, userRepo)
	go func() {
		log.Println("[Worker] Starting Redis worker...")
		redisWorker.Start()
	}()
	log.Println("[Worker] Redis worker started successfully.")

	// Khởi động worker Kafka (message queue)
	kafkaBroker, topic := config.GetKafkaConfig(cfg)
	kafkaWorker := message_queue.NewKafkaWorker(kafkaBroker, topic, userRepo, 3)

	go func() {
		log.Println("[Worker] Starting Kafka worker...")
		kafkaWorker.Start()
	}()
	log.Println("[Worker] Kafka worker started, waiting for messages...")
	select {}
}
