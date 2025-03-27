package models

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

// Cấu hình Redis
type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// Cấu hình Kafka
type KafkaConfig struct {
	Brokers []string `json:"brokers"`
	Topic   string   `json:"topic"`
}

// Cấu hình tổng hợp
type Config struct {
	Database DatabaseConfig `json:"database"`
	Redis    RedisConfig    `json:"redis"`
	Kafka    KafkaConfig    `json:"kafka"`
}
