package blackjack

import (
	"errors"
	"scrub/internal/entities/blackjack/utils"
	"scrub/internal/errorutils"
)

func (bj *Blackjack) Results(players []BlackjackPlayer, dealerHand DealerHand) error {
	bj.logger.Debug("calculating results")

	err := bj.DrawRemainingDealerCards(&dealerHand)
	if err != nil {
		return errorutils.ErrFailedSubMethod("DrawRemainingDealerCards", err)
	}

	err = bj.setPlayerResults(players, dealerHand)
	if err != nil {
		return errorutils.ErrFailedSubMethod("setPlayerResults", err)
	}

	err = updateCreditBalance(players)
	if err != nil {
		return errorutils.ErrFailedSubMethod("updateCreditBalance", err)
	}

	return nil
}

func (bj *Blackjack) setPlayerResults(players []BlackjackPlayer, dealerHand DealerHand) error {
	dealerBust := false
	if dealerHand.Bust() {
		bj.logger.Debug("dealer bust")
		dealerBust = true
	}

	for i, p := range players {
		for j, h := range p.Hands {
			if h.Bust() {
				bj.ResultPlayerBust(&players[i].Hands[j])
				continue
			}

			if h.Blackjack() && !dealerHand.Blackjack() {
				bj.ResultPlayerBlackjack(&players[i].Hands[j])
				continue
			}

			if dealerBust {
				bj.ResultDealerBust(&players[i].Hands[j])
				continue
			}

			if h.UpperValue() < dealerHand.UpperValue() {
				bj.ResultDealerWins(&players[i].Hands[j])
				continue
			}

			if h.UpperValue() > dealerHand.UpperValue() {
				bj.ResultPlayerWins(&players[i].Hands[j])
				continue
			}

			if h.UpperValue() == dealerHand.UpperValue() {
				bj.ResultPush(&players[i].Hands[j])
				continue
			}

			return errors.New("unexpected case")
		}
	}

	return nil
}

func updateCreditBalance(players []BlackjackPlayer) error {
	splitWinCount := 0

	for i, p := range players {
		for j, h := range p.Hands {
			if h.Result == nil {
				return errors.New("unexpected nil result")
			}

			res := *h.Result

			switch res {
			case utils.Win:
				if h.isSplit {
					splitWinCount++
				}

				players[i].Win(h.BetAmount)
			case utils.Loss:
				err := players[i].Lose(h.BetAmount)
				if err != nil {
					return errorutils.ErrFailedSubMethod("Lose", err)
				}
				if players[i].Credits == 0 {
					players[i].Hands[j].Result = &utils.Bankrupt
				}
			case utils.Push:
				players[i].Win(0)
			case utils.Blackjack:
				players[i].Win(uint64(float64(h.BetAmount) * 1.5))
			}

			if h.isSplit {
				splitRes := utils.SplitWon[splitWinCount]
				players[i].Hands[j].Result = &splitRes
			}
		}
	}

	return nil
}
