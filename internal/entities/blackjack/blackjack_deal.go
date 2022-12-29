package blackjack

import (
	"scrub/internal/entities/deck"
	"scrub/internal/errors"

	"go.uber.org/zap"
)

func (bj *Blackjack) DealCard(logger *zap.Logger) (*deck.Card, error) {
	card, err := bj.deck.TakeCardByIndex(0)
	if err != nil {
		return nil, errors.ErrFailedSubMethod("TakeCardByIndex", err)
	}

	card.Log(logger)

	return card, nil
}
