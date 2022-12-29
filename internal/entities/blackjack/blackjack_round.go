package blackjack

import (
	"scrub/internal/entities/deck"
	"scrub/internal/errors"

	"go.uber.org/zap"
)

func (bj *Blackjack) DealRound(logger *zap.Logger, numberOfHands uint8) (playerCards map[uint8][]deck.Card, dealerCards []deck.Card, err error) {
	playerCards = make(map[uint8][]deck.Card, numberOfHands)
	for i := uint8(0); i < numberOfHands; i++ {
		playerCards[i] = nil
	}

	// Deal first card to each player
	err = bj.DealRoundOfCards(logger, playerCards)
	if err != nil {
		return nil, nil, errors.ErrFailedSubMethod("DealRoundOfCards", err)
	}

	dealerCards = make([]deck.Card, 2)

	// Deal down card to dealer
	dealerCardDown, err := bj.DealCard()
	if err != nil {
		return nil, nil, errors.ErrFailedSubMethod("DealCard (dealer down)", err)
	}

	dealerCards[0] = *dealerCardDown

	// Deal second card to each player
	err = bj.DealRoundOfCards(logger, playerCards)
	if err != nil {
		return nil, nil, errors.ErrFailedSubMethod("DealRoundOfCards", err)
	}

	// Deal up card to dealer
	dealerCardUp, err := bj.DealCard()
	if err != nil {
		return nil, nil, errors.ErrFailedSubMethod("DealCard (dealer up)", err)
	}

	dealerCards[1] = *dealerCardUp

	return playerCards, dealerCards, err
}

func (bj *Blackjack) DealRoundOfCards(logger *zap.Logger, playerCards map[uint8][]deck.Card) error {
	for k := range playerCards {
		card, err := bj.DealCard()
		if err != nil {
			return err
		}

		playerCards[k] = append(playerCards[k], *card)
	}

	return nil
}
