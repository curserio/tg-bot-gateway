package tgbot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type File struct {
	tgbotapi.File
	filename string
	data     []byte
}
