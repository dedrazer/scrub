package blackjack

import (
	"scrub/internal/entities/deck"
	"scrub/internal/errorutils"

	"go.uber.org/zap"
)

var (
	acceptedInputs = map[string][]string{
		"first": {"hit", "double", "stand"},
		"hit":   {"hit", "stand"},
	}

	first = "first"

	stand  = "stand"
	split  = "split"
	hit    = "hit"
	double = "double"
)

func (bj *Blackjack) Play(players []BlackjackPlayer, dealerHand DealerHand) error {
	if dealerHand.Blackjack() {
		bj.logger.Debug("dealer has blackjack")

		bj.checkIfPlayersHaveBlackJack(players)
	} else {
		bj.logger.Debug("playing round")
		err := bj.playRound(players, dealerHand)
		if err != nil {
			return err
		}
	}

	if err := bj.Results(bj.logger, players, dealerHand); err != nil {
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
	for i, p := range players {
		for j := range p.Hands {
			if players[i].Hands[j].BetAmount > players[i].Credits {
				return errorutils.ErrInsufficientCredits
			}

			bj.logger.Debug("bet info", zap.Uint64("amount", players[i].Hands[j].BetAmount))

			if players[i].Hands[j].CanSplit(players[i].Credits) {
				action := bj.strategy(p.Hands[j], dealerHand, players[i].Credits)

				if action == split {
					bj.logger.Debug("splitting hand", zap.Int("player", i+1), zap.Int("hand", j+1))
					if p.Hands[j].cards[0].Symbol == p.Hands[j].cards[1].Symbol {
						splitHand := Hand{
							cards:     []deck.Card{p.Hands[j].cards[1]},
							isSplit:   true,
							BetAmount: p.Hands[j].BetAmount,
						}
						players[i].Hands[j].cards = p.Hands[j].cards[:1]
						players[i].Hands[j].isSplit = true
						players[i].Hands = append(p.Hands, splitHand)

						bj.SplitCount++
					} else {
						return errorutils.ErrCannotSplit
					}
				}
			}
		}

		for j := range players[i].Hands {
			bj.logger.Debug("turn", zap.Int("player", i+1), zap.Int("hand", j+1))

			var action string
			kind := hit

			dealerHand.DealerLog(bj.logger)
			for action != stand && !players[i].Hands[j].Bust() && players[i].Hands[j].UpperValue() != 21 {
				players[i].Hands[j].Log(bj.logger)

				kind = hit
				if len(players[i].Hands[j].cards) == 2 && !players[i].Hands[j].isSplit {
					kind = first

					if players[i].Hands[j].UpperValue() == 21 {
						bj.logger.Debug("player has blackjack")
						break
					}
				}

				action = bj.strategy(players[i].Hands[j], dealerHand, players[i].Credits)
				validInput := false
				for _, v := range acceptedInputs[kind] {
					if action == v {
						validInput = true
						break
					}
				}

				if !validInput {
					return errorutils.ErrInvalidInput(action)
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

					players[i].Hands[j].Double()

					players[i].Hands[j].AddCard(*c)
					players[i].Hands[j].Log(bj.logger)

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

					players[i].Hands[j].AddCard(*c)
					c.Log(bj.logger)
				}
			}

			if players[i].Hands[j].Bust() {
				bj.logger.Debug("player bust")
			}
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
