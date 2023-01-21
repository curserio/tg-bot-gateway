package tgbot

import (
	"time"
)

// Poller is a provider of Updates.
//
// All pollers must implement Poll(), which accepts bot
// pointer and subscription channel and start polling
// synchronously straight away.
type Poller interface {
	// Poll is supposed to take the bot object
	// subscription channel and start polling
	// for Updates immediately.
	//
	// Poller must listen for stop constantly and close
	// it as soon as it's done polling.
	Poll(b *Bot, updates chan Update, stop chan struct{})
}

// LongPoller is a classic LongPoller with timeout.
type LongPoller struct {
	Limit        int
	Timeout      time.Duration
	LastUpdateID int
}

// Poll does long polling
func (p *LongPoller) Poll(b *Bot, dest chan Update, stop chan struct{}) {
	for {
		select {
		case <-stop:
			return
		default:
		}

		updates, err := b.getUpdates(p.LastUpdateID+1, p.Limit, p.Timeout)
		if err != nil {
			b.debug(err)
			continue
		}

		for _, update := range updates {
			p.LastUpdateID = update.UpdateID
			dest <- update
		}
	}
}
