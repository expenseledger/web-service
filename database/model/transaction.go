package model

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/expenseledger/web-service/constant"
	"github.com/expenseledger/web-service/database"
	"github.com/shopspring/decimal"
)

// Transaction the structure represents a stored transaction on database
type Transaction struct {
	ID          string          `db:"id"`
	SrcWallet   string          `db:"src_wallet"`
	DstWallet   string          `db:"dst_wallet"`
	Amount      decimal.Decimal `db:"amount"`
	Type        string          `db:"type"`
	Category    string          `db:"category"`
	Description string          `db:"description"`
	OccurredAt  time.Time       `db:"occurred_at"`
}

type transaction struct {
	ID          string          `db:"id"`
	Wallet      string          `db:"wallet"`
	WalletRole  string          `db:"role"`
	Amount      decimal.Decimal `db:"amount"`
	Type        string          `db:"type"`
	Category    string          `db:"category"`
	Description string          `db:"description"`
	OccurredAt  time.Time       `db:"occurred_at"`
}

type affectedWallet struct {
	TransactionID string    `db:"transaction_id"`
	Wallet        string    `db:"wallet"`
	Role          string    `db:"role"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

// Transactions is defined just to be used as a receiver
type Transactions []Transaction

// Insert ...
func (tx *Transaction) Insert() error {
	query := tx.buildInsertSQLStmt()

	stmt, err := database.DB().Preparex(query)
	if err != nil {
		log.Println("Error inserting a transaction", err)
		return err
	}

	walletRole := constant.WalletRole()
	transactionType := constant.TransactionType()
	switch {
	case !tx.OccurredAt.IsZero() &&
		tx.Type == transactionType.Transfer:
		err = stmt.Get(
			tx,
			tx.Amount,
			tx.Type,
			tx.Category,
			tx.Description,
			tx.OccurredAt,
			tx.SrcWallet,
			walletRole.SrcWallet,
			tx.DstWallet,
			walletRole.DstWallet,
		)

	case !tx.OccurredAt.IsZero():
		wallet, role := tx.SrcWallet, walletRole.SrcWallet
		if tx.Type == transactionType.Income {
			wallet, role = tx.DstWallet, walletRole.DstWallet
		}
		err = stmt.Get(
			tx,
			tx.Amount,
			tx.Type,
			tx.Category,
			tx.Description,
			tx.OccurredAt,
			wallet,
			role,
		)

	case tx.Type == constant.TransactionType().Transfer:
		err = stmt.Get(
			tx,
			tx.Amount,
			tx.Type,
			tx.Category,
			tx.Description,
			tx.SrcWallet,
			walletRole.SrcWallet,
			tx.DstWallet,
			walletRole.DstWallet,
		)

	default:
		wallet, role := tx.SrcWallet, walletRole.SrcWallet
		if tx.Type == transactionType.Income {
			wallet, role = tx.DstWallet, walletRole.DstWallet
		}
		err = stmt.Get(
			tx,
			tx.Amount,
			tx.Type,
			tx.Category,
			tx.Description,
			wallet,
			role,
		)
	}

	if err != nil {
		log.Println("Error inserting a transaction", err)
		return err
	}
	return nil
}

// DeleteAll ...
func (txs *Transactions) DeleteAll() error {
	query := fmt.Sprintf(
		`
		WITH deleted_tx AS (
			DELETE FROM %s
			RETURNING *
		), deleted_tx_wallet AS (
			DELETE FROM %s
			RETURNING *
		)
		SELECT
			t.id AS id,
			w.wallet AS wallet,
			w.role AS role,
			t.amount AS amount,
			t.type AS type,
			t.category AS category,
			t.description AS description,
			t.occurred_at AS occurred_at
			FROM deleted_tx t, deleted_tx_wallet w
			WHERE t.id=w.transaction_id;
		`,
		database.Transaction,
		database.AffectedWallet,
	)

	stmt, err := database.DB().Preparex(query)
	if err != nil {
		log.Println("Error deleting all transactions", err)
		return err
	}

	var _txs []transaction
	if err := stmt.Select(&_txs); err != nil {
		log.Println("Error deleting all transactions", err)
		return err
	}

	role := constant.WalletRole()
	mapTx := make(map[string]*Transaction)
	for _, _tx := range _txs {
		tx, ok := mapTx[_tx.ID]
		if !ok {
			tx = _tx.toExtTx()
			mapTx[_tx.ID] = tx
		}

		if _tx.WalletRole == role.SrcWallet {
			tx.SrcWallet = _tx.Wallet
		} else {
			tx.DstWallet = _tx.Wallet
		}
	}

	transactions := make(Transactions, len(mapTx))
	i := 0
	for _, tx := range mapTx {
		transactions[i] = *tx
		i++
	}

	*txs = transactions

	return nil
}

// One ...
func (tx *Transaction) One() error {
	query := fmt.Sprintf(
		`
		SELECT
			t.id AS id,
			w.wallet AS wallet,
			w.role AS role,
			t.amount AS amount,
			t.type AS type,
			t.category AS category,
			t.description AS description,
			t.occurred_at AS occurred_at
			FROM %s t, %s w
			WHERE t.id=$1 AND t.id=w.transaction_id;
		`,
		database.Transaction,
		database.AffectedWallet,
	)

	stmt, err := database.DB().Preparex(query)
	if err != nil {
		log.Println("Error selecting a transactions", err)
		return err
	}

	var _txs []transaction
	if err := stmt.Select(&_txs, tx.ID); err != nil {
		log.Println("Error selecting a transactions", err)
		return err
	}

	if len(_txs) < 1 {
		log.Println("Error selecting a transactions record not found")
		return errors.New("record not found")
	}

	*tx = *_txs[0].toExtTx()
	if len(_txs) > 1 {
		role := constant.WalletRole()
		if _tx := _txs[1]; _tx.WalletRole == role.SrcWallet {
			tx.SrcWallet = _tx.Wallet
		} else {
			tx.DstWallet = _tx.Wallet
		}
	}

	return nil
}

// Delete ...
func (tx *Transaction) Delete() error {
	query := fmt.Sprintf(
		`
		WITH deleted_aw AS (
			DELETE FROM %s
			WHERE transaction_id=$1
			RETURNING *
		), deleted_tx AS (
			DELETE FROM %s
			WHERE id IN (SELECT transaction_id FROM deleted_aw)
			RETURNING *
		)
		SELECT
			t.id AS id,
			w.wallet AS wallet,
			w.role AS role,
			t.amount AS amount,
			t.type AS type,
			t.category AS category,
			t.description AS description,
			t.occurred_at AS occurred_at
			FROM deleted_tx t, deleted_aw w
			WHERE t.id=w.transaction_id;
		`,
		database.AffectedWallet,
		database.Transaction,
	)

	stmt, err := database.DB().Preparex(query)
	if err != nil {
		log.Println("Error deleting a transactions", err)
		return err
	}

	var _txs []transaction
	if err := stmt.Select(&_txs, tx.ID); err != nil {
		log.Println("Error deleting a transactions", err)
		return err
	}

	if len(_txs) < 1 {
		log.Println("Error selecting a transactions record not found")
		return errors.New("record not found")
	}

	*tx = *_txs[0].toExtTx()
	if len(_txs) > 1 {
		role := constant.WalletRole()
		if _tx := _txs[1]; _tx.WalletRole == role.SrcWallet {
			tx.SrcWallet = _tx.Wallet
		} else {
			tx.DstWallet = _tx.Wallet
		}
	}

	return nil
}

func (tx *Transaction) buildInsertSQLStmt() string {
	var query string

	switch {
	case !tx.OccurredAt.IsZero() &&
		tx.Type == constant.TransactionType().Transfer:
		query = fmt.Sprintf(
			`
			WITH inserted_tx AS (
				INSERT INTO %s
				(amount, type, category, description, occurred_at)
				VALUES
				($1, $2, $3, $4, $5)
				RETURNING *
			), tx_wallet AS (
				INSERT INTO %s
				(transaction_id, wallet, role)
				SELECT id, $6, CAST ($7 AS %s) FROM inserted_tx
				UNION ALL
				SELECT id, $8, CAST ($9 AS %s) FROM inserted_tx
				RETURNING *
			)
			SELECT id, occurred_at FROM inserted_tx;
			`,
			database.Transaction,
			database.AffectedWallet,
			database.WalletRole,
			database.WalletRole,
		)

	case !tx.OccurredAt.IsZero():
		query = fmt.Sprintf(
			`
			WITH inserted_tx AS (
				INSERT INTO %s
				(amount, type, category, description, occurred_at)
				VALUES
				($1, $2, $3, $4, $5)
				RETURNING *
			), tx_wallet AS (
				INSERT INTO %s
				(transaction_id, wallet, role)
				SELECT id, $6, CAST ($7 AS %s) FROM inserted_tx
				RETURNING *
			)
			SELECT id, occurred_at FROM inserted_tx;
			`,
			database.Transaction,
			database.AffectedWallet,
			database.WalletRole,
		)

	case tx.Type == constant.TransactionType().Transfer:
		query = fmt.Sprintf(
			`
			WITH inserted_tx AS (
				INSERT INTO %s
				(amount, type, category, description)
				VALUES
				($1, $2, $3, $4)
				RETURNING *
			), tx_wallet AS (
				INSERT INTO %s
				(transaction_id, wallet, role)
				SELECT id, $5, CAST ($6 AS %s) FROM inserted_tx
				UNION ALL
				SELECT id, $7, CAST ($8 AS %s) FROM inserted_tx
				RETURNING *
			)
			SELECT id, occurred_at FROM inserted_tx;
			`,
			database.Transaction,
			database.AffectedWallet,
			database.WalletRole,
			database.WalletRole,
		)

	default:
		query = fmt.Sprintf(
			`
			WITH inserted_tx AS (
				INSERT INTO %s
				(amount, type, category, description)
				VALUES
				($1, $2, $3, $4)
				RETURNING *
			), tx_wallet AS (
				INSERT INTO %s
				(transaction_id, wallet, role)
				SELECT id, $5, CAST ($6 AS %s) FROM inserted_tx
				RETURNING transaction_id, wallet, role
			)
			SELECT id, occurred_at FROM inserted_tx;
			`,
			database.Transaction,
			database.AffectedWallet,
			database.WalletRole,
		)
	}

	return query
}

func (tx *transaction) toExtTx() *Transaction {
	_tx := &Transaction{
		ID:          tx.ID,
		Amount:      tx.Amount,
		Type:        tx.Type,
		Category:    tx.Category,
		Description: tx.Description,
		OccurredAt:  tx.OccurredAt,
	}

	if role := constant.WalletRole(); tx.WalletRole == role.SrcWallet {
		_tx.SrcWallet = tx.Wallet
	} else {
		_tx.DstWallet = tx.Wallet
	}

	return _tx
}
