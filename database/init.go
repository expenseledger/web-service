package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

// Init MUST be called before any package's operations
// @TODO: fix this lame way to initial a package. It's highly depends on
// the order of execution because, somehow, it needs dbinfo.
func Init(dbinfo string) (err error) {
	db, err = sqlx.Open("postgres", dbinfo)
	if err != nil {
		log.Println("Error opening connection to the database", err)
	}
	return
}

// GetDB returns an (probably) initialized instance of sqlx.DB
func GetDB() *sqlx.DB {
	return db
}

// CreateTables creates (if not exists) all the required tables
func CreateTables() (err error) {
	err = createConstantTable("wallet_type")
	if err != nil {
		log.Println("Error creating table: wallet_type", err)
		return
	}

	err = createWalletTable()
	if err != nil {
		log.Println("Error creating table: wallet", err)
		return
	}

	err = createConstantTable("category")
	if err != nil {
		log.Println("Error creating table: category", err)
		return
	}

	err = createConstantTable("transaction_type")
	if err != nil {
		log.Println("Error creating table: transaction_type", err)
		return
	}

	err = createTransactionTable()
	if err != nil {
		log.Println("Error creating table: transaction", err)
		return
	}

	err = createTriggerSetUpdatedAt(
		"wallet_type",
		"wallet",
		"category",
		"transaction_type",
		"transaction",
	)
	if err != nil {
		log.Println("Error creating trigger for updated_at", err)
		return
	}

	return
}

func createConstantTable(tableName string) (err error) {
	query :=
		fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", tableName) +
			`name character varying(20) NOT NULL PRIMARY KEY,
			created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			deleted_at timestamp with time zone);`
	_, err = db.Exec(query)
	return
}

func createWalletTable() (err error) {
	query :=
		`
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		CREATE TABLE IF NOT EXISTS wallet (
			id uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
			name character varying(20) NOT NULL,
			type character varying(20) NOT NULL REFERENCES wallet_type,
			balance NUMERIC(11, 2) DEFAULT 0.00,
			created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			deleted_at timestamp with time zone
		);
		`

	_, err = db.Exec(query)
	return
}

func createTransactionTable() (err error) {
	query :=
		`
		CREATE TABLE IF NOT EXISTS transaction (
			id uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
			src_wallet uuid NOT NULL REFERENCES wallet,
			dst_wallet uuid REFERENCES wallet,
			amount NUMERIC(11, 2) DEFAULT 0.00,
			type character varying(20) NOT NULL REFERENCES transaction_type,
			category character varying(20) NOT NULL REFERENCES category,
			description text,
			created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			deleted_at timestamp with time zone,
			CHECK (dst_wallet <> src_wallet)
		);
		`

	_, err = db.Exec(query)
	return
}

func createTriggerSetUpdatedAt(tableNames ...string) (err error) {
	query := deleteExistingTriggers(tableNames)
	query += "CREATE EXTENSION IF NOT EXISTS moddatetime;"

	for _, tableName := range tableNames {
		query += fmt.Sprintf(
			`
			CREATE TRIGGER %s
			BEFORE UPDATE ON %s
			FOR EACH ROW
			EXECUTE PROCEDURE moddatetime (updated_at);
			`,
			"mdt_"+tableName,
			tableName,
		)
	}

	_, err = db.Exec(query)
	return
}

func deleteExistingTriggers(tableNames []string) string {
	var query string

	for _, tableName := range tableNames {
		query += fmt.Sprintf(
			"DROP TRIGGER IF EXISTS %s ON %s;",
			"mdt_"+tableName,
			tableName,
		)
	}

	return query
}
