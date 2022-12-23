package controllers

import (
	"discord-bot-dashboard-backend-go/discord"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

var tokenCookie = "access-token"

func AuthController(auth discord.OAuth2Config, scope string, router *gin.Engine) {

	router.GET("/login", func(c *gin.Context) {
		params := url.Values{
			"client_id":     {auth.ClientId},
			"scope":         {scope},
			"response_type": {"code"},
			"redirect_uri":  {getCallbackUrl(c)},
		}

		c.Redirect(http.StatusMovedPermanently, "https://discord.com/oauth2/authorize?"+params.Encode())
		c.Abort()
	})

	router.GET("/callback", func(c *gin.Context) {
		code := c.Query("code")

		if code == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		data, err := discord.GetToken(auth, getCallbackUrl(c), code)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.SetCookie(tokenCookie, data.AccessToken, data.ExpiresIn, "/", "", false, true)
		c.Redirect(http.StatusMovedPermanently, auth.RedirectUrl)
		c.Abort()
	})

	router.GET("/auth", func(c *gin.Context) {
		token, err := c.Cookie(tokenCookie)
		if err != nil || token == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if discord.CheckToken(token) {
			c.JSON(http.StatusOK, token)
		} else {
			c.JSON(http.StatusUnauthorized, "Invalid token")
		}
	})

	router.POST("/auth/signout", func(c *gin.Context) {
		token, err := c.Cookie(tokenCookie)
		if err != nil || token == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if err := discord.RevokeToken(auth, token); err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "access-token",
			Value:    "",
			Path:     "/",
			MaxAge:   0,
			HttpOnly: true,
		})
	})
}

func getCallbackUrl(c *gin.Context) string {
	scheme := "http://"
	if c.Request.TLS != nil {
		scheme = "https://"
	}

	return scheme + c.Request.Host + "/callback"
}
