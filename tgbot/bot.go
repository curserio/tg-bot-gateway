package tgbot

import (
	tbot "gopkg.in/telebot.v3"
)

// NewBot does try to build a Bot
func NewBot(pref tbot.Settings) (*tbot.Bot, error) {
	return tbot.NewBot(pref)
}
