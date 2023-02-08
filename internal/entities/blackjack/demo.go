package blackjack

import (
	"scrub/internal/entities/player"

	"go.uber.org/zap"
)

func Demo(logger *zap.Logger) {
	logger.Info("initialising blackjack")
	testBlackjack := NewBlackjack(6)

	players := []BlackJackPlayer{
		{
			Player: player.Player{
				Name:    "Martin",
				Credits: 1000,
			},
			Hands: []Hand{
				{
					betAmount: 50,
				},
				{
					betAmount: 100,
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
					betAmount: 25,
				},
			},
		},
	}

	logger.Info("dealing round", zap.Any("players", players))
	dealerHand, err := testBlackjack.DealRound(players)
	if err != nil {
		logger.Fatal("failed to deal round", zap.Error(err))
	}

	dealerHand.DealerLog(logger)

	err = testBlackjack.Play(logger, players, dealerHand, Strategy1)
	if err != nil {
		logger.Fatal("unexpected error", zap.Error(err))
	}
}
