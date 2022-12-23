package controllers

import (
	"discord-bot-dashboard-backend-go/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
)

type GuildInfo struct {
	EnabledFeatures []string `json:"enabledFeatures"`
	CustomField     string   `json:"customField"`
}

type WelcomeMessageOptions struct {
	Message string `json:"message"`
}

func GuildController(router *gin.Engine, db *gorm.DB) {
	router.GET("/guilds/:guild", func(c *gin.Context) {
		guild := c.Param("guild")
		if guild == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		var info *models.Guild
		if err := db.Model(&models.Guild{}).Find(&info, guild).Error; err != nil {
			info = nil
		}

		features := make([]string, 0)

		if info != nil && info.WelcomeMessage != nil {
			features = append(features, "welcome-message")
		}

		c.JSON(http.StatusOK, &GuildInfo{
			EnabledFeatures: features,
			CustomField:     "Hello World!",
		})
	})

	group := router.Group("/guilds/:guild/features")
	{
		group.GET("/welcome-message", func(c *gin.Context) {
			guild := c.Param("guild")
			if guild == "" {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			var info *models.Guild
			if err := db.Model(&models.Guild{}).Find(&info, guild).Error; err != nil {
				info = nil
			}

			if info == nil || info.WelcomeMessage == nil {
				c.AbortWithStatus(http.StatusNotFound)
			} else {
				c.JSON(http.StatusOK, WelcomeMessageOptions{
					Message: *info.WelcomeMessage,
				})
			}
		})

		group.PATCH("/welcome-message", func(c *gin.Context) {
			guild := c.Param("guild")
			if guild == "" {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			var body WelcomeMessageOptions
			if err := c.BindJSON(&body); err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			var updated models.Guild
			err := db.Model(&updated).
				Clauses(clause.Returning{}).
				Where("id = ?", guild).
				Updates(models.Guild{WelcomeMessage: &body.Message}).
				Error

			if err == nil {
				c.JSON(http.StatusOK, WelcomeMessageOptions{
					Message: *updated.WelcomeMessage,
				})
			} else {
				c.AbortWithStatus(http.StatusNotFound)
			}
		})

		group.POST("/welcome-message", func(c *gin.Context) {
			guild := c.Param("guild")
			if guild == "" {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			empty := ""
			err := db.Clauses(
				clause.OnConflict{
					Columns:   []clause.Column{{Name: "id"}},
					DoUpdates: clause.AssignmentColumns([]string{"welcome_message"}),
				},
			).Create(&models.Guild{
				Id:             guild,
				WelcomeMessage: &empty,
			}).Error

			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
			} else {
				c.AbortWithStatus(http.StatusOK)
			}
		})

		group.DELETE("/welcome-message", func(c *gin.Context) {
			guild := c.Param("guild")
			if guild == "" {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			db.Delete(&models.Guild{
				Id: guild,
			})

			c.AbortWithStatus(http.StatusOK)
		})
	}
}
