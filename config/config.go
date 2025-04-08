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

var RedisClient *redis.Client
var DB *gorm.DB

func ConnectDB(cfg *models.Config, caller string) (*gorm.DB, error) {
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
	log.Printf("%s Successfully connected to MySQL!", caller)
	return db, nil
}

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

func ConnectRedis(cfg *models.Config, caller string) *redis.Client {
	redisAddr := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("%s Redis connection error: %v", caller, err)
		return nil
	}

	log.Printf("[%s] Successfully connected to Redis!", caller)
	return client
}

func GetKafkaConfig(cfg *models.Config) (string, string) {
	return cfg.Kafka.Brokers[0], cfg.Kafka.Topic
}
