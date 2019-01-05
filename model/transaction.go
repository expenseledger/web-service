package model

import (

	// dbmodel "github.com/expenseledger/web-service/db/model"
	"errors"
	"time"

	"github.com/expenseledger/web-service/constant"
	"github.com/expenseledger/web-service/orm"
	"github.com/expenseledger/web-service/pkg/type/date"
	"github.com/shopspring/decimal"
)

// Transaction the structure represents a transaction in application layer
type Transaction struct {
	ID          string                   `json:"id" db:"id"`
	From        string                   `json:"src_wallet" db:"src_wallet"`
	To          string                   `json:"dst_wallet" db:"dst_wallet"`
	Amount      decimal.Decimal          `json:"amount" db:"amount"`
	Type        constant.TransactionType `json:"type" db:"type"`
	Category    string                   `json:"category" db:"category"`
	Description string                   `json:"description" db:"description"`
	Date        date.Date                `json:"date"`
	OccurredAt  time.Time                `json:"-" db:"occurred_at"`
}

// pass this to ORM
type _Transaction struct {
	ID          string                   `db:"id"`
	Wallet      string                   `db:"wallet"`
	Role        constant.WalletRole      `db:"role"`
	Type        constant.TransactionType `db:"type"`
	Amount      decimal.Decimal          `db:"amount"`
	Category    string                   `db:"category"`
	Description string                   `db:"description"`
	OccurredAt  time.Time                `db:"occurred_at"`
	CreatedAt   time.Time                `db:"created_at"`
}

func CreateTransction(
	amount decimal.Decimal,
	t constant.TransactionType,
	from string,
	to string,
	category string,
	description string,
	d date.Date,
) (*Transaction, error) {
	tim := time.Time(d)
	if tim.IsZero() {
		tim = time.Now()
	}

	if txTypes := constant.TransactionTypes(); t != txTypes.Transfer {
		var wallet string
		switch t {
		case txTypes.Expense:
			wallet = from
		case txTypes.Income:
			wallet = to
		}

		tx, err := createNonTransferTx(
			amount,
			t,
			wallet,
			category,
			description,
			tim,
		)
		if err != nil {
			return nil, err
		}
		return tx, nil
	}

	tx := Transaction{
		From:        from,
		To:          to,
		Amount:      amount,
		Type:        t,
		Category:    category,
		Description: description,
		OccurredAt:  tim,
	}

	mapper := orm.NewTxMapper(tx, t)

	tmp, err := mapper.Insert(&tx)
	if err != nil {
		return nil, err
	}

	tmpTx := tmp.(*Transaction)
	tmpTx.Date = date.Date(tmpTx.OccurredAt)
	return tmpTx, nil
}

func GetTransaction(id string) (*Transaction, error) {
	txTypes := constant.TransactionTypes()
	_tx := _Transaction{ID: id}
	mapper := orm.NewTxMapper(_tx, txTypes.Expense)

	tmp, err := mapper.One(&_tx)
	if err != nil {
		return nil, err
	}

	_txs := *(tmp.(*[]_Transaction))
	length := len(_txs)
	if length <= 0 {
		return nil, errors.New("transaction not found")
	}

	tx := _txs[0].toTransaction()
	for i := 1; i < length; i++ {
		tx.To = _txs[i].Wallet
	}

	return tx, nil
}

func ClearTransactions() ([]Transaction, error) {
	txTypes := constant.TransactionTypes()
	mapper := orm.NewTxMapper(_Transaction{}, txTypes.Expense)

	tmp, err := mapper.Clear()
	if err != nil {
		return nil, err
	}

	_txs := *(tmp.(*[]_Transaction))
	txs := make([]Transaction, 0, len(_txs))
	length := len(_txs)

	for i := 0; i < length; i++ {
		_tx := _txs[i]
		tx := _tx.toTransaction()

		if tx.Type == txTypes.Transfer {
			i++
			_tx = _txs[i]
			tx.To = _tx.Wallet
		}

		txs = append(txs, *tx)
	}

	return txs, nil
}

func createNonTransferTx(
	amount decimal.Decimal,
	t constant.TransactionType,
	wallet string,
	category string,
	description string,
	tim time.Time,
) (*Transaction, error) {
	tx := _Transaction{
		Wallet:      wallet,
		Type:        t,
		Amount:      amount,
		Category:    category,
		Description: description,
		OccurredAt:  tim,
	}

	if t == constant.TransactionTypes().Expense {
		tx.Role = constant.WalletRoles().SrcWallet
	} else {
		tx.Role = constant.WalletRoles().DstWallet
	}

	mapper := orm.NewTxMapper(tx, t)

	tmp, err := mapper.Insert(&tx)
	if err != nil {
		return nil, err
	}

	tmpTx := tmp.(*_Transaction)
	return tmpTx.toTransaction(), nil
}

func (tx *_Transaction) toTransaction() *Transaction {
	tmpTx := Transaction{
		ID:          tx.ID,
		Amount:      tx.Amount,
		Type:        tx.Type,
		Category:    tx.Category,
		Description: tx.Description,
		OccurredAt:  tx.OccurredAt,
	}

	if tx.Role == constant.WalletRoles().SrcWallet {
		tmpTx.From = tx.Wallet
	} else {
		tmpTx.To = tx.Wallet
	}

	tmpTx.Date = date.Date(tmpTx.OccurredAt)

	return &tmpTx
}

// Transactions is defined just to be used as a receiver
// type Transactions []Transaction

// Create ...
// func (tx *Transaction) Create() error {
// 	dbTx := tx.toDBCounterpart()
// 	if err := dbTx.Insert(); err != nil {
// 		return err
// 	}

// 	tx.fromDBCounterpart(dbTx)
// 	return nil
// }

// // Clear ...
// func (txs *Transactions) Clear() (int, error) {
// 	var dbTxs dbmodel.Transactions
// 	if err := dbTxs.DeleteAll(); err != nil {
// 		return 0, err
// 	}

// 	length := len(dbTxs)
// 	transactions := make(Transactions, length)
// 	for i, dbTx := range dbTxs {
// 		var tx Transaction
// 		tx.fromDBCounterpart(&dbTx)
// 		transactions[i] = tx
// 	}

// 	*txs = transactions
// 	return length, nil
// }

// // Get ...
// func (tx *Transaction) Get() error {
// 	dbTx := tx.toDBCounterpart()
// 	if err := dbTx.One(); err != nil {
// 		return err
// 	}

// 	tx.fromDBCounterpart(dbTx)
// 	return nil
// }

// // Delete ...
// func (tx *Transaction) Delete() error {
// 	dbTx := tx.toDBCounterpart()
// 	if err := dbTx.Delete(); err != nil {
// 		return err
// 	}

// 	tx.fromDBCounterpart(dbTx)
// 	return nil
// }
