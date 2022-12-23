package discord

import (
	"discord-bot-dashboard-backend-go/models"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

type BotConfig struct {
	Token string
}

func NewBot(config BotConfig, db *gorm.DB) *discordgo.Session {
	discord, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		panic(err.Error())
	}

	discord.AddHandler(GuildMemberAdd(db))

	return discord
}

func GuildMemberAdd(db *gorm.DB) func(s *discordgo.Session, event *discordgo.GuildMemberAdd) {
	return func(s *discordgo.Session, event *discordgo.GuildMemberAdd) {

		var result *models.Guild
		db.Model(models.Guild{Id: event.GuildID}).Find(result)

		if result == nil || result.WelcomeMessage == nil {
			return
		}

		//s.ChannelMessageSend()
	}
}
