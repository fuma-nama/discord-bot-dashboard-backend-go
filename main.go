package main

import (
	"discord-bot-dashboard-backend-go/controllers"
	"discord-bot-dashboard-backend-go/database"
	"discord-bot-dashboard-backend-go/discord"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

func main() {
	dbConfig := database.DataBaseConfig{
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

	db := database.Start(dbConfig)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(CORS())

	router.GET("/ping", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	controllers.AuthController(auth, scope, router)
	controllers.GuildController(router, db)

	if router.Run(":8080") != nil {
		return
	}
}

func CORS() gin.HandlerFunc {
	config := cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "OPTIONS", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type"},
	}

	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", strings.Join(config.AllowOrigins, ", "))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ", "))
		c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ", "))

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
