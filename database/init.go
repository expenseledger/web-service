package database

import (
	"database/sql"
	"log"
)

// InitTables create (if not exists) tables
func InitTables(db *sql.DB) {
	err := createWalletType(db)
	if err != nil {
		log.Println("Error creating table: wallet_type", err)
	}
}

func createWalletType(db *sql.DB) (err error) {
	query :=
		`CREATE TABLE IF NOT EXISTS wallet_type (
			name character varying(20) NOT NULL PRIMARY KEY,
			created_at timestamp with time zone,
			updated_at timestamp with time zone,
			deleted_at timestamp with time zone
		);`

	_, err = db.Exec(query)
	return
}
