package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/getynge/chatterbot/commands"
	"github.com/getynge/chatterbot/middleware"
	"github.com/getynge/chatterbot/routing"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func SetupRoutes(prefix string) *routing.Router {
	r := routing.NewRouter(prefix)

	r.Use(middleware.Logging)

	r.AddCommandFunc("echo", commands.Echo)

	return r
}

func main() {
	// configuring logging
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	// Configuring environment variables
	viper.SetEnvPrefix("CB")
	_ = viper.BindEnv("AuthenticationToken", "AUTHENTICATION_TOKEN")
	_ = viper.BindEnv("Prefix", "PREFIX")

	authToken := viper.GetString("AuthenticationToken")
	prefix := viper.GetString("Prefix")

	// Actual application logic starts here
	discord, err := discordgo.New("Bot " + authToken)

	if err != nil {
		zap.L().Panic("Could not authenticate with discord", zap.Error(err))
	}

	router := SetupRoutes(prefix)
	discord.AddHandler(router.HandlerBootstrap)

	if err = discord.Open(); err != nil {
		zap.L().Panic("could not open discord session", zap.Error(err))
	}

	zap.L().Info("Starting...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	err = discord.Close()

	if err != nil {
		zap.L().Error("Could not close discord connection", zap.Error(err))
	}

	zap.L().Info("Exiting...")
}
