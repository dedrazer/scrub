package blackjack

import (
	"scrub/internal/entities/player"

	"go.uber.org/zap"
)

func Demo(logger *zap.Logger) {
	logger.Debug("initialising blackjack")
	testBlackjack := NewBlackjack(6)

	players := []BlackjackPlayer{
		{
			Player: player.Player{
				Name:    "Martin",
				Credits: 1000,
			},
			Hands: []Hand{
				{
					BetAmount: 50,
				},
				{
					BetAmount: 100,
				},
			},
		},
		{
			Player: player.Player{
				Name:    "Fran",
				Credits: 1000,
			},
			Hands: []Hand{
				{
					BetAmount: 25,
				},
			},
		},
	}

	logger.Debug("dealing round", zap.Any("players", players))
	dealerHand, err := testBlackjack.DealRound(players)
	if err != nil {
		logger.Fatal("failed to deal round", zap.Error(err))
	}

	dealerHand.DealerLog(logger)

	err = testBlackjack.Play(logger, players, dealerHand, PlayingStrategy)
	if err != nil {
		logger.Fatal("unexpected error", zap.Error(err))
	}
}
