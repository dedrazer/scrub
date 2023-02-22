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
			case utils.Loss, utils.SplitWon0:
				if players[i].Credits >= players[i].Hands[j].BetAmount*2 {
					players[i].Hands[j].BetAmount *= 2
					lossStreak++
					continue
				}

				players[i].Hands[j].BetAmount = players[i].Credits
				logger.Debug("player all in", zap.Int("loss streak", lossStreak), zap.Uint64("next bet", players[i].Hands[j].BetAmount), zap.Int("round", round))
			case utils.Win, utils.Blackjack, utils.SplitWon2, utils.Bankrupt:
				lossStreak = 0
				players[i].Hands[j].BetAmount = oneCreditValue
			case utils.Push, utils.SplitWon1:
				continue
			default:
				return fmt.Errorf("invalid result: %s", *players[i].Hands[j].Result)
			}
		}

		round++
	}

	return nil
}
