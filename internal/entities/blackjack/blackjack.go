package blackjack

import (
	"scrub/internal/entities/deck"
)

type Blackjack struct {
	deck          *deck.Deck
	numberOfDecks uint
	BlackjackStatistics
}

type BlackjackStatistics struct {
	PlayerWins           uint64
	PlayerLosses         uint64
	Pushes               uint64
	DealerBust           uint64
	PlayerBust           uint64
	SplitCount           uint64
	PlayerBlackjackCount uint64
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
