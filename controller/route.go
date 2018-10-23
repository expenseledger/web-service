package controller

import "github.com/gin-gonic/gin"

// Route data structure to hold a path and the corresponding handler
type Route struct {
	Path    string
	Handler func(context *gin.Context)
}

// InitRoutes ...
func InitRoutes() *gin.Engine {
	router := gin.Default()

	walletRoute := router.Group("/wallet")
	walletRoute.POST("/create", walletCreate)
	walletRoute.POST("/get", walletGet)
	walletRoute.POST("/delete", walletDelete)
	walletRoute.POST("/list", walletList)
	walletRoute.POST("/listTypes", walletListTypes)

	categoryRoute := router.Group("/category")
	categoryRoute.POST("/create", categoryCreate)
	categoryRoute.POST("/get", categoryGet)
	categoryRoute.POST("/delete", categoryDelete)
	categoryRoute.POST("/list", categoryList)

	return router
}
