package cmd

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"

	"intent/config"
	"intent/controller"
	"intent/kafka"
	"intent/middleware"
)

func StartServer() {
	cfg := config.LoadConfig()
	db, err := config.ConnectDB(cfg, "Server")
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	redisClient := config.ConnectRedis(cfg, "Server")

	e := echo.New()
	e.Use(middleware.RequestIDMiddleware())
	e.Use(middleware.RateLimitMiddleware(redisClient, 50, 30*time.Second))
	kafkaProducer := kafka.NewProducer(cfg.Kafka.Brokers[0], cfg.Kafka.Topic)
	controller.RegisterUserRoutes(e, db, redisClient, kafkaProducer)

	port := ":8080"
	log.Printf("Server is running on port %s...", port)
	e.Logger.Fatal(e.Start(port))
}
