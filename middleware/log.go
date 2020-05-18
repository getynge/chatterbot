package middleware

import (
	"github.com/bwmarrin/discordgo"
	"github.com/getynge/chatterbot/routing"
	"go.uber.org/zap"
	"time"
)

func Logging(e routing.EventHandler) routing.EventHandler {
	return routing.EventHandlerFunc(func(discord *discordgo.Session, event *discordgo.MessageCreate, command *routing.Command) {
		t1 := time.Now()
		defer func() {
			zap.L().Info("handled command",
				zap.String("command", command.Command),
				zap.Duration("time elapsed", time.Since(t1)))
		}()

		e.Handler(discord, event, command)
	})
}
