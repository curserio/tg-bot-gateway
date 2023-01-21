package tgbot

import "log"

var defaultOnError = func(err error, c Context) {
	if c != nil {
		log.Println(c.Update().UpdateID, err)
	} else {
		log.Println(err)
	}
}

func (b *Bot) OnError(err error, c Context) {
	b.onError(err, c)
}

func (b *Bot) debug(err error) {
	if b.verbose {
		b.OnError(err, nil)
	}
}
