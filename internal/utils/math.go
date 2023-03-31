package utils

import "math"

func Round(f float64, decimalPlaces int) float64 {
	shift := math.Pow(10, float64(decimalPlaces))
	return math.Round(f*shift) / shift
}
