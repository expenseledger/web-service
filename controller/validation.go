package controller

import (
	"github.com/gin-gonic/gin"
)

// ValidateHeader validate header
func ValidateHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("X-Token")

		if token == "" {
			return
		}

		c.Next()
	}
}
