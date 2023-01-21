package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"time"
)

func (b *Bot) getUpdates(offset, limit int, timeout time.Duration) ([]Update, error) {
	u := tgbotapi.NewUpdate(offset)
	u.Limit = limit
	u.Timeout = int(timeout.Seconds())

	apiUpdates, err := b.api.GetUpdates(u)
	if err != nil {
		return nil, wrapError(err)
	}

	updates := make([]Update, len(apiUpdates))

	for i, apiUpdate := range apiUpdates {
		updates[i] = Update{Update: apiUpdate, AdditionalUpdateParams: &AdditionalUpdateParams{}}
	}

	return updates, nil
}

func (b *Bot) sendText(to Recipient, text string, opt *SendOptions) (*Message, error) {
	msg := tgbotapi.NewMessage(int64(to.ChatID()), text)

	b.embedSendOptions(&msg, opt)

	resp, err := b.api.Send(msg)

	return &Message{Message: &resp}, err
}
