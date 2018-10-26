package database

import (
	"fmt"
	"log"
	"strings"

	"github.com/expenseledger/web-service/config"
	configdb "github.com/expenseledger/web-service/config/database"
	"github.com/expenseledger/web-service/constant"
	"github.com/jmoiron/sqlx"

	// This is just a PostgreSQL driver for sqlx package
	_ "github.com/lib/pq"
)

var db *sqlx.DB

// Init MUST be called before any package's operations
// @TODO: fix this lame way to initial a package. It's highly depends on
// the order of execution because, somehow, it needs dbinfo.
func init() {
	var (
		dbinfo string
		err    error
	)

	if config.Mode == "PRODUCTION" {
		dbinfo = configdb.DBURL
	} else {
		dbinfo = fmt.Sprintf(
			"user=%s password=%s dbname=%s port=%s sslmode=disable",
			configdb.DBUser,
			configdb.DBPswd,
			configdb.DBName,
			configdb.DBPort,
		)
	}

	db, err = sqlx.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal("Error opening connection to the database", err)
	}
}

// GetDB returns an (probably) initialized instance of sqlx.DB
func GetDB() *sqlx.DB {
	return db
}

// CreateTables creates (if not exists) all the required tables
func CreateTables() (err error) {
	err = createWalletTypeEnum()
	if err != nil {
		log.Println("Error creating enum: wallet_type", err)
		return
	}

	err = createTransactionTypeEnum()
	if err != nil {
		log.Println("Error creating enum: transaction_type", err)
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

	err = createTransactionTable()
	if err != nil {
		log.Println("Error creating table: transaction", err)
		return
	}

	err = createTriggerSetUpdatedAt("wallet", "category", "transaction")
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

func createWalletTypeEnum() (err error) {
	query :=
		fmt.Sprintf(
			"CREATE TYPE wallet_type AS ENUM ('%s', '%s', '%s');",
			constant.WalletType.Cash,
			constant.WalletType.BankAccount,
			constant.WalletType.Credit,
		)
	_, err = db.Exec(query)
	return shouldIgnoreError(err)
}

func createTransactionTypeEnum() (err error) {
	query :=
		fmt.Sprintf(
			"CREATE TYPE transaction_type AS ENUM ('%s', '%s', '%s');",
			constant.TransactionType.Income,
			constant.TransactionType.Expense,
			constant.TransactionType.Transfer,
		)
	_, err = db.Exec(query)
	return shouldIgnoreError(err)
}

func createWalletTable() (err error) {
	query :=
		`
		CREATE TABLE IF NOT EXISTS wallet (
			name character varying(20) NOT NULL PRIMARY KEY,
			type wallet_type NOT NULL,
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
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		CREATE TABLE IF NOT EXISTS transaction (
			id uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
			src_wallet character varying(20) NOT NULL REFERENCES wallet,
			dst_wallet character varying(20) REFERENCES wallet,
			amount NUMERIC(11, 2) DEFAULT 0.00,
			type transaction_type NOT NULL,
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

func shouldIgnoreError(err error) error {
	if err != nil && strings.Contains(err.Error(), "already exists") {
		return nil
	}
	return err
}
