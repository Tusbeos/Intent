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

// Biến Redis Client
var RedisClient *redis.Client

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

// Hàm lưu cache vào Redis
func SetCache(key string, value string) error {
	ctx := context.Background()
	return RedisClient.Set(ctx, key, value, 0).Err()
}

// Hàm lấy dữ liệu từ Redis
func GetCache(key string) (string, error) {
	ctx := context.Background()
	return RedisClient.Get(ctx, key).Result()
}

// Hàm xóa cache
func DeleteCache(key string) error {
	ctx := context.Background()
	return RedisClient.Del(ctx, key).Err()
}
