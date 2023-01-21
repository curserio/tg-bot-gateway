package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
	"regexp"
	"time"
)

// NewBot does try to build a Bot with token `token`, which
// is a secret API key assigned to particular bot.
func NewBot(pref Settings) (*Bot, error) {
	if pref.Updates == 0 {
		pref.Updates = 100
	}

	client := pref.Client
	if client == nil {
		client = &http.Client{Timeout: time.Minute}
	}

	if pref.Poller == nil {
		pref.Poller = &LongPoller{}
	}
	if pref.OnError == nil {
		pref.OnError = defaultOnError
	}

	bot := &Bot{
		Token:   pref.Token,
		Poller:  pref.Poller,
		onError: pref.OnError,

		Updates:  make(chan Update, pref.Updates),
		handlers: make(map[string]HandlerFunc),
		stop:     make(chan chan struct{}),

		synchronous: pref.Synchronous,
		verbose:     pref.Verbose,
		parseMode:   pref.ParseMode,
		client:      client,
	}

	botApi, err := tgbotapi.NewBotAPI(bot.Token)
	if err != nil {
		return nil, err
	}

	botApi.Debug = pref.Verbose

	bot.api = botApi

	if pref.Offline {
		bot.Me = &User{}
	} else {
		apiUser, err := bot.api.GetMe()
		if err != nil {
			return nil, err
		}
		bot.Me = &User{&apiUser}
	}

	return bot, nil
}

type Bot struct {
	api *tgbotapi.BotAPI

	Me      *User
	Token   string
	URL     string
	Updates chan Update
	Poller  Poller
	onError func(error, Context)

	handlers    map[string]HandlerFunc
	synchronous bool
	verbose     bool
	parseMode   ParseMode
	stop        chan chan struct{}
	client      *http.Client
	stopClient  chan struct{}
}

// Settings represents a utility struct for passing certain
// properties of a bot around and is required to make bots.
type Settings struct {
	Token string

	// Updates channel capacity, defaulted to 100.
	Updates int

	// Poller is the provider of Updates.
	Poller Poller

	// Synchronous prevents handlers from running in parallel.
	// It makes ProcessUpdate return after the handler is finished.
	Synchronous bool

	// Verbose forces bot to log all upcoming requests.
	// Use for debugging purposes only.
	Verbose bool

	// ParseMode used to set default parse mode of all sent messages.
	// It attaches to every send, edit or whatever method. You also
	// will be able to override the default mode by passing a new one.
	ParseMode ParseMode

	// OnError is a callback function that will get called on errors
	// resulted from the handler. It is used as post-middleware function.
	// Notice that context can be nil.
	OnError func(error, Context)

	// HTTP Client used to make requests to telegram api
	Client *http.Client

	// Offline allows to create a bot without network for testing purposes.
	Offline bool
}

var (
	cmdRx   = regexp.MustCompile(`^(/\w+)(@(\w+))?(\s|$)(.+)?`)
	cbackRx = regexp.MustCompile(`^\f([-\w]+)(\|(.+))?$`)
)

// Handle lets you set the handler for some command name or
// one of the supported endpoints. It also applies middleware
// if such passed to the function.
//
// Example:
//
//	b.Handle("/start", func (c tele.Context) error {
//		return c.Reply("Hello!")
//	})
//
//	b.Handle(&inlineButton, func (c tele.Context) error {
//		return c.Respond(&tele.CallbackResponse{Text: "Hello!"})
//	})
//
// Middleware usage:
//
//	b.Handle("/ban", onBan, middleware.Whitelist(ids...))
func (b *Bot) Handle(endpoint interface{}, h HandlerFunc) {
	handler := func(c Context) error {
		return h(c)
	}

	switch end := endpoint.(type) {
	case string:
		b.handlers[end] = handler
	case CallbackEndpoint:
		b.handlers[end.CallbackUnique()] = handler
	default:
		panic("telebot: unsupported endpoint")
	}
}

// Start brings bot into motion by consuming incoming
// updates (see Bot.Updates channel).
func (b *Bot) Start() {
	if b.Poller == nil {
		panic("tgbot: can't start without a poller")
	}

	// do nothing if called twice
	if b.stopClient != nil {
		return
	}
	b.stopClient = make(chan struct{})

	stop := make(chan struct{})
	stopConfirm := make(chan struct{})

	go func() {
		b.Poller.Poll(b, b.Updates, stop)
		close(stopConfirm)
	}()

	for {
		select {
		// handle incoming updates
		case upd := <-b.Updates:
			b.ProcessUpdate(upd)
			// call to stop polling
		case confirm := <-b.stop:
			close(stop)
			<-stopConfirm
			close(confirm)
			b.stopClient = nil
			return
		}
	}
}

// Stop gracefully shuts the poller down.
func (b *Bot) Stop() {
	if b.stopClient != nil {
		close(b.stopClient)
	}
	confirm := make(chan struct{})
	b.stop <- confirm
	<-confirm
}

// NewContext returns a new native context object,
// field by the passed update.
func (b *Bot) NewContext(u Update) Context {
	return &tgContext{
		b: b,
		u: u,
	}
}

// Send accepts 2+ arguments, starting with destination chat, followed by
// some Sendable (or string!) and optional send options.
//
// NOTE:
//
//	Since most arguments are of type interface{}, but have pointer
//	method receivers, make sure to pass them by-pointer, NOT by-value.
//
// What is a send option exactly? It can be one of the following types:
//
//   - *SendOptions (the actual object accepted by Telegram API)
//   - *ReplyMarkup (a component of SendOptions)
//   - Option (a shortcut flag for popular options)
//   - ParseMode (HTML, Markdown, etc)
func (b *Bot) Send(to Recipient, what interface{}, opts ...interface{}) (*Message, error) {
	if to == nil {
		return nil, ErrBadRecipient
	}

	sendOpts := extractOptions(opts)

	switch object := what.(type) {
	case string:
		return b.sendText(to, object, sendOpts)
	case Sendable:
		return object.Send(b, to, sendOpts)
	default:
		return nil, ErrUnsupportedWhat
	}
}

// Reply behaves just like Send() with an exception of "reply-to" indicator.
// This function will panic upon nil Message.
func (b *Bot) Reply(to *Message, what interface{}, opts ...interface{}) (*Message, error) {
	sendOpts := extractOptions(opts)
	if sendOpts == nil {
		sendOpts = &SendOptions{}
	}

	sendOpts.ReplyTo = to
	return b.Send(&Chat{to.Chat}, what, sendOpts)
}
