package routing_test

import (
	"github.com/bwmarrin/discordgo"
	"github.com/getynge/chatterbot/routing"
)

// blank functions for use in examples
func userHandler(_ *discordgo.Session, _ *discordgo.MessageCreate, _ *routing.Command) {}

func ExampleRouter_AddSubcommand() {
	router := routing.NewRouter("$")

	router.AddSubcommand("add", func(r *routing.Router) {
		r.AddCommandFunc("user", userHandler)
	})
}

func ExampleRouter_AddSubcommand_nestedRoutes() {
	router := routing.NewRouter("$")

	router.AddSubcommand("add", func(r *routing.Router) {
		r.AddSubcommand("users", func(r2 *routing.Router) {
			r2.AddCommandFunc("named", userHandler)
		})
	})
}
