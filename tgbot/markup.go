package tgbot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type ReplyMarkup struct {
	tgbotapi.InlineKeyboardMarkup
}

func (r *ReplyMarkup) copy() *ReplyMarkup {
	cp := *r

	if len(r.InlineKeyboard) > 0 {
		cp.InlineKeyboard = make([][]tgbotapi.InlineKeyboardButton, len(r.InlineKeyboard))
		for i, row := range r.InlineKeyboard {
			cp.InlineKeyboard[i] = make([]tgbotapi.InlineKeyboardButton, len(row))
			copy(cp.InlineKeyboard[i], row)
		}
	}

	return &cp
}

type Btn struct {
	tgbotapi.KeyboardButton
	Unique string
}

type ReplyButton struct {
	tgbotapi.KeyboardButton
}

type InlineButton struct {
	tgbotapi.InlineKeyboardButton
	Unique string
}
