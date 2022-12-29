package blackjack

import (
	"errors"
	"scrub/internal/entities/deck"
	internal "scrub/internal/errors"

	"go.uber.org/zap"
)

func (bj *Blackjack) DealCard() (*deck.Card, error) {
	card, err := bj.deck.TakeCardByIndex(0)
	if err != nil {
		return nil, internal.ErrFailedSubMethod("TakeCardByIndex", err)
	}

	return card, nil
}

// When the dealer has served every player, the dealers face-down card is turned up.
// If the total is 17 or more, it must stand. If the total is 16 or under, they must take a card.
// The dealer must continue to take cards until the total is 17 or more, at which point the dealer must stand.
// If the dealer has an ace, and counting it as 11 would bring the total to 17 or more (but not over 21),
// the dealer must count the ace as 11 and stand. The dealer's decisions, then, are automatic on all plays,
// whereas the player always has the option of taking one or more cards.
func (bj *Blackjack) DrawDealerCards(logger *zap.Logger, h *Hand) error {
	value := h.Value()
	if len(value) == 0 {
		return errors.New("no value")
	}

	if len(value) == 2 {
		if value[1] < 17 {
			c, err := bj.DealCard()
			if err != nil {
				logger.Error("failed to deal card", zap.Error(err))
			}

			h.AddCard(*c)
			return bj.DrawDealerCards(logger, h)
		}

		return nil
	}

	if value[0] < 17 {
		c, err := bj.DealCard()
		if err != nil {
			logger.Error("failed to deal card", zap.Error(err))
		}

		h.AddCard(*c)
		return bj.DrawDealerCards(logger, h)
	}

	h.Log(logger)
	return nil
}
