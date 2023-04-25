package blackjack

import (
	"scrub/internal/entities/deck"

	"go.uber.org/zap"
)

type Blackjack struct {
	logger        *zap.Logger
	strategy      func(playerHand Hand, dealerHand DealerHand, playerCredits uint64) string
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

func NewBlackjack(logger *zap.Logger, numberOfDecks uint) *Blackjack {
	finalDeck := deck.NewShuffledDecks(numberOfDecks)

	return &Blackjack{
		logger:        logger,
		strategy:      PlayingStrategy,
		deck:          &finalDeck,
		numberOfDecks: numberOfDecks,
	}
}
