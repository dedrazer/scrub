package bettingstrategy

import (
	"scrub/internal/entities/blackjack"
	"scrub/internal/entities/blackjack/utils"
)

type Stern struct {
	level      int
	lossStreak int
	CommonStrategyVariables
}

func (s *Stern) BettingStrategy(players []blackjack.BlackjackPlayer) error {
	for i := range players {
		for j := range players[i].Hands {
			if players[i].Hands[j].Result == nil {
				continue
			}

			s.decide(&players[i], j)
		}
	}

	s.round++

	return nil
}

func (Stern) GetName() string {
	return "stern"
}

func (s *Stern) decide(p *blackjack.BlackjackPlayer, handIndex int) {
	switch *p.Hands[handIndex].Result {
	case utils.Win, utils.Blackjack, utils.SplitWon2, utils.Bankrupt:
		s.win(p, handIndex)
	case utils.Loss, utils.SplitWon0:
		s.lose(p, handIndex)
	case utils.Push, utils.SplitWon1:
		return
	}
}

func (s *Stern) lose(player *blackjack.BlackjackPlayer, handNumber int) {
	s.lossStreak++
	s.winStreak = 0
	if (s.level == 0 && s.lossStreak == 3) || (s.level > 0 && s.lossStreak == 2) {
		if player.Credits < player.Hands[handNumber].BetAmount*2 {
			playerAllIn(s.Logger, player, handNumber, s.round, s.lossStreak, utils.Loss)
		} else {
			player.Hands[handNumber].BetAmount *= 2
		}

		s.level++
		s.lossStreak = 0
	}
}

func (s *Stern) win(player *blackjack.BlackjackPlayer, handNumber int) {
	s.winStreak++

	// end of cycle
	if s.winStreak == 2 {
		s.endCycle(&player.Hands[handNumber])
	}

	if s.winStreak == 1 {
		if player.Credits >= player.Hands[handNumber].BetAmount*2 {
			player.Hands[handNumber].BetAmount *= 2
		} else {
			playerAllIn(s.Logger, player, handNumber, s.round, s.lossStreak, utils.Win)
		}
	}

	s.lossStreak = 0
}

func (s *Stern) endCycle(hand *blackjack.Hand) {
	s.winStreak = 0
	s.level = 0
	hand.BetAmount = s.OneCreditValue
}
