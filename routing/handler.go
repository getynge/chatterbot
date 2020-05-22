package routing

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"time"
)

type Command struct {
	// The trimmed command issued by the user (e.g. "grant")
	Command string
	// The text of the message with the command removed
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

// EventHandler is any type that can receive discordgo events
// The event handler is typically going to be either a router or a wrapper created by EventHandlerFunc
type EventHandler interface {
	// Handler is an event handler, see the discordgo documentation for details on how event handling works
	// This entry being the most relevant: https://godoc.org/github.com/bwmarrin/discordgo#Session.AddHandler
	Handler(*discordgo.Session, *discordgo.MessageCreate, *Command)
}

type eventHandlerWrapper struct {
	f func(*discordgo.Session, *discordgo.MessageCreate, *Command)
}

func (e eventHandlerWrapper) Handler(discord *discordgo.Session, event *discordgo.MessageCreate, command *Command) {
	e.f(discord, event, command)
}

// EventHandlerFunc creates a wrapper around a function capable of handling events.
// Use this function whenever you need an event handler that is not stateful
func EventHandlerFunc(f func(*discordgo.Session, *discordgo.MessageCreate, *Command)) EventHandler {
	return eventHandlerWrapper{f}
}
