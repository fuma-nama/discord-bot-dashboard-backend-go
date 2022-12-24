package main

import (
	"discord-bot-dashboard-backend-go/controllers"
	"discord-bot-dashboard-backend-go/database"
	"discord-bot-dashboard-backend-go/discord"
	"discord-bot-dashboard-backend-go/jwt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	jwtConfig := jwt.Config{
		Secret: os.Getenv("JWT_SECRET"),
	}
	botConfig := discord.BotConfig{
		Token: os.Getenv("BOT_TOKEN"),
	}
	dbConfig := database.Config{
		Host:     os.Getenv("PGHOST"),
		Name:     os.Getenv("PGDATABASE"),
		User:     os.Getenv("PGUSER"),
		Password: os.Getenv("PGPASSWORD"),
		Port:     os.Getenv("PGPORT"),
	}
	authConfig := discord.OAuth2Config{
		ClientId:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectUrl:  os.Getenv("API_ENDPOINT") + "/callback",
		ClientUrl:    os.Getenv("REDIRECT_URL"),
		Scope:        "identify guilds",
	}
	corsConfig := cors.Config{
		AllowOrigins:     []string{os.Getenv("CLIENT_URL")},
		AllowMethods:     []string{"GET", "OPTIONS", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}

	db := database.Start(dbConfig)
	bot := discord.NewBot(botConfig, db)
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(CORS(corsConfig))
	if err := router.SetTrustedProxies(nil); err != nil {
		panic("Failed to setup engine")
	}

	router.GET("/ping", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	controllers.AuthController(jwtConfig, authConfig, &router.RouterGroup)
	auth := router.Group("")
	{
		auth.Use(jwt.AuthMiddleware(jwtConfig))
		controllers.GuildController(auth, bot, db)
	}

	if router.Run(":8080") != nil {
		return
	}
}

func CORS(config cors.Config) gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", strings.Join(config.AllowOrigins, ", "))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(config.AllowCredentials))
		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ", "))
		c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ", "))

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
