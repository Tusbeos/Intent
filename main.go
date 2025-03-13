package main

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"

	"intent/config"
	"intent/controller"
	"intent/middleware"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Kết nối Database
	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	// Kết nối Redis
	redisClient := config.ConnectRedis(cfg)

	// Khởi tạo Echo server
	e := echo.New()

	cacheableRoutes := map[string]bool{
		"/users":     true,
		"/users/:id": true,
	}

	// Thêm middleware Redis Cache
	e.Use(middleware.RedisCache(redisClient, cacheableRoutes))

	// Thêm middleware Rate Limit
	e.Use(middleware.RateLimitMiddleware(redisClient, 10, 30*time.Second))

	// Đăng ký API routes
	controller.RegisterUserRoutes(e, db, redisClient)

	// Chạy server
	port := ":8080"
	log.Printf("Server is running on port %s...", port)
	e.Logger.Fatal(e.Start(port))
}
