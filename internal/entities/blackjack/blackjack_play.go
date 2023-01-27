package blackjack

import (
	"fmt"
	"scrub/internal/entities/deck"
	"scrub/internal/errors"
	"strings"

	"go.uber.org/zap"
)

var (
	acceptedInputs = map[string][]string{
		"first": {"h", "d", "s"},
		"hit":   {"h", "s"},
	}

	prompts = map[string]string{
		"first": "hit/double down/stand? (h/d/s)",
		"hit":   "hit/stand? (h/s)",
	}
)

func (bj *Blackjack) Play(logger *zap.Logger, players []BlackJackPlayer, dealerHand DealerHand) error {
	if dealerHand.Blackjack() {
		logger.Info("dealer has blackjack")
	} else {
		logger.Info("playing round")
		for i, p := range players {
			for j := range p.Hands {
				if p.Hands[j].cards[0].Symbol == p.Hands[j].cards[1].Symbol {
					p.Hands[j].Log(logger)
					fmt.Println("Split? (y/n)")
					var input string
					_, err := fmt.Scanln(&input)
					if err != nil {
						return errors.ErrFailedSubMethod("fmt.Scanln", err)
					}

					input = strings.ToLower(input)
					if input == "y" {
						splitHand := Hand{
							cards:   []deck.Card{p.Hands[j].cards[1]},
							isSplit: true,
						}
						players[i].Hands[j].cards = p.Hands[j].cards[:1]
						players[i].Hands[j].isSplit = true
						players[i].Hands = append(p.Hands, splitHand)
					}
				}
			}

			for j := range players[i].Hands {
				logger.Info("turn", zap.Int("player", i+1), zap.Int("hand", j+1))

				var input string
				kind := "hit"

				dealerHand.DealerLog(logger)
				for input != "s" && !players[i].Hands[j].Bust() && players[i].Hands[j].UpperValue() != 21 {
					players[i].Hands[j].Log(logger)

					kind = "hit"
					if len(players[i].Hands[j].cards) == 2 && !players[i].Hands[j].isSplit {
						kind = "first"

						if players[i].Hands[j].UpperValue() == 21 {
							logger.Info("player has blackjack")
							break
						}
					}

					fmt.Println(prompts[kind])
					_, err := fmt.Scanln(&input)
					if err != nil {
						return errors.ErrFailedSubMethod("fmt.Scanln", err)
					}

					input = strings.ToLower(input)
					validInput := false
					for _, v := range acceptedInputs[kind] {
						if input == v {
							validInput = true
							break
						}
					}

					if !validInput {
						return errors.ErrInvalidInput
					}

					if input == "d" {
						var c *deck.Card
						c, err = bj.DealCard()
						if err != nil {
							return errors.ErrFailedSubMethod("DealCard", err)
						}

						players[i].Hands[j].betAmount *= 2

						players[i].Hands[j].AddCard(*c)
						players[i].Hands[j].Log(logger)

						input = "n"
						break
					}

					if input == "h" {
						var c *deck.Card
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
