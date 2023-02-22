package blackjack

import (
	"scrub/internal/entities/deck"
	"scrub/internal/errors"

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

func (bj *Blackjack) Play(logger *zap.Logger, players []BlackjackPlayer, dealerHand DealerHand, strategy func(playerHand Hand, dealerHand DealerHand, playerCredits uint64) string) error {
	if dealerHand.Blackjack() {
		logger.Debug("dealer has blackjack")
	} else {
		logger.Debug("playing round")
		for i, p := range players {
			for j := range p.Hands {
				if players[i].Hands[j].BetAmount > players[i].Credits {
					return errors.ErrInsufficientCredits
				}

				logger.Debug("bet info", zap.Uint64("amount", players[i].Hands[j].BetAmount))

				if players[i].Hands[j].CanSplit(players[i].Credits) {
					action := strategy(p.Hands[j], dealerHand, players[i].Credits)

					if action == split {
						logger.Debug("splitting hand", zap.Int("player", i+1), zap.Int("hand", j+1))
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
							return errors.ErrCannotSplit
						}
					}
				}
			}

			for j := range players[i].Hands {
				logger.Debug("turn", zap.Int("player", i+1), zap.Int("hand", j+1))

				var action string
				kind := hit

				dealerHand.DealerLog(logger)
				for action != stand && !players[i].Hands[j].Bust() && players[i].Hands[j].UpperValue() != 21 {
					players[i].Hands[j].Log(logger)

					kind = hit
					if len(players[i].Hands[j].cards) == 2 && !players[i].Hands[j].isSplit {
						kind = first

						if players[i].Hands[j].UpperValue() == 21 {
							logger.Debug("player has blackjack")
							break
						}
					}

					action = strategy(players[i].Hands[j], dealerHand, players[i].Credits)
					validInput := false
					for _, v := range acceptedInputs[kind] {
						if action == v {
							validInput = true
							break
						}
					}

					if !validInput {
						return errors.ErrInvalidInput(action)
					}

					if action == double {
						var (
							c   *deck.Card
							err error
						)
						c, err = bj.DealCard()
						if err != nil {
							return errors.ErrFailedSubMethod("DealCard", err)
						}

						players[i].Hands[j].Double()

						players[i].Hands[j].AddCard(*c)
						players[i].Hands[j].Log(logger)

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
							return errors.ErrFailedSubMethod("DealCard", err)
						}

						players[i].Hands[j].AddCard(*c)
						c.Log(logger)
					}
				}

				if players[i].Hands[j].Bust() {
					logger.Debug("player bust")
				}
			}
		}
	}

	if err := bj.Results(logger, players, dealerHand); err != nil {
		return errors.ErrFailedSubMethod("Results", err)
	}

	if err := PrintAllResults(logger, players); err != nil {
		return errors.ErrFailedSubMethod("PrintAllResults", err)
	}

	for i := range players {
		if err := players[i].ResetHands(); err != nil {
			return errors.ErrFailedSubMethod("ResetHands", err)
		}
	}

	return nil
}
