package controller

import (
	"fmt"
	"net/http"

	"github.com/expenseledger/web-service/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

func validateHeader(c *gin.Context) {
	firebase, err := service.GetFirebaseInstance()

	if err != nil {
		buildAbortContext(c, fmt.Errorf("Cannot initialize firebase, %v", err), http.StatusInternalServerError)
		return
	}

	token := c.Request.Header.Get("Authorization")

	if token == "" {
		buildAbortContext(c, fmt.Errorf("Token cannot be empty"), http.StatusBadRequest)
		return
	}

	auth, err := firebase.Auth(context.Background())

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
