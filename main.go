package main

import (
	"discord-bot-dashboard-backend-go/db"
	"discord-bot-dashboard-backend-go/discord"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"os"
)

func main() {
	dbConfig := db.DataBaseConfig{
		Host:     os.Getenv("db_host"),
		Name:     os.Getenv("db_name"),
		User:     os.Getenv("db_username"),
		Password: os.Getenv("db_password"),
	}

	auth := discord.OAuth2Config{
		ClientId:     os.Getenv("client_id"),
		ClientSecret: os.Getenv("client_secret"),
		RedirectUrl:  os.Getenv("redirect_url"),
	}

	scope := "identify guilds"

	db.Start(dbConfig)
	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/ping", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/login", func(c *gin.Context) {
		c.Redirect(http.StatusOK, "https://discord.com/oauth2/authorize?response_type=code&client_id="+url.QueryEscape(auth.ClientId)+"&scope="+url.QueryEscape(scope))
	})

	router.GET("/callback", func(c *gin.Context) {
		code := c.Query("code")

		if code == "" {
			c.AbortWithStatus(http.StatusBadRequest)
		}

		data, err := discord.GetToken(auth, code)

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		c.SetCookie("access-token", data.AccessToken, data.ExpiresIn, "/", "", false, true)

		c.Redirect(http.StatusOK, auth.RedirectUrl)
	})

	router.GET("/auth", func(c *gin.Context) {
		token, err := c.Cookie("access-token")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		if discord.CheckToken(token) {
			c.JSON(http.StatusOK, token)
		} else {
			c.JSON(http.StatusUnauthorized, "Invalid token")
		}
	})

	if router.Run() != nil {
		return
	}
}
