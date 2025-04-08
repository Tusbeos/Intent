package cmd

import (
	"log"

	"intent/config"
	"intent/repository"
	"intent/worker/log_action"
	"intent/worker/message_queue"
)

func StartWorker() {
	log.Println("Loading configuration...")
	cfg := config.LoadConfig()
	db, err := config.ConnectDB(cfg, "Worker")

	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	redisClient := config.ConnectRedis(cfg, "Worker")
	if redisClient == nil {
		log.Fatal("Failed to connect to Redis")
	}

	userRepo := repository.NewUserRepository(db, redisClient)

	redisWorker := log_action.NewRedisWorker(redisClient, userRepo)
	go redisWorker.Start()

	processor := message_queue.NewProcessor(userRepo, 3)
	kafkaBroker, topic := config.GetKafkaConfig(cfg)
	kafkaWorker := message_queue.NewKafkaWorker(kafkaBroker, topic, processor)
	go kafkaWorker.Start()

	log.Println("Worker started and waiting for messages...")
	select {}
}
