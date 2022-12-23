package controllers

import (
	"discord-bot-dashboard-backend-go/discord"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

func AuthController(auth discord.OAuth2Config, scope string, router *gin.Engine) {

	router.GET("/login", func(c *gin.Context) {
		scheme := "http://"
		if c.Request.TLS != nil {
			scheme = "https://"
		}
		fmt.Println("LOGIN")

		params := url.Values{
			"client_id":     {auth.ClientId},
			"scope":         {scope},
			"response_type": {"code"},
			"redirect_uri":  {scheme + c.Request.Host + "/callback"},
		}

		c.Redirect(http.StatusMovedPermanently, "https://discord.com/oauth2/authorize?"+params.Encode())
		c.Abort()
	})

	router.GET("/callback", func(c *gin.Context) {
		println("Handle callback")

		code := c.Query("code")

		if code == "" {
			c.AbortWithStatus(http.StatusBadRequest)
		}

		data, err := discord.GetToken(auth, code)

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		println("Token", data.AccessToken)

		c.SetCookie("access-token", data.AccessToken, data.ExpiresIn, "/", "", false, true)

		c.Redirect(http.StatusOK, auth.RedirectUrl)
	})

	router.GET("/auth", func(c *gin.Context) {
		token, err := c.Cookie("access-token")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		println("Cookie Token", token)

		if discord.CheckToken(token) {
			c.JSON(http.StatusOK, token)
		} else {
			c.JSON(http.StatusUnauthorized, "Invalid token")
		}
	})
}
