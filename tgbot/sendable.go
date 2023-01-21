package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Recipient is any possible endpoint you can send
// messages to: either user, group or a channel
type Recipient interface {
	ChatID() int
}

// Sendable is any object that can send itself.
//
// This is pretty cool, since it lets bots implement
// custom Sendables for complex kind of media or
// chat objects spanning across multiple messages.
type Sendable interface {
	Send(*Bot, Recipient, *SendOptions) (*Message, error)
}

// Send delivers media through bot b to recipient.
func (d *Document) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	f := tgbotapi.FileBytes{
		Name:  d.File.filename,
		Bytes: d.File.data,
	}

	documentConfig := tgbotapi.NewDocumentUpload(int64(to.ChatID()), f)

	b.embedSendFileOptions(&documentConfig, opt)

	if d.FileSize != 0 {
		documentConfig.FileSize = d.FileSize
	}

	if d.Caption != "" {
		documentConfig.Caption = d.Caption
	}

	resp, err := b.api.Send(documentConfig)

	return &Message{Message: &resp}, err
}
