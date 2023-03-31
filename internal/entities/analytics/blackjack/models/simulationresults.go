package models

type SimulationResults struct {
	AverageRoundsSurvived      uint    `json:"average_rounds_survived"`
	EarliestBankruptcyRound    uint    `json:"earliest_bankruptcy_round"`
	HighestProfitPercentage    float64 `json:"highest_profit_percentage"`
	OneCreditPercentageOfTotal float64 `json:"one_credit_percentage_of_total"`
	StartingCredits            uint64  `json:"starting_credits"`
	EndingCredits              uint64  `json:"ending_credits"`
	RebuyCredits               uint64  `json:"rebuy_credits"`
	BankAtCredits              uint64  `json:"bank_at"`
	Score                      float64 `json:"score"`
}
