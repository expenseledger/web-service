package controller

import (
	"fmt"

	"github.com/expenseledger/web-service/constant"
	"github.com/expenseledger/web-service/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

// ValidateHeader validate header
func ValidateHeader() gin.HandlerFunc {
	firebase, err := service.GetFirebaseInstance()

	return func(c *gin.Context) {
		if err != nil {
			respondWithError(c, constant.HTTPStatusTypes().InternalServerError, fmt.Sprintf("Cannot initialize firebase, %s", err.Error()))
			return
		}

		token := c.Request.Header.Get("X-Token")

		if token == "" {
			respondWithError(c, constant.HTTPStatusTypes().BadRequest, "Token cannot be empty")
			return
		}

		auth, err := firebase.Auth(context.Background())

		if err != nil {
			respondWithError(c, constant.HTTPStatusTypes().InternalServerError, fmt.Sprintf("Cannot initialize firebase auth, %s", err.Error()))
			return
		}

		_, err = auth.VerifyIDToken(context.Background(), token)

		if err != nil {
			respondWithError(c, constant.HTTPStatusTypes().BadRequest, "Token is invalid")
			return
		}

		c.Next()
	}
}

func respondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"error": message})
}
