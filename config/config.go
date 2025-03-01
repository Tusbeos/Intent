package config

import (
	"encoding/json"
	"log"
	"os"
)

// Khai cau truc Config
type Config struct {
	Database struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		DBName   string `json:"dbname"`
	} `json:"database"`
}

func LoadConfig() *Config {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("Unable to open config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	log.Println("Config loaded successfully!")
	return config
}
