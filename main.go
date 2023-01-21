package main

import (
	"encoding/json"
	"flag"
	"github.com/curserio/tg-bot-gateway/tgbot"
	"log"
	"os"
	"time"
)

type Config struct {
	Token          string  `json:"token"`
	Verbose        bool    `json:"verbose"`
	UpdateInterval int     `json:"update_interval"`
	Admins         []int64 `json:"admins"`
}

func loadConfig(configFile string) (config Config) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalln("cannot read config:", err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalln("unmarshal config:", err)
	}

	return
}

func main() {
	var configFile string

	flag.StringVar(&configFile, "config", "config.json", "path to json config file")
	flag.StringVar(&configFile, "c", "config.json", "path to json config file")
	flag.Parse()

	config := loadConfig(configFile)
	config.Verbose = true

	pref := tgbot.Settings{
		Token:     config.Token,
		Poller:    &tgbot.LongPoller{Timeout: 10 * time.Second},
		Verbose:   config.Verbose,
		ParseMode: tgbot.ModeHTML,
	}

	b, err := tgbot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle(tgbot.OnText, func(c tgbot.Context) error {
		return c.Reply("test")
	})

	log.Println("bot started. Ready to get messages")

	b.Start()

	defer b.Stop()
}
