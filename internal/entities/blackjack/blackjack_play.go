package blackjack

import (
	bjutils "scrub/internal/entities/blackjack/utils"
	"scrub/internal/entities/deck"
	"scrub/internal/errorutils"
	"scrub/internal/utils"

	"go.uber.org/zap"
)

var (
	first = "first"

	stand  = "stand"
	split  = "split"
	hit    = "hit"
	double = "double"
)

func (bj *Blackjack) Play(players []BlackjackPlayer, dealerHand DealerHand) error {
	if dealerHand.Blackjack() {
		bj.logger.Debug("dealer has blackjack")

		err := bj.checkIfPlayersHaveBlackJack(players)
		if err != nil {
			return err
		}
	} else {
		bj.logger.Debug("playing round")
		err := bj.playRound(players, dealerHand)
		if err != nil {
			return err
		}
	}

	if err := bj.Results(players, dealerHand); err != nil {
		return errorutils.ErrFailedSubMethod("Results", err)
	}

	if err := PrintAllResults(bj.logger, players); err != nil {
		return errorutils.ErrFailedSubMethod("PrintAllResults", err)
	}

	for i := range players {
		if err := players[i].ResetHands(); err != nil {
			return errorutils.ErrFailedSubMethod("ResetHands", err)
		}
	}

	return nil
}

func (bj *Blackjack) playRound(players []BlackjackPlayer, dealerHand DealerHand) error {
	for i := range players {
		err := bj.handlePotentialSplits(&players[i], dealerHand)
		if err != nil {
			return errorutils.ErrFailedSubMethod("handlePotentialSplits", err)
		}

		err = bj.playHands(&players[i], dealerHand)
		if err != nil {
			return errorutils.ErrFailedSubMethod("playHands", err)
		}
	}

	return nil
}

func (bj *Blackjack) playHands(p *BlackjackPlayer, dealerHand DealerHand) error {
	for handIndex := range p.Hands {
		bj.logger.Debug("turn", zap.Any("player", p), zap.Int("hand", handIndex+1))

		var action string
		kind := hit

		dealerHand.DealerLog(bj.logger)
		for action != stand && !p.Hands[handIndex].Bust() && p.Hands[handIndex].UpperValue() != 21 {
			p.Hands[handIndex].Log(bj.logger)

			kind = hit
			if p.Hands[handIndex].IsUnplayed() {
				kind = first
			}

			action = bj.strategy(p.Hands[handIndex], dealerHand, p.Credits)

			err := bjutils.ValidateInput(kind, action)
			if err != nil {
				return err
			}

			if action == double {
				var (
					c   *deck.Card
					err error
				)
				c, err = bj.DealCard()
				if err != nil {
					return errorutils.ErrFailedSubMethod("DealCard", err)
				}

				p.Hands[handIndex].Double()

				p.Hands[handIndex].AddCard(*c)
				p.Hands[handIndex].Log(bj.logger)

				action = stand
				break
			}

			if action == hit {
				var (
					c   *deck.Card
					err error
				)
				c, err = bj.DealCard()
				if err != nil {
					return errorutils.ErrFailedSubMethod("DealCard", err)
				}

				p.Hands[handIndex].AddCard(*c)
				c.Log(bj.logger)
			}
		}

		if p.Hands[handIndex].Bust() {
			bj.logger.Debug("player bust")
		}
	}

	return nil
}

func (bj *Blackjack) checkIfPlayersHaveBlackJack(players []BlackjackPlayer) error {
	for i := range players {
		for j := range players[i].Hands {
			if players[i].Hands[j].UpperValue() >= 10 {
				var (
					c   *deck.Card
					err error
				)
				c, err = bj.DealCard()
				if err != nil {
					return errorutils.ErrFailedSubMethod("DealCard", err)
				}

				players[i].Hands[j].AddCard(*c)
				c.Log(bj.logger)
			}
		}
	}

	return nil
}

func (bj *Blackjack) handlePotentialSplits(p *BlackjackPlayer, dealerHand DealerHand) error {
	for handIndex := range p.Hands {
		err := utils.ValidateBetAmount(p.Hands[handIndex].BetAmount, p.Credits)
		if err != nil {
			return err
		}

		bj.logger.Debug("bet info", zap.Uint64("amount", p.Hands[handIndex].BetAmount))

		if p.Hands[handIndex].CanSplit(p.Credits) {
			action := bj.strategy(p.Hands[handIndex], dealerHand, p.Credits)

			if action == split {
				err = bj.split(p, handIndex)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
