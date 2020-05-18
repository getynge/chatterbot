package routing

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"time"
)

type Command struct {
	Command   string
	Arguments string
	Ctx       context.Context
	Cancel    context.CancelFunc
}

func NewCommand(timeout time.Duration) *Command {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return &Command{
		Ctx:    ctx,
		Cancel: cancel,
	}
}

// An event handler is any type that can receive discordgo events
// The event handler is typically going to be either a router or a wrapper created by EventHandlerFunc
type EventHandler interface {
	// An event handler, see the discordgo documentation for details on how event handling works
	// This entry being the most relevant: https://godoc.org/github.com/bwmarrin/discordgo#Session.AddHandler
	Handler(*discordgo.Session, *discordgo.MessageCreate, *Command)
}

type eventHandlerWrapper struct {
	f func(*discordgo.Session, *discordgo.MessageCreate, *Command)
}

func (e eventHandlerWrapper) Handler(discord *discordgo.Session, event *discordgo.MessageCreate, command *Command) {
	e.f(discord, event, command)
}

// creates a wrapper around
func EventHandlerFunc(f func(*discordgo.Session, *discordgo.MessageCreate, *Command)) EventHandler {
	return eventHandlerWrapper{f}
}
