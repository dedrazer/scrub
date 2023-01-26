package errors

import (
	"errors"
	"fmt"
)

const (
	errFailedSubMethod = "failed to %s: %w"
)

var (
	ErrActiveCardsIsEmpty  = errors.New("No active cards remainin in the deck")
	ErrIndexOutOfRange     = errors.New("Index is out of range")
	ErrInsufficientCredits = errors.New("Insufficient credits")
	ErrUnexpectedNil       = errors.New("Unexpected nil value")
)

func ErrFailedSubMethod(methodName string, err error) error {
	return fmt.Errorf(errFailedSubMethod, methodName, err)
}
