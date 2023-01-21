package tgbot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type User struct {
	*tgbotapi.User
}

func (u *User) ChatID() int {
	return u.ID
}

type Chat struct {
	*tgbotapi.Chat
}

func (c *Chat) ChatID() int {
	return int(c.ID)
}
