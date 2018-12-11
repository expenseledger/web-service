package orm

import (
	"fmt"
	"log"
	"time"

	"github.com/expenseledger/web-service/db"
	"github.com/shopspring/decimal"
)

// Wallet the structure represents a stored wallet on db
type Wallet struct {
	Name      string          `db:"name"`
	Type      string          `db:"type"`
	Balance   decimal.Decimal `db:"balance"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
}

// Wallets is defined just to be used as a receiver
type Wallets []Wallet

// Insert ...
func (wallet *Wallet) Insert() error {
	query :=
		`
		INSERT INTO wallet (name, type, balance)
		VALUES (:name, :type, :balance)
		RETURNING *;
		`

	stmt, err := db.Conn().PrepareNamed(query)
	if err != nil {
		log.Println("Error inserting a wallet", err)
		return err
	}

	if err := stmt.Get(wallet, wallet); err != nil {
		log.Println("Error inserting a wallet", err)
		return err
	}

	return nil
}

// One ...
func (wallet *Wallet) One() error {
	query :=
		`
		SELECT * FROM wallet
		WHERE name=:name;
		`

	namedStmt, err := db.Conn().PrepareNamed(query)
	if err != nil {
		log.Println("Error selecting a wallet", err)
		return err
	}

	if err := namedStmt.Get(wallet, wallet); err != nil {
		log.Println("Error selecting a wallet", err)
		return err
	}

	return nil
}

// Delete ...
func (wallet *Wallet) Delete(name string) error {
	query :=
		`
		DELETE FROM wallet
		WHERE name=$1
		RETURNING *;
		`

	stmt, err := db.Conn().Preparex(query)
	if err != nil {
		log.Println("Error deleting a wallet", err)
		return err
	}

	if err := stmt.Get(wallet, name); err != nil {
		log.Println("Error deleting a wallet", err)
		return err
	}

	return nil
}

// All gets all wallets
func (wallets *Wallets) All() (int, error) {
	query :=
		`
		SELECT * FROM wallet;
		`

	stmt, err := db.Conn().Preparex(query)
	if err != nil {
		log.Println("Error selecting all wallets", err)
		return 0, err
	}

	if err := stmt.Select(wallets); err != nil {
		log.Println("Error selecting all wallets", err)
		return 0, err
	}

	return len(*wallets), nil
}

// BatchInsert ...
func (wallets *Wallets) BatchInsert() (int, error) {
	var err error

	length := len(*wallets)
	insertedWallets := make(Wallets, length)
	for i, wallet := range *wallets {
		err = wallet.Insert()
		if err != nil {
			break
		}
		insertedWallets[i] = wallet
	}

	if err != nil {
		log.Println("Error doing batch insertion wallets", err)
		return 0, err
	}

	*wallets = insertedWallets
	return length, nil
}

// DeleteAll ...
func (wallets *Wallets) DeleteAll() (int, error) {
	query :=
		`
		DELETE FROM wallet
		RETURNING *;
		`

	stmt, err := db.Conn().Preparex(query)
	if err != nil {
		log.Println("Error deleting all wallets", err)
		return 0, err
	}

	if err := stmt.Select(wallets); err != nil {
		log.Println("Error deleting all wallets", err)
		return 0, err
	}

	return len(*wallets), nil
}

// Save replaces old fields with new ones
func (wallet *Wallet) Save() error {
	query :=
		`
		UPDATE wallet
		SET (balance, type)=(:balance, :type)
		WHERE name=:name
		RETURNING *;
		`

	namedStmt, err := db.Conn().PrepareNamed(query)
	if err != nil {
		log.Println("Error saving a wallet", err)
		return err
	}

	if err := namedStmt.Get(wallet, wallet); err != nil {
		log.Println("Error saving a wallet", err)
		return err
	}

	return nil
}

// Update updates only non-zero value
func (wallet *Wallet) Update() error {
	fields, names := wallet.buildUpdateSQLStmt()
	query := fmt.Sprintf(
		`
		UPDATE wallet
		SET %s=%s
		WHERE name=:name
		RETURNING *;
		`,
		fields,
		names,
	)

	namedStmt, err := db.Conn().PrepareNamed(query)
	if err != nil {
		log.Println("Error updating a wallet", err)
		return err
	}

	if err := namedStmt.Get(wallet, wallet); err != nil {
		log.Println("Error updating a wallet", err)
	}

	return nil
}

func (wallet *Wallet) buildUpdateSQLStmt() (fields, names string) {
	fields, names = "(", "("

	if !wallet.Balance.IsZero() {
		fields += "balance"
		names += ":balance"
	}

	if wallet.Type != "" {
		fields += ", type"
		names += ", :type"
	}

	fields += ")"
	names += ")"

	return
}
