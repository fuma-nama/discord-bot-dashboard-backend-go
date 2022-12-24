package controllers

import (
	"discord-bot-dashboard-backend-go/jwt"
	"discord-bot-dashboard-backend-go/models"
	"errors"
	"github.com/bwmarrin/discordgo"
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
	Channel *string `json:"channel"`
	Message *string `json:"message"`
}

func GuildController(router *gin.RouterGroup, bot *discordgo.Session, db *gorm.DB) {
	router.GET("/guilds/:guild", func(c *gin.Context) {
		guild, err := guildAndCheck(bot, c)
		if err != nil {
			return
		}

		var info models.Guild
		db.Find(&info, guild.ID)

		features := make([]string, 0)

		if info.Id != "" && info.WelcomeMessage != nil {
			features = append(features, "welcome-message")
		}

		c.JSON(http.StatusOK, &GuildInfo{
			EnabledFeatures: features,
			CustomField:     "Hello World!",
		})
	})

	router.GET("/guilds/:guild/roles", func(c *gin.Context) {
		guild, err := guild(c)
		if err != nil {
			return
		}

		roles, err := bot.GuildRoles(*guild)
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, roles)
	})

	router.GET("/guilds/:guild/channels", func(c *gin.Context) {
		guild, err := guild(c)
		if err != nil {
			return
		}

		channels, err := bot.GuildChannels(*guild)
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, channels)
	})

	group := router.Group("/guilds/:guild/features")
	{
		group.GET("/welcome-message", func(c *gin.Context) {
			guild, err := guild(c)
			if err != nil {
				return
			}

			var info models.Guild
			db.Find(&info, guild)

			if info.Id == "" || info.WelcomeMessage == nil {
				c.AbortWithStatus(http.StatusNotFound)
			} else {
				c.JSON(http.StatusOK, WelcomeMessageOptions{
					Channel: info.WelcomeChannel,
					Message: info.WelcomeMessage,
				})
			}
		})

		group.PATCH("/welcome-message", func(c *gin.Context) {
			guild, err := guildAndCheck(bot, c)
			if err != nil {
				return
			}

			var body WelcomeMessageOptions
			if err := c.BindJSON(&body); err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			if body.Channel != nil {
				channels, err := bot.GuildChannels(guild.ID)

				if err != nil || !containsChannel(channels, *body.Channel) {
					c.AbortWithStatus(http.StatusBadRequest)
					return
				}
			}

			var updated models.Guild
			err = db.Model(&updated).
				Clauses(clause.Returning{}).
				Where("id = ?", guild.ID).
				Updates(models.Guild{WelcomeMessage: body.Message, WelcomeChannel: body.Channel}).
				Error

			if err == nil {
				c.JSON(http.StatusOK, WelcomeMessageOptions{
					Channel: updated.WelcomeChannel,
					Message: updated.WelcomeMessage,
				})
			} else {
				c.AbortWithStatus(http.StatusNotFound)
			}
		})

		group.POST("/welcome-message", func(c *gin.Context) {
			guild, err := guild(c)
			if err != nil {
				return
			}

			empty := ""
			err = db.Clauses(
				clause.OnConflict{
					Columns:   []clause.Column{{Name: "id"}},
					DoUpdates: clause.AssignmentColumns([]string{"welcome_message"}),
				},
			).Create(&models.Guild{
				Id:             *guild,
				WelcomeMessage: &empty,
			}).Error

			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
			} else {
				c.AbortWithStatus(http.StatusOK)
			}
		})

		group.DELETE("/welcome-message", func(c *gin.Context) {
			guild, err := guild(c)
			if err != nil {
				return
			}

			db.Delete(&models.Guild{
				Id: *guild,
			})

			c.AbortWithStatus(http.StatusOK)
		})
	}
}

func guild(c *gin.Context) (*string, error) {
	guild := c.Param("guild")
	if guild == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return nil, errors.New("invalid request")
	}

	return &guild, nil
}

func guildAndCheck(bot *discordgo.Session, c *gin.Context) (id *discordgo.Guild, err error) {
	guild, err := guild(c)
	if err != nil {
		return
	}

	guildData, err := bot.Guild(*guild)
	if err != nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	principal := jwt.Principal(c)
	member, err := bot.GuildMember(*guild, principal.UserID)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if member.Permissions&8 == 0 && member.User.ID != guildData.OwnerID {
		c.AbortWithStatus(http.StatusBadRequest)
		return nil, errors.New("member missing permissions")
	}

	return guildData, nil
}

func containsChannel(s []*discordgo.Channel, channel string) bool {
	for _, v := range s {
		if v.ID == channel {
			return true
		}
	}

	return false
}
