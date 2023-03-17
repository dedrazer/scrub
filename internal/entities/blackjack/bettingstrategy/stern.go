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

func (s *Stern) Strategy(players []blackjack.BlackjackPlayer) error {
	for i := range players {
		for j := range players[i].Hands {
			if players[i].Hands[j].Result == nil {
				continue
			}

			switch *players[i].Hands[j].Result {
			case utils.Win, utils.Blackjack, utils.SplitWon2, utils.Bankrupt:
				s.winStreak++

				// end of cycle
				if s.winStreak == 2 {
					s.winStreak = 0
					s.level = 0
					players[i].Hands[j].BetAmount = s.OneCreditValue
				}

				if s.winStreak == 1 {
					if players[i].Credits >= players[i].Hands[j].BetAmount*2 {
						players[i].Hands[j].BetAmount *= 2
					} else {
						playerAllIn(s.Logger, &players[i], j, s.round, s.lossStreak, utils.Win)
					}
				}

				s.lossStreak = 0
			case utils.Push, utils.SplitWon1:
				continue
			case utils.Loss, utils.SplitWon0:
				s.lossStreak++
				s.winStreak = 0
				if (s.level == 0 && s.lossStreak == 3) || (s.level > 0 && s.lossStreak == 2) {
					if players[i].Credits < players[i].Hands[j].BetAmount*2 {
						playerAllIn(s.Logger, &players[i], j, s.round, s.lossStreak, utils.Loss)
					} else {
						players[i].Hands[j].BetAmount *= 2
					}

					s.level++
					s.lossStreak = 0
				}
			}
		}
	}

	s.round++

	return nil
}

func (s *Stern) GetName() string {
	return "stern"
}
