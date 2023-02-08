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

	//prompts = map[string]string{
	//	"first": "hit/double down/stand? (h/d/s)",
	//	"hit":   "hit/stand? (h/s)",
	//}

	first = "first"

	stand  = "stand"
	split  = "split"
	hit    = "hit"
	double = "double"
)

func (bj *Blackjack) Play(logger *zap.Logger, players []BlackJackPlayer, dealerHand DealerHand, strategy func(playerHand Hand, dealerHand DealerHand) string) error {
	if dealerHand.Blackjack() {
		logger.Info("dealer has blackjack")
	} else {
		logger.Info("playing round")
		for i, p := range players {
			for j := range p.Hands {
				if p.Hands[j].CanSplit() {
					action := strategy(p.Hands[j], dealerHand)

					if action == split {
						if p.Hands[j].cards[0].Symbol == p.Hands[j].cards[1].Symbol {
							splitHand := Hand{
								cards:   []deck.Card{p.Hands[j].cards[1]},
								isSplit: true,
							}
							players[i].Hands[j].cards = p.Hands[j].cards[:1]
							players[i].Hands[j].isSplit = true
							players[i].Hands = append(p.Hands, splitHand)
						} else {
							return errors.ErrCannotSplit
						}
					}
				}
			}

			for j := range players[i].Hands {
				logger.Info("turn", zap.Int("player", i+1), zap.Int("hand", j+1))

				var action string
				kind := hit

				dealerHand.DealerLog(logger)
				for action != stand && !players[i].Hands[j].Bust() && players[i].Hands[j].UpperValue() != 21 {
					players[i].Hands[j].Log(logger)

					kind = hit
					if len(players[i].Hands[j].cards) == 2 && !players[i].Hands[j].isSplit {
						kind = first

						if players[i].Hands[j].UpperValue() == 21 {
							logger.Info("player has blackjack")
							break
						}
					}

					//fmt.Println(prompts[kind])
					//_, err := fmt.Scanln(&input)
					//if err != nil {
					//	return errors.ErrFailedSubMethod("fmt.Scanln", err)
					//}

					//input = strings.ToLower(input)

					action = strategy(players[i].Hands[j], dealerHand)
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

						players[i].Hands[j].betAmount *= 2

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
					logger.Info("player bust")
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

	return nil
}
