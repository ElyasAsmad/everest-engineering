package utils

import "math"

// copied from https://stackoverflow.com/a/29786394
func Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(Round(num*output)) / output
}

func Truncate(val float64, precision int) float64 {
	pow := math.Pow(10, float64(precision))
	return math.Trunc(val*pow) / pow
}
