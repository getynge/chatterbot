package middleware

import (
	"github.com/bwmarrin/discordgo"
	"github.com/getynge/chatterbot/routing"
	"go.uber.org/zap"
)

// PanicHandling logs panics without crashing the program when they occur within the router
func PanicHandling(e routing.EventHandler) routing.EventHandler {
	return routing.EventHandlerFunc(func(session *discordgo.Session, create *discordgo.MessageCreate, command *routing.Command) {
		defer func() {
			if r := recover(); r != nil {
				zap.L().Error("Recovered from panic",
					zap.String("command", command.Command),
					zap.String("arguments", command.Arguments),
					zap.Any("r", r))
			}
		}()

		e.Handler(session, create, command)
	})
}
