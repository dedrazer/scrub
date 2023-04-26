package blackjack

import (
	"scrub/internal/entities/deck"
	"scrub/internal/errorutils"
)

func (bj *Blackjack) DealRound(players []BlackjackPlayer) (DealerHand, error) {
	resetHands(players)

	dealerHand, err := bj.dealCards(players)

	return dealerHand, err
}

func (bj *Blackjack) DealACardToEachPlayer(players []BlackjackPlayer) error {
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

func resetHands(players []BlackjackPlayer) {
	for i := range players {
		for j := range players[i].Hands {
			players[i].Hands[j].cards = []deck.Card{}
		}
	}
}

func (bj *Blackjack) dealCards(players []BlackjackPlayer) (DealerHand, error) {
	var dealerHand DealerHand

	err := bj.DealACardToEachPlayer(players)
	if err != nil {
		return dealerHand, errorutils.ErrFailedSubMethod("DealACardToEachPlayer", err)
	}

	dealerHand.cards = make([]deck.Card, 2)

	err = bj.dealCardToDealer(&dealerHand, true)
	if err != nil {
		return dealerHand, errorutils.ErrFailedSubMethod("dealDownCardToDealer", err)
	}

	err = bj.DealACardToEachPlayer(players)
	if err != nil {
		return dealerHand, errorutils.ErrFailedSubMethod("DealACardToEachPlayer", err)
	}

	err = bj.dealCardToDealer(&dealerHand, false)
	if err != nil {
		return dealerHand, errorutils.ErrFailedSubMethod("dealUpCardToDealer", err)
	}

	return dealerHand, nil
}

func (bj *Blackjack) dealCardToDealer(dealerHand *DealerHand, faceDown bool) error {
	card, err := bj.DealCard()
	if err != nil {
		return errorutils.ErrFailedSubMethod("DealCard (dealer down)", err)
	}

	if faceDown {
		dealerHand.cards[0] = *card
	} else {
		dealerHand.cards[1] = *card
	}

	return nil
}
