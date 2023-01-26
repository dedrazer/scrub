package blackjack

import (
	"scrub/internal/entities/deck"
	"scrub/internal/entities/player"
	"scrub/internal/errors"
)

func (bj *Blackjack) DealRound(playerBets []player.PlayerBet) (players []BlackJackPlayer, dealerHand DealerHand, err error) {
	players = make([]BlackJackPlayer, len(playerBets))
	for i := range players {
		players[i].PlayerBet = playerBets[i]
	}

	// Deal first card to each player
	err = bj.DealRoundOfCards(players)
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
	err = bj.DealRoundOfCards(players)
	if err != nil {
		return nil, dealerHand, errors.ErrFailedSubMethod("DealRoundOfCards", err)
	}

	// Deal up card to dealer
	dealerCardUp, err := bj.DealCard()
	if err != nil {
		return nil, dealerHand, errors.ErrFailedSubMethod("DealCard (dealer up)", err)
	}

	dealerHand.cards[1] = *dealerCardUp

	return players, dealerHand, err
}

func (bj *Blackjack) DealRoundOfCards(players []BlackJackPlayer) error {
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
