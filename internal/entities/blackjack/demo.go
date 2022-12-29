package blackjack

import "go.uber.org/zap"

func Demo(logger *zap.Logger) {
	logger.Info("initialising blackjack")
	testBlackjack := NewBlackjack(6)

	var numberOfHands uint8 = 1

	logger.Info("dealing round", zap.Uint8("numberOfHands", numberOfHands))
	playerHands, dealerHand, err := testBlackjack.DealRound(numberOfHands)
	if err != nil {
		logger.Fatal("failed to deal round", zap.Error(err))
	}

	for i, ph := range playerHands {
		logger.Info("player hand", zap.Int("player", i+1))
		ph.Log(logger)
	}

	dealerHand.DealerLog(logger)
}
