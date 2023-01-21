package tgbot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

// CallbackEndpoint is an interface any element capable
// of responding to a callback `\f<unique>`.
type CallbackEndpoint interface {
	CallbackUnique() string
}

type Callback struct {
	*tgbotapi.CallbackQuery
	// Unique displays an unique of the button from which the
	// callback was fired. Sets immediately before the handling,
	// while the Data field stores only with payload.
	Unique string
}

// CallbackUnique returns ReplyButton.Text.
func (t *ReplyButton) CallbackUnique() string {
	return t.Text
}

// CallbackUnique returns InlineButton.Unique.
func (t *InlineButton) CallbackUnique() string {
	return "\f" + t.Unique
}

// CallbackUnique implements CallbackEndpoint.
func (t *Btn) CallbackUnique() string {
	if t.Unique != "" {
		return "\f" + t.Unique
	}
	return t.Text
}
