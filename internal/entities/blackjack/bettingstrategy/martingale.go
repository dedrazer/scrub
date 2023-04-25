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

			err := m.decide(&players[i], j)
			if err != nil {
				return err
			}
		}
	}

	m.round++

	return nil
}

func (Martingale) GetName() string {
	return "martingale"
}

func (m *Martingale) decide(p *blackjack.BlackjackPlayer, handIndex int) error {
	switch *p.Hands[handIndex].Result {
	case utils.Loss, utils.SplitWon0:
		m.lose(p, handIndex)
	case utils.Win, utils.Blackjack, utils.SplitWon2, utils.Bankrupt:
		m.win(&p.Hands[handIndex])
	case utils.Push, utils.SplitWon1:
		return nil
	default:
		return fmt.Errorf("invalid result: %s", *p.Hands[handIndex].Result)
	}

	return nil
}

func (m *Martingale) lose(player *blackjack.BlackjackPlayer, handIndex int) {
	if player.Credits >= player.Hands[handIndex].BetAmount*2 {
		player.Hands[handIndex].BetAmount *= 2
		m.lossStreak++
		return
	}

	playerAllIn(m.Logger, player, handIndex, m.round, m.lossStreak, utils.Loss)
}

func (m *Martingale) win(hand *blackjack.Hand) {
	m.lossStreak = 0
	hand.BetAmount = m.OneCreditValue
}
