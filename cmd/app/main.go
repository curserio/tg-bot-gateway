package main

import (
	"flag"
	"github.com/curserio/tg-bot-gateway/config"
	"github.com/curserio/tg-bot-gateway/tgbot"
	tele "gopkg.in/telebot.v3"
	"log"
	"time"
)

func main() {
	var configFile string

	flag.StringVar(&configFile, "config", "config.yaml", "path to yaml config file")
	flag.StringVar(&configFile, "c", "config.yaml", "path to yaml config file")
	flag.Parse()

	// load config
	cfg := config.MustCreate(configFile)

	var appConfig config.AppConfig
	cfg.MustLoad("app", &appConfig)

	var telegramBotConfig tgbot.Config
	cfg.MustLoad("telegram_bot", &telegramBotConfig)

	// prepare telegram bot settings
	pref := tele.Settings{
		Token:     telegramBotConfig.Token,
		Poller:    &tele.LongPoller{Timeout: 10 * time.Second},
		Verbose:   telegramBotConfig.Verbose,
		ParseMode: tele.ModeHTML,
	}

	b, err := tgbot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	// prepare handlers
	b.Handle(tele.OnText, func(c tele.Context) error {
		return c.Reply("test")
	})

	log.Println("bot started. Ready to get messages")

	b.Start()

	defer b.Stop()
}
