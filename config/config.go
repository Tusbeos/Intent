package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"intent/models"
)

// Các  Biến
var RedisClient *redis.Client
var DB *gorm.DB

// Kết nối MySQL
func ConnectDB(cfg *models.Config) (*gorm.DB, error) {
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

// Load config từ file JSON
func LoadConfig() *models.Config {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
	}
	defer file.Close()

	config := &models.Config{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	return config
}

// Hàm kết nối Redis
func ConnectRedis(cfg *models.Config) *redis.Client {
	redisAddr := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Kiểm tra kết nối Redis
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis connection error: %v", err)
		return nil
	}

	log.Println("Successfully connected to Redis!")
	RedisClient = client
	return client
}
