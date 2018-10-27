package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type configFields struct {
	Mode string
	Port string
}

var configs configFields

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	configs.Mode = os.Getenv("MODE")
	configs.Port = os.Getenv("PORT")

	if configs.Port == "" {
		configs.Port = "3000"
	}
}

// GetConfigs ...
func GetConfigs() configFields {
	return configs
}
