package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	Mode string
	Port string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	Mode = os.Getenv("MODE")
	Port = os.Getenv("PORT")

	if Port == "" {
		Port = "3000"
	}
}
