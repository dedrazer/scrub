package errors

import (
	"errors"
	"fmt"
)

const (
	errFailedSubMethod = "Failed to %s: %w"
	errInvalidInput    = "Invalid input: %s"
)

var (
	ErrActiveCardsIsEmpty  = errors.New("No active cards remainin in the deck")
	ErrIndexOutOfRange     = errors.New("Index is out of range")
	ErrInsufficientCredits = errors.New("Insufficient credits")
	ErrUnexpectedNil       = errors.New("Unexpected nil value")
	ErrCannotSplit         = errors.New("Cannot split non-pair")
)

func ErrFailedSubMethod(methodName string, err error) error {
	return fmt.Errorf(errFailedSubMethod, methodName, err)
}

func ErrInvalidInput(input string) error {
	return fmt.Errorf(errInvalidInput, input)
}
