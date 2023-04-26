package utils

const (
	First = "first"

	Stand  = "stand"
	Split  = "split"
	Hit    = "hit"
	Double = "double"
)

var (
	Win       = "win"
	Loss      = "loss"
	Push      = "push"
	Blackjack = "blackjack"
	SplitWon0 = "split won 0"
	SplitWon1 = "split won 1"
	SplitWon2 = "split won 2"
	SplitWon  = map[int]string{
		0: SplitWon0,
		1: SplitWon1,
		2: SplitWon2,
	}
	Bankrupt = "bankrupt"
)
