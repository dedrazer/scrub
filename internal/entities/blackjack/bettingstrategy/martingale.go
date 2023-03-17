package bettingstrategy

import (
	"fmt"
	"scrub/internal/entities/blackjack"
	"scrub/internal/entities/blackjack/utils"
)

type Martingale struct {
	CommonStrategyVariables
}

func (m *Martingale) Strategy(players []blackjack.BlackjackPlayer, oneCreditValue uint64) error {
	for i := range players {
		for j := range players[i].Hands {
			if players[i].Hands[j].Result == nil {
				continue
			}

			switch *players[i].Hands[j].Result {
			case utils.Loss, utils.SplitWon0:
				if players[i].Credits >= players[i].Hands[j].BetAmount*2 {
					players[i].Hands[j].BetAmount *= 2
					m.lossStreak++
					continue
				}

				playerAllIn(m.Logger, &players[i], j, m.round, m.lossStreak, utils.Loss)
			case utils.Win, utils.Blackjack, utils.SplitWon2, utils.Bankrupt:
				m.lossStreak = 0
				players[i].Hands[j].BetAmount = oneCreditValue
			case utils.Push, utils.SplitWon1:
				continue
			default:
				return fmt.Errorf("invalid result: %s", *players[i].Hands[j].Result)
			}
		}
	}
	
	m.round++

	return nil
}
