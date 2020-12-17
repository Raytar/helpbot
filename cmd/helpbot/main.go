package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Raytar/helpbot"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	botName = "helpbot"
	cfgFile = ".helpbotrc"
)

var log = &logrus.Logger{
	Out:       os.Stderr,
	Formatter: new(logrus.TextFormatter),
	Hooks:     make(logrus.LevelHooks),
	Level:     logrus.InfoLevel,
}

var ag *helpbot.Autograder

func main() {
	err := initConfig()
	if err != nil {
		log.Fatalln("Failed to read config:", err)
	}

	if viper.GetBool("autograder") {
		ag, err = helpbot.NewAutograder(viper.GetInt("autograder-user-id"))
		if err != nil {
			log.Fatalln("Failed to init autograder:", err)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var bots []*helpbot.HelpBot

	for _, c := range cfg {
		bot, err := helpbot.New(c, log, ag)
		if err != nil {
			log.Fatalf("Failed to initialize bot: %v", err)
		}
		err = bot.Connect(ctx)
		if err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
		bots = append(bots, bot)
	}

	// run until interrupted
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	// cleanup
	if ag != nil {
		for _, b := range bots {
			b.Disconnect()
		}
		cancel()
		ag.Close()
	}
}
