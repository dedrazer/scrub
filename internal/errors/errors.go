package errors

import (
	"errors"
	"fmt"
)

const (
	errFailedSubMethod = "failed to %s: %w"
)

var (
	ErrActiveCardsIsEmpty = errors.New("No active cards remainin in the deck")
	ErrIndexOutOfRange    = errors.New("Index is out of range")
)

func ErrFailedSubMethod(methodName string, err error) error {
	return errors.New(fmt.Sprintf(errFailedSubMethod, methodName, err))
}
