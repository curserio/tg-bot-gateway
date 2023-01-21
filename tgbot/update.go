package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

type Update struct {
	tgbotapi.Update
	*AdditionalUpdateParams
}

type AdditionalUpdateParams struct {
	Payload string
}

// ProcessUpdate processes a single incoming update.
// A started bot calls this function automatically.
func (b *Bot) ProcessUpdate(u Update) {
	c := b.NewContext(u)

	if u.Message != nil {
		m := &Message{Message: u.Message}

		if m.PinnedMessage != nil {
			b.handle(OnPinned, c)
			return
		}

		// Commands
		if m.Text != "" {
			// Filtering malicious messages
			if m.Text[0] == '\a' {
				return
			}

			match := cmdRx.FindAllStringSubmatch(m.Text, -1)
			if match != nil {
				// Syntax: "</command>@<bot> <payload>"
				command, botName := match[0][1], match[0][3]

				if botName != "" && !strings.EqualFold(b.Me.UserName, botName) {
					return
				}

				m.Payload = match[0][5]
				u.Payload = m.Payload
				if b.handle(command, c) {
					return
				}
			}

			// 1:1 satisfaction
			if b.handle(m.Text, c) {
				return
			}

			b.handle(OnText, c)
			return
		}

		if b.handleMedia(c) {
			return
		}

		if m.Contact != nil {
			b.handle(OnContact, c)
			return
		}
		if m.Location != nil {
			b.handle(OnLocation, c)
			return
		}
		if m.Venue != nil {
			b.handle(OnVenue, c)
			return
		}
		if m.Game != nil {
			b.handle(OnGame, c)
			return
		}
		if m.Invoice != nil {
			b.handle(OnInvoice, c)
			return
		}
	}

	if u.EditedMessage != nil {
		b.handle(OnEdited, c)
		return
	}

	if u.ChannelPost != nil {
		m := u.ChannelPost

		if m.PinnedMessage != nil {
			b.handle(OnPinned, c)
			return
		}

		b.handle(OnChannelPost, c)
		return
	}

	if u.EditedChannelPost != nil {
		b.handle(OnEditedChannelPost, c)
		return
	}

	if u.CallbackQuery != nil {
		callback := &Callback{CallbackQuery: u.CallbackQuery}
		if data := callback.Data; data != "" && data[0] == '\f' {
			match := cbackRx.FindAllStringSubmatch(data, -1)
			if match != nil {
				unique, payload := match[0][1], match[0][3]
				if handler, ok := b.handlers["\f"+unique]; ok {
					callback.Unique = unique
					callback.Data = payload
					u.Payload = payload
					b.runHandler(handler, c)
					return
				}
			}
		}

		b.handle(OnCallback, c)
		return
	}

	if u.ShippingQuery != nil {
		b.handle(OnShipping, c)
		return
	}

	if u.PreCheckoutQuery != nil {
		b.handle(OnCheckout, c)
		return
	}
}

func (b *Bot) handle(end string, c Context) bool {
	if handler, ok := b.handlers[end]; ok {
		b.runHandler(handler, c)
		return true
	}
	return false
}

func (b *Bot) handleMedia(c Context) bool {
	var (
		m     = c.Message()
		fired = true
	)

	switch {
	case m.Photo != nil:
		fired = b.handle(OnPhoto, c)
	case m.Voice != nil:
		fired = b.handle(OnVoice, c)
	case m.Audio != nil:
		fired = b.handle(OnAudio, c)
	case m.Animation != nil:
		fired = b.handle(OnAnimation, c)
	case m.Document != nil:
		fired = b.handle(OnDocument, c)
	case m.Sticker != nil:
		fired = b.handle(OnSticker, c)
	case m.Video != nil:
		fired = b.handle(OnVideo, c)
	case m.VideoNote != nil:
		fired = b.handle(OnVideoNote, c)
	default:
		return false
	}

	if !fired {
		return b.handle(OnMedia, c)
	}

	return true
}

func (b *Bot) runHandler(h HandlerFunc, c Context) {
	f := func() {
		if err := h(c); err != nil {
			b.OnError(err, c)
		}
	}
	if b.synchronous {
		f()
	} else {
		go f()
	}
}
