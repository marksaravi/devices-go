package utils

import "math"

func ToRad(degree float64) float64 {
	return math.Pi / 180 * degree
}

func ToDeg(rad float64) float64 {
	return rad / math.Pi * 180
}
