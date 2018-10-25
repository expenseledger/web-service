package model

import (
	"github.com/expenseledger/web-service/database"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func init() {
	db = database.GetDB()
}
