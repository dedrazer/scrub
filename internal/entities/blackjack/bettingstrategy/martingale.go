package bettingstrategy

import (
	"fmt"
	"scrub/internal/entities/blackjack"
	"scrub/internal/entities/blackjack/utils"
)

func Martingale(players []blackjack.BlackjackPlayer, oneCreditValue uint64) error {
	for i := range players {
		for j := range players[i].Hands {
			if players[i].Hands[j].Result == nil {
				continue
			}

			switch *players[i].Hands[j].Result {
			case utils.Loss:
				players[i].Hands[j].BetAmount *= 2
			case utils.Win, utils.Blackjack:
				players[i].Hands[j].BetAmount = oneCreditValue
			case utils.Push:
				continue
			default:
				return fmt.Errorf("invalid result: %s", *players[i].Hands[j].Result)
			}
		}
	}

	return nil
}
