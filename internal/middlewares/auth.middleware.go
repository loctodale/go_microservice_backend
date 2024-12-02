package middlewares

import (
	"github.com/gin-gonic/gin"
	"go_microservice_backend_api/pkg/response"
)

func AuthenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "valid-token" {
			response.ErrorResponse(c, response.ErrInvalidToken)
			c.Abort()
			return
		}
		c.Next()
	}
}
