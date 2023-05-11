package utils_math

import "math"

func Sum[T int](arr ...T) T {
	var r T = 0
	for _, i := range arr {
		r += i
	}
	return r
}

func Abs(n float64) float64 {
	return math.Abs(n)
}
