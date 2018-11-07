package controller

type transactionIdentifyForm struct {
	ID string `json:"id" binding:"required"`
}
