package utils

import "scrub/internal/errorutils"

var acceptedInputs = map[string][]string{
	"first": {"hit", "double", "stand"},
	"hit":   {"hit", "stand"},
}

func ValidateInput(kind, action string) error {
	for _, v := range acceptedInputs[kind] {
		if action == v {
			return nil
		}
	}

	return errorutils.ErrInvalidInput(action)
}
