package bettingstrategy

import (
	"scrub/internal/entities/blackjack"
	"scrub/internal/entities/blackjack/utils"

	"go.uber.org/zap"
)

var (
	sternLevel      int
	sternLossStreak int
)

func Stern(logger *zap.Logger, players []blackjack.BlackjackPlayer, oneCreditValue uint64) error {
	for i := range players {
		for j := range players[i].Hands {
			if players[i].Hands[j].Result == nil {
				continue
			}

			switch *players[i].Hands[j].Result {
			case utils.Win, utils.Blackjack, utils.SplitWon2, utils.Bankrupt:
				winStreak++

				if winStreak == 1 {
					if players[i].Credits >= players[i].Hands[j].BetAmount*2 {
						players[i].Hands[j].BetAmount *= 2
						continue
					}

					playerAllIn(logger, &players[i], j)
				}

				// end of cycle
				if winStreak == 2 {
					winStreak = 0
					sternLevel = 0
					players[i].Hands[j].BetAmount = oneCreditValue
				}
			case utils.Push, utils.SplitWon1:
				continue
			case utils.Loss, utils.SplitWon0:
				sternLossStreak++
				if (sternLevel == 0 && sternLossStreak == 3) || (sternLevel > 0 && sternLossStreak == 2) {
					sternLevel++
					sternLossStreak = 0

					if players[i].Credits >= players[i].Hands[j].BetAmount*2 {
						players[i].Hands[j].BetAmount *= 2
						continue
					}

					playerAllIn(logger, &players[i], j)
				}
			}
		}
	}

	return nil
}
