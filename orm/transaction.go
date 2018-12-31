package orm

type TransactionMapper struct {
	BaseMapper
	transferStmt string
	txType       string
}
