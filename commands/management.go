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
		err := discord.GuildMemberDelete(event.GuildID, v.ID)

		if err != nil {
			zap.L().Error("failed to kick user",
				zap.String("user", v.Username+v.Discriminator),
				zap.Error(err))
		}
	}
}

func BanUser(discord *discordgo.Session, event *discordgo.MessageCreate, _ *routing.Command) {
	if !validateMemberPermissions(discord, event.Member, discordgo.PermissionAdministrator|discordgo.PermissionBanMembers) {
		zap.L().Info("member attempted kick command with insufficient permissions")
		return
	}

	for _, v := range event.Mentions {
		err := discord.GuildBanCreate(event.GuildID, v.ID, 0)

		if err != nil {
			zap.L().Error("failed to ban user",
				zap.String("user", v.Username+v.Discriminator),
				zap.Error(err))
		}
	}
}
