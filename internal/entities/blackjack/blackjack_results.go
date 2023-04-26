package blackjack

import "scrub/internal/entities/blackjack/utils"

func (bj *Blackjack) ResultPlayerBust(hand *Hand) {
	bj.PlayerBust++
	bj.PlayerLosses++
	hand.Result = &utils.Loss
}

func (bj *Blackjack) ResultPlayerBlackjack(hand *Hand) {
	bj.PlayerBlackjackCount++
	bj.PlayerWins++
	hand.Result = &utils.Blackjack
}

func (bj *Blackjack) ResultDealerBust(hand *Hand) {
	bj.DealerBust++
	bj.PlayerWins++
	hand.Result = &utils.Win
}

func (bj *Blackjack) ResultDealerWins(hand *Hand) {
	bj.PlayerLosses++
	hand.Result = &utils.Loss
}

func (bj *Blackjack) ResultPlayerWins(hand *Hand) {
	bj.PlayerWins++
	hand.Result = &utils.Win
}

func (bj *Blackjack) ResultPush(hand *Hand) {
	bj.Pushes++
	hand.Result = &utils.Push
}
