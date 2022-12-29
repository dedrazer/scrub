package blackjack

import (
	"scrub/internal/entities/deck"
	"strings"

	"go.uber.org/zap"
)

type Hand struct {
	cards []deck.Card
}

func (h *Hand) Print() string {
	cardStrings := make([]string, len(h.cards))
	for i, c := range h.cards {
		cardStrings[i] = c.Print()
	}

	return strings.Join(cardStrings, ", ")
}

func (h *Hand) Log(logger *zap.Logger) {
	if len(h.cards) == 0 {
		logger.Info("hand", zap.String("cards", "empty"))
		return
	}

	logger.Info("hand", zap.String("cards", h.Print()))
}

func (h *Hand) DealerLog(logger *zap.Logger) {
	logger.Info("dealer hand", zap.String("card", h.cards[1].Print()))
}

func (h *Hand) Value() []uint {
	var total uint

	containsAce := false
	for _, c := range h.cards {
		if c.Value == 1 {
			containsAce = true
		}

		total += c.Value
	}

	if total < 12 && containsAce {
		return []uint{total, total + 10}
	}

	return []uint{total}
}

func (h *Hand) AddCard(c deck.Card) {
	if h.cards == nil {
		h.cards = []deck.Card{}
	}
	h.cards = append(h.cards, c)
}
