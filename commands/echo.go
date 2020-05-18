package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/getynge/chatterbot/routing"
	"go.uber.org/zap"
)

func Echo(discord *discordgo.Session, event *discordgo.MessageCreate, command *routing.Command) {
	_, err := discord.ChannelMessageSend(event.ChannelID, command.Arguments)

	if err != nil {
		zap.L().Error("could not echo message", zap.Error(err))
	}
}
