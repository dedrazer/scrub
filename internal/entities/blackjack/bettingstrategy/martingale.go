package bettingstrategy

import (
	"fmt"
	"scrub/internal/entities/blackjack"
	"scrub/internal/entities/blackjack/utils"

	"go.uber.org/zap"
)

var (
	lossStreak int
	round      int
)

func Martingale(logger *zap.Logger, players []blackjack.BlackjackPlayer, oneCreditValue uint64) error {
	for i := range players {
		for j := range players[i].Hands {
			if players[i].Hands[j].Result == nil {
				continue
			}

			switch *players[i].Hands[j].Result {
			case utils.Loss:
				players[i].Hands[j].BetAmount *= 2
				lossStreak++
			case utils.Win, utils.Blackjack:
				lossStreak = 0
				players[i].Hands[j].BetAmount = oneCreditValue
			case utils.Push:
				continue
			default:
				return fmt.Errorf("invalid result: %s", *players[i].Hands[j].Result)
			}

			if players[i].Hands[j].BetAmount > players[i].Credits {
				players[i].Hands[j].BetAmount = players[i].Credits
				logger.Info("player all in", zap.Int("loss streak", lossStreak), zap.Uint64("next bet", players[i].Hands[j].BetAmount), zap.Int("round", round))
			}
		}

		round++
	}

	return nil
}
