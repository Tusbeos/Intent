package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		User     string
		Password string
		Host     string
		Port     int
		DBName   string
	}
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	return &cfg
}
