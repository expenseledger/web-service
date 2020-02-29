package controller

import (
	"fmt"
	"net/http"

	"github.com/expenseledger/web-service/config"
	"github.com/expenseledger/web-service/pkg"
	"github.com/expenseledger/web-service/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

var configs = config.GetConfigs()

func validateHeader(c *gin.Context) {
	if configs.Mode == "DEVELOPMENT" {
		c.Next()
		return
	}

	token, err := pkg.GetUserToken(c)

	if err != nil {
		buildAbortContext(c, err, http.StatusBadRequest)
		return
	}

	auth, err := service.GetFirebaseAuth(context.Background())

	if err != nil {
		buildAbortContext(c, fmt.Errorf("Cannot initialize firebase auth, %v", err), http.StatusInternalServerError)
		return
	}

	_, err = auth.VerifyIDToken(context.Background(), token)

	if err != nil {
		buildAbortContext(c, fmt.Errorf("Token is invalid"), http.StatusBadRequest)
		return
	}

	c.Next()
}
