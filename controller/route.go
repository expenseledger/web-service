package controller

import (
	"github.com/expenseledger/web-service/config"
	"github.com/gin-gonic/gin"
)

// Route data structure to hold a path and the corresponding handler
type Route struct {
	Path    string
	Handler func(context *gin.Context)
}

// InitRoutes ...
func InitRoutes() *gin.Engine {
	configs := config.GetConfigs()
	router := gin.Default()

	walletRoute := router.Group("/wallet")
	walletRoute.POST("/create", createWallet)
	walletRoute.POST("/get", getWallet)
	walletRoute.POST("/delete", deleteWallet)
	walletRoute.POST("/list", listWallets)
	walletRoute.POST("/listTypes", listWalletTypes)
	walletRoute.POST("/init", initWallets)
	// walletRoute.POST("/receive", walletReceive)
	// walletRoute.POST("/transfer", walletTransfer)

	categoryRoute := router.Group("/category")
	categoryRoute.POST("/create", createCategory)
	categoryRoute.POST("/get", getCategory)
	categoryRoute.POST("/delete", deleteCategory)
	categoryRoute.POST("/list", listCategories)
	categoryRoute.POST("/init", initCategories)

	transactionRoute := router.Group("/transaction")
	transactionRoute.POST("/createExpense", transactionCreateExpense)
	// transactionRoute.POST("/get", transactionGet)
	// transactionRoute.POST("/delete", transactionDelete)

	if configs.Mode != "PRODUCTION" {
		walletRoute.POST("/clear", clearWallets)
		categoryRoute.POST("/clear", clearCategories)
		// transactionRoute.POST("/clear", transactionClear)
	}

	return router
}
