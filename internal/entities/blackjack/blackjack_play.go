package blackjack

import (
	"fmt"
	"scrub/internal/entities/deck"
	"scrub/internal/errors"
	"strings"

	"go.uber.org/zap"
)

func (bj *Blackjack) Play(logger *zap.Logger, players []BlackJackPlayer, dealerHand DealerHand) error {
	logger.Info("playing round")
	for i, p := range players {
		for j, h := range p.Hands {
			logger.Info("turn", zap.Int("player", i+1), zap.Int("hand", j+1))

			input := "y"

			dealerHand.DealerLog(logger)
			for input == "y" && !h.Bust() {
				h.Log(logger)

				if len(h.cards) == 2 {
					var doubleDownInput string
					fmt.Println("Double down? (y/N)")
					_, err := fmt.Scanln(&doubleDownInput)
					if err != nil {
						return errors.ErrFailedSubMethod("fmt.Scanln", err)
					}
					if strings.ToLower(doubleDownInput) == "y" {
						var c *deck.Card
						c, err = bj.DealCard()
						if err != nil {
							return errors.ErrFailedSubMethod("DealCard", err)
						}

						// todo: one BetAmount per hand
						players[i].PlayerBet.BetAmount *= 2

						players[i].Hands[j].AddCard(*c)
						players[i].Hands[j].Log(logger)
						
						input = "n"
						break
					}
				}

				fmt.Println("Take card? (y/N)")
				_, err := fmt.Scanln(&input)
				if err != nil {
					return errors.ErrFailedSubMethod("fmt.Scanln", err)
				}

				if input == "y" {
					var c *deck.Card
					c, err = bj.DealCard()
					if err != nil {
						return errors.ErrFailedSubMethod("DealCard", err)
					}

					players[i].Hands[j].AddCard(*c)
					c.Log(logger)
				}
			}

			if h.Bust() {
				logger.Info("player bust")
			}
		}
	}

	if err := bj.Results(logger, players, dealerHand); err != nil {
		return errors.ErrFailedSubMethod("Results", err)
	}

	for _, p := range players {
		if err := p.PrintResult(logger); err != nil {
			return errors.ErrFailedSubMethod("PrintResult", err)
		}
	}

	return nil
}
