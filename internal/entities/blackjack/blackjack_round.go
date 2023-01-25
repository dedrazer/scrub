package blackjack

import (
	"scrub/internal/entities/deck"
	"scrub/internal/errors"
)

func (bj *Blackjack) DealRound(betAmounts []uint64) (playerHands []Hand, dealerHand DealerHand, err error) {
	playerHands = make([]Hand, len(betAmounts))
	for i := range betAmounts {
		playerHands[i].betAmount = betAmounts[i]
	}

	// Deal first card to each player
	err = bj.DealRoundOfCards(&playerHands)
	if err != nil {
		return nil, dealerHand, errors.ErrFailedSubMethod("DealRoundOfCards", err)
	}

	dealerHand.cards = make([]deck.Card, 2)

	// Deal down card to dealer
	dealerCardDown, err := bj.DealCard()
	if err != nil {
		return nil, dealerHand, errors.ErrFailedSubMethod("DealCard (dealer down)", err)
	}

	dealerHand.cards[0] = *dealerCardDown

	// Deal second card to each player
	err = bj.DealRoundOfCards(&playerHands)
	if err != nil {
		return nil, dealerHand, errors.ErrFailedSubMethod("DealRoundOfCards", err)
	}

	// Deal up card to dealer
	dealerCardUp, err := bj.DealCard()
	if err != nil {
		return nil, dealerHand, errors.ErrFailedSubMethod("DealCard (dealer up)", err)
	}

	dealerHand.cards[1] = *dealerCardUp

	return playerHands, dealerHand, err
}

func (bj *Blackjack) DealRoundOfCards(playerHands *[]Hand) error {
	for i := range *playerHands {
		card, err := bj.DealCard()
		if err != nil {
			return err
		}

		(*playerHands)[i].AddCard(*card)
	}

	return nil
}
