package controller

import (
	"net/http"

	"github.com/expenseledger/web-service/constant"
	"github.com/expenseledger/web-service/model"
	"github.com/gin-gonic/gin"
)

type transactionIdentifyForm struct {
	ID string `json:"id" binding:"required"`
}

func transactionClear(context *gin.Context) {
	var (
		eis model.ExpenseIncomes
		tfs model.Transfers
	)
	eisLength, err := eis.Clear()
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	tfsLength, err := tfs.Clear()
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	var (
		expenses model.ExpenseIncomes
		incomes  model.ExpenseIncomes
	)
	for _, ei := range eis {
		switch ei.Type {
		case constant.TransactionType.Expense:
			expenses = append(expenses, ei)
		case constant.TransactionType.Income:
			incomes = append(incomes, ei)
		}
	}

	data := map[string]interface{}{
		"transfers": tfs,
		"expenses":  expenses,
		"incomes":   incomes,
	}
	items := itemList{
		Length: eisLength + tfsLength,
		Items:  data,
	}

	context.JSON(
		http.StatusOK,
		buildSuccessResponse(items),
	)
	return
}
