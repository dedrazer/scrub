package errors

import "errors"

var (
	ErrActiveCardsIsEmpty = errors.New("No active cards remainin in the deck")
)
