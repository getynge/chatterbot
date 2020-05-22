package commands

import (
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func validateMemberPermissions(discord *discordgo.Session, member *discordgo.Member, permission int) bool {
	for _, v := range member.Roles {
		role, err := discord.State.Role(member.GuildID, v)
		if err != nil {
			zap.L().Error("could not get user roles", zap.Error(err))
			return false
		}
		if role.Permissions&permission != 0 {
			return true
		}
	}
	return false
}
