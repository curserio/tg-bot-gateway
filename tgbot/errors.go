package tgbot

import (
	"errors"
	"fmt"
)

var (
	ErrBadRecipient    = errors.New("tgbot: recipient is nil")
	ErrUnsupportedWhat = errors.New("tgbot: unsupported what argument")
	ErrCouldNotUpdate  = errors.New("tgbot: could not fetch new updates")
	ErrTrueResult      = errors.New("tgbot: result is True")
	ErrBadContext      = errors.New("tgbot: context does not contain message")
)

// wrapError returns new wrapped error
func wrapError(err error) error {
	return fmt.Errorf("tgbot: %w", err)
}
