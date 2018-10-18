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

	return router
}
