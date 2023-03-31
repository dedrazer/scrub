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
	finalDeck := deck.NewShuffledDecks(numberOfDecks)

	return &Blackjack{
		deck:          &finalDeck,
		numberOfDecks: numberOfDecks,
	}
}
