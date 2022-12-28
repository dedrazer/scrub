package main

import (
	"fmt"
	"scrub/internal/entities/deck"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("failed to init zap: %w ", err))
	}

	deck.Demo(logger)
}
