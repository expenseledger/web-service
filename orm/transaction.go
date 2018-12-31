package orm

import "errors"

type TxMapper struct {
	BaseMapper
	transferStmt string
	txType       string
}

// Insert ...
func (mapper *TxMapper) Insert(obj interface{}) (interface{}, error) {
	switch mapper.txType {
	case "TRANSFER":
		return worker(
			obj,
			mapper.modelType,
			mapper.transferStmt,
			"Error inserting",
		)
	case "EXPENSE":
		fallthrough
	case "INCOME":
		return worker(
			obj,
			mapper.modelType,
			mapper.insertStmt,
			"Error inserting",
		)
	}
	return nil, errors.New("unknown transaction type")
}
