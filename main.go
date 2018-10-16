package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/ExpenseLedger/expense-ledger-web-service/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file", err)
	}

	var dbinfo string
	if os.Getenv("MODE") == "DEVELOPMENT" {
		dbinfo = fmt.Sprintf(
			"user=%s password=%s dbname=%s port=%s sslmode=disable",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
		)
	} else {
		dbinfo = os.Getenv("DATABASE_URL")
	}

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal("Error opening connection to the database", err)
	}

	err = database.CreateTables(db)
	if err != nil {
		log.Println("Error creating tables")
	} else {
		log.Println("Successfully created tables")
	}
}
