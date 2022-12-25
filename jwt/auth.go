package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthPrincipal struct {
	AccessToken string
	UserID      string
	ExpireAt    int64
}

var contextKey = "jwt-auth-token"
var PrincipalCookie = "session"

func Principal(c *gin.Context) AuthPrincipal {
	principal, exist := c.Get(contextKey)

	if exist {
		return principal.(AuthPrincipal)
	} else {
		panic("using jwt principal outside of auth middleware")
	}
}

func AuthMiddleware(config Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawToken := ExtractToken(c)

		if rawToken == "" {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		token, err := jwt.Parse(ExtractToken(c), func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.Secret), nil
		})

		if err != nil {
			c.String(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if ok && token.Valid {
			c.Set(contextKey, AuthPrincipal{
				ExpireAt:    int64(int(claims["exp"].(float64))),
				AccessToken: claims["access_token"].(string),
				UserID:      claims["user_id"].(string),
			})

			c.Next()
		} else {
			InvalidateSession(c)
			c.String(http.StatusUnauthorized, "Invalid token")
			c.Abort()
		}
	}
}

func InvalidateSession(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     PrincipalCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   0,
		Domain:   "demo-discord-dashboard",
		SameSite: http.SameSiteNoneMode,
		HttpOnly: false,
		Secure:   true,
	})
}
