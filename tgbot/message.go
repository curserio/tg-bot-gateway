package tgbot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type Message struct {
	*tgbotapi.Message
	Payload string
}
