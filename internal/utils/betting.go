package utils

import "scrub/internal/errorutils"

func ValidateBetAmount(betAmount, credits uint64) error {
	if betAmount > credits {
		return errorutils.ErrInsufficientCredits
	}

	return nil
}
