package pkg

import (
	"context"
	"fmt"

	"github.com/expenseledger/web-service/service"
	"github.com/gin-gonic/gin"
)

func GetUserToken(c *gin.Context) (string, error) {
	var err error = nil
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		err = fmt.Errorf("Token cannot be empty")
	}

	return token, err
}

func GetUserId(c *gin.Context) (string, error) {
	var err error = nil
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		err = fmt.Errorf("Token cannot be empty")
		return "", err
	}

	auth, err := service.GetFirebaseAuth(context.Background())

	if err != nil {
		return "", err
	}

	t, err := auth.VerifyIDToken(context.Background(), token)

	if err != nil {
		return "", err
	}

	return t.UID, nil
}
