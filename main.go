package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"Intent/config"
	"Intent/controller"
)

func connectDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	log.Println("Successfully connected to MySQL!")
	return db, nil
}

func main() {
	cfg := config.LoadConfig()
	db, err := connectDB(cfg)
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
