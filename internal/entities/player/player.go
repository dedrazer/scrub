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

func (pb *PlayerBet) Lose(amount uint64) error {
	if pb.Player.Credits < amount {
		return errors.ErrInsufficientCredits
	}

	pb.Player.Credits -= amount
	pb.BetAmount = 0
	return nil
}
