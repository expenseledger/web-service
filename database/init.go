package database

import (
	"database/sql"
	"log"
)

// InitTables create (if not exists) tables
func InitTables(db *sql.DB) {
	var err error

	err = createWalletTypeTable(db)
	if err != nil {
		log.Println("Error creating table: wallet_type", err)
	}

	err = createWalletTable(db)
	if err != nil {
		log.Println("Error creating table: wallet", err)
	}

	err = createCategoryTable(db)
	if err != nil {
		log.Println("Error creating table: category", err)
	}
}

func createWalletTypeTable(db *sql.DB) (err error) {
	query :=
		`
		CREATE TABLE IF NOT EXISTS wallet_type (
			name character varying(20) NOT NULL PRIMARY KEY,
			created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			deleted_at timestamp with time zone
		);
		`

	_, err = db.Exec(query)
	return
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

func createCategoryTable(db *sql.DB) (err error) {
	query :=
		`
		CREATE TABLE IF NOT EXISTS category (
			name character varying(20) NOT NULL PRIMARY KEY,
			created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			deleted_at timestamp with time zone
		);
		`

	_, err = db.Exec(query)
	return
}
