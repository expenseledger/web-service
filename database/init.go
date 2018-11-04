package database

import (
	"fmt"
	"log"
	"strings"

	"github.com/expenseledger/web-service/config"
	dbconfig "github.com/expenseledger/web-service/config/database"
	"github.com/expenseledger/web-service/constant"
	"github.com/jmoiron/sqlx"

	// This is just a PostgreSQL driver for sqlx package
	_ "github.com/lib/pq"
)

// Table names
const (
	ExpenseIncome   = "expense_income"
	Transfer        = "transfer"
	Category        = "category"
	Wallet          = "wallet"
	WalletType      = "wallet_type"
	TransactionType = "transaction_type"
)

var db *sqlx.DB

func init() {
	var (
		dbinfo string
		err    error
	)

	configs := config.GetConfigs()
	dbconfigs := dbconfig.GetConfigs()

	if configs.Mode == "PRODUCTION" {
		dbinfo = dbconfigs.DBURL
	} else {
		dbinfo = fmt.Sprintf(
			"user=%s password=%s dbname=%s port=%s sslmode=disable",
			dbconfigs.DBUser,
			dbconfigs.DBPswd,
			dbconfigs.DBName,
			dbconfigs.DBPort,
		)
	}

	db, err = sqlx.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal("Error opening connection to the database", err)
	}
}

// DB returns an (probably) initialized instance of sqlx.DB
func DB() *sqlx.DB {
	return db
}

// CreateTables creates (if not exists) all the required tables
func CreateTables() (err error) {
	err = createWalletTypeEnum()
	if err != nil {
		log.Println("Error creating enum:", WalletType, err)
		return
	}

	err = createTransactionTypeEnum()
	if err != nil {
		log.Println("Error creating enum:", TransactionType, err)
		return
	}

	err = createWalletTable()
	if err != nil {
		log.Println("Error creating table:", Wallet, err)
		return
	}

	err = createConstantTable(Category)
	if err != nil {
		log.Println("Error creating table:", Category, err)
		return
	}

	err = createExpenseIncomeTable()
	if err != nil {
		log.Println("Error creating table:", ExpenseIncome, err)
		return
	}

	err = createTransferTable()
	if err != nil {
		log.Println("Error creating table:", Transfer, err)
		return
	}

	err = createTriggerSetUpdatedAt(
		Wallet,
		Category,
		ExpenseIncome,
		Transfer,
	)
	if err != nil {
		log.Println("Error creating trigger for updated_at", err)
		return
	}

	return
}

func createConstantTable(tableName string) (err error) {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", tableName)
	query +=
		`
		name character varying(20) PRIMARY KEY,
		created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		`
	_, err = db.Exec(query)
	return
}

func createWalletTypeEnum() (err error) {
	query :=
		fmt.Sprintf(
			"CREATE TYPE %s AS ENUM ('%s', '%s', '%s');",
			WalletType,
			constant.WalletType.Cash,
			constant.WalletType.BankAccount,
			constant.WalletType.Credit,
		)
	_, err = db.Exec(query)
	return filterError(err)
}

func createTransactionTypeEnum() (err error) {
	query :=
		fmt.Sprintf(
			"CREATE TYPE %s AS ENUM ('%s', '%s');",
			TransactionType,
			constant.TransactionType.Income,
			constant.TransactionType.Expense,
		)
	_, err = db.Exec(query)
	return filterError(err)
}

func createWalletTable() (err error) {
	query := fmt.Sprintf(
		`
		CREATE TABLE IF NOT EXISTS %s (
			name character varying(20) PRIMARY KEY,
			type %s NOT NULL,
			balance NUMERIC(11, 2) NOT NULL DEFAULT 0.00,
			created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		`,
		Wallet,
		WalletType,
	)

	_, err = db.Exec(query)
	return
}

func createExpenseIncomeTable() (err error) {
	query := fmt.Sprintf(
		`
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		CREATE TABLE IF NOT EXISTS %s (
			id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
			wallet character varying(20) NOT NULL REFERENCES %s,
			amount NUMERIC(11, 2) NOT NULL DEFAULT 0.00 CHECK (amount >= 0),
			type %s NOT NULL,
			category character varying(20) NOT NULL REFERENCES %s,
			description text NOT NULL DEFAULT '',
			occurred_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		`,
		ExpenseIncome,
		Wallet,
		TransactionType,
		Category,
	)

	_, err = db.Exec(query)
	return
}

func createTransferTable() (err error) {
	query := fmt.Sprintf(
		`
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		CREATE TABLE IF NOT EXISTS %s (
			id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
			src_wallet character varying(20) NOT NULL REFERENCES %s,
			dst_wallet character varying(20) REFERENCES %s,
			amount NUMERIC(11, 2) NOT NULL DEFAULT 0.00 CHECK (amount >= 0),
			description text NOT NULL DEFAULT '',
			occurred_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CHECK (dst_wallet <> src_wallet)
		);
		`,
		Transfer,
		Wallet,
		Wallet,
	)

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

func filterError(err error) error {
	if err != nil && strings.Contains(err.Error(), "already exists") {
		return nil
	}
	return err
}
