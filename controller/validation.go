package controller

import (
	"fmt"
	"net/http"

	"github.com/expenseledger/web-service/pkg"
	"github.com/expenseledger/web-service/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

func validateHeader(c *gin.Context) {
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
