package database

import (
	"os"

	"github.com/expenseledger/web-service/config"
)

var (
	DBUser string
	DBPswd string
	DBName string
	DBPort string
	DBURL  string
)

func init() {
	if config.Mode == "PRODUCTION" {
		DBURL = os.Getenv("DATABASE_URL")
	} else {
		DBUser = os.Getenv("DB_USER")
		DBPswd = os.Getenv("DB_PASSWORD")
		DBName = os.Getenv("DB_NAME")
		DBPort = os.Getenv("DB_PORT")
	}
}
