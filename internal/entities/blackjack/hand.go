package blackjack

import (
	"fmt"
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

	logger.Info("hand", zap.String("cards", h.Print()), zap.String("value", h.PrintValue()))
}

func (h *Hand) DealerLog(logger *zap.Logger) {
	logger.Info("dealer hand", zap.String("card", h.cards[1].Print()))
}

func (h *Hand) PrintValue() string {
	values := h.Value()

	stringValues := make([]string, len(values))
	for i, v := range values {
		stringValues[i] = fmt.Sprintf("%d", v)
	}

	return strings.Join(stringValues, "/")
}

func (h *Hand) LogValue(logger *zap.Logger) {
	logger.Info("hand value", zap.String("amount", h.PrintValue()))
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

func (h *Hand) UpperValue() uint {
	values := h.Value()
	return values[len(values)-1]
}

func (h *Hand) Bust() bool {
	return h.UpperValue() > 21
}

func (h *Hand) AddCard(c deck.Card) {
	if h.cards == nil {
		h.cards = []deck.Card{}
	}
	h.cards = append(h.cards, c)
}
