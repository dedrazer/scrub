package blackjack

import (
	"scrub/internal/entities/deck"
	"scrub/internal/errors"
)

func (bj *Blackjack) DealCard() (*deck.Card, error) {
	card, err := bj.deck.TakeCardByIndex(0)
	if err != nil {
		return nil, errors.ErrFailedSubMethod("TakeCardByIndex", err)
	}

	return card, nil
}
