package database

import (
	"database/sql"
	"fmt"
	"log"
)

// InitTables create (if not exists) tables
func InitTables(db *sql.DB) {
	var err error

	err = createConstantTable(db, "wallet_type")
	if err != nil {
		log.Println("Error creating table: wallet_type", err)
	}

	err = createWalletTable(db)
	if err != nil {
		log.Println("Error creating table: wallet", err)
	}

	err = createConstantTable(db, "category")
	if err != nil {
		log.Println("Error creating table: category", err)
	}

	err = createConstantTable(db, "transaction_type")
	if err != nil {
		log.Println("Error creating table: transaction_type", err)
	}
}

func createWalletTable(db *sql.DB) (err error) {
	query :=
		`
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		CREATE TABLE IF NOT EXISTS wallet (
			id uuid NOT NULL DEFAULT uuid_generate_v4(),
			name character varying(20) NOT NULL,
			type character varying(20) REFERENCES wallet_type,
			balance NUMERIC(11, 2) DEFAULT 0.00,
			created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			deleted_at timestamp with time zone
		);
		`

	_, err = db.Exec(query)
	return
}

func createConstantTable(db *sql.DB, tableName string) (err error) {
	query :=
		fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", tableName) +
			`name character varying(20) NOT NULL PRIMARY KEY,
			created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			deleted_at timestamp with time zone);`
	_, err = db.Exec(query)
	return
}
