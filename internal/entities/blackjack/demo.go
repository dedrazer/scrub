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
	}

	logger.Info("dealing round", zap.Any("playerBets", playerBets))
	playerHands, dealerHand, err := testBlackjack.DealRound(playerBets)
	if err != nil {
		logger.Fatal("failed to deal round", zap.Error(err))
	}

	dealerHand.DealerLog(logger)

	for i, ph := range playerHands {
		for j := range ph.Hands {
			logger.Info("player hand", zap.Int("player", i+1), zap.Int("hand", j+1))
			ph.Hands[j].Log(logger)
		}
	}
	err = testBlackjack.Play(logger, playerHands, dealerHand)
	if err != nil {
		logger.Fatal("unexpected error", zap.Error(err))
	}
}
