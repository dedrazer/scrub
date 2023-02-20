package blackjack

import (
	"scrub/internal/entities/deck"
	"scrub/internal/errors"
)

func (bj *Blackjack) DealRound(players []BlackjackPlayer) (dealerHand DealerHand, err error) {
	for i := range players {
		for j := range players[i].Hands {
			if len(players[i].Hands[j].cards) > 0 {
				players[i].Hands[j].cards = []deck.Card{}
			}
		}
	}

	// Deal first card to each player
	err = bj.DealRoundOfCards(players)
	if err != nil {
		return dealerHand, errors.ErrFailedSubMethod("DealRoundOfCards", err)
	}

	dealerHand.cards = make([]deck.Card, 2)

	// Deal down card to dealer
	dealerCardDown, err := bj.DealCard()
	if err != nil {
		return dealerHand, errors.ErrFailedSubMethod("DealCard (dealer down)", err)
	}

	dealerHand.cards[0] = *dealerCardDown

	// Deal second card to each player
	err = bj.DealRoundOfCards(players)
	if err != nil {
		return dealerHand, errors.ErrFailedSubMethod("DealRoundOfCards", err)
	}

	// Deal up card to dealer
	dealerCardUp, err := bj.DealCard()
	if err != nil {
		return dealerHand, errors.ErrFailedSubMethod("DealCard (dealer up)", err)
	}

	dealerHand.cards[1] = *dealerCardUp

	return dealerHand, err
}

func (bj *Blackjack) DealRoundOfCards(players []BlackjackPlayer) error {
	for i := range players {
		if len(players[i].Hands) == 0 {
			players[i].Hands = make([]Hand, 1)
		}
		for j := range players[i].Hands {
			card, err := bj.DealCard()
			if err != nil {
				return err
			}

			players[i].Hands[j].AddCard(*card)
		}
	}

	return nil
}
