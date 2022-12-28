package blackjack

import (
	"scrub/internal/errors"

	"go.uber.org/zap"
)

func (bj *Blackjack) DealRound(logger *zap.Logger) error {
	card, err := bj.deck.TakeCardByIndex(0)
	if err != nil {
		return errors.ErrFailedSubMethod("TakeCardByIndex", err)
	}

	card.Log(logger)

	return nil
}
