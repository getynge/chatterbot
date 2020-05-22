/*
Package routing implements a simple router for discord commands

This package provides two interfaces relevant to the user: EventHandler and Middleware.

Middleware is used to apply logic that is common to all routes within a Router.
A prime example of useful Middleware is the Logging middleware found in the middleware package, with the following source:
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
Middleware is not applied to every Router in a tree of Routers, only to the router on which Router.Use is called.

All Routers are EventHandlers as well as all routes within those Routers.
Routers can be nested within one another using the AddSubcommand function, which is useful for creating commands
such as "permission grant <user> <permission>."
It is generally good practice to break commands up into categories to make them more ergonomic for both users and
developers.
Prefer
 r := routing.NewRouter("$")
 r.AddSubcommand("token", func(r *routing.Router)){
	 r.AddCommandFunc("give", giveToken)
 }
over
 r := routing.NewRouter("$")
 r.AddCommandFunc("give-token", giveToken)
*/
package routing

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	"strings"
	"time"
)

// Router is A simple router for commands, typical of what you would find in an HTTP router.
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

// NewRouter creates a new router with default values and the specified prefixes.
// A new router created this way will have a timeout of one second, and will ignore all bots by default.
func NewRouter(prefixes ...string) *Router {
	return &Router{
		prefixes:   prefixes,
		routes:     make(map[string]EventHandler),
		Timeout:    time.Second,
		IgnoreBots: true,
	}
}

// HandlerBootstrap should never be called directly.
// HandlerBootstrap should instead be passed as an argument to discordgo.AddHandler in order to bootstrap the router
func (r *Router) HandlerBootstrap(discord *discordgo.Session, event *discordgo.MessageCreate) {
	command := NewCommand(r.Timeout)
	r.Handler(discord, event, command)
}

// Handler should never be called directly
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

// AddSubcommand adds another Router to the routes, where in "$x y z" y is a subcommand handled by the new router.
func (r *Router) AddSubcommand(command string, f func(r *Router)) {
	s := NewRouter("")
	r.routes[command] = s
	f(s)
}

// Use adds a middleware function to the chain.
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
