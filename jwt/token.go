package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Secret string
}

func GenerateToken(config Config, accessToken string, userId string, expireIn int) (string, error) {
	claims := jwt.MapClaims{}
	claims["access_token"] = accessToken
	claims["user_id"] = userId
	claims["exp"] = time.Now().Add(time.Second * time.Duration(expireIn)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.Secret))
}

func ExtractToken(c *gin.Context) string {
	var bearerToken = ""

	if cookie, err := c.Cookie(PrincipalCookie); err == nil {
		bearerToken = cookie
	}

	if header := strings.Split(c.Request.Header.Get("Authorization"), " "); len(header) == 2 {
		bearerToken = header[1]
	}

	return bearerToken
}
