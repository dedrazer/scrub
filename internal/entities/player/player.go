package player

import "scrub/internal/errors"

type Player struct {
	Name    string
	Credits uint64
}

type PlayerBet struct {
	Player    Player
	BetAmount uint64
}

func (pb *PlayerBet) Win(amount uint64) {
	pb.Player.Credits += amount
	pb.BetAmount = 0
}

func (pb *PlayerBet) Lose() error {
	if pb.Player.Credits < pb.BetAmount {
		return errors.ErrInsufficientCredits
	}

	pb.Player.Credits -= pb.BetAmount
	pb.BetAmount = 0
	return nil
}
