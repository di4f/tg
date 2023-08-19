package tg

import (
	"errors"
	"fmt"
)

type WrongUpdateType struct {
	Type string
}

var (
	ScreenNotExistErr    = errors.New("screen does not exist")
	SessionNotExistErr   = errors.New("session does not exist")
	KeyboardNotExistErr  = errors.New("keyboard does not exist")
	NotAvailableErr      = errors.New("the context is not available")
	EmptyKeyboardTextErr = errors.New("got empty text for a keyboard")
	ActionNotDefinedErr  = errors.New("action was not defined")
)

func (wut WrongUpdateType) Error() string {
	if wut.Type == "" {
		return "wrong update type"
	}
	return fmt.Sprintf("wrong update type '%s'", wut.Type)
}