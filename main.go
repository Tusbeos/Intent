package main

import (
	"log"

	"github.com/labstack/echo/v4"

	"Http_Management/config"
	"Http_Management/controller"
)

func main() {
	cfg := config.LoadConfig()
	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	// Khởi tạo Echo server
	e := echo.New()
	controller.RegisterUserRoutes(e, db)

	// Chạy server
	port := ":8080"
	log.Printf("Server is running on port %s...", port)
	e.Logger.Fatal(e.Start(port))

}
