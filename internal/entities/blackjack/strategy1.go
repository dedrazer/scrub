package blackjack

func Strategy1(playerHand Hand, dealerHand DealerHand) string {
	if playerHand.CanSplit() {
		switch playerHand.cards[0].Symbol {
		case "A", "8":
			return split
		case "9":
			if !(dealerHand.UpCardValue() == 7 || dealerHand.UpCardValue() == 10 || dealerHand.UpCardValue() == 1) {
				return split
			}
		case "7":
			if dealerHand.UpCardValue() < 8 {
				return split
			}
		case "6":
			if dealerHand.UpCardValue() < 7 && dealerHand.UpCardValue() > 2 {
				return split
			}
		case "2", "3":
			if dealerHand.UpCardValue() < 8 && dealerHand.UpCardValue() > 3 {
				return split
			}
		}
	}

	if playerHand.IsSoft() {
		switch playerHand.UpperValue() {
		case 20, 19:
			return stand
		case 18:
			if dealerHand.UpCardValue() < 9 {
				return stand
			}
			return hit
		case 17:
			if dealerHand.UpCardValue() < 3 || dealerHand.UpCardValue() > 6 {
				return hit
			}
			return DoubleIfPossible(playerHand)
		case 16, 15:
			if dealerHand.UpCardValue() < 4 || dealerHand.UpCardValue() > 6 {
				return hit
			}
			return DoubleIfPossible(playerHand)
		case 14, 13:
			if dealerHand.UpCardValue() < 5 || dealerHand.UpCardValue() > 6 {
				return DoubleIfPossible(playerHand)
			}
		}
	}

	if playerHand.UpperValue() > 16 {
		return stand
	}

	switch playerHand.UpperValue() {
	case 16, 15, 14, 13:
		if dealerHand.UpCardValue() > 6 {
			return hit
		}
		return stand
	case 12:
		if dealerHand.UpCardValue() < 4 || dealerHand.UpCardValue() > 6 {
			return hit
		}
		return stand
	case 11:
		return DoubleIfPossible(playerHand)
	case 10:
		return DoubleIfPossible(playerHand)
	case 9:
		if dealerHand.UpCardValue() < 3 || dealerHand.UpCardValue() > 6 {
			return hit
		}
		return DoubleIfPossible(playerHand)
	default:
		return hit
	}
}

func DoubleIfPossible(h Hand) string {
	if h.CanDouble() {
		return double
	}
	return hit
}
