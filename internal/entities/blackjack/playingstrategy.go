package blackjack

import bjutils "scrub/internal/entities/blackjack/utils"

func PlayingStrategy(playerHand Hand, dealerHand DealerHand, playerCredits uint64) string {
	if playerHand.CanSplit(playerCredits) {
		switch playerHand.cards[0].Symbol {
		case "A", "8":
			return bjutils.Split
		case "9":
			if !(dealerHand.UpCardValue() == 7 || dealerHand.UpCardValue() == 10 || dealerHand.UpCardValue() == 1) {
				return bjutils.Split
			}
		case "7":
			if dealerHand.UpCardValue() < 8 {
				return bjutils.Split
			}
		case "6":
			if dealerHand.UpCardValue() < 7 && dealerHand.UpCardValue() > 2 {
				return bjutils.Split
			}
		case "2", "3":
			if dealerHand.UpCardValue() < 8 && dealerHand.UpCardValue() > 3 {
				return bjutils.Split
			}
		}
	}

	if playerHand.IsSoft() {
		switch playerHand.UpperValue() {
		case 20, 19:
			return bjutils.Stand
		case 18:
			if dealerHand.UpCardValue() < 9 {
				return bjutils.Stand
			}
			return bjutils.Hit
		case 17:
			if dealerHand.UpCardValue() < 3 || dealerHand.UpCardValue() > 6 {
				return bjutils.Hit
			}
			return DoubleIfPossible(playerHand, playerCredits)
		case 16, 15:
			if dealerHand.UpCardValue() < 4 || dealerHand.UpCardValue() > 6 {
				return bjutils.Hit
			}
			return DoubleIfPossible(playerHand, playerCredits)
		case 14, 13:
			if dealerHand.UpCardValue() < 5 || dealerHand.UpCardValue() > 6 {
				return DoubleIfPossible(playerHand, playerCredits)
			}
		}
	}

	if playerHand.UpperValue() > 16 {
		return bjutils.Stand
	}

	switch playerHand.UpperValue() {
	case 16, 15, 14, 13:
		if dealerHand.UpCardValue() > 6 {
			return bjutils.Hit
		}
		return bjutils.Stand
	case 12:
		if dealerHand.UpCardValue() < 4 || dealerHand.UpCardValue() > 6 {
			return bjutils.Hit
		}
		return bjutils.Stand
	case 11:
		return DoubleIfPossible(playerHand, playerCredits)
	case 10:
		return DoubleIfPossible(playerHand, playerCredits)
	case 9:
		if dealerHand.UpCardValue() < 3 || dealerHand.UpCardValue() > 6 {
			return bjutils.Hit
		}
		return DoubleIfPossible(playerHand, playerCredits)
	default:
		return bjutils.Hit
	}
}

func DoubleIfPossible(h Hand, playerCredits uint64) string {
	if h.CanDouble(playerCredits) {
		return bjutils.Double
	}
	return bjutils.Hit
}
