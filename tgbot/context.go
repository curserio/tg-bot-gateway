package tgbot

import (
	"strings"
	"sync"
)

// HandlerFunc represents a handler function, which is
// used to handle actual endpoints.
type HandlerFunc func(Context) error

type Context interface {
	// Bot returns the bot instance.
	Bot() *Bot

	// Update returns the original update.
	Update() Update

	// Message returns stored message if such presented.
	Message() *Message

	// Callback returns stored callback if such presented.
	Callback() *Callback

	//// Query returns stored query if such presented.
	//Query() *Query
	//
	//// InlineResult returns stored inline result if such presented.
	//InlineResult() *InlineResult
	//
	//// ShippingQuery returns stored shipping query if such presented.
	//ShippingQuery() *ShippingQuery
	//
	//// PreCheckoutQuery returns stored pre checkout query if such presented.
	//PreCheckoutQuery() *PreCheckoutQuery
	//
	//// Poll returns stored poll if such presented.
	//Poll() *Poll
	//
	//// PollAnswer returns stored poll answer if such presented.
	//PollAnswer() *PollAnswer
	//
	//// ChatMember returns chat member changes.
	//ChatMember() *ChatMemberUpdate
	//
	//// ChatJoinRequest returns cha
	//ChatJoinRequest() *ChatJoinRequest
	//
	//// Migration returns both migration from and to chat IDs.
	//Migration() (int64, int64)

	// Sender returns the current recipient, depending on the context type.
	// Returns nil if user is not presented.
	Sender() *User

	// Chat returns the current chat, depending on the context type.
	// Returns nil if chat is not presented.
	Chat() *Chat

	// Recipient combines both Sender and Chat functions. If there is no user
	// the chat will be returned. The native context cannot be without sender,
	// but it is useful in the case when the context created intentionally
	// by the NewContext constructor and have only Chat field inside.
	Recipient() Recipient

	// Text returns the message text, depending on the context type.
	// In the case when no related data presented, returns an empty string.
	Text() string

	//// Entities returns the message entities, whether it's media caption's or the text's.
	//// In the case when no entities presented, returns a nil.
	//Entities() Entities

	// Data returns the current data, depending on the context type.
	// If the context contains command, returns its arguments string.
	// If the context contains payment, returns its payload.
	// In the case when no related data presented, returns an empty string.
	Data() string

	// Args returns a raw slice of command or callback arguments as strings.
	// The message arguments split by space, while the callback's ones by a "|" symbol.
	Args() []string

	// Send sends a message to the current recipient.
	// See Send from bot.go.
	Send(what interface{}, opts ...interface{}) error

	// Reply replies to the current message.
	// See Reply from bot.go.
	Reply(what interface{}, opts ...interface{}) error
	//
	//// Forward forwards the given message to the current recipient.
	//// See Forward from bot.go.
	//Forward(msg Editable, opts ...interface{}) error
	//
	//// ForwardTo forwards the current message to the given recipient.
	//// See Forward from bot.go
	//ForwardTo(to Recipient, opts ...interface{}) error
	//
	//// Edit edits the current message.
	//// See Edit from bot.go.
	//Edit(what interface{}, opts ...interface{}) error
	//
	//// EditCaption edits the caption of the current message.
	//// See EditCaption from bot.go.
	//EditCaption(caption string, opts ...interface{}) error
	//
	//// EditOrSend edits the current message if the update is callback,
	//// otherwise the content is sent to the chat as a separate message.
	//EditOrSend(what interface{}, opts ...interface{}) error
	//
	//// EditOrReply edits the current message if the update is callback,
	//// otherwise the content is replied as a separate message.
	//EditOrReply(what interface{}, opts ...interface{}) error
	//
	//// Delete removes the current message.
	//// See Delete from bot.go.
	//Delete() error
	//
	//// DeleteAfter waits for the duration to elapse and then removes the
	//// message. It handles an error automatically using b.OnError callback.
	//// It returns a Timer that can be used to cancel the call using its Stop method.
	//DeleteAfter(d time.Duration) *time.Timer
	//
	//// Notify updates the chat action for the current recipient.
	//// See Notify from bot.go.
	//Notify(action ChatAction) error
	//
	//// Ship replies to the current shipping query.
	//// See Ship from bot.go.
	//Ship(what ...interface{}) error
	//
	//// Accept finalizes the current deal.
	//// See Accept from bot.go.
	//Accept(errorMessage ...string) error
	//
	//// Answer sends a response to the current inline query.
	//// See Answer from bot.go.
	//Answer(resp *QueryResponse) error
	//
	//// Respond sends a response for the current callback query.
	//// See Respond from bot.go.
	//Respond(resp ...*CallbackResponse) error

	// Get retrieves data from the context.
	Get(key string) interface{}

	// Set saves data in the context.
	Set(key string, val interface{})
}

type tgContext struct {
	b     *Bot
	u     Update
	lock  sync.RWMutex
	store map[string]interface{}
}

func (c *tgContext) Bot() *Bot {
	return c.b
}

func (c *tgContext) Update() Update {
	return c.u
}

func (c *tgContext) Message() *Message {
	switch {
	case c.u.Message != nil:
		return &Message{Message: c.u.Message}
	case c.u.CallbackQuery != nil:
		return &Message{Message: c.u.CallbackQuery.Message}
	case c.u.EditedMessage != nil:
		return &Message{Message: c.u.EditedMessage}
	case c.u.ChannelPost != nil:
		if c.u.ChannelPost.PinnedMessage != nil {
			return &Message{Message: c.u.ChannelPost.PinnedMessage}
		}
		return &Message{Message: c.u.ChannelPost}
	case c.u.EditedChannelPost != nil:
		return &Message{Message: c.u.EditedChannelPost}
	default:
		return nil
	}
}

func (c *tgContext) Callback() *Callback {
	return &Callback{CallbackQuery: c.u.CallbackQuery}
}

func (c *tgContext) Sender() *User {
	switch {
	case c.u.CallbackQuery != nil:
		return &User{c.u.CallbackQuery.From}
	case c.Message() != nil:
		return &User{c.Message().From}
	default:
		return nil
	}
}

func (c *tgContext) Chat() *Chat {
	switch {
	case c.Message() != nil:
		return &Chat{c.Message().Chat}
	default:
		return nil
	}
}

func (c *tgContext) Recipient() Recipient {
	chat := c.Chat()
	if chat != nil {
		return chat
	}
	return c.Sender()
}

func (c *tgContext) Text() string {
	m := c.Message()
	if m == nil {
		return ""
	}
	if m.Caption != "" {
		return m.Caption
	}
	return m.Text
}

func (c *tgContext) Data() string {
	switch {
	case c.u.Message != nil:
		return c.u.Payload
	case c.u.CallbackQuery != nil:
		return c.u.Payload
	default:
		return ""
	}
}

func (c *tgContext) Args() []string {
	switch {
	case c.u.Message != nil:
		payload := strings.Trim(c.u.Payload, " ")
		if payload != "" {
			return strings.Split(payload, " ")
		}
	case c.u.CallbackQuery != nil:
		return strings.Split(c.u.CallbackQuery.Data, "|")
	}
	return nil
}

func (c *tgContext) Send(what interface{}, opts ...interface{}) error {
	_, err := c.b.Send(c.Recipient(), what, opts...)
	return err
}

func (c *tgContext) Reply(what interface{}, opts ...interface{}) error {
	msg := c.Message()
	if msg == nil {
		return ErrBadContext
	}
	_, err := c.b.Reply(msg, what, opts...)
	return err
}

func (c *tgContext) Get(key string) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.store[key]
}

func (c *tgContext) Set(key string, val interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.store == nil {
		c.store = make(map[string]interface{})
	}
	c.store[key] = val
}
