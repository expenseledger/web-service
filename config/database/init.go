package database

import (
	"os"

	"github.com/expenseledger/web-service/config"
)

type configFields struct {
	DBUser string
	DBPswd string
	DBName string
	DBPort string
	DBURL  string
}

var configs configFields

func init() {
	gConfigs := config.GetConfigs()

	if gConfigs.Mode == "PRODUCTION" {
		configs.DBURL = os.Getenv("DATABASE_URL")
	} else {
		configs.DBUser = os.Getenv("DB_USER")
		configs.DBPswd = os.Getenv("DB_PASSWORD")
		configs.DBName = os.Getenv("DB_NAME")
		configs.DBPort = os.Getenv("DB_PORT")
	}
}

// GetConfigs ...
func GetConfigs() configFields {
	return configs
}
