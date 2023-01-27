package blackjack

import (
	"scrub/internal/entities/player"

	"go.uber.org/zap"
)

func Demo(logger *zap.Logger) {
	logger.Info("initialising blackjack")
	testBlackjack := NewBlackjack(6)

	playerBets := []player.PlayerBet{
		{
			Player: player.Player{
				Name:    "Martin",
				Credits: 1000,
			},
			BetAmount: 50,
		},
		{
			Player: player.Player{
				Name:    "Fran",
				Credits: 1000,
			},
			BetAmount: 25,
		},
	}

	logger.Info("dealing round", zap.Any("playerBets", playerBets))
	players, dealerHand, err := testBlackjack.DealRound(playerBets)
	if err != nil {
		logger.Fatal("failed to deal round", zap.Error(err))
	}

	dealerHand.DealerLog(logger)

	err = testBlackjack.Play(logger, players, dealerHand)
	if err != nil {
		logger.Fatal("unexpected error", zap.Error(err))
	}
}
