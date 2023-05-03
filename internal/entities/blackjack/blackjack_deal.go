package blackjack

import (
	"errors"
	"scrub/internal/entities/deck"
	"scrub/internal/errorutils"

	"go.uber.org/zap"
)

func (bj *Blackjack) DealCard() (*deck.Card, error) {
	card, err := bj.deck.TakeCardByIndex(0)
	if err != nil {
		if bj.deck.IsFinished() {
			bj.deck.Shuffle()
			return bj.DealCard()
		}

		return nil, errorutils.ErrFailedSubMethod("TakeCardByIndex", err)
	}

	return card, nil
}

// When the dealer has served every player, the dealers face-down card is turned up.
// If the total is 17 or more, the dealer must stand. If the total is 16 or under, they must take a card.
// The dealer must continue to take cards until the total is 17 or more, at which point the dealer must stand.
// If the dealer has an ace, and counting it as 11 would bring the total to 17 or more (but not over 21),
// the dealer must count the ace as 11 and stand.
func (bj *Blackjack) DrawRemainingDealerCards(dh *DealerHand) error {
	if dh == nil {
		return errors.New("dealer hand is nil")
	}

	if dh.hasNoValue() {
		return errors.New("no value")
	}

	if dh.shouldDraw() {
		c, err := bj.DealCard()
		if err != nil {
			bj.logger.Error("failed to deal card", zap.Error(err))
		}

		dh.AddCard(*c)
		return bj.DrawRemainingDealerCards(dh)
	}

	dh.DealerResult(bj.logger)
	return nil
}
