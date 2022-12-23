package main

import (
	"discord-bot-dashboard-backend-go/controllers"
	"discord-bot-dashboard-backend-go/db"
	"discord-bot-dashboard-backend-go/discord"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
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
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "OPTIONS", "POST", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/ping", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	controllers.AuthController(auth, scope, router)

	if router.Run(":8080") != nil {
		return
	}
}
