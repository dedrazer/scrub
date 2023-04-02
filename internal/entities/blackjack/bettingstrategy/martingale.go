package bettingstrategy

import (
	"fmt"
	"scrub/internal/entities/blackjack"
	"scrub/internal/entities/blackjack/utils"
)

type Martingale struct {
	CommonStrategyVariables
}

func (m *Martingale) BettingStrategy(players []blackjack.BlackjackPlayer) error {
	for i := range players {
		for j := range players[i].Hands {
			if players[i].Hands[j].Result == nil {
				continue
			}

			switch *players[i].Hands[j].Result {
			case utils.Loss, utils.SplitWon0:
				m.lose(&players[i], j)
			case utils.Win, utils.Blackjack, utils.SplitWon2, utils.Bankrupt:
				m.win(&players[i].Hands[j])
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

func (s *Martingale) GetName() string {
	return "martingale"
}

func (m *Martingale) lose(player *blackjack.BlackjackPlayer, handNumber int) {
	if player.Credits >= player.Hands[handNumber].BetAmount*2 {
		player.Hands[handNumber].BetAmount *= 2
		m.lossStreak++
		return
	}

	playerAllIn(m.Logger, player, handNumber, m.round, m.lossStreak, utils.Loss)
}

func (m *Martingale) win(hand *blackjack.Hand) {
	m.lossStreak = 0
	hand.BetAmount = m.OneCreditValue
}
