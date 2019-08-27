package controller

import (
	"github.com/expenseledger/web-service/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

// ValidateHeader validate header
func ValidateHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		firebase, err := service.GetFirebaseInstance()

		if err != nil {
			return
		}

		token := c.Request.Header.Get("X-Token")

		if token == "" {
			return
		}

		auth, err := firebase.Auth(context.Background())

		if err != nil {
			return
		}

		_, err = auth.VerifyIDToken(context.Background(), token)

		if err != nil {
			return
		}

		c.Next()
	}
}
