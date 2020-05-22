package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/getynge/chatterbot/routing"
	"go.uber.org/zap"
)

func KickUser(discord *discordgo.Session, event *discordgo.MessageCreate, _ *routing.Command) {
	if !validateMemberPermissions(discord, event.Member, discordgo.PermissionAdministrator|discordgo.PermissionKickMembers) {
		zap.L().Info("member attempted kick command with insufficient permissions")
		return
	}

	for _, v := range event.Mentions {
		discord.GuildMemberDelete(event.GuildID, v.ID)
	}
}
