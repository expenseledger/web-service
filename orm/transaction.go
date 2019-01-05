package orm

import (
	"errors"

	"github.com/expenseledger/web-service/constant"
)

type TxMapper struct {
	BaseMapper
	transferStmt string
	txType       constant.TransactionType
}

func (mapper *TxMapper) Insert(obj interface{}) (interface{}, error) {
	txType := constant.TransactionTypes()
	switch mapper.txType {
	case txType.Transfer:
		return worker(
			obj,
			mapper.modelType,
			mapper.transferStmt,
			"Error inserting",
		)
	case txType.Expense:
		fallthrough
	case txType.Income:
		return worker(
			obj,
			mapper.modelType,
			mapper.insertStmt,
			"Error inserting",
		)
	}
	return nil, errors.New("unknown transaction type")
}

func (mapper *TxMapper) One(obj interface{}) (interface{}, error) {
	return sliceWorker(
		obj,
		mapper.modelType,
		mapper.oneStmt,
		"Error getting",
	)
}
