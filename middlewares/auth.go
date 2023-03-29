package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

var contextKey = "jwt-auth-token"

func GetToken(c *gin.Context) string {
	principal, exist := c.Get(contextKey)

	if exist {
		return principal.(string)
	} else {
		panic("using jwt principal outside of auth middleware")
	}
}

func extractToken(c *gin.Context) string {
	var bearerToken = ""

	if header := strings.Split(c.Request.Header.Get("Authorization"), " "); len(header) == 2 {
		bearerToken = header[1]
	}

	return bearerToken
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)

		if token == "" {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		c.Set(contextKey, token)
		c.Next()
	}
}
