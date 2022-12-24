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

	discord.Identify.Intents = discordgo.IntentsGuildMembers

	if err := discord.Open(); err != nil {
		panic(err.Error())
	}

	return discord
}

func GuildMemberAdd(db *gorm.DB) func(s *discordgo.Session, event *discordgo.GuildMemberAdd) {
	return func(s *discordgo.Session, event *discordgo.GuildMemberAdd) {
		var result models.Guild
		db.Find(&result, event.GuildID)

		if result.WelcomeMessage == nil || result.WelcomeChannel == nil {
			return
		}

		_, _ = s.ChannelMessageSend(*result.WelcomeChannel, *result.WelcomeMessage)
	}
}
