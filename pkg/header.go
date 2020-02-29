package pkg

import (
	"context"
	"fmt"

	"github.com/expenseledger/web-service/config"
	"github.com/expenseledger/web-service/service"
	"github.com/gin-gonic/gin"
)

var configs = config.GetConfigs()

const TestUserId = "t4ND1OuMCRVvM6sMuNS0c48Ad0o2"

func GetUserToken(c *gin.Context) (string, error) {
	var err error = nil
	authorizedHeader := c.Request.Header.Get("Authorization")

	if authorizedHeader == "" {
		err = fmt.Errorf("authorizedHeader cannot be empty")
	}

	token := authorizedHeader[7:]

	if token == "" {
		err = fmt.Errorf("token cannot be empty")
	}

	return token, err
}

func GetUserId(c *gin.Context) (string, error) {
	var err error = nil

	if configs.Mode == "DEVELOPMENT" {
		return TestUserId, nil
	}

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
