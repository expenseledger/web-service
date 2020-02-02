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
	UserId      string                   `json:"userId" db:"user_id"`
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
	UserId      string                   `db:"user_id"`
}

type transactions []_Transaction

func CreateTransction(
	amount decimal.Decimal,
	t constant.TransactionType,
	from string,
	to string,
	category string,
	description string,
	d date.Date,
	userId string,
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
			userId,
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
		UserId:      userId,
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

func GetTransaction(id string, userId string) (*Transaction, error) {
	return applyToTx(id, one, userId)
}

func DeleteTransaction(id string, userId string) (*Transaction, error) {
	return applyToTx(id, delete, userId)
}

func ListTransactions(walletName string, userId string) ([]Transaction, error) {
	txTypes := constant.TransactionTypes()
	mapper := orm.NewTxMapper(_Transaction{UserId: userId}, txTypes.Expense)

	_tx := _Transaction{Wallet: walletName, UserId: userId}

	tmp, err := mapper.Many(&_tx)
	if err != nil {
		return nil, err
	}

	_txs := (*transactions)(tmp.(*[]_Transaction))
	return _txs.toTransactions(), nil
}

func ClearTransactions(userId string) ([]Transaction, error) {
	txTypes := constant.TransactionTypes()
	mapper := orm.NewTxMapper(_Transaction{UserId: userId}, txTypes.Expense)

	tmp, err := mapper.Clear()
	if err != nil {
		return nil, err
	}

	_txs := (*transactions)(tmp.(*[]_Transaction))
	return _txs.toTransactions(), nil
}

func createNonTransferTx(
	amount decimal.Decimal,
	t constant.TransactionType,
	wallet string,
	category string,
	description string,
	tim time.Time,
	userId string,
) (*Transaction, error) {
	tx := _Transaction{
		Wallet:      wallet,
		Type:        t,
		Amount:      amount,
		Category:    category,
		Description: description,
		OccurredAt:  tim,
		UserId:      userId,
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
		UserId:      tx.UserId,
	}

	if tx.Role == constant.WalletRoles().SrcWallet {
		tmpTx.From = tx.Wallet
	} else {
		tmpTx.To = tx.Wallet
	}

	tmpTx.Date = date.Date(tmpTx.OccurredAt)

	return &tmpTx
}

func (tmpTxs *transactions) toTransactions() []Transaction {
	_txs := *tmpTxs
	length := len(_txs)
	txs := make([]Transaction, 0, length)
	txTypes := constant.TransactionTypes()

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

	return txs
}

func applyToTx(id string, op operation, userId string) (*Transaction, error) {
	_tx := _Transaction{ID: id, UserId: userId}
	mapper := orm.NewTxMapper(_tx, constant.TransactionTypes().Expense)

	var tmp interface{}
	var err error
	switch op {
	case delete:
		tmp, err = mapper.Delete(&_tx)
	case one:
		tmp, err = mapper.One(&_tx)
	}

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
