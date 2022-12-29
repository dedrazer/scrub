package blackjack

import "scrub/internal/entities/deck"

type Blackjack struct {
	deck          *deck.Deck
	numberOfDecks uint
}

func NewBlackjack(numberOfDecks uint) *Blackjack {
	var finalDeck deck.Deck

	for i := uint(0); i < numberOfDecks; i++ {
		d := deck.NewDeck()

		finalDeck = deck.Merge(&finalDeck, &d)
	}

	finalDeck.Shuffle()

	return &Blackjack{
		deck:          &finalDeck,
		numberOfDecks: numberOfDecks,
	}
}