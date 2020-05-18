package routing

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	"strings"
	"time"
)

// A simple router for commands, typical of what you would find in an HTTP router.
type Router struct {
	prefixes   []string
	middleware []Middleware
	routes     map[string]EventHandler
	// The timeout for commands handled by this router, defaults to time.Second
	Timeout time.Duration
	// Specify whether or not to ignore commands issued by bots.
	// Regardless of the setting of this value, the bot will always ignore it's own messages.
	// This defaults to true when the router is created via NewRouter
	IgnoreBots bool
}

func NewRouter(prefixes ...string) *Router {
	return &Router{
		prefixes:   prefixes,
		routes:     make(map[string]EventHandler),
		Timeout:    time.Second,
		IgnoreBots: true,
	}
}

// Never call this function directly
// Pass this function as an argument to discordgo.AddHandler in order to bootstrap the router
func (r *Router) HandlerBootstrap(discord *discordgo.Session, event *discordgo.MessageCreate) {
	command := NewCommand(r.Timeout)
	r.Handler(discord, event, command)
}

// Never call this function directly.
// All routes have leading and trailing whitespace removed, along with the prefix.
// Assuming the prefix is "$", then all of the following are the same:
//  "$ Hi"
//  "$Hi"
//  "$ Hi "
// Note that if you have multiple prefixes, only one of them will be recognized at a time.
func (r *Router) Handler(discord *discordgo.Session, event *discordgo.MessageCreate, command *Command) {
	prefix := ""
	hasPrefix := false
	handler := EventHandlerFunc(notFoundHandler)

	if (event.Author.Bot && r.IgnoreBots) || event.Author.ID == discord.State.User.ID {
		return
	}

	for _, v := range r.prefixes {
		if strings.HasPrefix(event.Content, v) || v == "" {
			prefix = v
			hasPrefix = true
			break
		}
	}

	if !hasPrefix {
		return
	}

	var deprefixed string
	if command.Command != "" {
		deprefixed = command.Arguments
	} else {
		deprefixed = strings.TrimSpace(strings.TrimPrefix(event.Content, prefix))
	}
	command.Command = strings.Split(deprefixed, " ")[0]

	if route, ok := r.routes[command.Command]; ok {
		handler = route
	}

	command.Arguments = strings.TrimSpace(strings.TrimPrefix(deprefixed, command.Command))

	for _, v := range r.middleware {
		handler = v(handler)
	}

	handler.Handler(discord, event, command)
}

func (r *Router) AddCommand(command string, handler EventHandler) {
	r.routes[command] = handler
}

func (r *Router) AddCommandFunc(command string, handler func(*discordgo.Session, *discordgo.MessageCreate, *Command)) {
	r.routes[command] = EventHandlerFunc(handler)
}

// Adds a subcommand to the routes, where in "$x y z" y is a subcommand.
func (r *Router) AddSubcommand(command string, f func(r *Router)) {
	s := NewRouter("")
	r.routes[command] = s
	f(s)
}

// Adds a middleware function to the chain.
// Middleware is called in the order that it is defined.
func (r *Router) Use(m Middleware) {
	r.middleware = append(r.middleware, m)
}

func notFoundHandler(discord *discordgo.Session, event *discordgo.MessageCreate, command *Command) {
	_, err := discord.ChannelMessageSend(event.ChannelID, fmt.Sprintf("%s command not found", command.Command))

	if err != nil {
		zap.L().Error("could not send message", zap.Error(err))
	}
}
