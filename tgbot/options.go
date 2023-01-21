package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Option is a shortcut flag type for certain message features
// (so-called options). It means that instead of passing
// fully-fledged SendOptions* to Send(), you can use these
// flags instead.
//
// Supported options are defined as iota-constants.
type Option int

const (
	// NoPreview = SendOptions.DisableWebPagePreview
	NoPreview Option = iota

	// Silent = SendOptions.DisableNotification
	Silent
)

// SendOptions has most complete control over in what way the message
// must be sent, providing an API-complete set of custom properties
// and options.
//
// Despite its power, SendOptions is rather inconvenient to use all
// the way through bot logic, so you might want to consider storing
// and re-using it somewhere or be using Option flags instead.
type SendOptions struct {
	// If the message is a reply, original message.
	ReplyTo *Message

	// See ReplyMarkup struct definition.
	ReplyMarkup *ReplyMarkup

	// For text messages, disables previews for links in this message.
	DisableWebPagePreview bool

	// Sends the message silently. iOS users will not receive a notification, Android users will receive a notification with no sound.
	DisableNotification bool

	// ParseMode controls how client apps render your message.
	ParseMode ParseMode
}

func (og *SendOptions) copy() *SendOptions {
	cp := *og
	if cp.ReplyMarkup != nil {
		cp.ReplyMarkup = cp.ReplyMarkup.copy()
	}
	return &cp
}

func extractOptions(how []interface{}) *SendOptions {
	opts := &SendOptions{}

	for _, prop := range how {
		switch opt := prop.(type) {
		case *SendOptions:
			opts = opt.copy()
		case *ReplyMarkup:
			if opt != nil {
				opts.ReplyMarkup = opt.copy()
			}
		case Option:
			switch opt {
			case NoPreview:
				opts.DisableWebPagePreview = true
			case Silent:
				opts.DisableNotification = true
			default:
				panic("telebot: unsupported flag-option")
			}
		case ParseMode:
			opts.ParseMode = opt
		default:
			panic("telebot: unsupported send-option")
		}
	}

	return opts
}

func (b *Bot) embedSendOptions(msg *tgbotapi.MessageConfig, opt *SendOptions) {
	if b.parseMode != ModeDefault {
		msg.ParseMode = b.parseMode
	}

	if opt == nil {
		return
	}

	if opt.ReplyTo != nil && opt.ReplyTo.MessageID != 0 {
		msg.ReplyToMessageID = opt.ReplyTo.MessageID
	}

	if opt.DisableWebPagePreview {
		msg.DisableWebPagePreview = true
	}

	if opt.DisableNotification {
		msg.DisableNotification = true
	}

	if opt.ParseMode != ModeDefault {
		msg.ParseMode = opt.ParseMode
	}

	if opt.ReplyMarkup != nil {
		inlineKeyboard := make([][]InlineButton, len(opt.ReplyMarkup.InlineKeyboard))

		for i, buttons := range opt.ReplyMarkup.InlineKeyboard {
			inlineKeyboardButtons := make([]InlineButton, len(buttons))

			for j, button := range buttons {
				inlineKeyboardButtons[j] = InlineButton{
					InlineKeyboardButton: button,
					Unique:               "",
				}
			}

			inlineKeyboard[i] = inlineKeyboardButtons
		}
		processButtons(inlineKeyboard)
		msg.ReplyMarkup = inlineKeyboard
	}
}

func (b *Bot) embedSendFileOptions(msg *tgbotapi.DocumentConfig, opt *SendOptions) {
	if b.parseMode != ModeDefault {
		msg.ParseMode = b.parseMode
	}

	if opt == nil {
		return
	}

	if opt.ReplyTo != nil && opt.ReplyTo.MessageID != 0 {
		msg.ReplyToMessageID = opt.ReplyTo.MessageID
	}

	if opt.DisableNotification {
		msg.DisableNotification = true
	}

	if opt.ParseMode != ModeDefault {
		msg.ParseMode = opt.ParseMode
	}

	if opt.ReplyMarkup != nil {
		inlineKeyboard := make([][]InlineButton, len(opt.ReplyMarkup.InlineKeyboard))

		for i, buttons := range opt.ReplyMarkup.InlineKeyboard {
			inlineKeyboardButtons := make([]InlineButton, len(buttons))

			for j, button := range buttons {
				inlineKeyboardButtons[j] = InlineButton{
					InlineKeyboardButton: button,
					Unique:               "",
				}
			}

			inlineKeyboard[i] = inlineKeyboardButtons
		}
		processButtons(inlineKeyboard)
		msg.ReplyMarkup = inlineKeyboard
	}
}

func processButtons(keys [][]InlineButton) {
	if keys == nil || len(keys) < 1 || len(keys[0]) < 1 {
		return
	}

	for i := range keys {
		for j := range keys[i] {
			key := &keys[i][j]
			if key.Unique != "" {
				callbackData := *key.CallbackData
				// Format: "\f<callback_name>|<data>"
				data := callbackData
				if data == "" {
					callbackData = "\f" + key.Unique
				} else {
					callbackData = "\f" + key.Unique + "|" + data
				}
				key.CallbackData = &callbackData
			}
		}
	}
}
